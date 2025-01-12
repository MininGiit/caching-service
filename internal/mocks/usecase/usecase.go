package usecase

import (
	"context"
	"time"
)

type MockUseCase struct {
	Key 		string
	Value 		interface{}
	ExpiresAt 	time.Time
	Err			error
	Keyes		[]string
	Values 		[]interface{}
}


func (m *MockUseCase) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return m.Err
}

func (m *MockUseCase) Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error) {
	return  m.Value, m.ExpiresAt, m.Err
} 

func (m *MockUseCase) GetAll(ctx context.Context) (keys []string, values []interface{}, err error) {
	return  m.Keyes, m.Values, m.Err
}

func (m *MockUseCase) Evict(ctx context.Context, key string) (value interface{}, err error) {
	return m.Value, m.Err
}

func (m *MockUseCase) EvictAll(ctx context.Context) error {
	return m.Err
}
