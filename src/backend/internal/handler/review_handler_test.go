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

// MockReviewService Mock审核服务
type MockReviewService struct {
	mock.Mock
}

func (m *MockReviewService) BusinessReview(paperID uint, result, comment string, reviewerID uint, ipAddress string) error {
	args := m.Called(paperID, result, comment, reviewerID, ipAddress)
	return args.Error(0)
}

func (m *MockReviewService) PoliticalReview(paperID uint, result, comment string, reviewerID uint, ipAddress string) error {
	args := m.Called(paperID, result, comment, reviewerID, ipAddress)
	return args.Error(0)
}

func (m *MockReviewService) GetReviewLogsByPaperID(paperID uint) ([]*entity.ReviewLog, error) {
	args := m.Called(paperID)
	return args.Get(0).([]*entity.ReviewLog), args.Error(1)
}

func (m *MockReviewService) GetPendingPapersForBusinessReview() ([]*entity.Paper, error) {
	args := m.Called()
	return args.Get(0).([]*entity.Paper), args.Error(1)
}

func (m *MockReviewService) GetPendingPapersForPoliticalReview() ([]*entity.Paper, error) {
	args := m.Called()
	return args.Get(0).([]*entity.Paper), args.Error(1)
}

func (m *MockReviewService) GetMyReviews(reviewerID uint) ([]*entity.ReviewLog, error) {
	args := m.Called(reviewerID)
	return args.Get(0).([]*entity.ReviewLog), args.Error(1)
}

// setupTestReviewRouter 设置测试路由
func setupTestReviewRouter(reviewService service.ReviewServiceInterface) *gin.Engine {
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

	handler := NewReviewHandler(reviewService)
	router.POST("/api/reviews/business/:paperId", handler.BusinessReview)
	router.POST("/api/reviews/political/:paperId", handler.PoliticalReview)
	router.GET("/api/reviews/:paperId/logs", handler.GetReviewLogs)
	router.GET("/api/reviews/pending/business", handler.GetPendingBusinessReviews)
	router.GET("/api/reviews/pending/political", handler.GetPendingPoliticalReviews)
	router.GET("/api/reviews/my", handler.GetMyReviews)

	return router
}

// 测试1: 业务审核 - 通过
func TestBusinessReview_Approve(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	mockService.On("BusinessReview", uint(1), "通过", "审核意见", uint(2), mock.Anything).Return(nil)

	requestBody := map[string]interface{}{
		"result":  "通过",
		"comment": "审核意见",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/reviews/business/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "2")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// 测试2: 业务审核 - 驳回
func TestBusinessReview_Reject(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	mockService.On("BusinessReview", uint(1), "驳回", "需要修改", uint(2), mock.Anything).Return(nil)

	requestBody := map[string]interface{}{
		"result":  "驳回",
		"comment": "需要修改",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/reviews/business/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "2")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// 测试3: 业务审核 - 无效的审核结果
func TestBusinessReview_InvalidResult(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	mockService.On("BusinessReview", uint(1), "无效结果", "", uint(2), mock.Anything).Return(service.ErrInvalidReviewResult)

	requestBody := map[string]interface{}{
		"result":  "无效结果",
		"comment": "",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/reviews/business/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "2")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["msg"], "无效")

	mockService.AssertExpectations(t)
}

// 测试4: 业务审核 - 论文不存在
func TestBusinessReview_PaperNotFound(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	mockService.On("BusinessReview", uint(999), "通过", "", uint(2), mock.Anything).Return(service.ErrPaperNotFound)

	requestBody := map[string]interface{}{
		"result":  "通过",
		"comment": "",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/reviews/business/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "2")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

// 测试5: 业务审核 - 论文状态不允许
func TestBusinessReview_InvalidStatus(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	mockService.On("BusinessReview", uint(1), "通过", "", uint(2), mock.Anything).Return(service.ErrInvalidStatus)

	requestBody := map[string]interface{}{
		"result":  "通过",
		"comment": "",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/reviews/business/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "2")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["msg"], "状态")

	mockService.AssertExpectations(t)
}

// 测试6: 政工审核 - 通过
func TestPoliticalReview_Approve(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	mockService.On("PoliticalReview", uint(1), "通过", "审核通过", uint(3), mock.Anything).Return(nil)

	requestBody := map[string]interface{}{
		"result":  "通过",
		"comment": "审核通过",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/reviews/political/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "3")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// 测试7: 政工审核 - 驳回
func TestPoliticalReview_Reject(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	mockService.On("PoliticalReview", uint(1), "驳回", "内容需要审查", uint(3), mock.Anything).Return(nil)

	requestBody := map[string]interface{}{
		"result":  "驳回",
		"comment": "内容需要审查",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/reviews/political/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "3")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// 测试8: 获取审核记录
func TestGetReviewLogs_Success(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	reviewer := &entity.User{
		ID:       2,
		Username: "reviewer1",
		Name:     "审核员1",
		Role:     "业务审核员",
	}

	reviewLogs := []*entity.ReviewLog{
		{
			ID:         1,
			PaperID:    1,
			ReviewerID: 2,
			ReviewType: "业务审核",
			Result:     "通过",
			Comment:    "审核通过",
			Reviewer:   reviewer,
		},
	}

	mockService.On("GetReviewLogsByPaperID", uint(1)).Return(reviewLogs, nil)

	req, _ := http.NewRequest("GET", "/api/reviews/1/logs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Equal(t, 1, len(data))

	mockService.AssertExpectations(t)
}

// 测试9: 获取待业务审核列表
func TestGetPendingBusinessReviews_Success(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	submitter := &entity.User{
		ID:   1,
		Name: "提交人",
	}

	pendingPapers := []*entity.Paper{
		{
			ID:        1,
			Title:     "待审核论文1",
			Status:    "待业务审核",
			Submitter: submitter,
		},
		{
			ID:        2,
			Title:     "待审核论文2",
			Status:    "待业务审核",
			Submitter: submitter,
		},
	}

	mockService.On("GetPendingPapersForBusinessReview").Return(pendingPapers, nil)

	req, _ := http.NewRequest("GET", "/api/reviews/pending/business", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))

	mockService.AssertExpectations(t)
}

// 测试10: 获取待政工审核列表
func TestGetPendingPoliticalReviews_Success(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	submitter := &entity.User{
		ID:   1,
		Name: "提交人",
	}

	pendingPapers := []*entity.Paper{
		{
			ID:        3,
			Title:     "待政审论文",
			Status:    "待政工审核",
			Submitter: submitter,
		},
	}

	mockService.On("GetPendingPapersForPoliticalReview").Return(pendingPapers, nil)

	req, _ := http.NewRequest("GET", "/api/reviews/pending/political", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Equal(t, 1, len(data))

	mockService.AssertExpectations(t)
}

// 测试11: 获取我的审核记录
func TestGetMyReviews_Success(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	reviewLogs := []*entity.ReviewLog{}

	mockService.On("GetMyReviews", uint(2)).Return(reviewLogs, nil)

	req, _ := http.NewRequest("GET", "/api/reviews/my", nil)
	req.Header.Set("X-User-ID", "2")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Equal(t, 0, len(data))

	mockService.AssertExpectations(t)
}

// 测试12: 审核参数校验 - 缺少result字段
func TestReview_MissingResultField(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	requestBody := map[string]interface{}{
		"comment": "只有评论没有结果",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/reviews/business/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "2")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["msg"], "参数错误")
}

// 测试13: 获取审核记录 - 格式错误的论文ID
func TestGetReviewLogs_InvalidPaperID(t *testing.T) {
	mockService := new(MockReviewService)
	router := setupTestReviewRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/reviews/invalid/logs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["msg"], "格式错误")
}
