package lib

import (
	"bytes"
	"fmt"
	"os/exec"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func init() {
	require.RegisterNativeModule("os/exec", func(runtime *js.Runtime, module *js.Object) {
		o := module.Get("exports").(*js.Object)
		o.Set("command", func(name string, args ...string) js.Value {
			cmd := exec.Command(name, args...)

			var (
				stdout bytes.Buffer
				stderr bytes.Buffer
			)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println("运行错误 ", err.Error())
				return MakeErrorValue(runtime, fmt.Errorf("cmd.Run() failed with %s", err.Error()))
			}

			outStr := ConvertByte2String(stdout.Bytes(), "GB18030")
			errStr := ConvertByte2String(stderr.Bytes(), "GB18030")
			if errStr != "" {
				return MakeErrorValue(runtime, fmt.Errorf("stderr %s", errStr))
			}

			return MakeReturnValue(runtime, outStr)
		})
	})
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
