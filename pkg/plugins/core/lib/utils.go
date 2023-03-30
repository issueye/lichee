package lib

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/google/uuid"
)

func GetNativeType(runtime *js.Runtime, call *js.FunctionCall, idx int) interface{} {
	return call.Argument(idx).ToObject(runtime).Get("nativeType").Export()
}

func GetGoType(runtime *js.Runtime, call *js.FunctionCall, idx int) js.Value {
	p := call.Argument(idx).ToObject(runtime)
	protoFunc, ok := js.AssertFunction(p.Get("getGoType"))
	if !ok {
		panic(runtime.NewTypeError("p%d not have getGoType() function", idx))
	}
	obj, err := protoFunc(p)
	if err != nil {
		panic(runtime.NewGoError(err))
	}
	return obj
}

func GetAllArgs(call *js.FunctionCall) []interface{} {
	args := make([]interface{}, len(call.Arguments))
	for i, v := range call.Arguments {
		args[i] = v.Export()
	}
	return args
}

func GetAllArgs_string(runtime *js.Runtime, call *js.FunctionCall) []string {
	args := make([]string, len(call.Arguments))
	for i, v := range call.Arguments {
		vv := v.Export()
		if s, ok := vv.(string); ok {
			args[i] = s
		} else {
			panic(runtime.NewTypeError("arg[%d] is not a string type:%T", i, v))
		}
	}
	return args
}

func MakeReturnValue(runtime *js.Runtime, value interface{}) js.Value {
	return runtime.ToValue(map[string]interface{}{
		"value": value,
	})
}

func MakeErrorValue(runtime *js.Runtime, err error) js.Value {
	return runtime.ToValue(map[string]interface{}{
		"err": NewError(runtime, err),
	})
}

func init() {
	require.RegisterNativeModule("utils", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		o.Set("print", func(call js.FunctionCall) js.Value {
			fmt.Print(call.Argument(0).String())
			return nil
		})

		o.Set("panic", func(call js.FunctionCall) js.Value {
			panic(call.Argument(0).String())
		})

		o.Set("toString", func(data []byte) js.Value {
			return runtime.ToValue(string(data))
		})

		o.Set("toBase64", func(call js.FunctionCall) js.Value {
			p0 := call.Argument(0).Export()
			var str string
			switch data := p0.(type) {
			case []byte:
				str = base64.StdEncoding.EncodeToString(data)
			case string:
				str = base64.StdEncoding.EncodeToString([]byte(data))
			default:
				panic(runtime.NewTypeError("p0 is not []byte or string type:%T", p0))
			}
			return runtime.ToValue(str)
		})

		o.Set("md5", func(call js.FunctionCall) js.Value {
			p0 := call.Argument(0).Export()
			var r []byte
			switch data := p0.(type) {
			case []byte:
				tmp := md5.Sum(data)
				r = tmp[:]
			case string:
				tmp := md5.Sum([]byte(data))
				r = tmp[:]
			default:
				panic(runtime.NewTypeError("p0 is not []byte or string type:%T", p0))
			}
			return runtime.ToValue(hex.EncodeToString(r))
		})

		o.Set("sha1", func(call js.FunctionCall) js.Value {
			p0 := call.Argument(0).Export()
			var r []byte
			switch data := p0.(type) {
			case []byte:
				tmp := sha1.Sum(data)
				r = tmp[:]
			case string:
				tmp := sha1.Sum([]byte(data))
				r = tmp[:]
			default:
				panic(runtime.NewTypeError("p0 is not []byte or string type:%T", p0))
			}
			return runtime.ToValue(hex.EncodeToString(r))
		})

		o.Set("uuid", func(call js.FunctionCall) js.Value {
			return runtime.ToValue(uuid.NewString())
		})
	})
}
