package task

import (
	"IntelligenceCenter/service/log"
	"strings"
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
			log.Info("创建爬虫成功:", task.TaskName, task.ID)
		} else {
			task.Crawler = crawler
			log.Info("获取爬虫成功:", task.TaskName, task.ID)
		}
		list := strings.Split(task.CrawlURL, "\n")
		for _, item := range list {
			log.Info("开始抓取:", item, task.TaskName, task.ID)
			err := task.Crawler.Visit(item)
			if err != nil {
				log.Info("抓取失败:", err, item, task.TaskName, task.ID)
			} else {
				log.Info("抓取成功:", item, task.TaskName, task.ID)
			}
		}
		task.Crawler.Wait()
		log.Info("所有子链接抓取任务完成:", task.TaskName, task.ID)
		insertTaskFlow(task.ID, TaskFlowStatusCompleted)

		if task.ExecType == 1 {
			updateTaskStatus(task.ID, false)
		}
	} else if task.CrawlMode == 2 {
		log.Info("开始执行智能抓取任务:", task.TaskName, task.ID)
		// TODO AI Agent

	}
}
