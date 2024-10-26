package common

import (
	"log"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	PageSize = 20
)

var (
	InviteExpir int64 = 14 * 24 * 3600 //14天有效期
)

// 分页信息
type PageInfo struct {
	PageNo      int    `json:"current"`
	TotalRecord int    `json:"total"`
	TotalPage   int    `json:"pages"`
	PageSize    int    `json:"size"`
	Keyword     string `json:"keyword"`
	Data        any    `json:"records"`
}

// totalPage 总页数
func totalPage(totalCount, maxPageSize int) int {
	var pageTotal int
	if totalCount < maxPageSize {
		return 1
	}
	pageTotal = int(math.Ceil(float64(totalCount) / float64(maxPageSize)))
	return pageTotal
}

// 记录分页信息
func Page(totalCount, pageSize, pageNo int) (*PageInfo, int) {
	pageInfo := &PageInfo{
		PageSize:    pageSize,
		PageNo:      pageNo,
		TotalRecord: totalCount,
	}
	pageInfo.TotalPage = totalPage(totalCount, pageSize)
	if pageInfo.PageNo > pageInfo.TotalPage {
		pageInfo.PageNo = pageInfo.TotalPage
		pageNo = pageInfo.TotalPage
	}
	return pageInfo, (pageNo - 1) * pageSize
}

func PageParams(c *gin.Context) (int, int) {
	curPage := c.Query("current")
	pageSize := c.Query("size")
	var cp, ps int
	var err error
	if len(curPage) == 0 {
		cp = 1
	} else {
		cp, err = strconv.Atoi(curPage)
		if err != nil {
			log.Println("curPage转换错误")
			cp = 1
		}
	}
	if len(pageSize) == 0 {
		ps = PageSize
	} else {
		ps, err = strconv.Atoi(pageSize)
		if err != nil {
			log.Println("pageSize转换错误")
			ps = PageSize
		}
	}
	return cp, ps
}
