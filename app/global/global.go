package global

import (
	"net/http"

	"github.com/issueye/lichee/app/model"

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

	LocalCfg   *model.Config // 本地配置信息
	ConfigPath string        // 配置文件路径
	IsHaveDb   bool          // 是否有数据库
)
