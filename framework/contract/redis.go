package contract

import (
	"fmt"

	"github.com/chsir-zy/anan/framework"
	"github.com/go-redis/redis/v8"
)

const RedisKey = "anan:redis"

type RedisConfig struct {
	*redis.Options
}

func (rc *RedisConfig) Uniqkey() string {
	return fmt.Sprintf("%v_%v_%v_%v", rc.Addr, rc.DB, rc.Username, rc.Network)
}

// 修改option选项的方法
type RedisOption func(container framework.Container, config *RedisConfig) error

// 一个redis服务
type RedisService interface {
	GetClient(option ...RedisOption) (*redis.Client, error)
}
