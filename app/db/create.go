package db

import (
	"IntelligenceCenter/common/sqlite"
	"IntelligenceCenter/service/log"
	"database/sql"
)

var (
	createSqlList = []string{optionTableSql, llmSettingTableSql, archiveTableSql, archiveDocsTableSql, docResourceTableSql, taskTableSql}

	checkTableSql = ""

	optionTableSql = `CREATE TABLE "option" (
						"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT ,    -- 主键
						"crawl_type" integer NOT NULL DEFAULT 1,             -- 抓取器选择 1-内置爬虫 2-headless浏览器 3-firecrawl
						"concurrent_count" integer NOT NULL,                 -- 并发数
						"scraping_interval" integer NOT NULL,                -- 抓取间隔(秒)
						"global_scraping_depth" integer NOT NULL,            -- 抓取深度
						"request_rate_limit" integer NOT NULL,               -- 每秒请求上限
						"use_proxy_ip_pool" boolean NOT NULL DEFAULT 0,      -- 使用代理IP池
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
							use_proxy_pool BOOLEAN NOT NULL DEFAULT 0,      -- 使用代理 IP 池 0否1是
							remark TEXT                                     -- 描述信息
						);`
	// "file_count" integer NOT NULL DEFAULT 0,                -- 档案文件数
	// "extraction_mode" integer NOT NULL,                     -- 提取模式 1-精准匹配 2-智能匹配
	// "api_key_id" integer NOT NULL,                          -- llm_api_settings 表ID
	// "extraction_model" varchar(128) NOT NULL,                  -- 提取模型
	archiveTableSql = `CREATE TABLE "archive" (
							"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,           -- 主键
							"archive_name" varchar(128) NOT NULL,                      -- 档案名称
							"created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
							"updated_at" datetime                                      -- 更新时间
						);`

	archiveDocsTableSql = `CREATE TABLE "archive_docs" (
							"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,           -- 主键
							"doc_name" varchar(128) NOT NULL,                          -- 文档名称
							"task_id" integer NOT NULL,                                -- 任务ID
							"archive_id" INTEGER NOT NULL,                             -- 所属的档案ID
							"origin_content" text,                                     -- 文档原始内容
							"extraction_content" text,                                 -- 提取后内容
							"translate_content" text,                                  -- 翻译后内容
							"extraction_mode" integer,                                 -- 提取模式 1-精准匹配 2-智能匹配
							"api_key_id" integer,                                      -- llm_api_settings 表ID
							"extraction_model" varchar(128),                           -- 提取模型
							"is_extracted" BOOLEAN NOT NULL DEFAULT 0,                 -- 是否被提取 0否1是
							"is_translated" BOOLEAN NOT NULL DEFAULT 0,                -- 是否被翻译 0否1是
							"src_url" text,                                            -- 来源网址/来源文档地址
							"created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
							"updated_at" datetime                                      -- 更新时间
						);`
	docResourceTableSql = `CREATE TABLE "doc_resource" (
							"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,           -- 主键
							"doc_id" INTEGER NOT NULL,                                 -- 文档ID
							"archive_id" INTEGER NOT NULL,                             -- 所属的档案ID
							"resource_type" integer NOT NULL,                          -- 资源类型 1-图片: 2-PDF 3-docs 4-PPT 5-Excel 6-magnet 7-telegram邀请链接
							"resource_path" text NOT NULL,                             -- 资源路径
							"resource_status" integer NOT NULL,                        -- 资源状态 1-未下载 2-已下载
							"resource_size" integer NOT NULL,                          -- 资源大小(字节数)
							"created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP   -- 创建时间
						);`
	taskTableSql = `CREATE TABLE "task" (
						"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,           -- 主键
						"archive_id" INTEGER NOT NULL,                             -- 指定归档的档案ID
						"task_name" varchar(128) NOT NULL,                         -- 任务名称
						"crawl_mode" INTEGER NOT NULL DEFAULT 1,                   -- 抓取模式 1地址抓取 2描述搜索抓取
						"crawl_url" text NOT NULL,                                 -- 抓取地址，多个地址换行分割
						"exec_type" INTEGER NOT NULL DEFAULT 1,                    -- 执行类型 1-立即执行 2-周期循环
						"cycle_type" INTEGER NOT NULL DEFAULT 1,                   -- 周期类型 1-每日 2-每周
						"week_days" varchar(20),                                   -- 指定周几执行，可多选，英文逗号隔开
						"exec_time" time,                                          -- 执行时间
						"enable_advanced_settings" BOOLEAN NOT NULL DEFAULT 0,     -- 启用进阶设置 0关闭 1启用
						"task_status" BOOLEAN NOT NULL DEFAULT 1,                  -- 任务状态 0关闭 1启用
						"enable_filter" BOOLEAN NOT NULL DEFAULT 0,                -- 启用匹配过滤器 0关闭 1启用
						"domain_match" text,                                       -- 域名匹配过滤器 为空则不生效
						"path_match" text,                                         -- 路径匹配过滤器 为空则不生效
						"crawl_option" BOOLEAN NOT NULL DEFAULT 1,                 -- 抓取器设置 0自定义 1全局
						"crawl_type" INTEGER NOT NULL DEFAULT 1,                   -- 抓取器选择 1 内置爬虫 2 headless浏览器 3 firecrawl
						"concurrent_count" INTEGER,                                -- 并发数
						"scraping_interval" INTEGER,                               -- 抓取间隔(秒)
						"global_scraping_depth" INTEGER,                           -- 抓取深度
						"request_rate_limit" INTEGER,                              -- 每秒请求上限
						"use_proxy_ip_pool" BOOLEAN DEFAULT 0,                     -- 使用代理IP池
						"api_settings_id" INTEGER,                                 -- API设置表ID
						"api_settings_id_list" varchar(128),                       -- API设置表ID(前端用)
						"api_model" varchar(128),                                  -- API指定LLM模型
						"created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
						"updated_at" datetime                                      -- 更新时间
					);`
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
