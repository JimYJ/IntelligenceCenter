package llm

import (
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
	if !save(r) {
		response.Err(c, 400, response.ErrOperationFailed)
		return
	}
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
}
