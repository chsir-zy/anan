package demo

import (
	"net/http"

	"github.com/chsir-zy/anan/app/provider/demo"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/chsir-zy/anan/framework/gin"
)

type DemoApi struct {
	service *Service
}

func Register(r *gin.Engine) error {
	api := NewDemoApi()
	r.Bind(&demo.DemoProvider{})
	r.GET("/demo/demo", api.Demo)
	r.GET("/demo/demo2", api.Demo2)
	r.POST("/demo/demo_post", api.DemoPost)
	return nil
}

func NewDemoApi() *DemoApi {
	service := NewService()
	return &DemoApi{
		service: service,
	}
}

func (api *DemoApi) Demo(c *gin.Context) {
	appService := c.MustMake(contract.AppKey).(contract.App)
	baseFolder := appService.BaseFolder()

	c.JSON(http.StatusOK, baseFolder)
}

func (api *DemoApi) Demo2(c *gin.Context) {
	demoProvider := c.MustMake(demo.DemoKey).(demo.IService)
	students := demoProvider.GetAllStudent()
	userDTO := StudentToUserDTOs(students)

	c.JSON(200, userDTO)
}

// func (api *DemoApi) DemoPost(c *gin.Context) {
// 	type Foo struct {
// 		Name string
// 	}

// 	foo := &Foo{}
// 	err := c.BindJSON(&foo)
// 	fmt.Println(err)
// 	if err != nil {
// 		c.AbortWithError(500, err)
// 	}

// 	c.JSON(200, nil)
// }

func (api *DemoApi) DemoPost(c *gin.Context) {
	type Foo struct {
		Name string
	}
	foo := Foo{}
	err := c.BindJSON(&foo)
	if err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(200, foo)
}
