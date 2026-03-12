package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/pkg/logger"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// ErrExportPermissionDenied 导出权限不足
	ErrExportPermissionDenied = fmt.Errorf("导出权限不足")
	// ErrExportFileFailed 导出文件失败
	ErrExportFileFailed = fmt.Errorf("导出文件失败")
	// ErrExportDirNotExist 导出目录不存在
	ErrExportDirNotExist = fmt.Errorf("导出目录不存在")
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = fmt.Errorf("用户不存在")
)

// ExportFieldMap 导出字段映射
var ExportFieldMap = map[string]string{
	"title":          "标题",
	"authors":        "作者",
	"journal_name":   "期刊名称",
	"publish_date":   "发表日期",
	"partition":      "收录类型",
	"impact_factor":  "影响因子",
	"citation_count": "引用次数",
	"doi":            "DOI",
	"abstract":       "摘要",
	"volume":         "卷",
	"issue":          "期",
	"start_page":     "起始页",
	"end_page":       "结束页",
}

// ExportServiceInterface 导出服务接口
type ExportServiceInterface interface {
	ExportPapersToExcel(papers []*entity.Paper, fields []string) (string, error)
	ExportPaperToPDF(paperID uint, userID uint) (string, error)
	ExportPaperToWord(paperID uint, userID uint) (string, error)
	ExportStatsToExcel(stats interface{}, title string) (string, error)
	ExportStatsToPDF(stats interface{}, title string) (string, error)
}

// ExportService 导出服务
type ExportService struct {
	db          *gorm.DB
	paperRepo   *repository.PaperRepository
	authorRepo  *repository.AuthorRepository
	projectRepo *repository.ProjectRepository
	reviewRepo  *repository.ReviewRepository
	userRepo    *repository.UserRepository
	exportDir   string
}

// NewExportService 创建导出服务实例
func NewExportService(
	db *gorm.DB,
	paperRepo *repository.PaperRepository,
	authorRepo *repository.AuthorRepository,
	projectRepo *repository.ProjectRepository,
	reviewRepo *repository.ReviewRepository,
	userRepo *repository.UserRepository,
) *ExportService {
	// 默认导出目录
	exportDir := "uploads/exports"
	return &ExportService{
		db:          db,
		paperRepo:   paperRepo,
		authorRepo:  authorRepo,
		projectRepo: projectRepo,
		reviewRepo:  reviewRepo,
		userRepo:    userRepo,
		exportDir:   exportDir,
	}
}

// ExportPapersToExcel 批量导出论文到Excel
func (s *ExportService) ExportPapersToExcel(papers []*entity.Paper, fields []string) (string, error) {
	// 确保导出目录存在
	if err := s.ensureExportDir(); err != nil {
		return "", err
	}

	// 创建Excel文件
	f := excelize.NewFile()
	sheetName := "论文数据"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		logger.GetLogger().Error("Failed to create sheet",
			zap.Error(err),
		)
		return "", ErrExportFileFailed
	}

	// 生成表头
	s.generateExcelHeader(f, sheetName, fields)

	// 填充数据
	for i, paper := range papers {
		row := i + 2 // 从第2行开始（跳过表头）
		for j, field := range fields {
			col, _ := excelize.ColumnNumberToName(j + 1)
			cell := fmt.Sprintf("%s%d", col, row)

			value := s.getPaperFieldValue(paper, field)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("papers_%s.xlsx", timestamp)
	filePath := filepath.Join(s.exportDir, fileName)

	// 保存文件
	if err := f.SaveAs(filePath); err != nil {
		logger.GetLogger().Error("Failed to save Excel file",
			zap.String("file_path", filePath),
			zap.Error(err),
		)
		return "", ErrExportFileFailed
	}

	logger.GetLogger().Info("Papers exported to Excel successfully",
		zap.String("file_path", filePath),
		zap.Int("count", len(papers)),
	)

	return filePath, nil
}

// ExportPaperToPDF 导出单篇论文到PDF
func (s *ExportService) ExportPaperToPDF(paperID uint, userID uint) (string, error) {
	// 权限校验
	if err := s.checkExportPermission(userID, paperID); err != nil {
		return "", err
	}

	// 查询论文详情
	paper, err := s.getPaperWithRelations(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to get paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return "", ErrPaperNotFound
	}

	if paper == nil {
		return "", ErrPaperNotFound
	}

	// 确保导出目录存在
	if err := s.ensureExportDir(); err != nil {
		return "", err
	}

	// 创建PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// 生成PDF内容
	s.generatePDFContent(pdf, paper)

	// 生成文件名
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("paper_%d_%s.pdf", paperID, timestamp)
	filePath := filepath.Join(s.exportDir, fileName)

	// 保存文件
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		logger.GetLogger().Error("Failed to save PDF file",
			zap.String("file_path", filePath),
			zap.Error(err),
		)
		return "", ErrExportFileFailed
	}

	logger.GetLogger().Info("Paper exported to PDF successfully",
		zap.String("file_path", filePath),
		zap.Uint("paper_id", paperID),
	)

	return filePath, nil
}

// ExportPaperToWord 导出单篇论文到Word（使用gofpdf模拟，输出PDF）
func (s *ExportService) ExportPaperToWord(paperID uint, userID uint) (string, error) {
	// 权限校验
	if err := s.checkExportPermission(userID, paperID); err != nil {
		return "", err
	}

	// 查询论文详情
	paper, err := s.getPaperWithRelations(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to get paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return "", ErrPaperNotFound
	}

	if paper == nil {
		return "", ErrPaperNotFound
	}

	// 确保导出目录存在
	if err := s.ensureExportDir(); err != nil {
		return "", err
	}

	// 创建PDF（模拟Word格式）
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// 生成PDF内容（类似Word的格式）
	s.generateWordLikePDFContent(pdf, paper)

	// 生成文件名（使用.doc扩展名作为标识，但实际输出PDF）
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("paper_%d_%s.doc", paperID, timestamp)
	filePath := filepath.Join(s.exportDir, fileName)

	// 保存文件（实际是PDF格式）
	pdfFilePath := strings.TrimSuffix(filePath, ".doc") + ".pdf"
	if err := pdf.OutputFileAndClose(pdfFilePath); err != nil {
		logger.GetLogger().Error("Failed to save Word file",
			zap.String("file_path", pdfFilePath),
			zap.Error(err),
		)
		return "", ErrExportFileFailed
	}

	logger.GetLogger().Info("Paper exported to Word successfully",
		zap.String("file_path", pdfFilePath),
		zap.Uint("paper_id", paperID),
	)

	return pdfFilePath, nil
}

// ExportStatsToExcel 导出统计结果到Excel
func (s *ExportService) ExportStatsToExcel(stats interface{}, title string) (string, error) {
	// 确保导出目录存在
	if err := s.ensureExportDir(); err != nil {
		return "", err
	}

	// 创建Excel文件
	f := excelize.NewFile()
	sheetName := "统计结果"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		logger.GetLogger().Error("Failed to create sheet",
			zap.Error(err),
		)
		return "", ErrExportFileFailed
	}

	// 设置标题
	f.SetCellValue(sheetName, "A1", title)
	titleStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", "A1", titleStyle)
	}

	// 根据stats类型处理数据
	row := 2
	switch v := stats.(type) {
	case map[string]interface{}:
		f.SetCellValue(sheetName, "A2", "统计项")
		f.SetCellValue(sheetName, "B2", "数值")
		for k, val := range v {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), k)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), val)
			row++
		}
	case []map[string]interface{}:
		if len(v) > 0 {
			// 设置表头
			headers := make([]string, 0)
			for k := range v[0] {
				headers = append(headers, k)
			}
			for j, header := range headers {
				col, _ := excelize.ColumnNumberToName(j + 1)
				f.SetCellValue(sheetName, fmt.Sprintf("%s2", col), header)
			}
			// 填充数据
			for i, item := range v {
				row := i + 3
				for j, header := range headers {
					col, _ := excelize.ColumnNumberToName(j + 1)
					f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, row), item[header])
				}
			}
		}
	default:
		f.SetCellValue(sheetName, "A2", "统计结果")
		f.SetCellValue(sheetName, "B2", fmt.Sprintf("%v", v))
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("stats_%s.xlsx", timestamp)
	filePath := filepath.Join(s.exportDir, fileName)

	// 保存文件
	if err := f.SaveAs(filePath); err != nil {
		logger.GetLogger().Error("Failed to save Excel file",
			zap.String("file_path", filePath),
			zap.Error(err),
		)
		return "", ErrExportFileFailed
	}

	logger.GetLogger().Info("Stats exported to Excel successfully",
		zap.String("file_path", filePath),
	)

	return filePath, nil
}

// ExportStatsToPDF 导出统计结果到PDF
func (s *ExportService) ExportStatsToPDF(stats interface{}, title string) (string, error) {
	// 确保导出目录存在
	if err := s.ensureExportDir(); err != nil {
		return "", err
	}

	// 创建PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// 设置标题
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, title)
	pdf.Ln(15)

	// 根据stats类型处理数据
	pdf.SetFont("Arial", "", 12)
	switch v := stats.(type) {
	case map[string]interface{}:
		for k, val := range v {
			pdf.Cell(190, 8, fmt.Sprintf("%s: %v", k, val))
			pdf.Ln(8)
		}
	case []map[string]interface{}:
		if len(v) > 0 {
			// 打印表头
			headers := make([]string, 0)
			for k := range v[0] {
				headers = append(headers, k)
			}
			for _, header := range headers {
				pdf.Cell(float64(190)/float64(len(headers)), 8, header)
			}
			pdf.Ln(8)
			// 打印数据
			for _, item := range v {
				for _, header := range headers {
					pdf.Cell(float64(190)/float64(len(headers)), 8, fmt.Sprintf("%v", item[header]))
				}
				pdf.Ln(8)
			}
		}
	default:
		pdf.Cell(190, 8, fmt.Sprintf("统计结果: %v", v))
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("stats_%s.pdf", timestamp)
	filePath := filepath.Join(s.exportDir, fileName)

	// 保存文件
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		logger.GetLogger().Error("Failed to save PDF file",
			zap.String("file_path", filePath),
			zap.Error(err),
		)
		return "", ErrExportFileFailed
	}

	logger.GetLogger().Info("Stats exported to PDF successfully",
		zap.String("file_path", filePath),
	)

	return filePath, nil
}

// generateExcelHeader 生成Excel表头
func (s *ExportService) generateExcelHeader(f *excelize.File, sheet string, fields []string) {
	for i, field := range fields {
		col, _ := excelize.ColumnNumberToName(i + 1)
		cell := fmt.Sprintf("%s1", col)

		// 使用中文字段名
		fieldName := ExportFieldMap[field]
		if fieldName == "" {
			fieldName = field
		}

		f.SetCellValue(sheet, cell, fieldName)
	}
}

// generatePDFContent 生成PDF内容
func (s *ExportService) generatePDFContent(pdf *gofpdf.Fpdf, paper *entity.Paper) {
	// 标题
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Paper Details")
	pdf.Ln(15)

	// 论文标题
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Title:")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(190, 6, paper.Title, "", "", false)
	pdf.Ln(5)

	// 作者
	if len(paper.Authors) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "Authors:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 11)
		authors := make([]string, len(paper.Authors))
		for i, author := range paper.Authors {
			authors[i] = author.Name
		}
		pdf.MultiCell(190, 6, strings.Join(authors, ", "), "", "", false)
		pdf.Ln(5)
	}

	// 期刊名称
	if paper.JournalName != "" {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "Journal:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(190, 6, paper.JournalName)
		pdf.Ln(10)
	}

	// 发表日期
	if paper.PublishDate != nil {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "Publish Date:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(190, 6, paper.PublishDate.Format("2006-01-02"))
		pdf.Ln(10)
	}

	// DOI
	if paper.DOI != "" {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "DOI:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(190, 6, paper.DOI)
		pdf.Ln(10)
	}

	// 影响因子
	if paper.ImpactFactor > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "Impact Factor:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(190, 6, fmt.Sprintf("%.2f", paper.ImpactFactor))
		pdf.Ln(10)
	}

	// 引用次数
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Citation Count:")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(190, 6, fmt.Sprintf("%d", paper.CitationCount))
	pdf.Ln(10)

	// 收录类型
	if paper.Partition != "" {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "Partition:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(190, 6, paper.Partition)
		pdf.Ln(10)
	}

	// 卷期页
	if paper.Volume != "" || paper.Issue != "" {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "Volume/Issue/Pages:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 11)
		volumeInfo := fmt.Sprintf("Vol: %s, Issue: %s, Pages: %s-%s",
			paper.Volume, paper.Issue, paper.StartPage, paper.EndPage)
		pdf.Cell(190, 6, volumeInfo)
		pdf.Ln(10)
	}

	// 摘要
	if paper.Abstract != "" {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "Abstract:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(190, 5, paper.Abstract, "", "", false)
	}
}

// generateWordLikePDFContent 生成类似Word格式的PDF内容
func (s *ExportService) generateWordLikePDFContent(pdf *gofpdf.Fpdf, paper *entity.Paper) {
	// 标题
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(190, 12, paper.Title)
	pdf.Ln(15)

	// 作者信息
	if len(paper.Authors) > 0 {
		pdf.SetFont("Arial", "I", 11)
		authors := make([]string, len(paper.Authors))
		for i, author := range paper.Authors {
			authors[i] = author.Name
		}
		pdf.MultiCell(190, 6, strings.Join(authors, ", "), "", "C", false)
		pdf.Ln(10)
	}

	// 期刊和日期信息
	pdf.SetFont("Arial", "", 11)
	if paper.JournalName != "" {
		pdf.Cell(190, 6, fmt.Sprintf("Journal: %s", paper.JournalName))
		pdf.Ln(6)
	}
	if paper.PublishDate != nil {
		pdf.Cell(190, 6, fmt.Sprintf("Published: %s", paper.PublishDate.Format("2006-01-02")))
		pdf.Ln(6)
	}
	if paper.DOI != "" {
		pdf.Cell(190, 6, fmt.Sprintf("DOI: %s", paper.DOI))
		pdf.Ln(6)
	}
	pdf.Ln(10)

	// 分隔线
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(5)

	// 详细信息
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Details:")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	if paper.ImpactFactor > 0 {
		pdf.Cell(190, 6, fmt.Sprintf("Impact Factor: %.2f", paper.ImpactFactor))
		pdf.Ln(6)
	}
	pdf.Cell(190, 6, fmt.Sprintf("Citations: %d", paper.CitationCount))
	pdf.Ln(6)
	if paper.Partition != "" {
		pdf.Cell(190, 6, fmt.Sprintf("Partition: %s", paper.Partition))
		pdf.Ln(6)
	}
	if paper.Volume != "" {
		pdf.Cell(190, 6, fmt.Sprintf("Volume: %s", paper.Volume))
		pdf.Ln(6)
	}
	if paper.Issue != "" {
		pdf.Cell(190, 6, fmt.Sprintf("Issue: %s", paper.Issue))
		pdf.Ln(6)
	}
	if paper.StartPage != "" && paper.EndPage != "" {
		pdf.Cell(190, 6, fmt.Sprintf("Pages: %s-%s", paper.StartPage, paper.EndPage))
		pdf.Ln(6)
	}
	pdf.Ln(10)

	// 摘要
	if paper.Abstract != "" {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "Abstract:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(190, 5, paper.Abstract, "", "", false)
	}
}

// checkExportPermission 导出权限校验
func (s *ExportService) checkExportPermission(userID uint, paperID uint) error {
	// 获取用户信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.GetLogger().Error("Failed to get user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if user == nil {
		return ErrUserNotFound
	}

	// 管理员可以导出所有数据
	if user.Role == "admin" {
		return nil
	}

	// 查询论文
	paper, err := s.paperRepo.FindByID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to get paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if paper == nil {
		return ErrPaperNotFound
	}

	// 普通用户只能导出自己提交的论文
	if paper.SubmitterID != userID {
		logger.GetLogger().Warn("Export permission denied",
			zap.Uint("user_id", userID),
			zap.Uint("paper_id", paperID),
			zap.Uint("submitter_id", paper.SubmitterID),
		)
		return ErrExportPermissionDenied
	}

	return nil
}

// filterAccessiblePapers 数据范围过滤
func (s *ExportService) filterAccessiblePapers(userID uint, papers []*entity.Paper) []*entity.Paper {
	// 获取用户信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.GetLogger().Error("Failed to get user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return []*entity.Paper{}
	}

	if user == nil {
		return []*entity.Paper{}
	}

	// 管理员可以查看所有数据
	if user.Role == "admin" {
		return papers
	}

	// 普通用户只能看到自己提交的论文
	filtered := make([]*entity.Paper, 0)
	for _, paper := range papers {
		if paper.SubmitterID == userID {
			filtered = append(filtered, paper)
		}
	}

	return filtered
}

// getPaperFieldValue 获取论文字段值
func (s *ExportService) getPaperFieldValue(paper *entity.Paper, field string) string {
	switch field {
	case "title":
		return paper.Title
	case "authors":
		if len(paper.Authors) > 0 {
			names := make([]string, len(paper.Authors))
			for i, author := range paper.Authors {
				names[i] = author.Name
			}
			return strings.Join(names, ", ")
		}
		return ""
	case "journal_name":
		return paper.JournalName
	case "publish_date":
		if paper.PublishDate != nil {
			return paper.PublishDate.Format("2006-01-02")
		}
		return ""
	case "partition":
		return paper.Partition
	case "impact_factor":
		return fmt.Sprintf("%.2f", paper.ImpactFactor)
	case "citation_count":
		return fmt.Sprintf("%d", paper.CitationCount)
	case "doi":
		return paper.DOI
	case "abstract":
		return paper.Abstract
	case "volume":
		return paper.Volume
	case "issue":
		return paper.Issue
	case "start_page":
		return paper.StartPage
	case "end_page":
		return paper.EndPage
	default:
		return ""
	}
}

// getPaperWithRelations 获取论文及其关联数据
func (s *ExportService) getPaperWithRelations(paperID uint) (*entity.Paper, error) {
	var paper entity.Paper
	err := s.db.Preload("Journal").
		Preload("Submitter").
		Preload("Authors").
		Preload("Projects").
		Preload("Attachments").
		First(&paper, paperID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &paper, nil
}

// ensureExportDir 确保导出目录存在
func (s *ExportService) ensureExportDir() error {
	if _, err := os.Stat(s.exportDir); os.IsNotExist(err) {
		if err := os.MkdirAll(s.exportDir, 0755); err != nil {
			logger.GetLogger().Error("Failed to create export directory",
				zap.String("dir", s.exportDir),
				zap.Error(err),
			)
			return ErrExportDirNotExist
		}
	}
	return nil
}

// getTitleStyle 获取标题样式
func (s *ExportService) getTitleStyle(f *excelize.File) (int, error) {
	return f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
}
