package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/chsir-zy/anan/framework"
)

func Test1() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("这个是中间件Test1 1部")
		c.Next()
		fmt.Println("这个是中间件Test1 2部")
		return nil
	}
}

func Test2() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("这个是中间件Test2 1部")
		c.Next()
		fmt.Println("这个是中间件Test2 2部")
		return nil
	}
}

//统计程序执行时长
func Cost() framework.ControllerHandler {
	return func(c *framework.Context) error {
		start := time.Now()

		c.Next()

		cost := time.Since(start)
		log.Printf("uri:%v, cost:%v", c.GetRequest().RequestURI, cost)

		return nil
	}
}
