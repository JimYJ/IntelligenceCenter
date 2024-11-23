package task

import (
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/service/log"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

// 任务流水状态枚举
const (
	_                            int = iota
	TaskFlowStatusCreated            // 创建触发
	TaskFlowStatusScheduled          // 定时触发
	TaskFlowStatusCompleted          // 执行结束
	TaskFlowStatusManuallyClosed     // 手动关闭
	TaskFlowStatusManuallyOpened     // 手动开启
)

var (
	taskPool = sync.Pool{
		New: func() interface{} {
			return &Task{}
		},
	}
)

type Task struct {
	ID                     int              `json:"id" db:"id"`                                             // 主键
	ArchiveOption          uint8            `json:"archive_option,string" db:"-"`                           // 1新建档案 2选择档案
	ArchiveID              int              `json:"archive_id" db:"archive_id"`                             // 指定归档的档案ID
	ArchiveName            string           `json:"archive_name" db:"archive_name"`                         // 档案名称
	TaskName               string           `json:"task_name" db:"task_name"`                               // 任务名称
	CrawlMode              uint8            `json:"crawl_mode,string" db:"crawl_mode"`                      // 抓取模式 1地址抓取 2描述搜索抓取
	CrawlURL               string           `json:"crawl_url" db:"crawl_url"`                               // 抓取地址或抓取描述
	ExecType               uint8            `json:"exec_type,string" db:"exec_type"`                        // 执行类型 1-立即执行 2-周期循环
	CycleType              uint8            `json:"cycle_type,string" db:"cycle_type"`                      // 周期类型 1-每日 2-每周
	WeekDays               []string         `json:"week_days" db:"-"`                                       // 指定周几执行，可多选
	WeekDaysStr            string           `json:"-" db:"week_days"`                                       // 指定周几执行，可多选，英文逗号隔开
	ExecTime               string           `json:"exec_time" db:"exec_time"`                               // 执行时间
	EnableAdvancedSettings bool             `json:"enable_advanced_settings" db:"enable_advanced_settings"` // 启用进阶设置
	TaskStatus             bool             `json:"task_status" db:"task_status"`                           // 任务状态
	EnableFilter           *bool            `json:"enable_filter" db:"enable_filter"`                       // 启用匹配过滤器
	DomainMatch            *string          `json:"domain_match" db:"domain_match"`                         // 域名匹配过滤器 为空则不生效
	PathMatch              *string          `json:"path_match" db:"path_match"`                             // 路径匹配过滤器 为空则不生效
	CrawlOption            *bool            `json:"crawl_option" db:"crawl_option"`                         // 抓取器设置,是否使用全局设置
	CrawlType              *uint8           `json:"crawl_type,string" db:"crawl_type"`                      // 抓取器选择 1 内置爬虫 2 headless浏览器 3 firecrawl
	ConcurrentCount        *int             `json:"concurrent_count" db:"concurrent_count"`                 // 并发数
	ScrapingInterval       *int             `json:"scraping_interval" db:"scraping_interval"`               // 抓取间隔(秒)
	GlobalScrapingDepth    *int             `json:"global_scraping_depth" db:"global_scraping_depth"`       // 抓取深度
	RequestRateLimit       *int             `json:"request_rate_limit" db:"request_rate_limit"`             // 每秒请求上限
	UseProxyIPPool         *bool            `json:"use_proxy_ip_pool" db:"use_proxy_ip_pool"`               // 使用代理IP池
	APISettingsIDList      []int            `json:"api_settings_id" db:"-"`                                 // API设置表ID
	APISettingsID          int              `json:"-" db:"api_settings_id"`                                 // API设置表ID
	APISettingsIDStr       string           `json:"-" db:"api_settings_id_list"`                            // API设置表ID 字符串 英文逗号隔开
	APIModel               *string          `json:"extraction_model" db:"extraction_model"`                 // API指定LLM模型
	ApiType                uint8            `json:"api_type" db:"api_type"`                                 // API类型 1-OpenAI API Api 2-Ollama
	LLMSettingName         string           `json:"llm_setting_name" db:"llm_setting_name"`                 // LLM设置名称
	ApiURL                 string           `json:"api_url" db:"api_url"`                                   // API 地址
	ApiKey                 string           `json:"api_key" db:"api_key"`                                   // API 密钥
	Timeout                int              `json:"timeout" db:"timeout"`                                   // 超时设置(秒),默认30秒
	APIRequestRateLimit    int              `json:"api_request_rate_limit" db:"api_request_rate_limit"`     // 每秒请求上限
	ExecTimeSec            int64            `db:"-" json:"-"`                                               // 执行时间戳
	CreatedAt              string           `json:"created_at" db:"created_at"`                             // 更新时间
	UpdatedAt              string           `json:"updated_at" db:"updated_at"`                             // 创建时间
	Crawler                *colly.Collector `json:"-" db:"-"`                                               // 爬虫
}

// Free 重置任务的状态
func (t *Task) Free() {
	t.ID = 0
	t.ArchiveOption = 0
	t.ArchiveID = 0
	t.ArchiveName = ""
	t.TaskName = ""
	t.CrawlMode = 0
	t.CrawlURL = ""
	t.ExecType = 0
	t.CycleType = 0
	t.WeekDays = nil
	t.WeekDaysStr = ""
	t.ExecTime = ""
	t.EnableAdvancedSettings = false
	t.TaskStatus = false
	t.EnableFilter = nil
	t.DomainMatch = nil
	t.PathMatch = nil
	t.CrawlOption = nil
	t.CrawlType = nil
	t.ConcurrentCount = nil
	t.ScrapingInterval = nil
	t.GlobalScrapingDepth = nil
	t.RequestRateLimit = nil
	t.UseProxyIPPool = nil
	t.APISettingsIDList = nil
	t.APISettingsID = 0
	t.APISettingsIDStr = ""
	t.APIModel = nil
	t.ApiType = 0
	t.LLMSettingName = ""
	t.ApiURL = ""
	t.ApiKey = ""
	t.Timeout = 0
	t.APIRequestRateLimit = 0
	t.CreatedAt = ""
	t.UpdatedAt = ""
	t.ExecTimeSec = 0
	taskPool.Put(t)
}

func (task *Task) CheckExecTime() bool {
	t := time.Now().Local()
	if task.ExecType == 1 { // 立即执行
		task.ExecTime = t.Format(time.DateTime)
		task.ExecTimeSec = t.Unix()
		return true
	}
	if len(task.ExecTime) == 0 {
		log.Info("脏数据!执行时间错误:", task.ExecTime, task.TaskName)
		return false
	}
	if task.ExecType == 2 {
		if task.CycleType == 2 {
			// 每周执行
			list := strings.Split(task.WeekDaysStr, ",")
			list2, err := utils.ConvertStringsToInts(list)
			if err != nil {
				log.Info("脏数据!每周日期错误:", task.WeekDaysStr, task.TaskName)
				return false
			}
			var confirm bool
			for _, item := range list2 {
				if time.Weekday(item) == t.Weekday() {
					confirm = true
				}
			}
			if !confirm {
				log.Info("不匹配指定每周日期，跳过:", task.WeekDaysStr, t.Weekday(), task.TaskName)
				return false
			}
		}
		timerDate := utils.JoinString(t.Format(time.DateOnly), " ", task.ExecTime)
		timer, err := time.ParseInLocation(time.DateTime, timerDate, time.Local)
		if err != nil {
			log.Info("脏数据!执行时间错误:", task.ExecTime, timerDate, task.TaskName)
			return false
		}
		task.ExecTime = timerDate
		task.ExecTimeSec = timer.Local().Unix()
		if task.ExecTimeSec < t.Unix() {
			log.Info("执行时间已经超过当前时间，之后再执行:", task.ExecTime, timerDate, task.TaskName)
			return false
		}
		return true
	}
	return false
}

// 实现排序
type tasklist []*Task

func (list tasklist) Len() int           { return len(list) }
func (list tasklist) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }
func (list tasklist) Less(i, j int) bool { return list[i].ExecTimeSec < list[j].ExecTimeSec }

type TaskData struct {
	ArchiveCount        int `json:"archive_count" db:"-"`          // 档案数量
	ArchiveDocsCount    int `json:"archive_docs_count" db:"-"`     // 档案文档数量
	ArchiveDocsResCount int `json:"archive_docs_res_count" db:"-"` // 档案文档资源数
	TaskCount           int `json:"task_count" db:"-"`             // 关联任务总数
	ActiveTaskCount     int `json:"active_task_count" db:"-"`      // 关联活跃任务总数
}

type RetryBody struct {
	Task  *Task  // 任务ID
	URL   string // 重试地址
	Count uint   // 重试次数
}

type extractionRules struct {
	MatchDomain []string          `yaml:"match-domain"`
	Content     map[string]string `yaml:"content"`
}
