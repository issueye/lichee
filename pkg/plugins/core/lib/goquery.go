package lib

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	resty "github.com/go-resty/resty/v2"
)

func init() {
	require.RegisterNativeModule("go/query", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		o.Set("do", func(call js.FunctionCall) js.Value {
			method := call.Argument(0).String()
			url := call.Argument(1).String()
			client := resty.New()

			var (
				resp *resty.Response
				err  error
			)
			switch strings.ToUpper(method) {
			case "POST":
				resp, err = client.R().Post(url)
			case "GET":
				resp, err = client.R().Get(url)
			}

			if err != nil {
				return NewError(runtime, err)
			}

			if resp.StatusCode() != 200 {
				fmt.Printf("请求页面失败 %d\n", resp.StatusCode())
				return nil
			}

			// data := resp.String()

			doc, err := goquery.NewDocumentFromReader(resp.RawBody())
			if err != nil {
				return NewError(runtime, err)
			}

			doc.Find("widget-flowList").Each(func(i int, s *goquery.Selection) {
				content := s.Find("s-left").Text()
				fmt.Printf("clearfix %d: %s\n", i, content)
			})

			// fmt.Printf("data %s \n", data)
			return o
		})
	})
}
