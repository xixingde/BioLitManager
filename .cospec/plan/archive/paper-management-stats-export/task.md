## 实施

### 后端实现

- [x] 1.1 创建统计仓储层
     【目标对象】`src/backend/internal/repository/stats_repository.go`
     【修改目的】实现统计数据查询
     【修改方式】新建文件，包含基础统计、按作者统计、按课题统计、按单位统计的查询方法
     【相关依赖】`src/backend/internal/model/entity/paper.go`, `src/backend/internal/model/entity/author.go`, `src/backend/internal/model/entity/project.go`
     【修改内容】
     - 创建 StatsRepository 结构体
     - 实现 GetBasicStats() 基础统计方法
     - 实现 GetAuthorStats(authorId) 按作者统计方法
     - 实现 GetProjectStats(projectId) 按课题统计方法
     - 实现 GetDepartmentStats(department) 按单位统计方法

- [x] 1.2 创建统计服务层
     【目标对象】`src/backend/internal/service/stats_service.go`
     【修改目的】实现统计业务逻辑
     【修改方式】新建文件，封装统计业务逻辑
     【相关依赖】`src/backend/internal/repository/stats_repository.go`
     【修改内容】
     - 创建 StatsService 结构体
     - 实现 GetBasicStats() 获取基础指标统计
     - 实现 GetStatsByAuthor(authorId) 按作者统计
     - 实现 GetStatsByProject(projectId) 按课题统计
     - 实现 GetStatsByDepartment(department) 按单位统计
     - 实现 GetYearlyStats() 年度统计
     - 实现 GetJournalStats() 期刊统计

- [x] 1.3 创建统计处理器
     【目标对象】`src/backend/internal/handler/stats_handler.go`
     【修改目的】处理统计 HTTP 请求
     【修改方式】新建文件，处理统计相关 API 请求
     【相关依赖】`src/backend/internal/service/stats_service.go`
     【修改内容】
     - 创建 StatsHandler 结构体
     - 实现 GetBasicStats() 处理基础统计请求
     - 实现 GetAuthorStats() 处理按作者统计请求
     - 实现 GetProjectStats() 处理按课题统计请求
     - 实现 GetDepartmentStats() 处理按单位统计请求

- [x] 1.4 创建导出服务层
     【目标对象】`src/backend/internal/service/export_service.go`
     【修改目的】实现数据导出业务逻辑
     【修改方式】新建文件，包含 Excel、PDF、Word 导出逻辑
     【相关依赖】`src/backend/internal/model/entity/paper.go`, excelize, gofpdf, unioffice 库
     【修改内容】
     - 创建 ExportService 结构体
     - 实现 ExportPapersToExcel(papers, fields) 查询结果导出 Excel
     - 实现 ExportPaperToPDF(paperId) 单篇论文导出 PDF
     - 实现 ExportPaperToWord(paperId) 单篇论文导出 Word
     - 实现 ExportStatsToExcel(stats) 统计结果导出 Excel
     - 实现 ExportStatsToPDF(stats) 统计结果导出 PDF
     - 实现 generateExcelHeader() 生成 Excel 表头
     - 实现 generatePDFContent() 生成 PDF 内容
     - 实现 generateWordContent() 生成 Word 内容
     - 实现 checkExportPermission(userId, paperId) 导出权限校验
     - 实现 filterAccessiblePapers(userId, papers) 数据范围过滤（普通用户仅返回自己提交的论文）

- [x] 1.5 创建导出处理器
     【目标对象】`src/backend/internal/handler/export_handler.go`
     【修改目的】处理导出 HTTP 请求
     【修改方式】新建文件，处理导出相关 API 请求
     【相关依赖】`src/backend/internal/service/export_service.go`
     【修改内容】
     - 创建 ExportHandler 结构体
     - 实现 ExportPapers() 处理查询结果导出请求
     - 实现 ExportPaper() 处理单篇论文导出请求（支持 PDF/Word 格式）
     - 实现 ExportStats() 处理统计结果导出请求
     - 实现 DownloadExportFile() 处理文件下载
     - 实现 GetExportFields() 获取可导出字段列表
     - 添加权限校验中间件：验证用户是否有导出权限
     - 添加数据范围过滤：普通用户仅可导出自己提交的论文

- [x] 1.6 注册统计和导出路由
     【目标对象】`src/backend/cmd/server/main.go`
     【修改目的】注册新的 API 路由
     【修改方式】在现有路由注册代码中添加统计和导出路由
     【相关依赖】`src/backend/internal/handler/stats_handler.go`, `src/backend/internal/handler/export_handler.go`
     【修改内容】
     - 初始化 StatsService 和 StatsHandler
     - 初始化 ExportService 和 ExportHandler
     - 注册 `/api/stats/*` 路由组
     - 注册 `/api/export/*` 路由组
     - 添加权限中间件：stats:view, stats:export, paper:export

### 前端实现

- [x] 2.1 扩展统计类型定义
     【目标对象】`src/frontend/src/types/statistics.ts`
     【修改目的】支持新的统计数据类型
     【修改方式】在现有类型文件中添加新的类型定义
     【相关依赖】无
     【修改内容】
     - 添加 BasicStats 接口（论文总数、年份统计、收录类型统计、期刊统计、平均影响因子、总引用次数、总他引次数）
     - 添加 AuthorStats 接口（作者信息、论文数量、第一作者数量、通讯作者数量、平均影响因子、总引用次数）
     - 添加 ProjectStats 接口（课题信息、论文数量、高影响因子论文数量、SCI 论文数量）
     - 添加 DepartmentStats 接口（单位信息、论文数量、总影响因子、总引用次数）
     - 添加 ChartData 接口（图表数据）
     - 添加 ExportRequest 接口（导出请求参数）
     - 添加 ExportField 接口（导出字段定义）
     - 添加 UserRole 接口（用户角色：admin, user）

- [x] 2.2 创建统计服务
     【目标对象】`src/frontend/src/services/statsService.ts`
     【修改目的】提供统计 API 调用
     【修改方式】新建文件，封装统计 API 调用
     【相关依赖】`src/frontend/src/services/api.ts`
     【修改内容】
     - 实现 getBasicStats() 获取基础统计
     - 实现 getAuthorStats(authorId) 获取作者统计
     - 实现 getProjectStats(projectId) 获取课题统计
     - 实现 getDepartmentStats(department) 获取单位统计
     - 实现 getYearlyStats() 获取年度统计
     - 实现 getJournalStats() 获取期刊统计

- [x] 2.3 创建导出服务
     【目标对象】`src/frontend/src/services/exportService.ts`
     【修改目的】提供导出 API 调用
     【修改方式】新建文件，封装导出 API 调用
     【相关依赖】`src/frontend/src/services/api.ts`
     【修改内容】
     - 实现 exportPapers(params) 导出查询结果（Excel 格式）
     - 实现 exportPaper(paperId, format) 导出单篇论文（支持 PDF、Word 格式）
     - 实现 exportStats(type, format) 导出统计结果（支持 Excel、PDF 格式）
     - 实现 getExportFields() 获取可导出字段列表
     - 实现 downloadFile(url, filename) 下载文件
     - 实现 checkExportPermission(paperId) 检查导出权限
     - 处理权限不足的错误提示

- [x] 2.4 创建导出按钮组件
     【目标对象】`src/frontend/src/components/common/ExportButton.tsx`
     【修改目的】提供统一的导出 UI
     【修改方式】新建文件，创建导出按钮组件
     【相关依赖】`src/frontend/src/services/exportService.ts`, antd 组件
     【修改内容】
     - 创建 ExportButton 组件
     - 支持导出类型：papers, paper, stats
     - 支持导出格式：Excel, PDF, Word
     - 支持自定义字段选择（查询结果导出）
     - 根据用户角色动态显示导出选项（普通用户仅显示自己提交的论文导出）
     - 根据权限控制导出按钮的显示/隐藏（stats:export 权限控制统计导出）
     - 处理导出加载状态和错误提示

- [x] 2.5 创建图表组件
     【目标对象】`src/frontend/src/components/common/StatsCharts.tsx`
     【修改目的】展示统计图表
     【修改方式】新建文件，创建 ECharts 图表组件
     【相关依赖】echarts, echarts-for-react
     【修改内容】
     - 创建 BarChart 组件（柱状图）
     - 创建 LineChart 组件（折线图）
     - 创建 PieChart 组件（饼图）
     - 创建 TableChart 组件（表格）
     - 实现图表导出功能（PNG、PDF）

- [x] 2.6 创建统计分析页面
     【目标对象】`src/frontend/src/pages/stats/StatsPage.tsx`
     【修改目的】展示统计分析页面
     【修改方式】新建文件，创建统计页面
     【相关依赖】`src/frontend/src/services/statsService.ts`, `src/frontend/src/components/common/StatsCharts.tsx`, `src/frontend/src/components/common/ExportButton.tsx`
     【修改内容】
     - 创建统计页面布局
     - 实现基础统计面板（数字卡片展示：论文总数、年份分布、收录类型分布、期刊分布、平均影响因子、总引用次数、总他引次数）
     - 实现多维度切换（作者、课题、单位）
     - 实现图表展示区域（柱状图、折线图、饼图、表格）
     - 实现图表导出功能（PNG、PDF 格式）
     - 添加权限控制：仅管理员（stats:export 权限）可见统计导出按钮，普通用户不可导出统计数据

- [x] 2.7 添加统计页面路由
     【目标对象】`src/frontend/src/router/index.tsx`
     【修改目的】添加统计页面路由
     【修改方式】在现有路由配置中添加统计页面路由
     【相关依赖】`src/frontend/src/pages/stats/StatsPage.tsx`
     【修改内容】
     - 懒加载 StatsPage 组件
     - 添加 /statistics 路由配置

- [x] 2.8 更新侧边栏菜单
     【目标对象】`src/frontend/src/components/Layout/MainLayout.tsx`
     【修改目的】添加统计入口菜单
     【修改方式】在现有菜单配置中添加统计菜单项
     【相关依赖】无
     【修改内容】
     - 添加"统计分析"菜单项
     - 添加权限控制（stats:view 权限可见）

- [x] 2.9 更新搜索页面导出功能
     【目标对象】`src/frontend/src/pages/search/SearchPage.tsx`
     【修改目的】添加查询结果导出功能
     【修改方式】在现有搜索页面中添加导出按钮
     【相关依赖】`src/frontend/src/components/common/ExportButton.tsx`
     【修改内容】
     - 在搜索结果区域添加导出按钮
     - 支持选择导出字段（标题、作者、期刊、发表年份、收录类型、影响因子、引用次数等）
     - 处理权限控制：普通用户仅可导出自己提交的论文，管理员可导出所有数据
     - 显示导出字段选择弹窗
     - 处理导出加载状态和下载

- [x] 2.10 更新论文详情页导出功能
     【目标对象】`src/frontend/src/pages/paper/PaperDetailPage.tsx`
     【修改目的】添加单篇论文导出功能
     【修改方式】在现有论文详情页面中添加导出按钮
     【相关依赖】`src/frontend/src/components/common/ExportButton.tsx`
     【修改内容】
     - 在论文详情页添加导出按钮（下拉菜单：PDF、Word）
     - 支持导出 PDF 格式（完整信息、作者信息、课题信息、审核记录）
     - 支持导出 Word 格式（完整信息、作者信息、课题信息、审核记录）
     - 处理权限控制：普通用户仅可导出自己提交的论文，管理员可导出所有数据
     - 非权限范围内的论文不显示导出按钮
