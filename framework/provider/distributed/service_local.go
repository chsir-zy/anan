package distributed

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
)

type LocalDistributedService struct {
	container framework.Container
}

func NewLocalDistributedService(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("param error")
	}

	container := params[0].(framework.Container)
	return &LocalDistributedService{container: container}, nil
}

func (s LocalDistributedService) Select(serviceName string, appID string, holdTime time.Duration) (selectAppID string, err error) {
	appService := s.container.MustMake(contract.AppKey).(contract.App)
	runtimeFolder := appService.RuntimeFolder()
	lockFile := filepath.Join(runtimeFolder, "distribute_"+serviceName)

	// 打开文件锁
	lock, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	// 尝试独占文件锁
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

	// 抢不到文件锁
	if err != nil {
		// 读取被选择的appid
		selectAppIDByt, err := ioutil.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(selectAppIDByt), err
	}

	go func() {
		defer func() {
			// 释放文件锁
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)

			// 释放文件
			lock.Close()

			os.Remove(lockFile)
		}()

		timer := time.NewTimer(holdTime)
		// 等待计时器结束
		<-timer.C
	}()

	// 这里已经是抢占到了，将抢占到的appID写入文件
	if _, err := lock.WriteString(appID); err != nil {
		return "", err
	}

	return appID, nil
}
