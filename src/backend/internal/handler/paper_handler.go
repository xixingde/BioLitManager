package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	paperresponse "biolitmanager/internal/model/dto/response"
	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/service"
	"biolitmanager/pkg/logger"
	"biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// PaperHandler 论文处理器
type PaperHandler struct {
	paperService service.PaperServiceInterface
}

// NewPaperHandler 创建论文处理器实例
func NewPaperHandler(paperService service.PaperServiceInterface) *PaperHandler {
	return &PaperHandler{
		paperService: paperService,
	}
}

// CreatePaper 创建论文
// POST /api/papers
func (h *PaperHandler) CreatePaper(c *gin.Context) {
	var req struct {
		Title        string           `json:"title" binding:"required"`
		Abstract     string           `json:"abstract"`
		JournalID    uint             `json:"journal_id"`
		DOI          string           `json:"doi"`
		ImpactFactor float64          `json:"impact_factor"`
		PublishDate  string           `json:"publish_date"`
		Authors      []*entity.Author `json:"authors"`
		Projects     []uint           `json:"projects"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid create paper request",
			zap.Error(err),
			zap.String("body", c.GetHeader("Content-Type")),
		)
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 获取操作者信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	// 解析日期
	var publishDate *time.Time
	if req.PublishDate != "" {
		if t, err := time.Parse("2006-01-02", req.PublishDate); err == nil {
			publishDate = &t
		}
	}

	paper, err := h.paperService.CreatePaper(
		req.Title,
		req.Abstract,
		req.JournalID,
		req.DOI,
		req.ImpactFactor,
		publishDate,
		userID.(uint),
		req.Authors,
		req.Projects,
		userID.(uint),
		ipAddress,
	)

	if err != nil {
		if err == service.ErrPaperDuplicate {
			response.Error(c, http.StatusBadRequest, "论文已存在重复")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, gin.H{
		"id": paper.ID,
	})
}

// GetPaper 获取论文详情
// GET /api/papers/:id
func (h *PaperHandler) GetPaper(c *gin.Context) {
	id := c.Param("id")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	paper, err := h.paperService.GetPaperByID(uint(paperID))
	if err != nil {
		if err == service.ErrPaperNotFound {
			response.Error(c, http.StatusNotFound, "论文不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	paperDTO := h.convertToPaperDTO(paper)

	response.Success(c, paperDTO)
}

// ListPapers 分页查询论文列表
// GET /api/papers
func (h *PaperHandler) ListPapers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	papers, total, err := h.paperService.ListPapers(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	paperList := make([]paperresponse.PaperDTO, len(papers))
	for i, paper := range papers {
		paperList[i] = *h.convertToPaperDTO(paper)
	}

	result := paperresponse.PaperListResponse{
		List:  paperList,
		Total: total,
		Page:  page,
		Size:  size,
	}

	response.Success(c, result)
}

// UpdatePaper 更新论文
// PUT /api/papers/:id
func (h *PaperHandler) UpdatePaper(c *gin.Context) {
	id := c.Param("id")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	var req struct {
		Title        string           `json:"title"`
		Abstract     string           `json:"abstract"`
		JournalID    uint             `json:"journal_id"`
		DOI          string           `json:"doi"`
		ImpactFactor float64          `json:"impact_factor"`
		PublishDate  string           `json:"publish_date"`
		Authors      []*entity.Author `json:"authors"`
		Projects     []uint           `json:"projects"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid update paper request",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取操作者信息
	operatorID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	// 解析日期
	var publishDate *time.Time
	if req.PublishDate != "" {
		if t, err := time.Parse("2006-01-02", req.PublishDate); err == nil {
			publishDate = &t
		}
	}

	if err := h.paperService.UpdatePaper(
		uint(paperID),
		req.Title,
		req.Abstract,
		req.JournalID,
		req.DOI,
		req.ImpactFactor,
		publishDate,
		req.Authors,
		req.Projects,
		operatorID.(uint),
		ipAddress,
	); err != nil {
		if err == service.ErrPaperNotFound {
			response.Error(c, http.StatusNotFound, "论文不存在")
			return
		}
		if err == service.ErrPaperDuplicate {
			response.Error(c, http.StatusBadRequest, "论文已存在重复")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// DeletePaper 删除论文
// DELETE /api/papers/:id
func (h *PaperHandler) DeletePaper(c *gin.Context) {
	id := c.Param("id")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	// 获取操作者信息
	operatorID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	if err := h.paperService.DeletePaper(uint(paperID), operatorID.(uint), ipAddress); err != nil {
		if err == service.ErrPaperNotFound {
			response.Error(c, http.StatusNotFound, "论文不存在")
			return
		}
		if err == service.ErrInvalidStatus {
			response.Error(c, http.StatusBadRequest, "论文状态不允许删除,仅草稿状态可删除")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// SubmitForReview 提交审核
// POST /api/papers/:id/submit
func (h *PaperHandler) SubmitForReview(c *gin.Context) {
	id := c.Param("id")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	// 获取操作者信息
	operatorID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	if err := h.paperService.SubmitForReview(uint(paperID), operatorID.(uint), ipAddress); err != nil {
		if err == service.ErrPaperNotFound {
			response.Error(c, http.StatusNotFound, "论文不存在")
			return
		}
		if err == service.ErrInvalidStatus {
			response.Error(c, http.StatusBadRequest, "论文状态不允许提交审核")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// SaveDraft 保存草稿
// POST /api/papers/:id/save-draft
func (h *PaperHandler) SaveDraft(c *gin.Context) {
	id := c.Param("id")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	var req struct {
		Title        string           `json:"title"`
		Abstract     string           `json:"abstract"`
		JournalID    uint             `json:"journal_id"`
		DOI          string           `json:"doi"`
		ImpactFactor float64          `json:"impact_factor"`
		PublishDate  string           `json:"publish_date"`
		Authors      []*entity.Author `json:"authors"`
		Projects     []uint           `json:"projects"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid save draft request",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取操作者信息
	operatorID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	// 解析日期
	var publishDate *time.Time
	if req.PublishDate != "" {
		if t, err := time.Parse("2006-01-02", req.PublishDate); err == nil {
			publishDate = &t
		}
	}

	if err := h.paperService.SaveDraft(
		uint(paperID),
		req.Title,
		req.Abstract,
		req.JournalID,
		req.DOI,
		req.ImpactFactor,
		publishDate,
		req.Authors,
		req.Projects,
		operatorID.(uint),
		ipAddress,
	); err != nil {
		if err == service.ErrPaperNotFound {
			response.Error(c, http.StatusNotFound, "论文不存在")
			return
		}
		if err == service.ErrInvalidStatus {
			response.Error(c, http.StatusBadRequest, "论文状态不允许保存草稿,仅草稿状态可保存")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// CheckDuplicate 检查重复
// POST /api/papers/check-duplicate
func (h *PaperHandler) CheckDuplicate(c *gin.Context) {
	var req struct {
		Title string `json:"title"`
		DOI   string `json:"doi"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid check duplicate request",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	papers, err := h.paperService.CheckDuplicate(req.Title, req.DOI)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, gin.H{
		"count":  len(papers),
		"papers": papers,
	})
}

// BatchImport 批量导入
// POST /api/papers/batch-import
func (h *PaperHandler) BatchImport(c *gin.Context) {
	// 获取上传的文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		logger.GetLogger().Warn("Failed to get file from request",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "文件上传失败")
		return
	}

	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		logger.GetLogger().Error("Failed to open uploaded file",
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "文件打开失败")
		return
	}
	defer file.Close()

	// 使用excelize读取Excel文件
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		logger.GetLogger().Error("Failed to open Excel file",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "Excel文件格式错误")
		return
	}
	defer xlsx.Close()

	// 获取操作者信息
	submitterID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	// 调用Service进行批量导入
	successCount, failedCount, errors := h.paperService.BatchImportPapers(xlsx, submitterID.(uint))

	response.Success(c, gin.H{
		"success_count": successCount,
		"failed_count":  failedCount,
		"errors":        errors,
	})
}

// GetMyPapers 获取我的论文
// GET /api/papers/my
func (h *PaperHandler) GetMyPapers(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	papers, total, err := h.paperService.GetMyPapers(userID.(uint), page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	paperList := make([]paperresponse.PaperDTO, len(papers))
	for i, paper := range papers {
		paperList[i] = *h.convertToPaperDTO(paper)
	}

	result := paperresponse.PaperListResponse{
		List:  paperList,
		Total: total,
		Page:  page,
		Size:  size,
	}

	response.Success(c, result)
}

// DownloadImportTemplate 下载导入模板
// GET /api/papers/import-template
func (h *PaperHandler) DownloadImportTemplate(c *gin.Context) {
	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()

	// 设置工作表名称
	sheetName := "论文导入"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		logger.GetLogger().Error("Failed to create sheet",
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "生成模板失败")
		return
	}
	f.SetActiveSheet(index)

	// 设置表头
	headers := []string{"论文标题", "摘要", "期刊ID", "DOI", "影响因子", "出版日期", "第一作者", "第一作者单位", "通讯作者", "通讯作者单位", "课题编号"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, header)
		// 设置表头样式
		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#E6E6FA"}, Pattern: 1},
		})
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// 添加示例数据
	exampleData := []interface{}{
		"示例论文标题",
		"这是论文摘要内容...",
		1,
		"10.1000/example",
		5.234,
		"2024-01-01",
		"张三",
		"北京大学",
		"李四",
		"清华大学",
		"NST-2024-001",
	}
	for i, data := range exampleData {
		cell := string(rune('A'+i)) + "2"
		f.SetCellValue(sheetName, cell, data)
	}

	// 设置列宽
	for i := range headers {
		col := string(rune('A' + i))
		f.SetColWidth(sheetName, col, col, 18)
	}

	// 生成Excel文件
	fileName := "论文导入模板.xlsx"
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	c.Header("Access-Control-Expose-Headers", "Content-Disposition")

	// 写入响应
	if err := f.Write(c.Writer); err != nil {
		logger.GetLogger().Error("Failed to write Excel file",
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "生成模板失败")
		return
	}
}

// convertToPaperDTO 将Paper实体转换为DTO格式
func (h *PaperHandler) convertToPaperDTO(paper *entity.Paper) *paperresponse.PaperDTO {
	paperDTO := &paperresponse.PaperDTO{
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
		paperDTO.Journal = &paperresponse.JournalDTO{
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
		paperDTO.Submitter = &paperresponse.UserDTO{
			ID:       paper.Submitter.ID,
			Username: paper.Submitter.Username,
			Name:     paper.Submitter.Name,
			Role:     paper.Submitter.Role,
		}
	}

	// 转换作者信息
	if paper.Authors != nil && len(paper.Authors) > 0 {
		authors := make([]paperresponse.AuthorDTO, len(paper.Authors))
		for i, author := range paper.Authors {
			authors[i] = paperresponse.AuthorDTO{
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
		projects := make([]paperresponse.ProjectDTO, len(paper.Projects))
		for i, project := range paper.Projects {
			projects[i] = paperresponse.ProjectDTO{
				ID:          project.ID,
				Name:        project.Name,
				Code:        project.Code,
				ProjectType: project.ProjectType,
				Source:      project.Source,
				Level:       project.Level,
			}
		}
		paperDTO.Projects = projects
	}

	// 转换附件信息
	if paper.Attachments != nil && len(paper.Attachments) > 0 {
		attachments := make([]paperresponse.AttachmentDTO, len(paper.Attachments))
		for i, attachment := range paper.Attachments {
			attachments[i] = paperresponse.AttachmentDTO{
				ID:        attachment.ID,
				FileType:  attachment.FileType,
				FileName:  attachment.FileName,
				FilePath:  attachment.FilePath,
				FileSize:  attachment.FileSize,
				MimeType:  attachment.MimeType,
				CreatedAt: attachment.CreatedAt,
			}
		}
		paperDTO.Attachments = attachments
	}

	return paperDTO
}
