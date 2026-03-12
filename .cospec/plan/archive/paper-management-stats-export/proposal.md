# 变更：统计分析 + 数据导出功能实现

## 原因
根据需求文档（spec.md）中的功能需求 FR-064 至 FR-073，需要实现统计分析（多维度统计、图表展示、图表导出）和数据导出（查询结果导出、论文导出、统计结果导出、权限控制）功能，为科研管理、成果申报、学术影响力分析提供数据支撑，满足用户存档、申报材料等实际使用需求。

## 变更内容

### 1. 后端 API 接口设计
- **统计服务** (`internal/service/stats_service.go`)
  - 基础指标统计接口：论文总数、各年份发表数量、各收录类型论文数量、各期刊发表数量、平均影响因子、总引用次数、总他引次数
  - 按作者统计接口：发表论文数量、第一作者论文数量、通讯作者论文数量、平均影响因子、总引用次数
  - 按课题统计接口：关联论文数量、高影响因子论文数量、SCI 收录论文数量
  - 按单位/部门统计接口：发表论文数量、总影响因子、总引用次数

- **导出服务** (`internal/service/export_service.go`)
  - 查询结果导出 Excel 接口：支持自定义导出字段
  - 单篇论文导出 PDF 接口：包含论文完整信息、作者信息、课题信息及审核记录
  - 统计结果导出 Excel/PDF 接口：包含图表数据和数据表格

### 2. 前端页面和组件设计
- **统计分析页面** (`src/frontend/src/pages/stats/StatsPage.tsx`)
  - 基础统计面板：展示论文总数、各年份发表数量、各收录类型论文数量、平均影响因子、总引用次数等
  - 多维度统计切换：作者、课题、单位/部门
  - 图表展示区域：柱状图、折线图、饼图、表格
  - 图表导出按钮：PNG、PDF 格式

- **数据导出组件** (`src/frontend/src/components/common/ExportButton.tsx`)
  - 查询结果导出：Excel 格式，支持字段选择
  - 论文详情导出：PDF 格式
  - 统计结果导出：Excel/PDF 格式

### 3. 图表组件选型和实现
- 使用 **ECharts 5.x** 作为图表库
- 图表类型：
  - 柱状图：各年份发表数量、各单位发表数量
  - 折线图：年度趋势
  - 饼图：收录类型分布、期刊分布
  - 表格：详细数据展示
- 图表导出：使用 ECharts 内置的导出功能，支持 PNG、PDF

### 4. Excel/Word/PDF 导出实现
- **Excel 导出**：使用 `excelize` 库，支持自定义字段、多Sheet、格式化
- **PDF 导出**：使用 `gofpdf` 库，支持图表嵌入、表格生成
- 导出文件存储在 `uploads/exports/` 目录，支持下载

### 5. 权限控制逻辑
- 基于现有 RBAC 权限系统
- `paper:export` 权限：普通用户仅可导出自己提交的论文，管理员可导出所有数据
- `stats:export` 权限：仅管理员可导出统计数据
- 导出接口增加权限校验和数据范围控制

## 影响
- **受影响的规范**：数据管理、统计分析
- **受影响的代码**：
  - `{src/backend/cmd/server/main.go}`: 添加统计和导出路由注册
  - `{src/backend/internal/handler/stats_handler.go}`: 新增统计处理器
  - `{src/backend/internal/handler/export_handler.go}`: 新增导出处理器
  - `{src/backend/internal/service/stats_service.go}`: 新增统计服务
  - `{src/backend/internal/service/export_service.go}`: 新增导出服务
  - `{src/backend/internal/repository/stats_repository.go}`: 新增统计仓储
  - `{src/frontend/src/router/index.tsx}`: 添加统计页面路由
  - `{src/frontend/src/pages/stats/StatsPage.tsx}`: 新增统计页面
  - `{src/frontend/src/services/statsService.ts}`: 扩展统计服务
  - `{src/frontend/src/services/exportService.ts}`: 新增导出服务
  - `{src/frontend/src/components/common/ExportButton.tsx}`: 新增导出组件
  - `{src/frontend/src/types/statistics.ts}`: 扩展统计类型定义
