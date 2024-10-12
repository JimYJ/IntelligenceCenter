package sqlite

import (
	"IntelligenceCenter/service/log"
	"sync"

	"github.com/jmoiron/sqlx"
)

const (
	databasePath = "./database/base.db"
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
