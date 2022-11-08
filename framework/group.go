package framework

// group方法返回的接口
type IGroup interface {
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)

	// 实现嵌套group
	Group(string) IGroup
	Use(...ControllerHandler)
}

type Group struct {
	core       *Core
	prefix     string
	parent     *Group
	middleware []ControllerHandler
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		prefix: prefix,
		parent: nil,
	}
}

func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}

	return g.parent.getAbsolutePrefix() + g.prefix
}

func (g *Group) getMiddlewares() []ControllerHandler {
	if g.parent == nil {
		return g.middleware
	}

	return append(g.parent.middleware, g.middleware...)
}

func (g *Group) Use(handlers ...ControllerHandler) {
	g.middleware = append(g.middleware, handlers...)
}

func (g *Group) Get(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Get(uri, allHandlers...)
}

func (g *Group) Post(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Post(uri, allHandlers...)
}

func (g *Group) Put(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Put(uri, allHandlers...)
}

func (g *Group) Delete(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Delete(uri, allHandlers...)
}

func (g *Group) Group(uri string) IGroup {
	ng := NewGroup(g.core, uri)
	ng.parent = g //设置上级group  在group多次调用的时候
	return ng
}
