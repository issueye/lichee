package router

import (
	"github.com/gin-gonic/gin"
)

type IRouters interface {
	Register(group *gin.RouterGroup)
}

func InitRouter(r *gin.Engine) {
	name := "api"
	apiName := r.Group(name)
	registerVersionRouter(apiName,
		&AutoJsRouter{}, // js脚本服务
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
