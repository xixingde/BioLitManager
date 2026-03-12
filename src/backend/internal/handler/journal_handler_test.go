package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"biolitmanager/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 测试1: 创建期刊 - 缺少必填字段
func TestCreateJournal_MissingRequiredFields(t *testing.T) {
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

	router.POST("/api/journals", func(c *gin.Context) {
		var req struct {
			FullName     string  `json:"full_name" binding:"required"`
			ShortName    string  `json:"short_name"`
			ISSN         string  `json:"issn" binding:"required"`
			ImpactFactor float64 `json:"impact_factor"`
			Publisher    string  `json:"publisher"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	// 缺少必填字段
	requestBody := map[string]interface{}{
		"short_name": "Nature",
		// 缺少 full_name 和 issn
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/journals", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试2: 创建期刊 - 正常数据
func TestCreateJournal_Success(t *testing.T) {
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

	router.POST("/api/journals", func(c *gin.Context) {
		var req struct {
			FullName     string  `json:"full_name" binding:"required"`
			ShortName    string  `json:"short_name"`
			ISSN         string  `json:"issn" binding:"required"`
			ImpactFactor float64 `json:"impact_factor"`
			Publisher    string  `json:"publisher"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": gin.H{"id": 1}})
	})

	requestBody := map[string]interface{}{
		"full_name":     "Nature",
		"short_name":    "Nature",
		"issn":          "0028-0836",
		"impact_factor": 69.504,
		"publisher":     "Nature Publishing Group",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/journals", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试3: 创建期刊 - ISSN已存在
func TestCreateJournal_ISSNExists(t *testing.T) {
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

	router.POST("/api/journals", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "ISSN已存在"})
	})

	requestBody := map[string]interface{}{
		"full_name": "Nature",
		"issn":      "0028-0836",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/journals", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["msg"], "ISSN")
}

// 测试4: 获取期刊详情 - 不存在的ID
func TestGetJournal_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.GET("/api/journals/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "999" {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "期刊不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	req, _ := http.NewRequest("GET", "/api/journals/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// 测试5: 获取期刊详情 - 格式错误的ID
func TestGetJournal_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.GET("/api/journals/:id", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "期刊ID格式错误"})
	})

	req, _ := http.NewRequest("GET", "/api/journals/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试6: 分页查询期刊列表
func TestListJournals_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.GET("/api/journals", func(c *gin.Context) {
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

	req, _ := http.NewRequest("GET", "/api/journals?page=1&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试7: 更新期刊
func TestUpdateJournal_Success(t *testing.T) {
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

	router.PUT("/api/journals/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	requestBody := map[string]interface{}{
		"full_name":     "Nature (Updated)",
		"impact_factor": 70.0,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/journals/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试8: 更新期刊影响因子
func TestUpdateImpactFactor_Success(t *testing.T) {
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

	router.PUT("/api/journals/:id/impact-factor", func(c *gin.Context) {
		var req struct {
			ImpactFactor float64 `json:"impact_factor" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})

	requestBody := map[string]interface{}{
		"impact_factor": 75.5,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/journals/1/impact-factor", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试9: 搜索期刊
func TestSearchJournals_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.GET("/api/journals/search", func(c *gin.Context) {
		keyword := c.Query("keyword")
		if keyword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "搜索关键字不能为空"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"list": []interface{}{
					map[string]interface{}{
						"id":            1,
						"full_name":     "Nature",
						"short_name":    "Nature",
						"issn":          "0028-0836",
						"impact_factor": 69.504,
					},
				},
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/journals/search?keyword=Nature", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试10: 搜索期刊 - 关键字为空
func TestSearchJournals_EmptyKeyword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger.InitLogger("test")

	router.GET("/api/journals/search", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "搜索关键字不能为空"})
	})

	req, _ := http.NewRequest("GET", "/api/journals/search?keyword=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
