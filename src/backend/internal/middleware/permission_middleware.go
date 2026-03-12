package middleware

import (
	"net/http"

	"biolitmanager/internal/security"
	"biolitmanager/pkg/logger"
	"biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PermissionMiddleware 权限中间件
// requiredPermission: 需要的权限
func PermissionMiddleware(requiredPermission security.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户信息
		role, exists := c.Get("role")
		if !exists {
			logger.GetLogger().Warn("Role not found in context",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
			response.Error(c, http.StatusUnauthorized, "用户未认证")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			logger.GetLogger().Warn("Invalid role type in context",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
			response.Error(c, http.StatusUnauthorized, "用户认证信息异常")
			c.Abort()
			return
		}

		// 超级管理员拥有所有权限
		if roleStr == string(security.RoleSuperAdmin) {
			logger.GetLogger().Debug("Super admin granted all permissions",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("role", roleStr),
			)
			c.Next()
			return
		}

		// 检查用户权限列表是否包含所需权限
		permissions, exists := c.Get("permissions")
		if !exists {
			logger.GetLogger().Warn("Permissions not found in context",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("role", roleStr),
			)
			response.Error(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		permList, ok := permissions.([]string)
		if !ok {
			logger.GetLogger().Warn("Invalid permissions type in context",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("role", roleStr),
			)
			response.Error(c, http.StatusForbidden, "权限信息异常")
			c.Abort()
			return
		}

		// 检查是否有所需权限
		requiredPermStr := string(requiredPermission)
		hasPermission := false
		for _, perm := range permList {
			if perm == requiredPermStr {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			logger.GetLogger().Warn("Permission denied",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("role", roleStr),
				zap.String("required_permission", requiredPermStr),
				zap.Strings("user_permissions", permList),
			)
			response.Error(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		logger.GetLogger().Debug("Permission granted",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("role", roleStr),
			zap.String("permission", requiredPermStr),
		)

		c.Next()
	}
}
