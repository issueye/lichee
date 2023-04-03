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
