package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/issueye/lichee/utils"
)

func Req() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("RQ_DATETIME", utils.Ltime{}.GetNowStr())
		c.Set("RQ_ID", utils.Lid{}.GetCryptId())
		c.Next()
	}
}
