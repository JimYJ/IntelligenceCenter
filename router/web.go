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
