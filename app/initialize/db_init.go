package initialize

import (
	"fmt"
	"time"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/global"
	"github.com/issueye/lichee/pkg/db"
	"github.com/issueye/lichee/utils"
	"go.etcd.io/bbolt"
)

func InitDB() {
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

	err = global.Bdb.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(common.AREA_BUCKET)
		return err
	})

	if err != nil {
		panic(fmt.Sprintf("创建参数域BUCKET失败，失败原因：%s", err.Error()))
	}

	fmt.Printf("【%s】初始化BUCKET完成...\n", utils.Ltime{}.GetNowStr())
}

func InitAdminUser() {
	isHave := false
	_ = global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.USER_BUCKET)
		data := b.Get(common.UserID(10000))
		if len(data) > 0 {
			isHave = true
		}

		return nil
	})

	if !isHave {
		global.Bdb.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket(common.USER_BUCKET)
			user := &model.User{
				Id:         10000,
				Account:    "admin",
				Name:       "管理员",
				Password:   "",
				Mark:       "系统生成的管理员账号",
				Enable:     1,
				CreateTime: time.Now(),
			}
			data, err := utils.GobBuff{}.StructToBytes(user)
			if err != nil {
				return err
			}
			return b.Put(common.UserID(10000), data)
		})
	}
}

func InitSysBucket() {
	// 查询是否存在系统参数域
	isHave := false
	global.Bdb.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(common.AREA_BUCKET)
		b := bucket.Get(common.AreaID(10000))
		if len(b) > 0 {
			isHave = true
		}

		return nil
	})

	if !isHave {
		err := global.Bdb.Update(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket(common.AREA_BUCKET)
			data := new(model.ParamArea)
			data.Id = 10000
			data.Name = common.SYS_AREA_NAME
			data.Mark = "系统配置的参数域"

			bufferData, err := utils.GobBuff{}.StructToBytes(data)
			if err != nil {
				return err
			}

			err = bucket.Put(common.AreaID(data.Id), bufferData)
			if err != nil {
				return err

			}

			_, err = bucket.CreateBucketIfNotExists(common.AreaBucketID(10000))
			return err
		})

		if err != nil {
			panic(fmt.Errorf("创建系统参数BUCKET失败，失败原因：%s", err.Error()))
		}
	}
}
