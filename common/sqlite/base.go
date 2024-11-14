package sqlite

import (
	"IntelligenceCenter/common"
	"IntelligenceCenter/service/log"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	databasePath = fmt.Sprintf("./%s/base.db", common.DBDir)
)

var (
	dbConn *sqlx.DB
	once   sync.Once
)

func Conn() *sqlx.DB {
	var err error
	once.Do(func() {
		dbConn, err = sqlx.Open("sqlite3", databasePath)
		if err != nil {
			log.Info("创建数据库连接报错:", err)
		}
	})
	return dbConn
}
