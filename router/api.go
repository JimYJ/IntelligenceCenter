package router

import (
	"IntelligenceCenter/app/archive"
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
	api.OPTIONS(":any")
	// LLM API
	llmseting := api.Group("/llm")
	llmseting.POST("/add", llm.Create)
	llmseting.GET("/del", llm.Del)
	llmseting.POST("/edit", llm.Edit)
	llmseting.POST("/list", llm.ListByPage)
	// 档案
	archiveDoc := api.Group("/archive")
	archiveDoc.POST("/list", archive.ArchiveListByPage)
	archiveDoc.POST("/doc/list", archive.DocListByPage)
	router.GET("/ping", func(c *gin.Context) {
		// log.Info("Handling /ping request", "thfghdf")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := router.Run(":6061")
	if err != nil {
		log.Logger.Println(err)
	}
}
