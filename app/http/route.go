package http

import (
	"github.com/chsir-zy/anan/app/http/module/demo"
	"github.com/chsir-zy/anan/framework/gin"
	"github.com/chsir-zy/anan/framework/middleware/static"
)

func Routes(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))
	demo.Register(r)
}
