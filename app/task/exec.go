package task

import (
	"IntelligenceCenter/service/log"
	"time"

	"github.com/gocolly/colly/v2"
)

var (
	taskCrawler    = make(map[int]*colly.Collector)
	requestTimeout = 30 * time.Second
)

// 任务执行器
func (task *Task) Exec() {
	if task.CrawlMode == 1 {
		log.Info("开始执行地址抓取任务:", task.TaskName, task.ID)
		if crawler, ok := taskCrawler[task.ID]; !ok {
			taskCrawler[task.ID] = task.CreateCrawler()
			task.Crawler = taskCrawler[task.ID]
		} else {
			task.Crawler = crawler
		}
	} else if task.CrawlMode == 2 {
		log.Info("开始执行智能抓取任务:", task.TaskName, task.ID)
		// TODO AI Agent

	}
}
