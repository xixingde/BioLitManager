package service

import (
	"fmt"
	"time"

	"biolitmanager/internal/model/dto/request"
	"biolitmanager/internal/model/dto/response"
	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SearchServiceInterface 搜索服务接口
type SearchServiceInterface interface {
	AdvancedSearch(req *request.SearchRequest) (*response.SearchResponse, error)
}

// SearchService 搜索服务
type SearchService struct {
	db               *gorm.DB
	paperRepo        *repository.PaperRepository
	authorRepo       *repository.AuthorRepository
	paperProjectRepo *repository.PaperProjectRepository
}

// NewSearchService 创建搜索服务实例
func NewSearchService(
	db *gorm.DB,
	paperRepo *repository.PaperRepository,
	authorRepo *repository.AuthorRepository,
	paperProjectRepo *repository.PaperProjectRepository,
) *SearchService {
	return &SearchService{
		db:               db,
		paperRepo:        paperRepo,
		authorRepo:       authorRepo,
		paperProjectRepo: paperProjectRepo,
	}
}

// AdvancedSearch 高级搜索
func (s *SearchService) AdvancedSearch(req *request.SearchRequest) (*response.SearchResponse, error) {
	logger.GetLogger().Info("Starting advanced search",
		zap.String("keywords", req.Keywords),
		zap.String("title", req.Title),
		zap.String("author_name", req.AuthorName),
	)

	// 构建查询条件
	queryReq := s.buildQueryConditions(req)

	// 应用排序
	s.applySorting(queryReq)

	// 权限过滤（仅审核通过的论文可公开查询）
	s.filterByPermission(queryReq)

	// 按作者类型过滤
	if req.AuthorType != "" {
		s.filterByAuthorType(queryReq, req.AuthorType)
	}

	// 调用repository的高级搜索
	papers, total, err := s.paperRepo.AdvancedSearch(queryReq)
	if err != nil {
		logger.GetLogger().Error("Failed to execute advanced search",
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	// 转换为响应结构
	results := s.convertToSearchResults(papers)

	// 构建分页信息
	page := req.Pagination.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.Pagination.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	searchResponse := &response.SearchResponse{
		Pagination: response.PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalCount: total,
			TotalPages: totalPages,
		},
		Results: results,
	}

	logger.GetLogger().Info("Advanced search completed",
		zap.Int64("total", total),
		zap.Int("result_count", len(results)),
	)

	return searchResponse, nil
}

// buildQueryConditions 构建查询条件
func (s *SearchService) buildQueryConditions(req *request.SearchRequest) *request.SearchRequest {
	// 复制请求对象，避免修改原始请求
	queryReq := &request.SearchRequest{
		Keywords:         req.Keywords,
		Title:            req.Title,
		Abstract:         req.Abstract,
		DOI:              req.DOI,
		JournalID:        req.JournalID,
		JournalName:      req.JournalName,
		AuthorName:       req.AuthorName,
		AuthorType:       req.AuthorType,
		Department:       req.Department,
		ProjectCode:      req.ProjectCode,
		ProjectType:      req.ProjectType,
		PublishDateStart: req.PublishDateStart,
		PublishDateEnd:   req.PublishDateEnd,
		ImpactFactorMin:  req.ImpactFactorMin,
		ImpactFactorMax:  req.ImpactFactorMax,
		Status:           req.Status,
		QueryGroup:       req.QueryGroup,
		SortField:        req.SortField,
		SortOrder:        req.SortOrder,
		Pagination:       req.Pagination,
	}

	// 处理复杂查询条件组
	if req.QueryGroup != nil {
		s.buildGroupConditions(queryReq, req.QueryGroup)
	}

	// 处理部门搜索
	if req.Department != "" {
		// 部门搜索通过AuthorName间接处理，在repository中会通过作者表关联查询
		req.AuthorName = req.Department
	}

	return queryReq
}

// buildGroupConditions 递归构建条件组
func (s *SearchService) buildGroupConditions(req *request.SearchRequest, group *request.QueryGroup) {
	// 处理当前层级的条件
	for _, condition := range group.Conditions {
		switch condition.Field {
		case "title":
			req.Title = condition.Value
		case "abstract":
			req.Abstract = condition.Value
		case "keywords":
			req.Keywords = condition.Value
		case "doi":
			req.DOI = condition.Value
		case "journal_name":
			req.JournalName = condition.Value
		case "author_name":
			req.AuthorName = condition.Value
		case "project_code":
			req.ProjectCode = condition.Value
		case "publish_date_start":
			req.PublishDateStart = condition.Value
		case "publish_date_end":
			req.PublishDateEnd = condition.Value
		case "impact_factor_min":
			fmt.Sscanf(condition.Value, "%f", &req.ImpactFactorMin)
		case "impact_factor_max":
			fmt.Sscanf(condition.Value, "%f", &req.ImpactFactorMax)
		}
	}

	// 递归处理嵌套条件组
	for _, subGroup := range group.Groups {
		s.buildGroupConditions(req, &subGroup)
	}
}

// applySorting 应用排序
func (s *SearchService) applySorting(req *request.SearchRequest) {
	// 设置默认排序
	if req.SortField == "" {
		req.SortField = request.SortFieldPublishDate
	}
	if req.SortOrder == "" {
		req.SortOrder = request.SortOrderDesc
	}

	logger.GetLogger().Debug("Sorting applied",
		zap.String("sort_field", string(req.SortField)),
		zap.String("sort_order", string(req.SortOrder)),
	)
}

// filterByPermission 权限过滤（仅审核通过的论文可公开查询）
func (s *SearchService) filterByPermission(req *request.SearchRequest) {
	// 默认只返回审核通过的论文
	// 如果请求中指定了状态，则使用请求的状态
	if req.Status == "" {
		req.Status = "审核通过"
	}

	logger.GetLogger().Debug("Permission filter applied",
		zap.String("status", req.Status),
	)
}

// filterByAuthorType 按作者类型过滤
func (s *SearchService) filterByAuthorType(req *request.SearchRequest, authorType request.AuthorType) {
	switch authorType {
	case request.AuthorTypeFirst:
		// 第一作者：rank = 1，在repository中通过Author表关联查询
		logger.GetLogger().Debug("Filtering by first author (rank = 1)")
	case request.AuthorTypeCoFirst:
		// 共同第一作者：通过Paper.IsCoFirstAuthor字段过滤
		logger.GetLogger().Debug("Filtering by co-first author (Paper.IsCoFirstAuthor)")
	case request.AuthorTypeCorresponding:
		// 通讯作者：通过Paper.IsCorrespondingAuthor字段过滤
		logger.GetLogger().Debug("Filtering by corresponding author (Paper.IsCorrespondingAuthor)")
	}

	logger.GetLogger().Info("Author type filter applied",
		zap.String("author_type", string(authorType)),
	)
}

// convertToSearchResults 转换为搜索结果
func (s *SearchService) convertToSearchResults(papers []*entity.Paper) []response.PaperSearchResult {
	results := make([]response.PaperSearchResult, 0, len(papers))

	for _, paper := range papers {
		result := response.PaperSearchResult{
			ID:           paper.ID,
			Title:        paper.Title,
			Abstract:     paper.Abstract,
			DOI:          paper.DOI,
			ImpactFactor: paper.ImpactFactor,
			PublishDate:  paper.PublishDate,
			Partition:    paper.Partition,
			Status:       paper.Status,
			Citation:     paper.CitationCount,
			CreatedAt:    paper.CreatedAt,
		}

		// 填充期刊信息
		if paper.Journal != nil {
			result.JournalName = paper.Journal.FullName
			result.JournalShort = paper.Journal.ShortName
		}

		// 填充作者信息
		authors := make([]response.AuthorInfo, 0, len(paper.Authors))
		for _, author := range paper.Authors {
			authorType := "coauthor"
			if author.Rank == 1 {
				authorType = "first"
			}
			if paper.IsCoFirstAuthor && author.Rank == 1 {
				authorType = "co_first"
			}
			if paper.IsCorrespondingAuthor && author.Rank == len(paper.Authors) {
				authorType = "corresponding"
			}

			authors = append(authors, response.AuthorInfo{
				Name:       author.Name,
				Department: author.Department,
				AuthorType: authorType,
				Rank:       author.Rank,
			})
		}
		result.Authors = authors

		// 填充课题信息
		projects := make([]response.ProjectInfo, 0, len(paper.Projects))
		for _, proj := range paper.Projects {
			projects = append(projects, response.ProjectInfo{
				Code:        proj.Code,
				ProjectType: proj.ProjectType,
				Name:        proj.Name,
			})
		}
		result.Projects = projects

		results = append(results, result)
	}

	return results
}

// buildComplexQuery 构建复杂查询（支持AND/OR/NOT逻辑）
func (s *SearchService) buildComplexQuery(db *gorm.DB, group *request.QueryGroup) *gorm.DB {
	if group == nil {
		return db
	}

	// 处理当前层级的条件
	conditions := group.Conditions
	groups := group.Groups

	// 递归构建子查询
	for _, subGroup := range groups {
		db = s.buildComplexQuery(db, &subGroup)
	}

	// 应用当前层级的逻辑
	switch group.Logic {
	case request.LogicAnd:
		// AND逻辑：所有条件都需要满足
		for _, cond := range conditions {
			db = s.applyCondition(db, cond, "AND")
		}
	case request.LogicOr:
		// OR逻辑：满足任一条件即可
		for _, cond := range conditions {
			db = s.applyCondition(db, cond, "OR")
		}
	case request.LogicNot:
		// NOT逻辑：排除满足条件的记录
		for _, cond := range conditions {
			db = s.applyCondition(db, cond, "NOT")
		}
	}

	return db
}

// applyCondition 应用单个查询条件
func (s *SearchService) applyCondition(db *gorm.DB, cond request.QueryCondition, logic string) *gorm.DB {
	field := cond.Field
	value := cond.Value
	operator := cond.Operator

	if operator == "" {
		operator = "eq"
	}

	switch operator {
	case "eq":
		// 精确匹配
		db = db.Where(fmt.Sprintf("%s = ?", field), value)
	case "like":
		// 模糊匹配
		db = db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%")
	case "gt":
		// 大于
		db = db.Where(fmt.Sprintf("%s > ?", field), value)
	case "lt":
		// 小于
		db = db.Where(fmt.Sprintf("%s < ?", field), value)
	case "gte":
		// 大于等于
		db = db.Where(fmt.Sprintf("%s >= ?", field), value)
	case "lte":
		// 小于等于
		db = db.Where(fmt.Sprintf("%s <= ?", field), value)
	case "in":
		// IN查询
		db = db.Where(fmt.Sprintf("%s IN (?)", field), value)
	}

	return db
}

// ValidateSearchRequest 验证搜索请求参数
func (s *SearchService) ValidateSearchRequest(req *request.SearchRequest) error {
	// 验证分页参数
	if req.Pagination.Page < 0 {
		return fmt.Errorf("页码不能为负数")
	}
	if req.Pagination.PageSize < 0 {
		return fmt.Errorf("每页条数不能为负数")
	}
	if req.Pagination.PageSize > 100 {
		return fmt.Errorf("每页条数不能超过100")
	}

	// 验证影响因子范围
	if req.ImpactFactorMin < 0 {
		return fmt.Errorf("影响因子最小值不能为负数")
	}
	if req.ImpactFactorMax < 0 {
		return fmt.Errorf("影响因子最大值不能为负数")
	}
	if req.ImpactFactorMin > req.ImpactFactorMax {
		return fmt.Errorf("影响因子最小值不能大于最大值")
	}

	// 验证日期范围
	if req.PublishDateStart != "" {
		if _, err := time.Parse("2006-01-02", req.PublishDateStart); err != nil {
			return fmt.Errorf("开始日期格式无效，应为YYYY-MM-DD")
		}
	}
	if req.PublishDateEnd != "" {
		if _, err := time.Parse("2006-01-02", req.PublishDateEnd); err != nil {
			return fmt.Errorf("结束日期格式无效，应为YYYY-MM-DD")
		}
	}

	logger.GetLogger().Debug("Search request validated",
		zap.String("keywords", req.Keywords),
	)

	return nil
}
