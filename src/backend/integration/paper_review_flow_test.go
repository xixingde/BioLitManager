package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"biolitmanager/internal/handler"
	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// PaperReviewFlowIntegrationTestSuite 论文审核流程集成测试套件
type PaperReviewFlowIntegrationTestSuite struct {
	suite.Suite
	db              *gorm.DB
	router          *gin.Engine
	testUserID      uint
	businessUserID  uint
	politicalUserID uint
	journalID       uint
	paperID         uint
}

// SetupSuite 测试套件初始化
func (s *PaperReviewFlowIntegrationTestSuite) SetupSuite() {
	// 使用内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)

	s.db = db

	// 自动迁移所有表
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Paper{},
		&entity.Author{},
		&entity.Journal{},
		&entity.Project{},
		&entity.PaperProject{},
		&entity.Attachment{},
		&entity.ReviewLog{},
		&entity.Archive{},
		&entity.OperationLog{},
	)
	s.Require().NoError(err)

	// 初始化测试数据
	s.initTestData()

	// 初始化所有repositories
	paperRepo := repository.NewPaperRepository(db)
	authorRepo := repository.NewAuthorRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)
	paperProjectRepo := repository.NewPaperProjectRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	userRepo := repository.NewUserRepository(db)
	journalRepo := repository.NewJournalRepository(db)
	archiveRepo := repository.NewArchiveRepository(db)
	operationLogRepo := repository.NewOperationLogRepository(db)

	// 初始化基础服务（无依赖或依赖少）
	operationLogService := service.NewOperationLogService(operationLogRepo)
	notificationService := service.NewNotificationService(userRepo, operationLogService)
	archiveService := service.NewArchiveService(archiveRepo, paperRepo, operationLogService)
	journalService := service.NewJournalService(journalRepo, operationLogService)

	// 初始化核心服务
	paperService := service.NewPaperService(db, paperRepo, authorRepo, attachmentRepo, paperProjectRepo, operationLogService)
	reviewService := service.NewReviewService(db, reviewRepo, paperRepo, userRepo, operationLogService, notificationService, archiveService)

	// 设置路由
	gin.SetMode(gin.TestMode)
	s.router = gin.New()

	// 论文路由
	paperHandler := handler.NewPaperHandler(paperService)
	s.router.POST("/api/papers", paperHandler.CreatePaper)
	s.router.GET("/api/papers/:id", paperHandler.GetPaper)
	s.router.PUT("/api/papers/:id", paperHandler.UpdatePaper)
	s.router.DELETE("/api/papers/:id", paperHandler.DeletePaper)
	s.router.POST("/api/papers/:id/submit", paperHandler.SubmitForReview)
	s.router.POST("/api/papers/:id/save-draft", paperHandler.SaveDraft)

	// 审核路由
	reviewHandler := handler.NewReviewHandler(reviewService)
	s.router.POST("/api/reviews/business/:paperId", reviewHandler.BusinessReview)
	s.router.POST("/api/reviews/political/:paperId", reviewHandler.PoliticalReview)
	s.router.GET("/api/reviews/:paperId/logs", reviewHandler.GetReviewLogs)
	s.router.GET("/api/reviews/pending/business", reviewHandler.GetPendingBusinessReviews)
	s.router.GET("/api/reviews/pending/political", reviewHandler.GetPendingPoliticalReviews)

	// 期刊路由
	journalHandler := handler.NewJournalHandler(journalService)
	s.router.GET("/api/journals/:id", journalHandler.GetJournal)
}

// initTestData 初始化测试数据
func (s *PaperReviewFlowIntegrationTestSuite) initTestData() {
	// 创建测试用户
	testUser := &entity.User{
		Username:     "testuser",
		PasswordHash: "hashed_password",
		Name:         "测试用户",
		Role:         "用户",
	}
	s.db.Create(testUser)
	s.testUserID = testUser.ID

	businessUser := &entity.User{
		Username:     "business_reviewer",
		PasswordHash: "hashed_password",
		Name:         "业务审核员",
		Role:         "业务审核员",
	}
	s.db.Create(businessUser)
	s.businessUserID = businessUser.ID

	politicalUser := &entity.User{
		Username:     "political_reviewer",
		PasswordHash: "hashed_password",
		Name:         "政工审核员",
		Role:         "政工审核员",
	}
	s.db.Create(politicalUser)
	s.politicalUserID = politicalUser.ID

	// 创建测试期刊
	journal := &entity.Journal{
		FullName:     "Nature",
		ShortName:    "Nature",
		ISSN:         "0028-0836",
		ImpactFactor: 69.504,
		Publisher:    "Nature Publishing Group",
	}
	s.db.Create(journal)
	s.journalID = journal.ID

	// 创建测试项目
	project := &entity.Project{
		Name:        "国家重点实验室项目",
		Code:        "NST-2024-001",
		ProjectType: "科研项目",
		Source:      "国家级",
		Level:       "重点",
	}
	s.db.Create(project)
}

// TearDownSuite 测试套件清理
func (s *PaperReviewFlowIntegrationTestSuite) TearDownSuite() {
	// 清理数据库
	sqlDB, _ := s.db.DB()
	sqlDB.Close()
}

// TestCompleteFlow 测试完整的论文审核流程
func (s *PaperReviewFlowIntegrationTestSuite) TestCompleteFlow() {
	t := s.T()

	// 步骤1: 创建论文（草稿状态）
	s.Run("Step1_CreatePaper", func() {
		requestBody := map[string]interface{}{
			"title":         "人工智能在生物学中的应用研究",
			"abstract":      "本文探讨了人工智能技术在生物学领域的最新应用...",
			"journal_id":    s.journalID,
			"doi":           "10.1038/s41586-024-xxxx",
			"impact_factor": 5.234,
			"publish_date":  "2024-01-15",
			"authors": []map[string]interface{}{
				{
					"name":        "张三",
					"author_type": "第一作者",
					"rank":        1,
					"department":  "北京大学",
				},
			},
			"projects": []interface{}{},
		}

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/papers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.testUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		s.paperID = uint(response["id"].(float64))
		assert.Greater(t, s.paperID, uint(0))
	})

	// 步骤2: 验证论文状态为草稿
	s.Run("Step2_VerifyDraftStatus", func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/papers/%d", s.paperID), nil)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		paper := response["data"].(map[string]interface{})
		assert.Equal(t, "draft", paper["status"])
	})

	// 步骤3: 更新草稿
	s.Run("Step3_UpdateDraft", func() {
		requestBody := map[string]interface{}{
			"title":      "人工智能在生物学中的应用研究（修订版）",
			"abstract":   "本文探讨了人工智能技术在生物学领域的最新应用，重点分析了机器学习在基因组学中的应用...",
			"journal_id": s.journalID,
		}

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/papers/%d", s.paperID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.testUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 步骤4: 提交审核
	s.Run("Step4_SubmitForReview", func() {
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/papers/%d/submit", s.paperID), nil)
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.testUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 步骤5: 验证论文状态为"待业务审核"
	s.Run("Step5_VerifyPendingBusinessReview", func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/papers/%d", s.paperID), nil)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		paper := response["data"].(map[string]interface{})
		assert.Equal(t, "待业务审核", paper["status"])
	})

	// 步骤6: 获取待业务审核列表
	s.Run("Step6_GetPendingBusinessReviews", func() {
		req, _ := http.NewRequest("GET", "/api/reviews/pending/business", nil)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].([]interface{})
		assert.GreaterOrEqual(t, len(data), 1)
	})

	// 步骤7: 业务审核通过
	s.Run("Step7_BusinessReviewApprove", func() {
		requestBody := map[string]interface{}{
			"result":  "通过",
			"comment": "论文内容完整，格式符合要求",
		}

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/reviews/business/%d", s.paperID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.businessUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 步骤8: 验证论文状态为"待政工审核"
	s.Run("Step8_VerifyPendingPoliticalReview", func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/papers/%d", s.paperID), nil)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		paper := response["data"].(map[string]interface{})
		assert.Equal(t, "待政工审核", paper["status"])
	})

	// 步骤9: 获取待政工审核列表
	s.Run("Step9_GetPendingPoliticalReviews", func() {
		req, _ := http.NewRequest("GET", "/api/reviews/pending/political", nil)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].([]interface{})
		assert.GreaterOrEqual(t, len(data), 1)
	})

	// 步骤10: 政工审核通过
	s.Run("Step10_PoliticalReviewApprove", func() {
		requestBody := map[string]interface{}{
			"result":  "通过",
			"comment": "内容符合政治审核要求",
		}

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/reviews/political/%d", s.paperID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.politicalUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 步骤11: 验证论文状态为"审核通过"并已归档
	s.Run("Step11_VerifyApprovedAndArchived", func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/papers/%d", s.paperID), nil)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		paper := response["data"].(map[string]interface{})
		assert.Equal(t, "审核通过", paper["status"])
	})

	// 步骤12: 查看完整的审核记录
	s.Run("Step12_GetReviewLogs", func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/reviews/%d/logs", s.paperID), nil)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		logs := response["data"].([]interface{})
		assert.Equal(t, 2, len(logs)) // 业务审核 + 政工审核
	})
}

// TestRejectFlow 测试驳回流程
func (s *PaperReviewFlowIntegrationTestSuite) TestRejectFlow() {
	t := s.T()

	var rejectedPaperID uint

	// 步骤1: 创建论文
	s.Run("RejectStep1_CreatePaper", func() {
		requestBody := map[string]interface{}{
			"title":      "待驳回的论文",
			"abstract":   "这是一篇测试论文",
			"journal_id": s.journalID,
		}

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/papers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.testUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		rejectedPaperID = uint(response["id"].(float64))
	})

	// 步骤2: 提交审核
	s.Run("RejectStep2_SubmitForReview", func() {
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/papers/%d/submit", rejectedPaperID), nil)
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.testUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 步骤3: 业务审核驳回
	s.Run("RejectStep3_BusinessReviewReject", func() {
		requestBody := map[string]interface{}{
			"result":  "驳回",
			"comment": "论文格式不符合要求，请重新整理",
		}

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/reviews/business/%d", rejectedPaperID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.businessUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 步骤4: 验证论文状态为"驳回"
	s.Run("RejectStep4_VerifyRejectedStatus", func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/papers/%d", rejectedPaperID), nil)

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		paper := response["data"].(map[string]interface{})
		assert.Equal(t, "驳回", paper["status"])
	})
}

// TestInvalidStatusTransitions 测试无效的状态转换
func (s *PaperReviewFlowIntegrationTestSuite) TestInvalidStatusTransitions() {
	t := s.T()

	var testPaperID uint

	// 创建并提交论文
	s.Run("InvalidStep1_CreateAndSubmitPaper", func() {
		requestBody := map[string]interface{}{
			"title":      "状态转换测试论文",
			"abstract":   "测试论文",
			"journal_id": s.journalID,
		}

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/papers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.testUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		testPaperID = uint(response["id"].(float64))

		// 提交审核
		req, _ = http.NewRequest("POST", fmt.Sprintf("/api/papers/%d/submit", testPaperID), nil)
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.testUserID))

		w = httptest.NewRecorder()
		s.router.ServeHTTP(w, req)
	})

	// 尝试编辑非草稿状态的论文（应失败）
	s.Run("InvalidStep2_TryEditNonDraft", func() {
		requestBody := map[string]interface{}{
			"title": "修改后的标题",
		}

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/papers/%d", testPaperID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.testUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// 尝试删除非草稿状态的论文（应失败）
	s.Run("InvalidStep3_TryDeleteNonDraft", func() {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/papers/%d", testPaperID), nil)
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.testUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// 尝试重复提交审核（应失败）
	s.Run("InvalidStep4_TrySubmitAgain", func() {
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/papers/%d/submit", testPaperID), nil)
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", s.testUserID))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestIntegration 运行集成测试
func TestPaperReviewFlowIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PaperReviewFlowIntegrationTestSuite))
}
