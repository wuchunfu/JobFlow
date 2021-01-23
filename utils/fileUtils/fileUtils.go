package fileUtils

import (
	"github.com/sirupsen/logrus"
	"os"
)

// 递归创建目录
func MkdirAll(filePath string) bool {
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		logrus.Errorf("mkdir %s failed. %s", filePath, err)
		return false
	}
	return true
}

// 判断文件夹是否存在
func IsExistPath(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// 判断文件是否存在及是否有权限访问
func FileExist(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	if os.IsPermission(err) {
		return false
	}
	return true
}