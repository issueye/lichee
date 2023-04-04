package service

import (
	"fmt"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/global"
	"github.com/issueye/lichee/utils"
	"go.etcd.io/bbolt"
)

type ParamService struct{}

func NewParamService() *ParamService {
	return new(ParamService)
}

// Save
// 写入数据
func (p ParamService) Save(data *model.Param) error {
	return global.Bdb.Update(func(tx *bbolt.Tx) error {
		// 获取BUCKET
		b := tx.Bucket(common.AREA_BUCKET)

		// 获取对应参数域的BUCKET
		areaBucket := b.Bucket(common.AreaBucketID(data.AreaId))

		if areaBucket == nil {
			return fmt.Errorf("参数域【%d】不存在，请先创建", data.AreaId)
		}

		// 将数据转换为字节切片
		byteData, err := utils.GobBuff{}.StructToBytes(data)
		if err != nil {
			return err
		}

		// 存入数据
		return areaBucket.Put(common.ParamID(data.Id), byteData)
	})
}

func (p ParamService) GetAreaById(id int64) (*model.ParamArea, error) {
	data := new(model.ParamArea)
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		areaBucket := tx.Bucket(common.AREA_BUCKET)
		b := areaBucket.Get(common.AreaID(id))
		return utils.GobBuff{}.BytesToStruct(b, data)
	})

	return data, err
}

func (p ParamService) GetAreaList(req *model.ReqQueryArea) ([]*model.ParamArea, error) {
	list := make([]*model.ParamArea, 0)
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		areaBucket := tx.Bucket(common.AREA_BUCKET)
		err := areaBucket.ForEach(func(k, v []byte) error {
			data := new(model.ParamArea)

			if req.Name != "" {
				err := Find(req.Name, v, data, &list)
				if err != nil {
					return err
				}
				return nil
			}

			if len(v) > 0 {
				err := utils.GobBuff{}.BytesToStruct(v, data)
				if err != nil {
					common.Log.Errorf("将字节转换为对象失败，失败原因：%s", err.Error())
					return err
				}

				list = append(list, data)
			}

			return nil
		})

		return err
	})
	return list, err
}

// Del
// 写入数据
// func (p ParamService) Del(id int64) error {
// 	return global.Bdb.Update(func(tx *bbolt.Tx) error {
// 		b := tx.Bucket(common.USER_BUCKET)
// 		return b.Delete(common.UserID(id))
// 	})
// }

// DelArea
// 写入数据
func (p ParamService) DelArea(id int64) error {
	return global.Bdb.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.AREA_BUCKET)

		// 删除对应的参数域 BUCKET
		err := b.DeleteBucket(common.AreaBucketID(id))
		if err != nil {
			return err
		}

		// 参数参数域信息
		return b.Delete(common.AreaID(id))
	})
}

// Delete
// 写入数据
func (p ParamService) GetById(id int64) (*model.Param, error) {
	data := new(model.Param)
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.PARAM_BUCKET)
		byteData := b.Get(common.ParamID(id))
		err := utils.GobBuff{}.BytesToStruct(byteData, data)
		if err != nil {
			return err
		}
		return nil
	})

	return data, err
}

// ModifyArea
// 写入数据
func (p ParamService) ModifyArea(data *model.ParamArea) error {
	err := global.Bdb.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.AREA_BUCKET)

		bufferData, err := utils.GobBuff{}.StructToBytes(data)
		if err != nil {
			return err
		}

		return b.Put(common.AreaID(data.Id), bufferData)
	})

	return err
}

// DelParam
// 写入数据
func (p ParamService) DelParam(area, id int64) error {
	return global.Bdb.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.AREA_BUCKET)

		paramBucket := b.Bucket(common.AreaBucketID(area))

		// 参数参数域信息
		return paramBucket.Delete(common.ParamID(id))
	})
}

// ModifyParam
// 写入数据
func (p ParamService) ModifyParam(data *model.Param) error {
	err := global.Bdb.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.AREA_BUCKET)

		paramBucket := b.Bucket(common.AreaBucketID(data.AreaId))
		bufferData, err := utils.GobBuff{}.StructToBytes(data)
		if err != nil {
			return err
		}

		return paramBucket.Put(common.ParamID(data.Id), bufferData)
	})

	return err
}

// Query
// 写入数据
func (p ParamService) Query(req *model.ReqQueryParam) ([]*model.Param, error) {
	list := make([]*model.Param, 0)
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		area := tx.Bucket(common.AREA_BUCKET)
		// 读取对应参数域的参数
		param := area.Bucket(common.AreaBucketID(req.AreaId))
		if param == nil {
			return fmt.Errorf("编码【%d】的参数域没有找到", req.AreaId)
		}
		err := param.ForEach(func(k, v []byte) error {
			data := new(model.Param)
			if req.Name != "" {
				err := Find(req.Name, v, data, &list)
				if err != nil {
					return err
				}
				return nil
			}

			if len(v) > 0 {
				err := utils.GobBuff{}.BytesToStruct(v, data)
				if err != nil {
					common.Log.Errorf("将字节转换为对象失败，失败原因：%s", err.Error())
					return err
				}

				list = append(list, data)
			}

			return nil
		})

		return err
	})

	// 数据的条数
	req.Total = int64(len(list))
	return list, err
}

func (p ParamService) CreateArea(area *model.ParamArea) error {
	return global.Bdb.Update(func(tx *bbolt.Tx) error {
		areaBucket := tx.Bucket(common.AREA_BUCKET)

		// 添加数据
		data, err := utils.GobBuff{}.StructToBytes(area)
		if err != nil {
			return err
		}
		err = areaBucket.Put(common.AreaID(area.Id), data)
		if err != nil {
			return err
		}

		_, err = areaBucket.CreateBucket(common.AreaBucketID(area.Id))
		return err
	})
}
