package demo

import (
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
	r.GET("/demo/orm", api.DemoOrm)
	r.GET("/demo/demo2", api.Demo2)
	r.POST("/demo/demo_post", api.DemoPost)
	r.POST("/demo/redis", api.DemoRedis)
	return nil
}

func NewDemoApi() *DemoApi {
	service := NewService()
	return &DemoApi{
		service: service,
	}
}

// Demo godoc
// @Summary 获取所有用户
// @tag.description.markdown demo.md
// @Produce  json
// @Tags demo
// @Success 200
// @Router /demo/demo [get]
func (api *DemoApi) Demo(c *gin.Context) {
	// appService := c.MustMake(contract.AppKey).(contract.App)
	// baseFolder := appService.BaseFolder()

	// c.JSON(http.StatusOK, baseFolder)

	// 获取password
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	password := configService.GetString("database.mysql.password")
	// 打印出来
	c.JSON(200, password)

}

// Demo2  for godoc
// @Summary 获取所有学生
// @Description 获取所有学生,不进行分页
// @Produce  json
// @Tags demo
// @Success 200
// @Router /demo/demo2 [get]
func (api *DemoApi) Demo2(c *gin.Context) {
	// testService := c.MustMake(user.UserKey).(user.Service)
	// foo := testService.Foo()
	// fmt.Println(foo)
	// c.JSON(200, foo)

	demoProvider := c.MustMake(demo.DemoKey).(demo.IService)
	students := demoProvider.GetAllStudent()
	userDTO := StudentToUserDTOs(students)

	c.JSON(200, userDTO)
}

type DTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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
