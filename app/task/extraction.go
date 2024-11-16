package task

import (
	"IntelligenceCenter/app/archive"
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/service/log"
	"encoding/json"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/yaml.v2"
)

type ExtractionBody struct {
	URL     string
	Content string
	DocID   int
}

var (
	extractionChan = make(chan *ExtractionBody, 65535)
)

func ListenMatch() {
	for {
		body := <-extractionChan
		Match(body.URL, body.Content, body.DocID)
	}
}

func Match(url, body string, docID int) {
	rules := getRules()
	domain := GetDomainFromURL(url)
	matchContent := make(map[string]string)
	var extractionMode uint8 = 2
	for _, rule := range rules {
		if matchPatterns(rule.MatchDomain, domain) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
			if err != nil {
				log.Info("解析HTML错误:", err)
				return
			}
			for k, selection := range rule.Content {
				doc.Find(selection).Each(func(i int, s *goquery.Selection) {
					// 在这里处理匹配到的内容
					log.Info("匹配到的内容:", s.Text())
					matchContent[k] = s.Text()
				})
			}
			extractionMode = 1
			updateContent(matchContent, extractionMode, docID, -1, "-")
			break
		}
	}
	if extractionMode == 2 {
		// TODO 智能匹配
	}
}

func updateContent(matchContent map[string]string, extractionMode uint8, docID, apiKeyID int, extractionModel string) {
	var matchContentJSON []byte
	var err error
	if len(matchContent) > 0 {
		// 保存到数据库
		matchContentJSON, err = json.Marshal(matchContent)
		if err != nil {
			log.Info("转换matchContent为JSON字符串时出错:", err)
			return
		}
		matchContentJSONStr := string(matchContentJSON)
		archive.UpdateDocByExtraction(docID, apiKeyID, extractionMode, &matchContentJSONStr, extractionModel)
	} else {
		archive.UpdateDocByExtraction(docID, apiKeyID, extractionMode, nil, extractionModel)
	}
}

// 匹配多个模式
func matchPatterns(patterns []string, domain string) bool {
	for _, pattern := range patterns {
		if matchDomain(pattern, domain) {
			return true
		}
	}
	return false
}

// 匹配函数
func matchDomain(pattern, domain string) bool {
	// 如果模式是以*开头并且长度大于1，则去掉*进行匹配
	if strings.HasPrefix(pattern, "*") && len(pattern) > 1 {
		pattern = pattern[2:]
		if strings.HasSuffix(domain, pattern) && domain != pattern {
			return true
		}
	}
	// 将模式和域名分割成各个部分
	patternParts := strings.Split(pattern, ".")
	domainParts := strings.Split(domain, ".")
	// 如果模式部分的数量大于域名部分的数量，则不匹配
	if len(patternParts) > len(domainParts) {
		return false
	}
	// 从右到左进行匹配
	for i := 0; i < len(patternParts); i++ {
		patternPart := patternParts[len(patternParts)-1-i]
		domainPart := domainParts[len(domainParts)-1-i]
		// 如果模式部分不是*且不等于域名部分，则不匹配
		if patternPart != "*" && patternPart != domainPart {
			return false
		}
	}
	// 如果所有部分都匹配，则返回true
	return true
}

// 添加以下函数
func GetDomainFromURL(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Info("解析URL错误:", err)
		return ""
	}
	return parsedURL.Host
}

func getRules() []*extractionRules {
	list, err := utils.FindFilesBySuffix("./extraction-rules", "yml", "yaml")
	log.Info(list, err)
	rulesList := make([]*extractionRules, 0)
	for _, item := range list {
		rules := getRulesFiles(item)
		if rules != nil {
			rulesList = append(rulesList, rules)
		}
	}
	return rulesList
}

func getRulesFiles(path string) *extractionRules {
	rules := &extractionRules{}
	file, err := os.ReadFile(path)
	if err != nil {
		log.Info("读取文件错误:", err)
		return nil
	}
	err = yaml.Unmarshal(file, rules)
	if err != nil {
		log.Info("解析文件错误:", err)
		return nil
	}
	return rules
}
