package router

import (
	"IntelligenceCenter/service/log"
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Web(static embed.FS) {
	router := gin.New()
	router.Use(log.Logs())
	router.Use(log.Recovery())
	st, err := fs.Sub(static, "static/dist")
	if err != nil {
		log.Info("加载静态资源失败:", err)
	}
	router.StaticFS("/", http.FS(st))
	router.NoRoute(func(c *gin.Context) {
		data, err := static.ReadFile("static/dist/index.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})
	err = router.Run(":6060")
	if err != nil {
		log.Info("启动静态页面失败:", err)
	}
}
