package kernel

import (
	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/chsir-zy/anan/framework/gin"
)

type AnanKernelProvider struct {
	HttpEngine *gin.Engine
}

func (provider *AnanKernelProvider) Register(c framework.Container) framework.NewInstance {
	return NewAnanKernelService
}

func (provider *AnanKernelProvider) Boot(c framework.Container) error {
	if provider.HttpEngine == nil {
		provider.HttpEngine = gin.Default()
	}

	provider.HttpEngine.SetContainer(c)
	return nil
}

func (provider *AnanKernelProvider) IsDefer() bool {
	return false
}

func (provider *AnanKernelProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.HttpEngine}
}

func (provider *AnanKernelProvider) Name() string {
	return contract.KernelKey
}
