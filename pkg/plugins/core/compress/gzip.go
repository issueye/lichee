package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
	"github.com/issueye/lichee/utils"
)

func init() {
	require.RegisterNativeModule("std/compress/gzip", func(runtime *goja.Runtime, module *goja.Object) {
		o := module.Get("exports").(*goja.Object)
		// writer
		o.Set("newWriter", func(call goja.FunctionCall) goja.Value {
			var buf bytes.Buffer
			zw := gzip.NewWriter(&buf)
			return NewGzipWriter(runtime, zw)
		})

		// reader
		o.Set("newReader", func(call goja.FunctionCall) goja.Value {
			var buf bytes.Buffer
			zr, err := gzip.NewReader(&buf)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}

			return NewGzipReader(runtime, zr, &buf)
		})

		// 解压gzip
		o.Set("unCompress", func(call goja.FunctionCall) goja.Value {
			name := call.Argument(0).String()
			unName := call.Argument(0).String()

			// 读取文件
			rf, err := os.Open(name)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}

			// 创建一个 reader
			r, err := gzip.NewReader(rf)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}

			defer r.Close()

			// 创建一个writer
			wf, err := os.Create(unName)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}

			defer wf.Close()

			// 将压缩包数据 拷贝到普通文件
			written, err := io.Copy(wf, rf)
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}

			// 返回内容
			return lib.MakeReturnValue(runtime, written)
		})
	})
}

func NewGzipReader(rt *goja.Runtime, r *gzip.Reader, buf *bytes.Buffer) *goja.Object {
	o := rt.NewObject()

	// header
	NewHeader(rt, o, &r.Header)

	// reset 重置
	o.Set("reset", func(call goja.FunctionCall) goja.Value {
		r.Reset(buf)
		return nil
	})

	// close 关闭
	o.Set("close", func(call goja.FunctionCall) goja.Value {
		err := r.Close()
		if err != nil {
			return lib.MakeErrorValue(rt, err)
		}

		return nil
	})

	// read
	o.Set("read", func(call goja.FunctionCall) goja.Value {
		s := call.Argument(0).String()
		i, err := r.Read([]byte(s))
		if err != nil {
			return lib.MakeErrorValue(rt, err)
		}

		return lib.MakeReturnValue(rt, i)
	})

	return o
}

func NewGzipWriter(rt *goja.Runtime, w *gzip.Writer) *goja.Object {
	o := rt.NewObject()

	// header
	NewHeader(rt, o, &w.Header)

	// Write 写入数据
	o.Set("write", func(call goja.FunctionCall) goja.Value {
		s := call.Argument(0).String()
		i, err := w.Write([]byte(s))
		if err != nil {
			return lib.MakeErrorValue(rt, err)
		}
		return lib.MakeReturnValue(rt, i)
	})

	// Reset 重置
	o.Set("reset", func(call goja.FunctionCall) goja.Value {
		var buf bytes.Buffer
		w.Reset(&buf)
		return nil
	})

	// flush 添加到IO层
	o.Set("flush", func(call goja.FunctionCall) goja.Value {
		err := w.Flush()
		if err != nil {
			return lib.MakeErrorValue(rt, err)
		}

		return nil
	})

	// close 关闭
	o.Set("close", func(call goja.FunctionCall) goja.Value {
		err := w.Close()
		if err != nil {
			return lib.MakeErrorValue(rt, err)
		}

		return nil
	})

	return o
}

func NewHeader(rt *goja.Runtime, o *goja.Object, header *gzip.Header) {
	// 文件名 Name
	o.Set("setName", func(call goja.FunctionCall) goja.Value {
		s := call.Argument(0).String()
		header.Name = s
		return nil
	})

	// 注释 Comment
	o.Set("setComment", func(call goja.FunctionCall) goja.Value {
		s := call.Argument(0).String()
		header.Comment = s
		return nil
	})

	// Extra 额外数据
	o.Set("setExtra", func(call goja.FunctionCall) goja.Value {
		s := call.Argument(0).Export().([]byte)
		header.Extra = s
		return nil
	})

	// ModTime  修改时间
	o.Set("setModTime, ", func(call goja.FunctionCall) goja.Value {
		s := call.Argument(0).String()
		t, err := time.Parse(utils.FormatDateTimeMs, s)
		if err != nil {
			return lib.MakeErrorValue(rt, err)
		}
		header.ModTime = t
		return nil
	})
}
