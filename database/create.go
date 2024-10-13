package database

import (
	"IntelligenceCenter/common/sqlite"
	"IntelligenceCenter/service/log"
	"database/sql"
)

var (
	createSqlList = []string{optionTableSql, llmSettingTableSql, archiveTableSql}

	checkTableSql = ""

	optionTableSql = `CREATE TABLE "option" (
						"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT ,    -- 主键
						"crawl_type" integer NOT NULL DEFAULT 1,             -- 抓取器选择 1-内置爬虫 2-headless浏览器 3-firecrawl
						"concurrent_count" integer NOT NULL,                 -- 并发数
						"scraping_interval" integer NOT NULL,                -- 抓取间隔(秒)
						"global_scraping_depth" integer NOT NULL,            -- 抓取深度
						"request_rate_limit" integer NOT NULL,               -- 每秒请求上限
						"use_proxy_ip_pool" boolean NOT NULL,                -- 使用代理IP池
						"use_global_concurrency_pool" boolean NOT NULL       -- 使用全局并发池
					);`

	llmSettingTableSql = `CREATE TABLE llm_api_settings (
							id INTEGER PRIMARY KEY AUTOINCREMENT,           -- 主键
							name varchar ( 128 ) NOT NULL,                  -- 配置名称
							api_type integer NOT NULL,                      -- API类型 1-OpenAI Api 2-Ollama 3-Siliconflow API
							api_url TEXT NOT NULL,                          -- API 地址
							api_key TEXT NOT NULL,                          -- API 密钥
							timeout INTEGER NOT NULL DEFAULT 30,            -- 超时设置(秒),默认30秒
							request_rate_limit INTEGER NOT NULL,            -- 每秒请求上限
							use_proxy_pool BOOLEAN NOT NULL DEFAULT 0,      -- 使用代理 IP 池
							remark TEXT                                     -- 描述信息
						);`
	archiveTableSql = `CREATE TABLE "archive" (
							"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,           -- 主键
							"archive_name" varchar(128) NOT NULL,                      -- 档案名称
							"file_count" integer NOT NULL DEFAULT 0,                   -- 档案文件数
							"extraction_mode" integer NOT NULL,                        -- 提取模式 1-精准匹配 2-智能匹配
							"api_key_id" integer NOT NULL,                             -- llm_api_settings 表ID
							"extraction_model" varchar(128) NOT NULL,                  -- 提取模型
							"created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
							"updated_at" datetime                                      -- 更新时间
						);
						  `
)

// 初始化数据库
func createDatabase() {
	for _, item := range createSqlList {
		_, err := sqlite.Conn().Exec(item)
		if err != nil {
			log.Info("创建初始化表失败", err)
			return
		}
	}
}

// 检查数据库
func CheckDatabase() {
	var tableExists bool
	err := sqlite.Conn().Get(&tableExists, checkTableSql)
	if err == sql.ErrNoRows {
		log.Info("找不到数据文件，准备创建...")
		createDatabase()
		return
	} else if err != nil {
		log.Info(err)
		return
	} else {
		log.Info("数据库已经存在。", tableExists)
		return
	}
}
