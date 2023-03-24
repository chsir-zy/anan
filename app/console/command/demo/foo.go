package demo

import (
	"fmt"
	"log"
	"os"

	"github.com/chsir-zy/anan/framework/cobra"
	"github.com/chsir-zy/anan/framework/contract"
)

func InitFoo() *cobra.Command {
	// FooCommand.AddCommand(Foo1Command)
	// FooCommand.AddCommand(envCommand)
	// return FooCommand

	ConfigCommand.AddCommand(getConfigCommand)
	return ConfigCommand
}

var FooCommand = &cobra.Command{
	Use:     "foo",
	Short:   "foo简要说明",
	Long:    "foo的详细描述",
	Aliases: []string{"fo", "f"},
	Example: "foo的命令例子",
	RunE: func(cmd *cobra.Command, args []string) error {
		// cmd.Help()
		log.Println("execute foo command")
		return nil
	},
}

var Foo1Command = &cobra.Command{
	Use:     "foo1",
	Short:   "foo1简要说明",
	Long:    "foo1的详细描述",
	Aliases: []string{"fo1", "f1"},
	Example: "foo1的命令例子",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		log.Println(container)
		return nil
	},
}

// envCommand 获取当前的 App 环境
var envCommand = &cobra.Command{
	Use:   "env",
	Short: "获取当前的 App 环境",
	Run: func(c *cobra.Command, args []string) {
		// 获取 env 环境
		container := c.GetContainer()
		envService := container.MustMake(contract.EnvKey).(contract.Env)
		// 打印环境
		fmt.Println("environment:", envService.AppEnv())
	},
}

var ConfigCommand = &cobra.Command{
	Use:   "config",
	Short: "获取配置文件",
	Long:  "获取配置文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("execute config command")
		return nil
	},
}

var getConfigCommand = &cobra.Command{
	Use:   "get",
	Short: "获取某个配置文件",
	Long:  "获取某个配置文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := os.Args[3]
		c := cmd.GetContainer()
		log.Println(params)

		conf := c.MustMake(contract.ConfigKey).(contract.Config)
		p := conf.GetString(params)
		log.Println(p)
		return nil
	},
}
