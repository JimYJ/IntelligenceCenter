package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Api() {
	r := gin.Default()
	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "hello world",
		})
	})
	err := r.Run(":6061")
	if err != nil {
		fmt.Println(err)
	}
}