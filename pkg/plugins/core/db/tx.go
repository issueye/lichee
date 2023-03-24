package db

import (
	"database/sql"

	js "github.com/dop251/goja"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
)

func (pdb *PDB) NewTx(runtime *js.Runtime, tx *sql.Tx) js.Value {
	o := runtime.NewObject()
	o.Set("commit", func(call js.FunctionCall) js.Value {
		err := tx.Commit()
		return lib.MakeErrorValue(runtime, err)
	})

	o.Set("rollback", func(call js.FunctionCall) js.Value {
		err := tx.Rollback()
		return lib.MakeErrorValue(runtime, err)
	})

	o.Set("exec", func(call js.FunctionCall) js.Value {
		pSql := call.Arguments[0]
		r, err := tx.Exec(pSql.String())
		if err != nil {
			return lib.MakeErrorValue(runtime, err)
		}
		return lib.MakeReturnValue(runtime, pdb.NewExecResult(runtime, r))
	})

	o.Set("query", func(call js.FunctionCall) js.Value {
		pSql := call.Arguments[0]
		rows, err := tx.Query(pSql.String())
		if err != nil {
			return lib.MakeErrorValue(runtime, err)
		}

		data, err := MakeData(rows)
		if err != nil {
			return lib.MakeErrorValue(runtime, err)
		}

		return lib.MakeReturnValue(runtime, data)
	})

	return o
}
