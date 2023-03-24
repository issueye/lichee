package http

import (
	"mime/multipart"

	js "github.com/dop251/goja"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
)

func NewMultipartFile(runtime *js.Runtime, file multipart.File) *js.Object {
	o := runtime.NewObject()
	o.Set("nativeType", file)

	o.Set("read", func(call js.FunctionCall) js.Value {
		p0 := call.Argument(0).Export()
		if buf, ok := p0.([]byte); ok {
			n, err := file.Read(buf)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}
			return lib.MakeReturnValue(runtime, n)
		}
		panic(runtime.NewTypeError("p0 is not []byte type:%T", p0))
	})

	o.Set("readAt", func(call js.FunctionCall) js.Value {
		p0 := call.Argument(0).Export()
		off := call.Argument(1).ToInteger()
		if buf, ok := p0.([]byte); ok {
			n, err := file.ReadAt(buf, off)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}
			return lib.MakeReturnValue(runtime, n)
		}
		panic(runtime.NewTypeError("p0 is not []byte type:%T", p0))
	})

	o.Set("seek", func(call js.FunctionCall) js.Value {
		offset := call.Argument(0).ToInteger()
		whence := call.Argument(1).ToInteger()
		v, err := file.Seek(offset, int(whence))
		if err != nil {
			return lib.MakeErrorValue(runtime, err)
		}
		return lib.MakeReturnValue(runtime, v)
	})

	o.Set("close", func(call js.FunctionCall) js.Value {
		err := file.Close()
		if err != nil {
			return runtime.ToValue(lib.NewError(runtime, err))
		}
		return nil
	})

	return o
}
