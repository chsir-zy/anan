package command

import (
	"fmt"
	"path/filepath"

	"github.com/chsir-zy/anan/framework/cobra"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/swaggo/swag/gen"
)

func initSwaggerCommand() *cobra.Command {
	swaggerCommand.AddCommand(swaggerGenCommand)
	return swaggerCommand
}

var swaggerCommand = &cobra.Command{
	Use:   "swagger",
	Short: "swagger对应命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var swaggerGenCommand = &cobra.Command{
	Use:   "gen",
	Short: "生成对应的swagger文件，包括swagger.yaml, docs.go",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		outputDir := filepath.Join(appService.AppFolder(), "http", "swagger")
		httpFolder := filepath.Join(appService.AppFolder(), "http")

		conf := &gen.Config{
			SearchDir:          httpFolder, //表示在哪个目录下找注释代码
			Excludes:           "",
			OutputDir:          outputDir,
			MainAPIFile:        "swagger.go",
			PropNamingStrategy: "",
			ParseVendor:        false,
			ParseDependency:    false,
			ParseInternal:      false,
			MarkdownFilesDir:   httpFolder,
			GeneratedTime:      false,
		}

		err := gen.New().Build(conf)
		if err != nil {
			fmt.Println(err)
		}
	},
}
