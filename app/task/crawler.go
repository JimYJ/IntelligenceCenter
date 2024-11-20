package task

import (
	"IntelligenceCenter/app/archive"
	"IntelligenceCenter/app/option"
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/service/log"
	"regexp"
	"strings"
	"sync"
	"time"

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
			log.Info(hostList, len(hostList))
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
		DomainRegexp: `.+`,
		Parallelism:  *task.ConcurrentCount,
	}
	if task.ScrapingInterval != nil && *task.ScrapingInterval > 0 {
		limitRule.RandomDelay = time.Duration(*task.ScrapingInterval) * time.Second
	}
	log.Info("并发:", limitRule.Parallelism, "随机延迟:", limitRule.RandomDelay)
	if err := c.Limit(limitRule); err != nil {
		log.Info("设置并发错误:", err)
	}
	var docID int64
	c.OnHTML("title", func(e *colly.HTMLElement) {
		docID = archive.CreateDoc(task.ID, task.ArchiveID, e.Text, string(e.Response.Body), e.Request.URL.String())
		if docID != -1 {
			extractionChan <- &ExtractionBody{
				URL:     e.Request.URL.String(),
				Content: string(e.Response.Body),
				DocID:   int(docID),
			}
		}
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = e.Request.AbsoluteURL(link)
		if srcType, _ := checkSrcLink(link); srcType != -1 {
			archive.DocResourceChan <- &archive.DocResource{
				DocID:          docID,
				ArchiveID:      task.ArchiveID,
				ResourceType:   srcType,
				ResourcePath:   link,
				ResourceStatus: 1,
				ResourceSize:   0,
			}
		} else {
			// log.Info("准备访问链接", link)
			e.Request.Visit(link)
			// c.Visit(e.Request.AbsoluteURL(link))
		}
	})
	c.OnRequest(func(r *colly.Request) {
		if task.UseProxyIPPool != nil && !*task.UseProxyIPPool {
			ips, err := utils.GeneratePublicIPs(1, 1)
			if err == nil {
				r.Headers.Set("X-Real-IP", ips[0].String())
				r.Headers.Set("x-forwarded-for", ips[0].String())
			}
		}
		if ref := r.Ctx.Get("_referer"); ref != "" {
			r.Headers.Set("Referer", ref)
		}
	})
	// 处理响应
	c.OnResponse(func(r *colly.Response) {
		r.Ctx.Put("_referer", r.Request.URL.String())
	})
	// 设置代理
	task.setProxy(c)
	// 设置重试
	taskRetryCounter[task.ID] = NewRetryCounter(3) // TODO 可以改成界面配置
	c.OnError(OnError(c, task))
	taskCrawler[task.ID] = c
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

func checkSrcLink(link string) (int8, string) {
	extensionsMap := map[string]int8{
		".jpg":  1,
		".jpeg": 1,
		".pdf":  2,
		".doc":  3,
		".docx": 3,
		".wps":  3,
		".ppt":  4,
		".pptx": 4,
		".xls":  5,
		".xlsx": 5,
		".zip":  8,
		".rar":  8,
		".7z":   8,
		".gz":   8,
		".psd":  9,
		".cdr":  9,
	}
	if strings.HasPrefix(link, "magnet:") {
		return 6, link
	}
	if strings.HasPrefix(link, "https://t.me/") {
		return 7, link
	}
	for ext, enum := range extensionsMap {
		if strings.HasSuffix(strings.ToLower(link), ext) {
			return enum, strings.ToLower(link)
		}
	}
	return -1, ""
}
