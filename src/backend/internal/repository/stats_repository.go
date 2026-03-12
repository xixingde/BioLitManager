package repository

import (
	"biolitmanager/internal/model/entity"

	"gorm.io/gorm"
)

// StatsRepository 统计仓储
type StatsRepository struct {
	db *gorm.DB
}

// NewStatsRepository 创建统计仓储实例
func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{
		db: db,
	}
}

// BasicStats 基础统计数据
type BasicStats struct {
	TotalPapers        int64            // 论文总数
	YearlyCounts       map[int]int64    // 各年份论文数量
	TypeCounts         map[string]int64 // 各收录类型论文数量
	JournalCounts      map[string]int64 // 各期刊论文数量
	AvgImpactFactor    float64          // 平均影响因子
	TotalCitations     int64            // 总引用次数
	TotalSelfCitations int64            // 总他引次数
}

// GetBasicStats 获取基础统计数据
func (r *StatsRepository) GetBasicStats() (*BasicStats, error) {
	stats := &BasicStats{
		YearlyCounts:  make(map[int]int64),
		TypeCounts:    make(map[string]int64),
		JournalCounts: make(map[string]int64),
	}

	// 论文总数
	if err := r.db.Model(&entity.Paper{}).Count(&stats.TotalPapers).Error; err != nil {
		return nil, err
	}

	// 各年份发表论文数量
	var yearlyResults []struct {
		Year  int
		Count int64
	}
	if err := r.db.Model(&entity.Paper{}).
		Select("YEAR(publish_date) as year, COUNT(*) as count").
		Group("YEAR(publish_date)").
		Order("year").
		Scan(&yearlyResults).Error; err != nil {
		return nil, err
	}
	for _, result := range yearlyResults {
		stats.YearlyCounts[result.Year] = result.Count
	}

	// 各收录类型论文数量
	var sciCount, eiCount, ciCount, diCount, coreCount int64
	r.db.Model(&entity.Paper{}).Where("is_sci = ?", true).Count(&sciCount)
	r.db.Model(&entity.Paper{}).Where("is_ei = ?", true).Count(&eiCount)
	r.db.Model(&entity.Paper{}).Where("is_ci = ?", true).Count(&ciCount)
	r.db.Model(&entity.Paper{}).Where("is_di = ?", true).Count(&diCount)
	r.db.Model(&entity.Paper{}).Where("is_core = ?", true).Count(&coreCount)
	stats.TypeCounts["SCI"] = sciCount
	stats.TypeCounts["EI"] = eiCount
	stats.TypeCounts["CI"] = ciCount
	stats.TypeCounts["DI"] = diCount
	stats.TypeCounts["CORE"] = coreCount

	// 各期刊发表论文数量
	var journalResults []struct {
		JournalName string
		Count       int64
	}
	if err := r.db.Model(&entity.Paper{}).
		Select("journal_name, COUNT(*) as count").
		Where("journal_name != ''").
		Group("journal_name").
		Order("count DESC").
		Scan(&journalResults).Error; err != nil {
		return nil, err
	}
	for _, result := range journalResults {
		stats.JournalCounts[result.JournalName] = result.Count
	}

	// 平均影响因子
	var avgIF float64
	if err := r.db.Model(&entity.Paper{}).Select("AVG(impact_factor)").Scan(&avgIF).Error; err != nil {
		return nil, err
	}
	stats.AvgImpactFactor = avgIF

	// 总引用次数
	if err := r.db.Model(&entity.Paper{}).Select("COALESCE(SUM(citation_count), 0)").Scan(&stats.TotalCitations).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// AuthorStats 作者统计数据
type AuthorStats struct {
	Author             *entity.Author // 作者信息
	PaperCount         int64          // 论文数量
	FirstAuthorCount   int64          // 第一作者数量
	CorrespondingCount int64          // 通讯作者数量
	AvgImpactFactor    float64        // 平均影响因子
	TotalCitations     int64          // 总引用次数
}

// GetAuthorStats 获取作者统计数据
func (r *StatsRepository) GetAuthorStats(authorId uint) (*AuthorStats, error) {
	stats := &AuthorStats{}

	// 查询作者信息
	var author entity.Author
	if err := r.db.First(&author, authorId).Error; err != nil {
		return nil, err
	}
	stats.Author = &author

	// 论文数量（该作者参与的所有论文）
	if err := r.db.Model(&entity.Author{}).Where("user_id = ?", author.UserID).Count(&stats.PaperCount).Error; err != nil {
		return nil, err
	}

	// 第一作者数量
	if err := r.db.Model(&entity.Author{}).
		Where("user_id = ? AND author_type = ?", author.UserID, "first").
		Count(&stats.FirstAuthorCount).Error; err != nil {
		return nil, err
	}

	// 通讯作者数量
	if err := r.db.Model(&entity.Author{}).
		Where("user_id = ? AND author_type = ?", author.UserID, "corresponding").
		Count(&stats.CorrespondingCount).Error; err != nil {
		return nil, err
	}

	// 获取该作者参与的论文ID列表
	var paperIDs []uint
	if err := r.db.Model(&entity.Author{}).
		Select("paper_id").
		Where("user_id = ?", author.UserID).
		Pluck("paper_id", &paperIDs).Error; err != nil {
		return nil, err
	}

	if len(paperIDs) > 0 {
		// 平均影响因子
		var avgIF float64
		if err := r.db.Model(&entity.Paper{}).
			Where("id IN ?", paperIDs).
			Select("AVG(impact_factor)").
			Scan(&avgIF).Error; err != nil {
			return nil, err
		}
		stats.AvgImpactFactor = avgIF

		// 总引用次数
		if err := r.db.Model(&entity.Paper{}).
			Where("id IN ?", paperIDs).
			Select("COALESCE(SUM(citation_count), 0)").
			Scan(&stats.TotalCitations).Error; err != nil {
			return nil, err
		}
	}

	return stats, nil
}

// ProjectStats 课题统计数据
type ProjectStats struct {
	Project         *entity.Project // 课题信息
	PaperCount      int64           // 论文数量
	HighImpactCount int64           // 高影响因子论文数量（IF >= 5）
	SCIPaperCount   int64           // SCI论文数量
}

// GetProjectStats 获取课题统计数据
func (r *StatsRepository) GetProjectStats(projectId uint) (*ProjectStats, error) {
	stats := &ProjectStats{}

	// 查询课题信息
	var project entity.Project
	if err := r.db.First(&project, projectId).Error; err != nil {
		return nil, err
	}
	stats.Project = &project

	// 论文数量（该课题关联的所有论文）
	if err := r.db.Model(&entity.Paper{}).
		Joins("JOIN paper_projects ON papers.id = paper_projects.paper_id").
		Where("paper_projects.project_id = ?", projectId).
		Count(&stats.PaperCount).Error; err != nil {
		return nil, err
	}

	// 高影响因子论文数量（IF >= 5）
	if err := r.db.Model(&entity.Paper{}).
		Joins("JOIN paper_projects ON papers.id = paper_projects.paper_id").
		Where("paper_projects.project_id = ? AND impact_factor >= ?", projectId, 5.0).
		Count(&stats.HighImpactCount).Error; err != nil {
		return nil, err
	}

	// SCI论文数量
	if err := r.db.Model(&entity.Paper{}).
		Joins("JOIN paper_projects ON papers.id = paper_projects.paper_id").
		Where("paper_projects.project_id = ? AND is_sci = ?", projectId, true).
		Count(&stats.SCIPaperCount).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// DepartmentStats 单位统计数据
type DepartmentStats struct {
	Department        string  // 单位名称
	PaperCount        int64   // 论文数量
	TotalImpactFactor float64 // 总影响因子
	TotalCitations    int64   // 总引用次数
}

// GetDepartmentStats 获取单位统计数据
func (r *StatsRepository) GetDepartmentStats(department string) (*DepartmentStats, error) {
	stats := &DepartmentStats{
		Department: department,
	}

	// 获取该单位所有作者的UserID
	var userIDs []uint
	if err := r.db.Model(&entity.User{}).
		Select("id").
		Where("department = ?", department).
		Pluck("id", &userIDs).Error; err != nil {
		return nil, err
	}

	if len(userIDs) == 0 {
		return stats, nil
	}

	// 获取该单位作者参与的论文ID
	var paperIDs []uint
	if err := r.db.Model(&entity.Author{}).
		Select("DISTINCT paper_id").
		Where("user_id IN ?", userIDs).
		Pluck("paper_id", &paperIDs).Error; err != nil {
		return nil, err
	}

	if len(paperIDs) == 0 {
		return stats, nil
	}

	// 论文数量
	if err := r.db.Model(&entity.Paper{}).
		Where("id IN ?", paperIDs).
		Count(&stats.PaperCount).Error; err != nil {
		return nil, err
	}

	// 总影响因子
	if err := r.db.Model(&entity.Paper{}).
		Where("id IN ?", paperIDs).
		Select("COALESCE(SUM(impact_factor), 0)").
		Scan(&stats.TotalImpactFactor).Error; err != nil {
		return nil, err
	}

	// 总引用次数
	if err := r.db.Model(&entity.Paper{}).
		Where("id IN ?", paperIDs).
		Select("COALESCE(SUM(citation_count), 0)").
		Scan(&stats.TotalCitations).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// YearlyStats 年度统计
type YearlyStats struct {
	Year  int   // 年份
	Count int64 // 论文数量
}

// GetYearlyStats 获取年度统计数据
func (r *StatsRepository) GetYearlyStats() ([]*YearlyStats, error) {
	var results []*YearlyStats

	if err := r.db.Model(&entity.Paper{}).
		Select("YEAR(publish_date) as year, COUNT(*) as count").
		Group("YEAR(publish_date)").
		Order("year DESC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// JournalStats 期刊统计
type JournalStats struct {
	JournalName     string  // 期刊名称
	PaperCount      int64   // 论文数量
	AvgImpactFactor float64 // 平均影响因子
}

// GetJournalStats 获取期刊统计数据
func (r *StatsRepository) GetJournalStats() ([]*JournalStats, error) {
	var results []*JournalStats

	if err := r.db.Model(&entity.Paper{}).
		Select("journal_name, COUNT(*) as paper_count, AVG(impact_factor) as avg_impact_factor").
		Where("journal_name != ''").
		Group("journal_name").
		Order("paper_count DESC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
