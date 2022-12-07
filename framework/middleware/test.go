package middleware

import (
	"log"
	"time"

	"github.com/chsir-zy/anan/framework/gin"
)

//统计程序执行时长
func Cost() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		cost := time.Since(start)
		log.Printf("uri:%v, cost:%v", c.Request.RequestURI, cost)
		time.Sleep(2 * time.Second)

		c.Next()

	}
}
