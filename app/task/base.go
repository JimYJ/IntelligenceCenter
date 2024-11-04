package task

import (
	"IntelligenceCenter/response"
	"IntelligenceCenter/service/log"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	task := &Task{}
	err := c.ShouldBindJSON(task)
	if err != nil {
		log.Info(err)
		response.Err(c, 400, response.ErrInvalidRequestParam)
		return
	}
	if len(task.TaskName) == 0 {
		response.Err(c, 400, "任务名称不可为空")
		return
	}
	if len(task.CrawlURL) == 0 {
		response.Err(c, 400, "信息抓取网址不可为空")
		return
	}
	if task.ExecType > 2 {
		response.Err(c, 400, "执行类型不正确")
		return
	}
	if createtask(task) {
		response.Success(c, nil)
	} else {
		response.Err(c, 500, response.ErrOperationFailed)
	}
}
