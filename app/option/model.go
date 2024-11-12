package option

type Option struct {
	// ID                       int  `db:"id" json:"id"`
	CrawlType                uint8 `db:"crawl_type" json:"crawl_type"`                                   // 抓取器选择 1 内置爬虫 2 headless浏览器 3 firecrawl
	ConcurrentCount          int   `db:"concurrent_count" json:"concurrent_count"`                       // 并发数
	ScrapingInterval         int   `db:"scraping_interval" json:"scraping_interval"`                     // 抓取间隔(秒)
	GlobalScrapingDepth      int   `db:"global_scraping_depth" json:"global_scraping_depth"`             // 抓取深度
	RequestRateLimit         int   `db:"request_rate_limit" json:"request_rate_limit"`                   // 每秒请求上限
	UseProxyIPPool           bool  `db:"use_proxy_ip_pool" json:"use_proxy_ip_pool"`                     // 使用代理IP池
	UseGlobalConcurrencyPool bool  `db:"use_global_concurrency_pool" json:"use_global_concurrency_pool"` // 使用全局并发池
}
