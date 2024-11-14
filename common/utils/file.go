package utils

import (
	"IntelligenceCenter/service/log"
	"os"
	"path/filepath"
	"strings"
)

// 查找指定目录下的指定后缀文件
func FindFilesBySuffix(dir string, suffixes ...string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileName := info.Name()
		log.Info(fileName)
		for _, suffix := range suffixes {
			if strings.HasSuffix(fileName, suffix) {
				files = append(files, path)
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
