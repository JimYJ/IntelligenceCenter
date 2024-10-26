package common

import (
	"fmt"
	"time"
)

var (
	now       = time.Now()
	_, offset = now.Zone() // 获取时区名称和偏移量
)

func GetTimeZone() string {
	if offset == 0 {
		return ""
	} else if offset > 0 {
		return fmt.Sprintf("+%d hours", offset/3600)
	} else {
		return fmt.Sprintf("-%d hours", offset/3600)
	}
}
