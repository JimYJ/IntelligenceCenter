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
		response.Err(c, 400, response.ErrInvalidRequestParam)
	}
	pageNo, pageSize := common.PageParams(c)
	totalCount := countRecord(k.Keyword)
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
	list := listByPage(start, pageSize, k.Keyword)
	pager.Data = list
	pager.Keyword = k.Keyword
	response.Success(c, pager)
}

func ListByGroup(c *gin.Context) {
	list := list()
	temp := make(map[uint8][]*Request)
	for _, item := range list {
		temp[item.ApiType] = append(temp[item.ApiType], item)
	}
	tempList := make([]*LLMType, 0)
	for k, v := range temp {
		llmType := &LLMType{
			Value:    k,
			Label:    getType(k),
			Children: make([]*Setting, 0),
		}
		tempList = append(tempList, llmType)
		for _, item := range v {
			llmType.Children = append(llmType.Children, &Setting{
				Value: item.ID,
				Label: item.Name,
			})
		}
	}
	response.Success(c, tempList)
}

func getType(t uint8) string {
	if t == 1 {
		return "OpenAI API"
	} else if t == 2 {
		return "OLlama API"
	}
	return ""
}
