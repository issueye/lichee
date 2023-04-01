package v1

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/app/service"
	"github.com/issueye/lichee/utils"
	"github.com/spf13/cast"
)

type JobController struct{}

func NewJobController() *JobController {
	return &JobController{}
}

// Create doc
// @tags        定时任务管理
// @Summary     添加定时任务数据
// @Description 添加定时任务数据
// @Produce     json
// @Param       data body     model.ReqCreateJob true "添加定时任务数据"
// @Success     200  {object} common.Base        true "code: 200 成功"
// @Failure     500  {object} common.Base        true "错误返回内容"
// @Router      /api/job [post]
// @Security    ApiKeyAuth
func (job *JobController) Create(ctx *gin.Context) {
	req := new(model.ReqCreateJob)
	err := ctx.Bind(req)
	if err != nil {
		common.FailBind(ctx, err)
		return
	}

	data := new(model.Job)
	data.Name = req.Name
	data.Enable = false
	data.Expr = req.Expr
	data.Id = utils.Lid{}.GenID()
	data.Mark = req.Mark
	data.CreateTime = utils.LongDateTime{Time: time.Now()}

	err = service.NewJobService().Save(data)
	if err != nil {
		common.Log.Errorf("创建定时任务失败，失败原因：%s", err.Error())
		common.FailByMsg(ctx, "创建定时任务失败")
		return
	}

	common.SuccessByMsg(ctx, "创建定时任务成功")
}

// List doc
// @tags        定时任务管理
// @Summary     获取定时任务列表
// @Description 获取定时任务列表
// @Produce     json
// @Param       isNotPaging query    string                        false "是否需要分页， 默认需要， 如果不分页 传 true"
// @Param       pageNum     query    string                        true  "页码， 如果不分页 传 0"
// @Param       pageSize    query    string                        true  "一页大小， 如果不分页 传 0"
// @Param       name        query    string                        false "任务名称"
// @Param       mark        query    string                        false "任务描述"
// @Success     200         {object} common.Full{data=[]model.Job} true  "code: 200 成功"
// @Failure     500         {object} common.Base                   true  "错误返回内容"
// @Router      /api/job [get]
// @Security    ApiKeyAuth
func (job *JobController) List(ctx *gin.Context) {
	req := new(model.ReqQueryJob)
	err := ctx.BindQuery(req)
	if err != nil {
		common.FailBind(ctx, err)
		return
	}

	fmt.Println("请求内容   ", utils.Ljson{}.Struct2Json(req))

	list, err := service.NewJobService().Query(req)
	if err != nil {
		common.Log.Errorf("查询定时任务失败，失败原因：%s", err.Error())
		common.FailByMsg(ctx, "查询定时任务失败")
		return
	}

	common.SuccessAutoData(ctx, req, list)
}

// Del doc
// @tags        定时任务管理
// @Summary     删除定时任务
// @Description 删除定时任务
// @Produce     json
// @Param       id  path     string      true "任务ID"
// @Success     200 {object} common.Base true "code: 200 成功"
// @Failure     500 {object} common.Base true "错误返回内容"
// @Router      /api/job [delete]
// @Security    ApiKeyAuth
func (job *JobController) Del(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		common.FailByMsg(ctx, "任务ID不能为空")
		return
	}

	err := service.NewJobService().Del(cast.ToInt64(id))
	if err != nil {
		common.Log.Errorf("查询定时任务失败，失败原因：%s", err.Error())
		common.FailByMsg(ctx, "查询定时任务失败")
		return
	}

	common.Success(ctx)
}

// Modify doc
// @tags        定时任务管理
// @Summary     修改定时任务数据
// @Description 修改定时任务数据
// @Produce     json
// @Param       data body     model.ReqModifyJob true "修改定时任务数据"
// @Success     200  {object} common.Base        true "code: 200 成功"
// @Failure     500  {object} common.Base        true "错误返回内容"
// @Router      /api/job [put]
// @Security    ApiKeyAuth
func (job *JobController) Modify(ctx *gin.Context) {
	req := new(model.ReqModifyJob)
	err := ctx.Bind(req)
	if err != nil {
		common.FailBind(ctx, err)
		return
	}

	// 从数据库查询定时任务信息
	j, err := service.NewJobService().GetById(req.Id)
	if err != nil {
		common.Log.Errorf("获取定时任务失败，失败原因：%s", err.Error())
		common.FailByMsg(ctx, "获取定时任务失败")
		return
	}

	j.Name = req.Name
	j.Expr = req.Expr
	j.Mark = req.Mark

	// 保存定时任务信息
	err = service.NewJobService().Save(j)
	if err != nil {
		common.Log.Errorf("修改定时任务失败，失败原因：%s", err.Error())
		common.FailByMsg(ctx, "修改定时任务失败")
		return
	}

	common.SuccessByMsg(ctx, "修改定时任务成功")
}
