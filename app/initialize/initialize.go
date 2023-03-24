package initialize

func Initialize() {
	InitConfig()
	// 初始化日志
	InitLogger()
	// 初始化数据库
	InitDB()
	// 初始化JS虚拟机
	InitPlugins()
	// 初始化定时任务
	InitTaskJob()
	// 初始化http服务
	InitHttpServer()
}
