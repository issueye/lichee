package lib

import (
	"os"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func init() {
	require.RegisterNativeModule("os", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		o.Set("O_CREATE", os.O_CREATE)
		o.Set("O_WRONLY", os.O_WRONLY)
		o.Set("O_RDONLY", os.O_RDONLY)
		o.Set("O_RDWR", os.O_RDWR)
		o.Set("O_APPEND", os.O_APPEND)
		o.Set("O_EXCL", os.O_EXCL)
		o.Set("O_SYNC", os.O_SYNC)
		o.Set("O_TRUNC", os.O_TRUNC)

		o.Set("args", os.Args)

		o.Set("tempDir", func(call js.FunctionCall) js.Value {
			value := os.TempDir()
			return runtime.ToValue(value)
		})

		o.Set("hostname", func(call js.FunctionCall) js.Value {
			name, err := os.Hostname()
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			return MakeReturnValue(runtime, name)
		})

		o.Set("getEnv", func(call js.FunctionCall) js.Value {
			key := call.Argument(0).String()
			value := os.Getenv(key)
			return runtime.ToValue(value)
		})

		o.Set("remove", func(call js.FunctionCall) js.Value {
			name := call.Argument(0).String()
			err := os.Remove(name)
			if err != nil {
				return runtime.ToValue(NewError(runtime, err))
			}
			return nil
		})

		o.Set("removeAll", func(call js.FunctionCall) js.Value {
			path := call.Argument(0).String()
			err := os.RemoveAll(path)
			if err != nil {
				return runtime.ToValue(NewError(runtime, err))
			}
			return nil
		})

		o.Set("mkdir", func(call js.FunctionCall) js.Value {
			name := call.Argument(0).String()
			perm := call.Argument(1).ToInteger()
			err := os.Mkdir(name, os.FileMode(perm))
			if err != nil {
				return runtime.ToValue(NewError(runtime, err))
			}
			return nil
		})

		o.Set("mkdirAll", func(call js.FunctionCall) js.Value {
			path := call.Argument(0).String()
			perm := call.Argument(1).ToInteger()
			err := os.MkdirAll(path, os.FileMode(perm))
			if err != nil {
				return runtime.ToValue(NewError(runtime, err))
			}
			return nil
		})

		o.Set("getwd", func(call js.FunctionCall) js.Value {
			dir, err := os.Getwd()
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			return MakeReturnValue(runtime, dir)
		})

		o.Set("chdir", func(call js.FunctionCall) js.Value {
			dir := call.Argument(0).String()
			err := os.Chdir(dir)
			if err != nil {
				return runtime.ToValue(NewError(runtime, err))
			}
			return nil
		})

		o.Set("openFile", func(call js.FunctionCall) js.Value {
			name := call.Argument(0).String()
			flag := call.Argument(1).ToInteger()
			perm := call.Argument(2).ToInteger()

			file, err := os.OpenFile(name, int(flag), os.FileMode(perm))
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			return MakeReturnValue(runtime, NewFile(runtime, file))
		})

		o.Set("create", func(call js.FunctionCall) js.Value {
			name := call.Argument(0).String()
			file, err := os.Create(name)
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			return MakeReturnValue(runtime, NewFile(runtime, file))
		})

		o.Set("open", func(call js.FunctionCall) js.Value {
			name := call.Argument(0).String()
			file, err := os.Open(name)
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			return MakeReturnValue(runtime, NewFile(runtime, file))
		})

		o.Set("stat", func(call js.FunctionCall) js.Value {
			name := call.Argument(0).String()
			fileInfo, err := os.Stat(name)
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			// todo
			return MakeReturnValue(runtime, fileInfo)
		})

		o.Set("isExist", func(call js.FunctionCall) js.Value {
			p0 := GetNativeType(runtime, &call, 0)
			if err, ok := p0.(error); ok {
				return runtime.ToValue(os.IsExist(err))
			}
			panic(runtime.NewTypeError("p0 is not error type:%T", p0))
		})

		o.Set("isNotExist", func(call js.FunctionCall) js.Value {
			p0 := GetNativeType(runtime, &call, 0)
			if err, ok := p0.(error); ok {
				return runtime.ToValue(os.IsNotExist(err))
			}
			panic(runtime.NewTypeError("p0 is not error type:%T", p0))
		})
	})
}
