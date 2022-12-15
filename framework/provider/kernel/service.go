package kernel

import (
	"net/http"

	"github.com/chsir-zy/anan/framework/gin"
)

type AnanKernelService struct {
	engine *gin.Engine
}

func NewAnanKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &AnanKernelService{engine: httpEngine}, nil
}

func (s *AnanKernelService) HttpEngine() http.Handler {
	return s.engine
}
