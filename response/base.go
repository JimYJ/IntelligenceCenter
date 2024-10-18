package response

import (
	"github.com/gin-gonic/gin"
)

// Success 统一响应成功结构体
func Success(c *gin.Context, data interface{}) {
	respon := &Respon{}
	respon.Code = 200
	respon.Success = true
	respon.Message = "Success"
	respon.Data = data
	c.JSON(200, respon)
	c.Abort()
}

// Err400 请求错误
func Err(c *gin.Context, code int, msg string) {
	respon := &Respon{}
	respon.Code = code
	respon.Success = false
	respon.Message = msg
	c.JSON(200, respon)
	c.Abort()
}
