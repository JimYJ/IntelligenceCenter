package utils

import (
	"IntelligenceCenter/service/log"
	"net/url"
	"strings"
)

func CheckURL(str string) bool {
	parsedURL, err := url.Parse(str)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func GetHost(str *string) []string {
	list := make([]string, 0)
	if str == nil {
		return list
	}
	urlList := strings.Split(*str, "\n")
	for _, item := range urlList {
		parsedURL, err := url.Parse(item)
		if err != nil {
			log.Info("解析网址错误:", err)
			continue
		}
		list = append(list, parsedURL.Scheme+parsedURL.Host)
	}
	return list
}
