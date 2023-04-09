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

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

// Create doc
//	@tags			用户管理
//	@Summary		添加用户
//	@Description	添加用户
//	@Produce		json
//	@Param			data	body		model.ReqCreateUser	true	"添加用户"
//	@Success		200		{object}	res.Base			"code: 200 成功"
//	@Failure		500		{object}	res.Base			"错误返回内容"
//	@Router			/api/user [post]
//	@Security		ApiKeyAuth
func (user *UserController) Create(ctx *gin.Context) {
	req := new(model.ReqCreateUser)
	err := ctx.Bind(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	data := new(model.User)
	data.Id = utils.Lid{}.GenID()
	data.Name = req.Name
	data.Account = req.Account
	data.Password = req.Password
	data.Enable = 0
	data.Mark = req.Mark
	data.CreateTime = time.Now()

	err = service.NewUserService().Save(data)
	if err != nil {
		common.Log.Errorf("创建用户信息失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "创建用户信息失败")
		return
	}

	res.SuccessByMsg(ctx, "创建用户成功")
}

// List doc
//	@tags			用户管理
//	@Summary		获取用户列表
//	@Description	获取用户列表
//	@Produce		json
//	@Param			isNotPaging	query		string								false	"是否需要分页， 默认需要， 如果不分页 传 true"
//	@Param			pageNum		query		string								false	"页码， 如果不分页 传 0"
//	@Param			pageSize	query		string								false	"一页大小， 如果不分页 传 0"
//	@Param			name		query		string								false	"任务名称"
//	@Param			account		query		string								false	"任务名称"
//	@Param			mark		query		string								false	"任务描述"
//	@Success		200			{object}	res.Full{data=[]model.ReqQueryUser}	"code: 200 成功"
//	@Failure		500			{object}	res.Base							"错误返回内容"
//	@Router			/api/user [get]
//	@Security		ApiKeyAuth
func (user *UserController) List(ctx *gin.Context) {
	req := new(model.ReqQueryUser)
	err := ctx.BindQuery(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	fmt.Println("请求内容   ", utils.Ljson{}.Struct2Json(req))

	list, err := service.NewUserService().Query(req)
	if err != nil {
		common.Log.Errorf("查询用户信息失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "查询用户信息失败")
		return
	}

	res.SuccessAutoData(ctx, req, list)
}

// Del doc
//	@tags			用户管理
//	@Summary		删除用户
//	@Description	删除用户
//	@Produce		json
//	@Param			id	path		string		true	"用户ID"
//	@Success		200	{object}	res.Base	"code: 200 成功"
//	@Failure		500	{object}	res.Base	"错误返回内容"
//	@Router			/api/user/{id} [delete]
//	@Security		ApiKeyAuth
func (user *UserController) Del(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res.FailByMsg(ctx, "用户ID不能为空")
		return
	}

	// 删除用户
	err := service.NewUserService().Del(cast.ToInt64(id))
	if err != nil {
		common.Log.Errorf("查询用户信息失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "查询用户信息失败")
		return
	}

	res.Success(ctx)
}

// Modify doc
//	@tags			用户管理
//	@Summary		修改用户信息
//	@Description	修改用户信息
//	@Produce		json
//	@Param			data	body		model.ReqModifyUser	true	"修改用户信息"
//	@Success		200		{object}	res.Base			"code: 200 成功"
//	@Failure		500		{object}	res.Base			"错误返回内容"
//	@Router			/api/user [put]
//	@Security		ApiKeyAuth
func (user *UserController) Modify(ctx *gin.Context) {
	req := new(model.ReqModifyUser)
	err := ctx.Bind(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	// 从数据库查询定时任务信息
	u, err := service.NewUserService().GetById(req.Id)
	if err != nil {
		common.Log.Errorf("获取用户信息失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "获取用户信息失败")
		return
	}

	u.Account = req.Account
	u.Name = req.Name
	u.Mark = req.Mark
	u.Password = req.Password

	// 保存定时任务信息
	err = service.NewUserService().Save(u)
	if err != nil {
		common.Log.Errorf("修改用户信息失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "修改用户信息失败")
		return
	}

	res.SuccessByMsg(ctx, "修改用户信息成功")
}

// ModifyStatus doc
//	@tags			用户管理
//	@Summary		停用/启用用户
//	@Description	停用/启用用户
//	@Produce		json
//	@Param			id	path		string		true	"用户ID"
//	@Success		200	{object}	res.Base	"code: 200 成功"
//	@Failure		500	{object}	res.Base	"错误返回内容"
//	@Router			/api/user/status/{id} [put]
//	@Security		ApiKeyAuth
func (user *UserController) ModifyStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res.FailByMsg(ctx, "用户ID不能为空")
		return
	}

	// 从数据库查询用户信息
	u, err := service.NewUserService().GetById(cast.ToInt64(id))
	if err != nil {
		common.Log.Errorf("获取用户信息失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "获取用户信息失败")
		return
	}

	status := u.Enable != 1
	if status {
		u.Enable = 1
	} else {
		u.Enable = 0
	}

	// 保存用户信息
	err = service.NewUserService().Save(u)
	if err != nil {
		common.Log.Errorf("修改用户信息失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "修改用户信息失败")
		return
	}

	message := ""
	if status {
		message = "启用用户成功"
	} else {
		message = "停用用户成功"
	}

	res.SuccessByMsg(ctx, message)
}
