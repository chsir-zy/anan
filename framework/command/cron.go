package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/chsir-zy/anan/framework/cobra"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/chsir-zy/anan/framework/util"
	"github.com/sevlyar/go-daemon"
)

var cronDaemon = false

func initCronCommand() *cobra.Command {
	cronStartCommand.Flags().BoolVarP(&cronDaemon, "daemon", "d", false, "start serve daemon")
	cronCommand.AddCommand(cronStartCommand)
	cronCommand.AddCommand(cronListCommand)
	return cronCommand
}

var cronCommand = &cobra.Command{
	Use:   "cron",
	Short: "定时任务相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}

		return nil
	},
}

var cronStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer() //获取容器

		//获取app服务
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 设置cron日志地址和进程id地址
		pidFolder := appService.RuntimeFolder()
		serverPidFile := filepath.Join(pidFolder, "cron.pid")
		logFolder := appService.LogFolder()
		serverLogFile := filepath.Join(logFolder, "cron.log")
		currentFolder := appService.BaseFolder()

		if cronDaemon {
			cntxt := &daemon.Context{
				//设置pid
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				//日志文件
				LogFileName: serverLogFile,
				LogFilePerm: 0664,
				WorkDir:     currentFolder,
				// 设置所有设置文件的mask
				Umask: 027,
				Args:  []string{"", "cron", "start", "--daemon=true"},
			}

			// 启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}

			if d != nil {
				// 父进程直接打印启动成功信息
				fmt.Println("cron serve started, pid", d.Pid)
				fmt.Println("log file: ", serverLogFile)
				return nil
			}

			// 子进程执行cron.run
			defer cntxt.Release()
			fmt.Println("daemon started")
			// gspt.SetProcTitle("anan cron")

			cmd.Root().Cron.Run()
			return nil
		}

		fmt.Println("start cron job")
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := ioutil.WriteFile(serverPidFile, []byte(content), 0664)
		if err != nil {
			return err
		}
		cmd.Root().Cron.Run()
		return nil
	},
}

var cronListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有的定时任务",
	RunE: func(cmd *cobra.Command, args []string) error {
		cronSpecs := cmd.Root().CronSpecs
		ps := [][]string{}
		for _, cronSpec := range cronSpecs {
			line := []string{cronSpec.Type, cronSpec.Spec, cronSpec.Cmd.Use, cronSpec.Cmd.Short, cronSpec.ServiceName}
			ps = append(ps, line)
		}
		util.PrettyPrint(ps)

		return nil
	},
}

/* var cronRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重启cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					return err
				}

				for i := 0; i < 10; i++ {
					if util.CheckProcessExist(pid) == false {
						break
					}
					time.Sleep(1 * time.Second)
				}
				fmt.Println("kill process:" + strconv.Itoa(pid))
			}
		}

		cronDaemon = true
		return cronStartCommand.RunE(cmd, args)
	},
} */
