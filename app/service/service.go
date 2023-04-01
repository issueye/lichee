package service

import (
	"bytes"
	"reflect"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/utils"
)

func Find(condition string, v []byte, data any, list any) error {
	ok := bytes.Contains(v, []byte(condition))
	if ok {
		err := utils.GobBuff{}.BytesToStruct(v, data)
		if err != nil {
			common.Log.Errorf("将字节转换为对象失败，失败原因：%s", err.Error())
			return err
		}

		listRef := reflect.ValueOf(list)
		if listRef.Kind() == reflect.Ptr {
			listRef = listRef.Elem()
		}

		// 判断是否是切片，如果是切片则获取切片的对应数据
		if listRef.Kind() == reflect.Slice {
			// 判断传入的切片长度，如果传入的长度和 size相等，则不需要再重新切
			listRef.Set(reflect.Append(listRef, reflect.ValueOf(data)))
		}
	}

	return nil
}
