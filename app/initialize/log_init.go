package initialize

import (
	"fmt"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/pkg/logger"
	"github.com/issueye/lichee/utils"
)

func InitLogger() {
	common.Log, common.Logger = logger.InitLogger(common.LocalCfg.Log)
	fmt.Printf("【%s】初始化日志完成...\n", utils.Ltime{}.GetNowStr())
}
