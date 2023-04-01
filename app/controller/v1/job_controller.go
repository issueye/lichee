package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/issueye/lichee/app/model"
)

type JobController struct{}

func (job *JobController) Create(ctx *gin.Context) {
	req := new(model.ReqCreateJob)
	err := ctx.Bind(req)
	if err != nil {
		return
	}
}
