package middleware

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
)

var TokenHeadName = "Bearer" // Token 认证方式

// InitAuth
// 初始化 JWT 中间件
func InitAuth(ai AuthInterface) (*jwt.GinJWTMiddleware, error) {
	Authdleware, err := jwt.New(&jwt.GinJWTMiddleware{
		DisabledAbort:   true,                                               //禁止此三方包内部Abort()
		Realm:           ai.GetJwtRealm(),                                   // jwt 标识
		Key:             []byte(ai.GetJwtKey()),                             // 服务端密钥
		Timeout:         time.Hour * time.Duration(ai.GetJwtTimeOut()),      // token 过期时间
		MaxRefresh:      time.Hour * time.Duration(ai.GetJwtMaxRefresh()),   // token 最大刷新时间(RefreshToken 过期时间 = Timeout+MaxRefresh)
		PayloadFunc:     ai.PayloadFunc,                                     // 有效载荷处理
		IdentityHandler: ai.IdentityHandler,                                 // 解析 Claims
		Authenticator:   ai.Login,                                           // 校验 token 的正确性, 处理登录逻辑
		Authorizator:    ai.Authorizator,                                    // 用户登录校验成功处理
		Unauthorized:    ai.Unauthorized,                                    // 用户登录校验失败处理
		LoginResponse:   ai.LoginResponse,                                   // 登录成功后的响应
		LogoutResponse:  ai.LogoutResponse,                                  // 登出后的响应
		RefreshResponse: ai.RefreshResponse,                                 // 刷新 token 后的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt", // 自动在这几个地方寻找请求中的 token
		TokenHeadName:   TokenHeadName,                                      // header 名称
		TimeFunc:        time.Now,
	})
	return Authdleware, err
}
