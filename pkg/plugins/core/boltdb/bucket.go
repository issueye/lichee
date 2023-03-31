package boltdb

import (
	js "github.com/dop251/goja"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
	bolt "go.etcd.io/bbolt"
)

func NewBucket(runtime *js.Runtime, bucket *bolt.Bucket) *js.Object {
	o := runtime.NewObject()
	// create
	o.Set("create", func(call js.FunctionCall) js.Value {
		name := call.Argument(0).String()
		b, err := bucket.CreateBucket([]byte(name))
		if err != nil {
			return lib.MakeErrorValue(runtime, err)
		}

		return NewBucket(runtime, b)
	})

	// get
	o.Set("get", func(call js.FunctionCall) js.Value {
		key := call.Argument(0).String()
		b := bucket.Get([]byte(key))
		return lib.MakeReturnValue(runtime, string(b))
	})

	// put
	o.Set("put", func(call js.FunctionCall) js.Value {
		key := call.Argument(0).String()
		value := call.Argument(1).String()
		err := bucket.Put([]byte(key), []byte(value))
		if err != nil {
			return lib.MakeErrorValue(runtime, err)
		}
		return nil
	})

	// delete
	o.Set("delete", func(call js.FunctionCall) js.Value {
		key := call.Argument(0).String()
		err := bucket.Delete([]byte(key))
		if err != nil {
			return lib.MakeErrorValue(runtime, err)
		}
		return nil
	})

	// forEach
	o.Set("foreach", func(call js.FunctionCall) js.Value {
		callback := call.Argument(0).Export().(CallBackFunc)
		bucket.ForEach(func(k, v []byte) error {
			callback(js.FunctionCall{
				This: NewBucket(runtime, bucket),
				Arguments: []js.Value{
					js.New().ToValue(string(k)),
					js.New().ToValue(string(v)),
				},
			})
			return nil
		})
		return nil
	})

	// forEachBucket
	o.Set("foreachBucket", func(call js.FunctionCall) js.Value {
		callback := call.Argument(0).Export().(CallBackFunc)
		bucket.ForEachBucket(func(k []byte) error {
			callback(js.FunctionCall{
				This: NewBucket(runtime, bucket),
				Arguments: []js.Value{
					js.New().ToValue(string(k)),
				},
			})
			return nil
		})
		return nil
	})

	// cursor
	o.Set("cursor", func(call js.FunctionCall) js.Value {
		corsor := bucket.Cursor()
		return NewCursor(runtime, corsor)
	})

	return o
}

func NewCursor(runtime *js.Runtime, cursor *bolt.Cursor) *js.Object {
	o := runtime.NewObject()

	// first
	o.Set("first", func(call js.FunctionCall) js.Value {
		key, value := cursor.First()
		return lib.MakeReturnValue(runtime, map[string]string{
			string(key): string(value),
		})
	})

	// last
	o.Set("last", func(call js.FunctionCall) js.Value {
		key, value := cursor.First()
		return lib.MakeReturnValue(runtime, map[string]string{
			string(key): string(value),
		})
	})

	// next
	o.Set("next", func(call js.FunctionCall) js.Value {
		key, value := cursor.Next()
		return lib.MakeReturnValue(runtime, map[string]string{
			string(key): string(value),
		})
	})

	// seek
	o.Set("seek", func(call js.FunctionCall) js.Value {
		prefix := call.Argument(0).String()
		key, value := cursor.Seek([]byte(prefix))
		return lib.MakeReturnValue(runtime, map[string]string{
			string(key): string(value),
		})
	})

	// prev
	o.Set("prev", func(call js.FunctionCall) js.Value {
		key, value := cursor.Prev()
		return lib.MakeReturnValue(runtime, map[string]string{
			string(key): string(value),
		})
	})

	return o
}
