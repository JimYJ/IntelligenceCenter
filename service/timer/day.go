package timer

import (
	"IntelligenceCenter/app/task"
	"IntelligenceCenter/service/log"
	"time"
)

var (
	nextDay time.Time
	dt      *time.Timer
)

// DayTaskFor2AM 每天凌晨0点执行一次
// 本函数不得在主线程执行
func DayTaskFor0AM() {
	for {
		nextDay = time.Now().Local().Add(time.Hour * 24)
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
		dt = time.NewTimer(nextDay.Sub(time.Now().Local()))
		<-dt.C
		log.Info("=========开始每日00:00AM任务==========")
		task.Scan()
		log.Info("=========结束每日00:00AM任务==========")
	}
}
