package main

import (
	"github.com/chsir-zy/anan/framework"
)

func registerRouter(core *framework.Core) {
	// core.Get("foo", framework.TimeoutHandler(FooControllerHandler, time.Second*1))
	// core.Get("foo", FooControllerHandler)
	// core.Use(middleware.Test1())
	core.Get("/user", UserControllerHandler)
	core.Get("/info", SubjectSubGetControllerHandler)

	// groupApi := core.Group("/subject")
	// {
	// 	// groupApi.Use(middleware.Test2())
	// 	groupApi.Get("/name", middleware.Test2(), SubjectGetControllerHandler)
	// 	groupApi.Put("/name", SubjectPutControllerHandler)
	// 	groupApi.Post("/:id", SubjectPostControllerHandler)

	// 	// groupUser := groupApi.Group("/user")
	// 	// {
	// 	// 	groupUser.Get("/subname", SubjectSubGetControllerHandler)
	// 	// }

	// }

}
