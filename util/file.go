package util

import (
	"fmt"
	"os"
)

// GetProjectAbsPath 获取当前项目的绝对路径
func GetProjectAbsPath() (string, error) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取当前目录失败: %v", err)
	}
	return currentDir, nil
}

// PathExists 判断路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
