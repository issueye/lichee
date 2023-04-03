package common

import (
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/pkg/plugins/core"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	LocalDb    *gorm.DB              // 本地数据库连接
	Log        *zap.SugaredLogger    // 日志记录
	Logger     *zap.Logger           // 日志记录
	Router     *gin.Engine           // 路由对象
	HttpServer *http.Server          // http 服务
	LocalCfg   *model.Config         // 本地配置信息
	ConfigPath string                // 配置文件路径
	IsHaveDb   bool                  // 是否有数据库
	Auth       *jwt.GinJWTMiddleware // 鉴权
)

var (
	JOB_BUCKET  = []byte("JOB_BUCKET")
	USER_BUCKET = []byte("USER_BUCKET")
)

const TokenHeadName = "Bearer" // Token 认证方式

type JOB_TYPE int64

const (
	JOB_MODIFY JOB_TYPE = iota
	JOB_DEL
	JOB_ADD
)

type JobChan struct {
	model.Job
	Type JOB_TYPE
}

var (
	TASK_CHAN = make(chan JobChan, 10)
)

func JobGo(j model.Job, t JOB_TYPE) {
	TASK_CHAN <- JobChan{
		Job: model.Job{
			Id:     j.Id,
			Name:   j.Name,
			Expr:   j.Expr,
			Mark:   j.Mark,
			Enable: j.Enable,
			Path:   j.Path,
		},
		Type: t,
	}
}

func JobID(id int64) []byte {
	return []byte(fmt.Sprintf("JOB:ID:%d", id))
}

func UserID(id int64) []byte {
	return []byte(fmt.Sprintf("USER:ID:%d", id))
}

// GetInitCore
// 初始化JS插件内容
func GetInitCore() *core.Core {
	vm := core.NewCore(core.OptionLog(Logger))
	return vm
}
