package usecase

import (
	"context"
	"errors"
	"time"
)

type MockCache struct {
	Size      int
	Key       string
	Value     interface{}
	ExpiresAt time.Time
	Err       error
	Keyes     []string
	Values    []interface{}
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
		return errors.New("the data is not valid")
	}
}

func (m *MockCache) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return validator(value)
}

func (m *MockCache) Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error) {
	if key != m.Key {
		return nil, time.Now(), errors.New("not found")
	}
	return m.Value, m.ExpiresAt, m.Err
}

func (m *MockCache) GetAll(ctx context.Context) (keys []string, values []interface{}, err error) {
	if m.Size == 0 {
		return nil, nil, errors.New("cache is empty")
	}
	return m.Keyes, m.Values, m.Err
}

func (m *MockCache) Evict(ctx context.Context, key string) (value interface{}, err error) {
	if key != m.Key {
		return nil, errors.New("not found")
	}
	return m.Value, m.Err
}

func (m *MockCache) EvictAll(ctx context.Context) error {
	return nil
}
