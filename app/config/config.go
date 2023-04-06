package config

import (
	"fmt"
	"strings"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/app/service"
	"github.com/issueye/lichee/global"
	"github.com/issueye/lichee/utils"
	"github.com/spf13/cast"
	"go.etcd.io/bbolt"
)

type Result struct {
	value  any    // 参数值
	name   string // 参数名称
	mark   string // 参数说明
	areaId int64  // 参数域
}

func (r *Result) Name() string {
	return r.name
}

func (r *Result) String() string {
	return cast.ToString(r.value)
}

func (r *Result) Int64() int64 {
	return cast.ToInt64(r.value)
}

func (r *Result) Int() int {
	return cast.ToInt(r.value)
}

func (r *Result) Float() float64 {
	return cast.ToFloat64(r.value)
}

func (r *Result) Bool() bool {
	return cast.ToBool(r.value)
}

// 获取参数
func GetSysParam(name string) *Result {
	r := new(Result)
	global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.AREA_BUCKET)
		sysParamBucket := b.Bucket(common.AreaBucketID(common.SYS_AREA))

		err := sysParamBucket.ForEach(func(k, v []byte) error {
			data := new(model.Param)
			err := utils.GobBuff{}.BytesToStruct(v, data)
			if err != nil {
				return err
			}

			// 判断是否等于
			if strings.EqualFold(name, data.Name) {
				r.name = data.Name
				r.value = data.Value
				r.mark = data.Mark
				r.areaId = data.AreaId
			}

			return nil
		})

		return err
	})

	return r
}

// 获取参数
func SetSysParam(name, value, mark string) error {
	data := new(model.Param)
	data.Id = utils.Lid{}.GenID()
	data.Name = name
	data.Mark = mark
	data.Value = value
	data.AreaId = common.SYS_AREA
	data.Area = common.SYS_AREA_NAME
	return service.NewParamService().Save(data)
}

// 判断如果没有则添加
func IsNotExitSetSysParam(name, value, mark string) {
	var isHave bool
	err := global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.AREA_BUCKET)
		sysParamBucket := b.Bucket(common.AreaBucketID(common.SYS_AREA))

		err := sysParamBucket.ForEach(func(k, v []byte) error {
			data := new(model.Param)
			err := utils.GobBuff{}.BytesToStruct(v, data)
			if err != nil {
				return err
			}
			if strings.EqualFold(name, data.Name) {
				isHave = true
			}

			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("【%s】参数出现异常\n", name)
	}
	if !isHave {
		err := SetSysParam(name, value, mark)
		if err != nil {
			fmt.Printf("【%s】参数写入失败，失败原因：%s", name, err.Error())
		}
	}
}
