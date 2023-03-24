package initialize

import (
	"fmt"

	"github.com/issueye/lichee/app/global"
	pDB "github.com/issueye/lichee/pkg/plugins/core/db"
)

// InitPlugins
// 初始化插件对象
func InitPlugins() {
	global.JSPlugin = global.GetInitCore()

	// 如果有数据库则添加数据库连接对象
	if global.IsHaveDb {
		localDB, err := global.LocalDb.DB()
		if err != nil {
			panic(fmt.Errorf("获取数据库连接失败，失败原因：%s", err.Error()))
		}
		pDB.NewLocalDB(localDB)
	}
}
