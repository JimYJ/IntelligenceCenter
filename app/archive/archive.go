package archive

import (
	"IntelligenceCenter/app/common"
	"IntelligenceCenter/response"
	"IntelligenceCenter/service/log"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	DocResourceChan = make(chan *DocResource, 65535)
	DocResourceQuit = make(chan bool)
)

func init() {
	go SaveDocResourceBatch()
}

// 档案分页
func ListByPage(c *gin.Context) {
	k := &common.Keyword{}
	err := c.ShouldBindJSON(k)
	if err != nil {
		response.Err(c, 400, "请求参数不正确")
	}
	pageNo, pageSize := common.PageParams(c)
	totalCount := CountRecord(k.Keyword)
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
	list := archiveListByPage(start, pageSize, k.Keyword)
	pager.Data = list
	pager.Keyword = k.Keyword
	response.Success(c, pager)
}

// 档案列表
func List(c *gin.Context) {
	response.Success(c, archiveList())
}

// 档案信息
func Info(c *gin.Context) {
	id := c.Query("id")
	data := archiveInfo(id)
	data.FileCount = DocCountRecord(id, "")
	data.TaskCount = archiveTask(id, -1)
	data.ActiveTaskCount = archiveTask(id, 1)
	response.Success(c, data)
}

// 文档分页
func DocListByPage(c *gin.Context) {
	k := &common.Keyword{}
	err := c.ShouldBindJSON(k)
	if err != nil {
		response.Err(c, 400, "请求参数不正确")
	}
	id := c.Query("id")
	pageNo, pageSize := common.PageParams(c)
	totalCount := DocCountRecord(id, k.Keyword)
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
	list := docListByPage(start, pageSize, id, k.Keyword)
	pager.Data = list
	pager.Keyword = k.Keyword
	response.Success(c, pager)
}

// 分批保存资源
func SaveDocResourceBatch() {
	var resources []*DocResource
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case resource := <-DocResourceChan:
			resources = append(resources, resource)
			if len(resources) >= 1000 {
				err := SaveDocResource(resources)
				if err != nil {
					log.Info("批量保存资源失败:", err)
				}
				resources = resources[:0]
				ticker.Reset(10 * time.Second)
			}
		case <-ticker.C:
			if len(resources) > 0 {
				err := SaveDocResource(resources)
				if err != nil {
					log.Info("批量保存资源失败:", err)
				}
				resources = resources[:0]
			}
		case <-DocResourceQuit:
			if len(resources) > 0 {
				err := SaveDocResource(resources)
				if err != nil {
					log.Info("批量保存资源失败:", err)
				}
			}
			return // 退出函数
		}
	}
}
