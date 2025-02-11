package llm

type Request struct {
	// UseProxyPool     bool    `db:"use_proxy_pool" json:"use_proxy_pool"`         // 使用代理 IP 池 0否1是
	ID               int     `db:"id" json:"id"`                                        // 主键
	Name             string  `db:"name" json:"name"`                                    // 配置名称
	ApiType          uint8   `db:"api_type" json:"api_type,string"`                     // API类型 1-OpenAI Api 2-Ollama
	ApiURL           string  `db:"api_url" json:"api_url"`                              // API 地址
	ApiKey           string  `db:"api_key" json:"api_key"`                              // API 密钥
	Timeout          int     `db:"timeout" json:"timeout,string"`                       // 超时设置(秒),默认30秒
	RequestRateLimit int     `db:"request_rate_limit" json:"request_rate_limit,string"` // 每秒请求上限
	Remark           *string `db:"remark" json:"remark"`                                // 描述信息
}

type Keyword struct {
	Keyword string `db:"-" json:"keyword"` // 描述信息
}

type LLMType struct {
	Value    uint8      `json:"value"`
	Label    string     `json:"label"`
	Children []*Setting `json:"children"`
}

type Setting struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}
