package main

import (
	"fmt"

	"github.com/chsir-zy/anan/framework/gin"
	"github.com/chsir-zy/anan/provider/demo"
)

func UserControllerHandler(c *gin.Context) {
	fmt.Println("UserControllerHandler1")
	// ff, _ := c.FormFile("file")
	// fmt.Println(ff.Filename)

	// time.Sleep(10 * time.Second)
	c.ISetOkStatus()
	c.IJson("ok,UserControllerHandler")
	c.IJsonp("abc")
}

func SubjectListController(c *gin.Context) {
	demoService := c.MustMake(demo.Key).(demo.Service)
	fmt.Println(demoService)
	foo := demoService.GetFoo()

	c.ISetOkStatus().IJson(foo)
}
