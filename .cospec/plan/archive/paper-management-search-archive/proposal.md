# 变更：多维度组合查询检索 + 论文归档管理功能实现

## 原因
根据需求文档，系统需要实现论文查询检索和归档管理功能。当前系统已有基本的论文管理和审核流程，但缺少：
1. 多维度组合查询检索功能（FR-045至FR-055）
2. 完整的论文归档管理功能（FR-056至FR-063）

这些功能是用户使用频率最高且对科研管理至关重要的核心功能。

## 变更内容

### 1. 后端 - 搜索服务（Search Service）
- 新增 `search_service.go`：实现多维度组合查询逻辑
- 新增 `search_handler.go`：提供搜索API接口
- 新增搜索请求/响应DTO

### 2. 后端 - 归档服务扩展（Archive Service）
- 扩展 `archive_service.go`：增加分类查询、归档修改申请等功能
- 扩展 `archive_handler.go`：提供归档管理API接口
- 扩展 `archive_repository.go`：增加分类查询方法

### 3. 数据库扩展
- 扩展 papers 表：添加缺失字段（DOI、影响因子、分区、收录类型等）
- 扩展 archives 表：添加 status 字段（公开/隐藏）

### 4. 前端 - 搜索页面
- 新增 `SearchPage.tsx`：多维度组合查询页面
- 新增 `searchService.ts`：搜索API服务
- 新增搜索相关类型定义

### 5. 前端 - 归档管理页面
- 新增 `ArchivePage.tsx`：归档管理页面（分类查看）
- 新增 `archiveService.ts`：归档API服务

### 6. 路由配置
- 更新 `router/index.tsx`：添加搜索和归档路由

## 影响

### 受影响的规范
- 数据管理：论文查询检索、归档管理

### 受影响的代码
- **后端**：
  - `src/backend/internal/handler/search_handler.go`：新增搜索处理器
  - `src/backend/internal/service/search_service.go`：新增搜索服务
  - `src/backend/internal/repository/paper_repository.go`：扩展查询方法
  - `src/backend/internal/service/archive_service.go`：扩展归档服务
  - `src/backend/internal/handler/archive_handler.go`：扩展归档处理器
  - `src/backend/internal/repository/archive_repository.go`：扩展归档仓储
  - `src/backend/internal/model/entity/paper.go`：扩展论文实体
  - `src/backend/internal/model/entity/archive.go`：扩展归档实体
  - `src/backend/cmd/server/main.go`：注册路由

- **前端**：
  - `src/frontend/src/pages/search/SearchPage.tsx`：新增搜索页面
  - `src/frontend/src/pages/archive/ArchivePage.tsx`：新增归档页面
  - `src/frontend/src/services/searchService.ts`：新增搜索服务
  - `src/frontend/src/services/archiveService.ts`：新增归档服务
  - `src/frontend/src/router/index.tsx`：更新路由配置
  - `src/frontend/src/types/paper.ts`：扩展类型定义

## 功能实现详情

### 查询检索功能（FR-045至FR-055）
1. **多维度组合查询**：支持论文ID、标题、作者、期刊、日期、DOI、影响因子、课题等条件
2. **查询结果展示**：论文ID、标题、作者列表、期刊、出版日期、影响因子、分区、审核状态
3. **排序功能**：按出版日期、影响因子、引用次数升序/降序
4. **精准检索**：DOI、PubMedID、ISSN精确匹配
5. **组合逻辑**：且/或/非逻辑组合
6. **作者类型检索**：第一作者、共同第一作者、通讯作者
7. **课题关联检索**：通过课题编号、项目类型等检索
8. **权限控制**：仅审核通过的论文可公开查询
9. **分页展示**：每页20条，超过1000条正常分页
10. **无结果提示**：显示"未找到符合条件的论文"

### 归档管理功能（FR-056至FR-063）
1. **自动归档**：审核通过后自动归档
2. **归档编号生成**：年份+论文ID+随机3位数字
3. **按年份分类查看**
4. **按收录类型分类**：SCI/EI/中文核心等
5. **按课题分类**
6. **按作者分类**
7. **归档修改申请**：二次审核流程
8. **归档论文隐藏**：不可删除，仅可隐藏
