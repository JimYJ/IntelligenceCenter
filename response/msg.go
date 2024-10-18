package response

const (
	ErrInvalidRequestParam     = "请求参数不正确"
	ErrIDCannotBeZeroOrEmpty   = "ID不可为0或空"
	ErrApiKeyLengthExceeded    = "API Key长度不得为空或超过256"
	ErrApiURLLengthExceeded    = "API 地址长度不得超过2048"
	ErrRemarkLengthExceeded    = "备注长度不得大于1024"
	ErrTimeoutInvalid          = "超时时间不得为0或大于1小时"
	ErrRequestRateLimitInvalid = "每秒请求数不得为0或大于10000"
	ErrOperationFailed         = "操作失败，请前往https://github.com/JimYJ/IntelligenceCenter/issues 提供错误日志"

	MaxApiKeyLength     = 256
	MaxApiURLLength     = 2048
	MaxRemarkLength     = 1024
	MinTimeout          = 1     // 超时不得为0
	MaxTimeout          = 3600  // 超时不得大于1小时
	MinRequestRateLimit = 1     // 每秒请求数不得为0
	MaxRequestRateLimit = 10000 // 每秒请求数不得大于10000
)
