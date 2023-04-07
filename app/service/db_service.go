package service

import (
	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/global"
	"github.com/issueye/lichee/utils"
	"go.etcd.io/bbolt"
)

type DbSourceService struct{}

func NewDbSourceService() *DbSourceService {
	return new(DbSourceService)
}

// Create
// 写入数据
func (dbSource DbSourceService) Create(data *model.DbSource) error {
	return global.Bdb.Update(func(tx *bbolt.Tx) error {
		// 创建 bucket
		b := tx.Bucket(common.DB_SOURCE_BUCKET)

		// 将数据转换为字节切片
		byteData, err := utils.GobBuff{}.StructToBytes(data)
		if err != nil {
			return err
		}

		// 存入数据
		err = b.Put(common.DbSourceID(data.Id), byteData)
		if err != nil {
			return err
		}

		return nil
	})
}

// Delete
// 写入数据
func (dbSource DbSourceService) Del(id int64) error {
	return global.Bdb.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.DB_SOURCE_BUCKET)
		return b.Delete(common.DbSourceID(id))
	})
}

// Delete
// 写入数据
func (dbSource DbSourceService) GetById(id int64) (*model.DbSource, error) {
	data := new(model.DbSource)
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.DB_SOURCE_BUCKET)
		byteData := b.Get(common.DbSourceID(id))
		err := utils.GobBuff{}.BytesToStruct(byteData, data)
		if err != nil {
			return err
		}
		return nil
	})

	return data, err
}

// Query
// 写入数据
func (dbSource DbSourceService) Query(req *model.ReqQueryDbSource) ([]*model.DbSource, error) {
	list := make([]*model.DbSource, 0)
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.DB_SOURCE_BUCKET)
		return b.ForEach(func(k, v []byte) error {
			data := new(model.DbSource)

			// 任务名称
			if req.Name != "" {
				err := Find(req.Name, v, data, &list)
				if err != nil {
					return err
				}

				return nil
			}

			err := utils.GobBuff{}.BytesToStruct(v, data)
			if err != nil {
				common.Log.Errorf("将字节转换为对象失败，失败原因：%s", err.Error())
				return err
			}

			list = append(list, data)
			return nil
		})
	})

	// 数据的条数
	req.Total = int64(len(list))
	return list, err
}
