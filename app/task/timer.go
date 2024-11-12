package task

import (
	"IntelligenceCenter/service/log"
	"sort"
	"time"
)

var (
	NewTaskChan = make(chan bool, chanSize)
	TaskChan    = make(chan *Task, chanSize)
	cancelChan  = make(chan bool, 1)
	chanSize    = 1024
)

// 监听延时任务
func Listen() {
	for {
		task := <-TaskChan
		now := time.Now().Local()
		if task.ExecType == 1 || task.ExecTimeSec >= now.Unix() {
			log.Info("任务类型是单次立即执行任务或已经超过指定时间，立即执行任务:", task.TaskName, task.ID)
			go task.Exec()
		}
		timer := time.NewTimer(time.Unix(task.ExecTimeSec, 0).Sub(now))
		log.Info("开始等待任务:", task.TaskName, "准备执行时间:", task.ExecTime, "解析后的时间:", time.Unix(task.ExecTimeSec, 0).Format(time.DateTime))
		select {
		case <-timer.C:
			log.Info("已经到指定时间，开始执行任务:", task.TaskName, task.ID)
			go task.Exec()
			continue
		case <-cancelChan:
			log.Info("有新任务或任务被删除或结束了,重新加载,被取消等待的策略是:", task.TaskName, task.ID)
			continue
		}
	}
}

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
