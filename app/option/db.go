package option

import (
	"IntelligenceCenter/common/sqlite"
	"IntelligenceCenter/service/log"
)

func GetOption() Option {
	sql := `SELECT
				id,
				crawl_type,
				concurrent_count,
				scraping_interval,
				global_scraping_depth,
				request_rate_limit,
				use_proxy_ip_pool,
				use_global_concurrency_pool
			FROM
				option 
			WHERE
				id = 1;`
	var o Option
	err := sqlite.Conn().Get(&o, sql)
	if err != nil {
		log.Info("查询全局设置出错:", err)
		return o
	}
	return o
}
