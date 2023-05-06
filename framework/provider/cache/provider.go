package cache

import (
	"strings"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/chsir-zy/anan/framework/provider/cache/services"
)

type AnanCacheProvider struct {
	framework.ServiceProvider

	Driver string
}

func (l *AnanCacheProvider) Register(c framework.Container) framework.NewInstance {
	if l.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			return services.NewRedisCache
		}

		cs := tcs.(contract.Config)
		l.Driver = strings.ToLower(cs.GetString("cacht.driver"))
	}

	return services.NewRedisCache
}

// Boot 启动的时候注入
func (l *AnanCacheProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer 是否延迟加载
func (l *AnanCacheProvider) IsDefer() bool {
	return true
}

// Params 定义要传递给实例化方法的参数
func (l *AnanCacheProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

// Name 定义对应的服务字符串凭证
func (l *AnanCacheProvider) Name() string {
	return contract.CacheKey
}
