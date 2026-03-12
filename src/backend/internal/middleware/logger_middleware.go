package middleware

import (
	"time"

	"biolitmanager/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware 日志中间件
// 记录请求方法、路径、状态码、响应时间、客户端 IP
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 记录结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 获取请求信息
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

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

		// 记录日志
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.Uint("user_id", userID),
			zap.String("username", username),
		}

		// 根据状态码选择日志级别
		log := logger.GetLogger()
		if statusCode >= 500 {
			// 服务器错误
			log.Error("Server error", fields...)
		} else if statusCode >= 400 {
			// 客户端错误
			log.Warn("Client error", fields...)
		} else {
			// 成功请求
			log.Info("Request completed", fields...)
		}
	}
}
