package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/issueye/lichee/app/controller/v1"
)

type JobRouter struct{}

func (job *JobRouter) Register(group *gin.RouterGroup) {
	r := v1.NewJobController()
	g := group.Group("job")
	{
		g.POST("", r.Create)
		g.DELETE("/:id", r.Del)
		g.PUT("", r.Modify)
		g.GET("", r.List)
		g.GET("atOnceRun/:id", r.AtOnceRun)
		g.PUT("/status/:id", r.ModifyStatus)
	}
}
