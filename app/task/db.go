package task

import (
	"IntelligenceCenter/app/common"
	"IntelligenceCenter/common/sqlite"
	"IntelligenceCenter/service/log"
	"fmt"
)

func createtask(task *Task) bool {
	sql := `INSERT INTO task (
            archive_id,
            task_name,
			crawl_mode,
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
			api_settings_id_list,
            api_model
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,  ?, ?, ?, ?, ?)`
	_, err := sqlite.Conn().Exec(sql,
		task.ArchiveID, task.TaskName, task.CrawlMode, task.CrawlURL, task.ExecType, task.CycleType, task.WeekDaysStr, task.ExecTime,
		task.EnableFilter, task.DomainMatch, task.PathMatch, task.CrawlOption, task.CrawlType, task.ConcurrentCount,
		task.ScrapingInterval, task.GlobalScrapingDepth, task.RequestRateLimit, task.UseProxyIPPool, task.EnableAdvancedSettings,
		task.APISettingsID, task.APISettingsIDStr, task.APIModel)
	if err != nil {
		log.Info("创建任务出错:", err)
		return false
	}
	return true
}

func taskListByPage(start, pageSize int, keyword string) []*Task {
	var searchSql string
	if len(keyword) != 0 {
		searchSql = "where task_status Like CONCAT('%',?,'%')"
	}
	sql := `SELECT
				t.id,
				t.task_name,
				t.task_status,
				t.crawl_mode,
				t.exec_type,
				t.cycle_type,
				t.week_days,
				t.exec_time,
				t.enable_advanced_settings,
				t.enable_filter,
				t.domain_match,
				t.path_match,
				t.crawl_option,
				t.crawl_type,
				t.concurrent_count,
				t.scraping_interval,
				t.global_scraping_depth,
				t.request_rate_limit,
				t.use_proxy_ip_pool,
				t.api_model extraction_model,
				t.api_settings_id_list,
				a.archive_name,
				las.api_type,
				las.name llm_setting_name,
				strftime('%s', t.created_at, '%s') created_at,
				IFNULL(strftime('%s', t.updated_at, '%s'),"") updated_at
			FROM
				task t
				LEFT JOIN archive a ON a.id = t.archive_id
				LEFT JOIN llm_api_settings las ON las.id = t.api_settings_id 
			%s
			GROUP BY
				t.id
			LIMIT ? ,?`
	format := "%Y-%m-%d %H:%M:%S"
	sql = fmt.Sprintf(sql, format, common.GetTimeZone(), format, common.GetTimeZone(), searchSql)
	var list []*Task
	params := make([]any, 0)
	if len(keyword) != 0 {
		params = append(params, keyword)
	}
	params = append(params, start, pageSize)
	err := sqlite.Conn().Select(&list, sql, params...)
	if err != nil {
		log.Info("查询任务表出错:", err)
		return list
	}
	return list
}

// 获取记录总数
func taskCount(status int8, keyword string) int {
	var statusSql string
	if status > 0 {
		statusSql = " AND task_status = ?"
	}
	var searchSql string
	if len(keyword) != 0 {
		searchSql = " AND task_name Like CONCAT('%',?,'%')"
	}
	sql := `SELECT
				count(1)
			FROM
				task 
			where 1 = 1
				%s
				%s;`
	sql = fmt.Sprintf(sql, statusSql, searchSql)
	var num int
	params := make([]any, 0)
	if status > 0 {
		params = append(params, status)
	}
	if len(keyword) != 0 {
		params = append(params, keyword)
	}
	err := sqlite.Conn().Get(&num, sql, params...)
	if err != nil {
		log.Info("查询任务数出错:", err)
		return num
	}
	return num
}

// 抓取全部可执行任务
func allTaskForExec() []*Task {
	sql := `SELECT
				t.id,
				t.task_name,
				t.task_status,
				t.crawl_mode,
				t.exec_type,
				t.cycle_type,
				t.week_days,
				t.exec_time,
				t.enable_advanced_settings,
				t.enable_filter,
				t.domain_match,
				t.path_match,
				t.crawl_option,
				t.crawl_type,
				t.crawl_url,
				t.concurrent_count,
				t.scraping_interval,
				t.global_scraping_depth,
				t.request_rate_limit,
				t.use_proxy_ip_pool,
				t.api_model extraction_model,
				t.api_settings_id_list,
				t.archive_id,
				las.api_type,
				las.api_url,
				las.api_key,
				las.timeout,
				las.request_rate_limit api_request_rate_limit,
				las.name llm_setting_name
			FROM
				task t
			LEFT JOIN llm_api_settings las ON las.id = t.api_settings_id
			WHERE t.task_status = ?
			GROUP BY t.id
			ORDER BY t.id;`
	var list []*Task
	err := sqlite.Conn().Select(&list, sql, 1)
	if err != nil {
		log.Info("查询可执行任务表出错:", err)
		return list
	}
	return list
}

// 插入任务流水
func insertTaskFlow(taskID int, status int) error {
	sql := `INSERT INTO task_flow (task_id, status, created_at) VALUES (?, ?, CURRENT_TIMESTAMP)`
	_, err := sqlite.Conn().Exec(sql, taskID, status)
	if err != nil {
		log.Info("插入任务流水失败:", err)
		return err
	}
	return nil
}

// 修改任务状态
func updateTaskStatus(taskID int, status bool) error {
	sql := `UPDATE task SET task_status = ? WHERE id = ?`
	_, err := sqlite.Conn().Exec(sql, status, taskID)
	if err != nil {
		log.Info("更新任务状态失败:", err)
		return err
	}
	return nil
}
