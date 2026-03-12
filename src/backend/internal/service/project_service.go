package service

import (
	"errors"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/pkg/logger"

	"go.uber.org/zap"
)

var (
	// ErrProjectNotFound 课题不存在
	ErrProjectNotFound = errors.New("课题不存在")
	// ErrProjectCodeExists 课题编号已存在
	ErrProjectCodeExists = errors.New("课题编号已存在")
	// ErrProjectLinked 课题已关联论文，无法删除
	ErrProjectLinked = errors.New("课题已关联论文，无法删除")
)

// ProjectService 课题服务
type ProjectService struct {
	projectRepo         *repository.ProjectRepository
	paperProjectRepo    *repository.PaperProjectRepository
	operationLogService *OperationLogService
}

// NewProjectService 创建课题服务实例
func NewProjectService(
	projectRepo *repository.ProjectRepository,
	paperProjectRepo *repository.PaperProjectRepository,
	operationLogService *OperationLogService,
) *ProjectService {
	return &ProjectService{
		projectRepo:         projectRepo,
		paperProjectRepo:    paperProjectRepo,
		operationLogService: operationLogService,
	}
}

// CreateProject 创建课题
func (s *ProjectService) CreateProject(
	name string,
	code string,
	projectType string,
	source string,
	level string,
	status string,
	operatorID uint,
	ipAddress string,
) (*entity.Project, error) {
	// 校验课题编号唯一性
	existingProject, err := s.projectRepo.FindByCode(code)
	if err != nil {
		logger.GetLogger().Error("Failed to check project code",
			zap.String("code", code),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	if existingProject != nil {
		return nil, ErrProjectCodeExists
	}

	project := &entity.Project{
		Name:        name,
		Code:        code,
		ProjectType: projectType,
		Source:      source,
		Level:       level,
		Status:      status,
	}

	if err := s.projectRepo.Create(project); err != nil {
		logger.GetLogger().Error("Failed to create project",
			zap.String("name", name),
			zap.String("code", code),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "create", "project", code, "创建课题 "+name, "成功", ipAddress)

	logger.GetLogger().Info("Project created successfully",
		zap.Uint("project_id", project.ID),
		zap.String("name", name),
		zap.String("code", code),
	)

	return project, nil
}

// UpdateProject 更新课题
func (s *ProjectService) UpdateProject(
	id uint,
	name string,
	code string,
	projectType string,
	source string,
	level string,
	status string,
	operatorID uint,
	ipAddress string,
) error {
	// 查询课题
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		logger.GetLogger().Error("Failed to find project",
			zap.Uint("project_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if project == nil {
		return ErrProjectNotFound
	}

	// 校验课题编号唯一性（排除当前课题）
	existingProject, err := s.projectRepo.FindByCode(code)
	if err != nil {
		logger.GetLogger().Error("Failed to check project code",
			zap.String("code", code),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if existingProject != nil && existingProject.ID != id {
		return ErrProjectCodeExists
	}

	// 更新课题信息
	project.Name = name
	project.Code = code
	project.ProjectType = projectType
	project.Source = source
	project.Level = level
	project.Status = status

	if err := s.projectRepo.Update(project); err != nil {
		logger.GetLogger().Error("Failed to update project",
			zap.Uint("project_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "update", "project", code, "更新课题 "+name, "成功", ipAddress)

	logger.GetLogger().Info("Project updated successfully",
		zap.Uint("project_id", id),
		zap.String("name", name),
	)

	return nil
}

// DeleteProject 删除课题
func (s *ProjectService) DeleteProject(id uint, operatorID uint, ipAddress string) error {
	// 查询课题
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		logger.GetLogger().Error("Failed to find project",
			zap.Uint("project_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if project == nil {
		return ErrProjectNotFound
	}

	// 检查课题是否已关联论文
	linkCount, err := s.projectRepo.CheckIsLinked(id)
	if err != nil {
		logger.GetLogger().Error("Failed to check project links",
			zap.Uint("project_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if linkCount > 0 {
		return ErrProjectLinked
	}

	if err := s.projectRepo.Delete(id); err != nil {
		logger.GetLogger().Error("Failed to delete project",
			zap.Uint("project_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "delete", "project", project.Code, "删除课题 "+project.Name, "成功", ipAddress)

	logger.GetLogger().Info("Project deleted successfully",
		zap.Uint("project_id", id),
		zap.String("name", project.Name),
	)

	return nil
}

// GetProjectByID 获取课题详情
func (s *ProjectService) GetProjectByID(id uint) (*entity.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		logger.GetLogger().Error("Failed to find project",
			zap.Uint("project_id", id),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	if project == nil {
		return nil, ErrProjectNotFound
	}

	return project, nil
}

// ListProjects 分页查询课题列表
func (s *ProjectService) ListProjects(page, size int, name, code, projectType, level string) ([]*entity.Project, int64, error) {
	// 这里暂时使用repository的List方法，后续可以在repository层添加更多查询条件
	projects, total, err := s.projectRepo.List(page, size)
	if err != nil {
		logger.GetLogger().Error("Failed to list projects",
			zap.Int("page", page),
			zap.Int("size", size),
			zap.Error(err),
		)
		return nil, 0, ErrSystemError
	}

	// TODO: 在repository层添加更灵活的查询条件（按名称、编号、类型、级别筛选）

	return projects, total, nil
}
