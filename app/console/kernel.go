package console

import (
	"time"

	"github.com/chsir-zy/anan/app/console/command/demo"
	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/cobra"
	"github.com/chsir-zy/anan/framework/command"
)

func RunCommand(container framework.Container) error {
	var rootCmd = &cobra.Command{
		Use:   "anan",
		Short: "anan 命令",
		Long:  "anan框架提供的命令行工具",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	rootCmd.SetContainer(container)

	// 框架 启动web服务
	command.AddKernelCommands(rootCmd)

	// 业务相关
	AddAppCommand(rootCmd)

	return rootCmd.Execute()
}

// 绑定业务命令
func AddAppCommand(rootCmd *cobra.Command) {
	// rootCmd.AddCommand(demo.InitFoo())
	// rootCmd.AddCommand(demo.FooCommand)
	// rootCmd.AddCronCommand("* * * * * *", demo.FooCommand)

	rootCmd.AddDistributedCronCommand("foo_func_for_test", "*/5 * * * * *", demo.FooCommand, 2*time.Second)
}
