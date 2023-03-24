package distributed

import (
	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
)

type LocalDistributedProvider struct {
}

func (p *LocalDistributedProvider) Register(container framework.Container) framework.NewInstance {
	return NewLocalDistributedService
}

func (p *LocalDistributedProvider) Boot(container framework.Container) error {
	return nil
}

func (p *LocalDistributedProvider) IsDefer() bool {
	return false
}

func (p *LocalDistributedProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (p *LocalDistributedProvider) Name() string {
	return contract.DistributedKey
}
