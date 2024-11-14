package common

var (
	// 检查初始化目录
	ProxyipDir = "proxy-ip"
	LogsDir    = "logs"
	rulesDir   = "extraction-rules"
	DBDir      = "database"
	NeedDir    = []string{LogsDir, rulesDir, DBDir, ProxyipDir}
)
