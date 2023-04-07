package initialize

import (
	"fmt"
	"time"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/app/service"
	"github.com/issueye/lichee/global"
	"github.com/issueye/lichee/pkg/db"
	licheeDB "github.com/issueye/lichee/pkg/plugins/core/db"
	"github.com/issueye/lichee/utils"
	"go.etcd.io/bbolt"
)

func InitDB() {
	req := new(model.ReqQueryDbSource)
	list, err := service.NewDbSourceService().Query(req)
	if err != nil {
		common.Log.Errorf("获取数据库源列表失败，失败原因：%s", err.Error())
		return
	}

	for _, data := range list {

		cfg := &db.Config{
			Username: data.User,
			Password: data.Password,
			Host:     data.Host,
			Database: data.Database,
			Port:     int(data.Port),
			LogMode:  true,
		}

		d, err := common.GetDb(cfg, int(data.Type))
		if err != nil {
			common.Log.Errorf("连接数据库【%s】失败，失败原因：%s", data.Name, err.Error())
			continue
		}

		info := &global.DbInfo{
			Cfg:  cfg,
			Name: data.Name,
			DB:   d,
		}

		// 将数据库注册到JS虚拟机
		licheeDB.RegisterDB(fmt.Sprintf("db/%s", info.Name), info.DB)
		global.GdbMap[data.Name] = info
	}

	fmt.Printf("【%s】初始化数据库完成...\n", utils.Ltime{}.GetNowStr())
}

func InitBoltDb() {
	CreateBucket(common.JOB_BUCKET, "定时任务")       // JOB_BUCKET
	CreateBucket(common.USER_BUCKET, "用户")        // USER_BUCKET
	CreateBucket(common.AREA_BUCKET, "参数域")       // AREA_BUCKET
	CreateBucket(common.DB_SOURCE_BUCKET, "数据库源") // DB_SOURCE_BUCKET

	fmt.Printf("【%s】初始化BUCKET完成...\n", utils.Ltime{}.GetNowStr())
}

func CreateBucket(name []byte, describe string) {
	err := global.Bdb.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})
	if err != nil {
		panic(fmt.Sprintf("创建[%s]->BUCKET 失败，失败原因：%s", describe, err.Error()))
	}
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
