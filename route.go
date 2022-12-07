package main

import (
	"github.com/chsir-zy/anan/framework/gin"
)

func registerRouter(core *gin.Engine) {
	// core.Get("foo", framework.TimeoutHandler(FooControllerHandler, time.Second*1))
	// core.Get("foo", FooControllerHandler)
	// core.Use(middleware.Test1())
	core.GET("/user", UserControllerHandler)
	// core.Get("/info", SubjectSubGetControllerHandler)
	core.GET("/list/all", SubjectListController)
	groupApi := core.Group("/subject")
	{
		// groupApi.Use(middleware.Test2())
		// groupApi.Get("/name", middleware.Test2(), SubjectGetControllerHandler)
		// groupApi.Put("/name", SubjectPutControllerHandler)
		// groupApi.Post("/:id", SubjectPostControllerHandler)

		groupApi.GET("/list/all1", SubjectListController)

	}

}
