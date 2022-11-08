package main

import (
	"net/http"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/middleware"
)

func main() {
	core := framework.NewCore()
	core.Use(middleware.Recovery(), middleware.Cost())
	registerRouter(core)

	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	server.ListenAndServe()

	// a := []int{1, 2, 3}
	// b := append(a, 4)
	// fmt.Println(a, b)
}
