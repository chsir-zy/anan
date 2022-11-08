package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router     map[string]*Tree
	middleware []ControllerHandler
}

func NewCore() *Core {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()

	return &Core{router: router}
}

func (c *Core) Use(handlers ...ControllerHandler) {
	c.middleware = append(c.middleware, handlers...)
}

func (c *Core) Get(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middleware, handlers...)
	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add route err: ", err)
	}
}

func (c *Core) Post(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middleware, handlers...)
	if err := c.router["POST"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add route err: ", err)
	}
}

func (c *Core) Put(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middleware, handlers...)
	if err := c.router["PUT"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add route err: ", err)
	}
}

func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middleware, handlers...)
	if err := c.router["DELETE"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add route err: ", err)
	}
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (c *Core) FindRouteByRequest(request *http.Request) []ControllerHandler {
	method := request.Method
	uri := request.URL.Path
	upMethod := strings.ToUpper(method)
	upUri := strings.ToUpper(uri)

	if methodHandler, ok := c.router[upMethod]; ok {
		return methodHandler.FindControllerHandler(upUri)
	}

	return nil
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core.serveHTTP")
	ctx := NewContext(request, response)

	handlers := c.FindRouteByRequest(request)
	if handlers == nil {
		ctx.Json(404, "router not found")
		return
	}
	// err1 := SubjectGetControllerHandler(ctx)

	ctx.SetHandlers(handlers)
	if err := ctx.Next(); err != nil {
		ctx.Json(500, "service error")
		return
	}

}
