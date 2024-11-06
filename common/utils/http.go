package utils

import (
	"net/url"
)

func CheckURL(str string) bool {
	parsedURL, err := url.Parse(str)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}
