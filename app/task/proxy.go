package task

import (
	"IntelligenceCenter/common"
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/service/log"
)

var (
	dirList = []string{"/" + common.ProxyipDir + "/http", "/" + common.ProxyipDir + "/https", "/" + common.ProxyipDir + "/socks5"}
)

func getProxyIP() []string {
	ipList := make([]string, 0)
	for i, item := range dirList {
		ips, err := utils.ReadProxy(item, uint8(i)+1)
		if err != nil {
			log.Info("获取代理IP地址失败")
			continue
		}
		ipList = append(ipList, ips...)
	}
	return ipList
}
