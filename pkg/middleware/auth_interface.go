package middleware

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// type User struct {
// 	Id         int64     `json:"id"`          // 用户ID
// 	Account    string    `json:"account"`     // 登录名
// 	Name       string    `json:"name"`        // 用户名
// 	Password   string    `json:"password"`    // 用户密码
// 	Mark       string    `json:"mark"`        // 备注
// 	Enable     int64     `json:"enable"`      // 启用
// 	LoginTime  time.Time `json:"login_time"`  // 登录时间
// 	CreateTime time.Time `json:"create_time"` // 创建时间
// }

// type LoginUser struct {
// 	Account  string `json:"account"`  // 登录名
// 	Password string `json:"password"` // 用户密码
// }

type AuthInterface interface {
	// 用户鉴权
	// UserAuth(*LoginUser) (*User, error)
	// 获取jwt标识
	GetJwtRealm() string
	// jwt 秘钥
	GetJwtKey() string
	// 超时
	GetJwtTimeOut() int64
	// 刷新时间
	GetJwtMaxRefresh() int64
	// 获取用户信息
	// GetUser(*gin.Context) (*User, error)
	// 有效载荷处理
	PayloadFunc(data interface{}) jwt.MapClaims
	// 解析Claims
	IdentityHandler(c *gin.Context) interface{}
	// 用户登录
	Login(c *gin.Context) (interface{}, error)
	// 用户登录校验成功处理
	Authorizator(data interface{}, c *gin.Context) bool
	// 用户登录校验失败处理
	Unauthorized(ctx *gin.Context, code int, message string)
	// 登录成功后的响应
	LoginResponse(ctx *gin.Context, _ int, token string, expires time.Time)
	// 用户登出
	LogoutResponse(ctx *gin.Context, _ int)
	// 刷新token
	RefreshResponse(ctx *gin.Context, _ int, token string, expires time.Time)
}
