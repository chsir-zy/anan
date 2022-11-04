package framework

import "fmt"

// group方法返回的接口
type IGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Delete(string, ControllerHandler)
	Put(string, ControllerHandler)

	// 实现嵌套group
	Group(string) IGroup
}

type Group struct {
	core   *Core
	prefix string
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}

func (g *Group) Get(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.core.Get(uri, handler)
}

func (g *Group) Post(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.core.Post(uri, handler)
}

func (g *Group) Put(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.core.Put(uri, handler)
}

func (g *Group) Delete(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.core.Delete(uri, handler)
}

func (g *Group) Group(uri string) IGroup {
	ng := NewGroup(g.core, g.prefix+uri)
	fmt.Println(ng.prefix)
	return ng
}
