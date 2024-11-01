package router

import (
	"IntelligenceCenter/app/archive"
	"IntelligenceCenter/app/common"
	"IntelligenceCenter/app/llm"
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
	llmseting.POST("/add", llm.Create)
	llmseting.GET("/del", llm.Del)
	llmseting.POST("/edit", llm.Edit)
	llmseting.POST("/list", llm.ListByPage)
	// 档案
	archiveDoc := api.Group("/archive")
	archiveDoc.OPTIONS(":any", common.Ok)
	archiveDoc.OPTIONS("/doc/:any", common.Ok)
	archiveDoc.POST("/list", archive.ArchiveListByPage)
	archiveDoc.GET("/info", archive.ArchiveInfo)
	archiveDoc.POST("/doc/list", archive.DocListByPage)
	err := router.Run(":6061")
	if err != nil {
		log.Logger.Println(err)
	}
}
