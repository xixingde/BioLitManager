package handler

import (
	"net/http"
	"strconv"
	"time"

	reviewresponse "biolitmanager/internal/model/dto/response"
	"biolitmanager/internal/service"
	"biolitmanager/pkg/logger"
	"biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReviewHandler 审核处理器
type ReviewHandler struct {
	reviewService service.ReviewServiceInterface
}

// NewReviewHandler 创建审核处理器实例
func NewReviewHandler(reviewService service.ReviewServiceInterface) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// BusinessReview 业务审核
// POST /api/reviews/business/:paperId
func (h *ReviewHandler) BusinessReview(c *gin.Context) {
	id := c.Param("paperId")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	var req struct {
		Result  string `json:"result" binding:"required"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid business review request",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取审核人信息
	reviewerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	if err := h.reviewService.BusinessReview(uint(paperID), req.Result, req.Comment, reviewerID.(uint), ipAddress); err != nil {
		if err == service.ErrPaperNotFound {
			response.Error(c, http.StatusNotFound, "论文不存在")
			return
		}
		if err == service.ErrInvalidReviewResult {
			response.Error(c, http.StatusBadRequest, "无效的审核结果")
			return
		}
		if err == service.ErrInvalidStatus {
			response.Error(c, http.StatusBadRequest, "论文状态不允许审核")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// PoliticalReview 政工审核
// POST /api/reviews/political/:paperId
func (h *ReviewHandler) PoliticalReview(c *gin.Context) {
	id := c.Param("paperId")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	var req struct {
		Result  string `json:"result" binding:"required"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid political review request",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取审核人信息
	reviewerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	if err := h.reviewService.PoliticalReview(uint(paperID), req.Result, req.Comment, reviewerID.(uint), ipAddress); err != nil {
		if err == service.ErrPaperNotFound {
			response.Error(c, http.StatusNotFound, "论文不存在")
			return
		}
		if err == service.ErrInvalidReviewResult {
			response.Error(c, http.StatusBadRequest, "无效的审核结果")
			return
		}
		if err == service.ErrInvalidStatus {
			response.Error(c, http.StatusBadRequest, "论文状态不允许审核")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// GetReviewLogs 获取审核记录
// GET /api/reviews/:paperId/logs
func (h *ReviewHandler) GetReviewLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("paperId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	reviewLogs, err := h.reviewService.GetReviewLogsByPaperID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	reviewLogDTOs := make([]reviewresponse.ReviewLogDTO, len(reviewLogs))
	for i, log := range reviewLogs {
		reviewLogDTOs[i] = reviewresponse.ReviewLogDTO{
			ID:         log.ID,
			PaperID:    log.PaperID,
			ReviewType: log.ReviewType,
			Result:     log.Result,
			Comment:    log.Comment,
			ReviewTime: log.ReviewTime,
			CreatedAt:  log.CreatedAt,
		}
		if log.Reviewer != nil {
			reviewLogDTOs[i].Reviewer = &reviewresponse.UserDTO{
				ID:       log.Reviewer.ID,
				Username: log.Reviewer.Username,
				Name:     log.Reviewer.Name,
				Role:     log.Reviewer.Role,
			}
		}
	}

	response.Success(c, reviewLogDTOs)
}

// GetPendingBusinessReviews 获取待业务审核的论文列表
// GET /api/reviews/pending/business
func (h *ReviewHandler) GetPendingBusinessReviews(c *gin.Context) {
	papers, err := h.reviewService.GetPendingPapersForBusinessReview()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	pendingReviewDTOs := make([]reviewresponse.PendingReviewDTO, len(papers))
	for i, paper := range papers {
		pendingReviewDTOs[i] = reviewresponse.PendingReviewDTO{
			ID:              paper.ID,
			Title:           paper.Title,
			SubmitterName:   paper.Submitter.Name,
			SubmitTime:      paper.SubmitTime,
			Status:          paper.Status,
			DaysSinceSubmit: int(time.Since(paper.SubmitTime).Hours() / 24),
		}
	}

	response.Success(c, pendingReviewDTOs)
}

// GetPendingPoliticalReviews 获取待政工审核的论文列表
// GET /api/reviews/pending/political
func (h *ReviewHandler) GetPendingPoliticalReviews(c *gin.Context) {
	papers, err := h.reviewService.GetPendingPapersForPoliticalReview()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	pendingReviewDTOs := make([]reviewresponse.PendingReviewDTO, len(papers))
	for i, paper := range papers {
		pendingReviewDTOs[i] = reviewresponse.PendingReviewDTO{
			ID:              paper.ID,
			Title:           paper.Title,
			SubmitterName:   paper.Submitter.Name,
			SubmitTime:      paper.SubmitTime,
			Status:          paper.Status,
			DaysSinceSubmit: int(time.Since(paper.SubmitTime).Hours() / 24),
		}
	}

	response.Success(c, pendingReviewDTOs)
}

// GetMyReviews 获取我的审核记录
// GET /api/reviews/my
func (h *ReviewHandler) GetMyReviews(c *gin.Context) {
	// 获取当前用户ID
	reviewerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	reviewLogs, err := h.reviewService.GetMyReviews(reviewerID.(uint))
	if err != nil {
		logger.GetLogger().Error("Failed to get my reviews",
			zap.Uint("reviewer_id", reviewerID.(uint)),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, reviewLogs)
}
