package lib

import (
	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/spf13/cast"
	ini "gopkg.in/ini.v1"
)

var iniMap = make(map[string]*ini.File)

func NewIni(runtime *js.Runtime, cfg *ini.File, path string) js.Value {
	o := runtime.NewObject()
	// get
	o.Set("getStr", func(call js.FunctionCall) js.Value {
		sectionStr := call.Argument(0).String()
		key := call.Argument(1).String()
		value := cfg.Section(sectionStr).Key(key).String()
		return runtime.ToValue(value)
	})

	o.Set("getInt", func(call js.FunctionCall) js.Value {
		sectionStr := call.Argument(0).String()
		key := call.Argument(1).String()
		value := cfg.Section(sectionStr).Key(key).MustInt64(-1)
		return runtime.ToValue(value)
	})

	o.Set("getBool", func(call js.FunctionCall) js.Value {
		sectionStr := call.Argument(0).String()
		key := call.Argument(1).String()
		value := cfg.Section(sectionStr).Key(key).MustBool(false)
		return runtime.ToValue(value)
	})

	o.Set("getSection", func(call js.FunctionCall) js.Value {
		sectionStr := call.Argument(0).String()
		section, err := cfg.GetSection(sectionStr)
		if err != nil {
			return MakeErrorValue(runtime, err)
		}
		keys := section.KeysHash()
		return MakeReturnValue(runtime, keys)
	})

	// set
	o.Set("setStr", func(call js.FunctionCall) js.Value {
		sectionStr := call.Argument(0).String()
		key := call.Argument(1).String()
		value := call.Argument(2).String()
		cfg.Section(sectionStr).Key(key).SetValue(value)
		return nil
	})

	o.Set("setInt", func(call js.FunctionCall) js.Value {
		sectionStr := call.Argument(0).String()
		key := call.Argument(1).String()
		value := call.Argument(2).ToInteger()
		cfg.Section(sectionStr).Key(key).SetValue(cast.ToString(value))
		return nil
	})

	o.Set("setBool", func(call js.FunctionCall) js.Value {
		sectionStr := call.Argument(0).String()
		key := call.Argument(1).String()
		value := call.Argument(2).ToBoolean()
		cfg.Section(sectionStr).Key(key).SetValue(cast.ToString(value))
		return nil
	})

	o.Set("save", func(call js.FunctionCall) js.Value {
		err := cfg.SaveTo(path)
		if err != nil {
			return MakeErrorValue(runtime, err)
		}
		return nil
	})
	return o
}

func init() {
	require.RegisterNativeModule("ini", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		o.Set("create", func(call js.FunctionCall) js.Value {
			path := call.Argument(0).String()
			iniCfg, err := ini.Load(path)
			if err != nil {
				return MakeErrorValue(runtime, err)
			}
			iniMap[path] = iniCfg
			return MakeReturnValue(runtime, NewIni(runtime, iniMap[path], path))
		})
	})
}
