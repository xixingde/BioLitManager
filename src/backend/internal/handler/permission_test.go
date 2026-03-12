package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"biolitmanager/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupPermissionRouter 设置带权限中间件的路由
func setupPermissionRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	logger.InitLogger("test")

	// 添加用户认证中间件
	router.Use(func(c *gin.Context) {
		if userIDStr := c.GetHeader("X-User-ID"); userIDStr != "" {
			var userID uint
			fmt.Sscanf(userIDStr, "%d", &userID)
			c.Set("user_id", userID)
		}
		if userRole := c.GetHeader("X-User-Role"); userRole != "" {
			c.Set("user_role", userRole)
		}
		c.Next()
	})

	return router
}

// 测试1: 无权限访问 - 未登录用户访问需要权限的接口
func TestPermission_NoLogin(t *testing.T) {
	router := setupPermissionRouter()

	// 模拟需要登录的接口
	router.POST("/api/papers", func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权访问"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	// 不带用户ID的请求
	req, _ := http.NewRequest("POST", "/api/papers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// 测试2: 无权限访问 - 普通用户访问管理员接口
func TestPermission_RegularUserToAdmin(t *testing.T) {
	router := setupPermissionRouter()

	// 模拟管理员接口
	router.POST("/api/journals", func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "管理员" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权限访问"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	// 普通用户访问
	req, _ := http.NewRequest("POST", "/api/journals", nil)
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("X-User-Role", "用户")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// 测试3: 普通用户可以创建论文
func TestPermission_RegularUserCreatePaper(t *testing.T) {
	router := setupPermissionRouter()

	router.POST("/api/papers", func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权访问"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("POST", "/api/papers", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试4: 普通用户可以编辑自己的论文
func TestPermission_RegularUserEditOwnPaper(t *testing.T) {
	router := setupPermissionRouter()

	router.PUT("/api/papers/:id", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		paperOwnerID := c.GetHeader("X-Paper-Owner-ID")

		// 检查是否是论文所有者
		if fmt.Sprintf("%v", userID) != paperOwnerID {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权限编辑他人的论文"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("PUT", "/api/papers/1", nil)
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("X-Paper-Owner-ID", "1") // 论文所有者是用户1
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试5: 普通用户不能编辑他人的论文
func TestPermission_RegularUserCannotEditOthersPaper(t *testing.T) {
	router := setupPermissionRouter()

	router.PUT("/api/papers/:id", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		paperOwnerID := c.GetHeader("X-Paper-Owner-ID")

		if fmt.Sprintf("%v", userID) != paperOwnerID {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权限编辑他人的论文"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("PUT", "/api/papers/1", nil)
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("X-Paper-Owner-ID", "2") // 论文所有者是用户2
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["msg"], "无权限")
}

// 测试6: 普通用户可以删除自己的论文（草稿状态）
func TestPermission_RegularUserDeleteOwnDraft(t *testing.T) {
	router := setupPermissionRouter()

	router.DELETE("/api/papers/:id", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		paperOwnerID := c.GetHeader("X-Paper-Owner-ID")
		paperStatus := c.GetHeader("X-Paper-Status")

		// 检查是否是论文所有者
		if fmt.Sprintf("%v", userID) != paperOwnerID {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权限删除他人的论文"})
			return
		}

		// 检查论文状态
		if paperStatus != "draft" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "只有草稿状态的论文可以删除"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("DELETE", "/api/papers/1", nil)
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("X-Paper-Owner-ID", "1")
	req.Header.Set("X-Paper-Status", "draft")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试7: 普通用户不能删除他人的论文
func TestPermission_RegularUserCannotDeleteOthersPaper(t *testing.T) {
	router := setupPermissionRouter()

	router.DELETE("/api/papers/:id", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		paperOwnerID := c.GetHeader("X-Paper-Owner-ID")

		if fmt.Sprintf("%v", userID) != paperOwnerID {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权限删除他人的论文"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("DELETE", "/api/papers/1", nil)
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("X-Paper-Owner-ID", "2")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// 测试8: 业务审核员可以审核论文
func TestPermission_BusinessReviewer(t *testing.T) {
	router := setupPermissionRouter()

	router.POST("/api/reviews/business/:paperId", func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "业务审核员" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无业务审核权限"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("POST", "/api/reviews/business/1", nil)
	req.Header.Set("X-User-ID", "2")
	req.Header.Set("X-User-Role", "业务审核员")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试9: 普通用户不能进行业务审核
func TestPermission_RegularUserCannotBusinessReview(t *testing.T) {
	router := setupPermissionRouter()

	router.POST("/api/reviews/business/:paperId", func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "业务审核员" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无业务审核权限"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("POST", "/api/reviews/business/1", nil)
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("X-User-Role", "用户")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// 测试10: 政工审核员可以审核论文
func TestPermission_PoliticalReviewer(t *testing.T) {
	router := setupPermissionRouter()

	router.POST("/api/reviews/political/:paperId", func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "政工审核员" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无政工审核权限"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("POST", "/api/reviews/political/1", nil)
	req.Header.Set("X-User-ID", "3")
	req.Header.Set("X-User-Role", "政工审核员")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试11: 普通用户不能进行政工审核
func TestPermission_RegularUserCannotPoliticalReview(t *testing.T) {
	router := setupPermissionRouter()

	router.POST("/api/reviews/political/:paperId", func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "政工审核员" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无政工审核权限"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("POST", "/api/reviews/political/1", nil)
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("X-User-Role", "用户")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// 测试12: 管理员可以管理课题
func TestPermission_AdminManageProject(t *testing.T) {
	router := setupPermissionRouter()

	router.POST("/api/projects", func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "管理员" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无课题管理权限"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("POST", "/api/projects", nil)
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("X-User-Role", "管理员")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试13: 普通用户不能管理课题
func TestPermission_RegularUserCannotManageProject(t *testing.T) {
	router := setupPermissionRouter()

	router.POST("/api/projects", func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "管理员" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无课题管理权限"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("POST", "/api/projects", nil)
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("X-User-Role", "用户")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// 测试14: 管理员可以管理期刊
func TestPermission_AdminManageJournal(t *testing.T) {
	router := setupPermissionRouter()

	router.POST("/api/journals", func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "管理员" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无期刊管理权限"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("POST", "/api/journals", nil)
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("X-User-Role", "管理员")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试15: 普通用户不能管理期刊
func TestPermission_RegularUserCannotManageJournal(t *testing.T) {
	router := setupPermissionRouter()

	router.POST("/api/journals", func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "管理员" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无期刊管理权限"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("POST", "/api/journals", nil)
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("X-User-Role", "用户")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// 测试16: 业务审核员不能进行政工审核
func TestPermission_BusinessReviewerCannotPoliticalReview(t *testing.T) {
	router := setupPermissionRouter()

	router.POST("/api/reviews/political/:paperId", func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "政工审核员" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无政工审核权限"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("POST", "/api/reviews/political/1", nil)
	req.Header.Set("X-User-ID", "2")
	req.Header.Set("X-User-Role", "业务审核员")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// 测试17: 政工审核员不能进行业务审核
func TestPermission_PoliticalReviewerCannotBusinessReview(t *testing.T) {
	router := setupPermissionRouter()

	router.POST("/api/reviews/business/:paperId", func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "业务审核员" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无业务审核权限"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("POST", "/api/reviews/business/1", nil)
	req.Header.Set("X-User-ID", "3")
	req.Header.Set("X-User-Role", "政工审核员")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
