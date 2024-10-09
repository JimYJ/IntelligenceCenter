package router

import (
	"IntelligenceCenter/service/log"

	"github.com/gin-gonic/gin"
)

func Api() {
	router:=gin.New()
	router.Use(log.Logs())
	router.Use(log.Recovery())
	api:=router.Group("/api")
	api.POST("/test")
	router.GET("/ping", func(c *gin.Context) {
        log.Info("Handling /ping request")
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
	err := router.Run(":6061")
	if err != nil {
		log.Logger.Println(err)
	}
}