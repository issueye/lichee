package service

import (
	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/global"
	"github.com/issueye/lichee/utils"
	"go.etcd.io/bbolt"
)

type JobService struct{}

func NewJobService() *JobService {
	return new(JobService)
}

// Save
// 写入数据
func (job JobService) Save(data *model.Job) error {
	return global.Bdb.Update(func(tx *bbolt.Tx) error {
		// 创建 bucket
		b := tx.Bucket(common.JOB_BUCKET)

		// 将数据转换为字节切片
		byteData, err := utils.GobBuff{}.StructToBytes(data)
		if err != nil {
			return err
		}

		// 存入数据
		err = b.Put(common.JobID(data.Id), byteData)
		if err != nil {
			return err
		}

		return nil
	})
}

// Delete
// 写入数据
func (job JobService) Del(id int64) error {
	return global.Bdb.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.JOB_BUCKET)
		return b.Delete(common.JobID(id))
	})
}

// Delete
// 写入数据
func (job JobService) GetById(id int64) (*model.Job, error) {
	data := new(model.Job)
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.JOB_BUCKET)
		byteData := b.Get(common.JobID(id))
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
func (job JobService) Query(req *model.ReqQueryJob) ([]*model.ResQueryJob, error) {
	list := make([]*model.Job, 0)
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.JOB_BUCKET)
		return b.ForEach(func(k, v []byte) error {
			data := new(model.Job)

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

	resList := make([]*model.ResQueryJob, 0)

	for _, data := range list {
		// 获取参数域
		pa, err := NewParamService().GetAreaById(data.AreaId)
		if err != nil {
			common.Log.Errorf("【%d】未找到参数域信息", data.AreaId)
			continue
		}

		res := new(model.ResQueryJob)
		res.Id = data.Id
		res.Name = data.Name
		res.Expr = data.Expr
		res.Mark = data.Mark
		res.Enable = data.Enable
		res.Path = data.Path
		res.AreaId = pa.Id
		res.Area = pa.Name
		res.CreateTime = data.CreateTime.Format(utils.FormatDateTimeMs)
		resList = append(resList, res)
	}

	return resList, err
}
