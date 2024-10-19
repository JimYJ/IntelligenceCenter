package router

import (
	"IntelligenceCenter/app/llm"
	"IntelligenceCenter/service/log"

	"github.com/gin-gonic/gin"
)

func Api() {
	router := gin.New()
	router.Use(log.Logs())
	router.Use(log.Recovery())
	api := router.Group("/api")
	llmseting := api.Group("/llm")
	llmseting.POST("/add", llm.Create)
	llmseting.GET("/del", llm.Del)
	llmseting.POST("/edit", llm.Edit)
	llmseting.POST("/list", llm.ListByPage)
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
