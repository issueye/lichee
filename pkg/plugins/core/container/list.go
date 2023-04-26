package container

import (
	"container/list"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
)

// 写入数据的结构体
type ElemData struct {
	Key  string
	Data any
}

func init() {
	require.RegisterNativeModule("std/container/list", func(runtime *goja.Runtime, module *goja.Object) {
		o := module.Get("exports").(*goja.Object)
		// writer
		o.Set("new", func(call goja.FunctionCall) goja.Value {
			l := list.New()
			return NewList(runtime, l)
		})
	})
}

func NewList(rt *goja.Runtime, l *list.List) *goja.Object {
	o := rt.NewObject()
	// Back() *Element
	o.Set("back", func(call goja.FunctionCall) goja.Value {
		e := l.Back()
		return NewElem(rt, e)
	})

	// Front() *list.Element
	o.Set("front", func(call goja.FunctionCall) goja.Value {
		e := l.Front()
		return NewElem(rt, e)
	})

	// Init() *list.List
	o.Set("init", func(call goja.FunctionCall) goja.Value {
		initList := l.Init()
		return NewList(rt, initList)
	})

	// Len() int
	o.Set("len", func(call goja.FunctionCall) goja.Value {
		len := l.Len()
		return lib.MakeReturnValue(rt, len)
	})

	// PushBack(v any) *list.Element
	o.Set("pushBack", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		value := call.Argument(1).Export()

		// 组装结构体
		data := &ElemData{
			Key:  key,
			Data: value,
		}

		elem := l.PushBack(data)
		return NewElem(rt, elem)
	})

	// PushFront(v any) *list.Element
	o.Set("pushFront", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		value := call.Argument(1).Export()

		// 组装结构体
		data := &ElemData{
			Key:  key,
			Data: value,
		}

		elem := l.PushFront(data)
		return NewElem(rt, elem)
	})

	// PushBack(v any) *list.Element
	o.Set("pushBack", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		value := call.Argument(1).Export()

		// 组装结构体
		data := &ElemData{
			Key:  key,
			Data: value,
		}

		elem := l.PushBack(data)
		return NewElem(rt, elem)
	})

	// InsertAfter(v any, mark *list.Element) *list.Element
	o.Set("insertAfter", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		value := call.Argument(1).Export()

		e, ok := find(l, key)
		if ok {
			// 组装结构体
			data := &ElemData{
				Key:  key,
				Data: value,
			}

			elem := l.InsertAfter(data, e)
			return NewElem(rt, elem)
		}

		return nil
	})

	// InsertBefore(v any, mark *list.Element) *list.Element
	o.Set("insertBefore", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		value := call.Argument(1).Export()

		e, ok := find(l, key)
		if ok {
			// 组装结构体
			data := &ElemData{
				Key:  key,
				Data: value,
			}

			elem := l.InsertBefore(data, e)
			return NewElem(rt, elem)
		}

		return nil
	})

	// MoveToBack(e *list.Element)
	o.Set("moveToBack", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()

		e, ok := find(l, key)
		if ok {
			// 组装结构体
			l.MoveToBack(e)
		}
		return nil
	})

	// MoveToFront(e *list.Element)
	o.Set("moveToFront", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()

		e, ok := find(l, key)
		if ok {
			// 组装结构体
			l.MoveToFront(e)
		}
		return nil
	})

	// MoveBefore(e *list.Element, mark *list.Element)
	o.Set("moveBefore", func(call goja.FunctionCall) goja.Value {
		key1 := call.Argument(0).String()
		key2 := call.Argument(0).String()

		e1, ok1 := find(l, key1)
		e2, ok2 := find(l, key2)
		if ok1 && ok2 {
			l.MoveBefore(e1, e2)
		}

		return nil
	})

	// MoveAfter(e *list.Element, mark *list.Element)
	o.Set("moveAfter", func(call goja.FunctionCall) goja.Value {
		key1 := call.Argument(0).String()
		key2 := call.Argument(0).String()

		e1, ok1 := find(l, key1)
		e2, ok2 := find(l, key2)
		if ok1 && ok2 {
			l.MoveAfter(e1, e2)
		}

		return nil
	})

	// Remove(e *list.Element) any
	o.Set("remove", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()

		e, ok := find(l, key)
		if ok {
			value := l.Remove(e)
			return lib.MakeReturnValue(rt, value)
		}

		return nil
	})

	// 查找元素
	o.Set("find", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()

		// 遍历
		for e := l.Front(); e != nil; e = e.Next() {
			data := e.Value.(*ElemData)
			if strings.EqualFold(data.Key, key) {
				return NewElem(rt, e)
			}
		}

		return nil
	})

	return o
}

func find(l *list.List, key string) (*list.Element, bool) {
	// 遍历
	for e := l.Front(); e != nil; e = e.Next() {
		data := e.Value.(*ElemData)
		if strings.EqualFold(data.Key, key) {
			return e, true
		}
	}

	return nil, false
}

func NewElem(rt *goja.Runtime, elem *list.Element) *goja.Object {
	o := rt.NewObject()
	// Next() *list.Element
	o.Set("next", func(call goja.FunctionCall) goja.Value {
		e := elem.Next()
		return NewElem(rt, e)
	})

	// Prev() *list.Element
	o.Set("prev", func(call goja.FunctionCall) goja.Value {
		e := elem.Prev()
		return NewElem(rt, e)
	})

	// Value any
	o.Set("value", func(call goja.FunctionCall) goja.Value {
		e := elem.Value
		return lib.MakeReturnValue(rt, e)
	})

	return o
}
