package repository

import (
	"errors"

	"biolitmanager/internal/model/entity"

	"gorm.io/gorm"
)

// ReviewRepository 审核记录仓储接口
type ReviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository 创建审核记录仓储实例
func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

// Create 创建审核记录
func (r *ReviewRepository) Create(reviewLog *entity.ReviewLog) error {
	return r.db.Create(reviewLog).Error
}

// FindByID 根据ID查询审核记录
func (r *ReviewRepository) FindByID(id uint) (*entity.ReviewLog, error) {
	var reviewLog entity.ReviewLog
	err := r.db.First(&reviewLog, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reviewLog, nil
}

// FindByPaperID 根据论文ID查询该论文的所有审核记录，按时间倒序
func (r *ReviewRepository) FindByPaperID(paperID uint) ([]*entity.ReviewLog, error) {
	var reviewLogs []*entity.ReviewLog
	err := r.db.Where("paper_id = ?", paperID).Order("review_time DESC").Find(&reviewLogs).Error
	if err != nil {
		return nil, err
	}
	return reviewLogs, nil
}

// FindLatestByPaperIDAndType 根据论文ID和审核类型查询最新的审核记录
func (r *ReviewRepository) FindLatestByPaperIDAndType(paperID uint, reviewType string) (*entity.ReviewLog, error) {
	var reviewLog entity.ReviewLog
	err := r.db.Where("paper_id = ? AND review_type = ?", paperID, reviewType).
		Order("review_time DESC").
		First(&reviewLog).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reviewLog, nil
}

// ListByReviewer 根据审核人员ID查询该人员的审核记录
func (r *ReviewRepository) ListByReviewer(reviewerID uint) ([]*entity.ReviewLog, error) {
	var reviewLogs []*entity.ReviewLog
	err := r.db.Where("reviewer_id = ?", reviewerID).Order("review_time DESC").Find(&reviewLogs).Error
	if err != nil {
		return nil, err
	}
	return reviewLogs, nil
}

// ListPendingReview 查询待审核的论文列表（状态为"待业务审核"或"待政工审核"）
func (r *ReviewRepository) ListPendingReview(reviewType string) ([]*entity.Paper, error) {
	var papers []*entity.Paper

	// 根据审核类型确定待审核状态
	var status string
	if reviewType == "业务审核" {
		status = "待业务审核"
	} else if reviewType == "政工审核" {
		status = "待政工审核"
	} else {
		return nil, nil
	}

	err := r.db.Where("status = ?", status).Order("submit_time ASC").Find(&papers).Error
	if err != nil {
		return nil, err
	}
	return papers, nil
}
