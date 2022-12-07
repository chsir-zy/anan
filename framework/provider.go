package framework

type NewInstance func(...interface{}) (interface{}, error)

// 服务提供者接口
type ServiceProvider interface {
	//在服务容器中测试一个 实例化服务的方法
	//是否在注册的时候就实例化 要根据IsDefer方法来判断
	Register(Container) NewInstance

	// 在调用实例化服务之前调用
	Boot(Container) error

	//是否在注册时候实例化服务
	IsDefer() bool

	// 传递给NewInstance的参数
	Params(Container) []interface{}

	// 代表这个服务提供者的凭证
	Name() string
}
