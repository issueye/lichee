package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/app/service"
	"github.com/issueye/lichee/global"
	"github.com/issueye/lichee/pkg/db"
	licheeDB "github.com/issueye/lichee/pkg/plugins/core/db"
	"github.com/issueye/lichee/pkg/res"
	"github.com/issueye/lichee/utils"
	"github.com/spf13/cast"
)

type DbController struct{}

func NewDbController() *DbController {
	return &DbController{}
}

// Create doc
//
//	@tags			数据库源管理
//	@Summary		添加数据库源
//	@Description	添加数据库源
//	@Produce		json
//	@Param			data	body		model.ReqCreateDbSource	true	"添加参数"
//	@Success		200		{object}	res.Base				"code: 200 成功"
//	@Failure		500		{object}	res.Base				"错误返回内容"
//	@Router			/api/dbSource [post]
//	@Security		ApiKeyAuth
func (control *DbController) Create(ctx *gin.Context) {
	req := new(model.ReqCreateDbSource)
	err := ctx.Bind(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	ds := new(model.DbSource)
	ds.Id = utils.Lid{}.GenID()
	ds.Name = req.Name
	ds.Host = req.Host
	ds.Port = req.Port
	ds.Database = req.Database
	ds.User = req.User
	ds.Password = req.Password
	ds.Type = req.Type
	ds.Mark = req.Mark

	err = service.NewDbSourceService().Create(ds)
	if err != nil {
		common.Log.Errorf("添加数据库源失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "添加数据库源失败")
		return
	}

	// 尝试连接数据库
	cfg := new(db.Config)
	cfg.Host = req.Host
	cfg.Port = int(req.Port)
	cfg.Database = req.Database
	cfg.Username = req.User
	cfg.Password = req.Password

	d, err := common.GetDb(cfg, int(req.Type))
	if err != nil {
		common.Log.Errorf("数据库连接失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "数据库连接失败")
		return
	}

	info := &global.DbInfo{
		Cfg:  cfg,
		Name: ds.Name,
		DB:   d,
	}

	// 将数据库注册到JS虚拟机
	licheeDB.RegisterDB(fmt.Sprintf("db/%s", info.Name), info.DB)
	global.GdbMap[ds.Name] = info

	res.SuccessByMsg(ctx, "添加数据库源成功")
}

// List doc
//
//	@tags			数据库源管理
//	@Summary		获取数据库源列表
//	@Description	获取数据库源列表
//	@Produce		json
//	@Param			isNotPaging	query		string									false	"是否需要分页， 默认需要， 如果不分页 传 true"
//	@Param			pageNum		query		string									false	"页码， 如果不分页 传 0"
//	@Param			pageSize	query		string									false	"一页大小， 如果不分页 传 0"
//	@Param			name		query		string									false	"数据源名称"
//	@Success		200			{object}	res.Full{data=[]model.ReqQueryDbSource}	"code: 200 成功"
//	@Failure		500			{object}	res.Base								"错误返回内容"
//	@Router			/api/dbSource [get]
//	@Security		ApiKeyAuth
func (control *DbController) List(ctx *gin.Context) {
	req := new(model.ReqQueryDbSource)
	err := ctx.BindQuery(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	fmt.Println("请求内容   ", utils.Ljson{}.Struct2Json(req))

	list, err := service.NewDbSourceService().Query(req)
	if err != nil {
		common.Log.Errorf("获取数据库源列表失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "获取数据库源列表失败")
		return
	}

	res.SuccessAutoData(ctx, req, list)
}

// Del doc
//
//	@tags			数据库源管理
//	@Summary		删除数据库源
//	@Description	删除数据库源
//	@Produce		json
//	@Param			id	path		string		true	"ID"
//	@Success		200	{object}	res.Base	"code: 200 成功"
//	@Failure		500	{object}	res.Base	"错误返回内容"
//	@Router			/api/dbSource/{id} [delete]
//	@Security		ApiKeyAuth
func (control *DbController) Del(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res.FailByMsg(ctx, "数据库源ID不能为空")
		return
	}

	// 先查询数据
	ds, err := service.NewDbSourceService().GetById(cast.ToInt64(id))
	if err != nil {
		common.Log.Errorf("查询数据库源失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "查询数据库源失败")
		return
	}

	// 从数据库查询定时任务信息
	err = service.NewDbSourceService().Del(cast.ToInt64(id))
	if err != nil {
		common.Log.Errorf("删除数据库源失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "删除数据库源失败")
		return
	}

	// 删除全局中的对象
	_, ok := global.GdbMap[ds.Name]
	if ok {
		delete(global.GdbMap, ds.Name)
	}

	res.Success(ctx)
}

// Modify doc
//
//	@tags			数据库源管理
//	@Summary		修改数据库源数据
//	@Description	修改数据库源数据
//	@Produce		json
//	@Param			data	body		model.ReqModifyDbSource	true	"修改定时任务数据"
//	@Success		200		{object}	res.Base				"code: 200 成功"
//	@Failure		500		{object}	res.Base				"错误返回内容"
//	@Router			/api/dbSource [put]
//	@Security		ApiKeyAuth
func (control *DbController) Modify(ctx *gin.Context) {
	req := new(model.ReqModifyDbSource)
	err := ctx.Bind(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	// 修改数据库源信息
	db, err := service.NewDbSourceService().GetById(req.Id)
	if err != nil {
		common.Log.Errorf("修改数据库源失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "修改数据库源失败")
		return
	}

	db.Name = req.Name
	db.Host = req.Host
	db.Port = req.Port
	db.Database = req.Database
	db.User = req.User
	db.Password = req.Password
	db.Type = req.Type
	db.Mark = req.Mark

	// 保存定时任务信息
	err = service.NewDbSourceService().Create(db)
	if err != nil {
		common.Log.Errorf("修改数据库源失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "修改数据库源失败")
		return
	}

	res.SuccessByMsg(ctx, "修改数据库源成功")
}

// TestLink doc
//
//	@tags			数据库源管理
//	@Summary		修改数据库源数据
//	@Description	修改数据库源数据
//	@Produce		json
//	@Param			data	body		model.ReqModifyDbSource	true	"修改定时任务数据"
//	@Success		200		{object}	res.Base				"code: 200 成功"
//	@Failure		500		{object}	res.Base				"错误返回内容"
//	@Router			/api/dbSource/testLink [post]
//	@Security		ApiKeyAuth
func (control *DbController) TestLink(ctx *gin.Context) {
	req := new(model.ReqModifyDbSource)
	err := ctx.Bind(req)
	if err != nil {
		res.FailBind(ctx, err)
		return
	}

	// 尝试连接数据库
	cfg := new(db.Config)
	cfg.Host = req.Host
	cfg.Port = int(req.Port)
	cfg.Database = req.Database
	cfg.Username = req.User
	cfg.Password = req.Password

	_, err = common.GetDb(cfg, int(req.Type))
	if err != nil {
		common.Log.Errorf("数据库连接失败，失败原因：%s", err.Error())
		res.FailByMsg(ctx, "数据库连接失败")
		return
	}

	res.SuccessByMsg(ctx, "测试连接成功")
}
