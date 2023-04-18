package util

import (
	"os"
	"syscall"
)

//获取当前程序执行所在的目录
func GetExecDirectory() string {
	file, err := os.Getwd()
	if err == nil {
		return file + string(os.PathSeparator)
	}

	return ""
}

func CheckProcessExist(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}

	return true
}
