package demo

import (
	"time"

	"github.com/chsir-zy/anan/framework/contract"
	"github.com/chsir-zy/anan/framework/gin"
	"github.com/chsir-zy/anan/framework/provider/redis"
)

func (api *DemoApi) DemoRedis(c *gin.Context) {
	logger := c.MustMakeLog()
	logger.Info(c, "request start", nil)

	// 初始化redis服务
	redisService := c.MustMake(contract.RedisKey).(contract.RedisService)
	client, err := redisService.GetClient(redis.WithConfigPath("cache"), redis.WithRedisConfig(func(options *contract.RedisConfig) {
		options.MaxRetries = 3
	}))

	if err != nil {
		logger.Error(c, err.Error(), nil)
		c.AbortWithError(50001, err)
		return
	}

	if err := client.Set(c, "foo", "bar", 1*time.Hour).Err(); err != nil {
		c.AbortWithError(500, err)
	}

	val := client.Get(c, "foo").String()
	logger.Info(c, "redis get", map[string]interface{}{
		"val": val,
	})

	c.JSON(200, "ok")
}
