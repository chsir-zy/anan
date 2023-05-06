package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/chsir-zy/anan/framework/provider/redis"
	redisv8 "github.com/go-redis/redis/v8"
)

type RedisCache struct {
	container framework.Container
	client    *redisv8.Client
	lock      *sync.RWMutex
}

func NewRedisCache(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	if !container.IsBind(contract.RedisKey) {
		err := container.Bind(&redis.RedisProvider{})
		if err != nil {
			return nil, err
		}
	}

	redisService := container.MustMake(contract.RedisKey).(contract.RedisService)
	client, err := redisService.GetClient(redis.WithConfigPath("cache"))
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		container: container,
		client:    client,
		lock:      &sync.RWMutex{},
	}, nil
}

// GetObj 获取某个key对应的对象, 对象必须实现 https://pkg.go.dev/encoding#BinaryUnMarshaler
func (r *RedisCache) GetObj(ctx context.Context, key string, model interface{}) error {
	cmd := r.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redisv8.Nil) {
		return ErrKeyNotFound
	}

	err := cmd.Scan(model)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	pipeline := r.client.Pipeline()
	vals := make(map[string]string)
	cmds := make([]*redisv8.StringCmd, 0, len(keys))

	for _, key := range keys {
		cmds = append(cmds, pipeline.Get(ctx, key))
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}

	for _, cmd := range cmds {
		val, _ := cmd.Result()
		key := cmd.Args()[1].(string)

		vals[key] = val
	}

	return vals, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, val string, timeout time.Duration) error {
	return r.client.Set(ctx, key, val, timeout).Err()
}

func (r *RedisCache) SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return r.client.Set(ctx, key, val, timeout).Err()
}

func (r *RedisCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	pipline := r.client.Pipeline()
	cmds := make([]*redisv8.StatusCmd, 0, len(data))
	for k, v := range data {
		cmds = append(cmds, pipline.Set(ctx, k, v, timeout))
	}

	_, err := pipline.Exec(ctx)
	return err
}

func (r *RedisCache) SetForever(ctx context.Context, key string, val string) error {
	return r.client.Set(ctx, key, val, NoneDuration).Err()
}

func (r *RedisCache) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	return r.client.Set(ctx, key, val, NoneDuration).Err()
}

func (r *RedisCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	return r.client.Expire(ctx, key, timeout).Err()
}

func (r *RedisCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *RedisCache) Remember(ctx context.Context, key string, timeout time.Duration, rememberFucn contract.RememberFunc, obj interface{}) error {
	err := r.GetObj(ctx, key, obj)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return err
	}

	objNew, err := rememberFucn(ctx, r.container)
	if err != nil {
		return err
	}

	if err := r.SetObj(ctx, key, objNew, timeout); err != nil {
		return err
	}

	if err := r.GetObj(ctx, key, obj); err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	return r.client.IncrBy(ctx, key, step).Result()
}

func (r *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.IncrBy(ctx, key, 1).Result()
}

func (r *RedisCache) Decrement(ctx context.Context, key string) (int64, error) {
	return r.client.IncrBy(ctx, key, -1).Result()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) DelMany(ctx context.Context, keys []string) error {
	pipeline := r.client.Pipeline()
	cmds := make([]*redisv8.IntCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipeline.Del(ctx, key))
	}
	_, err := pipeline.Exec(ctx)
	return err
}
