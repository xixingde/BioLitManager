# 变更：论文管理与双审核流程

## 原因
实现论文全生命周期管理的核心功能，包括论文信息录入（手动录入和批量导入）、作者信息管理、课题信息管理、双审核流程（业务审核+政工审核）、附件上传、草稿保存、重复校验、审核时限提醒、审核驳回重提等功能。这是论文管理系统的核心业务模块，是用户故事2-8和FR-008至FR-044的具体实现。

## 变更内容

### 后端实现
- **新增数据库表**：
  - `papers` 表：论文基础信息
  - `authors` 表：论文作者信息
  - `projects` 表：课题信息
  - `paper_projects` 表：论文-课题关联表
  - `journals` 表：期刊信息
  - `review_logs` 表：审核记录
  - `archives` 表：归档记录
  - `attachments` 表：附件信息

- **新增后端模块**：
  - 论文管理模块（`paper_handler.go`, `paper_service.go`, `paper_repository.go`）
  - 作者管理模块（`author_handler.go`, `author_service.go`, `author_repository.go`）
  - 课题管理模块（`project_handler.go`, `project_service.go`, `project_repository.go`）
  - 期刊管理模块（`journal_handler.go`, `journal_service.go`, `journal_repository.go`）
  - 审核管理模块（`review_handler.go`, `review_service.go`, `review_repository.go`）
  - 附件管理模块（`file_handler.go`, `file_service.go`, `file_repository.go`）
  - 通知服务（`notification_service.go`）

- **API接口**：
  - 论文CRUD接口：`GET/POST /api/papers`, `GET/PUT/DELETE /api/papers/:id`
  - 论文提交审核：`POST /api/papers/:id/submit`
  - 论文草稿保存：`POST /api/papers/:id/save-draft`
  - 重复校验：`POST /api/papers/check-duplicate`
  - 业务审核接口：`POST /api/reviews/business/:paperId`
  - 政工审核接口：`POST /api/reviews/political/:paperId`
  - 批量导入接口：`POST /api/papers/batch-import`
  - 文件上传接口：`POST /api/files/upload`

### 前端实现
- **新增页面组件**：
  - 论文列表页（`PaperListPage.tsx`）
  - 论文录入页（`PaperCreatePage.tsx`）
  - 论文详情页（`PaperDetailPage.tsx`）
  - 批量导入页（`BatchImportPage.tsx`）
  - 业务审核列表页（`BusinessReviewListPage.tsx`）
  - 政工审核列表页（`PoliticalReviewListPage.tsx`）
  - 审核页面（`ReviewPage.tsx`）

- **新增业务组件**：
  - 论文表单组件（`PaperForm.tsx`）
  - 作者列表组件（`AuthorList.tsx`）
  - 课题选择器（`ProjectSelector.tsx`）
  - 审核表单组件（`ReviewForm.tsx`）
  - 文件上传组件（`FileUpload.tsx`）
  - 批量导入组件（`BatchImportForm.tsx`）

- **新增服务层**：
  - `paperService.ts`：论文相关API
  - `authorService.ts`：作者相关API
  - `projectService.ts`：课题相关API
  - `journalService.ts`：期刊相关API
  - `reviewService.ts`：审核相关API
  - `fileService.ts`：文件上传API

- **新增状态管理**：
  - `paperStore.ts`：论文状态管理

- **新增类型定义**：
  - `paper.ts`：论文相关类型
  - `author.ts`：作者相关类型
  - `project.ts`：课题相关类型
  - `review.ts`：审核相关类型

### 核心功能
1. **论文信息录入**：支持手动录入论文信息，包含字段校验（DOI、ISSN、日期格式等）
2. **批量导入论文**：通过Excel模板批量导入，支持数据校验和错误提示
3. **作者信息管理**：支持关联人员库，支持作者类型设置（第一作者、共同第一作者、通讯作者等）
4. **课题信息管理**：支持关联课题库，支持手动录入新课题
5. **论文提交审核**：提交后状态变更为"待业务审核"，不可再修改
6. **业务审核**：审核学术规范性和真实性，审核通过后进入政工审核
7. **政工审核**：审核政治合规性，审核通过后自动归档
8. **附件上传**：支持上传论文全文PDF、首页、期刊封面、审批件，最大100MB
9. **草稿保存**：支持保存未完成的录入内容为草稿
10. **重复校验**：根据论文标题+DOI判断是否重复
11. **审核时限提醒**：3个工作日时限，逾期自动发送提醒
12. **审核驳回重提**：驳回后可修改信息重新提交

## 影响
- **受影响的规范**：论文管理（FR-008至FR-044）
- **受影响的代码**：
  - 后端：新增 `src/backend/internal/handler/paper_handler.go`、`review_handler.go`、`file_handler.go` 等8个Handler文件
  - 后端：新增 `src/backend/internal/service/paper_service.go`、`review_service.go`、`file_service.go` 等8个Service文件
  - 后端：新增 `src/backend/internal/repository/paper_repository.go`、`review_repository.go`、`file_repository.go` 等8个Repository文件
  - 后端：新增 `src/backend/internal/model/entity/paper.go`、`author.go`、`project.go`、`journal.go`、`review_log.go`、`archive.go`、`attachment.go` 等7个实体文件
  - 后端：新增 `src/backend/internal/model/dto/request/paper_request.go`、`review_request.go` 等请求DTO文件
  - 后端：新增 `src/backend/internal/model/dto/response/paper_response.go`、`review_response.go` 等响应DTO文件
  - 后端：修改 `src/backend/cmd/server/main.go` 新增路由注册
  - 后端：修改 `src/backend/internal/database/migration.go` 新增数据库迁移
  - 后端：修改 `src/backend/pkg/errors/errors.go` 新增错误码定义
  - 后端：修改 `src/backend/go.mod` 新增 `github.com/xuri/excelize/v2` 依赖
  - 前端：新增 `src/frontend/src/pages/paper/` 目录下的5个页面组件
  - 前端：新增 `src/frontend/src/components/paper/` 目录下的6个业务组件
  - 前端：新增 `src/frontend/src/services/paperService.ts`、`reviewService.ts` 等6个服务文件
  - 前端：新增 `src/frontend/src/stores/paperStore.ts` 状态管理文件
  - 前端：新增 `src/frontend/src/types/paper.ts`、`review.ts` 等5个类型定义文件
  - 前端：新增 `src/frontend/src/router/index.tsx` 路由配置
  - 前端：修改 `src/frontend/src/components/Layout/Layout.tsx` 新增菜单项
  - 前端：修改 `src/frontend/package.json` 新增 `react-hook-form`、`echarts`、`@types/react-hook-form` 依赖

## 数据库变更
- 新增 7 张数据表：`papers`、`authors`、`projects`、`paper_projects`、`journals`、`review_logs`、`archives`、`attachments`
- 修改 `users` 表：新增 `is_email_notify` 字段（邮件通知配置）
- 修改 `operation_logs` 表：支持审核类型操作记录

## 依赖变更
- 后端新增：`github.com/xuri/excelize/v2`（Excel处理）
- 前端新增：`react-hook-form`、`@types/react-hook-form`（表单管理）、`echarts`（图表）

## 验证方法
1. **功能测试**：使用不同角色账号测试论文录入、编辑、删除、提交审核、审核、驳回等完整流程
2. **字段校验测试**：测试DOI、ISSN、日期等字段的格式校验
3. **重复校验测试**：测试重复论文的校验功能
4. **附件上传测试**：测试PDF文件上传（正常大小、超过100MB）
5. **批量导入测试**：测试Excel批量导入功能（正常数据、错误数据）
6. **审核流程测试**：测试业务审核、政工审核的完整流程，包括驳回重提
7. **审核时限提醒测试**：测试逾期提醒功能
8. **权限测试**：测试不同角色对论文和审核功能的访问权限
