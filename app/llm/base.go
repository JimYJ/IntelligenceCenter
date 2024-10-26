package llm

import (
	"IntelligenceCenter/app/common"
	"IntelligenceCenter/response"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	r := &Request{}
	err := c.ShouldBindJSON(r)
	if err != nil {
		response.Err(c, 400, response.ErrInvalidRequestParam)
		return
	}
	if len(r.ApiKey) > response.MaxApiKeyLength {
		response.Err(c, 400, response.ErrApiKeyLengthExceeded)
		return
	}

	if len(r.ApiURL) > response.MaxApiURLLength {
		response.Err(c, 400, response.ErrApiURLLengthExceeded)
		return
	}

	if r.Remark != nil && len(*r.Remark) > response.MaxRemarkLength {
		response.Err(c, 400, response.ErrRemarkLengthExceeded)
		return
	}

	if r.Timeout < response.MinTimeout || r.Timeout > response.MaxTimeout {
		response.Err(c, 400, response.ErrTimeoutInvalid)
		return
	}

	if r.RequestRateLimit < response.MinRequestRateLimit || r.RequestRateLimit > response.MaxRequestRateLimit {
		response.Err(c, 400, response.ErrRequestRateLimitInvalid)
		return
	}
	if !create(r) {
		response.Err(c, 400, response.ErrOperationFailed)
		return
	}
	response.Success(c, nil)
}

func Del(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		response.Err(c, 400, response.ErrIDCannotBeZeroOrEmpty)
		return
	}
	if !del(id) {
		response.Err(c, 400, response.ErrOperationFailed)
		return
	}
	response.Success(c, nil)
}

func Edit(c *gin.Context) {
	r := &Request{}
	err := c.ShouldBindJSON(r)
	if err != nil {
		response.Err(c, 400, response.ErrInvalidRequestParam)
		return
	}
	if r.ID == 0 {
		response.Err(c, 400, response.ErrIDCannotBeZeroOrEmpty)
		return
	}
	if len(r.ApiKey) > response.MaxApiKeyLength {
		response.Err(c, 400, response.ErrApiKeyLengthExceeded)
		return
	}

	if len(r.ApiURL) > response.MaxApiURLLength {
		response.Err(c, 400, response.ErrApiURLLengthExceeded)
		return
	}

	if r.Remark != nil && len(*r.Remark) > response.MaxRemarkLength {
		response.Err(c, 400, response.ErrRemarkLengthExceeded)
		return
	}

	if r.Timeout < response.MinTimeout || r.Timeout > response.MaxTimeout {
		response.Err(c, 400, response.ErrTimeoutInvalid)
		return
	}

	if r.RequestRateLimit < response.MinRequestRateLimit || r.RequestRateLimit > response.MaxRequestRateLimit {
		response.Err(c, 400, response.ErrRequestRateLimitInvalid)
		return
	}
	if !edit(r) {
		response.Err(c, 400, response.ErrOperationFailed)
		return
	}
	response.Success(c, nil)
}

func ListByPage(c *gin.Context) {
	k := &Keyword{}
	err := c.ShouldBindJSON(k)
	if err != nil {
		response.Err(c, 400, "请求参数不正确")
	}
	pageNo, pageSize := common.PageParams(c)
	totalCount := countRecord(k.Keyword)
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
	list := listByPage(start, pageSize, k.Keyword)
	pager.Data = list
	response.Success(c, pager)
}
