package initialize

import (
	"fmt"

	"github.com/issueye/lichee/app/common"
	pDB "github.com/issueye/lichee/pkg/plugins/core/db"
)

// InitPlugins
// 初始化插件对象
func InitPlugins() {
	// 如果有数据库则添加数据库连接对象
	if common.LocalCfg.UseDB {
		localDB, err := common.LocalDb.DB()
		if err != nil {
			panic(fmt.Errorf("获取数据库连接失败，失败原因：%s", err.Error()))
		}
		pDB.NewLocalDB(localDB)
	}
}
