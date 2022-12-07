package middleware

import (
	"net/http"

	"github.com/chsir-zy/anan/framework/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.ISetStatus(http.StatusInternalServerError).IJson("panic innser error")
			}
		}()
		c.Next()

	}
}
