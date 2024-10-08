package router

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)


func Web(static embed.FS) {
	r := gin.Default()
	st, _ := fs.Sub(static, "static/dist")
	r.StaticFS("/", http.FS(st))
	r.NoRoute(func(c *gin.Context) {
		data, err := static.ReadFile("static/dist/index.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	err := r.Run(":6060")
	if err != nil {
		fmt.Println(err)
	}
}

func Api() {
	r := gin.Default()
	r.GET("/api/test",func(c *gin.Context) {
		c.JSON(200, gin.H{
		  "msg": "hello world",
		})
	})
	err := r.Run(":6061")
	if err != nil {
		fmt.Println(err)
	}
}