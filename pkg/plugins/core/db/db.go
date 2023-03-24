package db

import (
	"database/sql"
	"fmt"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
)

type PluginsDB interface {
	Query(sqlStr string) ([]map[string]interface{}, error)
	Exec(sqlStr string) (sql.Result, error)
	Trans() (*sql.Tx, error)
	NewExecResult(runtime *js.Runtime, r sql.Result) js.Value
	NewTx(runtime *js.Runtime, tx *sql.Tx) js.Value
}

type PDB struct {
	DB *sql.DB
}

// Query 查询
func (pdb *PDB) Query(sqlStr string) ([]map[string]interface{}, error) {

	//产生查询语句的Statement
	stmt, err := pdb.DB.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("生成语句错误，错误原因：%s\n", err.Error())
		return nil, err
	}
	defer stmt.Close()

	//通过Statement执行查询
	rows, err := stmt.Query()
	if err != nil {
		fmt.Printf("查询失败，失败原因：%s\n", err.Error())
		return nil, err
	}

	data, err := MakeData(rows)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (pdb *PDB) NewExecResult(runtime *js.Runtime, r sql.Result) js.Value {
	o := runtime.NewObject()
	o.Set("lastInsertId", func(call js.FunctionCall) js.Value {
		id, err := r.LastInsertId()
		if err != nil {
			return lib.MakeErrorValue(runtime, err)
		}
		return lib.MakeReturnValue(runtime, id)
	})

	o.Set("rowsAffected", func(call js.FunctionCall) js.Value {
		count, err := r.RowsAffected()
		if err != nil {
			return lib.MakeErrorValue(runtime, err)
		}
		return lib.MakeReturnValue(runtime, count)
	})
	return o
}

// Exec 执行
func (pdb *PDB) Exec(sqlStr string) (sql.Result, error) {
	r, err := pdb.DB.Exec(sqlStr)
	if err != nil {
		fmt.Printf("执行脚本失败，失败原因：%s\n", err.Error())
		return nil, err
	}
	return r, nil
}

// Trans
// 执行
func (pdb *PDB) Trans() (*sql.Tx, error) {
	r, err := pdb.DB.Begin()
	if err != nil {
		fmt.Printf("生成事务对象失败，失败原因：%s\n", err.Error())
		return nil, err
	}
	return r, nil
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

func RegisterDB(moduleName string, pdb PluginsDB) {
	require.RegisterNativeModule(moduleName, func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		// query 查询
		o.Set("query", func(call js.FunctionCall) js.Value {
			pSql := call.Arguments[0]
			rows, err := pdb.Query(pSql.String())
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}
			return lib.MakeReturnValue(runtime, rows)
		})

		// exec 执行语句 增删改
		o.Set("exec", func(call js.FunctionCall) js.Value {
			pSql := call.Arguments[0]
			r, err := pdb.Exec(pSql.String())
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}
			return lib.MakeReturnValue(runtime, pdb.NewExecResult(runtime, r))
		})

		// 事务
		// begin 开启事务
		o.Set("begin", func(call js.FunctionCall) js.Value {
			tx, err := pdb.Trans()
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}
			return lib.MakeReturnValue(runtime, pdb.NewTx(runtime, tx))
		})
	})
}
