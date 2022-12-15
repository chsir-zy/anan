package command

import "github.com/chsir-zy/anan/framework/cobra"

func AddKernelCommands(root *cobra.Command) {
	root.AddCommand(DemoCommand)

	root.AddCommand(initAppCommand())
}
