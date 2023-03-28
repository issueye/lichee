package test

import (
	"testing"

	"github.com/issueye/lichee/pkg/plugins/core"
)

func Test_GoQuery(t *testing.T) {
	src := `
    var query = require('go/query')
    let doc = query.do('get', 'https://m.163.com/touch/tech/#adaptation=pc&refer=https%3A%2F%2Fnews.163.com%2F')
    console.log(doc)
    `

	c := core.NewCore()
	err := c.RunString(src)
	if err != nil {
		t.Errorf("运行脚本失败，失败原因：%s", err.Error())
	}
}
