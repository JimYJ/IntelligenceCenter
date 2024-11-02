package archive

import (
	"IntelligenceCenter/app/common"
	"IntelligenceCenter/response"

	"github.com/gin-gonic/gin"
)

// 档案分页
func ArchiveListByPage(c *gin.Context) {
	k := &Keyword{}
	err := c.ShouldBindJSON(k)
	if err != nil {
		response.Err(c, 400, "请求参数不正确")
	}
	pageNo, pageSize := common.PageParams(c)
	totalCount := archiveCountRecord(k.Keyword)
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
func ArchiveList(c *gin.Context) {
	response.Success(c, archiveList)
}

// 档案信息
func ArchiveInfo(c *gin.Context) {
	id := c.Query("id")
	data := archiveInfo(id)
	data.FileCount = docCountRecord(id, "")
	data.TaskCount = archiveTask(id, -1)
	data.ActiveTaskCount = archiveTask(id, 1)
	response.Success(c, data)
}

// 文档分页
func DocListByPage(c *gin.Context) {
	k := &Keyword{}
	err := c.ShouldBindJSON(k)
	if err != nil {
		response.Err(c, 400, "请求参数不正确")
	}
	id := c.Query("id")
	pageNo, pageSize := common.PageParams(c)
	totalCount := docCountRecord(id, k.Keyword)
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
