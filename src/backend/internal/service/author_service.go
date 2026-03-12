package service

import (
	"errors"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/pkg/logger"

	"go.uber.org/zap"
)

var (
	// ErrAuthorNotFound 作者不存在
	ErrAuthorNotFound = errors.New("作者不存在")
	// ErrInvalidAuthorRank 无效的作者排名
	ErrInvalidAuthorRank = errors.New("无效的作者排名")
	// ErrDuplicateAuthorType 作者类型重复
	ErrDuplicateAuthorType = errors.New("作者类型重复")
)

// AuthorService 作者服务
type AuthorService struct {
	authorRepo *repository.AuthorRepository
}

// NewAuthorService 创建作者服务实例
func NewAuthorService(authorRepo *repository.AuthorRepository) *AuthorService {
	return &AuthorService{
		authorRepo: authorRepo,
	}
}

// CreateAuthor 创建作者
func (s *AuthorService) CreateAuthor(
	paperID uint,
	name string,
	authorType string,
	rank int,
	department string,
	userID *uint,
) (*entity.Author, error) {
	// 校验作者数据
	if err := s.ValidateAuthorData(authorType, rank); err != nil {
		return nil, err
	}

	author := &entity.Author{
		PaperID:    paperID,
		Name:       name,
		AuthorType: authorType,
		Rank:       rank,
		Department: department,
		UserID:     userID,
	}

	if err := s.authorRepo.Create(author); err != nil {
		logger.GetLogger().Error("Failed to create author",
			zap.Uint("paper_id", paperID),
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	logger.GetLogger().Info("Author created successfully",
		zap.Uint("author_id", author.ID),
		zap.Uint("paper_id", paperID),
		zap.String("name", name),
	)

	return author, nil
}

// UpdateAuthor 更新作者
func (s *AuthorService) UpdateAuthor(
	id uint,
	name string,
	authorType string,
	rank int,
	department string,
	userID *uint,
) error {
	// 查询作者
	author, err := s.authorRepo.FindByID(id)
	if err != nil {
		logger.GetLogger().Error("Failed to find author",
			zap.Uint("author_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if author == nil {
		return ErrAuthorNotFound
	}

	// 校验作者数据
	if err := s.ValidateAuthorData(authorType, rank); err != nil {
		return err
	}

	// 更新作者信息
	author.Name = name
	author.AuthorType = authorType
	author.Rank = rank
	author.Department = department
	author.UserID = userID

	if err := s.authorRepo.Update(author); err != nil {
		logger.GetLogger().Error("Failed to update author",
			zap.Uint("author_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	logger.GetLogger().Info("Author updated successfully",
		zap.Uint("author_id", id),
		zap.String("name", name),
	)

	return nil
}

// DeleteAuthor 删除作者
func (s *AuthorService) DeleteAuthor(id uint) error {
	// 查询作者是否存在
	author, err := s.authorRepo.FindByID(id)
	if err != nil {
		logger.GetLogger().Error("Failed to find author",
			zap.Uint("author_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if author == nil {
		return ErrAuthorNotFound
	}

	if err := s.authorRepo.Delete(id); err != nil {
		logger.GetLogger().Error("Failed to delete author",
			zap.Uint("author_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	logger.GetLogger().Info("Author deleted successfully",
		zap.Uint("author_id", id),
	)

	return nil
}

// GetAuthorsByPaperID 获取某论文的所有作者
func (s *AuthorService) GetAuthorsByPaperID(paperID uint) ([]*entity.Author, error) {
	authors, err := s.authorRepo.ListByPaperID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to get authors by paper id",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	return authors, nil
}

// BatchCreateAuthors 批量创建作者
func (s *AuthorService) BatchCreateAuthors(authors []*entity.Author) error {
	if len(authors) == 0 {
		return nil
	}

	// 校验作者数据
	for _, author := range authors {
		if err := s.ValidateAuthorData(author.AuthorType, author.Rank); err != nil {
			logger.GetLogger().Warn("Invalid author data",
				zap.String("name", author.Name),
				zap.Error(err),
			)
			return err
		}
	}

	// 校验作者类型互斥性（第一作者、共同第一作者、通讯作者不能重复）
	if err := s.validateAuthorTypes(authors); err != nil {
		return err
	}

	if err := s.authorRepo.CreateBatch(authors); err != nil {
		logger.GetLogger().Error("Failed to batch create authors",
			zap.Int("count", len(authors)),
			zap.Error(err),
		)
		return ErrSystemError
	}

	logger.GetLogger().Info("Authors batch created successfully",
		zap.Int("count", len(authors)),
	)

	return nil
}

// UpdateRankings 更新作者排名
func (s *AuthorService) UpdateRankings(paperID uint, authorIDs []uint) error {
	// 查询该论文的所有作者
	authors, err := s.authorRepo.ListByPaperID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to get authors by paper id",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	// 校验排名连续性
	if len(authorIDs) != len(authors) {
		return ErrInvalidAuthorRank
	}

	// 创建作者ID到作者的映射
	authorMap := make(map[uint]*entity.Author)
	for _, author := range authors {
		authorMap[author.ID] = author
	}

	// 更新作者排名
	for i, authorID := range authorIDs {
		author, exists := authorMap[authorID]
		if !exists {
			return ErrAuthorNotFound
		}
		author.Rank = i + 1
		if err := s.authorRepo.Update(author); err != nil {
			logger.GetLogger().Error("Failed to update author ranking",
				zap.Uint("author_id", author.ID),
				zap.Int("rank", i+1),
				zap.Error(err),
			)
			return ErrSystemError
		}
	}

	logger.GetLogger().Info("Author rankings updated successfully",
		zap.Uint("paper_id", paperID),
		zap.Int("count", len(authorIDs)),
	)

	return nil
}

// ValidateAuthorData 校验作者数据
func (s *AuthorService) ValidateAuthorData(authorType string, rank int) error {
	// 校验作者类型
	validTypes := map[string]bool{
		"first_author":         true,
		"co_first_author":      true,
		"corresponding_author": true,
		"author":               true,
	}

	if !validTypes[authorType] {
		return errors.New("无效的作者类型")
	}

	// 校验排名
	if rank < 1 {
		return ErrInvalidAuthorRank
	}

	return nil
}

// validateAuthorTypes 校验作者类型互斥性
func (s *AuthorService) validateAuthorTypes(authors []*entity.Author) error {
	firstAuthorCount := 0
	correspondingAuthorCount := 0

	for _, author := range authors {
		if author.AuthorType == "first_author" {
			firstAuthorCount++
		}
		if author.AuthorType == "corresponding_author" {
			correspondingAuthorCount++
		}
	}

	// 第一作者只能有一个
	if firstAuthorCount > 1 {
		return ErrDuplicateAuthorType
	}

	// 通讯作者只能有一个
	if correspondingAuthorCount > 1 {
		return ErrDuplicateAuthorType
	}

	return nil
}
