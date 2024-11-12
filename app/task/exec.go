package task

import (
	"IntelligenceCenter/app/option"
	"IntelligenceCenter/common/utils"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
)

var (
	taskCrawler = make(map[int]*colly.Collector)
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

// 创建爬虫
func (task Task) CreateCrawler() *colly.Collector {
	task.setGroupOption()
	c := colly.NewCollector()
	c.Async = true
	extensions.RandomUserAgent(c)
	hostList := utils.GetHost(task.DomainMatch)
	if task.EnableFilter != nil && *task.EnableFilter {
		// 域名过滤
		if len(hostList) != 0 {
			c.AllowedDomains = hostList
		}
		// 路径过滤
		if task.PathMatch != nil && len(*task.PathMatch) != 0 {
			pathList := strings.Split(*task.PathMatch, "\n")
			for _, item := range pathList {
				c.URLFilters = append(c.URLFilters, regexp.MustCompile(item))
			}
		}
	}
	if task.GlobalScrapingDepth != nil && *task.GlobalScrapingDepth > 0 {
		c.MaxDepth = *task.GlobalScrapingDepth
	}
	// 并发和延迟
	limitRule := &colly.LimitRule{
		DomainRegexp: `*`,
		Parallelism:  *task.ConcurrentCount,
	}
	if task.ScrapingInterval != nil && *task.ScrapingInterval > 0 {
		limitRule.RandomDelay = time.Duration(*task.ScrapingInterval) * time.Second
	}
	c.Limit(limitRule)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// c.Visit(e.Request.AbsoluteURL(link))
		e.Request.Visit(e.Request.AbsoluteURL(link))
	})

	// 精准匹配
	c.OnHTML("div.entry-content.th-content", func(e *colly.HTMLElement) {
		// 获取 div 的内容
		// content := e.Text
		e.Request.URL.String()
	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
		r.Headers.Set("Referer", "ip") // ip
		if ref := r.Ctx.Get("_referer"); ref != "" {
			r.Headers.Set("Referer", ref)
		}
	})

	c.OnResponse(func(r *colly.Response) {

	})

	// 设置代理
	task.setProxy(c)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("请求URL:", r.Request.URL, "响应失败:", string(r.Body), "\nError:", err)
	})

	// c.Visit("https://misif.org.my/directory/")
	return c
}

// 判断及获取全局设置
func (task *Task) setGroupOption() {
	if task.CrawlOption != nil && *task.CrawlOption {
		op := option.GetOption()
		*task.RequestRateLimit = op.RequestRateLimit
		*task.GlobalScrapingDepth = op.GlobalScrapingDepth
		*task.ConcurrentCount = op.ConcurrentCount
		*task.ScrapingInterval = op.ScrapingInterval
		*task.UseProxyIPPool = op.UseProxyIPPool
		*task.CrawlType = op.CrawlType
	}
}

func (task Task) setProxy(c *colly.Collector) {
	if task.UseProxyIPPool != nil && *task.UseProxyIPPool {
		rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1337", "socks5://127.0.0.1:1338")
		if err != nil {
			log.Fatal(err)
		}
		c.SetProxyFunc(rp)
	}
}
