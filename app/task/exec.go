package task

import (
	"IntelligenceCenter/service/log"
	"sort"
)

var (
	NewTaskChan = make(chan bool, chanSize)
	TaskChan    = make(chan *Task, chanSize)
	cancelChan  = make(chan bool, 1)
	chanSize    = 1024
)

// 新任务监听
func ListenNewTask() {
	for {
		<-NewTaskChan
		Scan()
		if len(cancelChan) == 0 {
			cancelChan <- true
		}
	}
}

// 扫描任务
func Scan() {
	list := allTaskForExec()
	execList := make(tasklist, 0)
	// 检查执行时间
	for _, item := range list {
		log.Info("读取任务:", item.TaskName)
		if item.CheckExecTime() {
			execList = append(execList, item)
		}
	}
	if len(execList) > 0 {
		sort.Sort(execList)
		coverChan(execList)
	}
}

// 清空并重建通道内容
func coverChan(execList tasklist) {
	close(TaskChan)
	for range TaskChan {
		<-TaskChan
	}
	TaskChan = make(chan *Task, chanSize)
	for i := len(execList) - 1; i >= 0; i-- {
		TaskChan <- execList[i]
	}
}
