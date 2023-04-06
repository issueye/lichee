package test

import (
	"testing"

	"github.com/issueye/lichee/pkg/plugins/core"
)

func TestGoQuery(t *testing.T) {
	core := core.NewCore()
	err := core.Run("goquery_test", "goquery_test.js")
	if err != nil {
		t.Logf("运行代码失败，失败原因：%s", err.Error())
	}
}
