package lib

import (
	"time"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func NewTime(runtime *js.Runtime, t *time.Time) *js.Object {
	o := runtime.NewObject()
	o.Set("string", func(call js.FunctionCall) js.Value {
		str := t.String()
		return runtime.ToValue(str)
	})

	return o
}

func GetNowDateTime(runtime *js.Runtime, formatStr string) js.Value {
	cst := time.FixedZone("CST", 28800)
	str := time.Now().In(cst).Format(formatStr)
	return runtime.ToValue(str)
}

func init() {
	require.RegisterNativeModule("time", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		// 睡眠
		o.Set("sleep", func(call js.FunctionCall) js.Value {
			d := call.Argument(0).ToInteger()
			<-time.After(time.Duration(d) * time.Millisecond)
			return nil
		})

		// 当前时间字符串
		o.Set("nowString", func(call js.FunctionCall) js.Value {
			const DATE_TIME_FORMAT = "2006-01-02 15:04:05"
			return GetNowDateTime(runtime, DATE_TIME_FORMAT)
		})

		o.Set("nowDate", func(call js.FunctionCall) js.Value {
			const DATE_TIME_FORMAT = "2006-01-02"
			return GetNowDateTime(runtime, DATE_TIME_FORMAT)
		})

		o.Set("nowYear", func(call js.FunctionCall) js.Value {
			const DATE_TIME_FORMAT = "2006"
			return GetNowDateTime(runtime, DATE_TIME_FORMAT)
		})

		o.Set("nowMonth", func(call js.FunctionCall) js.Value {
			const DATE_TIME_FORMAT = "01"
			return GetNowDateTime(runtime, DATE_TIME_FORMAT)
		})

		o.Set("nowDay", func(call js.FunctionCall) js.Value {
			const DATE_TIME_FORMAT = "02"
			return GetNowDateTime(runtime, DATE_TIME_FORMAT)
		})

		o.Set("nowHour", func(call js.FunctionCall) js.Value {
			const DATE_TIME_FORMAT = "15"
			return GetNowDateTime(runtime, DATE_TIME_FORMAT)
		})

		o.Set("nowMinute", func(call js.FunctionCall) js.Value {
			const DATE_TIME_FORMAT = "04"
			return GetNowDateTime(runtime, DATE_TIME_FORMAT)
		})

		o.Set("nowSecond", func(call js.FunctionCall) js.Value {
			const DATE_TIME_FORMAT = "05"
			return GetNowDateTime(runtime, DATE_TIME_FORMAT)
		})
	})
}
