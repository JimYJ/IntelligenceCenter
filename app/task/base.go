package task

import (
	"IntelligenceCenter/app/archive"
	"IntelligenceCenter/app/common"
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
		response.Err(c, 400, response.ErrTaskNameCannotBeEmpty)
		return
	}
	if task.CrawlMode > 2 {
		response.Err(c, 400, response.ErrInvalidCrawlMode)
		return
	}
	if len(task.CrawlURL) == 0 {
		response.Err(c, 400, response.ErrCrawlURLCannotBeEmpty)
		return
	}
	if task.ExecType > 2 {
		response.Err(c, 400, response.ErrInvalidExecType)
		return
	}
	if task.ExecType == 2 && task.CycleType > 2 {
		response.Err(c, 400, response.ErrInvalidCycleType)
		return
	}
	if task.CycleType == 2 && len(task.WeekDays) == 0 {
		response.Err(c, 400, response.ErrWeekDaysCannotBeEmpty)
		return
	}
	if len(task.WeekDays) > 0 {
		task.WeekDaysStr = strings.Join(task.WeekDays, ",")
	}
	if task.ExecType == 2 && len(task.ExecTime) == 0 {
		response.Err(c, 400, response.ErrExecTimeCannotBeEmpty)
		return
	}
	if len(task.APISettingsIDList) == 0 {
		response.Err(c, 400, response.ErrAPISettingsCannotBeEmpty)
		return
	}
	task.APISettingsIDStr = strings.Join(utils.ConvertIntsToStrings(task.APISettingsIDList), ",")
	task.APISettingsID = task.APISettingsIDList[len(task.APISettingsIDList)-1]
	if task.APIModel == nil || len(*task.APIModel) == 0 {
		response.Err(c, 400, response.ErrAPIModelCannotBeEmpty)
		return
	}
	if task.CrawlMode == 1 {
		list := strings.Split(task.CrawlURL, "\n")
		for _, item := range list {
			if !utils.CheckURL(item) {
				response.Err(c, 400, response.ErrInvalidCrawlURLPrefix)
				return
			}
		}
	}
	if task.ArchiveOption == 1 {
		task.ArchiveID = int(archive.Create(task.TaskName))
	}
	if createtask(task) {
		task.Free()
		NewTaskChan <- true
		response.Success(c, nil)
	} else {
		response.Err(c, 500, response.ErrOperationFailed)
	}
}

// 档案分页
func ListByPage(c *gin.Context) {
	k := &common.Keyword{}
	err := c.ShouldBindJSON(k)
	if err != nil {
		response.Err(c, 400, response.ErrInvalidRequestParam)
	}
	pageNo, pageSize := common.PageParams(c)
	totalCount := taskCount(-1, k.Keyword)
	if totalCount == 0 {
		response.Success(c, &common.PageInfo{
			PageNo:      pageNo,
			TotalRecord: 0,
			TotalPage:   0,
			PageSize:    pageSize,
			Keyword:     k.Keyword,
		})
		return
	}
	pager, start := common.Page(totalCount, pageSize, pageNo)
	list := taskListByPage(start, pageSize, k.Keyword)
	for _, item := range list {
		if len(item.WeekDaysStr) != 0 {
			item.WeekDays = strings.Split(item.WeekDaysStr, ",")
		}
	}
	for _, item := range list {
		if len(item.APISettingsIDStr) != 0 {
			item.APISettingsIDList, err = utils.ConvertStringsToInts(strings.Split(item.APISettingsIDStr, ","))
			if err != nil {
				response.Err(c, 500, response.ErrOperationFailed)
			}
		}
	}
	pager.Data = list
	pager.Keyword = k.Keyword
	response.Success(c, pager)
}

// 任务系统信息
func TaskInfo(c *gin.Context) {
	data := &TaskData{}
	data.TaskCount = taskCount(-1, "")
	data.ActiveTaskCount = taskCount(1, "")
	data.ArchiveCount = archive.CountRecord("")
	data.ArchiveDocsCount = archive.DocCountRecord("", "")
	data.ArchiveDocsResCount = archive.DocResCountRecord("")
	response.Success(c, data)
}
