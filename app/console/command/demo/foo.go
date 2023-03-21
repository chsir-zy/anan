package demo

import (
	"fmt"
	"log"

	"github.com/chsir-zy/anan/framework/cobra"
	"github.com/chsir-zy/anan/framework/contract"
)

func InitFoo() *cobra.Command {
	// FooCommand.AddCommand(Foo1Command)
	FooCommand.AddCommand(envCommand)
	return FooCommand
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
