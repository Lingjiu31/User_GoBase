package repository

import (
	"Project-WeBook/webook/internal/repository/cache"
	"context"
)

var (
	ErrCodeSendTooMany       = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTime = cache.ErrCodeVerifyTooManyTime
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (repo *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, code)
}
