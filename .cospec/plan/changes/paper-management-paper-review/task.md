## 实施

### 1. 后端基础设施准备

#### 1.1 新增数据库实体模型
- [x] 1.1.1 创建论文实体模型
     【目标对象】`src/backend/internal/model/entity/paper.go`
     【修改目的】定义论文数据表结构
     【修改方式】创建新文件
     【相关依赖】GORM框架
     【修改内容】
        - 定义 Paper 结构体，包含以下字段：ID、Title（标题）、Abstract（摘要）、JournalID（期刊ID）、DOI（数字对象唯一标识符）、ImpactFactor（影响因子）、PublishDate（出版日期）、Status（状态：草稿/待业务审核/待政工审核/审核通过/驳回）、SubmitterID（提交人ID）、SubmitTime（提交时间）、CreatedAt、UpdatedAt、DeletedAt
        - 添加 GORM 标签定义表名和字段约束（标题必填、状态默认为草稿）
        - 添加外键关联 Journal 实体和 User 实体

- [x] 1.1.2 创建作者实体模型
     【目标对象】`src/backend/internal/model/entity/author.go`
     【修改目的】定义作者数据表结构
     【修改方式】创建新文件
     【相关依赖】GORM框架
     【修改内容】
        - 定义 Author 结构体，包含以下字段：ID、PaperID（论文ID）、Name（姓名）、AuthorType（作者类型：第一作者/共同第一作者/通讯作者/普通作者）、Rank（排名）、Department（单位）、UserID（关联人员库用户ID，可选）、CreatedAt、UpdatedAt、DeletedAt
        - 添加 GORM 标签定义表名和字段约束
        - 添加外键关联 Paper 实体和 User 实体（可选）

- [x] 1.1.3 创建课题实体模型
     【目标对象】`src/backend/internal/model/entity/project.go`
     【修改目的】定义课题数据表结构
     【修改方式】创建新文件
     【相关依赖】GORM框架
     【修改内容】
        - 定义 Project 结构体，包含以下字段：ID、Name（课题名称）、Code（课题编号）、ProjectType（项目类型：纵向/横向）、Source（来源）、Level（级别：国家级/省部级/市级）、Status（状态：进行中/已结题）、CreatedAt、UpdatedAt、DeletedAt
        - 添加 GORM 标签定义表名和字段约束（课题名称必填、课题编号唯一）
        - 添加唯一索引到 Code 字段

- [x] 1.1.4 创建期刊实体模型
     【目标对象】`src/backend/internal/model/entity/journal.go`
     【修改目的】定义期刊数据表结构
     【修改方式】创建新文件
     【相关依赖】GORM框架
     【修改内容】
        - 定义 Journal 结构体，包含以下字段：ID、FullName（期刊全称）、ShortName（期刊简称）、ISSN（国际标准连续出版物编号）、ImpactFactor（影响因子）、Publisher（出版社）、CreatedAt、UpdatedAt、DeletedAt
        - 添加 GORM 标签定义表名和字段约束（ISSN唯一）
        - 添加唯一索引到 ISSN 字段

- [x] 1.1.5 创建论文课题关联实体模型
     【目标对象】`src/backend/internal/model/entity/paper_project.go`
     【修改目的】定义论文与课题的多对多关联表结构
     【修改方式】创建新文件
     【相关依赖】GORM框架
     【修改内容】
        - 定义 PaperProject 结构体，包含以下字段：ID、PaperID（论文ID）、ProjectID（课题ID）、CreatedAt
        - 添加 GORM 标签定义表名和字段约束
        - 添加联合唯一索引到 PaperID 和 ProjectID 字段
        - 添加外键关联 Paper 实体和 Project 实体

- [x] 1.1.6 创建审核记录实体模型
     【目标对象】`src/backend/internal/model/entity/review_log.go`
     【修改目的】定义审核记录数据表结构
     【修改方式】创建新文件
     【相关依赖】GORM框架
     【修改内容】
        - 定义 ReviewLog 结构体，包含以下字段：ID、PaperID（论文ID）、ReviewType（审核类型：业务审核/政工审核）、Result（审核结果：通过/驳回）、Comment（审核意见）、ReviewerID（审核人ID）、ReviewTime（审核时间）、CreatedAt、UpdatedAt
        - 添加 GORM 标签定义表名和字段约束
        - 添加外键关联 Paper 实体和 User 实体

- [x] 1.1.7 创建归档实体模型
     【目标对象】`src/backend/internal/model/entity/archive.go`
     【修改目的】定义归档记录数据表结构
     【修改方式】创建新文件
     【相关依赖】GORM框架
     【修改内容】
        - 定义 Archive 结构体，包含以下字段：ID、PaperID（论文ID）、ArchiveNumber（归档编号）、ArchiveDate（归档时间）、ArchiverID（归档人ID）、CreatedAt
        - 添加 GORM 标签定义表名和字段约束
        - 添加外键关联 Paper 实体和 User 实体
        - 添加唯一索引到 ArchiveNumber 字段

- [x] 1.1.8 创建附件实体模型
     【目标对象】`src/backend/internal/model/entity/attachment.go`
     【修改目的】定义附件数据表结构
     【修改方式】创建新文件
     【相关依赖】GORM框架
     【修改内容】
        - 定义 Attachment 结构体，包含以下字段：ID、PaperID（论文ID）、FileType（文件类型：全文/首页/期刊封面/审批件）、FileName（原始文件名）、FilePath（存储路径）、FileSize（文件大小字节）、MimeType（MIME类型）、UploaderID（上传人ID）、CreatedAt
        - 添加 GORM 标签定义表名和字段约束
        - 添加外键关联 Paper 实体和 User 实体

#### 1.2 创建请求和响应DTO
- [x] 1.2.1 创建论文请求DTO
     【目标对象】`src/backend/internal/model/dto/request/paper_request.go`
     【修改目的】定义论文相关API请求数据结构
     【修改方式】创建新文件
     【相关依赖】无
     【修改内容】
        - 定义 CreatePaperRequest 结构体：Title、Abstract、JournalID、DOI、ImpactFactor、PublishDate、Authors（作者列表）、Projects（课题ID列表）、Attachments（附件信息）
        - 定义 UpdatePaperRequest 结构体：Title、Abstract、JournalID、DOI、ImpactFactor、PublishDate、Authors、Projects
        - 定义 SubmitForReviewRequest 结构体：PaperID
        - 定义 SaveDraftRequest 结构体：PaperID
        - 定义 CheckDuplicateRequest 结构体：Title、DOI
        - 定义 BatchImportRequest 结构体：File（Excel文件）

- [x] 1.2.2 创建审核请求DTO
     【目标对象】`src/backend/internal/model/dto/request/review_request.go`
     【修改目的】定义审核相关API请求数据结构
     【修改方式】创建新文件
     【相关依赖】无
     【修改内容】
        - 定义 BusinessReviewRequest 结构体：Result（通过/驳回）、Comment（审核意见）
        - 定义 PoliticalReviewRequest 结构体：Result（通过/驳回）、Comment（审核意见）

- [x] 1.2.3 创建作者请求DTO
     【目标对象】`src/backend/internal/model/dto/request/author_request.go`
     【修改目的】定义作者相关API请求数据结构
     【修改方式】创建新文件
     【相关依赖】无
     【修改内容】
        - 定义 CreateAuthorRequest 结构体：PaperID、Name、AuthorType、Rank、Department、UserID
        - 定义 UpdateAuthorRequest 结构体：Name、AuthorType、Rank、Department、UserID

- [x] 1.2.4 创建课题请求DTO
     【目标对象】`src/backend/internal/model/dto/request/project_request.go`
     【修改目的】定义课题相关API请求数据结构
     【修改方式】创建新文件
     【相关依赖】无
     【修改内容】
        - 定义 CreateProjectRequest 结构体：Name、Code、ProjectType、Source、Level
        - 定义 UpdateProjectRequest 结构体：Name、Code、ProjectType、Source、Level、Status

- [x] 1.2.5 创建期刊请求DTO
     【目标对象】`src/backend/internal/model/dto/request/journal_request.go`
     【修改目的】定义期刊相关API请求数据结构
     【修改方式】创建新文件
     【相关依赖】无
     【修改内容】
        - 定义 CreateJournalRequest 结构体：FullName、ShortName、ISSN、ImpactFactor、Publisher
        - 定义 UpdateJournalRequest 结构体：FullName、ShortName、ImpactFactor、Publisher

- [x] 1.2.6 创建文件上传请求DTO
     【目标对象】`src/backend/internal/model/dto/request/file_upload_request.go`
     【修改目的】定义文件上传API请求数据结构
     【修改方式】创建新文件
     【相关依赖】无
     【修改内容】
        - 定义 UploadFileRequest 结构体：PaperID、FileType、File（文件对象）

- [x] 1.2.7 创建论文响应DTO
     【目标对象】`src/backend/internal/model/dto/response/paper_response.go`
     【修改目的】定义论文相关API响应数据结构
     【修改方式】创建新文件
     【相关依赖】无
     【修改内容】
        - 定义 PaperDTO 结构体：ID、Title、Abstract、Journal（期刊信息）、DOI、ImpactFactor、PublishDate、Status、Submitter（提交人信息）、SubmitTime、Authors（作者列表）、Projects（课题列表）、Attachments（附件列表）、CreatedAt、UpdatedAt
        - 定义 PaperListResponse 结构体：List（PaperDTO数组）、Total、Page、Size

- [x] 1.2.8 创建审核响应DTO
     【目标对象】`src/backend/internal/model/dto/response/review_response.go`
     【修改目的】定义审核相关API响应数据结构
     【修改方式】创建新文件
     【相关依赖】无
     【修改内容】
        - 定义 ReviewLogDTO 结构体：ID、PaperID、ReviewType、Result、Comment、Reviewer（审核人信息）、ReviewTime、CreatedAt
        - 定义 PendingReviewDTO 结构体：PaperID、Title、SubmitterName、SubmitTime、Status、DaysSinceSubmit（距离提交天数）

#### 1.3 新增错误码定义
- [x] 1.3.1 在 errors.go 中新增论文相关错误码
     【目标对象】`src/backend/pkg/errors/errors.go`
     【修改目的】定义论文管理相关错误码
     【修改方式】在 errors.go 文件末尾追加错误码定义
     【相关依赖】无
     【修改内容】
        - 定义论文错误码 102001-102010：ErrPaperNotFound（论文不存在）、ErrPaperDuplicate（论文重复）、ErrPaperNotAllowModify（论文不允许修改）、ErrPaperStatusInvalid（论文状态不允许操作）等

- [x] 1.3.2 在 errors.go 中新增审核相关错误码
     【目标对象】`src/backend/pkg/errors/errors.go`
     【修改目的】定义审核管理相关错误码
     【修改方式】在 errors.go 文件末尾追加错误码定义
     【相关依赖】无
     【修改内容】
        - 定义审核错误码 103001-103010：ErrNoReviewPermission（无审核权限）、ErrAlreadyReviewed（已审核）、ErrPaperStatusNotAllowedReview（论文状态不允许审核）等

- [x] 1.3.3 在 errors.go 中新增文件相关错误码
     【目标对象】`src/backend/pkg/errors/errors.go`
     【修改目的】定义文件管理相关错误码
     【修改方式】在 errors.go 文件末尾追加错误码定义
     【相关依赖】无
     【修改内容】
        - 定义文件错误码 104001-104010：ErrFileTooLarge（文件过大）、ErrInvalidFileType（文件格式错误）、ErrUploadFailed（上传失败）等

#### 1.4 安装Excelize依赖
- [x] 1.4.1 安装 Excelize 依赖
     【目标对象】`src/backend/go.mod`
     【修改目的】添加 Excel 处理依赖
     【修改方式】执行命令添加依赖
     【相关依赖】无
     【修改内容】
        - 执行 `go get github.com/xuri/excelize/v2`
        - 执行 `go mod tidy` 整理依赖
        - 验证 go.mod 中已添加 `github.com/xuri/excelize/v2` 依赖

#### 1.5 更新数据库迁移
- [x] 1.5.1 注册论文管理相关实体到数据库迁移
     【目标对象】`src/backend/internal/database/migration.go`
     【修改目的】在数据库迁移中注册新的实体模型
     【修改方式】在 AutoMigrate 函数的 db.AutoMigrate 参数列表中新增实体
     【相关依赖】entity 目录下的所有实体文件
     【修改内容】
        - 在 db.AutoMigrate 函数中追加：&entity.Paper{}、&entity.Author{}、&entity.Project{}、&entity.Journal{}、&entity.PaperProject{}、&entity.ReviewLog{}、&entity.Archive{}、&entity.Attachment{}
        - 保持原有 User 和 OperationLog 实体的注册

- [x] 1.5.2 为 users 表新增 is_email_notify 字段
     【目标对象】`src/backend/internal/model/entity/user.go`
     【修改目的】支持邮件通知配置
     【修改方式】在 User 结构体中新增字段
     【相关依赖】无
     【修改内容】
        - 在 User 结构体中新增 IsEmailNotify 字段（bool类型，默认true）
        - 添加 GORM 标签定义字段约束

- [x] 1.5.3 为 operation_logs 表支持审核类型操作
     【目标对象】`src/backend/internal/model/entity/operation_log.go`
     【修改目的】扩展操作日志类型以支持审核操作
     【修改方式】检查现有结构，确保支持审核类型（如需要修改则更新）
     【相关依赖】无
     【修改内容】
        - 检查 OperationType 字段是否支持审核类型
        - 如需修改，添加审核相关的操作类型常量（如：业务审核、政工审核、提交审核等）

### 2. 后端Repository层实现

#### 2.1 创建论文Repository
- [x] 2.1.1 创建 PaperRepository 结构体和构造函数
     【目标对象】`src/backend/internal/repository/paper_repository.go`
     【修改目的】实现论文数据访问基础结构
     【修改方式】创建新文件，定义 PaperRepository 结构体和 NewPaperRepository 函数
     【相关依赖】GORM框架
     【修改内容】
        - 定义 PaperRepository 结构体，包含 db *gorm.DB 字段
        - 创建 NewPaperRepository 构造函数，接收 db 参数并返回实例

- [x] 2.1.2 实现 PaperRepository 的 CRUD 方法
     【目标对象】`src/backend/internal/repository/paper_repository.go`
     【修改目的】实现论文的增删改查方法
     【修改方式】在 PaperRepository 中添加方法
     【相关依赖】entity.Paper
     【修改内容】
        - Create 方法：接收 Paper 指针，调用 db.Create 创建记录
        - Update 方法：接收 Paper 指针，调用 db.Save 更新记录
        - Delete 方法：接收 id 参数，调用 db.Delete 删除记录（软删除）
        - FindByID 方法：接收 id 参数，返回 Paper 或 nil
        - List 方法：接收 page、size 参数，返回 Paper 列表和总数，按创建时间倒序

- [x] 2.1.3 实现 PaperRepository 的业务查询方法
     【目标对象】`src/backend/internal/repository/paper_repository.go`
     【修改目的】实现论文的业务查询方法
     【修改方式】在 PaperRepository 中添加方法
     【相关依赖】entity.Paper
     【修改内容】
        - ListByStatus 方法：接收 status 参数，返回该状态的论文列表
        - FindDuplicate 方法：接收 title 和 doi 参数，返回是否存在重复的论文
        - ListBySubmitter 方法：接收 submitterID 参数，返回该用户提交的论文列表
        - UpdateStatus 方法：接收 id 和 status 参数，更新论文状态

#### 2.2 创建作者Repository
- [x] 2.2.1 创建 AuthorRepository 结构体和基础方法
     【目标对象】`src/backend/internal/repository/author_repository.go`
     【修改目的】实现作者数据访问基础结构
     【修改方式】创建新文件，定义 AuthorRepository 结构体和基础方法
     【相关依赖】GORM框架
     【修改内容】
        - 定义 AuthorRepository 结构体，包含 db *gorm.DB 字段
        - 创建 NewAuthorRepository 构造函数
        - Create 方法：创建作者记录
        - Update 方法：更新作者记录
        - Delete 方法：删除作者记录
        - FindByID 方法：根据ID查询作者

- [x] 2.2.2 实现 AuthorRepository 的批量操作方法
     【目标对象】`src/backend/internal/repository/author_repository.go`
     【修改目的】实现作者的批量操作和查询方法
     【修改方式】在 AuthorRepository 中添加方法
     【相关依赖】entity.Author
     【修改内容】
        - ListByPaperID 方法：接收 paperID 参数，返回该论文的所有作者，按 Rank 排序
        - DeleteByPaperID 方法：接收 paperID 参数，删除该论文的所有作者
        - CreateBatch 方法：接收 Author 数组，批量创建作者记录（使用事务）

#### 2.3 创建课题Repository
- [x] 2.3.1 创建 ProjectRepository 结构体和基础方法
     【目标对象】`src/backend/internal/repository/project_repository.go`
     【修改目的】实现课题数据访问基础结构
     【修改方式】创建新文件，定义 ProjectRepository 结构体和基础方法
     【相关依赖】GORM框架
     【修改内容】
        - 定义 ProjectRepository 结构体，包含 db *gorm.DB 字段
        - 创建 NewProjectRepository 构造函数
        - Create 方法：创建课题记录
        - Update 方法：更新课题记录
        - Delete 方法：删除课题记录
        - FindByID 方法：根据ID查询课题

- [x] 2.3.2 实现 ProjectRepository 的业务查询方法
     【目标对象】`src/backend/internal/repository/project_repository.go`
     【修改目的】实现课题的业务查询方法
     【修改方式】在 ProjectRepository 中添加方法
     【相关依赖】entity.Project
     【修改内容】
        - FindByCode 方法：接收 code 参数，根据课题编号查询课题
        - List 方法：接收 page、size 参数，分页查询课题列表
        - CheckIsLinked 方法：接收 projectID 参数，检查课题是否已关联论文，返回关联数量

#### 2.4 创建期刊Repository
- [x] 2.4.1 创建 JournalRepository 结构体和基础方法
     【目标对象】`src/backend/internal/repository/journal_repository.go`
     【修改目的】实现期刊数据访问基础结构
     【修改方式】创建新文件，定义 JournalRepository 结构体和基础方法
     【相关依赖】GORM框架
     【修改内容】
        - 定义 JournalRepository 结构体，包含 db *gorm.DB 字段
        - 创建 NewJournalRepository 构造函数
        - Create 方法：创建期刊记录
        - Update 方法：更新期刊记录
        - FindByID 方法：根据ID查询期刊

- [x] 2.4.2 实现 JournalRepository 的查询方法
     【目标对象】`src/backend/internal/repository/journal_repository.go`
     【修改目的】实现期刊的查询方法
     【修改方式】在 JournalRepository 中添加方法
     【相关依赖】entity.Journal
     【修改内容】
        - FindByName 方法：接收 name 参数，模糊匹配期刊名称
        - FindByISSN 方法：接收 issn 参数，根据ISSN查询期刊
        - List 方法：接收 page、size 参数，分页查询期刊列表

#### 2.5 创建论文课题关联Repository
- [x] 2.5.1 创建 PaperProjectRepository 结构体和基础方法
     【目标对象】`src/backend/internal/repository/paper_project_repository.go`
     【修改目的】实现论文-课题关联数据访问基础结构
     【修改方式】创建新文件，定义 PaperProjectRepository 结构体和基础方法
     【相关依赖】GORM框架
     【修改内容】
        - 定义 PaperProjectRepository 结构体，包含 db *gorm.DB 字段
        - 创建 NewPaperProjectRepository 构造函数
        - Create 方法：创建论文-课题关联记录

- [x] 2.5.2 实现 PaperProjectRepository 的关联操作方法
     【目标对象】`src/backend/internal/repository/paper_project_repository.go`
     【修改目的】实现论文-课题关联的查询和删除方法
     【修改方式】在 PaperProjectRepository 中添加方法
     【相关依赖】entity.PaperProject
     【修改内容】
        - FindByPaperID 方法：接收 paperID 参数，返回该论文关联的所有课题
        - DeleteByPaperID 方法：接收 paperID 参数，删除该论文的所有课题关联
        - CreateBatch 方法：接收 PaperProject 数组，批量创建关联记录（使用事务）

#### 2.6 创建审核记录Repository
- [x] 2.6.1 创建 ReviewRepository 结构体和基础方法
     【目标对象】`src/backend/internal/repository/review_repository.go`
     【修改目的】实现审核记录数据访问基础结构
     【修改方式】创建新文件，定义 ReviewRepository 结构体和基础方法
     【相关依赖】GORM框架
     【修改内容】
        - 定义 ReviewRepository 结构体，包含 db *gorm.DB 字段
        - 创建 NewReviewRepository 构造函数
        - Create 方法：创建审核记录
        - FindByID 方法：根据ID查询审核记录

- [x] 2.6.2 实现 ReviewRepository 的业务查询方法
     【目标对象】`src/backend/internal/repository/review_repository.go`
     【修改目的】实现审核记录的业务查询方法
     【修改方式】在 ReviewRepository 中添加方法
     【相关依赖】entity.ReviewLog
     【修改内容】
        - FindByPaperID 方法：接收 paperID 参数，返回该论文的所有审核记录，按时间倒序
        - FindLatestByPaperIDAndType 方法：接收 paperID 和 reviewType 参数，返回该论文该类型的最新审核记录
        - ListByReviewer 方法：接收 reviewerID 参数，返回该审核人员的审核记录
        - ListPendingReview 方法：接收 reviewType 参数，查询待审核的论文列表（状态为"待业务审核"或"待政工审核"）

#### 2.7 创建附件Repository
- [x] 2.7.1 创建 AttachmentRepository 结构体和基础方法
     【目标对象】`src/backend/internal/repository/attachment_repository.go`
     【修改目的】实现附件数据访问基础结构
     【修改方式】创建新文件，定义 AttachmentRepository 结构体和基础方法
     【相关依赖】GORM框架
     【修改内容】
        - 定义 AttachmentRepository 结构体，包含 db *gorm.DB 字段
        - 创建 NewAttachmentRepository 构造函数
        - Create 方法：创建附件记录
        - FindByID 方法：根据ID查询附件
        - Delete 方法：删除附件记录（仅数据库记录）

- [x] 2.7.2 实现 AttachmentRepository 的查询方法
     【目标对象】`src/backend/internal/repository/attachment_repository.go`
     【修改目的】实现附件的查询方法
     【修改方式】在 AttachmentRepository 中添加方法
     【相关依赖】entity.Attachment
     【修改内容】
        - ListByPaperID 方法：接收 paperID 参数，返回该论文的所有附件
        - DeleteByPaperID 方法：接收 paperID 参数，删除该论文的所有附件记录

#### 2.8 创建归档Repository
- [x] 2.8.1 创建 ArchiveRepository 结构体和基础方法
     【目标对象】`src/backend/internal/repository/archive_repository.go`
     【修改目的】实现归档记录数据访问基础结构
     【修改方式】创建新文件，定义 ArchiveRepository 结构体和基础方法
     【相关依赖】GORM框架
     【修改内容】
        - 定义 ArchiveRepository 结构体，包含 db *gorm.DB 字段
        - 创建 NewArchiveRepository 构造函数
        - Create 方法：创建归档记录
        - FindByPaperID 方法：接收 paperID 参数，返回该论文的归档记录

### 3. 后端Service层实现

#### 3.1 创建论文Service
- [x] 3.1.1 创建 PaperService 结构体和构造函数
     【目标对象】`src/backend/internal/service/paper_service.go`
     【修改目的】实现论文业务逻辑基础结构
     【修改方式】创建新文件，定义 PaperService 结构体和 NewPaperService 函数
     【相关依赖】paper_repository, author_repository, paper_project_repository, attachment_repository, journal_repository
     【修改内容】
        - 定义 PaperService 结构体，包含 paperRepo、authorRepo、projectRepo、attachmentRepo、journalRepo 字段
        - 创建 NewPaperService 构造函数，接收所有 repository 参数并返回实例

- [x] 3.1.2 实现论文创建和更新方法
     【目标对象】`src/backend/internal/service/paper_service.go`
     【修改目的】实现论文创建、更新和删除业务逻辑
     【修改方式】在 PaperService 中添加方法
     【相关依赖】request.PaperRequest, entity.Paper
     【修改内容】
        - CreatePaper 方法：接收创建请求，使用事务创建论文记录、作者记录、课题关联、附件记录，校验论文数据格式，校验重复，返回论文ID
        - UpdatePaper 方法：接收论文ID和更新请求，校验论文状态（仅草稿状态可修改），使用事务更新论文、作者、课题关联，校验重复
        - DeletePaper 方法：接收论文ID，校验论文状态（仅草稿状态可删除），使用事务删除论文及其关联数据

- [x] 3.1.3 实现论文查询方法
     【目标对象】`src/backend/internal/service/paper_service.go`
     【修改目的】实现论文查询业务逻辑
     【修改方式】在 PaperService 中添加方法
     【相关依赖】entity.Paper, dto.PaperDTO
     【修改内容】
        - GetPaperByID 方法：接收论文ID，查询论文详情，包含作者、课题、期刊、附件信息，转换为 DTO 返回
        - ListPapers 方法：接收 page、size、status、keyword 参数，分页查询论文列表，支持按标题、作者、期刊、状态筛选，转换为 DTO 列表返回
        - GetMyPapers 方法：接收 userID、page、size 参数，查询该用户提交的论文列表

- [x] 3.1.4 实现论文提交和草稿保存方法
     【目标对象】`src/backend/internal/service/paper_service.go`
     【修改目的】实现论文提交审核和草稿保存业务逻辑
     【修改方式】在 PaperService 中添加方法
     【相关依赖】entity.Paper
     【修改内容】
        - SubmitForReview 方法：接收论文ID和操作人ID，校验论文状态（仅草稿状态可提交），校验必填字段完整性，更新状态为"待业务审核"，记录提交时间和提交人，记录操作日志
        - SaveDraft 方法：接收论文ID和更新请求，保存论文草稿，不进行严格校验和重复校验，状态保持为草稿

- [x] 3.1.5 实现论文重复校验方法
     【目标对象】`src/backend/internal/service/paper_service.go`
     【修改目的】实现论文重复校验业务逻辑
     【修改方式】在 PaperService 中添加方法
     【相关依赖】entity.Paper
     【修改内容】
        - CheckDuplicate 方法：接收标题和DOI参数，查询是否存在重复论文（相同标题或相同DOI），返回是否重复和重复论文信息

- [x] 3.1.6 实现论文批量导入方法
     【目标对象】`src/backend/internal/service/paper_service.go`
     【修改目的】实现Excel批量导入论文业务逻辑
     【修改方式】在 PaperService 中添加方法
     【相关依赖】excelize库
     【修改内容】
        - BatchImportPapers 方法：接收Excel文件，使用 excelize 库读取Excel内容，校验数据格式和重复性，批量创建论文记录，返回成功数量、失败数量和错误详情列表

 - [x] 3.1.7 实现论文数据校验方法
      【目标对象】`src/backend/internal/service/paper_service.go`
      【修改目的】实现论文数据格式校验逻辑
      【修改方式】在 PaperService 中添加方法
      【相关依赖】无
      【修改内容】
         - ValidatePaperData 方法：接收论文数据，校验DOI格式（正则表达式）、ISSN格式（正则表达式）、日期格式、必填字段，返回校验错误列表

#### 3.2 创建作者Service
- [x] 3.2.1 创建 AuthorService 结构体和构造函数
     【目标对象】`src/backend/internal/service/author_service.go`
     【修改目的】实现作者管理业务逻辑基础结构
     【修改方式】创建新文件，定义 AuthorService 结构体和 NewAuthorService 函数
     【相关依赖】author_repository
     【修改内容】
        - 定义 AuthorService 结构体，包含 authorRepo 字段
        - 创建 NewAuthorService 构造函数

- [x] 3.2.2 实现作者管理方法
     【目标对象】`src/backend/internal/service/author_service.go`
     【修改目的】实现作者的增删改查业务逻辑
     【修改方式】在 AuthorService 中添加方法
     【相关依赖】request.AuthorRequest, entity.Author
     【修改内容】
        - CreateAuthor 方法：创建作者记录，校验作者类型和排名
        - UpdateAuthor 方法：更新作者记录，校验作者类型和排名
        - DeleteAuthor 方法：删除作者记录
        - GetAuthorsByPaperID 方法：获取某论文的所有作者，按排名排序

- [x] 3.2.3 实现作者批量操作方法
     【目标对象】`src/backend/internal/service/author_service.go`
     【修改目的】实现作者批量操作和校验业务逻辑
     【修改方式】在 AuthorService 中添加方法
     【相关依赖】entity.Author
     【修改内容】
        - BatchCreateAuthors 方法：接收作者数组，使用事务批量创建作者，校验作者排名唯一性和作者类型互斥性（第一作者、共同第一作者、通讯作者不能重复）
        - UpdateRankings 方法：接收论文ID和作者顺序数组，更新作者排名，校验排名连续性
        - ValidateAuthorData 方法：校验作者数据，检查作者类型互斥性和排名合理性

#### 3.3 创建课题Service
- [x] 3.3.1 创建 ProjectService 结构体和构造函数
     【目标对象】`src/backend/internal/service/project_service.go`
     【修改目的】实现课题管理业务逻辑基础结构
     【修改方式】创建新文件，定义 ProjectService 结构体和 NewProjectService 函数
     【相关依赖】project_repository
     【修改内容】
        - 定义 ProjectService 结构体，包含 projectRepo 字段
        - 创建 NewProjectService 构造函数

- [x] 3.3.2 实现课题管理方法
     【目标对象】`src/backend/internal/service/project_service.go`
     【修改目的】实现课题的增删改查业务逻辑
     【修改方式】在 ProjectService 中添加方法
     【相关依赖】request.ProjectRequest, entity.Project
     【修改内容】
        - CreateProject 方法：创建课题记录，校验课题编号唯一性
        - UpdateProject 方法：更新课题记录，校验课题编号唯一性
        - DeleteProject 方法：接收课题ID，检查是否已关联论文，如已关联则不允许删除
        - GetProjectByID 方法：获取课题详情
        - ListProjects 方法：分页查询课题列表，支持按名称、编号、类型筛选

#### 3.4 创建期刊Service
- [x] 3.4.1 创建 JournalService 结构体和构造函数
     【目标对象】`src/backend/internal/service/journal_service.go`
     【修改目的】实现期刊管理业务逻辑基础结构
     【修改方式】创建新文件，定义 JournalService 结构体和 NewJournalService 函数
     【相关依赖】journal_repository
     【修改内容】
        - 定义 JournalService 结构体，包含 journalRepo 字段
        - 创建 NewJournalService 构造函数

- [x] 3.4.2 实现期刊管理方法
     【目标对象】`src/backend/internal/service/journal_service.go`
     【修改目的】实现期刊的增删改查业务逻辑
     【修改方式】在 JournalService 中添加方法
     【相关依赖】request.JournalRequest, entity.Journal
     【修改内容】
        - CreateJournal 方法：创建期刊记录，校验ISSN唯一性
        - UpdateJournal 方法：更新期刊记录，校验ISSN唯一性
        - GetJournalByID 方法：获取期刊详情
        - SearchJournals 方法：接收关键字，搜索期刊（按名称或ISSN模糊匹配）
        - ListJournals 方法：分页查询期刊列表
        - UpdateImpactFactor 方法：更新期刊影响因子

#### 3.5 创建审核Service
- [x] 3.5.1 创建 ReviewService 结构体和构造函数
     【目标对象】`src/backend/internal/service/review_service.go`
     【修改目的】实现审核业务逻辑基础结构
     【修改方式】创建新文件，定义 ReviewService 结构体和 NewReviewService 函数
     【相关依赖】review_repository, paper_repository, notification_service
     【修改内容】
        - 定义 ReviewService 结构体，包含 reviewRepo、paperRepo、notificationService 字段
        - 创建 NewReviewService 构造函数

- [x] 3.5.2 实现业务审核方法
     【目标对象】`src/backend/internal/service/review_service.go`
     【修改目的】实现业务审核业务逻辑
     【修改方式】在 ReviewService 中添加方法
     【相关依赖】entity.ReviewLog
     【修改内容】
        - BusinessReview 方法：接收论文ID、审核结果、审核意见、审核人ID，校验审核权限和论文状态（必须是"待业务审核"），如通过则更新论文状态为"待政工审核"，如驳回则更新论文状态为"草稿"，创建审核记录，发送审核结果通知

- [x] 3.5.3 实现政工审核方法
     【目标对象】`src/backend/internal/service/review_service.go`
     【修改目的】实现政工审核业务逻辑
     【修改方式】在 ReviewService 中添加方法
     【相关依赖】entity.ReviewLog
     【修改内容】
        - PoliticalReview 方法：接收论文ID、审核结果、审核意见、审核人ID，校验审核权限和论文状态（必须是"待政工审核"），如通过则更新论文状态为"审核通过"并调用归档流程，如驳回则更新论文状态为"草稿"，创建审核记录，发送审核结果通知

- [x] 3.5.4 实现审核查询方法
     【目标对象】`src/backend/internal/service/review_service.go`
     【修改目的】实现审核记录查询业务逻辑
     【修改方式】在 ReviewService 中添加方法
     【相关依赖】entity.ReviewLog
     【修改内容】
        - GetReviewLogsByPaperID 方法：接收论文ID，返回该论文的所有审核记录
        - GetPendingPapersForBusinessReview 方法：返回待业务审核的论文列表，包含提交时间、距提交天数等信息
        - GetPendingPapersForPoliticalReview 方法：返回待政工审核的论文列表，包含提交时间、距提交天数、业务审核意见等信息

- [x] 3.5.5 实现审核辅助方法
     【目标对象】`src/backend/internal/service/review_service.go`
     【修改目的】实现审核权限校验、时限检查等辅助逻辑
     【修改方式】在 ReviewService 中添加方法
     【相关依赖】无
     【修改内容】
        - ValidateReviewPermission 方法：接收审核人ID和审核类型，校验用户是否有审核权限（检查用户角色）
        - CheckReviewDeadline 方法：接收论文提交时间，计算距离提交的工作日数，判断是否超过3个工作日
        - SendReviewReminder 方法：接收审核人ID，发送审核时限提醒通知（调用 notification_service）
        - ProcessRejectPaper 方法：处理驳回论文，重置状态为草稿，发送驳回通知（包含驳回原因）
        - ProcessApprovePaper 方法：处理审核通过论文，调用归档Service创建归档记录，发送通过通知

#### 3.6 创建附件Service
- [x] 3.6.1 创建 FileService 结构体和构造函数
     【目标对象】`src/backend/internal/service/file_service.go`
     【修改目的】实现文件管理业务逻辑基础结构
     【修改方式】创建新文件，定义 FileService 结构体和 NewFileService 函数
     【相关依赖】attachment_repository
     【修改内容】
        - 定义 FileService 结构体，包含 attachmentRepo 字段和 uploadDir 字段（上传目录）
        - 创建 NewFileService 构造函数，接收 repository 和上传目录路径

- [x] 3.6.2 实现文件上传方法
     【目标对象】`src/backend/internal/service/file_service.go`
     【修改目的】实现文件上传业务逻辑
     【修改方式】在 FileService 中添加方法
     【相关依赖】entity.Attachment
     【修改内容】
        - UploadFile 方法：接收论文ID、文件类型、文件对象，校验文件大小（不超过100MB）和格式（仅允许PDF、JPG、PNG），生成唯一文件名和存储路径，保存文件到磁盘，创建附件记录，返回附件ID
        - ValidateFile 方法：校验文件大小和格式，返回校验错误
        - GenerateFilePath 方法：根据论文ID、文件类型、时间戳生成文件存储路径
        - CheckFileExists 方法：检查文件路径是否存在

- [x] 3.6.3 实现文件查询和删除方法
     【目标对象】`src/backend/internal/service/file_service.go`
     【修改目的】实现文件查询和删除业务逻辑
     【修改方式】在 FileService 中添加方法
     【相关依赖】entity.Attachment
     【修改内容】
        - GetFileByID 方法：接收附件ID，返回附件信息
        - GetFilesByPaperID 方法：接收论文ID，返回该论文的所有附件
        - DeleteFile 方法：接收附件ID，删除磁盘上的物理文件和数据库记录，校验文件是否存在

#### 3.7 创建通知Service
- [x] 3.7.1 创建 NotificationService 结构体和构造函数
     【目标对象】`src/backend/internal/service/notification_service.go`
     【修改目的】实现通知发送业务逻辑基础结构
     【修改方式】创建新文件，定义 NotificationService 结构体和 NewNotificationService 函数
     【相关依赖】user_repository
     【修改内容】
        - 定义 NotificationService 结构体，包含 userRepo 字段
        - 创建 NewNotificationService 构造函数

- [x] 3.7.2 实现通知发送方法
     【目标对象】`src/backend/internal/service/notification_service.go`
     【修改目的】实现各类通知发送业务逻辑
     【修改方式】在 NotificationService 中添加方法
     【相关依赖】无
     【修改内容】
        - SendSubmitNotification 方法：接收论文ID和提交人ID，发送提交审核通知给业务审核人员，记录操作日志
        - SendReviewNotification 方法：接收论文ID、审核类型、审核结果，发送审核结果通知给提交人（通过或驳回），记录操作日志
        - SendRejectNotification 方法：接收论文ID和驳回原因，发送驳回通知给提交人（包含驳回原因），记录操作日志
        - SendDeadlineReminder 方法：接收审核人ID和论文ID，发送审核时限提醒通知给审核人员，记录操作日志
        - SendApprovalNotification 方法：接收论文ID，发送审核通过通知给提交人，记录操作日志
        - SendBusinessReviewPassedNotification 方法：接收论文ID，发送业务审核通过通知给政工审核人员，记录操作日志

#### 3.8 创建归档Service
- [x] 3.8.1 创建 ArchiveService 结构体和构造函数
     【目标对象】`src/backend/internal/service/archive_service.go`
     【修改目的】实现归档业务逻辑基础结构
     【修改方式】创建新文件，定义 ArchiveService 结构体和 NewArchiveService 函数
     【相关依赖】archive_repository
     【修改内容】
        - 定义 ArchiveService 结构体，包含 archiveRepo 字段
        - 创建 NewArchiveService 构造函数

- [x] 3.8.2 实现归档方法
     【目标对象】`src/backend/internal/service/archive_service.go`
     【修改目的】实现论文归档业务逻辑
     【修改方式】在 ArchiveService 中添加方法
     【相关依赖】entity.Archive
     【修改内容】
        - ArchivePaper 方法：接收论文ID和归档人ID，生成唯一归档编号（格式：ARCH-YYYYMMDD-XXXX），创建归档记录，更新论文状态为"已归档"，记录操作日志
        - GetArchiveByPaperID 方法：接收论文ID，返回归档记录

### 4. 后端Handler层实现

#### 4.1 创建论文Handler
- [x] 4.1.1 创建 PaperHandler 结构体和构造函数
     【目标对象】`src/backend/internal/handler/paper_handler.go`
     【修改目的】实现论文HTTP接口基础结构
     【修改方式】创建新文件，定义 PaperHandler 结构体、请求结构体和 NewPaperHandler 函数
     【相关依赖】paper_service
     【修改内容】
        - 定义 CreatePaperRequest、UpdatePaperRequest 请求结构体（包含标签验证）
        - 定义 PaperHandler 结构体，包含 paperService 字段
        - 创建 NewPaperHandler 构造函数

- [x] 4.1.2 实现论文CRUD接口
     【目标对象】`src/backend/internal/handler/paper_handler.go`
     【修改目的】实现论文增删改查HTTP接口
     【修改方式】在 PaperHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - CreatePaper 方法：POST /api/papers，绑定请求参数，调用 paperService.CreatePaper，返回成功或错误响应
        - GetPaper 方法：GET /api/papers/:id，解析论文ID，调用 paperService.GetPaperByID，返回论文详情或404
        - ListPapers 方法：GET /api/papers，解析分页和查询参数，调用 paperService.ListPapers，返回论文列表
        - UpdatePaper 方法：PUT /api/papers/:id，解析论文ID和请求参数，调用 paperService.UpdatePaper，返回成功或错误
        - DeletePaper 方法：DELETE /api/papers/:id，解析论文ID，调用 paperService.DeletePaper，返回成功或错误

- [x] 4.1.3 实现论文提交和草稿保存接口
     【目标对象】`src/backend/internal/handler/paper_handler.go`
     【修改目的】实现论文提交审核和草稿保存HTTP接口
     【修改方式】在 PaperHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - SubmitForReview 方法：POST /api/papers/:id/submit，解析论文ID和操作人ID，调用 paperService.SubmitForReview，返回成功或错误
        - SaveDraft 方法：POST /api/papers/:id/save-draft，解析论文ID和请求参数，调用 paperService.SaveDraft，返回成功或错误

- [x] 4.1.4 实现论文校验和批量导入接口
     【目标对象】`src/backend/internal/handler/paper_handler.go`
     【修改目的】实现论文重复校验和批量导入HTTP接口
     【修改方式】在 PaperHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - CheckDuplicate 方法：POST /api/papers/check-duplicate，绑定请求参数，调用 paperService.CheckDuplicate，返回重复检查结果
        - BatchImport 方法：POST /api/papers/batch-import，接收Excel文件，调用 paperService.BatchImportPapers，返回导入结果（成功数、失败数、错误详情）
        - DownloadImportTemplate 方法：GET /api/papers/import-template，返回Excel导入模板文件

- [x] 4.1.5 实现我的论文接口
     【目标对象】`src/backend/internal/handler/paper_handler.go`
     【修改目的】实现获取当前用户提交的论文HTTP接口
     【修改方式】在 PaperHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - GetMyPapers 方法：GET /api/papers/my，从上下文获取当前用户ID，调用 paperService.GetMyPapers，返回论文列表

#### 4.2 创建审核Handler
- [x] 4.2.1 创建 ReviewHandler 结构体和构造函数
     【目标对象】`src/backend/internal/handler/review_handler.go`
     【修改目的】实现审核HTTP接口基础结构
     【修改方式】创建新文件，定义 ReviewHandler 结构体、请求结构体和 NewReviewHandler 函数
     【相关依赖】review_service
     【修改内容】
        - 定义 ReviewRequest 请求结构体（Result、Comment）
        - 定义 ReviewHandler 结构体，包含 reviewService 字段
        - 创建 NewReviewHandler 构造函数

- [x] 4.2.2 实现审核操作接口
     【目标对象】`src/backend/internal/handler/review_handler.go`
     【修改目的】实现业务审核和政工审核HTTP接口
     【修改方式】在 ReviewHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - BusinessReview 方法：POST /api/reviews/business/:paperId，解析论文ID和请求参数，获取当前审核人ID，调用 reviewService.BusinessReview，返回成功或错误
        - PoliticalReview 方法：POST /api/reviews/political/:paperId，解析论文ID和请求参数，获取当前审核人ID，调用 reviewService.PoliticalReview，返回成功或错误

- [x] 4.2.3 实现审核查询接口
     【目标对象】`src/backend/internal/handler/review_handler.go`
     【修改目的】实现审核记录查询HTTP接口
     【修改方式】在 ReviewHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - GetReviewLogs 方法：GET /api/reviews/:paperId/logs，解析论文ID，调用 reviewService.GetReviewLogsByPaperID，返回审核记录列表
        - GetPendingBusinessReviews 方法：GET /api/reviews/pending/business，获取当前审核人待业务审核的论文列表
        - GetPendingPoliticalReviews 方法：GET /api/reviews/pending/political，获取当前审核人待政工审核的论文列表
        - GetMyReviews 方法：GET /api/reviews/my，获取当前审核人的审核历史记录

#### 4.3 创建课题Handler
- [x] 4.3.1 创建 ProjectHandler 结构体和构造函数
     【目标对象】`src/backend/internal/handler/project_handler.go`
     【修改目的】实现课题管理HTTP接口基础结构
     【修改方式】创建新文件，定义 ProjectHandler 结构体、请求结构体和 NewProjectHandler 函数
     【相关依赖】project_service
     【修改内容】
        - 定义 CreateProjectRequest、UpdateProjectRequest 请求结构体
        - 定义 ProjectHandler 结构体，包含 projectService 字段
        - 创建 NewProjectHandler 构造函数

- [x] 4.3.2 实现课题CRUD接口
     【目标对象】`src/backend/internal/handler/project_handler.go`
     【修改目的】实现课题增删改查HTTP接口
     【修改方式】在 ProjectHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - CreateProject 方法：POST /api/projects，绑定请求参数，调用 projectService.CreateProject，返回成功或错误
        - GetProject 方法：GET /api/projects/:id，解析课题ID，调用 projectService.GetProjectByID，返回课题详情或404
        - ListProjects 方法：GET /api/projects，解析分页和查询参数，调用 projectService.ListProjects，返回课题列表
        - UpdateProject 方法：PUT /api/projects/:id，解析课题ID和请求参数，调用 projectService.UpdateProject，返回成功或错误
        - DeleteProject 方法：DELETE /api/projects/:id，解析课题ID，调用 projectService.DeleteProject，返回成功或错误

#### 4.4 创建期刊Handler
- [x] 4.4.1 创建 JournalHandler 结构体和构造函数
     【目标对象】`src/backend/internal/handler/journal_handler.go`
     【修改目的】实现期刊管理HTTP接口基础结构
     【修改方式】创建新文件，定义 JournalHandler 结构体、请求结构体和 NewJournalHandler 函数
     【相关依赖】journal_service
     【修改内容】
        - 定义 CreateJournalRequest、UpdateJournalRequest、UpdateImpactFactorRequest 请求结构体
        - 定义 JournalHandler 结构体，包含 journalService 字段
        - 创建 NewJournalHandler 构造函数

- [x] 4.4.2 实现期刊CRUD接口
     【目标对象】`src/backend/internal/handler/journal_handler.go`
     【修改目的】实现期刊增删改查HTTP接口
     【修改方式】在 JournalHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - CreateJournal 方法：POST /api/journals，绑定请求参数，调用 journalService.CreateJournal，返回成功或错误
        - GetJournal 方法：GET /api/journals/:id，解析期刊ID，调用 journalService.GetJournalByID，返回期刊详情或404
        - ListJournals 方法：GET /api/journals，解析分页参数，调用 journalService.ListJournals，返回期刊列表
        - UpdateJournal 方法：PUT /api/journals/:id，解析期刊ID和请求参数，调用 journalService.UpdateJournal，返回成功或错误
        - UpdateImpactFactor 方法：PUT /api/journals/:id/impact-factor，解析期刊ID和请求参数，调用 journalService.UpdateImpactFactor，返回成功或错误

- [x] 4.4.3 实现期刊搜索接口
     【目标对象】`src/backend/internal/handler/journal_handler.go`
     【修改目的】实现期刊搜索HTTP接口
     【修改方式】在 JournalHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - SearchJournals 方法：GET /api/journals/search，解析搜索关键字，调用 journalService.SearchJournals，返回匹配的期刊列表

#### 4.5 创建文件上传Handler
- [x] 4.5.1 创建 FileHandler 结构体和构造函数
     【目标对象】`src/backend/internal/handler/file_handler.go`
     【修改目的】实现文件管理HTTP接口基础结构
     【修改方式】创建新文件，定义 FileHandler 结构体和 NewFileHandler 函数
     【相关依赖】file_service
     【修改内容】
        - 定义 FileHandler 结构体，包含 fileService 字段
        - 创建 NewFileHandler 构造函数

- [x] 4.5.2 实现文件上传接口
     【目标对象】`src/backend/internal/handler/file_handler.go`
     【修改目的】实现文件上传HTTP接口
     【修改方式】在 FileHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - UploadFile 方法：POST /api/files/upload，接收multipart/form-data格式的文件，解析文件类型和文件对象，调用 fileService.UploadFile，返回附件ID和文件路径

- [x] 4.5.3 实现文件查询和下载接口
     【目标对象】`src/backend/internal/handler/file_handler.go`
     【修改目的】实现文件查询和下载HTTP接口
     【修改方式】在 FileHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - GetFile 方法：GET /api/files/:id，解析附件ID，调用 fileService.GetFileByID，返回附件信息或404
        - DownloadFile 方法：GET /api/files/:id/download，解析附件ID，调用 fileService.GetFileByID 获取附件信息，读取文件内容，以二进制流方式返回文件，设置正确的 Content-Type 和 Content-Disposition 头
        - DeleteFile 方法：DELETE /api/files/:id，解析附件ID，调用 fileService.DeleteFile，返回成功或错误

#### 4.6 创建归档Handler
 - [x] 4.6.1 创建 ArchiveHandler 结构体和构造函数
     【目标对象】`src/backend/internal/handler/archive_handler.go`
     【修改目的】实现归档管理HTTP接口基础结构
     【修改方式】创建新文件，定义 ArchiveHandler 结构体和 NewArchiveHandler 函数
     【相关依赖】archive_service
     【修改内容】
        - 定义 ArchiveHandler 结构体，包含 archiveService 字段
        - 创建 NewArchiveHandler 构造函数

 - [x] 4.6.2 实现归档查询接口
     【目标对象】`src/backend/internal/handler/archive_handler.go`
     【修改目的】实现归档记录查询HTTP接口
     【修改方式】在 ArchiveHandler 中添加方法
     【相关依赖】gin框架
     【修改内容】
        - GetArchiveByPaper 方法：GET /api/archives/paper/:paperId，解析论文ID，调用 archiveService.GetArchiveByPaperID，返回归档记录或404
        - ListArchives 方法：GET /api/archives，解析分页参数，返回归档记录列表

### 5. 后端路由和中间件配置

#### 5.1 注册论文管理路由
 - [x] 5.1.1 在 main.go 中创建 Handler 实例
     【目标对象】`src/backend/cmd/server/main.go`
     【修改目的】初始化所有论文管理相关的 Handler 实例
     【修改方式】在 main 函数中创建 Handler 实例
     【相关依赖】所有 Service 和 Handler 文件
     【修改内容】
        - 在创建 Handler 的代码块中，添加：paperHandler := handler.NewPaperHandler(paperService)、reviewHandler := handler.NewReviewHandler(reviewService)、projectHandler := handler.NewProjectHandler(projectService)、journalHandler := handler.NewJournalHandler(journalService)、fileHandler := handler.NewFileHandler(fileService)、archiveHandler := handler.NewArchiveHandler(archiveService)

 - [x] 5.1.2 注册论文相关API路由
     【目标对象】`src/backend/cmd/server/main.go`
     【修改目的】注册论文管理、审核、文件上传等API路由组
     【修改方式】在路由配置部分新增路由组
     【相关依赖】所有 Handler 实例和权限中间件
     【修改内容】
        - 在路由配置中添加 papers 路由组：
          * POST /api/papers（创建论文），应用 paper:create 权限
          * GET /api/papers（论文列表），应用 paper:view 权限
          * GET /api/papers/:id（论文详情），应用 paper:view 权限
          * PUT /api/papers/:id（更新论文），应用 paper:edit 权限
          * DELETE /api/papers/:id（删除论文），应用 paper:delete 权限
          * POST /api/papers/:id/submit（提交审核），应用 paper:submit 权限
          * POST /api/papers/:id/save-draft（保存草稿），应用 paper:edit 权限
          * POST /api/papers/check-duplicate（检查重复），应用 paper:view 权限
          * POST /api/papers/batch-import（批量导入），应用 paper:import 权限
          * GET /api/papers/my（我的论文），应用 paper:view 权限
          * GET /api/papers/import-template（下载导入模板），应用 paper:view 权限
        - 添加 reviews 路由组：
          * POST /api/reviews/business/:paperId（业务审核），应用 review:business 权限
          * POST /api/reviews/political/:paperId（政工审核），应用 review:political 权限
          * GET /api/reviews/:paperId/logs（审核记录），应用 review:view 权限
          * GET /api/reviews/pending/business（待业务审核），应用 review:business 权限
          * GET /api/reviews/pending/political（待政工审核），应用 review:political 权限
          * GET /api/reviews/my（我的审核），应用 review:view 权限
        - 添加 projects 路由组：
          * POST /api/projects（创建课题），应用 project:manage 权限
          * GET /api/projects（课题列表），应用 project:view 权限
          * GET /api/projects/:id（课题详情），应用 project:view 权限
          * PUT /api/projects/:id（更新课题），应用 project:manage 权限
          * DELETE /api/projects/:id（删除课题），应用 project:manage 权限
        - 添加 journals 路由组：
          * POST /api/journals（创建期刊），应用 journal:manage 权限
          * GET /api/journals（期刊列表），应用 journal:view 权限
          * GET /api/journals/:id（期刊详情），应用 journal:view 权限
          * PUT /api/journals/:id（更新期刊），应用 journal:manage 权限
          * PUT /api/journals/:id/impact-factor（更新影响因子），应用 journal:manage 权限
          * GET /api/journals/search（搜索期刊），应用 journal:view 权限
        - 添加 files 路由组：
          * POST /api/files/upload（上传文件），应用 file:upload 权限
          * GET /api/files/:id（文件信息），应用 file:view 权限
          * GET /api/files/:id/download（下载文件），应用 file:download 权限
          * DELETE /api/files/:id（删除文件），应用 file:delete 权限
        - 添加 archives 路由组：
          * GET /api/archives（归档列表），应用 archive:view 权限
          * GET /api/archives/paper/:paperId（论文归档信息），应用 archive:view 权限

#### 5.2 配置文件上传中间件
 - [x] 5.2.1 创建文件大小限制中间件
      【目标对象】`src/backend/internal/middleware/file_size_middleware.go`
      【修改目的】创建限制上传文件大小的中间件
      【修改方式】创建新文件，实现文件大小限制逻辑
      【相关依赖】gin框架
      【修改内容】
         - 定义 FileSizeLimit 函数，接收 maxMB 参数（最大文件大小MB），返回 gin.HandlerFunc
         - 在中间件中检查请求头 Content-Length，如果超过限制则返回 413 Payload Too Large 错误
         - 设置 gin.Context 的 MaxBytesReader

 - [x] 5.2.2 应用文件大小限制中间件到文件上传路由
      【目标对象】`src/backend/cmd/server/main.go`
      【修改目的】在文件上传路由应用文件大小限制中间件
      【修改方式】在文件上传路由中应用 FileSizeLimitMiddleware
      【相关依赖】file_size_middleware
      【修改内容】
         - 在 POST /api/files/upload 路由前添加中间件：middleware.FileSizeLimit(100)，限制最大100MB

### 6. 前端基础设施准备

#### 6.1 安装前端依赖
 - [x] 6.1.1 安装前端依赖
     【目标对象】`src/frontend/package.json`
     【修改目的】添加表单管理和图表库依赖
     【修改方式】执行 npm install 命令安装依赖
     【相关依赖】无
     【修改内容】
        - 执行 `npm install react-hook-form @types/react-hook-form echarts`
        - 验证 package.json 中已添加 `react-hook-form`、`@types/react-hook-form`、`echarts` 依赖
        - 执行 `npm install` 安装依赖

#### 6.2 创建论文相关类型定义
- [x] 6.2.1 创建论文类型定义文件
     【目标对象】`src/frontend/src/types/paper.ts`
     【修改目的】定义论文相关的TypeScript类型
     【修改方式】创建新文件，导出类型定义
     【相关依赖】无
     【修改内容】
        - 定义 Paper 接口：id、title、abstract、journalId、journalName、doi、impactFactor、publishDate、status、submitterId、submitterName、submitTime、createdAt、updatedAt
        - 定义 PaperForm 接口：title、abstract、journalId、doi、impactFactor、publishDate
        - 定义 PaperStatus 类型：'draft' | 'pending_business' | 'pending_political' | 'approved' | 'rejected'
        - 导出类型和类型转换函数

- [x] 6.2.2 创建作者类型定义文件
     【目标对象】`src/frontend/src/types/author.ts`
     【修改目的】定义作者相关的TypeScript类型
     【修改方式】创建新文件，导出类型定义
     【相关依赖】无
     【修改内容】
        - 定义 Author 接口：id、paperId、name、authorType、rank、department、userId
        - 定义 AuthorForm 接口：name、authorType、rank、department、userId
        - 定义 AuthorType 类型：'first_author' | 'co_first_author' | 'corresponding_author' | 'author'
        - 导出类型和类型转换函数

- [x] 6.2.3 创建课题类型定义文件
     【目标对象】`src/frontend/src/types/project.ts`
     【修改目的】定义课题相关的TypeScript类型
     【修改方式】创建新文件，导出类型定义
     【相关依赖】无
     【修改内容】
        - 定义 Project 接口：id、name、code、projectType、source、level、status
        - 定义 ProjectForm 接口：name、code、projectType、source、level
        - 定义 ProjectType 类型：'vertical' | 'horizontal'
        - 定义 ProjectLevel 类型：'national' | 'provincial' | 'municipal'
        - 导出类型和类型转换函数

- [x] 6.2.4 创建审核类型定义文件
     【目标对象】`src/frontend/src/types/review.ts`
     【修改目的】定义审核相关的TypeScript类型
     【修改方式】创建新文件，导出类型定义
     【相关依赖】无
     【修改内容】
        - 定义 ReviewLog 接口：id、paperId、reviewType、result、comment、reviewerId、reviewerName、reviewTime、createdAt
        - 定义 ReviewForm 接口：result、comment
        - 定义 ReviewType 类型：'business' | 'political'
        - 定义 ReviewResult 类型：'approved' | 'rejected'
        - 定义 PendingReview 接口：paperId、title、submitterName、submitTime、status、daysSinceSubmit
        - 导出类型和类型转换函数

- [x] 6.2.5 创建通用类型定义文件
     【目标对象】`src/frontend/src/types/common.ts`
     【修改目的】定义通用TypeScript类型
     【修改方式】创建新文件，导出类型定义
     【相关依赖】无
     【修改内容】
        - 定义 PageResponse 接口：list、total、page、size
        - 定义 UploadFile 接口：file、onProgress、onSuccess、onError
        - 定义 FileItem 接口：id、fileName、filePath、fileSize、mimeType
        - 导出类型

#### 6.3 创建论文相关Service层
- [x] 6.3.1 创建论文Service
     【目标对象】`src/frontend/src/services/paperService.ts`
     【修改目的】封装论文相关API调用
     【修改方式】创建新文件，实现论文API调用方法
     【相关依赖】src/frontend/src/services/api.ts
     【修改内容】
        - 定义 createPaper 方法：POST /api/papers，创建论文
        - 定义 getPaper 方法：GET /api/papers/:id，获取论文详情
        - 定义 getPapers 方法：GET /api/papers，分页查询论文列表，支持查询参数
        - 定义 updatePaper 方法：PUT /api/papers/:id，更新论文
        - 定义 deletePaper 方法：DELETE /api/papers/:id，删除论文
        - 定义 submitForReview 方法：POST /api/papers/:id/submit，提交审核
        - 定义 saveDraft 方法：POST /api/papers/:id/save-draft，保存草稿
        - 定义 checkDuplicate 方法：POST /api/papers/check-duplicate，检查重复
        - 定义 batchImport 方法：POST /api/papers/batch-import，批量导入
        - 定义 getMyPapers 方法：GET /api/papers/my，获取我的论文
        - 定义 downloadImportTemplate 方法：GET /api/papers/import-template，下载导入模板
        - 导出所有方法

- [x] 6.3.2 创建作者Service
     【目标对象】`src/frontend/src/services/authorService.ts`
     【修改目的】封装作者相关API调用
     【修改方式】创建新文件，实现作者API调用方法
     【相关依赖】src/frontend/src/services/api.ts
     【修改内容】
        - 定义 getAuthorsByPaperId 方法：GET /api/papers/:paperId/authors，获取论文的作者列表
        - 定义 createAuthor 方法：POST /api/authors，创建作者
        - 定义 updateAuthor 方法：PUT /api/authors/:id，更新作者
        - 定义 deleteAuthor 方法：DELETE /api/authors/:id，删除作者
        - 定义 updateRankings 方法：PUT /api/papers/:paperId/authors/rankings，更新作者排名
        - 定义 searchUsers 方法：GET /api/users/search，从人员库搜索用户（用于关联作者）
        - 导出所有方法

- [x] 6.3.3 创建课题Service
     【目标对象】`src/frontend/src/services/projectService.ts`
     【修改目的】封装课题相关API调用
     【修改方式】创建新文件，实现课题API调用方法
     【相关依赖】src/frontend/src/services/api.ts
     【修改内容】
        - 定义 createProject 方法：POST /api/projects，创建课题
        - 定义 getProject 方法：GET /api/projects/:id，获取课题详情
        - 定义 getProjects 方法：GET /api/projects，分页查询课题列表
        - 定义 updateProject 方法：PUT /api/projects/:id，更新课题
        - 定义 deleteProject 方法：DELETE /api/projects/:id，删除课题
        - 定义 searchProjects 方法：GET /api/projects/search，搜索课题
        - 导出所有方法

- [x] 6.3.4 创建期刊Service
     【目标对象】`src/frontend/src/services/journalService.ts`
     【修改目的】封装期刊相关API调用
     【修改方式】创建新文件，实现期刊API调用方法
     【相关依赖】src/frontend/src/services/api.ts
     【修改内容】
        - 定义 createJournal 方法：POST /api/journals，创建期刊
        - 定义 getJournal 方法：GET /api/journals/:id，获取期刊详情
        - 定义 getJournals 方法：GET /api/journals，分页查询期刊列表
        - 定义 updateJournal 方法：PUT /api/journals/:id，更新期刊
        - 定义 searchJournals 方法：GET /api/journals/search，搜索期刊
        - 定义 updateImpactFactor 方法：PUT /api/journals/:id/impact-factor，更新影响因子
        - 导出所有方法

- [x] 6.3.5 创建审核Service
     【目标对象】`src/frontend/src/services/reviewService.ts`
     【修改目的】封装审核相关API调用
     【修改方式】创建新文件，实现审核API调用方法
     【相关依赖】src/frontend/src/services/api.ts
     【修改内容】
        - 定义 businessReview 方法：POST /api/reviews/business/:paperId，业务审核
        - 定义 politicalReview 方法：POST /api/reviews/political/:paperId，政工审核
        - 定义 getReviewLogs 方法：GET /api/reviews/:paperId/logs，获取审核记录
        - 定义 getPendingBusinessReviews 方法：GET /api/reviews/pending/business，获取待业务审核列表
        - 定义 getPendingPoliticalReviews 方法：GET /api/reviews/pending/political，获取待政工审核列表
        - 定义 getMyReviews 方法：GET /api/reviews/my，获取我的审核记录
        - 导出所有方法

- [x] 6.3.6 创建文件Service
     【目标对象】`src/frontend/src/services/fileService.ts`
     【修改目的】封装文件上传和下载API调用
     【修改方式】创建新文件，实现文件API调用方法
     【相关依赖】src/frontend/src/services/api.ts
     【修改内容】
        - 定义 uploadFile 方法：POST /api/files/upload，上传文件，接收 FormData
        - 定义 getFile 方法：GET /api/files/:id，获取文件信息
        - 定义 downloadFile 方法：GET /api/files/:id/download，下载文件
        - 定义 deleteFile 方法：DELETE /api/files/:id，删除文件
        - 导出所有方法

#### 6.4 创建论文状态管理Store
 - [x] 6.4.1 创建论文Store
     【目标对象】`src/frontend/src/stores/paperStore.ts`
     【修改目的】管理论文相关状态
     【修改方式】创建新文件，使用 Zustand 创建 Store
     【相关依赖】Zustand
     【修改内容】
        - 定义 PaperState 接口：currentPaper、paperList、total、page、size、loading、error
        - 定义 PaperActions 接口：setCurrentPaper、fetchPapers、fetchPaperById、createPaper、updatePaper、deletePaper、submitForReview、saveDraft、reset、setLoading、setError
        - 使用 create 创建 paperStore，实现状态管理和 Actions 方法
        - 在 fetchPapers 方法中调用 paperService.getPapers
        - 在 fetchPaperById 方法中调用 paperService.getPaper
        - 在 createPaper 方法中调用 paperService.createPaper
        - 在 updatePaper 方法中调用 paperService.updatePaper
        - 在 deletePaper 方法中调用 paperService.deletePaper
        - 在 submitForReview 方法中调用 paperService.submitForReview
        - 在 saveDraft 方法中调用 paperService.saveDraft
        - 导出 usePaperStore Hook

### 7. 前端页面组件实现

#### 7.1 创建论文列表页
- [x] 7.1.1 创建论文列表页组件基础结构
     【目标对象】`src/frontend/src/pages/paper/PaperListPage.tsx`
     【修改目的】实现论文列表查询和展示基础结构
     【修改方式】创建新文件，使用 React 和 Ant Design 组件
     【相关依赖】paperService, paperStore, useAuth
     【修改内容】
        - 导入必要的依赖和类型
        - 定义 PaperListPage 函数组件
        - 使用 usePaperStore 获取论文列表、total、page、size、loading
        - 定义状态：searchParams（搜索参数）、selectedRowKeys（选中的行）
        - 定义搜索表单字段：title、author、journal、status

- [x] 7.1.2 实现论文列表查询和展示
     【目标对象】`src/frontend/src/pages/paper/PaperListPage.tsx`
     【修改目的】实现论文列表查询和展示功能
     【修改方式】在 PaperListPage 组件中添加查询和展示逻辑
     【相关依赖】paperService, paperStore
     【修改内容】
        - 使用 useEffect 在组件挂载时调用 fetchPapers 方法
        - 实现搜索表单，包含标题、作者、期刊、状态等查询条件
        - 使用 Ant Design Table 展示论文列表，列包括：论文ID、标题、作者、期刊、状态、提交时间、操作
        - 实现分页功能，使用 Ant Design Pagination 组件
        - 根据状态显示不同颜色的标签（草稿-灰色、待审核-蓝色、审核通过-绿色、驳回-红色）
        - 实现表格刷新功能

- [x] 7.1.3 实现论文列表操作功能
     【目标对象】`src/frontend/src/pages/paper/PaperListPage.tsx`
     【修改目的】实现论文列表的操作按钮和功能
     【修改方式】在 PaperListPage 组件中添加操作按钮
     【相关依赖】useNavigate, paperStore
     【修改内容】
        - 实现操作列，包含：查看、编辑、删除、提交审核按钮
        - 查看按钮：跳转到论文详情页 /papers/:id
        - 编辑按钮：跳转到论文编辑页 /papers/:id/edit（仅草稿状态可编辑）
        - 删除按钮：调用 deletePaper 方法，确认后删除（仅草稿状态可删除）
        - 提交审核按钮：调用 submitForReview 方法（仅草稿状态可提交）
        - 集成权限控制，根据用户角色显示不同操作按钮
        - 实现批量操作（批量删除等）

#### 7.2 创建论文录入页
 - [x] 7.2.1 创建论文录入页组件基础结构
     【目标对象】`src/frontend/src/pages/paper/PaperCreatePage.tsx`
     【修改目的】实现论文信息录入基础结构
     【修改方式】创建新文件，使用 React Hook Form 管理表单
     【相关依赖】React Hook Form, paperService
     【修改内容】
        - 导入必要的依赖和类型
        - 定义 PaperCreatePage 函数组件
        - 使用 useForm 创建表单实例，定义表单字段：title、abstract、journalId、doi、impactFactor、publishDate
        - 定义状态：authors（作者列表）、projects（课题列表）、attachments（附件列表）、isSubmitting

 - [x] 7.2.2 实现论文基础字段表单
     【目标对象】`src/frontend/src/pages/paper/PaperCreatePage.tsx`
     【修改目的】实现论文基础字段表单和校验
     【修改方式】在 PaperCreatePage 组件中添加表单字段
     【相关依赖】React Hook Form, Ant Design Form
     【修改内容】
        - 使用 Ant Design Form.Item 实现标题输入框（必填，最大长度500）
        - 实现摘要输入框（必填，多行文本框，最大长度2000）
        - 实现期刊选择器（下拉框，支持搜索）
        - 实现 DOI 输入框（选填，格式校验正则表达式）
        - 实现影响因子输入框（选填，数字类型，保留3位小数）
        - 实现出版日期选择器（选填，日期格式）
        - 设置表单校验规则

 - [x] 7.2.3 集成作者、课题、附件和提交功能
     【目标对象】`src/frontend/src/pages/paper/PaperCreatePage.tsx`
     【修改目的】集成作者列表、课题选择器、文件上传和提交功能
     【修改方式】在 PaperCreatePage 组件中添加组件集成
     【相关依赖】AuthorList, ProjectSelector, FileUpload 组件
     【修改内容】
        - 集成 AuthorList 组件，管理作者信息（添加、删除、调整顺序）
        - 集成 ProjectSelector 组件，管理课题信息（搜索、选择、手动录入）
        - 集成 FileUpload 组件，实现附件上传（支持全文、首页、期刊封面、审批件）
        - 实现重复校验功能（在提交时调用 checkDuplicate）
        - 实现保存草稿按钮（调用 saveDraft 方法）
        - 实现提交审核按钮（调用 submitForReview 方法，包含完整校验）
        - 实现表单重置功能

#### 7.3 创建论文详情页
 - [x] 7.3.1 创建论文详情页组件基础结构
     【目标对象】`src/frontend/src/pages/paper/PaperDetailPage.tsx`
     【修改目的】展示论文完整信息基础结构
     【修改方式】创建新文件，使用 React 和 Ant Design 组件
     【相关依赖】paperService, reviewService
     【修改内容】
        - 导入必要的依赖和类型
        - 定义 PaperDetailPage 函数组件
        - 使用 useParams 获取论文ID
        - 定义状态：paper（论文详情）、reviewLogs（审核记录）、loading

 - [x] 7.3.2 实现论文信息展示
     【目标对象】`src/frontend/src/pages/paper/PaperDetailPage.tsx`
     【修改目的】展示论文基础信息、作者、课题、附件
     【修改方式】在 PaperDetailPage 组件中添加信息展示逻辑
     【相关依赖】Ant Design 组件
     【修改内容】
        - 使用 useEffect 调用 getPaper 方法获取论文详情
        - 使用 Ant Design Descriptions 组件展示论文基础信息（标题、摘要、期刊、DOI、影响因子、出版日期、状态、提交时间、提交人）
        - 展示作者列表（包含作者类型、排名、单位），使用 Table 或 List 组件
        - 展示课题信息（课题名称、编号、类型、级别）
        - 展示附件列表，支持下载功能（使用 fileService.downloadFile）
        - 根据状态显示不同颜色的标签

 - [x] 7.3.3 实现审核记录和操作按钮
     【目标对象】`src/frontend/src/pages/paper/PaperDetailPage.tsx`
     【修改目的】展示审核记录和操作按钮
     【修改方式】在 PaperDetailPage 组件中添加审核记录和操作按钮
     【相关依赖】reviewService, useNavigate
     【修改内容】
        - 使用 useEffect 调用 getReviewLogs 方法获取审核记录
        - 使用 Timeline 或 Table 展示审核记录（审核类型、审核结果、审核意见、审核人、审核时间）
        - 根据论文状态显示不同操作按钮：
          * 草稿状态：显示编辑、提交审核按钮
          * 驳回状态：显示编辑、重新提交按钮
          * 审核通过状态：无操作按钮
          * 待审核状态：无操作按钮（审核人员在审核列表页操作）
        - 编辑按钮：跳转到编辑页 /papers/:id/edit
        - 集成权限控制，根据用户角色显示不同操作按钮

#### 7.4 创建批量导入页
 - [x] 7.4.1 创建批量导入页组件基础结构
     【目标对象】`src/frontend/src/pages/paper/BatchImportPage.tsx`
     【修改目的】实现Excel批量导入论文基础结构
     【修改方式】创建新文件，使用 React 和 Ant Design 组件
     【相关依赖】paperService, fileService
     【修改内容】
        - 导入必要的依赖和类型
        - 定义 BatchImportPage 函数组件
        - 定义状态：file（上传的文件）、importResult（导入结果）、isUploading

 - [x] 7.4.2 实现导入模板下载和文件上传
     【目标对象】`src/frontend/src/pages/paper/BatchImportPage.tsx`
     【修改目的】实现导入模板下载和Excel文件上传功能
     【修改方式】在 BatchImportPage 组件中添加下载和上传逻辑
     【相关依赖】paperService, Ant Design Upload
     【修改内容】
        - 实现下载导入模板按钮，调用 paperService.downloadImportTemplate
        - 使用 Ant Design Upload 组件上传 Excel 文件
        - 限制文件类型为 .xlsx 和 .xls
        - 限制文件大小为 10MB
        - 实现文件上传前校验
        - 实现上传进度显示

 - [x] 7.4.3 实现批量导入和结果展示
     【目标对象】`src/frontend/src/pages/paper/BatchImportPage.tsx`
     【修改目的】实现批量导入和结果展示功能
     【修改方式】在 BatchImportPage 组件中添加导入和结果展示逻辑
     【相关依赖】paperService
     【修改内容】
        - 实现开始导入按钮，调用 paperService.batchImport
        - 展示导入结果：成功数量、失败数量
        - 展示错误详情表格（错误行号、错误原因、错误字段）
        - 实现重新导入功能
        - 实现导入历史记录（可选）
        - 实现返回论文列表页按钮

#### 7.5 创建业务审核列表页
- [x] 7.5.1 创建业务审核列表页组件基础结构
     【目标对象】`src/frontend/src/pages/review/BusinessReviewListPage.tsx`
     【修改目的】展示待业务审核的论文列表基础结构
     【修改方式】创建新文件，使用 React 和 Ant Design 组件
     【相关依赖】reviewService, useAuth
     【修改内容】
        - 导入必要的依赖和类型
        - 定义 BusinessReviewListPage 函数组件
        - 定义状态：papers（待审核论文列表）、loading、searchParams
        - 使用 useEffect 在组件挂载时获取待审核列表

- [x] 7.5.2 实现待审核论文列表展示
     【目标对象】`src/frontend/src/pages/review/BusinessReviewListPage.tsx`
     【修改目的】展示待业务审核的论文列表
     【修改方式】在 BusinessReviewListPage 组件中添加列表展示逻辑
     【相关依赖】reviewService, Ant Design Table
     【修改内容】
        - 调用 reviewService.getPendingBusinessReviews 获取待业务审核列表
        - 使用 Ant Design Table 展示论文列表，列包括：论文ID、标题、提交人、提交时间、距提交天数、状态
        - 使用 useEffect 定时刷新列表（每5分钟）
        - 实现搜索和筛选功能（按标题、提交人、提交时间）
        - 根据距提交天数显示不同颜色（超过2天显示红色）

- [x] 7.5.3 实现去审核功能
     【目标对象】`src/frontend/src/pages/review/BusinessReviewListPage.tsx`
     【修改目的】实现跳转到审核页面功能
     【修改方式】在 BusinessReviewListPage 组件中添加跳转逻辑
     【相关依赖】useNavigate
     【修改内容】
        - 在表格中添加"去审核"操作按钮
        - 点击按钮跳转到审核页面 /reviews/business/:paperId
        - 仅业务审核人员可访问（通过权限中间件控制）
        - 实现快速筛选功能（今日待审、逾期待审等）

#### 7.6 创建政工审核列表页
- [x] 7.6.1 创建政工审核列表页组件基础结构
     【目标对象】`src/frontend/src/pages/review/PoliticalReviewListPage.tsx`
     【修改目的】展示待政工审核的论文列表基础结构
     【修改方式】创建新文件，使用 React 和 Ant Design 组件
     【相关依赖】reviewService, useAuth
     【修改内容】
        - 导入必要的依赖和类型
        - 定义 PoliticalReviewListPage 函数组件
        - 定义状态：papers（待审核论文列表）、loading、searchParams
        - 使用 useEffect 在组件挂载时获取待审核列表

- [x] 7.6.2 实现待审核论文列表展示
     【目标对象】`src/frontend/src/pages/review/PoliticalReviewListPage.tsx`
     【修改目的】展示待政工审核的论文列表
     【修改方式】在 PoliticalReviewListPage 组件中添加列表展示逻辑
     【相关依赖】reviewService, Ant Design Table
     【修改内容】
        - 调用 reviewService.getPendingPoliticalReviews 获取待政工审核列表
        - 使用 Ant Design Table 展示论文列表，列包括：论文ID、标题、提交人、提交时间、距提交天数、状态、业务审核意见
        - 使用 useEffect 定时刷新列表（每5分钟）
        - 实现搜索和筛选功能（按标题、提交人、提交时间）
        - 根据距提交天数显示不同颜色（超过2天显示红色）

- [x] 7.6.3 实现去审核功能
     【目标对象】`src/frontend/src/pages/review/PoliticalReviewListPage.tsx`
     【修改目的】实现跳转到审核页面功能
     【修改方式】在 PoliticalReviewListPage 组件中添加跳转逻辑
     【相关依赖】useNavigate
     【修改内容】
        - 在表格中添加"去审核"操作按钮
        - 点击按钮跳转到审核页面 /reviews/political/:paperId
        - 仅政工审核人员可访问（通过权限中间件控制）
        - 实现快速筛选功能（今日待审、逾期待审等）

#### 7.7 创建审核页面
 - [x] 7.7.1 创建审核页面组件基础结构
     【目标对象】`src/frontend/src/pages/review/ReviewPage.tsx`
     【修改目的】实现审核操作基础结构
     【修改方式】创建新文件，使用 React Hook Form 管理表单
     【相关依赖】reviewService, paperService
     【修改内容】
        - 导入必要的依赖和类型
        - 定义 ReviewPage 函数组件
        - 使用 useParams 获取论文ID和审核类型
        - 使用 useForm 创建表单实例，定义表单字段：result（通过/驳回）、comment（审核意见）
        - 定义状态：paper（论文详情）、reviewLogs（历史审核记录）、loading

 - [x] 7.7.2 展示论文信息和附件
     【目标对象】`src/frontend/src/pages/review/ReviewPage.tsx`
     【修改目的】展示论文完整信息和附件列表
     【修改方式】在 ReviewPage 组件中添加信息展示逻辑
     【相关依赖】paperService
     【修改内容】
        - 使用 useEffect 调用 getPaper 方法获取论文详情
        - 展示论文基础信息（标题、摘要、期刊、DOI、影响因子、出版日期等）
        - 展示作者列表
        - 展示课题信息
        - 展示附件列表，支持预览和下载
        - 根据审核类型显示不同内容（业务审核显示学术规范信息，政工审核显示政治合规信息）

 - [x] 7.7.3 实现审核表单和提交功能
     【目标对象】`src/frontend/src/pages/review/ReviewPage.tsx`
     【修改目的】实现审核表单和提交功能
     【修改方式】在 ReviewPage 组件中添加审核表单和提交逻辑
     【相关依赖】React Hook Form, reviewService
     【修改内容】
        - 使用 Ant Design Radio 实现审核结果选择（通过/驳回）
        - 使用 Ant Design TextArea 实现审核意见输入框（驳回时必填）
        - 设置表单校验规则
        - 实现提交审核按钮，调用 reviewService.businessReview 或 reviewService.politicalReview（根据审核类型）
        - 实现取消按钮，返回审核列表页
        - 显示提交中的加载状态

 - [x] 7.7.4 展示历史审核记录
     【目标对象】`src/frontend/src/pages/review/ReviewPage.tsx`
     【修改目的】展示历史审核记录
     【修改方式】在 ReviewPage 组件中添加历史审核记录展示逻辑
     【相关依赖】reviewService
     【修改内容】
        - 使用 useEffect 调用 getReviewLogs 方法获取审核记录
        - 使用 Timeline 组件展示历史审核记录
        - 显示审核类型、审核结果、审核意见、审核人、审核时间
        - 根据审核结果显示不同颜色的标签

### 8. 前端业务组件实现

#### 8.1 创建论文表单组件
 - [x] 8.1.1 创建论文表单组件基础结构
     【目标对象】`src/frontend/src/components/paper/PaperForm.tsx`
     【修改目的】封装论文表单，供录入页和编辑页复用
     【修改方式】创建新文件，使用 React Hook Form 实现表单
     【相关依赖】React Hook Form, Ant Design Form
     【修改内容】
        - 定义 PaperFormProps 接口：mode（create/edit/view）、initialData、onSubmit、onCancel
        - 定义 PaperForm 函数组件，接收 props
        - 使用 useForm 创建表单实例，接收 initialData 初始化表单
        - 使用 useEffect 监听 initialData 变化，重置表单

 - [x] 8.1.2 实现论文基础字段表单
     【目标对象】`src/frontend/src/components/paper/PaperForm.tsx`
     【修改目的】实现论文基础字段表单和校验
     【修改方式】在 PaperForm 组件中添加表单字段
     【相关依赖】Ant Design Form
     【修改内容】
        - 实现标题输入框（必填，最大长度500，禁用/只读根据mode）
        - 实现摘要输入框（必填，多行文本框，最大长度2000，禁用/只读根据mode）
        - 实现期刊选择器（下拉框，支持搜索，禁用/只读根据mode）
        - 实现 DOI 输入框（选填，格式校验正则表达式，禁用/只读根据mode）
        - 实现影响因子输入框（选填，数字类型，保留3位小数，禁用/只读根据mode）
        - 实现出版日期选择器（选填，日期格式，禁用/只读根据mode）
        - 设置表单校验规则

 - [x] 8.1.3 暴露表单实例和回调函数
     【目标对象】`src/frontend/src/components/paper/PaperForm.tsx`
     【修改目的】暴露表单实例和回调函数供父组件调用
     【修改方式】在 PaperForm 组件中添加暴露逻辑
     【相关依赖】React Hook Form
     【修改内容】
        - 使用 useImperativeHandle 暴露表单实例和方法（submit、reset、validate）
        - 实现表单提交逻辑，调用 onSubmit 回调函数
        - 实现取消按钮，调用 onCancel 回调函数
        - 根据mode显示不同的操作按钮（创建/编辑显示保存按钮，查看模式不显示）
        - 导出 PaperForm 组件和 PaperFormRef 类型

#### 8.2 创建作者列表组件
 - [x] 8.2.1 创建作者列表组件基础结构
     【目标对象】`src/frontend/src/components/paper/AuthorList.tsx`
     【修改目的】管理论文作者列表（添加、删除、调整顺序）
     【修改方式】创建新文件，使用 React 和 Ant Design 实现
     【相关依赖】authorService
     【修改内容】
        - 定义 AuthorListProps 接口：authors（作者列表）、onChange（作者列表变化回调）、disabled（是否禁用）
        - 定义 AuthorList 函数组件，接收 props
        - 定义状态：isAddingAuthor（是否正在添加作者）、users（人员库用户列表）

 - [x] 8.2.2 实现作者列表展示和操作
     【目标对象】`src/frontend/src/components/paper/AuthorList.tsx`
     【修改目的】实现作者列表展示和操作功能
     【修改方式】在 AuthorList 组件中添加列表展示和操作逻辑
     【相关依赖】Ant Design List, Table
     【修改内容】
        - 使用 Ant Design Table 展示作者列表，列包括：排名、姓名、作者类型、单位、操作
        - 实现"添加作者"按钮（disabled时禁用）
        - 实现"删除"按钮（disabled时禁用）
        - 实现"上移"/"下移"按钮调整作者顺序（disabled时禁用）
        - 根据排名显示序号
        - 实现作者类型显示（第一作者、共同第一作者、通讯作者、普通作者）

 - [x] 8.2.3 实现作者添加和编辑功能
     【目标对象】`src/frontend/src/components/paper/AuthorList.tsx`
     【修改目的】实现作者添加和编辑功能
     【修改方式】在 AuthorList 组件中添加添加和编辑逻辑
     【相关依赖】Ant Design Modal, Form
     【修改内容】
        - 实现添加作者弹窗，包含表单：姓名（支持从人员库选择或手动输入）、作者类型（下拉选择）、排名（自动生成）、单位（选填）
        - 实现从人员库搜索用户功能（调用 authorService.searchUsers）
        - 实现作者类型选择器（第一作者、共同第一作者、通讯作者、普通作者）
        - 实现作者类型互斥性校验（第一作者、共同第一作者、通讯作者不能重复）
        - 实现作者排名自动调整（添加后自动调整排名）
        - 调用 onChange 回调函数更新父组件的作者列表

#### 8.3 创建课题选择器组件
 - [x] 8.3.1 创建课题选择器组件基础结构
     【目标对象】`src/frontend/src/components/paper/ProjectSelector.tsx`
     【修改目的】管理论文关联的课题
     【修改方式】创建新文件，使用 React 和 Ant Design 实现
     【相关依赖】projectService
     【修改内容】
        - 定义 ProjectSelectorProps 接口：selectedProjects（已选课题列表）、onChange（课题列表变化回调）、disabled（是否禁用）
        - 定义 ProjectSelector 函数组件，接收 props
        - 定义状态：projects（课题列表）、searchKeyword（搜索关键字）、isAddingProject（是否正在添加课题）

 - [x] 8.3.2 实现课题搜索和选择功能
     【目标对象】`src/frontend/src/components/paper/ProjectSelector.tsx`
     【修改目的】实现课题搜索和选择功能
     【修改方式】在 ProjectSelector 组件中添加搜索和选择逻辑
     【相关依赖】Ant Design Select, Modal
     【修改内容】
        - 实现课题搜索框，输入关键字后调用 projectService.searchProjects 搜索课题
        - 实现课题选择下拉框，支持多选（disabled时禁用）
        - 显示已选课题列表，使用 Tag 组件展示
        - 实现"删除课题"按钮，点击删除该课题（disabled时禁用）
        - 调用 onChange 回调函数更新父组件的课题列表

 - [x] 8.3.3 实现手动录入新课题功能
     【目标对象】`src/frontend/src/components/paper/ProjectSelector.tsx`
     【修改目的】支持手动录入新课题
     【修改方式】在 ProjectSelector 组件中添加手动录入逻辑
     【相关依赖】Ant Design Modal, Form
     【修改内容】
        - 实现"添加新课题"按钮（disabled时禁用）
        - 实现添加新课题弹窗，包含表单：课题名称（必填）、课题编号（必填）、项目类型（纵向/横向）、来源、级别（国家级/省部级/市级）
        - 实现表单校验
        - 调用 projectService.createProject 创建新课题
        - 创建成功后将新课题添加到已选列表
        - 显示课题详细信息（课题名称、编号、项目类型、级别）

#### 8.4 创建审核表单组件
 - [x] 8.4.1 创建审核表单组件基础结构
     【目标对象】`src/frontend/src/components/review/ReviewForm.tsx`
     【修改目的】封装审核表单，供审核页面复用
     【修改方式】创建新文件，使用 React Hook Form 实现表单
     【相关依赖】React Hook Form
     【修改内容】
        - 定义 ReviewFormProps 接口：reviewType（审核类型）、onSubmit、onCancel
        - 定义 ReviewForm 函数组件，接收 props
        - 使用 useForm 创建表单实例

 - [x] 8.4.2 实现审核表单和校验
     【目标对象】`src/frontend/src/components/review/ReviewForm.tsx`
     【修改目的】实现审核表单和校验逻辑
     【修改方式】在 ReviewForm 组件中添加表单和校验逻辑
     【相关依赖】Ant Design Radio, TextArea, Form
     【修改内容】
        - 使用 Ant Design Radio 实现"通过"和"驳回"单选按钮
        - 使用 Ant Design TextArea 实现审核意见输入框
        - 设置审核意见输入框的校验规则（驳回时必填，最多500字）
        - 使用 useEffect 监听审核结果变化，驳回时聚焦到审核意见输入框
        - 实现表单提交逻辑，调用 onSubmit 回调函数
        - 实现取消按钮，调用 onCancel 回调函数
        - 根据审核类型显示不同的提示信息（业务审核提示关注学术规范性，政工审核提示关注政治合规性）
        - 导出 ReviewForm 组件

#### 8.5 创建文件上传组件
 - [x] 8.5.1 创建文件上传组件基础结构
     【目标对象】`src/frontend/src/components/common/FileUpload.tsx`
     【修改目的】封装文件上传功能
     【修改方式】创建新文件，使用 Ant Design Upload 实现
     【相关依赖】fileService, Ant Design Upload
     【修改内容】
        - 定义 FileUploadProps 接口：accept（允许的文件类型）、maxSize（最大文件大小）、maxCount（最大文件数量）、onChange（文件列表变化回调）
        - 定义 FileUpload 函数组件，接收 props
        - 定义状态：fileList（文件列表）、uploading（上传中）

 - [x] 8.5.2 实现文件上传功能
     【目标对象】`src/frontend/src/components/common/FileUpload.tsx`
     【修改目的】实现文件选择、上传和校验功能
     【修改方式】在 FileUpload 组件中添加上传逻辑
     【相关依赖】Ant Design Upload
     【修改内容】
        - 使用 Ant Design Upload 组件实现文件选择和上传
        - 设置 accept 属性限制文件类型（默认接受所有文件）
        - 设置 beforeUpload 回调函数，校验文件大小（超过maxSize则拒绝）
        - 调用 fileService.uploadFile 上传文件
        - 显示上传进度条
        - 显示上传状态（成功、失败）

 - [x] 8.5.3 实现已上传文件展示和删除功能
     【目标对象】`src/frontend/src/components/common/FileUpload.tsx`
     【修改目的】展示已上传文件列表和实现删除功能
     【修改方式】在 FileUpload 组件中添加文件展示和删除逻辑
     【相关依赖】Ant Design List, Upload
     【修改内容】
        - 使用 Ant Design Upload.List 展示已上传文件列表
        - 显示文件名称、文件大小、上传状态
        - 实现"删除"按钮，点击后从fileList中移除该文件
        - 限制最大文件数量（超过maxCount则不能再上传）
        - 调用 onChange 回调函数更新父组件的文件列表

#### 8.6 创建批量导入表单组件
 - [x] 8.6.1 创建批量导入表单组件基础结构
      【目标对象】`src/frontend/src/components/paper/BatchImportForm.tsx`
      【修改目的】封装批量导入表单
      【修改方式】创建新文件，使用 React 和 Ant Design 实现
      【相关依赖】paperService
      【修改内容】
         - 定义 BatchImportFormProps 接口：onImportComplete（导入完成回调）
         - 定义 BatchImportForm 函数组件，接收 props
         - 定义状态：file（上传的文件）、importResult（导入结果）、isImporting

 - [x] 8.6.2 实现导入模板下载和文件上传
      【目标对象】`src/frontend/src/components/paper/BatchImportForm.tsx`
      【修改目的】实现导入模板下载和Excel文件上传功能
      【修改方式】在 BatchImportForm 组件中添加下载和上传逻辑
      【相关依赖】paperService, Ant Design Upload
      【修改内容】
         - 实现下载导入模板按钮，调用 paperService.downloadImportTemplate
         - 使用 Ant Design Upload 组件上传 Excel 文件
         - 限制文件类型为 .xlsx 和 .xls
         - 限制文件大小为 10MB
         - 实现文件上传前校验

 - [x] 8.6.3 实现批量导入和结果展示
      【目标对象】`src/frontend/src/components/paper/BatchImportForm.tsx`
      【修改目的】实现批量导入和结果展示功能
      【修改方式】在 BatchImportForm 组件中添加导入和结果展示逻辑
      【相关依赖】paperService, Ant Design Table, Alert
      【修改内容】
         - 实现开始导入按钮，调用 paperService.batchImport
         - 显示导入中的加载状态
         - 使用 Ant Design Alert 展示导入结果（成功数量、失败数量）
         - 使用 Ant Design Table 展示错误详情（错误行号、错误原因、错误字段）
         - 支持修正后重新导入
         - 调用 onImportComplete 回调函数通知父组件导入完成

### 9. 前端路由和菜单配置

#### 9.1 配置论文管理相关路由
 - [x] 9.1.1 在路由配置中新增论文管理相关路由
     【目标对象】`src/frontend/src/router/index.tsx`
     【修改目的】注册论文管理和审核相关路由
     【修改方式】在路由配置中新增路由
     【相关依赖】各页面组件
     【修改内容】
        - 导入所有论文和审核相关页面组件
        - 在路由配置中添加以下路由：
          * /papers - 论文列表页（PaperListPage）
          * /papers/create - 论文录入页（PaperCreatePage）
          * /papers/:id - 论文详情页（PaperDetailPage）
          * /papers/:id/edit - 论文编辑页（PaperCreatePage，带edit模式）
          * /papers/batch-import - 批量导入页（BatchImportPage）
          * /reviews/business - 业务审核列表页（BusinessReviewListPage）
          * /reviews/business/:paperId - 业务审核页面（ReviewPage，类型为business）
          * /reviews/political - 政工审核列表页（PoliticalReviewListPage）
          * /reviews/political/:paperId - 政工审核页面（ReviewPage，类型为political）

 - [x] 9.1.2 为路由添加权限控制
     【目标对象】`src/frontend/src/router/index.tsx`
     【修改目的】为论文管理和审核相关路由添加权限控制
     【修改方式】在路由配置中添加权限守卫
     【相关依赖】useAuth
     【修改内容】
        - 为 /papers/create 路由添加 paper:create 权限检查
        - 为 /papers/:id/edit 路由添加 paper:edit 权限检查
        - 为 /papers/:id/delete 路由添加 paper:delete 权限检查
        - 为 /papers/batch-import 路由添加 paper:import 权限检查
        - 为 /reviews/business/* 路由添加 review:business 权限检查
        - 为 /reviews/political/* 路由添加 review:political 权限检查
        - 无权限访问时跳转到403页面或显示权限错误提示
        - 使用路由守卫（RouteGuard组件或高阶组件）统一处理权限检查

#### 9.2 更新Layout菜单项
 - [x] 9.2.1 在Layout中新增论文管理菜单组
     【目标对象】`src/frontend/src/components/Layout/Layout.tsx`
     【修改目的】在菜单中新增论文管理相关菜单项
     【修改方式】在菜单配置中新增菜单组
     【相关依赖】useAuth
     【修改内容】
        - 在菜单配置中添加"论文管理"菜单组
        - 在"论文管理"菜单组下添加子菜单项：
          * "论文列表" - 链接到 /papers（所有用户可访问）
          * "论文录入" - 链接到 /papers/create（所有用户可访问）
          * "批量导入" - 链接到 /papers/batch-import（所有用户可访问）
        - 使用 hasPermission 方法实现权限控制（如有需要）

 - [x] 9.2.2 在Layout中新增审核管理菜单组
     【目标对象】`src/frontend/src/components/Layout/Layout.tsx`
     【修改目的】在菜单中新增审核管理相关菜单项
     【修改方式】在菜单配置中新增菜单组
     【相关依赖】useAuth
     【修改内容】
        - 在菜单配置中添加"审核管理"菜单组
        - 在"审核管理"菜单组下添加子菜单项：
          * "业务审核" - 链接到 /reviews/business（仅业务审核人员可访问，使用 hasPermission('review:business')）
          * "政工审核" - 链接到 /reviews/political（仅政工审核人员可访问，使用 hasPermission('review:political')）
        - 使用 hasPermission 方法实现权限控制，无权限的用户不显示对应的菜单项

- [x] 9.2.3 在Layout中新增课题管理菜单项（可选）
     【目标对象】`src/frontend/src/components/Layout/Layout.tsx`
     【修改目的】在菜单中新增课题管理菜单项
     【修改方式】在菜单配置中新增菜单项
     【相关依赖】useAuth
     【修改内容】
        - 在"论文管理"菜单组下添加"课题管理"子菜单项
        - "课题管理"链接到 /projects（管理员和课题负责人可访问，使用 hasPermission('project:manage')）
        - 使用 hasPermission 方法实现权限控制，无权限的用户不显示对应的菜单项

### 10. 后端定时任务实现（可选）

#### 10.1 实现审核时限提醒任务
 - [x] 10.1.1 在 main.go 中启动审核时限提醒定时任务
      【目标对象】`src/backend/cmd/server/main.go`
      【修改目的】启动定时任务检查待审核论文，发送逾期提醒
      【修改方式】在 main 函数中启动定时任务
      【相关依赖】review_service, notification_service
      【修改内容】
         - 在 main 函数启动服务前，使用 go 关键字启动一个 goroutine
         - 使用 time.Ticker 创建定时器，每小时执行一次
         - 在 goroutine 中调用 reviewService.SendReviewReminderForOverdue 方法
         - 添加日志记录定时任务的执行情况

 - [x] 10.1.2 实现检查逾期审核的方法
      【目标对象】`src/backend/internal/service/review_service.go`
      【修改目的】查询待审核且超过2个工作日的论文，发送提醒通知
      【修改方式】在 ReviewService 中添加方法
      【相关依赖】review_repository, notification_service
      【修改内容】
         - 实现 SendReviewReminderForOverdue 方法：
           * 查询所有状态为"待业务审核"和"待政工审核"的论文
           * 计算距离提交时间的工作日数
           * 筛选出超过2个工作日且未发送过提醒的论文
           * 对每篇论文调用 notification_service 发送提醒通知给对应的审核人员
           * 记录已发送提醒的标记（可添加字段或在操作日志中记录）
         - 添加日志记录提醒发送情况

### 11. 测试和验证

#### 11.1 后端接口测试
- [x] 11.1.1 测试论文CRUD接口
     【目标对象】所有论文相关Handler接口
     【修改目的】验证论文CRUD接口功能正确性
     【修改方式】使用 Postman 或 cURL 测试各接口
     【相关依赖】无
     【修改内容】
        - 测试 POST /api/papers 创建论文接口（正常数据、缺少必填字段、格式错误）
        - 测试 GET /api/papers/:id 获取论文详情接口（正常ID、不存在的ID）
        - 测试 GET /api/papers 分页查询接口（正常分页、带查询条件、空结果）
        - 测试 PUT /api/papers/:id 更新论文接口（草稿状态更新、非草稿状态更新）
        - 测试 DELETE /api/papers/:id 删除论文接口（草稿状态删除、非草稿状态删除）

- [x] 11.1.2 测试论文提交和校验接口
     【目标对象】论文提交和校验相关Handler接口
     【修改目的】验证论文提交和校验接口功能正确性
     【修改方式】使用 Postman 或 cURL 测试各接口
     【相关依赖】无
     【修改内容】
        - 测试 POST /api/papers/:id/submit 提交审核接口（草稿状态提交、非草稿状态提交、重复提交）
        - 测试 POST /api/papers/:id/save-draft 保存草稿接口
        - 测试 POST /api/papers/check-duplicate 检查重复接口（正常数据、重复数据）
        - 测试 GET /api/papers/my 获取我的论文接口

- [x] 11.1.3 测试审核接口
     【目标对象】审核相关Handler接口
     【修改目的】验证审核接口功能正确性
     【修改方式】使用 Postman 或 cURL 测试各接口
     【相关依赖】无
     【修改内容】
        - 测试 POST /api/reviews/business/:paperId 业务审核接口（通过、驳回、无权限）
        - 测试 POST /api/reviews/political/:paperId 政工审核接口（通过、驳回、无权限）
        - 测试 GET /api/reviews/:paperId/logs 获取审核记录接口
        - 测试 GET /api/reviews/pending/business 获取待业务审核列表接口
        - 测试 GET /api/reviews/pending/political 获取待政工审核列表接口
        - 测试 GET /api/reviews/my 获取我的审核记录接口

- [ ] 11.1.4 测试文件上传和课题、期刊接口
     【目标对象】文件上传、课题、期刊相关Handler接口
     【修改目的】验证文件上传、课题、期刊接口功能正确性
     【修改方式】使用 Postman 或 cURL 测试各接口
     【相关依赖】无
     【修改内容】
        - 测试 POST /api/files/upload 文件上传接口（正常文件、超大文件、错误格式）
        - 测试 GET /api/files/:id 获取文件信息接口
        - 测试 GET /api/files/:id/download 下载文件接口
        - 测试 DELETE /api/files/:id 删除文件接口
        - 测试课题CRUD接口
        - 测试期刊CRUD接口
        - 测试 POST /api/papers/batch-import 批量导入接口（正常Excel、错误Excel）

- [ ] 11.1.5 测试权限控制
     【目标对象】所有Handler接口
     【修改目的】验证权限控制功能正确性
     【修改方式】使用不同角色的用户测试接口
     【相关依赖】无
     【修改内容】
        - 测试无权限访问接口应返回403错误
        - 测试不同角色对论文的访问权限（创建、编辑、删除）
        - 测试不同角色对审核的访问权限（业务审核、政工审核）
        - 测试不同角色对课题和期刊的管理权限

#### 11.2 前端功能测试
- [ ] 11.2.1 测试论文录入功能
     【目标对象】PaperCreatePage 组件
     【修改目的】验证论文录入功能正确性
     【修改方式】手动测试页面功能
     【相关依赖】无
     【修改内容】
        - 测试论文基础字段录入（标题、摘要、期刊、DOI、影响因子、出版日期）
        - 测试字段校验（必填字段、格式校验、长度限制）
        - 测试作者列表管理（添加、删除、调整顺序、类型互斥性校验）
        - 测试课题选择和手动录入
        - 测试文件上传（正常文件、超大文件、错误格式）
        - 测试重复校验功能
        - 测试保存草稿功能
        - 测试提交审核功能

- [ ] 11.2.2 测试论文列表和详情功能
     【目标对象】PaperListPage 和 PaperDetailPage 组件
     【修改目的】验证论文列表和详情功能正确性
     【修改方式】手动测试页面功能
     【相关依赖】无
     【修改内容】
        - 测试论文列表查询（搜索、分页、排序）
        - 测试论文列表操作（查看、编辑、删除、提交审核）
        - 测试论文详情展示（基础信息、作者、课题、附件、审核记录）
        - 测试根据状态显示不同操作按钮
        - 测试附件下载功能
        - 测试权限控制（不同角色看到不同操作）

- [ ] 11.2.3 测试批量导入功能
     【目标对象】BatchImportPage 组件
     【修改目的】验证批量导入功能正确性
     【修改方式】手动测试页面功能
     【相关依赖】无
     【修改内容】
        - 测试导入模板下载
        - 测试Excel文件上传（正常数据、错误数据、超大文件）
        - 测试批量导入功能（全部成功、部分失败、全部失败）
        - 测试导入结果展示（成功数量、失败数量、错误详情）
        - 测试修正后重新导入

- [ ] 11.2.4 测试审核功能
     【目标对象】BusinessReviewListPage、PoliticalReviewListPage、ReviewPage 组件
     【修改目的】验证审核功能正确性
     【修改方式】手动测试页面功能
     【相关依赖】无
     【修改内容】
        - 测试待审核列表展示（业务审核、政工审核）
        - 测试审核页面展示（论文信息、附件、历史审核记录）
        - 测试审核通过功能（业务审核通过、政工审核通过）
        - 测试审核驳回功能（必填校验、驳回原因）
        - 测试审核驳回后重新提交功能
        - 测试审核时限提醒（如实现）
        - 测试权限控制（业务审核人员、政工审核人员）

#### 11.3 端到端流程测试
- [ ] 11.3.1 测试论文录入到归档的完整流程
     【目标对象】完整业务流程
     【修改目的】验证论文从录入到归档的完整流程
     【修改方式】模拟真实用户操作，走通完整流程
     【相关依赖】无
     【修改内容】
        - 流程1：普通用户录入论文 → 提交审核 → 业务审核通过 → 政工审核通过 → 归档
          * 验证每个步骤的状态变化
          * 验证通知是否正确发送
          * 验证审核记录是否正确记录
        - 流程2：普通用户录入论文 → 提交审核 → 业务审核驳回 → 修改后重新提交 → 审核通过
          * 验证驳回后状态是否重置为草稿
          * 验证修改后重新提交流程
        - 流程3：管理员批量导入论文 → 审核通过 → 归档
          * 验证批量导入的论文正确入库
          * 验证批量导入后的审核流程

- [ ] 11.3.2 测试审核时限提醒和通知功能
     【目标对象】审核时限提醒和通知功能
     【修改目的】验证审核时限提醒和通知功能
     【修改方式】模拟逾期场景，测试通知功能
     【相关依赖】无
     【修改内容】
        - 流程4：测试审核时限提醒功能
          * 创建待审核论文，等待超过2个工作日
          * 验证审核人员是否收到提醒通知
          * 验证提醒通知的内容是否正确
        - 测试通知到达率（系统消息、邮件通知）
        - 测试不同审核人员的通知分配

- [ ] 11.3.3 验证操作日志和数据一致性
     【目标对象】操作日志和数据一致性
     【修改目的】验证所有操作日志是否正确记录，数据一致性是否保持
     【修改方式】检查数据库和日志
     【相关依赖】无
     【修改内容】
        - 验证所有操作是否正确记录到 operation_logs 表
        - 验证审核记录是否正确记录到 review_logs 表
        - 验证归档记录是否正确记录到 archive 表
        - 验证数据关联是否正确（论文-作者、论文-课题）
        - 验证软删除是否正确执行
        - 验证并发场景下的数据一致性
