package user

import (
	"github.com/chsir-zy/anan/framework"
)

type UserProvider struct {
	framework.ServiceProvider
	c framework.Container
}

func (sp *UserProvider) Name() string {
	return UserKey
}

func (sp *UserProvider) Register(c framework.Container) framework.NewInstance {
	return NewUserService
}

func (sp *UserProvider) IsDefer() bool {
	return false
}

func (sp *UserProvider) Params(c framework.Container) []interface{} {
	return []interface{}{sp.c}
}

func (sp *UserProvider) Boot(c framework.Container) error {
	sp.c = c
	return nil
}
