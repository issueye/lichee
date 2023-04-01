package initialize

import (
	"fmt"
	"mime"
	"net/http"
	"runtime"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/router"
	"github.com/issueye/lichee/pkg/middleware"
	orange_validator "github.com/issueye/lichee/pkg/validator"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/dimiro1/banner"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

func InitHttpServer() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	orange_validator.RegisterValidator()

	// 加载中间件
	r.Use(middleware.Req())
	r.Use(middleware.GinLogger(common.Logger))
	r.Use(middleware.GinRecovery(common.Logger, true))
	r.Use(middleware.CORSMiddleware([]string{}))

	// 设置一个静态文件服务器
	r.Static("/www", "./static")

	// 告诉服务文件的MIME类型
	_ = mime.AddExtensionType(".js", "application/javascript")
	_ = mime.AddExtensionType(".css", "text/css")
	_ = mime.AddExtensionType(".woff", "application/font-woff")
	_ = mime.AddExtensionType(".woff2", "application/font-woff2")
	_ = mime.AddExtensionType(".ttf", "application/font-ttf")
	_ = mime.AddExtensionType(".eot", "application/vnd.ms-fontobject")

	// 设置 swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 初始化路由
	router.InitRouter(r)

	common.Router = r
	common.HttpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", common.LocalCfg.LocalPort),
		Handler: common.Router,
	}

	ShowInfo()
	err := common.HttpServer.ListenAndServe()
	if err != nil {
		panic("http服务开启失败，失败原因：" + err.Error())
	}
}

func ShowInfo() {
	bannerStr := `{{ .Title "lichee" "" 4 }}`
	banner.InitString(colorable.NewColorableStdout(), true, true, bannerStr)

	info := `
	ＧＯ　版本: %s
	运行　平台: %s
	平台　架构: %s
	处理器数量: %d
	程序　版本: V0.1.0
	运行端口号：%d`
	fmt.Printf(
		info,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		common.LocalCfg.LocalPort)

	fmt.Println("")
}
