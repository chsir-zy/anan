package redis

import (
	"sync"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/go-redis/redis/v8"
)

type AnanRedis struct {
	container framework.Container
	clients   map[string]*redis.Client

	lock *sync.RWMutex
}

func NewAnanRedis(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	clients := make(map[string]*redis.Client)
	lock := &sync.RWMutex{}

	return &AnanRedis{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil
}

func (ar *AnanRedis) GetClient(options ...contract.RedisOption) (*redis.Client, error) {
	config := GetBaseConfig(ar.container)

	for _, opt := range options {
		err := opt(ar.container, config)
		if err != nil {
			return nil, err
		}
	}

	// 查看key是否已经存在
	key := config.Uniqkey()

	ar.lock.RLock()
	if client, ok := ar.clients[key]; ok {
		ar.lock.RUnlock()
		return client, nil
	}
	ar.lock.RUnlock()

	ar.lock.Lock()
	defer ar.lock.Unlock()

	client := redis.NewClient(config.Options)
	ar.clients[key] = client

	return client, nil
}
