package initialize

import (
	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/pkg/ws"
)

func Initialize() {
	// 初始化boltdb
	InitBoltDb()
	// 初始化系统操作员
	InitAdminUser()
	// 初始化系统参数 bucket
	InitSysBucket()
	// 初始化配置文件
	InitConfig()
	// 初始化日志
	InitLogger()
	// 初始化数据库
	InitDB()
	// 初始化ws
	ws.RunWs(common.Log)
	// 初始化JS虚拟机
	InitPlugins()
	// 初始化定时任务
	InitTaskJob()
	// 初始化http服务
	InitHttpServer()
}
