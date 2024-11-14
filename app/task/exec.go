package task

import (
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
		if crawler, ok := taskCrawler[task.ID]; !ok {
			taskCrawler[task.ID] = task.CreateCrawler()
			task.Crawler = taskCrawler[task.ID]
		} else {
			task.Crawler = crawler
		}

	} else if task.CrawlMode == 2 {
		// TODO AI Agent

	}
}
