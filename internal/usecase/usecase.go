package usecase

import (
	"cachingService/internal/repository/cache"
	"context"
	"time"
)

type UseCase struct {
	cache cache.ILRUCache
}

func New(cache cache.ILRUCache) *UseCase {
	return &UseCase{cache: cache}
}

func (uc *UseCase) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return uc.cache.Put(ctx, key, value, ttl)
}

func (uc *UseCase) Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error) {
	return uc.cache.Get(ctx, key)
}

func (uc *UseCase) GetAll(ctx context.Context) (keys []string, values []interface{}, err error) {
	return uc.cache.GetAll(ctx)
}

func (uc *UseCase) Evict(ctx context.Context, key string) (value interface{}, err error) {
	return uc.cache.Evict(ctx, key)
} 

func (uc *UseCase) EvictAll(ctx context.Context) error {
	return uc.cache.EvictAll(ctx)
}