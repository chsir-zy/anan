package command

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/cobra"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/chsir-zy/anan/framework/util"
	"github.com/erikdubbelboer/gspt"
	"github.com/pkg/errors"
	"github.com/sevlyar/go-daemon"
)

const CLOST_WAIT = 5

var appAddress = ""   //app启动地址
var appDaemon = false //是否以daemond的方式启动

func initAppCommand() *cobra.Command {
	appStartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", false, "start app daemon")
	appStartCommand.Flags().StringVar(&appAddress, "address", "", "设置app启动的地址，默认是:8888")
	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appStateCommand)
	appCommand.AddCommand(appStopCommand)
	appCommand.AddCommand(appRestartCommand)
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

// 启动app(web)服务 会阻塞在当前的goroutine
func startAppServe(server *http.Server, c framework.Container) error {
	// 先启动一个goroutine启动服务
	go func() {
		server.ListenAndServe()
	}()

	// 监听中断信号 阻塞在这里
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	closeWait := CLOST_WAIT
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	if configService.IsExist("app.close_wait") {
		closeWait = configService.GetInt("app.close_wait")
	}
	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Duration(closeWait)*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutContext); err != nil {
		return err
	}
	return nil
}

var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个web服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()

		// kernel服务实例 获取服务引擎(gin.Engine)
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		core := kernelService.HttpEngine()

		if appAddress == "" {
			envService := container.MustMake(contract.EnvKey).(contract.Env)
			if envAddress := envService.Get("ADDRESS"); envAddress != "" {
				appAddress = envAddress
			} else {
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configAddress := configService.GetString("app.address"); configAddress != "" {
					appAddress = configAddress
				} else {
					appAddress = ":8888"
				}
			}
		}
		server := &http.Server{
			Handler: core,
			Addr:    appAddress,
		}

		appService := container.MustMake(contract.AppKey).(contract.App)
		pidFolder := appService.RuntimeFolder()
		if !util.Exists(pidFolder) {
			if err := os.MkdirAll(pidFolder, os.ModePerm); err != nil {
				return nil
			}
		}
		serverPidFile := filepath.Join(pidFolder, "app.pid")
		logFolder := appService.LogFolder()
		if !util.Exists(logFolder) {
			if err := os.MkdirAll(logFolder, os.ModePerm); err != nil {
				return nil
			}
		}
		serverLogFile := filepath.Join(logFolder, "app.log")
		currentFolder := util.GetExecDirectory()

		if appDaemon {
			cntxt := &daemon.Context{
				// 设置pid文件
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				// 设置日志文件
				LogFileName: serverLogFile,
				// 工作路径
				WorkDir: currentFolder,
				Umask:   027,
				Args:    []string{"", "app", "start", "--daemon=true"},
			}
			// 如果d不为空 就为父进程  为空则为子进程
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}

			if d != nil {
				// 父进程直接打印成功信息，
				fmt.Println("app启动成功，pid:", d.Pid)
				fmt.Println("日志文件：", serverLogFile)
				return nil
			}

			defer cntxt.Release()

			// 子进程真正执行app启动操作
			fmt.Println("daemon started")
			gspt.SetProcTitle("anan app")
			if err := startAppServe(server, container); err != nil {
				fmt.Println(err)
			}
			return nil
		}

		// 非daemon模式, 直接执行
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := ioutil.WriteFile(serverPidFile, []byte(content), 0664)
		if err != nil {
			return err
		}

		gspt.SetProcTitle("anan app")
		fmt.Println("app serve url: ", appAddress)
		if err := startAppServe(server, container); err != nil {
			fmt.Println(err)
		}

		return nil
	},
}

var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "查看app的运行状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		serverFolder := appService.RuntimeFolder()
		serverFolderFile := filepath.Join(serverFolder, "app.pid")
		strPid, err := ioutil.ReadFile(serverFolderFile)
		if err != nil {
			fmt.Println("app进程没有运行")
			return err
		}

		// 判断进程文件里面是否有数据
		if strPid == nil || len(strPid) == 0 {
			fmt.Println("app进程没有运行")
			return err
		}

		pid, err := strconv.Atoi(string(strPid))
		if err != nil {
			fmt.Println("app进程没有运行")
			return err
		}
		isExist := util.CheckProcessExist(pid)
		if !isExist {
			fmt.Println("app进程没有运行")
			return err
		}

		fmt.Println("app进程正在运行，进程号PID: ", pid)
		return nil
	},
}

var appStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止运行app进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		serverFolderFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
		strPid, err := ioutil.ReadFile(serverFolderFile)
		if err != nil {
			fmt.Println("app进程没有运行")
			return err
		}

		// 判断进程文件里面是否有数据
		if strPid == nil || len(strPid) == 0 {
			fmt.Println("app进程没有运行")
			return err
		}

		pid, err := strconv.Atoi(string(strPid))
		if err != nil {
			fmt.Println("app进程没有运行")
			return err
		}
		isExist := util.CheckProcessExist(pid)
		if !isExist {
			fmt.Println("app进程没有运行")
			return err
		}

		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			return nil
		}

		if err := ioutil.WriteFile(serverFolderFile, []byte{}, 0664); err != nil {
			return nil
		}

		fmt.Println("app进程停止运行，进程号PID: ", pid)
		return nil
	},
}

var appRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重新启动app进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		serverFolderFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
		strPid, err := ioutil.ReadFile(serverFolderFile)
		if err != nil {
			fmt.Println("app进程没有运行")
			return err
		}

		// 判断进程文件里面是否有数据
		if strPid != nil || len(strPid) != 0 {
			pid, err := strconv.Atoi(string(strPid))
			if err != nil {
				fmt.Println("app进程没有运行")
				return err
			}
			if util.CheckProcessExist(pid) { //检测到进程正在运行
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					return nil
				}

				// 等待两个closeWait的时间让原来的进程先介绍
				closeWait := CLOST_WAIT * 2
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configService.IsExist("app.close_wait") {
					closeWait = configService.GetInt("app.close_wait")
				}

				for i := 0; i < closeWait; i++ {
					if !util.CheckProcessExist(pid) { //如果进程结束了则跳出循环
						break
					}
					time.Sleep(time.Second * 1)
				}

				// 循环结束后进程还未结束 则返回
				if util.CheckProcessExist(pid) {
					fmt.Println("结束进程失败, Pid:", pid)
					return errors.New("结束进程失败")
				}

				if err := ioutil.WriteFile(serverFolderFile, []byte{}, 0664); err != nil {
					return nil
				}
				fmt.Println("结束进程成功，PID:", pid)
			}
		}

		//以daemon的方式启动新进程
		appDaemon = true
		return appStartCommand.RunE(cmd, args)
	},
}
