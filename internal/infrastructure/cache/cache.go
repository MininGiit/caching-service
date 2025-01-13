/*
Пакет cache имплементирует "cachingService/internal/repository/cache", содержит реализацию кеша с использованием алгоритма LRU
*/
package cache

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrKeyNotFound  = errors.New("key not found")
	ErrCacheEmpty   = errors.New("cache is empty")
	ErrDataNotValid = errors.New("data is not valid")
)

type item struct {
	key            string
	value          interface{}
	expiresAt      time.Time
	keyInTimeQueue *list.Element
}

// LRUCahe представляет структуру, реализующую LRU кеш
type LRUCache struct {
	size             int
	maxSize          int
	defaultTtl       time.Duration
	timeQueue        *list.List
	mutex            sync.Mutex
	data             map[string]*item
	cleaning         atomic.Bool
	intervalCleaning time.Duration
}

// New() создаёт новый экземпляр кеша,
func New(maxSize int, defaultTtl time.Duration) *LRUCache {
	data := make(map[string]*item, maxSize)
	timeQueue := list.New()
	return &LRUCache{
		size:             0,
		maxSize:          maxSize,
		defaultTtl:       defaultTtl,
		timeQueue:        timeQueue,
		data:             data,
		intervalCleaning: time.Second * 5,
	}
}

func validator(value interface{}) error {
	switch value.(type) {
	case int:
		return nil
	case float64:
		return nil
	case string:
		return nil
	default:
		return ErrDataNotValid
	}
}

// Put запись данных в кэш
func (c *LRUCache) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	err := validator(value)
	if err != nil {
		return err
	}
	var element *list.Element
	if item, ok := c.data[key]; ok {
		element = item.keyInTimeQueue
		c.timeQueue.MoveToBack(element)
		element = c.timeQueue.Back()
	} else {
		element = c.timeQueue.PushBack(key)
		if c.size == c.maxSize {
			front := c.timeQueue.Front()
			key := front.Value.(string)
			c.timeQueue.Remove(front)
			delete(c.data, key)
			c.size--
		}
		c.size++
	}
	if ttl == 0 {
		ttl = c.defaultTtl
	}
	expiresAt := time.Now().Add(ttl)
	newItem := &item{
		key:            key,
		value:          value,
		expiresAt:      expiresAt,
		keyInTimeQueue: element,
	}
	c.data[key] = newItem
	return nil
}

// Get получение данных из кэша по ключу
func (c *LRUCache) Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, ok := c.data[key]
	if !ok {
		return nil, time.Now(), ErrKeyNotFound
	}
	if item.expiresAt.Before(time.Now()) {
		element := item.keyInTimeQueue
		c.timeQueue.Remove(element)
		c.size--
		delete(c.data, key)
		return nil, time.Now(), ErrKeyNotFound
	}
	return item.value, item.expiresAt, nil
}

// GetAll получение всего наполнения кэша в виде двух слайсов: слайса ключей и слайса значений. Пары ключ-значения из кэша распол
func (c *LRUCache) GetAll(ctx context.Context) (keys []string, values []interface{}, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.size == 0 {
		return nil, nil, ErrCacheEmpty
	}
	keys = make([]string, 0, c.size)
	values = make([]interface{}, 0, c.size)
	for key, item := range c.data {
		if item.expiresAt.Before(time.Now()) {
			element := item.keyInTimeQueue
			c.timeQueue.Remove(element)
			c.size--
			delete(c.data, key)
			continue
		}
		keys = append(keys, key)
		values = append(values, item.value)
	}
	return keys, values, nil
}

// Evict ручное удаление данных по ключу
func (c *LRUCache) Evict(ctx context.Context, key string) (value interface{}, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, ok := c.data[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	element := item.keyInTimeQueue
	c.timeQueue.Remove(element)
	c.size--
	delete(c.data, key)
	return item.value, nil
}

// EvictAll ручная инвалидация всего кэша
func (c *LRUCache) EvictAll(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.size = 0
	c.data = make(map[string]*item)
	c.timeQueue = list.New()
	return nil
}

func (c *LRUCache) StartCollector() {
	c.cleaning.Store(true)
	go c.rottenDataCollector()
}

func (c *LRUCache) StopCollector() {
	c.cleaning.Store(false)
}

func (c *LRUCache) roundCleaning() {
	for key, item := range c.data {
		if !c.cleaning.Load() {
			return
		}
		if item.expiresAt.Before(time.Now()) {
			element := item.keyInTimeQueue
			c.timeQueue.Remove(element)
			c.size--
			delete(c.data, key)
		}
	}
}

func (c *LRUCache) rottenDataCollector() {
	ticker := time.NewTicker(c.intervalCleaning)
	defer ticker.Stop()
	for {
		if !c.cleaning.Load() {
			return
		}
		select {
		case <-ticker.C:
			c.roundCleaning()
		default:
			continue
		}
	}
}
