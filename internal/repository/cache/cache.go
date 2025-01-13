package cache

import (
	"context"
	"time"
)

type ILRUCache interface {
	// Put запись данных в кэш
	Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	// Get получение данных из кэша по ключу
	Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error)
	// GetAll получение всего наполнения кэша в виде двух слайсов: слайса ключей и слайса значений. Пары ключ-значения из кэша распол
	GetAll(ctx context.Context) (keys []string, values []interface{}, err error)
	// Evict ручное удаление данных по ключу
	Evict(ctx context.Context, key string) (value interface{}, err error)
	// EvictAll ручная инвалидация всего кэша
	EvictAll(ctx context.Context) error
}
