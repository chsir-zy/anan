package command

import (
	"fmt"

	"github.com/chsir-zy/anan/framework/cobra"
	"github.com/chsir-zy/anan/framework/contract"
)

func initEnvCommand() *cobra.Command {
	// envCommand.AddCommand()
	return envCommand
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
