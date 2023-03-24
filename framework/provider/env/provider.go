package env

import (
	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
)

type AnanEnvProvider struct {
	Floder string
}

func (provider *AnanEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewAnanEnv
}

func (provider *AnanEnvProvider) Boot(c framework.Container) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	provider.Floder = app.BaseFolder()
	return nil
}

func (provider *AnanEnvProvider) IsDefer() bool {
	return false
}

func (provider *AnanEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.Floder}
}

func (provider *AnanEnvProvider) Name() string {
	return contract.EnvKey
}
