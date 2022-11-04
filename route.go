package main

import (
	"github.com/chsir-zy/anan/framework"
)

func registerRouter(core *framework.Core) {
	// core.Get("foo", framework.TimeoutHandler(FooControllerHandler, time.Second*1))
	// core.Get("foo", FooControllerHandler)

	core.Get("/user/login", UserControllerHandler)

	groupApi := core.Group("/subject")
	{
		groupApi.Get("/name", SubjectGetControllerHandler)
		groupApi.Put("/name", SubjectPutControllerHandler)
		groupApi.Post("/:id", SubjectPostControllerHandler)

		groupUser := groupApi.Group("/user")
		{
			groupUser.Get("/subname", SubjectSubGetControllerHandler)
		}

		groupInfo := groupApi.Group("/info")
		{
			groupInfo.Get("/ginfo", SubjectSubInfoGetControllerHandler)

			subGroupInfo := groupInfo.Group("/sun")
			subGroupInfo.Get("/subSun", SubjectSubInfoSunGetControllerHandler)
		}
	}

}
