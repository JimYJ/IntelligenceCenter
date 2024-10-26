package common

import (
	"github.com/gin-gonic/gin"
)

// Ok options专用响应
func Ok(c *gin.Context) {
	c.String(200, "")
}
