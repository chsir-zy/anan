package orm

import (
	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
)

// 提供App的具体实现方法
type GormProvider struct {
}

func (a *GormProvider) Register(container framework.Container) framework.NewInstance {
	return NewAnanGorm
}

func (a *GormProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (a *GormProvider) Boot(container framework.Container) error {
	return nil
}

func (a *GormProvider) IsDefer() bool {
	return true
}

func (a *GormProvider) Name() string {
	return contract.ORMKey
}
