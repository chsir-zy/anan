package command

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chsir-zy/anan/framework/cobra"
	"github.com/chsir-zy/anan/framework/contract"
)

func initAppCommand() *cobra.Command {
	appCommand.AddCommand(appStartCommand)
	return appCommand
}

var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个web服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()

		// kernel服务实例 获取服务引擎(gin.Engine)
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		core := kernelService.HttpEngine()

		server := &http.Server{
			Handler: core,
			Addr:    ":8888",
		}

		go func() {
			server.ListenAndServe()
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-quit

		timeoutctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(timeoutctx); err != nil {
			log.Fatal("server shutdown:", err)
		}

		return nil
	},
}
