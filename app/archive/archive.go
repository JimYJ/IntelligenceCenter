package archive

import (
	"IntelligenceCenter/app/common"
	"IntelligenceCenter/response"

	"github.com/gin-gonic/gin"
)

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
		})
		return
	}
	pager, start := common.Page(totalCount, pageSize, pageNo)
	list := archiveListByPage(start, pageSize, k.Keyword)
	pager.Data = list
	response.Success(c, pager)
}

func DocListByPage(c *gin.Context) {
	k := &Keyword{}
	err := c.ShouldBindJSON(k)
	if err != nil {
		response.Err(c, 400, "请求参数不正确")
	}
	id := c.Query("id")
	pageNo, pageSize := common.PageParams(c)
	totalCount := docCountRecord(k.Keyword)
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
	list := docListByPage(start, pageSize, id, k.Keyword)
	pager.Data = list
	response.Success(c, pager)
}
