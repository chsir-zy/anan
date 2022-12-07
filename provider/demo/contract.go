package demo

// 注入到服务容器的凭证
const Key = "anan:demo"

type Service interface {
	GetFoo() Foo
}

type Foo struct {
	Name string
}
