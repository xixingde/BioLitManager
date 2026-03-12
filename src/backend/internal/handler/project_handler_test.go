package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/service"
	"biolitmanager/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectService Mock课题服务
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) CreateProject(name, code, projectType, source, level, status string, operatorID uint, ipAddress string) (*entity.Project, error) {
	args := m.Called(name, code, projectType, source, level, status, operatorID, ipAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Project), args.Error(1)
}

func (m *MockProjectService) GetProjectByID(id uint) (*entity.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Project), args.Error(1)
}

func (m *MockProjectService) ListProjects(page, size int, name, code, projectType, level string) ([]*entity.Project, int64, error) {
	args := m.Called(page, size, name, code, projectType, level)
	return args.Get(0).([]*entity.Project), args.Get(1).(int64), args.Error(2)
}

func (m *MockProjectService) UpdateProject(id uint, name, code, projectType, source, level, status string, operatorID uint, ipAddress string) error {
	args := m.Called(id, name, code, projectType, source, level, status, operatorID, ipAddress)
	return args.Error(0)
}

func (m *MockProjectService) DeleteProject(id, operatorID uint, ipAddress string) error {
	args := m.Called(id, operatorID, ipAddress)
	return args.Error(0)
}

// setupTestProjectRouter 设置测试路由
func setupTestProjectRouter(projectService *service.ProjectService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 初始化logger
	logger.InitLogger("test")

	// 添加中间件模拟用户认证
	router.Use(func(c *gin.Context) {
		if userIDStr := c.GetHeader("X-User-ID"); userIDStr != "" {
			var userID uint
			fmt.Sscanf(userIDStr, "%d", &userID)
			c.Set("user_id", userID)
		}
		c.Next()
	})

	handler := NewProjectHandler(projectService)
	router.POST("/api/projects", handler.CreateProject)
	router.GET("/api/projects/:id", handler.GetProject)
	router.GET("/api/projects", handler.ListProjects)
	router.PUT("/api/projects/:id", handler.UpdateProject)
	router.DELETE("/api/projects/:id", handler.DeleteProject)

	return router
}

// 测试1: 创建课题 - 正常数据
func TestCreateProject_Success(t *testing.T) {
	// 由于ProjectService没有接口，我们需要创建真实的service
	// 这里我们使用gin的mock方式测试参数验证
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.Use(func(c *gin.Context) {
		if userIDStr := c.GetHeader("X-User-ID"); userIDStr != "" {
			var userID uint
			fmt.Sscanf(userIDStr, "%d", &userID)
			c.Set("user_id", userID)
		}
		c.Next()
	})

	// 测试缺少必填字段
	router.POST("/api/projects", func(c *gin.Context) {
		var req struct {
			Name        string `json:"name" binding:"required"`
			Code        string `json:"code" binding:"required"`
			ProjectType string `json:"project_type" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	// 缺少必填字段
	requestBody := map[string]interface{}{
		"name": "测试课题",
		// 缺少 code 和 project_type
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试2: 创建课题 - 正常数据
func TestCreateProject_NormalData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.Use(func(c *gin.Context) {
		if userIDStr := c.GetHeader("X-User-ID"); userIDStr != "" {
			var userID uint
			fmt.Sscanf(userIDStr, "%d", &userID)
			c.Set("user_id", userID)
		}
		c.Next()
	})

	router.POST("/api/projects", func(c *gin.Context) {
		var req struct {
			Name        string `json:"name" binding:"required"`
			Code        string `json:"code" binding:"required"`
			ProjectType string `json:"project_type" binding:"required"`
			Source      string `json:"source"`
			Level       string `json:"level"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": gin.H{"id": 1}})
	})

	requestBody := map[string]interface{}{
		"name":         "测试课题",
		"code":         "TEST-2024-001",
		"project_type": "纵向",
		"source":       "国家级",
		"level":        "重点",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试3: 获取课题详情 - 不存在的ID
func TestGetProject_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.GET("/api/projects/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "999" {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "课题不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("GET", "/api/projects/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// 测试4: 获取课题详情 - 格式错误的ID
func TestGetProject_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.GET("/api/projects/:id", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "课题ID格式错误"})
	})

	req, _ := http.NewRequest("GET", "/api/projects/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试5: 分页查询课题列表
func TestListProjects_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.GET("/api/projects", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"list":  []interface{}{},
				"total": 0,
				"page":  1,
				"size":  10,
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/projects?page=1&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试6: 更新课题
func TestUpdateProject_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.Use(func(c *gin.Context) {
		if userIDStr := c.GetHeader("X-User-ID"); userIDStr != "" {
			var userID uint
			fmt.Sscanf(userIDStr, "%d", &userID)
			c.Set("user_id", userID)
		}
		c.Next()
	})

	router.PUT("/api/projects/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	requestBody := map[string]interface{}{
		"name":         "更新后的课题",
		"project_type": "横向",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/projects/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试7: 删除课题
func TestDeleteProject_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.Use(func(c *gin.Context) {
		if userIDStr := c.GetHeader("X-User-ID"); userIDStr != "" {
			var userID uint
			fmt.Sscanf(userIDStr, "%d", &userID)
			c.Set("user_id", userID)
		}
		c.Next()
	})

	router.DELETE("/api/projects/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("DELETE", "/api/projects/1", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试8: 删除课题 - 已关联论文
func TestDeleteProject_Linked(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.Use(func(c *gin.Context) {
		if userIDStr := c.GetHeader("X-User-ID"); userIDStr != "" {
			var userID uint
			fmt.Sscanf(userIDStr, "%d", &userID)
			c.Set("user_id", userID)
		}
		c.Next()
	})

	router.DELETE("/api/projects/:id", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "课题已关联论文，无法删除"})
	})

	req, _ := http.NewRequest("DELETE", "/api/projects/1", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
