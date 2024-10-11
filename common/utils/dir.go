package utils

import (
	"IntelligenceCenter/service/log"
	"os"
)

// IsDir 判断文件是否存在，是否目录 返回：文件是否存在，是否目录
func IsDir(path string) (bool, bool) {
	fileinfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false
		}
		log.Info(err)
	}
	return true, fileinfo.IsDir()
}

// CreateDir 创建目录
func CreateDir(logPath string) {
	_, isDir := IsDir(logPath)
	if !isDir {
		err := os.Mkdir(logPath, os.ModePerm)
		if err != nil {
			log.Info("创建日志目录失败：", err)
		}
	}
}
