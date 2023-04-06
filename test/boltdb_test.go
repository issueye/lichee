package test

import (
	"testing"

	"github.com/issueye/lichee/pkg/plugins/core"
)

func Test_boltdb(t *testing.T) {
	c := core.NewCore()
	err := c.Run("boltdb_test", "boltdb_test.js")
	if err != nil {
		t.Logf("运行脚本失败，失败原因：%s", err.Error())
	}
}
