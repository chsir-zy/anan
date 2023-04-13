// Package http API.
// @title hade
// @version 1.1
// @description.markdown api.md
// @termsOfService https://github.com/swaggo/swag

// @contact.name chisr
// @contact.email chisr

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @x-extension-openapi {"example": "value on a json format"}

package http

import (
	_ "github.com/chsir-zy/anan/app/http/swagger"
)
