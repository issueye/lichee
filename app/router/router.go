package router

import (
	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
	"github.com/issueye/lichee/app/common"
	v1 "github.com/issueye/lichee/app/controller/v1"
	"github.com/issueye/lichee/app/middleware"
)

type IRouters interface {
	Register(group *gin.RouterGroup)
}

func InitRouter(r *gin.Engine) {
	name := "api"
	common.Auth = middleware.NewAuth()
	apiName := r.Group(name)

	r.GET("socket", v1.WsLogMonitor)

	r.GET("/debug/statsviz/*filepath", func(context *gin.Context) {
		if context.Param("filepath") == "/ws" {
			statsviz.Ws(context.Writer, context.Request)
			return
		}
		statsviz.IndexAtRoot("/debug/statsviz").ServeHTTP(context.Writer, context.Request)
	})

	// 鉴权API
	apiName.POST("/login", common.Auth.LoginHandler)
	apiName.GET("/refresh", common.Auth.RefreshHandler)
	apiName.GET("/logout", common.Auth.LogoutHandler)

	user := apiName.Group("user", common.Auth.MiddlewareFunc())
	{
		user.POST("", v1.NewUserController().Create)
		user.DELETE("/:id", v1.NewUserController().Del)
		user.PUT("", v1.NewUserController().Modify)
		user.GET("", v1.NewUserController().List)
		user.PUT("/status/:id", v1.NewUserController().ModifyStatus)
	}

	// 参数
	param := apiName.Group("param", common.Auth.MiddlewareFunc())
	{
		param.POST("", v1.NewParamController().Create)
		param.GET("", v1.NewParamController().List)
		param.PUT("", v1.NewParamController().Modify)
		param.DELETE("/:areaid/:id", v1.NewParamController().Del)
	}

	// 参数域
	area := apiName.Group("area", common.Auth.MiddlewareFunc())
	{
		area.POST("", v1.NewParamController().CreateArea)
		area.GET("", v1.NewParamController().AreaList)
		area.DELETE("/:id", v1.NewParamController().DelArea)
		area.PUT("", v1.NewParamController().ModifyArea)
	}

	// 数据库源
	dbSource := apiName.Group("dbSource", common.Auth.MiddlewareFunc())
	{
		dbSource.POST("", v1.NewDbController().Create)
		dbSource.GET("", v1.NewDbController().List)
		dbSource.DELETE("/:id", v1.NewDbController().Del)
		dbSource.PUT("", v1.NewDbController().Modify)
		dbSource.POST("/testLink", v1.NewDbController().TestLink)
	}

	registerVersionRouter(apiName,
		&AutoJsRouter{}, // js脚本服务
		&JobRouter{},    // 定时任务
	)

	registerStatic(&r.RouterGroup)
}

func registerStatic(r *gin.RouterGroup, iRouters ...IRouters) {
	for _, iRouter := range iRouters {
		iRouter.Register(r)
	}
}

// registerRouter 注册路由
func registerVersionRouter(r *gin.RouterGroup, iRouters ...IRouters) {
	for _, iRouter := range iRouters {
		iRouter.Register(r)
	}
}
