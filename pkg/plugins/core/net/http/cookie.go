package http

import (
	"net/http"

	js "github.com/dop251/goja"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
)

func NewCookie(runtime *js.Runtime, cookie *http.Cookie) *js.Object {
	o := runtime.NewObject()

	o.Set("string", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(cookie.String())
	})

	o.Set("getDomain", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(cookie.Domain)
	})

	o.Set("getExpires", func(call js.FunctionCall) js.Value {
		return lib.NewTime(runtime, &cookie.Expires)
	})

	o.Set("getHttpOnly", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(cookie.HttpOnly)
	})

	o.Set("getMaxAge", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(cookie.MaxAge)
	})

	o.Set("getName", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(cookie.Name)
	})

	o.Set("getPath", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(cookie.Path)
	})

	o.Set("getRaw", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(cookie.Raw)
	})

	o.Set("getRawExpires", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(cookie.RawExpires)
	})

	o.Set("getSecure", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(cookie.Secure)
	})

	o.Set("getUnparsed", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(cookie.Unparsed)
	})

	o.Set("getValue", func(call js.FunctionCall) js.Value {
		return runtime.ToValue(cookie.Value)
	})

	return o
}
