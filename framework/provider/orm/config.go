package orm

import (
	"context"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
)

// 从database.yaml文件中读取数据
func GetBaseConfig(c framework.Container) *contract.DBConfig {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	logService := c.MustMake(contract.LogKey).(contract.Log)

	baseConfig := &contract.DBConfig{}
	err := configService.Load("database", baseConfig)
	if err != nil {
		logService.Error(context.Background(), "parse database config error", nil)
		return nil
	}

	return baseConfig
}

// 加载配置文件地址
func WithConfigPath(configPath string) contract.DBoption {
	return func(container framework.Container, config *contract.DBConfig) error {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		if err := configService.Load(configPath, config); err != nil {
			return err
		}
		return nil
	}
}
