package middleware

import (
	"github.com/gin-gonic/gin"
)

// Cors 允许跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		// origin := c.Request.Header.Get("origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Headers", "Access-Token")
		c.Next()
	}
}
