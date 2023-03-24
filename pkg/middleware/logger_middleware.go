package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/**
 * @author:  Yeh
 * @date：   2022年2月23日09:28:45
 * @note:   对业务对象进行访问日志收集
 */

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger, args ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var showAuth bool
		if len(args) > 0 {
			showAuth = args[0].(bool)
		}

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		auth := c.Request.Header["Authorization"]
		c.Next()
		cost := time.Since(start)
		if showAuth {
			logger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", cost),
				zap.String("auth", strings.Join(auth, " ")),
			)
		} else {
			logger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", cost),
			)
		}
	}
}
