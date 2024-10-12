package database

import (
	"IntelligenceCenter/common/sqlite"
	"IntelligenceCenter/service/log"
	"database/sql"
)

var (
	createSqlList = []string{optionTableSql}

	checkTableSql = ""

	optionTableSql = ""
)

// 初始化数据库
func createDatabase() {
	for _, item := range createSqlList {
		_, err := sqlite.Conn().Exec(item)
		if err != nil {
			log.Info("创建初始化表失败", err)
			return
		}
	}

}

// 检查数据库
func CheckDatabase() {
	var tableExists bool
	err := sqlite.Conn().Get(&tableExists, checkTableSql)
	if err == sql.ErrNoRows {
		log.Info("找不到数据文件，准备创建...")
		createDatabase()
		return
	} else if err != nil {
		log.Info(err)
		return
	} else {
		log.Info("数据库已经存在。", tableExists)
		return
	}
}
