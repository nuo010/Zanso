package handler

import (
	"zanso/util"

	"github.com/gin-gonic/gin"
)

// 自定义 Gin 日志中间件，使用你的 util.Log()
func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求（可选）
		util.Log().Info("%s %s %s", c.ClientIP(), c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}
