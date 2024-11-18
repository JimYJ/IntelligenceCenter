package utils

import (
	"net/url"
	"strings"
)

func CheckURL(str string) bool {
	parsedURL, err := url.Parse(str)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func GetHost(str *string) []string {
	if str == nil || len(*str) == 0 {
		return nil
	}
	list := make([]string, 0)
	urlList := strings.Split(*str, "\n")
	for _, item := range urlList {
		list = append(list, strings.TrimSpace(item))
	}
	return list
}
