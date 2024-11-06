package task

import (
	"IntelligenceCenter/app/archive"
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/response"
	"IntelligenceCenter/service/log"
	"strings"

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
	if task.CrawlMode > 2 {
		response.Err(c, 400, "抓取模式不正确")
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
	if task.ExecType == 2 && task.CycleType > 2 {
		response.Err(c, 400, "执行周期设置不正确")
		return
	}
	if task.CycleType == 2 && len(task.WeekDays) == 0 {
		response.Err(c, 400, "执行周期是每周执行时，执行的每周日期不可为空")
		return
	}
	if len(task.WeekDays) > 0 {
		task.WeekDaysStr = strings.Join(task.WeekDays, ",")
	}
	if task.ExecType == 2 && len(task.ExecTime) == 0 {
		response.Err(c, 400, "执行周期是周期循环时，执行时间不可为空")
		return
	}
	if task.APISettingsID == nil || *task.APISettingsID == 0 {
		response.Err(c, 400, "选择内容提取模型的API设置不可为空")
		return
	}
	if task.APIModel == nil || len(*task.APIModel) == 0 {
		response.Err(c, 400, "提取模型不可为空")
		return
	}
	if task.CrawlMode == 1 {
		list := strings.Split(task.CrawlURL, "\n")
		for _, item := range list {
			if !utils.CheckURL(item) {
				response.Err(c, 400, "使用地址抓取模式时，抓取网页地址每行必须是http://或https://为前缀")
				return
			}
		}
	}
	if task.ArchiveOption == 1 {
		task.ArchiveID = int(archive.Create(task.TaskName))
	}
	if createtask(task) {
		response.Success(c, nil)
	} else {
		response.Err(c, 500, response.ErrOperationFailed)
	}
}
