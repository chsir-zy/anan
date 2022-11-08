package middleware

import (
	"github.com/chsir-zy/anan/framework"
)

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.Json(500, "panic innser error")
			}
		}()
		c.Next()

		return nil
	}
}
