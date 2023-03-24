package initialize

import (
	"github.com/issueye/lichee/app/global"
	"github.com/issueye/lichee/pkg/logger"
)

func InitLogger() {
	global.Log, global.Logger = logger.InitLogger(global.LocalCfg.Log)
}
