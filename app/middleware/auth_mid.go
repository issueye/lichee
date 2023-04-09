package middleware

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/app/service"
	"github.com/issueye/lichee/pkg/middleware"
	"github.com/issueye/lichee/pkg/res"
	"github.com/issueye/lichee/utils"
)

type Auth struct{}

func NewAuth() *jwt.GinJWTMiddleware {
	auth := new(Auth)
	jwt, err := middleware.InitAuth(auth)
	if err != nil {
		panic(fmt.Sprintf("初始化鉴权中间件失败，失败原因：%s", err.Error()))
	}
	return jwt
}

// PayloadFunc
// 有效载荷处理
func (auth *Auth) PayloadFunc(data interface{}) jwt.MapClaims {
	mapClaims := make(jwt.MapClaims)
	v, ok := data.(map[string]interface{})
	if ok {
		user := new(model.User)
		// 将用户json转为结构体
		utils.Ljson{}.JsonI2Struct(v["user"], user)

		mapClaims[jwt.IdentityKey] = user.Id
		mapClaims["user"] = v["user"]
	}
	return mapClaims
}

// IdentityHandler
// 解析Claims
func (auth *Auth) IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	// 此处返回值类型 map[string]interface{}
	// 与 payloadFunc 和 authorizator 的 data 类型必须一致, 否则会导致授权失败还不容易找到原因
	mapData := make(map[string]interface{})
	mapData["IdentityKey"] = claims[jwt.IdentityKey]
	mapData["user"] = claims["user"]
	return mapData
}

// Login godoc
//	@tags			基本接口
//	@Summary		用户登录
//	@Produce		json
//	@Description	```
//	@Description	用户登录
//	@Description	```
//	@Param			data	body		model.LoginUser	true	"登录信息"
//	@Success		200		{object}	res.Full
//	@Failure		500		{object}	res.Base	"错误返回内容"
//	@Router			/api/login [post]
func (auth *Auth) Login(c *gin.Context) (interface{}, error) {
	req := new(model.LoginUser)
	// 请求json绑定
	err := c.ShouldBind(req)
	if err != nil {
		return "", err
	}

	user, err := auth.UserAuth(req)
	if err != nil {
		return nil, err
	}

	if user.Enable == 0 {
		return nil, errors.New("当前账户已停用")
	}

	// 修改用户登录时间
	user.LoginTime = time.Now()
	err = service.NewUserService().Save(user)
	if err != nil {
		return nil, err
	}

	MapData := make(map[string]interface{})
	MapData["user"] = utils.Ljson{}.Struct2Json(user)
	// 将用户信息写入到上下文中，在后面登录成功处理时，需要用到
	c.Set("user", user)
	// 将用户以json格式写入, payloadFunc/authorizator会使用到
	return MapData, nil
}

// Authorizator
// 用户登录校验成功处理
func (auth *Auth) Authorizator(data interface{}, c *gin.Context) bool {
	v, ok := data.(map[string]interface{})
	if ok {
		userStr := v["user"].(string)
		user := new(model.User)
		// 将用户json转为结构体
		utils.Ljson{}.Json2Struct(userStr, &user)
		// 将用户保存到context, api调用时取数据方便
		c.Set("user", user)
		return true
	}
	return false
}

type JwtToken struct {
	ID      int64  `json:"id"`      // id
	UID     string `json:"uid"`     // 用户ID
	Name    string `json:"name"`    // 用户名
	Token   string `json:"token"`   // token
	Expires string `json:"expires"` // 时间
}

// Unauthorized
// 用户登录校验失败处理
func (auth *Auth) Unauthorized(ctx *gin.Context, code int, message string) {
	common.Log.Debugf("JWT认证失败, 错误码: %d, 错误信息: %s", code, message)
	res.FailByMsgAndCode(ctx, fmt.Sprintf("认证失败，错误原因：%s", message), res.UNAUTHORIZED)
	ctx.Abort()
}

// LoginResponse
// 登录成功后的响应
func (auth *Auth) LoginResponse(ctx *gin.Context, _ int, token string, expires time.Time) {
	jt := new(JwtToken)
	jt.Token = strings.Join([]string{common.TokenHeadName, token}, " ")
	jt.Expires = expires.Format(utils.FormatDateTimeMs)
	user, err := auth.GetUser(ctx)
	if err != nil {
		common.Log.Errorf("获取用户信息失败，失败原因：%s", err.Error())
		return
	}

	jt.ID = user.Id
	jt.UID = user.Account
	jt.Name = user.Name
	res.SuccessData(ctx, jt)
}

// LogoutResponse godoc
//	@tags			基本接口
//	@Summary		用户登出
//	@Description	用户登出时，调用此接口
//	@Produce		json
//	@Success		200	{object}	res.Base
//	@Failure		500	{object}	res.Base	"错误返回内容"
//	@Router			/api/logout [get]
//	@Security		ApiKeyAuth
func (auth *Auth) LogoutResponse(ctx *gin.Context, _ int) {
	common.Log.Info("退出成功")
	res.Success(ctx)
}

// RefreshResponse godoc
//	@tags			基本接口
//	@Summary		刷新token
//	@Description	当token即将获取或者过期时刷新token
//	@Produce		json
//	@Success		200	{object}	res.Full{data=JwtToken}	"code:200 成功"
//	@Failure		500	{object}	res.Base				"错误返回内容"
//	@Router			/api/refreshToken [get]
//	@Security		ApiKeyAuth
func (auth *Auth) RefreshResponse(ctx *gin.Context, _ int, token string, expires time.Time) {
	jt := new(JwtToken)
	jt.Token = strings.Join([]string{common.TokenHeadName, token}, " ")
	jt.Expires = expires.Format(utils.FormatDateTimeMs)
	res.SuccessData(ctx, jt)
}

// UserAuth
// 用户鉴权
func (auth *Auth) UserAuth(info *model.LoginUser) (*model.User, error) {
	req := new(model.ReqQueryUser)
	req.Account = info.Account
	user, err := service.NewUserService().FindUser(info)
	if err != nil {
		return nil, err
	}

	if user.Id == 0 {
		return nil, fmt.Errorf("未查找到用户[%s]信息", info.Account)
	}
	return user, nil
}

// GetJwtRealm
// 获取 jwt标识
func (auth *Auth) GetJwtRealm() string {
	return "042f7a4b82bb4c48a9cb3082a47818532765c0cc"
}

// GetJwtKey
// jwt 秘钥
func (auth *Auth) GetJwtKey() string {
	return "6046ce088ad7283fc513733974f97cbae2f71282"
}

// GetJwtTimeOut
// 超时
func (auth *Auth) GetJwtTimeOut() int64 {
	return 24
}

// GetJwtMaxRefresh
// 刷新时间
func (auth *Auth) GetJwtMaxRefresh() int64 {
	return 5
}

// GetUser
// 获取用户信息
func (auth *Auth) GetUser(ctx *gin.Context) (*model.User, error) {
	user, ok := ctx.Get("user")
	if !ok {
		return nil, errors.New("未获取到用户信息")
	}

	u := user.(*model.User)
	return u, nil
}
