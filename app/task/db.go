package task

import (
	"IntelligenceCenter/common/sqlite"
	"IntelligenceCenter/service/log"
)

func createtask(task *Task) bool {
	sql := `INSERT INTO task (
            archive_id,
            task_name,
            crawl_url,
            exec_type,
            cycle_type,
            week_days,
            exec_time,
            enable_filter,
            domain_match,
            path_match,
            crawl_option,
            crawl_type,
            concurrent_count,
            scraping_interval,
            global_scraping_depth,
            request_rate_limit,
            use_proxy_ip_pool,
            enable_advanced_settings,
            api_settings_id,
            api_model
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := sqlite.Conn().Exec(sql,
		task.ArchiveID, task.TaskName, task.CrawlURL, task.ExecType, task.CycleType, task.WeekDays, task.ExecTime,
		task.EnableFilter, task.DomainMatch, task.PathMatch, task.CrawlOption, task.CrawlType, task.ConcurrentCount,
		task.ScrapingInterval, task.GlobalScrapingDepth, task.RequestRateLimit, task.UseProxyIPPool, task.EnableAdvancedSettings,
		task.APISettingsID, task.APIModel)
	if err != nil {
		log.Info("创建任务出错:", err)
		return false
	}
	return true
}
