package initialize

import (
	"fmt"
	"mime"
	"net"
	"net/http"
	"runtime"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/config"
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
	mode := config.GetSysParam("sysRunMode").String()
	gin.SetMode(mode)
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
	// 显示 banner
	ShowInfo()
	// 启动 监听
	ListenAndServe(fmt.Sprintf(":%d", common.LocalCfg.LocalPort), r)
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func ListenAndServe(addr string, handler http.Handler) error {
	srv := &http.Server{Addr: addr, Handler: handler}
	addr = srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp4", addr) // 仅指定 IPv4
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
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
