package middleware

import (
	"net/http"

	"biolitmanager/pkg/logger"
	"biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware 异常恢复中间件
// 使用 defer + recover 捕获 panic，记录错误日志，返回统一的错误响应
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				logger.GetLogger().Error("Panic recovered",
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
					zap.Stack("stack"),
				)

				// 获取用户信息（如果已认证）
				var userID uint
				var username string
				if uid, exists := c.Get("user_id"); exists {
					if u, ok := uid.(uint); ok {
						userID = u
					}
				}
				if uname, exists := c.Get("username"); exists {
					if u, ok := uname.(string); ok {
						username = u
					}
				}

				// 记录更多上下文信息
				logger.GetLogger().Error("Panic context info",
					zap.Uint("user_id", userID),
					zap.String("username", username),
					zap.String("query", c.Request.URL.RawQuery),
				)

				// 返回统一的错误响应
				response.Error(c, http.StatusInternalServerError, "系统异常，请稍后重试")

				// 终止请求处理
				c.Abort()
			}
		}()

		c.Next()
	}
}
