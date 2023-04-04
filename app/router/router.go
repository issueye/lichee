package router

import (
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

	// 鉴权API
	apiName.POST("/login", common.Auth.LoginHandler)
	apiName.GET("/refresh", common.Auth.RefreshHandler)
	apiName.GET("/logout", common.Auth.LogoutHandler)

	user := apiName.Group("user")
	{
		user.POST("", v1.NewUserController().Create)
		user.DELETE("/:id", v1.NewUserController().Del)
		user.PUT("", v1.NewUserController().Modify)
		user.GET("", v1.NewUserController().List)
		user.PUT("/status/:id", v1.NewUserController().ModifyStatus)
	}

	// 参数
	param := apiName.Group("param")
	{
		param.POST("", v1.NewParamController().Create)
		param.GET("", v1.NewParamController().List)
		param.PUT("", v1.NewParamController().Modify)
		param.DELETE("/:areaid/:id", v1.NewParamController().Del)
	}

	// 参数域
	area := apiName.Group("area")
	{
		area.POST("", v1.NewParamController().CreateArea)
		area.GET("", v1.NewParamController().AreaList)
		area.DELETE("/:id", v1.NewParamController().DelArea)
		area.PUT("", v1.NewParamController().ModifyArea)
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
