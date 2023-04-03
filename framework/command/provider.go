package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/cobra"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/chsir-zy/anan/framework/util"
)

/*
	自动生成服务相关的文件 contract.go  provider.go  service.go
*/

func InitProviderCommand() *cobra.Command {
	providerCommand.AddCommand(listProviderCommand)
	providerCommand.AddCommand(createProviderCommand)
	return providerCommand
}

var providerCommand = &cobra.Command{
	Use:   "provider",
	Short: "生成服务",
	Long:  "自动构建生成服务需要的基础文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

var listProviderCommand = &cobra.Command{
	Use:   "list",
	Short: "显式已有的服务列表",
	Long:  "显式已有的服务列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer().(*framework.AnanContainer)
		nameList := container.NameList()
		for _, name := range nameList {
			println(name)
		}
		return nil
	},
}

var createProviderCommand = &cobra.Command{
	Use:   "new",
	Short: "创建一个服务",
	Long:  "创建一个服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer().(*framework.AnanContainer)

		var name string   // 服务名
		var folder string // 目录

		{
			prompt := &survey.Input{
				Message: "请输入服务名称(服务凭证)：",
			}
			survey.AskOne(prompt, &name)
		}

		{
			prompt := &survey.Input{
				Message: "请输入服务目录名(默认服务名一致)：",
			}
			survey.AskOne(prompt, &folder)
		}

		nameList := container.NameList()
		// 检查服务是否已经存在了
		strMap := util.ConvertStrSlice2Map(nameList)
		if _, ok := strMap[name]; ok {
			println("服务名" + name + " 已经存在了,请换个")
			return nil
		}

		if folder == "" {
			folder = name
		}

		//校验 文件夹/目录 是否已经存在
		appSer := container.MustMake(contract.AppKey).(contract.App)
		pFolder := appSer.ProviderFolder() // provider所在的目录

		fs, err := os.ReadDir(pFolder)
		if err != nil {
			println("系统出错了")
			return nil
		}

		var checkFolder []string
		for _, f := range fs {
			if f.IsDir() {
				checkFolder = append(checkFolder, f.Name())
			}
		}
		checkFolderMap := util.ConvertStrSlice2Map(checkFolder)
		if _, ok := checkFolderMap[folder]; ok {
			println("文件夹" + folder + " 已经存在了,请换个")
			return nil
		}

		err = os.Mkdir(filepath.Join(pFolder, folder), 0777)
		if err != nil {
			println(err)
			return nil
		}

		funcs := template.FuncMap{"title": strings.Title}
		{
			file := filepath.Join(pFolder, folder, "contract.go")
			fs, err := os.Create(file)
			if err != nil {
				println(err)
				return nil
			}

			t := template.Must(template.New("contract").Funcs(funcs).Parse(ContractTmp))
			t.Execute(fs, name)
		}

		{
			file := filepath.Join(pFolder, folder, "provider.go")
			fs, err := os.Create(file)
			if err != nil {
				println(err)
				return nil
			}

			t := template.Must(template.New("provider").Funcs(funcs).Parse(ProviderTmp))
			t.Execute(fs, name)
		}

		{
			file := filepath.Join(pFolder, folder, "service.go")
			fs, err := os.Create(file)
			if err != nil {
				println(err)
				return nil
			}

			t := template.Must(template.New("service").Funcs(funcs).Parse(ServiceTmp))
			t.Execute(fs, name)
		}

		fmt.Println("创建服务成功, 文件夹地址:", filepath.Join(pFolder, folder))
		fmt.Println("请不要忘记挂载新创建的服务")
		return nil
	},
}

var ContractTmp = `
package {{.}}

// 定义app字符串凭证
const {{.|title}}Key = "{{.}}"

type Service interface {
	// 请在这里定义你的方法
	Foo() string
}
`

var ProviderTmp = `
package {{.}}

import (
	"github.com/chsir-zy/anan/framework"
)

type {{.|title}}Provider struct {
	framework.ServiceProvider
	c framework.Container
}

func (sp *{{.|title}}Provider) Name() string {
	return {{.|title}}Key
}

func (sp *{{.|title}}Provider) Register(c framework.Container) framework.NewInstance {
	return New{{.|title}}Service
}

func (sp *{{.|title}}Provider) IsDefer() bool {
	return false
}

func (sp *{{.|title}}Provider) Params(c framework.Container) []interface{} {
	return []interface{}{sp.c}
}

func (sp *{{.|title}}Provider) Boot(c framework.Container) error {
	sp.c = c
	return nil
}
`

var ServiceTmp = `
package {{.}}

import (
	"github.com/chsir-zy/anan/framework"
)

type {{.|title}}Service struct {
	container framework.Container
}

func New{{.|title}}Service(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &{{.|title}}Service{container: container}, nil
}

func (s *{{.|title}}Service) Foo() string {
    return ""
}
`
