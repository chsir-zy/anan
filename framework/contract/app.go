package contract

// 定义app字符串凭证
const AppKey = "anan:app"

type App interface {
	// 当前版本
	Version() string

	// 项目的根目录
	BaseFolder() string

	// 配置目录
	ConfigFolder() string

	// 日志目录
	LogFolder() string

	// 服务提供者目录
	ProviderFolder() string

	// 业务中间件目录
	MiddlewareFolder() string

	CommandFolder() string

	RuntimeFolder() string

	TestFolder() string
}
