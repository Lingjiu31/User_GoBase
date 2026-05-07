package cache

import (
	"Project-WeBook/webook/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ErrKeyNotExist = errors.New("key不存在")

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Hour,
	}
}

func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.Key(id)
	// 如果数据不存在,会返回未查找到 key
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, ErrKeyNotExist
	}
	var user domain.User
	err = json.Unmarshal(val, &user)
	return user, err
}

func (cache *UserCache) Set(ctx context.Context, user domain.User) error {
	u, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.Key(user.Id)

	return cache.client.Set(ctx, key, u, cache.expiration).Err()
}

func (cache *UserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
