package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/**
 * @author:  Yeh
 * @date：   2022年2月23日09:28:45
 * @note:   恐慌处理，在服务出现报错时能够继续运行服务
 */

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				// 打印调用栈
				if stack {
					logger.Sugar().Errorf("出现异常恐慌，堆栈信息:%s", debug.Stack())

					// logger.Error("[Recovery from panic]",
					// 	zap.Any("error", err),
					// 	zap.String("request", string(httpRequest)),
					// 	zap.String("stack", string(debug.Stack())),
					// )
				} else {
					logger.Sugar().Errorf("出现异常恐慌：%s", err)

					// logger.Error("[Recovery from panic]",
					// 	zap.Any("error", err),
					// 	// zap.String("request", string(httpRequest)),
					// )
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
