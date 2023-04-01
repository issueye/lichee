package common

import (
	"fmt"
	"net/http"
	"os"

	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/pkg/plugins/core"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	LocalDb    *gorm.DB           // 本地数据库连接
	Log        *zap.SugaredLogger // 日志记录
	Logger     *zap.Logger        // 日志记录
	Router     *gin.Engine        // 路由对象
	HttpServer *http.Server       // http 服务
	LocalCfg   *model.Config      // 本地配置信息
	ConfigPath string             // 配置文件路径
	IsHaveDb   bool               // 是否有数据库
)

var (
	JOB_BUCKET = []byte("JOB_BUCKET")
)

func JobID(id int64) []byte {
	return []byte(fmt.Sprintf("JOB:ID:%d", id))
}

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
