package initialize

import (
	"github.com/issueye/lichee/app/global"
	"github.com/issueye/lichee/pkg/db"
)

func InitDB() {
	global.LocalDb = db.InitSqlServer(global.LocalCfg.Db, global.Log)
}
