package main

import (
	"IntelligenceCenter/app/db"
	"IntelligenceCenter/app/task"
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/router"
	"IntelligenceCenter/service/timer"
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
	go task.ListenNewTask()
	go task.Scan()
	go timer.DayTaskFor0AM()
	go router.Web(static)
	router.Api()
}

// 初始化必要目录
func initDir() {
	for _, item := range needDir {
		utils.CreateDir(item)
	}
}
