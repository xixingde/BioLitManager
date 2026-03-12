package middleware

import (
	"net/http"
	"strings"

	"biolitmanager/internal/cache"
	"biolitmanager/internal/config"
	"biolitmanager/internal/security"
	"biolitmanager/pkg/logger"
	"biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.GetLogger().Warn("Authorization header is missing",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
			response.Error(c, http.StatusUnauthorized, "未授权访问")
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.GetLogger().Warn("Invalid Authorization format",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
			response.Error(c, http.StatusUnauthorized, "认证格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		// 解析 Token 验证有效性
		claims, err := security.ParseToken(token)
		if err != nil {
			logger.GetLogger().Warn("Invalid token",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Error(err),
			)
			response.Error(c, http.StatusUnauthorized, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 检查缓存中的会话是否存在
		cfg := config.GetConfig()
		if cfg != nil {
			cacheInstance := cache.GetInstance()
			logger.GetLogger().Debug("Auth middleware checking cache",
				zap.Any("cache_instance", cacheInstance),
			)
			if cacheInstance != nil {
				sessionKey := cache.SessionKey(token)
				session, exists := cacheInstance.Get(sessionKey)
				if !exists {
					logger.GetLogger().Warn("Session not found in cache",
						zap.String("path", c.Request.URL.Path),
						zap.String("method", c.Request.Method),
						zap.Uint("user_id", claims.UserID),
					)
					response.Error(c, http.StatusUnauthorized, "会话已过期")
					c.Abort()
					return
				}

				// 验证会话数据是否匹配
				if sessionData, ok := session.(*security.Claims); ok {
					if sessionData.UserID != claims.UserID {
						logger.GetLogger().Warn("Session data mismatch",
							zap.String("path", c.Request.URL.Path),
							zap.String("method", c.Request.Method),
							zap.Uint("token_user_id", claims.UserID),
							zap.Uint("session_user_id", sessionData.UserID),
						)
						response.Error(c, http.StatusUnauthorized, "会话数据异常")
						c.Abort()
						return
					}
				}
			}
		}

		// 将用户信息设置到 gin.Context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)
		c.Set("claims", claims)

		logger.GetLogger().Debug("Authentication successful",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("role", claims.Role),
		)

		c.Next()
	}
}
