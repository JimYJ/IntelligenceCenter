package task

import (
	"IntelligenceCenter/app/option"
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/service/log"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
)

var (
	RetryChan        = make(chan *RetryBody, 65535)
	taskRetryCounter = make(map[int]*RetryCounter)
)

// 创建爬虫
func (task *Task) CreateCrawler() *colly.Collector {
	task.setGroupOption()
	c := colly.NewCollector(
		colly.Async(),
		colly.TraceHTTP(),
		colly.IgnoreRobotsTxt(),
	)
	c.SetRequestTimeout(requestTimeout)
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
	// c.OnHTML("div.entry-content.th-content", func(e *colly.HTMLElement) {
	// 	// 获取 div 的内容
	// 	// content := e.Text
	// 	e.Request.URL.String()
	// })
	c.OnRequest(func(r *colly.Request) {
		if task.UseProxyIPPool != nil && !*task.UseProxyIPPool {
			ips, err := utils.GeneratePublicIPs(1, 1)
			if err != nil {
				r.Headers.Set("X-Real-IP", ips[0].String())
				r.Headers.Set("x-forwarded-for", ips[0].String())
			}
		}
		if ref := r.Ctx.Get("_referer"); ref != "" {
			r.Headers.Set("Referer", ref)
		}
	})
	// 处理响应
	c.OnResponse(task.OnResponse)
	// 设置代理
	task.setProxy(c)
	// 设置重试
	taskRetryCounter[task.ID] = NewRetryCounter(3) // TODO 可以改成界面配置
	c.OnError(OnError(c, task))
	taskCrawler[task.ID] = c
	// c.Visit("https://misif.org.my/directory/")
	return c
}

func (task *Task) OnResponse(r *colly.Response) {
	r.Ctx.Put("_referer", r.Request.URL.String())

	// 精准解析
	html, err := goquery.NewDocumentFromReader(r.Request.Body)
	if err != nil {
		log.Info("解析HTML错误:", err)
	}
	// Find the review items
	html.Find(".left-content article .post-title").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("a").Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
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

// 设置代理
func (task Task) setProxy(c *colly.Collector) {
	if task.UseProxyIPPool != nil && *task.UseProxyIPPool {
		ips := getProxyIP()
		if len(ips) > 0 {
			rp, err := proxy.RoundRobinProxySwitcher(ips...)
			if err != nil {
				log.Info("设置代理错误:", err)
				return
			}
			c.SetProxyFunc(rp)
		}
	}
}

// NewRetryBody 创建一个新的RetryInfo实例
func NewRetryBody(url string, task *Task, retry uint) *RetryBody {
	return &RetryBody{
		URL:   url,
		Task:  task,
		Count: retry,
	}
}

// OnError 处理请求错误，并决定是否通过通道传递重试信息
func OnError(c *colly.Collector, task *Task) colly.ErrorCallback {
	return func(response *colly.Response, err error) {
		if err != nil {
			log.Info("请求错误:", err, response.Request.URL.String())
			retryCounter, ok := taskRetryCounter[task.ID]
			var retryCount uint = 1
			if ok && retryCounter != nil {
				// 超过重试次数就退出
				if retryCounter.Limit(response.Request.URL.String()) {
					return
				}
				retryCount = retryCounter.Add(response.Request.URL.String())
			}
			RetryChan <- NewRetryBody(response.Request.URL.String(), task, retryCount)
		}
	}
}

// RetryCounter 结构体用于记录每个URL的重试次数
type RetryCounter struct {
	URLCount map[string]uint
	Max      uint
	Lock     sync.RWMutex
}

func NewRetryCounter(max uint) *RetryCounter {
	return &RetryCounter{
		URLCount: make(map[string]uint),
		Max:      max,
	}
}

func (r *RetryCounter) Add(url string) uint {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	return r.URLCount[url] + 1
}

func (r *RetryCounter) Limit(url string) bool {
	r.Lock.RLock()
	defer r.Lock.RUnlock()
	return r.URLCount[url]+1 > r.Max
}

func Retry() {
	for {
		retryInfo, ok := <-RetryChan
		if ok && taskCrawler[retryInfo.Task.ID] != nil {
			taskCrawler[retryInfo.Task.ID].Visit(retryInfo.URL)
		}
	}
}
