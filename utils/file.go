package utils

import "os"

type LFile struct{}

// 获取程序运行目录
func (l LFile) GetWorkDir() string {
	pwd, _ := os.Getwd()
	return pwd
}

func (l LFile) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return false, err
}
