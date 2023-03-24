package lib

import (
	"errors"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func NewError(runtime *js.Runtime, err error) *js.Object {
	o := runtime.NewObject()
	o.Set("error", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(err.Error())
	})

	o.Set("nativeType", err)

	return o
}

func init() {
	require.RegisterNativeModule("error", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		o.Set("new", func(call js.FunctionCall) js.Value {
			text := call.Argument(0).String()
			err := errors.New(text)
			return NewError(runtime, err)
		})
	})
}
