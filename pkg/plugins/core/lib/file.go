package lib

import (
	"io/ioutil"
	"os"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func NewFile(runtime *js.Runtime, file *os.File) *js.Object {
	o := runtime.NewObject()
	o.Set("writeString", func(call js.FunctionCall) js.Value {
		s := call.Argument(0).String()
		n, err := file.WriteString(s)
		if err != nil {
			return MakeErrorValue(runtime, err)
		}
		return MakeReturnValue(runtime, n)
	})

	o.Set("close", func(call js.FunctionCall) js.Value {
		err := file.Close()
		if err != nil {
			return MakeErrorValue(runtime, err)
		}
		return nil
	})

	o.Set("nativeType", file)

	return o
}

func init() {
	require.RegisterNativeModule("file", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		o.Set("write", func(call js.FunctionCall) js.Value {
			filename := call.Argument(0).String()
			data := call.Argument(1).String()
			err := ioutil.WriteFile(filename, []byte(data), 0666)
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			return nil
		})

		o.Set("read", func(call js.FunctionCall) js.Value {
			filename := call.Argument(0).String()
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			return MakeReturnValue(runtime, string(data))
		})
	})
}
