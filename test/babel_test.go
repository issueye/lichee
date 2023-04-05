package test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/issueye/lichee/pkg/plugins/core"
	"github.com/issueye/lichee/pkg/plugins/core/compiler"
)

func Test_BabelTrans(t *testing.T) {
	c := compiler.New()
	c.Options.CompatibilityMode = compiler.CompatibilityModeExtended
	c2 := core.NewCore()
	b, err := ioutil.ReadFile("babel_test.js")
	if err != nil {
		t.Errorf("读取文件失败，失败原因：%s", err.Error())
	}

	t.Run("transform", func(t *testing.T) {

		// _, code, err := c.Compile(`import "something"`, "script.js", true)
		// require.NoError(t, err)
		// assert.Equal(t, `"use strict";require("something");`, code)

		_, code, err := c.Compile(string(b), "babel_test.js", true)
		if err != nil {
			t.Errorf("babel 编译失败，失败原因：%s", err.Error())
		}

		fmt.Printf("编译之后的代码 = \n%s\n", code)

		// code = strings.ReplaceAll(code, `"use strict";`, "")

		// fmt.Printf("处理之后的代码 = \n%s\n", code)

		err = c2.RunString(code)
		if err != nil {
			t.Errorf("运行代码失败，失败原因：%s", err.Error())
		}
	})
}
