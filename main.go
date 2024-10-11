package main

import (
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/router"
	"embed"
)

var (
	// 检查初始化目录
	needDir = []string{"logs", "extraction-rules"}
)

//go:embed static/dist
var static embed.FS

func main() {
	initDir()
	go router.Web(static)
	router.Api()
}

// 初始化必要目录
func initDir() {
	for _, item := range needDir {
		utils.CreateDir(item)
	}
}
