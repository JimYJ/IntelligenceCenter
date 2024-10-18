package main

import (
	"IntelligenceCenter/app/db"
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/router"
	"embed"

	_ "github.com/mattn/go-sqlite3"
)

var (
	// 检查初始化目录
	needDir = []string{"logs", "extraction-rules", "database", "proxy-ip"}
)

//go:embed static/dist
var static embed.FS

func main() {
	initDir()
	db.CheckDatabase()
	go router.Web(static)
	router.Api()
}

// 初始化必要目录
func initDir() {
	for _, item := range needDir {
		utils.CreateDir(item)
	}
}
