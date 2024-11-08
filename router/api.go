package router

import (
	"IntelligenceCenter/app/archive"
	"IntelligenceCenter/app/common"
	"IntelligenceCenter/app/llm"
	"IntelligenceCenter/app/task"
	"IntelligenceCenter/router/middleware"
	"IntelligenceCenter/service/log"

	"github.com/gin-gonic/gin"
)

func Api() {
	router := gin.New()
	router.Use(log.Logs())
	router.Use(log.Recovery())
	router.Use(middleware.Cors())
	api := router.Group("/api")
	// LLM API
	llmseting := api.Group("/llm")
	llmseting.OPTIONS(":any", common.Ok)
	llmseting.POST("/create", llm.Create)
	llmseting.GET("/del", llm.Del)
	llmseting.POST("/edit", llm.Edit)
	llmseting.POST("/list", llm.ListByPage)
	llmseting.GET("/list/grouping", llm.ListByGroup)
	// 档案
	archiveDoc := api.Group("/archive")
	archiveDoc.OPTIONS(":any", common.Ok)
	archiveDoc.OPTIONS("/doc/:any", common.Ok)
	archiveDoc.POST("/list", archive.ListByPage)
	archiveDoc.GET("/list/select", archive.List)
	archiveDoc.GET("/info", archive.Info)
	archiveDoc.POST("/doc/list", archive.DocListByPage)
	// 任务
	taskApi := api.Group("/task")
	taskApi.OPTIONS(":any", common.Ok)
	taskApi.POST("/create", task.Create)
	taskApi.POST("/list", task.ListByPage)
	taskApi.GET("/info", task.TaskInfo)
	err := router.Run(":6061")
	if err != nil {
		log.Logger.Println(err)
	}
}
