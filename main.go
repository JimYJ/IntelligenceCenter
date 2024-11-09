package main

import (
	"IntelligenceCenter/app/db"
	"IntelligenceCenter/app/task"
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/router"
	"embed"
	"fmt"
	"regexp"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// 检查初始化目录
	needDir = []string{"logs", "extraction-rules", "database", "proxy-ip"}
)

//go:embed static/dist
var static embed.FS

func main() {
	initDir()
	db.CheckDatabase()
	go task.ListenNewTask()
	go router.Web(static)
	router.Api()
}

// 初始化必要目录
func initDir() {
	for _, item := range needDir {
		utils.CreateDir(item)
	}
}

func main2() {
	c := colly.NewCollector(
		colly.AllowedDomains("misif.org.my"),
		colly.URLFilters(
			regexp.MustCompile("https://misif.org.my/directory/*"),
		),
	)
	extensions.RandomUserAgent(c)
	c.Limit(&colly.LimitRule{
		DomainRegexp: `misif\.org\.my`,
		RandomDelay:  500 * time.Millisecond,
		Parallelism:  1,
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// log.Info("发现链接: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// 在访问网页时调用的处理函数
	c.OnHTML("div.entry-content.th-content", func(e *colly.HTMLElement) {
		// 获取 div 的内容
		content := e.Text
		fmt.Println("内容:", content)
	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {

	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("请求URL:", r.Request.URL, "响应失败:", string(r.Body), "\nError:", err)
	})

	c.Visit("https://misif.org.my/directory/")
}
