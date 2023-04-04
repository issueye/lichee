package service

import (
	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/global"
	"github.com/issueye/lichee/utils"
	"go.etcd.io/bbolt"
)

type UserService struct{}

func NewUserService() *UserService {
	return new(UserService)
}

// Save
// 写入数据
func (user UserService) Save(data *model.User) error {
	return global.Bdb.Update(func(tx *bbolt.Tx) error {
		// 创建 bucket
		b := tx.Bucket(common.USER_BUCKET)

		// 将数据转换为字节切片
		byteData, err := utils.GobBuff{}.StructToBytes(data)
		if err != nil {
			return err
		}

		// 存入数据
		err = b.Put(common.UserID(data.Id), byteData)
		if err != nil {
			return err
		}

		return nil
	})
}

func (user UserService) FindUser(lu *model.LoginUser) (*model.User, error) {
	data := new(model.User)
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.USER_BUCKET)
		err := b.ForEach(func(k, v []byte) error {
			tmpData := new(model.User)
			err := utils.GobBuff{}.BytesToStruct(v, tmpData)
			if err != nil {
				common.Log.Errorf("将字节转换为对象失败，失败原因：%s", err.Error())
				return err
			}

			// 判断账号和密码相等的用户
			if tmpData.Account == lu.Account {
				data = tmpData
			}
			return nil
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}

// Delete
// 写入数据
func (user UserService) Del(id int64) error {
	return global.Bdb.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.USER_BUCKET)
		return b.Delete(common.UserID(id))
	})
}

// Delete
// 写入数据
func (user UserService) GetById(id int64) (*model.User, error) {
	data := new(model.User)
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.USER_BUCKET)
		byteData := b.Get(common.UserID(id))
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
func (user UserService) Query(req *model.ReqQueryUser) ([]*model.ResQueryUser, error) {
	list := make([]*model.User, 0)
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.USER_BUCKET)
		return b.ForEach(func(k, v []byte) error {
			data := new(model.User)

			// 用户名
			if req.Name != "" {
				err := Find(req.Name, v, data, &list)
				if err != nil {
					return err
				}

				return nil
			}

			// 登录名
			if req.Account != "" {
				err := Find(req.Account, v, data, &list)
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

	resList := make([]*model.ResQueryUser, 0)

	for _, data := range list {
		res := new(model.ResQueryUser)
		res.Id = data.Id
		res.Name = data.Name
		res.Account = data.Account
		res.Mark = data.Mark
		res.Enable = data.Enable
		res.CreateTime = data.CreateTime.Format(utils.FormatDateTimeMs)
		res.LoginTime = data.LoginTime.Format(utils.FormatDateTimeMs)
		resList = append(resList, res)
	}

	return resList, err
}
