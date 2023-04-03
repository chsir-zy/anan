package app

import (
	"errors"
	"path/filepath"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/util"
	"github.com/google/uuid"
)

type AnanApp struct {
	container  framework.Container
	baseFolder string

	appId string // 表示当前这个app的唯一id  可以用于分布式锁

	configMap map[string]string // 配置加载
}

func (a AnanApp) Version() string {
	return "0.0.1"
}

func (a AnanApp) BaseFolder() string {
	if a.baseFolder != "" {
		return a.baseFolder
	}

	return util.GetExecDirectory()
}

func (a AnanApp) ConfigFolder() string {
	if val, ok := a.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "config")
}

func (a AnanApp) LogFolder() string {
	if val, ok := a.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(a.StorageFolder(), "log")
}

func (a AnanApp) HttpFolder() string {
	if val, ok := a.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "http")
}

func (a AnanApp) ConsoleFolder() string {
	if val, ok := a.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "console")
}

func (a AnanApp) StorageFolder() string {
	if val, ok := a.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "storage")
}

func (a AnanApp) ProviderFolder() string {
	if val, ok := a.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "app", "provider")
}

func (a AnanApp) MiddlewareFolder() string {
	if val, ok := a.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(a.HttpFolder(), "middleware")
}

func (a AnanApp) CommandFolder() string {
	if val, ok := a.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(a.ConsoleFolder(), "command")
}

func (a AnanApp) RuntimeFolder() string {
	if val, ok := a.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(a.StorageFolder(), "runtime")
}

func (a AnanApp) TestFolder() string {
	if val, ok := a.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "test")
}

func NewAnanApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}

	container := params[0].(framework.Container)
	baseFolder := params[1].(string)

	appId := uuid.New().String()
	return &AnanApp{baseFolder: baseFolder, container: container, appId: appId}, nil
}

func (a AnanApp) AppID() string {
	return a.appId
}

// LoadAppConfig 加载配置map
func (app *AnanApp) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		app.configMap[key] = val
	}
}
