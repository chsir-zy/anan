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

	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.AnanKernelProvider{HttpEngine: engine})
	}

	console.RunCommand(container)

	/* var a interface{}
	var f = make(map[int]interface{}, 0)
	f[123] = "123"

	a = f
	switch a.(type) {
	case map[interface{}]interface{}:
		fmt.Println(123)
	default:
		fmt.Println(456)
	} */
}
