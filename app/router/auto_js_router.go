package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/issueye/lichee/app/controller/v1"
)

type AutoJsRouter struct{}

func (auto *AutoJsRouter) Register(group *gin.RouterGroup) {
	autoJs := &v1.AutoJsController{
		ScriptTimeoutSec: 10,
	}
	group.Any("lichee/*any", autoJs.AutoJsReceiveServer)
}
