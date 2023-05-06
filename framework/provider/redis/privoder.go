package redis

import (
	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
)

// 提供App的具体实现方法
type RedisProvider struct {
}

func (a *RedisProvider) Register(container framework.Container) framework.NewInstance {
	return NewAnanRedis
}

func (a *RedisProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (a *RedisProvider) Boot(container framework.Container) error {
	return nil
}

func (a *RedisProvider) IsDefer() bool {
	return true
}

func (a *RedisProvider) Name() string {
	return contract.RedisKey
}
