## 实施

- [x] 1.1 扩展论文实体（Paper Entity）
     【目标对象】`src/backend/internal/model/entity/paper.go`
     【修改目的】添加缺失的字段以支持完整查询功能
     【修改方式】添加新字段并更新GORM标签
     【修改内容】
        - 添加 DOI 字段（string, varchar(100)）
        - 添加 PubMedID 字段（string, varchar(50)）
        - 添加 ISSN 字段（string, varchar(20)）
        - 添加 ImpactFactor 字段（float64）
        - 添加 Partition 字段（分区：string）
        - 添加 IsSCI、IsEI、IsCI、IsDI 字段（boolean）
        - 添加 IsCore 字段（中文核心：boolean）
        - 添加 CitationCount 字段（引用次数：int）
        - 添加 Language 字段（语言：string）
        - 添加 JournalName 字段（期刊名称，用于非关联查询）
        - 添加 Volume、Issue、StartPage、EndPage 字段
        - 添加 AuthorType 字段（作者类型：string，用于区分第一作者/共同第一作者/通讯作者）
        - 添加 IsFirstAuthor、IsCoFirstAuthor、IsCorrespondingAuthor 字段（boolean）

- [x] 1.2 扩展归档实体（Archive Entity）
     【目标对象】`src/backend/internal/model/entity/archive.go`
     【修改目的】支持归档状态管理
     【修改方式】添加新字段
     【修改内容】
        - 添加 Status 字段（string, 公开/隐藏）
        - 添加 IsHidden 字段（boolean, 隐藏状态）

- [x] 2.1 创建搜索请求DTO
     【目标对象】`src/backend/internal/model/dto/request/search_request.go`
     【修改目的】定义搜索请求参数
     【修改方式】新建文件
     【修改内容】
        - SearchRequest 结构体：包含所有查询条件字段
        - QueryCondition 结构体：单个查询条件
        - QueryGroup 结构体：查询条件组，支持嵌套组合
        - LogicType 枚举：AND（且）/OR（或）/NOT（非）逻辑类型
        - AuthorType 枚举：第一作者/共同第一作者/通讯作者
        - SortField 枚举：排序字段（出版日期、影响因子、引用次数）
        - SortOrder 枚举：排序方向（升序/降序）
        - Pagination 结构体：分页参数（page, pageSize 默认20）

- [x] 2.2 创建搜索响应DTO
     【目标对象】`src/backend/internal/model/dto/response/search_response.go`
     【修改目的】定义搜索响应格式
     【修改方式】新建文件
     【修改内容】
         - SearchResponse 结构体：包含分页信息和结果列表
         - PaperSearchResult 结构体：论文搜索结果项，包含论文ID、标题、作者列表、期刊、出版日期、影响因子、分区、审核状态、作者类型
         - AuthorInfo 结构体：作者信息（姓名、单位、作者类型）
         - ProjectInfo 结构体：课题信息（课题编号、项目类型）
         - PaginationInfo 结构体：分页信息（当前页、每页条数、总条数、总页数）
         - ApiResponse 结构体：统一响应格式

- [x] 2.3 扩展论文仓储（Paper Repository）
     【目标对象】`src/backend/internal/repository/paper_repository.go`
     【修改目的】支持复杂查询
     【修改方式】添加新查询方法
     【修改内容】
         - 添加 AdvancedSearch 方法：支持多维度组合查询，传入查询条件组
         - 添加 FindByDOI 方法：DOI精确查询
         - 添加 FindByPubMedID 方法：PubMedID精确查询
         - 添加 FindByISSN 方法：ISSN精确查询
         - 添加 FindByAuthorName 方法：作者名查询
         - 添加 FindByAuthorType 方法：按作者类型查询（第一作者/共同第一作者/通讯作者）
         - 添加 FindByProjectCode 方法：课题编号查询
         - 添加 ListByYear 方法：按年份查询
         - 添加 ListByType 方法：按收录类型查询
         - 添加 Count 方法：统计总数（用于分页）

- [x] 2.4 创建搜索服务（Search Service）
     【目标对象】`src/backend/internal/service/search_service.go`
     【修改目的】实现搜索业务逻辑
     【修改方式】新建文件
     【修改内容】
        - SearchService 结构体
        - AdvancedSearch 方法：多维度组合查询
        - buildQueryConditions 方法：构建查询条件，支持且(AND)/或(OR)/非(NOT)逻辑组合
        - applySorting 方法：应用排序（按出版日期、影响因子、引用次数升序/降序）
        - filterByPermission 方法：权限过滤（仅审核通过的论文可公开查询）
        - filterByAuthorType 方法：按作者类型过滤（第一作者/共同第一作者/通讯作者）

- [x] 2.5 创建搜索处理器（Search Handler）
     【目标对象】`src/backend/internal/handler/search_handler.go`
     【修改目的】处理搜索HTTP请求
     【修改方式】新建文件
     【修改内容】
        - SearchHandler 结构体
        - Search 方法：处理搜索请求
        - GetPaperDetail 方法：获取论文详情

- [x] 2.5.1 扩展审核服务（审核通过自动归档）
     【目标对象】`src/backend/internal/service/paper_service.go`
     【修改目的】实现审核通过后自动归档功能
     【修改方式】修改审核通过方法
     【修改内容】
        - 在 ApprovePaper 方法中，审核通过后自动调用归档服务
        - 调用 archive_service.CreateArchive 方法创建归档记录
        - 生成归档编号（年份+论文ID+随机3位数字）
        - 设置归档状态为"公开"
        - 处理归档失败的情况（记录日志但不影响审核流程）

- [x] 2.6 扩展归档仓储（Archive Repository）
     【目标对象】`src/backend/internal/repository/archive_repository.go`
     【修改目的】支持分类查询
     【修改方式】添加新查询方法
     【修改内容】
        - 添加 ListByYear 方法：按年份查询
        - 添加 ListByType 方法：按收录类型查询
        - 添加 ListByAuthor 方法：按作者查询
        - 添加 ListByProject 方法：按课题查询
        - 添加 ListArchived 方法：查询已归档论文
        - 添加 UpdateStatus 方法：更新归档状态

- [x] 2.7 扩展归档服务（Archive Service）
     【目标对象】`src/backend/internal/service/archive_service.go`
     【修改目的】实现完整的归档管理功能
     【修改方式】扩展现有方法
     【修改内容】
        - 修正 generateArchiveNumber 方法：年份+论文ID+随机3位数字
        - 添加 GetArchivedPapers 方法：获取已归档论文列表
        - 添加 GetArchivedPapersByYear 方法：按年份获取
        - 添加 GetArchivedPapersByType 方法：按收录类型获取
        - 添加 GetArchivedPapersByAuthor 方法：按作者获取
        - 添加 GetArchivedPapersByProject 方法：按课题获取
        - 添加 HideArchive 方法：隐藏归档论文
        - 添加 SubmitArchiveModifyRequest 方法：提交归档修改申请

- [x] 2.8 扩展归档处理器（Archive Handler）
     【目标对象】`src/backend/internal/handler/archive_handler.go`
     【修改目的】处理归档管理HTTP请求
     【修改方式】扩展现有方法
     【修改内容】
        - 添加 GetArchiveList 方法：获取归档列表（支持分类筛选）
        - 添加 GetArchiveByPaperID 方法：获取单个归档记录
        - 添加 HideArchive 方法：隐藏归档论文
        - 添加 SubmitModifyRequest 方法：提交修改申请

- [x] 2.9 注册搜索和归档路由
     【目标对象】`src/backend/cmd/server/main.go`
     【修改目的】注册新的API路由
     【修改方式】添加路由配置
     【修改内容】
        - 添加搜索API路由：GET /api/search
        - 添加论文详情API：GET /api/papers/:id/detail
        - 添加归档列表API：GET /api/archives
        - 添加归档操作API：PUT /api/archives/:paperId/hide
        - 添加归档修改申请API：POST /api/archives/:paperId/modify

- [x] 3.1 扩展前端类型定义
     【目标对象】`src/frontend/src/types/paper.ts`
     【修改目的】支持搜索和归档相关类型
     【修改方式】扩展现有类型
     【修改内容】
         - 添加 SearchRequest 类型：包含查询条件、逻辑组合、分页、排序参数
         - 添加 QueryCondition 类型：单个查询条件
         - 添加 LogicType 类型：AND/OR/NOT 逻辑类型
         - 添加 AuthorTypeFilter 类型：第一作者/共同第一作者/通讯作者/全部
         - 添加 PaperSearchResult 类型：搜索结果项
         - 添加 PaginationInfo 类型：分页信息
         - 添加 ArchiveInfo 类型：归档信息
         - 添加 SortField 类型：排序字段
         - 添加 SortOrder 类型：排序方向

- [x] 3.2 创建搜索服务（前端）
     【目标对象】`src/frontend/src/services/searchService.ts`
     【修改目的】提供搜索API调用
     【修改方式】新建文件
     【修改内容】
        - advancedSearch 函数：多维度组合查询
        - getPaperDetail 函数：获取论文详情

- [x] 3.3 创建归档服务（前端）
     【目标对象】`src/frontend/src/services/archiveService.ts`
     【修改目的】提供归档管理API调用
     【修改方式】新建文件
     【修改内容】
        - getArchiveList 函数：获取归档列表
        - getArchiveByPaperId 函数：获取归档详情
        - hideArchive 函数：隐藏归档论文
        - submitModifyRequest 函数：提交修改申请

- [x] 3.4 创建搜索页面
     【目标对象】`src/frontend/src/pages/search/SearchPage.tsx`
     【修改目的】多维度组合查询界面
     【修改方式】新建文件
     【修改内容】
         - 搜索表单组件：支持论文ID、标题、作者、期刊、日期范围、DOI、影响因子范围、课题等条件
         - 逻辑组合组件：支持且(AND)/或(OR)/非(NOT)逻辑组合多个查询条件
         - 作者类型筛选：支持下拉选择"全部"/"第一作者"/"共同第一作者"/"通讯作者"
         - 结果表格组件：展示论文ID、标题、作者列表、期刊、出版日期、影响因子、分区、审核状态
         - 排序功能：支持下拉选择按出版日期、影响因子、引用次数排序，支持升序/降序
         - 分页组件：每页20条，超过1000条正常分页，显示总页数和当前页
         - 详情查看：点击行查看论文详情（完整信息、作者、课题、附件）
         - 无结果提示：当查询结果为空时显示"未找到符合条件的论文"

- [x] 3.5 创建归档管理页面
     【目标对象】`src/frontend/src/pages/archive/ArchivePage.tsx`
     【修改目的】归档论文管理界面
     【修改方式】新建文件
     【修改内容】
        - 分类筛选组件：按年份、收录类型、课题、作者分类
        - 归档列表表格：展示归档论文
        - 归档详情查看：查看归档论文完整信息
        - 隐藏功能：隐藏归档论文
        - 修改申请功能：提交归档修改申请

- [x] 3.6 更新前端路由配置
     【目标对象】`src/frontend/src/router/index.tsx`
     【修改目的】添加搜索和归档页面路由
     【修改方式】添加新的路由配置
     【修改内容】
        - 添加 /search 路由：搜索页面
        - 添加 /archives 路由：归档管理页面

- [x] 3.7 更新布局菜单
     【目标对象】`src/frontend/src/components/Layout/MainLayout.tsx`
     【修改目的】添加搜索和归档菜单入口
     【修改方式】添加菜单项
     【修改内容】
        - 添加"论文查询"菜单项
        - 添加"归档管理"菜单项
