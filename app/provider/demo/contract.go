package demo

// 注入到服务容器的凭证
const DemoKey = "demo"

type IService interface {
	GetAllStudent() []Student
}

type Student struct {
	Id   int
	Name string
}
