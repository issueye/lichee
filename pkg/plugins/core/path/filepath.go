package path

import (
	"path/filepath"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
)

func init() {
	require.RegisterNativeModule("path/filepath", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		o.Set("abs", func(call js.FunctionCall) js.Value {
			path := call.Argument(0).String()
			str, err := filepath.Abs(path)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}
			return lib.MakeReturnValue(runtime, str)
		})

		o.Set("join", func(call js.FunctionCall) js.Value {
			arr := lib.GetAllArgs_string(runtime, &call)
			str := filepath.Join(arr...)
			return runtime.ToValue(str)
		})

		o.Set("ext", func(call js.FunctionCall) js.Value {
			path := call.Argument(0).String()
			str := filepath.Ext(path)
			return runtime.ToValue(str)
		})

	})
}
