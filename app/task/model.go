package task

type Task struct {
	// ExtractionMode         bool     `json:"extraction_mode" db:"extraction_mode"`                   // 抽取模式 1精准抽取 2智能抽取
	ID                     int      `json:"id" db:"id"`                                             // 主键
	ArchiveOption          uint8    `json:"archive_option,string" db:"-"`                           // 1新建档案 2选择档案
	ArchiveID              int      `json:"archive_id" db:"archive_id"`                             // 指定归档的档案ID
	ArchiveName            string   `json:"archive_name" db:"archive_name"`                         // 档案名称
	TaskName               string   `json:"task_name" db:"task_name"`                               // 任务名称
	CrawlMode              uint8    `json:"crawl_mode,string" db:"crawl_mode"`                      // 抓取模式 1地址抓取 2描述搜索抓取
	CrawlURL               string   `json:"crawl_url" db:"crawl_url"`                               // 抓取地址或抓取描述
	ExecType               uint8    `json:"exec_type,string" db:"exec_type"`                        // 执行类型 1-立即执行 2-周期循环
	CycleType              uint8    `json:"cycle_type,string" db:"cycle_type"`                      // 周期类型 1-每日 2-每周
	WeekDays               []string `json:"week_days" db:"-"`                                       // 指定周几执行，可多选
	WeekDaysStr            string   `json:"-" db:"week_days"`                                       // 指定周几执行，可多选，英文逗号隔开
	ExecTime               string   `json:"exec_time" db:"exec_time"`                               // 执行时间
	EnableAdvancedSettings bool     `json:"enable_advanced_settings" db:"enable_advanced_settings"` // 启用进阶设置
	TaskStatus             bool     `json:"task_status" db:"task_status"`                           // 任务状态
	EnableFilter           *bool    `json:"enable_filter" db:"enable_filter"`                       // 启用匹配过滤器
	DomainMatch            *string  `json:"domain_match" db:"domain_match"`                         // 域名匹配过滤器 为空则不生效
	PathMatch              *string  `json:"path_match" db:"path_match"`                             // 路径匹配过滤器 为空则不生效
	CrawlOption            *bool    `json:"crawl_option" db:"crawl_option"`                         // 抓取器设置 0自定义 1全局
	CrawlType              *uint8   `json:"crawl_type,string" db:"crawl_type"`                      // 抓取器选择 1 内置爬虫 2 headless浏览器 3 firecrawl
	ConcurrentCount        *int     `json:"concurrent_count" db:"concurrent_count"`                 // 并发数
	ScrapingInterval       *int     `json:"scraping_interval" db:"scraping_interval"`               // 抓取间隔(秒)
	GlobalScrapingDepth    *int     `json:"global_scraping_depth" db:"global_scraping_depth"`       // 抓取深度
	RequestRateLimit       *int     `json:"request_rate_limit" db:"request_rate_limit"`             // 每秒请求上限
	UseProxyIPPool         *bool    `json:"use_proxy_ip_pool" db:"use_proxy_ip_pool"`               // 使用代理IP池
	APISettingsID          *int     `json:"api_settings_id" db:"api_settings_id"`                   // API设置表ID
	APIModel               *string  `json:"extraction_model" db:"extraction_model"`                 // API指定LLM模型
	ApiType                uint8    `json:"api_type" db:"api_type"`                                 // API类型 1-OpenAI API Api 2-Ollama
	LLMSettingName         string   `json:"llm_setting_name" db:"llm_setting_name"`                 // LLM设置名称
	CreatedAt              string   `json:"created_at" db:"created_at"`                             // 更新时间
	UpdatedAt              string   `json:"updated_at" db:"updated_at"`                             // 创建时间
}