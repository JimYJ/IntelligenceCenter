package response

const (
	ErrInvalidRequestParam      = "请求参数不正确"
	ErrIDCannotBeZeroOrEmpty    = "ID不可为0或空"
	ErrApiKeyLengthExceeded     = "API Key长度不得为空或超过256"
	ErrApiURLLengthExceeded     = "API 地址长度不得超过2048"
	ErrRemarkLengthExceeded     = "备注长度不得大于1024"
	ErrTimeoutInvalid           = "超时时间不得为0或大于1小时"
	ErrRequestRateLimitInvalid  = "每秒请求数不得为0或大于10000"
	ErrOperationFailed          = "操作失败，请前往https://github.com/JimYJ/IntelligenceCenter/issues 提供错误日志"
	ErrTaskNameCannotBeEmpty    = "任务名称不可为空"
	ErrInvalidCrawlMode         = "抓取模式不正确"
	ErrCrawlURLCannotBeEmpty    = "信息抓取网址不可为空"
	ErrInvalidExecType          = "执行类型不正确"
	ErrInvalidCycleType         = "执行周期设置不正确"
	ErrWeekDaysCannotBeEmpty    = "执行周期是每周执行时，执行的每周日期不可为空"
	ErrExecTimeCannotBeEmpty    = "执行周期是周期循环时，执行时间不可为空"
	ErrAPISettingsCannotBeEmpty = "选择内容提取模型的API设置不可为空"
	ErrInvalidCrawlURLPrefix    = "使用地址抓取模式时，抓取网页地址每行必须是http://或https://为前缀"
	ErrAPIModelCannotBeEmpty    = "提取模型不可为空"

	MaxApiKeyLength     = 256
	MaxApiURLLength     = 2048
	MaxRemarkLength     = 1024
	MinTimeout          = 1     // 超时不得为0
	MaxTimeout          = 3600  // 超时不得大于1小时
	MinRequestRateLimit = 1     // 每秒请求数不得为0
	MaxRequestRateLimit = 10000 // 每秒请求数不得大于10000
)
