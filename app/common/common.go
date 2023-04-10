package common

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/pkg/db"
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
	JOB_BUCKET       = []byte("JOB_BUCKET")
	DB_SOURCE_BUCKET = []byte("DB_SOURCE_BUCKET")
	USER_BUCKET      = []byte("USER_BUCKET")
	AREA_BUCKET      = []byte("AREA_BUCKET")
	PARAM_BUCKET     = []byte("PARAM_BUCKET")
	PARAM_SYS_BYCKET = []byte("PARAM_SYS_BUCKET")

	SYS_USER = int64(10000)
	SYS_AREA = int64(10000)

	SYS_AREA_NAME = "lichee"
)

const (
	WS_LOG_GROUP = "LOG_GROUP"
)

const TokenHeadName = "Bearer" // Token 认证方式

type JOB_TYPE int64

const (
	JOB_MODIFY JOB_TYPE = iota
	JOB_DEL
	JOB_ADD
	JOB_AT_ONCE_RUN
	JOB_DELAY_ONCE_RUN
)

type JobChan struct {
	model.Job
	Delay time.Duration // 延迟时间
	Type  JOB_TYPE
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
			AreaId: j.AreaId,
		},
		Type: t,
	}
}

func JobID(id int64) []byte {
	return []byte(fmt.Sprintf("JOB:ID:%d", id))
}

func DbSourceID(id int64) []byte {
	return []byte(fmt.Sprintf("DB:SOURCE:ID:%d", id))
}

func UserID(id int64) []byte {
	return []byte(fmt.Sprintf("USER:ID:%d", id))
}

func AreaID(id int64) []byte {
	return []byte(fmt.Sprintf("AREA:ID:%d", id))
}

func AreaBucketID(id int64) []byte {
	return []byte(fmt.Sprintf("AREA:BUCKET:ID:%d", id))
}

func ParamID(id int64) []byte {
	return []byte(fmt.Sprintf("AREA:PARAM:ID:%d", id))
}

// GetInitCore
// 初始化JS插件内容
func GetInitCore() *core.Core {
	vm := core.NewCore(core.OptionLog(Logger))
	return vm
}

func GetDb(cfg *db.Config, t int) (*gorm.DB, error) {
	var (
		d   *gorm.DB
		err error
	)

	switch t {
	case 0:
		{
			// 尝试连接
			d, err = db.InitSqlServer(cfg, Log)
			if err != nil {
				Log.Errorf("数据库连接失败，失败原因：%s", err.Error())
				return nil, err
			}
		}
	case 1:
		{
			// 尝试连接
			d, err = db.InitMysql(cfg, Log)
			if err != nil {
				Log.Errorf("数据库连接失败，失败原因：%s", err.Error())
				return nil, err
			}
		}
	case 2:
		{
			// 尝试连接
			d, err = db.InitOracle(cfg, Log)
			if err != nil {
				Log.Errorf("数据库连接失败，失败原因：%s", err.Error())
				return nil, err
			}
		}
	}

	return d, err
}
