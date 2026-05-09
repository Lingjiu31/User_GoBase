package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	ErrCodeSendTooMany       = errors.New("发送过于频繁")
	ErrCodeVerifyTooManyTime = errors.New("错误次数过多")
	ErrUnknowForCode         = errors.New("未知错误")
)

// 编译器会在编译时把代码放到变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{
		client: client,
	}
}

func (c *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		// 没有问题
		return nil
	case -1:
		// 发送频繁
		return ErrCodeSendTooMany
	default:
		// 系统错误
		return errors.New("系统错误")
	}
}

func (c *CodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		// 没问题
		return true, nil
	case -1:
		// 一直出错
		return false, ErrCodeVerifyTooManyTime
	case -2:
		// 出错一次
		return false, nil
	default:
		return false, ErrUnknowForCode
	}
}

func (c *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone:code:%s:%s", biz, phone)
}
