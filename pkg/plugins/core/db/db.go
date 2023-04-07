package db

import (
	"database/sql"
	"fmt"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
	"gorm.io/gorm"
)

type sqlResult struct {
	LastInsertId int64
	RowsAffected int64
}

func NewExecResult(runtime *js.Runtime, r sqlResult) js.Value {
	o := runtime.NewObject()
	o.Set("rowsAffected", func(call js.FunctionCall) js.Value {
		return lib.MakeReturnValue(runtime, r.RowsAffected)
	})
	return o
}

// MakeData
// 生成数据
func MakeData(rows *sql.Rows) ([]map[string]interface{}, error) {
	data := make([]map[string]interface{}, 0)
	cols, err := rows.Columns()
	if err != nil {
		fmt.Printf("查询失败，失败原因：%s\n", err.Error())
		return nil, err
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		pointers := make([]interface{}, len(cols))

		for i := range columns {
			pointers[i] = &columns[i]
		}
		err := rows.Scan(pointers...)
		if err != nil {
			fmt.Printf("绑定数据失败，失败原因：%s\n", err.Error())
			return nil, err
		}

		oneData := make(map[string]interface{})
		for i, v := range cols {
			valueP := pointers[i].(*interface{})
			value := *valueP
			oneData[v] = value
		}
		data = append(data, oneData)
	}
	return data, nil
}

func RegisterDB(moduleName string, gdb *gorm.DB) {
	require.RegisterNativeModule(moduleName, func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		// query 查询
		o.Set("query", func(call js.FunctionCall) js.Value {
			sqlStr := call.Argument(0).String()
			// 查询数据
			result := gdb.Raw(sqlStr)

			if result.Error != nil {
				return lib.MakeErrorValue(runtime, result.Error)
			}

			rows, err := result.Rows()
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}

			// 生成数据
			data, err := MakeData(rows)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}

			return lib.MakeReturnValue(runtime, data)
		})

		// exec 执行语句 增删改
		o.Set("exec", func(call js.FunctionCall) js.Value {
			sqlStr := call.Argument(0).String()
			result := gdb.Exec(sqlStr)
			if result.Error != nil {
				return lib.MakeErrorValue(runtime, result.Error)
			}

			// 获取返回
			return lib.MakeReturnValue(runtime, NewExecResult(runtime, sqlResult{
				RowsAffected: result.RowsAffected,
			}))
		})

		// 事务
		// begin 开启事务
		o.Set("begin", func(call js.FunctionCall) js.Value {
			tx := gdb.Begin()
			if tx.Error != nil {
				return lib.MakeErrorValue(runtime, tx.Error)
			}
			return lib.MakeReturnValue(runtime, NewTx(runtime, tx))
		})
	})
}
