package usecase

import (
	"cachingService/internal/repository/cache"
	"context"
	"time"
)

type UseCase struct {
	repository cache.ILRUCache
}

func New(cache cache.ILRUCache) *UseCase {
	return &UseCase{repository: cache}
}

func (uc *UseCase) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return uc.repository.Put(ctx, key, value, ttl)
}

func (uc *UseCase) Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error) {
	return uc.repository.Get(ctx, key)
}

func (uc *UseCase) GetAll(ctx context.Context) (keys []string, values []interface{}, err error) {
	return uc.repository.GetAll(ctx)
}

func (uc *UseCase) Evict(ctx context.Context, key string) (value interface{}, err error) {
	return uc.repository.Evict(ctx, key)
}

func (uc *UseCase) EvictAll(ctx context.Context) error {
	return uc.repository.EvictAll(ctx)
}
