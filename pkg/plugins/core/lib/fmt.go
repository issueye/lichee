package lib

import (
	"fmt"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func init() {
	require.RegisterNativeModule("fmt", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		o.Set("sprintf", func(call js.FunctionCall) js.Value {
			format := call.Argument(0).String()
			args := GetAllArgs(&call)
			str := fmt.Sprintf(format, args[1:]...)
			return runtime.ToValue(str)
		})

		o.Set("printf", func(call js.FunctionCall) js.Value {
			format := call.Argument(0).String()
			args := GetAllArgs(&call)
			n, err := fmt.Printf(format, args[1:]...)
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			return MakeReturnValue(runtime, n)
		})

		o.Set("println", func(call js.FunctionCall) js.Value {
			args := GetAllArgs(&call)
			n, err := fmt.Println(args...)
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			return MakeReturnValue(runtime, n)
		})

		o.Set("print", func(call js.FunctionCall) js.Value {
			args := GetAllArgs(&call)
			n, err := fmt.Print(args...)
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			return MakeReturnValue(runtime, n)
		})
	})
}
