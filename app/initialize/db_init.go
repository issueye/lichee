package initialize

import (
	"github.com/issueye/lichee/app/global"
	"github.com/issueye/lichee/pkg/db"
)

func InitDB() {
	d := global.LocalCfg.Db
	if d == nil {
		return
	}

	if d.Database == "" || d.Host == "" || d.Password == "" || d.Username == "" {
		return
	}

	global.LocalDb = db.InitSqlServer(global.LocalCfg.Db, global.Log)
	global.IsHaveDb = true
}
