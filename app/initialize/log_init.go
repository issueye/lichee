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
	cfg.Path = config.GetSysParam("log-path").String()
	cfg.Level = config.GetSysParam("log-level").Int()
	cfg.MaxAge = config.GetSysParam("log-max-age").Int()
	cfg.MaxBackups = config.GetSysParam("log-max-backups").Int()
	cfg.MaxSize = config.GetSysParam("log-max-size").Int()
	cfg.Compress = config.GetSysParam("log-compress").Bool()
	common.Log, common.Logger = logger.InitLogger(cfg)
	fmt.Printf("【%s】初始化日志完成...\n", utils.Ltime{}.GetNowStr())
}
