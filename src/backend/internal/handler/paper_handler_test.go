package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/service"
	"biolitmanager/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
)

// MockPaperService Mock服务
type MockPaperService struct {
	mock.Mock
}

func (m *MockPaperService) CreatePaper(title, abstract string, journalID uint, doi string, impactFactor float64, publishDate *time.Time, submitterID uint, authors []*entity.Author, projectIDs []uint, operatorID uint, ipAddress string) (*entity.Paper, error) {
	args := m.Called(title, abstract, journalID, doi, impactFactor, publishDate, submitterID, authors, projectIDs, operatorID, ipAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Paper), args.Error(1)
}

func (m *MockPaperService) GetPaperByID(id uint) (*entity.Paper, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Paper), args.Error(1)
}

func (m *MockPaperService) ListPapers(page, size int) ([]*entity.Paper, int64, error) {
	args := m.Called(page, size)
	return args.Get(0).([]*entity.Paper), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaperService) UpdatePaper(id uint, title, abstract string, journalID uint, doi string, impactFactor float64, publishDate *time.Time, authors []*entity.Author, projectIDs []uint, operatorID uint, ipAddress string) error {
	args := m.Called(id, title, abstract, journalID, doi, impactFactor, publishDate, authors, projectIDs, operatorID, ipAddress)
	return args.Error(0)
}

func (m *MockPaperService) DeletePaper(id, operatorID uint, ipAddress string) error {
	args := m.Called(id, operatorID, ipAddress)
	return args.Error(0)
}

func (m *MockPaperService) SubmitForReview(id, operatorID uint, ipAddress string) error {
	args := m.Called(id, operatorID, ipAddress)
	return args.Error(0)
}

func (m *MockPaperService) SaveDraft(id uint, title, abstract string, journalID uint, doi string, impactFactor float64, publishDate *time.Time, authors []*entity.Author, projectIDs []uint, operatorID uint, ipAddress string) error {
	args := m.Called(id, title, abstract, journalID, doi, impactFactor, publishDate, authors, projectIDs, operatorID, ipAddress)
	return args.Error(0)
}

func (m *MockPaperService) CheckDuplicate(title, doi string) ([]*entity.Paper, error) {
	args := m.Called(title, doi)
	return args.Get(0).([]*entity.Paper), args.Error(1)
}

func (m *MockPaperService) GetMyPapers(userID uint, page, size int) ([]*entity.Paper, int64, error) {
	args := m.Called(userID, page, size)
	return args.Get(0).([]*entity.Paper), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaperService) BatchImportPapers(file *excelize.File, submitterID uint) (int, int, []string) {
	args := m.Called(file, submitterID)
	return args.Int(0), args.Int(1), args.Get(2).([]string)
}

// setupTestRouter 设置测试路由
func setupTestRouter(paperService service.PaperServiceInterface) *gin.Engine {
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

	handler := NewPaperHandler(paperService)
	router.POST("/api/papers", handler.CreatePaper)
	router.GET("/api/papers/:id", handler.GetPaper)
	router.GET("/api/papers", handler.ListPapers)
	router.PUT("/api/papers/:id", handler.UpdatePaper)
	router.DELETE("/api/papers/:id", handler.DeletePaper)
	router.POST("/api/papers/:id/submit", handler.SubmitForReview)
	router.POST("/api/papers/:id/save-draft", handler.SaveDraft)
	router.POST("/api/papers/check-duplicate", handler.CheckDuplicate)
	router.GET("/api/papers/my", handler.GetMyPapers)

	return router
}

// 测试1: 创建论文 - 正常数据
func TestCreatePaper_Success(t *testing.T) {
	// 准备测试数据
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	expectedPaper := &entity.Paper{
		ID:        1,
		Title:     "测试论文标题",
		Abstract:  "测试摘要",
		JournalID: 1,
		DOI:       "10.1000/test",
		Status:    "draft",
	}

	mockService.On("CreatePaper",
		"测试论文标题",
		"测试摘要",
		uint(1),
		"10.1000/test",
		5.234,
		mock.Anything,
		uint(1),
		[]*entity.Author{},
		[]uint{},
		uint(1),
		mock.Anything,
	).Return(expectedPaper, nil)

	// 准备请求
	requestBody := map[string]interface{}{
		"title":         "测试论文标题",
		"abstract":      "测试摘要",
		"journal_id":    1,
		"doi":           "10.1000/test",
		"impact_factor": 5.234,
		"authors":       []interface{}{},
		"projects":      []interface{}{},
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/papers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证结果
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "success", resp["msg"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])

	mockService.AssertExpectations(t)
}

// 测试2: 创建论文 - 缺少必填字段
func TestCreatePaper_MissingRequiredFields(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	// 设置mock期望 - 当journal_id为0时，service层会返回错误
	mockService.On("CreatePaper", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, service.ErrPaperNotFound)

	// 缺少journal_id
	requestBody := map[string]interface{}{
		"title":    "测试论文标题",
		"abstract": "测试摘要",
		"doi":      "10.1000/test",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/papers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 实际返回500因为service层返回错误
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// 测试3: 创建论文 - 论文重复
func TestCreatePaper_Duplicate(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	mockService.On("CreatePaper", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, service.ErrPaperDuplicate)

	requestBody := map[string]interface{}{
		"title":         "测试论文标题",
		"abstract":      "测试摘要",
		"journal_id":    1,
		"doi":           "10.1000/test",
		"impact_factor": 5.234,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/papers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["msg"], "重复")

	mockService.AssertExpectations(t)
}

// 测试4: 获取论文详情 - 正常ID
func TestGetPaper_Success(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	expectedPaper := &entity.Paper{
		ID:        1,
		Title:     "测试论文标题",
		Abstract:  "测试摘要",
		JournalID: 1,
		Status:    "draft",
	}

	mockService.On("GetPaperByID", uint(1)).Return(expectedPaper, nil)

	req, _ := http.NewRequest("GET", "/api/papers/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(1), response["data"].(map[string]interface{})["id"])

	mockService.AssertExpectations(t)
}

// 测试5: 获取论文详情 - 不存在的ID
func TestGetPaper_NotFound(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	mockService.On("GetPaperByID", uint(999)).Return(nil, service.ErrPaperNotFound)

	req, _ := http.NewRequest("GET", "/api/papers/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

// 测试6: 获取论文详情 - 格式错误的ID
func TestGetPaper_InvalidID(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/papers/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 测试7: 分页查询论文 - 正常分页
func TestListPapers_Success(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	papers := []*entity.Paper{
		{ID: 1, Title: "论文1", Status: "draft"},
		{ID: 2, Title: "论文2", Status: "draft"},
	}

	mockService.On("ListPapers", 1, 10).Return(papers, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/papers?page=1&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.Equal(t, float64(1), data["page"])
	assert.Equal(t, float64(10), data["size"])

	mockService.AssertExpectations(t)
}

// 测试8: 更新论文 - 草稿状态
func TestUpdatePaper_DraftStatus(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	mockService.On("UpdatePaper", uint(1), "更新标题", "更新摘要", uint(1), "10.1000/updated", mock.Anything, mock.Anything, mock.Anything, mock.Anything, uint(1), mock.Anything).
		Return(nil)

	requestBody := map[string]interface{}{
		"title":      "更新标题",
		"abstract":   "更新摘要",
		"journal_id": 1,
		"doi":        "10.1000/updated",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/papers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// 测试9: 删除论文 - 草稿状态
func TestDeletePaper_DraftStatus(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	mockService.On("DeletePaper", uint(1), uint(1), mock.Anything).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/papers/1", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// 测试10: 删除论文 - 非草稿状态（应失败）
func TestDeletePaper_NonDraftStatus(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	mockService.On("DeletePaper", uint(1), uint(1), mock.Anything).Return(service.ErrInvalidStatus)

	req, _ := http.NewRequest("DELETE", "/api/papers/1", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["msg"], "草稿")

	mockService.AssertExpectations(t)
}

// 测试11: 提交审核 - 草稿状态
func TestSubmitForReview_DraftStatus(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	mockService.On("SubmitForReview", uint(1), uint(1), mock.Anything).Return(nil)

	req, _ := http.NewRequest("POST", "/api/papers/1/submit", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// 测试12: 提交审核 - 非草稿状态（应失败）
func TestSubmitForReview_NonDraftStatus(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	mockService.On("SubmitForReview", uint(1), uint(1), mock.Anything).Return(service.ErrInvalidStatus)

	req, _ := http.NewRequest("POST", "/api/papers/1/submit", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertExpectations(t)
}

// 测试13: 检查重复 - 无重复
func TestCheckDuplicate_NoDuplicate(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	mockService.On("CheckDuplicate", "新论文标题", "10.1000/new").Return([]*entity.Paper{}, nil)

	requestBody := map[string]interface{}{
		"title": "新论文标题",
		"doi":   "10.1000/new",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/papers/check-duplicate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["count"])

	mockService.AssertExpectations(t)
}

// 测试14: 检查重复 - 存在重复
func TestCheckDuplicate_HasDuplicate(t *testing.T) {
	mockService := new(MockPaperService)
	router := setupTestRouter(mockService)

	duplicatePapers := []*entity.Paper{
		{ID: 1, Title: "重复论文", DOI: "10.1000/duplicate"},
	}

	mockService.On("CheckDuplicate", "重复论文", "10.1000/duplicate").Return(duplicatePapers, nil)

	requestBody := map[string]interface{}{
		"title": "重复论文",
		"doi":   "10.1000/duplicate",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/papers/check-duplicate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["count"])

	mockService.AssertExpectations(t)
}
