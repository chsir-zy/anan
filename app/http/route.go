package http

import (
	"github.com/chsir-zy/anan/app/http/module/demo"
	"github.com/chsir-zy/anan/framework/gin"
)

func Routes(r *gin.Engine) {
	r.Static("/dist", "./dist/")
	demo.Register(r)
}
