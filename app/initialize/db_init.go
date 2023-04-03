package initialize

import (
	"fmt"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/global"
	"github.com/issueye/lichee/pkg/db"
	"github.com/issueye/lichee/utils"
	"go.etcd.io/bbolt"
)

func InitDB() {

	InitBoltDb()

	var flag string
	if common.LocalCfg.UseDB {
		flag = "是"
	} else {
		flag = "否"
	}
	fmt.Printf("【%s】是否需要使用SQL数据库：%s\n", utils.Ltime{}.GetNowStr(), flag)

	if !common.LocalCfg.UseDB {
		return
	}

	common.LocalDb = db.InitSqlServer(common.LocalCfg.Db, common.Log)

	fmt.Printf("【%s】初始化数据库完成...\n", utils.Ltime{}.GetNowStr())
}

func InitBoltDb() {
	err := global.Bdb.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(common.JOB_BUCKET)
		return err
	})

	if err != nil {
		panic(fmt.Sprintf("创建定时任务BUCKET失败，失败原因：%s", err.Error()))
	}

	err = global.Bdb.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(common.USER_BUCKET)
		return err
	})

	if err != nil {
		panic(fmt.Sprintf("创建用户BUCKET失败，失败原因：%s", err.Error()))
	}

	fmt.Printf("【%s】初始化[JOB]BUCKET完成...\n", utils.Ltime{}.GetNowStr())
}
