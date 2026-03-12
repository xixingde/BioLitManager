package main

import (
	"fmt"
	"os"
	"path/filepath"

	"biolitmanager/internal/cache"
	"biolitmanager/internal/config"
	"biolitmanager/internal/database"
	"biolitmanager/internal/handler"
	"biolitmanager/internal/middleware"
	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/internal/security"
	"biolitmanager/internal/service"
	"biolitmanager/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// 初始化日志（最先初始化，以便后续可以使用日志）
	if err := logger.InitLogger("debug"); err != nil {
		panic(fmt.Sprintf("Failed to init logger: %v", err))
	}
	logger.GetLogger().Info("Logger initialized")

	// 加载配置
	if err := config.InitConfig(); err != nil {
		logger.GetLogger().Fatal("Failed to load config", zap.Error(err))
	}

	cfg := config.GetConfig()
	logger.GetLogger().Info("Config loaded successfully",
		zap.String("port", cfg.Server.Port),
		zap.String("mode", cfg.Server.Mode),
	)

	// 确保数据库目录存在
	dataDir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		logger.GetLogger().Fatal("Failed to create data directory",
			zap.String("path", dataDir),
			zap.Error(err),
		)
	}
	logger.GetLogger().Info("Data directory created", zap.String("path", dataDir))

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		logger.GetLogger().Fatal("Failed to init database",
			zap.Error(err),
		)
	}
	database.SetDB(db)
	logger.GetLogger().Info("Database initialized successfully")

	// 执行数据库迁移
	if err := database.AutoMigrate(db); err != nil {
		logger.GetLogger().Fatal("Failed to migrate database",
			zap.Error(err),
		)
	}
	logger.GetLogger().Info("Database migration completed")

	// 初始化内存缓存
	memoryCache := cache.InitCache()
	logger.GetLogger().Info("Cache initialized")

	// 创建初始管理员账户
	if err := createInitialAdmin(db); err != nil {
		logger.GetLogger().Fatal("Failed to create initial admin",
			zap.Error(err),
		)
	}

	// 初始化 repositories
	userRepo := repository.NewUserRepository(db)
	operationLogRepo := repository.NewOperationLogRepository(db)
	authorRepo := repository.NewAuthorRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	journalRepo := repository.NewJournalRepository(db)
	paperRepo := repository.NewPaperRepository(db)
	paperProjectRepo := repository.NewPaperProjectRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	archiveRepo := repository.NewArchiveRepository(db)

	// 初始化 services
	operationLogService := service.NewOperationLogService(operationLogRepo)
	authService := service.NewAuthService(userRepo, memoryCache, operationLogService)
	userService := service.NewUserService(userRepo, operationLogService)
	authorService := service.NewAuthorService(authorRepo)
	projectService := service.NewProjectService(projectRepo, paperProjectRepo, operationLogService)
	journalService := service.NewJournalService(journalRepo, operationLogService)
	fileService := service.NewFileService(attachmentRepo)
	notificationService := service.NewNotificationService(userRepo, operationLogService)
	reviewService := service.NewReviewService(db, reviewRepo, paperRepo, operationLogService)
	archiveService := service.NewArchiveService(archiveRepo, paperRepo, operationLogService)
	paperService := service.NewPaperService(db, paperRepo, authorRepo, attachmentRepo, paperProjectRepo, operationLogService)

	// 初始化 handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	paperHandler := handler.NewPaperHandler(paperService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	projectHandler := handler.NewProjectHandler(projectService)
	journalHandler := handler.NewJournalHandler(journalService)
	fileHandler := handler.NewFileHandler(fileService)

	// 创建 Gin 实例
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// 注册中间件
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// 注册 CORS 中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 注册路由
	api := router.Group("/api")
	{
		// 认证相关路由（不需要认证）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// 需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// 认证路由
			authGroup := authenticated.Group("/auth")
			{
				authGroup.POST("/logout", authHandler.Logout)
				authGroup.GET("/profile", authHandler.GetProfile)
			}

			// 系统管理路由（需要权限）
			system := authenticated.Group("/system")
			system.Use(middleware.PermissionMiddleware(security.PermissionSystemUserManage))
			{
				users := system.Group("/users")
				{
					users.POST("", userHandler.CreateUser)
					users.PUT("/:id", userHandler.UpdateUser)
					users.PUT("/:id/disable", userHandler.DisableUser)
					users.PUT("/:id/enable", userHandler.EnableUser)
					users.GET("", userHandler.ListUsers)
				}
			}

			// 论文管理路由
			papers := authenticated.Group("/papers")
			{
				papers.POST("", paperHandler.CreatePaper)
				papers.GET("", paperHandler.ListPapers)
				papers.GET("/my", paperHandler.GetMyPapers)
				papers.GET("/:id", paperHandler.GetPaper)
				papers.PUT("/:id", paperHandler.UpdatePaper)
				papers.DELETE("/:id", paperHandler.DeletePaper)
				papers.POST("/:id/submit", paperHandler.SubmitForReview)
				papers.POST("/:id/save-draft", paperHandler.SaveDraft)
				papers.POST("/check-duplicate", paperHandler.CheckDuplicate)
				papers.POST("/batch-import", paperHandler.BatchImport)
			}

			// 审核管理路由
			reviews := authenticated.Group("/reviews")
			{
				reviews.POST("/business/:paperId", reviewHandler.BusinessReview)
				reviews.POST("/political/:paperId", reviewHandler.PoliticalReview)
				reviews.GET("/:paperId/logs", reviewHandler.GetReviewLogs)
				reviews.GET("/pending/business", reviewHandler.GetPendingBusinessReviews)
				reviews.GET("/pending/political", reviewHandler.GetPendingPoliticalReviews)
				reviews.GET("/my", reviewHandler.GetMyReviews)
			}

			// 课题管理路由
			projects := authenticated.Group("/projects")
			{
				projects.POST("", projectHandler.CreateProject)
				projects.GET("", projectHandler.ListProjects)
				projects.GET("/:id", projectHandler.GetProject)
				projects.PUT("/:id", projectHandler.UpdateProject)
				projects.DELETE("/:id", projectHandler.DeleteProject)
			}

			// 期刊管理路由
			journals := authenticated.Group("/journals")
			{
				journals.POST("", journalHandler.CreateJournal)
				journals.GET("", journalHandler.ListJournals)
				journals.GET("/search", journalHandler.SearchJournals)
				journals.GET("/:id", journalHandler.GetJournal)
				journals.PUT("/:id", journalHandler.UpdateJournal)
				journals.PUT("/:id/impact-factor", journalHandler.UpdateImpactFactor)
			}

			// 文件管理路由
			files := authenticated.Group("/files")
			{
				files.POST("/upload", fileHandler.UploadFile)
				files.GET("/:id", fileHandler.GetFile)
				files.GET("/:id/download", fileHandler.DownloadFile)
				files.DELETE("/:id", fileHandler.DeleteFile)
			}
		}
	}

	// 启动 HTTP 服务器
	addr := ":" + cfg.Server.Port
	logger.GetLogger().Info("Starting server",
		zap.String("addr", addr),
		zap.String("mode", cfg.Server.Mode),
	)

	if err := router.Run(addr); err != nil {
		logger.GetLogger().Fatal("Failed to start server",
			zap.Error(err),
		)
	}
}

// createInitialAdmin 创建初始管理员账户
func createInitialAdmin(db *gorm.DB) error {
	userRepo := repository.NewUserRepository(db)

	// 检查是否已存在管理员账户
	admin, err := userRepo.FindByUsername("admin")
	if err != nil {
		return fmt.Errorf("failed to check admin existence: %w", err)
	}

	// 如果已存在，则跳过
	if admin != nil {
		logger.GetLogger().Info("Admin account already exists")
		return nil
	}

	// 创建初始管理员账户
	passwordHash, err := security.HashPassword("Admin@123")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		Username:     "admin",
		PasswordHash: passwordHash,
		Name:         "超级管理员",
		Role:         string(security.RoleSuperAdmin),
		Department:   "系统管理",
	}

	if err := userRepo.Create(user); err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	logger.GetLogger().Info("Initial admin account created successfully",
		zap.String("username", "admin"),
		zap.String("role", user.Role),
	)

	return nil
}
