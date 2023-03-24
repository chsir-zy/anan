package util

import "os"

// 判断文件、文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}

	return true
}
