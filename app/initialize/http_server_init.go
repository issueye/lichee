package initialize

import (
	"fmt"
	"mime"
	"net/http"
	"runtime"

	"github.com/issueye/lichee/app/global"
	"github.com/issueye/lichee/app/router"
	"github.com/issueye/lichee/pkg/middleware"

	"github.com/dimiro1/banner"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

func InitHttpServer() {
	gin.SetMode(gin.ReleaseMode)
	global.Router = gin.New()

	// 加载中间件
	global.Router.Use(middleware.GinLogger(global.Logger))
	global.Router.Use(middleware.GinRecovery(global.Logger, false))
	global.Router.Use(middleware.CORSMiddleware([]string{}))

	// 设置一个静态文件服务器
	global.Router.Static("/www", "./static")

	// 告诉服务文件的MIME类型
	_ = mime.AddExtensionType(".js", "application/javascript")
	_ = mime.AddExtensionType(".css", "text/css")
	_ = mime.AddExtensionType(".woff", "application/font-woff")
	_ = mime.AddExtensionType(".woff2", "application/font-woff2")
	_ = mime.AddExtensionType(".ttf", "application/font-ttf")
	_ = mime.AddExtensionType(".eot", "application/vnd.ms-fontobject")

	// 初始化路由
	router.InitRouter(global.Router)
	global.HttpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", global.LocalCfg.LocalPort),
		Handler: global.Router,
	}

	ShowInfo()
	err := global.HttpServer.ListenAndServe()
	if err != nil {
		panic("http服务开启失败，失败原因：" + err.Error())
	}
}

func ShowInfo() {
	bannerStr := `{{ .Title "GO-PLUGINS" "" 4 }}`
	banner.InitString(colorable.NewColorableStdout(), true, true, bannerStr)

	info := `
	ＧＯ　版本: %s
	运行　平台: %s
	平台　架构: %s
	处理器数量: %d
	程序　版本: V0.1.0`
	fmt.Printf(info,
		runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.NumCPU())

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
}
