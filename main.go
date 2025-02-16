package main

import (
	"IntelligenceCenter/app/db"
	"IntelligenceCenter/app/proxy"
	"IntelligenceCenter/app/task"
	"IntelligenceCenter/common"
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/router"
	"IntelligenceCenter/service/timer"
	"embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed static/dist
var static embed.FS

func main() {
	initDir()
	db.CheckDatabase()
	task.Scan()
	go task.ListenNewTask()
	go task.Listen()
	go timer.DayTaskFor0AM()
	go router.Web(static)
	go task.Retry()
	go task.ListenMatch()
	go proxy.Run()
	router.Api()
}

// 初始化必要目录
func initDir() {
	for _, item := range common.NeedDir {
		utils.CreateDir(item)
	}
}
