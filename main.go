package main

import (
	"github.com/chsir-zy/anan/app/console"
	"github.com/chsir-zy/anan/app/http"
	"github.com/chsir-zy/anan/app/provider/demo"
	"github.com/chsir-zy/anan/app/provider/user"
	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/provider/app"
	"github.com/chsir-zy/anan/framework/provider/config"
	"github.com/chsir-zy/anan/framework/provider/distributed"
	"github.com/chsir-zy/anan/framework/provider/env"
	"github.com/chsir-zy/anan/framework/provider/kernel"
	"github.com/chsir-zy/anan/framework/provider/log"
	"github.com/chsir-zy/anan/framework/provider/orm"
)

func main() {
	// 初始化服务容器
	container := framework.NewAnanContainer()
	container.Bind(&app.AnanAppProvider{})
	container.Bind(&demo.DemoProvider{})
	container.Bind(&user.UserProvider{})
	container.Bind(&env.AnanEnvProvider{})

	container.Bind(&distributed.LocalDistributedProvider{})
	container.Bind(&config.AnanConfigProvider{})
	container.Bind(&orm.GormProvider{})
	container.Bind(&log.AnanLogServiceProvider{})

	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.AnanKernelProvider{HttpEngine: engine})
	}

	console.RunCommand(container)

}

// func NewClient() {
// 	var ctx = context.Background()
// 	redis.NewClient(&redis.Options{
// 		Addr: "localhost:6739",
// 	})

// }
