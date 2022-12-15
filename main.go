package main

import (
	"github.com/chsir-zy/anan/app/console"
	appHttp "github.com/chsir-zy/anan/app/http"
	"github.com/chsir-zy/anan/app/provider/demo"
	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/provider/app"
	"github.com/chsir-zy/anan/framework/provider/kernel"
)

func main() {
	// 初始化服务容器
	container := framework.NewAnanContainer()
	container.Bind(&app.AnanAppProvider{})
	container.Bind(&demo.DemoProvider{})

	if engine, err := appHttp.NewHttpEngine(); err == nil {
		container.Bind(&kernel.AnanKernelProvider{HttpEngine: engine})
	}

	console.RunCommand(container)

	// core := gin.New()
	// core.Bind(&app.AnanAppProvider{})
	// core.Bind(&demo.DemoProvider{})

	// core.Use(middleware.Recovery(), middleware.Cost())
	// appHttp.Routes(core)

	// server := &http.Server{
	// 	Handler: core,
	// 	Addr:    ":8888",
	// }
	// server.ListenAndServe()

	/* go func() {
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println(123)
		log.Fatal("all goroutine is over ", err)

	} */

}
