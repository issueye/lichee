package boltdb

import (
	"log"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
	bolt "go.etcd.io/bbolt"
)

var (
	Bdb *bolt.DB
)

type CallBackFunc = func(js.FunctionCall) js.Value

func init() {
	if Bdb == nil {
		var err error
		Bdb, err = bolt.Open("lichee.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	require.RegisterNativeModule("db/bolt", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		o.Set("createBucketIfNotExists", func(call js.FunctionCall) js.Value {
			name := call.Argument(0).String()
			err := Bdb.Update(func(tx *bolt.Tx) error {
				_, err := tx.CreateBucketIfNotExists([]byte(name))
				if err != nil {
					return err
				}
				return nil
			})

			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}

			return nil
		})

		o.Set("view", func(call js.FunctionCall) js.Value {
			name := call.Argument(0).String()
			callBack := call.Argument(1).Export().(CallBackFunc)
			Bdb.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(name))

				callBack(js.FunctionCall{
					This:      NewBucket(runtime, b),
					Arguments: []js.Value{NewBucket(runtime, b)},
				})

				return nil
			})
			return nil
		})

		o.Set("update", func(call js.FunctionCall) js.Value {
			name := call.Argument(0).String()
			callBack := call.Argument(1).Export().(CallBackFunc)
			Bdb.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(name))

				// 回调
				callBack(js.FunctionCall{
					This:      NewBucket(runtime, b),
					Arguments: []js.Value{NewBucket(runtime, b)},
				})

				return nil
			})
			return nil
		})
	})
}
