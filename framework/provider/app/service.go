package app

import (
	"errors"
	"flag"
	"path/filepath"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/util"
)

type AnanApp struct {
	container  framework.Container
	baseFolder string
}

func (a AnanApp) Version() string {
	return "0.0.1"
}

func (a AnanApp) BaseFolder() string {
	if a.baseFolder != "" {
		return a.baseFolder
	}

	var baseFolder string
	flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数为当前默认路劲")
	flag.Parse()

	if baseFolder != "" {
		return baseFolder
	}

	return util.GetExecDirectory()
}

func (a AnanApp) ConfigFolder() string {
	return filepath.Join(a.BaseFolder(), "config")
}
func (a AnanApp) LogFolder() string {
	return filepath.Join(a.StorageFolder(), "log")
}
func (a AnanApp) HttpFolder() string {
	return filepath.Join(a.BaseFolder(), "http")
}
func (a AnanApp) ConsoleFolder() string {
	return filepath.Join(a.BaseFolder(), "console")
}
func (a AnanApp) StorageFolder() string {
	return filepath.Join(a.BaseFolder(), "storage")
}
func (a AnanApp) ProviderFolder() string {
	return filepath.Join(a.BaseFolder(), "provider")
}
func (a AnanApp) MiddlewareFolder() string {
	return filepath.Join(a.HttpFolder(), "middleware")
}
func (a AnanApp) CommandFolder() string {
	return filepath.Join(a.ConsoleFolder(), "command")
}
func (a AnanApp) RuntimeFolder() string {
	return filepath.Join(a.StorageFolder(), "runtime")
}
func (a AnanApp) TestFolder() string {
	return filepath.Join(a.BaseFolder(), "test")
}

func NewAnanApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}

	container := params[0].(framework.Container)
	baseFolder := params[1].(string)

	return &AnanApp{baseFolder: baseFolder, container: container}, nil
}
