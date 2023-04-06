package initialize

import (
	"fmt"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/config"
	"github.com/issueye/lichee/pkg/logger"
	"github.com/issueye/lichee/utils"
)

func InitLogger() {
	cfg := new(logger.Config)
	cfg.Path = config.GetSysParam("logPath").String()
	cfg.Level = config.GetSysParam("logLevel").Int()
	cfg.MaxAge = config.GetSysParam("logMaxAge").Int()
	cfg.MaxBackups = config.GetSysParam("logMaxBackups").Int()
	cfg.MaxSize = config.GetSysParam("logMaxSize").Int()
	cfg.Compress = config.GetSysParam("logCompress").Bool()
	common.Log, common.Logger = logger.InitLogger(cfg)
	fmt.Printf("【%s】初始化日志完成...\n", utils.Ltime{}.GetNowStr())
}
