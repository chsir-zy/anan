package app

import (
	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
)

// 提供App的具体实现方法
type AnanAppProvider struct {
	BaseFolder string
}

func (a *AnanAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewAnanApp
}

func (a *AnanAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, a.BaseFolder}
}

func (a *AnanAppProvider) Boot(container framework.Container) error {
	return nil
}

func (a *AnanAppProvider) IsDefer() bool {
	return false
}

func (a *AnanAppProvider) Name() string {
	return contract.AppKey
}
