package demo

import (
	"log"

	"github.com/chsir-zy/anan/framework/cobra"
)

func InitFoo() *cobra.Command {
	FooCommand.AddCommand(Foo1Command)
	return FooCommand
}

var FooCommand = &cobra.Command{
	Use:     "foo",
	Short:   "foo简要说明",
	Long:    "foo的详细描述",
	Aliases: []string{"fo", "f"},
	Example: "foo的命令例子",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		container := cmd.GetContainer()
		log.Println(container)
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
