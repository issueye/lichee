package v1

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/app/service"
	"github.com/issueye/lichee/pkg/res"
	"github.com/issueye/lichee/utils"
	"github.com/spf13/cast"
)

type JobController struct{}

func NewJobController() *JobController {
	return &JobController{}
}

// Create doc
//
//	@tags			定时任务管理
//	@Summary		添加定时任务数据
//	@Description	添加定时任务数据
//	@Produce		json
//	@Param			data	body		model.ReqCreateJob	true	"添加定时任务数据"
//	@Success		200		{object}	res.Base			"code: 200 成功"
//	@Failure		500		{object}	res.Base			"错误返回内容"
//	@Router			/api/job [post]
//	@Security		ApiKeyAuth
func (job *JobController) Create(ctx *gin.Context) {
	req := new(model.ReqCreateJob)
	err := ctx.Bind(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	data := new(model.Job)
	data.Name = req.Name
	data.Enable = false
	data.Expr = req.Expr
	data.Id = utils.Lid{}.GenID()
	data.Mark = req.Mark
	data.Path = req.Path
	data.AreaId = req.AreaId // 参数域
	data.CreateTime = time.Now()

	err = service.NewJobService().Save(data)
	if err != nil {
		common.Log.Errorf("创建定时任务失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "创建定时任务失败")
		return
	}

	res.SuccessByMsg(ctx, "创建定时任务成功")
}

// List doc
//
//	@tags			定时任务管理
//	@Summary		获取定时任务列表
//	@Description	获取定时任务列表
//	@Produce		json
//	@Param			isNotPaging	query		string								false	"是否需要分页， 默认需要， 如果不分页 传 true"
//	@Param			pageNum		query		string								false	"页码， 如果不分页 传 0"
//	@Param			pageSize	query		string								false	"一页大小， 如果不分页 传 0"
//	@Param			name		query		string								false	"任务名称"
//	@Param			mark		query		string								false	"任务描述"
//	@Success		200			{object}	res.Full{data=[]model.ResQueryJob}	"code: 200 成功"
//	@Failure		500			{object}	res.Base							"错误返回内容"
//	@Router			/api/job [get]
//	@Security		ApiKeyAuth
func (job *JobController) List(ctx *gin.Context) {
	req := new(model.ReqQueryJob)
	err := ctx.BindQuery(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	fmt.Println("请求内容   ", utils.Ljson{}.Struct2Json(req))

	list, err := service.NewJobService().Query(req)
	if err != nil {
		common.Log.Errorf("查询定时任务失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "查询定时任务失败")
		return
	}

	res.SuccessAutoData(ctx, req, list)
}

// Del doc
//
//	@tags			定时任务管理
//	@Summary		删除定时任务
//	@Description	删除定时任务
//	@Produce		json
//	@Param			id	path		string		true	"任务ID"
//	@Success		200	{object}	res.Base	"code: 200 成功"
//	@Failure		500	{object}	res.Base	"错误返回内容"
//	@Router			/api/job/{id} [delete]
//	@Security		ApiKeyAuth
func (job *JobController) Del(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res.FailByMsg(ctx, "任务ID不能为空")
		return
	}

	// 从数据库查询定时任务信息
	j, err := service.NewJobService().GetById(cast.ToInt64(id))
	if err != nil {
		common.Log.Errorf("获取定时任务失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "获取定时任务失败")
		return
	}

	err = service.NewJobService().Del(cast.ToInt64(id))
	if err != nil {
		common.Log.Errorf("查询定时任务失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "查询定时任务失败")
		return
	}

	if j.Enable {
		common.JobGo(*j, common.JOB_DEL)
	}

	res.Success(ctx)
}

// Modify doc
//
//	@tags			定时任务管理
//	@Summary		修改定时任务数据
//	@Description	修改定时任务数据
//	@Produce		json
//	@Param			data	body		model.ReqModifyJob	true	"修改定时任务数据"
//	@Success		200		{object}	res.Base			"code: 200 成功"
//	@Failure		500		{object}	res.Base			"错误返回内容"
//	@Router			/api/job [put]
//	@Security		ApiKeyAuth
func (job *JobController) Modify(ctx *gin.Context) {
	req := new(model.ReqModifyJob)
	err := ctx.Bind(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	// 从数据库查询定时任务信息
	j, err := service.NewJobService().GetById(req.Id)
	if err != nil {
		common.Log.Errorf("获取定时任务失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "获取定时任务失败")
		return
	}

	j.Name = req.Name
	j.Expr = req.Expr
	j.Mark = req.Mark
	j.Path = req.Path
	j.AreaId = req.AreaId

	// 保存定时任务信息
	err = service.NewJobService().Save(j)
	if err != nil {
		common.Log.Errorf("修改定时任务失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "修改定时任务失败")
		return
	}

	if j.Enable {
		common.JobGo(*j, common.JOB_ADD)
	}

	res.SuccessByMsg(ctx, "修改定时任务成功")
}

// AtOnceRun doc
//
//	@tags			定时任务管理
//	@Summary		马上运行一次任务
//	@Description	马上运行一次任务
//	@Produce		json
//	@Param			id	path		string		true	"任务ID"
//	@Success		200	{object}	res.Base	"code: 200 成功"
//	@Failure		500	{object}	res.Base	"错误返回内容"
//	@Router			/api/job/atOnceRun/{id} [get]
//	@Security		ApiKeyAuth
func (job *JobController) AtOnceRun(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res.FailByMsg(ctx, "任务ID不能为空")
		return
	}

	// 从数据库查询定时任务信息
	j, err := service.NewJobService().GetById(cast.ToInt64(id))
	if err != nil {
		common.Log.Errorf("获取定时任务失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "获取定时任务失败")
		return
	}

	// 运行一次任务
	common.JobGo(*j, common.JOB_AT_ONCE_RUN)

	res.SuccessByMsg(ctx, "将任务添加到队列中成功，即将运行任务")
}

// ModifyStatus doc
//
//	@tags			定时任务管理
//	@Summary		修改定时任务数据
//	@Description	修改定时任务数据
//	@Produce		json
//	@Param			id	path		string		true	"任务ID"
//	@Success		200	{object}	res.Base	"code: 200 成功"
//	@Failure		500	{object}	res.Base	"错误返回内容"
//	@Router			/api/job/status/{id} [put]
//	@Security		ApiKeyAuth
func (job *JobController) ModifyStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res.FailByMsg(ctx, "任务ID不能为空")
		return
	}

	// 从数据库查询定时任务信息
	j, err := service.NewJobService().GetById(cast.ToInt64(id))
	if err != nil {
		common.Log.Errorf("获取定时任务失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "获取定时任务失败")
		return
	}

	status := !j.Enable
	j.Enable = status

	// 保存定时任务信息
	err = service.NewJobService().Save(j)
	if err != nil {
		common.Log.Errorf("修改定时任务失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "修改定时任务失败")
		return
	}

	message := ""
	if status {
		message = "启用定时任务成功"
		common.JobGo(*j, common.JOB_ADD)
	} else {
		message = "停用定时任务成功"
		common.JobGo(*j, common.JOB_DEL)
	}

	res.SuccessByMsg(ctx, message)
}
