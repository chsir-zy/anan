package main

import (
	"net/http"

	"github.com/chsir-zy/anan/framework/gin"
	"github.com/chsir-zy/anan/provider/demo"
)

func main() {
	core := gin.New()
	core.Bind(&demo.DemoServiceProvider{})

	// core.Use(middleware.Recovery(), middleware.Cost())
	registerRouter(core)

	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	server.ListenAndServe()

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
