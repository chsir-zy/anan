package contract

const (
	// 生产环境
	EnvProduction = "production"
	// 测试环境
	EnvTesting = "testing"
	// 开发环境
	EnvDevelopment = "development"

	// 环境变量服务的字符串凭证
	EnvKey = "anna:key"
)

type Env interface {
	AppEnv() string         // 获取当前的环境 development/testing/production
	IsExist(string) bool    //判断一个变量是否存在
	Get(string) string      // 获取一个环境变量
	All() map[string]string // 获取所有环境变量 .env文件中的 和运行环境融合后的
}
