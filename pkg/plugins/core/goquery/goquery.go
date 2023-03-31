package goquery

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	resty "github.com/go-resty/resty/v2"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
	"github.com/spf13/cast"
)

func init() {
	require.RegisterNativeModule("go/query", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		o.Set("do", func(call js.FunctionCall) js.Value {
			method := call.Argument(0).String()
			url := call.Argument(1).String()
			header := call.Argument(2).Export().(map[string]interface{})
			body := call.Argument(3).String()
			timeout := call.Argument(4).ToInteger()

			client := resty.New()
			r := client.R()
			// header
			for key, value := range header {
				client.Header.Add(key, cast.ToString(value))
			}
			// timeout
			client.SetTimeout(time.Duration(timeout) * time.Millisecond)
			// body
			r.SetBody(body)

			var (
				resp *resty.Response
				err  error
			)
			switch strings.ToUpper(method) {
			case "POST":
				resp, err = r.Post(url)
			case "GET":
				resp, err = r.Get(url)
			}

			if err != nil {
				return lib.NewError(runtime, err)
			}

			if resp.StatusCode() != 200 {
				fmt.Printf("请求页面失败 %d\n", resp.StatusCode())
				return nil
			}

			data := resp.Body()
			reader := bytes.NewReader(data)
			doc, err := goquery.NewDocumentFromReader(reader)
			if err != nil {
				return lib.NewError(runtime, err)
			}

			return NewDocument(runtime, doc)
		})
	})
}
