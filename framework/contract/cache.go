package contract

import (
	"context"
	"time"

	"github.com/chsir-zy/anan/framework"
)

const CacheKey = "anan:cache"

// Cache-aside模式对应的对象生成方法
type RememberFunc func(ctx context.Context, container framework.Container) (interface{}, error)

// 缓存服务
type CacheService interface {
	// get方法
	Get(ctx context.Context, key string) (string, error)
	GetObj(ctx context.Context, key string, model interface{}) error
	GetMany(ctx context.Context, keys []string) (map[string]string, error)

	// set 方法
	Set(ctx context.Context, key string, val string, timeout time.Duration) error
	SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error
	SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error
	SetForever(ctx context.Context, key string, val string) error
	SetForeverObj(ctx context.Context, key string, val interface{}) error

	// 设置某个key的超时时间
	SetTTL(ctx context.Context, key string, timeout time.Duration) error
	GetTTL(ctx context.Context, key string) (time.Duration, error)

	Remember(ctx context.Context, key string, timeout time.Duration, rememberFunc RememberFunc, model interface{}) error

	// 往key对应的值增加step
	Calc(ctx context.Context, key string, step int64) (int64, error)
	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)

	Del(ctx context.Context, key string) error
	DelMany(ctx context.Context, keys []string) error
}
