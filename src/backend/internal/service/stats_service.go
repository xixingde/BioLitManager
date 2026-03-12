package service

import (
	"errors"

	"biolitmanager/internal/repository"

	stderrors "errors"

	"gorm.io/gorm"
)

var (
	// ErrStatsNotFound 统计数据不存在
	ErrStatsNotFound = stderrors.New("统计数据不存在")
)

// StatsServiceInterface 统计服务接口
type StatsServiceInterface interface {
	GetBasicStats() (*repository.BasicStats, error)
	GetStatsByAuthor(authorId uint) (*repository.AuthorStats, error)
	GetStatsByProject(projectId uint) (*repository.ProjectStats, error)
	GetStatsByDepartment(department string) (*repository.DepartmentStats, error)
	GetYearlyStats() ([]*repository.YearlyStats, error)
	GetJournalStats() ([]*repository.JournalStats, error)
}

// StatsService 统计服务
type StatsService struct {
	statsRepo *repository.StatsRepository
}

// NewStatsService 创建统计服务实例
func NewStatsService(statsRepo *repository.StatsRepository) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
	}
}

// GetBasicStats 获取基础指标统计
func (s *StatsService) GetBasicStats() (*repository.BasicStats, error) {
	stats, err := s.statsRepo.GetBasicStats()
	if err != nil {
		return nil, ErrStatsNotFound
	}
	return stats, nil
}

// GetStatsByAuthor 按作者统计
func (s *StatsService) GetStatsByAuthor(authorId uint) (*repository.AuthorStats, error) {
	stats, err := s.statsRepo.GetAuthorStats(authorId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthorNotFound
		}
		return nil, ErrStatsNotFound
	}
	return stats, nil
}

// GetStatsByProject 按课题统计
func (s *StatsService) GetStatsByProject(projectId uint) (*repository.ProjectStats, error) {
	stats, err := s.statsRepo.GetProjectStats(projectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, ErrStatsNotFound
	}
	return stats, nil
}

// GetStatsByDepartment 按单位统计
func (s *StatsService) GetStatsByDepartment(department string) (*repository.DepartmentStats, error) {
	stats, err := s.statsRepo.GetDepartmentStats(department)
	if err != nil {
		return nil, ErrStatsNotFound
	}
	return stats, nil
}

// GetYearlyStats 年度统计
func (s *StatsService) GetYearlyStats() ([]*repository.YearlyStats, error) {
	stats, err := s.statsRepo.GetYearlyStats()
	if err != nil {
		return nil, ErrStatsNotFound
	}
	return stats, nil
}

// GetJournalStats 期刊统计
func (s *StatsService) GetJournalStats() ([]*repository.JournalStats, error) {
	stats, err := s.statsRepo.GetJournalStats()
	if err != nil {
		return nil, ErrStatsNotFound
	}
	return stats, nil
}
