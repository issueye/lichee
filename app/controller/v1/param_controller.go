package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/app/service"
	"github.com/issueye/lichee/pkg/res"
	"github.com/issueye/lichee/utils"
	"github.com/spf13/cast"
)

type ParamController struct{}

func NewParamController() *ParamController {
	return &ParamController{}
}

// Create doc
//
//	@tags			参数管理
//	@Summary		添加参数
//	@Description	添加参数
//	@Produce		json
//	@Param			data	body		model.ReqCreateParam	true	"添加参数"
//	@Success		200		{object}	res.Base				"code: 200 成功"
//	@Failure		500		{object}	res.Base				"错误返回内容"
//	@Router			/api/param [post]
//	@Security		ApiKeyAuth
func (control *ParamController) Create(ctx *gin.Context) {
	req := new(model.ReqCreateParam)
	err := ctx.Bind(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	area, err := service.NewParamService().GetAreaById(req.AreaId)
	if err != nil {
		common.Log.Errorf("根据参数域ID参数参数域信息失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "根据参数域ID参数参数域信息失败")
		return
	}

	data := new(model.Param)
	data.Id = utils.Lid{}.GenID()
	data.Name = req.Name
	data.AreaId = req.AreaId
	data.Area = area.Name
	data.Mark = req.Mark
	data.Value = req.Value

	err = service.NewParamService().Save(data)
	if err != nil {
		common.Log.Errorf("创建参数失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "创建参数失败")
		return
	}

	res.SuccessByMsg(ctx, "创建参数成功")
}

// Del doc
//
//	@tags			参数管理
//	@Summary		删除参数
//	@Description	删除参数
//	@Produce		json
//	@Param			areaid	path		string		true	"参数域ID"
//	@Param			id		path		string		true	"参数ID"
//	@Success		200		{object}	res.Base	"code: 200 成功"
//	@Failure		500		{object}	res.Base	"错误返回内容"
//	@Router			/api/param/{areaid}/{id} [delete]
//	@Security		ApiKeyAuth
func (control *ParamController) Del(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res.FailByMsg(ctx, "参数ID不能为空")
		return
	}

	areaid := ctx.Param("areaid")
	if areaid == "" {
		res.FailByMsg(ctx, "参数域ID不能为空")
		return
	}

	// 删除用户
	err := service.NewParamService().DelParam(cast.ToInt64(areaid), cast.ToInt64(id))
	if err != nil {
		common.Log.Errorf("查询用户信息失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "查询用户信息失败")
		return
	}

	res.Success(ctx)
}

// Modify doc
//
//	@tags			参数管理
//	@Summary		修改参数信息
//	@Description	修改参数信息
//	@Produce		json
//	@Param			data	body		model.ReqModifyParam	true	"修改参数信息"
//	@Success		200		{object}	res.Base				"code: 200 成功"
//	@Failure		500		{object}	res.Base				"错误返回内容"
//	@Router			/api/param [put]
//	@Security		ApiKeyAuth
func (control *ParamController) Modify(ctx *gin.Context) {
	req := new(model.ReqModifyParam)
	err := ctx.Bind(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	// 查询参数域
	pa, err := service.NewParamService().GetAreaById(req.AreaId)
	if err != nil {
		common.Log.Errorf("根据参数域编码【%s】查找参数域信息失败，失败原因：%s", req.AreaId, err.Error())
		res.FailByMsgf(ctx, "根据参数域编码【%s】查找参数域信息失败", req.AreaId)
		return
	}

	// 从数据库查询定时任务信息
	data := new(model.Param)
	data.Id = req.Id
	data.Name = req.Name
	data.Value = req.Value
	data.Mark = req.Mark
	data.AreaId = req.AreaId
	data.Area = pa.Name

	err = service.NewParamService().ModifyParam(data)
	if err != nil {
		common.Log.Errorf("修改参数信息失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "修改参数信息失败")
		return
	}

	res.SuccessByMsg(ctx, "修改参数信息成功")
}

// List doc
//
//	@tags			参数管理
//	@Summary		获取参数列表
//	@Description	获取参数列表
//	@Produce		json
//	@Param			isNotPaging	query		string							false	"是否需要分页， 默认需要， 如果不分页 传 true"
//	@Param			pageNum		query		string							false	"页码， 如果不分页 传 0"
//	@Param			pageSize	query		string							false	"一页大小， 如果不分页 传 0"
//	@Param			name		query		string							false	"查询内容"
//	@Param			area_id		query		int								true	"参数域ID"
//	@Success		200			{object}	res.Full{data=[]model.Param}	"code: 200 成功"
//	@Failure		500			{object}	res.Base						"错误返回内容"
//	@Router			/api/param [get]
//	@Security		ApiKeyAuth
func (control *ParamController) List(ctx *gin.Context) {
	req := new(model.ReqQueryParam)
	err := ctx.BindQuery(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	// 如果没有参数域的话，就直接查系统参数域
	if req.AreaId == 0 {
		req.AreaId = 10000
	}

	list, err := service.NewParamService().Query(req)
	if err != nil {
		common.Log.Errorf("查询参数列表失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "查询参数列表失败")
		return
	}

	res.SuccessAutoData(ctx, req, list)
}

/**************************area*****************************/

// Create doc
//
//	@tags			参数管理
//	@Summary		添加参数
//	@Description	添加参数
//	@Produce		json
//	@Param			data	body		model.ReqCreateArea	true	"添加参数"
//	@Success		200		{object}	res.Base			"code: 200 成功"
//	@Failure		500		{object}	res.Base			"错误返回内容"
//	@Router			/api/area [post]
//	@Security		ApiKeyAuth
func (control *ParamController) CreateArea(ctx *gin.Context) {
	req := new(model.ReqCreateArea)
	err := ctx.Bind(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	area := new(model.ParamArea)
	area.Id = utils.Lid{}.GenID()
	area.Name = req.Name
	area.Mark = req.Mark
	err = service.NewParamService().CreateArea(area)
	if err != nil {
		common.Log.Errorf("创建参数域失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "创建参数域失败")
		return
	}

	res.SuccessByMsg(ctx, "创建参数域成功")
}

// AreaList doc
//
//	@tags			参数管理
//	@Summary		获取参数域列表
//	@Description	获取参数域列表
//	@Produce		json
//	@Param			isNotPaging	query		string								false	"是否需要分页， 默认需要， 如果不分页 传 true"
//	@Param			pageNum		query		string								false	"页码， 如果不分页 传 0"
//	@Param			pageSize	query		string								false	"一页大小， 如果不分页 传 0"
//	@Param			name		query		string								false	"检索内容"
//	@Success		200			{object}	res.Full{data=[]model.ParamArea}	"code: 200 成功"
//	@Failure		500			{object}	res.Base							"错误返回内容"
//	@Router			/api/area [get]
//	@Security		ApiKeyAuth
func (control *ParamController) AreaList(ctx *gin.Context) {
	req := new(model.ReqQueryArea)
	err := ctx.BindQuery(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	list, err := service.NewParamService().GetAreaList(req)
	if err != nil {
		common.Log.Errorf("查询参数域列表失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "查询参数域列表失败")
		return
	}

	res.SuccessAutoData(ctx, req, list)
}

// DelArea doc
//
//	@tags			参数管理
//	@Summary		获取参数域列表
//	@Description	获取参数域列表
//	@Produce		json
//	@Param			id	path		string		true	"参数域ID"
//	@Success		200	{object}	res.Base	"code: 200 成功"
//	@Failure		500	{object}	res.Base	"错误返回内容"
//	@Router			/api/area/{id} [DELETE]
//	@Security		ApiKeyAuth
func (control *ParamController) DelArea(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res.FailByMsg(ctx, "参数域ID不能为空")
		return
	}

	err := service.NewParamService().DelArea(cast.ToInt64(id))
	if err != nil {
		common.Log.Errorf("删除参数域失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "删除参数域失败")
		return
	}

	res.SuccessByMsg(ctx, "删除参数域成功")
}

// ModifyArea doc
//
//	@tags			参数管理
//	@Summary		修改参数域信息
//	@Description	修改参数域信息
//	@Produce		json
//	@Param			data	body		model.ReqModifyArea	true	"修改参数域信息"
//	@Success		200		{object}	res.Base			"code: 200 成功"
//	@Failure		500		{object}	res.Base			"错误返回内容"
//	@Router			/api/area [put]
//	@Security		ApiKeyAuth
func (control *ParamController) ModifyArea(ctx *gin.Context) {
	req := new(model.ReqModifyArea)
	err := ctx.Bind(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	data := new(model.ParamArea)
	data.Id = req.Id
	data.Name = req.Name
	data.Mark = req.Mark
	err = service.NewParamService().ModifyArea(data)
	if err != nil {
		common.Log.Errorf("修改参数域失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "修改参数域失败")
		return
	}

	res.SuccessByMsg(ctx, "修改参数域成功")
}
