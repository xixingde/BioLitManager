package handler

import (
	"net/http"
	"strconv"

	"biolitmanager/internal/model/dto/request"
	"biolitmanager/internal/model/dto/response"
	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/service"
	"biolitmanager/pkg/logger"
	resp "biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SearchHandler 搜索处理器
type SearchHandler struct {
	searchService service.SearchServiceInterface
	paperService  service.PaperServiceInterface
}

// NewSearchHandler 创建搜索处理器实例
func NewSearchHandler(searchService service.SearchServiceInterface, paperService service.PaperServiceInterface) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		paperService:  paperService,
	}
}

// Search 处理搜索请求
// GET /api/search
func (h *SearchHandler) Search(c *gin.Context) {
	// 解析查询参数
	req, err := h.parseSearchRequest(c)
	if err != nil {
		logger.GetLogger().Warn("Failed to parse search request",
			zap.Error(err),
		)
		resp.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 调用搜索服务
	result, err := h.searchService.AdvancedSearch(req)
	if err != nil {
		logger.GetLogger().Error("Search failed",
			zap.Error(err),
		)
		resp.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	resp.Success(c, result)
}

// GetPaperDetail 获取论文详情
// GET /api/search/papers/:id
func (h *SearchHandler) GetPaperDetail(c *gin.Context) {
	id := c.Param("id")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	// 调用paperService获取论文详情
	paper, err := h.paperService.GetPaperByID(uint(paperID))
	if err != nil {
		if err == service.ErrPaperNotFound {
			resp.Error(c, http.StatusNotFound, "论文不存在")
			return
		}
		logger.GetLogger().Error("Failed to get paper detail",
			zap.Error(err),
			zap.Uint("paper_id", uint(paperID)),
		)
		resp.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为详情响应格式
	paperDetail := h.convertToPaperDTO(paper)

	resp.Success(c, paperDetail)
}

// parseSearchRequest 解析搜索请求参数
func (h *SearchHandler) parseSearchRequest(c *gin.Context) (*request.SearchRequest, error) {
	req := &request.SearchRequest{
		Pagination: request.Pagination{
			Page:     1,
			PageSize: 20,
		},
	}

	// 解析分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Pagination.Page = page
		}
	}
	if sizeStr := c.Query("page_size"); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 && size <= 100 {
			req.Pagination.PageSize = size
		}
	}

	// 解析基础查询条件
	req.Keywords = c.Query("keywords")
	req.Title = c.Query("title")
	req.Abstract = c.Query("abstract")
	req.DOI = c.Query("doi")

	// 解析期刊信息
	if journalIDStr := c.Query("journal_id"); journalIDStr != "" {
		if journalID, err := strconv.ParseUint(journalIDStr, 10, 32); err == nil {
			req.JournalID = uint(journalID)
		}
	}
	req.JournalName = c.Query("journal_name")

	// 解析作者信息
	req.AuthorName = c.Query("author_name")
	req.AuthorType = request.AuthorType(c.Query("author_type"))
	req.Department = c.Query("department")

	// 解析项目信息
	req.ProjectCode = c.Query("project_code")
	req.ProjectType = c.Query("project_type")

	// 解析日期范围
	req.PublishDateStart = c.Query("publish_date_start")
	req.PublishDateEnd = c.Query("publish_date_end")

	// 解析影响因子范围
	if ifMinStr := c.Query("impact_factor_min"); ifMinStr != "" {
		if ifMin, err := strconv.ParseFloat(ifMinStr, 64); err == nil {
			req.ImpactFactorMin = ifMin
		}
	}
	if ifMaxStr := c.Query("impact_factor_max"); ifMaxStr != "" {
		if ifMax, err := strconv.ParseFloat(ifMaxStr, 64); err == nil {
			req.ImpactFactorMax = ifMax
		}
	}

	// 解析状态
	req.Status = c.Query("status")

	// 解析排序参数
	req.SortField = request.SortField(c.Query("sort_field"))
	req.SortOrder = request.SortOrder(c.Query("sort_order"))

	return req, nil
}

// convertToPaperDTO 将Paper实体转换为DTO格式
func (h *SearchHandler) convertToPaperDTO(paper *entity.Paper) *response.PaperDTO {
	paperDTO := &response.PaperDTO{
		ID:           paper.ID,
		Title:        paper.Title,
		Abstract:     paper.Abstract,
		DOI:          paper.DOI,
		ImpactFactor: paper.ImpactFactor,
		PublishDate:  paper.PublishDate,
		Status:       paper.Status,
		SubmitTime:   paper.SubmitTime,
		CreatedAt:    paper.CreatedAt,
		UpdatedAt:    paper.UpdatedAt,
	}

	// 转换期刊信息
	if paper.Journal != nil {
		paperDTO.Journal = &response.JournalDTO{
			ID:           paper.Journal.ID,
			FullName:     paper.Journal.FullName,
			ShortName:    paper.Journal.ShortName,
			ISSN:         paper.Journal.ISSN,
			ImpactFactor: paper.Journal.ImpactFactor,
			Publisher:    paper.Journal.Publisher,
		}
	}

	// 转换提交人信息
	if paper.Submitter != nil {
		paperDTO.Submitter = &response.UserDTO{
			ID:       paper.Submitter.ID,
			Username: paper.Submitter.Username,
			Name:     paper.Submitter.Name,
			Role:     paper.Submitter.Role,
		}
	}

	// 转换作者信息
	if paper.Authors != nil && len(paper.Authors) > 0 {
		authors := make([]response.AuthorDTO, len(paper.Authors))
		for i, author := range paper.Authors {
			authors[i] = response.AuthorDTO{
				ID:         author.ID,
				Name:       author.Name,
				AuthorType: author.AuthorType,
				Rank:       author.Rank,
				Department: author.Department,
				UserID:     author.UserID,
			}
		}
		paperDTO.Authors = authors
	}

	// 转换课题信息
	if paper.Projects != nil && len(paper.Projects) > 0 {
		projects := make([]response.ProjectDTO, len(paper.Projects))
		for i, proj := range paper.Projects {
			projects[i] = response.ProjectDTO{
				ID:          proj.ID,
				Name:        proj.Name,
				Code:        proj.Code,
				ProjectType: proj.ProjectType,
				Source:      proj.Source,
				Level:       proj.Level,
			}
		}
		paperDTO.Projects = projects
	}

	// 转换附件信息
	if paper.Attachments != nil && len(paper.Attachments) > 0 {
		attachments := make([]response.AttachmentDTO, len(paper.Attachments))
		for i, att := range paper.Attachments {
			attachments[i] = response.AttachmentDTO{
				ID:        att.ID,
				FileType:  att.FileType,
				FileName:  att.FileName,
				FilePath:  att.FilePath,
				FileSize:  att.FileSize,
				MimeType:  att.MimeType,
				CreatedAt: att.CreatedAt,
			}
		}
		paperDTO.Attachments = attachments
	}

	return paperDTO
}
