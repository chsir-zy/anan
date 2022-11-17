package middleware

import (
	"net/http"

	"github.com/chsir-zy/anan/framework"
)

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.SetStatus(http.StatusInternalServerError).Json("panic innser error")
			}
		}()
		c.Next()

		return nil
	}
}
