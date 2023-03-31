package goquery

import (
	"github.com/PuerkitoBio/goquery"
	js "github.com/dop251/goja"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
)

type EachFunc = func(js.FunctionCall) js.Value

func NewDocument(runtime *js.Runtime, doc *goquery.Document) *js.Object {
	o := runtime.NewObject()
	o.Set("find", func(call js.FunctionCall) js.Value {
		param := call.Argument(0).String()
		s := doc.Find(param)
		return NewSelection(runtime, s)
	})
	return o
}

func NewSelection(runtime *js.Runtime, seletcion *goquery.Selection) *js.Object {
	o := runtime.NewObject()
	// Find
	o.Set("find", func(call js.FunctionCall) js.Value {
		param := call.Argument(0).String()
		s := seletcion.Find(param)
		return NewSelection(runtime, s)
	})

	//Attr
	o.Set("attr", func(call js.FunctionCall) js.Value {
		param := call.Argument(0).String()
		val, exists := seletcion.Attr(param)
		if exists {
			return lib.MakeReturnValue(runtime, val)
		}

		return nil
	})

	// Text
	o.Set("text", func(call js.FunctionCall) js.Value {
		value := seletcion.Text()
		return lib.MakeReturnValue(runtime, value)
	})

	// Each
	o.Set("each", func(call js.FunctionCall) js.Value {
		arg := call.Argument(0).Export()
		eachBackCall := arg.(EachFunc)
		seletcion.Each(func(i int, s *goquery.Selection) {
			eachBackCall(js.FunctionCall{
				Arguments: []js.Value{js.New().ToValue(i), NewSelection(runtime, s)},
				This:      NewSelection(runtime, s),
			})
		})
		return nil
	})

	return o
}
