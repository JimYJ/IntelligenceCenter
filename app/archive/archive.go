package archive

import (
	"IntelligenceCenter/app/common"
	"IntelligenceCenter/response"

	"github.com/gin-gonic/gin"
)

func ArchiveListByPage(c *gin.Context) {
	keyword := c.Query("keyword")
	pageNo, pageSize := common.PageParams(c)
	totalCount := archiveCountRecord(keyword)
	if totalCount == 0 {
		response.Success(c, &common.PageInfo{
			PageNo:      pageNo,
			TotalRecord: 0,
			TotalPage:   0,
			PageSize:    pageSize,
		})
		return
	}
	pager, start := common.Page(totalCount, pageSize, pageNo)
	list := archiveListByPage(start, pageSize, keyword)
	pager.Data = list
	response.Success(c, pager)
}

func DocListByPage(c *gin.Context) {
	keyword := c.Query("keyword")
	pageNo, pageSize := common.PageParams(c)
	totalCount := docCountRecord(keyword)
	if totalCount == 0 {
		response.Success(c, &common.PageInfo{
			PageNo:      pageNo,
			TotalRecord: 0,
			TotalPage:   0,
			PageSize:    pageSize,
		})
		return
	}
	pager, start := common.Page(totalCount, pageSize, pageNo)
	list := docListByPage(start, pageSize, keyword)
	pager.Data = list
	response.Success(c, pager)
}
