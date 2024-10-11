package log

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger = setupLogger()
)

// 设置日志记录器
func setupLogger() *logrus.Logger {
	logger := logrus.New()
	// 获取当前日期并格式化
	currentDate := time.Now().Format("2006-01-02")
	logFileName := "logs/" + currentDate + ".log"

	// 设置日志输出为文件，按日期分文件
	logger.SetOutput(&lumberjack.Logger{
		Filename:   logFileName, // 日志文件位置
		MaxSize:    10,          // 每个文件最大尺寸 (MB)
		MaxBackups: 3,           // 保留旧文件的最大数量
		MaxAge:     30,          // 旧文件保留天数
		Compress:   true,        // 是否压缩/归档旧文件
	})

	// 设置自定义日志格式
	logger.SetFormatter(&CustomFormatter{})
	// 可以根据需要设置日志级别
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

// 自定义日志中间件
func Logs() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// 处理请求
		c.Next()
		duration := time.Since(start)
		Logger.Infof("%d | %s | %s | %s | %q",
			c.Writer.Status(),  // HTTP 状态码
			duration,           // 请求处理时间
			c.ClientIP(),       // 客户端 IP 地址
			c.Request.Method,   // 请求方法
			c.Request.URL.Path, // 请求路径
		)
	}
}

// 自定义日志格式
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 格式化日志输出
	logLine := entry.Time.Format("2006-01-02 15:04:05")
	return []byte("[GIN] " + logLine + " " + entry.Message + "\n"), nil
}

// 获取调用者的文件和行号
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2) // 获取上层调用信息
	if ok {
		return file + ":" + fmt.Sprintf("%d", line)
	}
	return "unknown"
}

// 自定义日志记录函数，包含位置信息
func Info(msg ...any) {
	location := getCallerInfo()
	var gap any = " "
	tempMsg := make([]any, 0)
	tempMsg = append(tempMsg, location)
	for _, item := range msg {
		tempMsg = append(tempMsg, gap, item)
	}
	Logger.Log(logrus.InfoLevel, tempMsg...)
}
