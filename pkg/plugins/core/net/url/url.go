package url

import (
	"encoding/json"
	"net/url"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
)

func NewURL(runtime *goja.Runtime, u *url.URL) *goja.Object {
	// TODO
	o := runtime.NewObject()
	o.Set("getForceQuery", func(call goja.FunctionCall) goja.Value {
		return runtime.ToValue(u.ForceQuery)
	})

	o.Set("getFragment", func(call goja.FunctionCall) goja.Value {
		return runtime.ToValue(u.Fragment)
	})

	o.Set("getHost", func(call goja.FunctionCall) goja.Value {
		return runtime.ToValue(u.Host)
	})

	o.Set("getOpaque", func(call goja.FunctionCall) goja.Value {
		return runtime.ToValue(u.Opaque)
	})

	o.Set("getPath", func(call goja.FunctionCall) goja.Value {
		return runtime.ToValue(u.Path)
	})

	o.Set("getRawPath", func(call goja.FunctionCall) goja.Value {
		return runtime.ToValue(u.RawPath)
	})

	o.Set("getRawQuery", func(call goja.FunctionCall) goja.Value {
		return runtime.ToValue(u.RawQuery)
	})

	o.Set("getScheme", func(call goja.FunctionCall) goja.Value {
		return runtime.ToValue(u.Scheme)
	})

	o.Set("getPort", func(call goja.FunctionCall) goja.Value {
		return runtime.ToValue(u.Port())
	})

	return o
}

func init() {
	require.RegisterNativeModule("url", func(runtime *goja.Runtime, module *goja.Object) {
		o := module.Get("exports").(*goja.Object)
		o.Set("parse", func(call goja.FunctionCall) goja.Value {
			rawurl := call.Argument(0).String()
			u, err := url.Parse(rawurl)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}
			return lib.MakeReturnValue(runtime, NewURL(runtime, u))
		})

		o.Set("queryEscape", func(call goja.FunctionCall) goja.Value {
			str := url.QueryEscape(call.Argument(0).String())
			return runtime.ToValue(str)
		})

		o.Set("queryUnescape", func(call goja.FunctionCall) goja.Value {
			str, err := url.QueryUnescape(call.Argument(0).String())
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}
			return lib.MakeReturnValue(runtime, str)
		})

		o.Set("parseRequestURI", func(call goja.FunctionCall) goja.Value {
			rawurl := call.Argument(0).String()
			mUrl, err := url.ParseRequestURI(rawurl)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}
			return lib.MakeReturnValue(runtime, NewURL(runtime, mUrl))
		})

		o.Set("parseQuery", func(call goja.FunctionCall) goja.Value {
			query := call.Argument(0).String()
			if query == "" {
				return lib.MakeReturnValue(runtime, "{}")
			}
			values, err := url.ParseQuery(query)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}
			if len(values) > 0 {
				val, err := json.Marshal(values)
				if err != nil {
					panic("parseQuery Marshal json error:" + err.Error())
				}
				return lib.MakeReturnValue(runtime, string(val))
			}
			return lib.MakeReturnValue(runtime, "{}")
		})

		o.Set("newValues", func(call goja.FunctionCall) goja.Value {
			return NewValues(runtime, make(url.Values))
		})
	})
}
