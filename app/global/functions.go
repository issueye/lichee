package global

import (
	"fmt"
	"os"

	"github.com/issueye/lichee/pkg/plugins/core"
)

func SetGlobalPathOption() func(c *core.Core) {
	return func(c *core.Core) {
		path, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("获取程序当前所在目录失败，失败原因：%s", err.Error()))
		}
		loadPath := fmt.Sprintf(`%s/runtime/js`, path)
		c.SetGlobalPath(loadPath)
	}
}

// GetInitCore
// 初始化JS插件内容
func GetInitCore() *core.Core {
	vm := core.NewCore(core.OptionLog(Logger), SetGlobalPathOption())
	return vm
}
