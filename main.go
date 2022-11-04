package main

import (
	"net/http"

	"github.com/chsir-zy/anan/framework"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	server.ListenAndServe()

	// s := "a/b/c"
	// split := strings.SplitN(s, "/", 2)
	// fmt.Println(split)
}
