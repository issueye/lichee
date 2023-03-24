package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/**
 * @author:  Yeh
 * @date：   2022年2月23日09:28:45
 * @note:   CORS 跨域中间件
 */

var (
	methods      = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	allowMethods = strings.Join(methods, ", ")
)

// CORSMiddleware 跨域中间件
func CORSMiddleware(headerMap []string) gin.HandlerFunc {
	headerMap = append(headerMap, "Authorization")
	headerMap = append(headerMap, "Content-Length")
	headerMap = append(headerMap, "X-CSRF-Token")
	headerMap = append(headerMap, "Token")
	headerMap = append(headerMap, "session")
	headerMap = append(headerMap, "X_Requested_With")
	headerMap = append(headerMap, "Accept")
	headerMap = append(headerMap, "Origin")
	headerMap = append(headerMap, "Host")
	headerMap = append(headerMap, "Connection")
	headerMap = append(headerMap, "Accept-Encoding")
	headerMap = append(headerMap, "Accept-Language")
	headerMap = append(headerMap, "DNT")
	headerMap = append(headerMap, "X-CustomHeader")
	headerMap = append(headerMap, "Keep-Alive")
	headerMap = append(headerMap, "User-Agent")
	headerMap = append(headerMap, "X-Requested-With")
	headerMap = append(headerMap, "If-Modified-Since")
	headerMap = append(headerMap, "Cache-Control")
	headerMap = append(headerMap, "Content-Type")
	headerMap = append(headerMap, "Pragma")
	headerMap = append(headerMap, "access_token")
	headerMap = append(headerMap, "Queue-Type")
	headerMap = append(headerMap, "screenType")
	crosHeaderStr := strings.Join(headerMap, ",")

	return func(c *gin.Context) {
		if c.Request.Header.Get("Origin") != "" {
			// 允许跨域设置
			c.Header("Access-Control-Allow-Origin", "*")           // 跨域请求可以从中执行的来源列表。如果设置为特殊的 "*" 值，则将允许所有来源的跨域请求
			c.Header("Access-Control-Allow-Methods", allowMethods) // 允许客户端与跨域请求一起使用的方法的列表
			// 允许客户端与跨域请求一起使用的非简单标头的列表
			// Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader,
			// Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma, access_token, Queue-Type, screenType

			c.Header("Access-Control-Allow-Headers", crosHeaderStr)
			// 指示哪些标头可以安全地暴露给 CORS API 规范的 API
			c.Header("Access-Control-Expose-Headers", "Content-Disposition, Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")

			c.Header("Access-Control-Max-Age", "172800") // 指示预检请求的结果可以缓存多长时间 (以秒为单位)
			// 指示请求是否可以包括诸如 Cookie，HTTP 身份验证或客户端 SSL 证书之类的用户凭据，为 true 时 Access-Control-Allow-Origin 必须指定一个确定的域名
			c.Header("Access-Control-Allow-Credentials", "false") // 设置返回格式是json
		}
		// 安全相关自定义响应头
		c.Header("X-Content-Type-Options", "nosniff") // 设置 Content-Type 响应头不会被自动更改
		c.Header("X-XSS-Protection", "1; mode=block") // 防止 XSS 攻击

		// OPTIONS 请求直接返回 HTTP 204 状态码，并终止本次请求，不进行后续处理
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next() //  处理请求
	}
}
