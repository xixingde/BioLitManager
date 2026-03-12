# BioLitManager - 生命科学论文管理数据库系统

## 项目简介

BioLitManager 是一套面向生命科学领域的论文管理数据库系统，主要服务于蛋白质相关研究的科研团队。系统实现论文全生命周期的规范化管理，包括论文录入、审核流程、查询检索、归档管理、统计分析等功能。

## 技术栈

### 后端
- **语言**: Go 1.25+
- **框架**: Gin
- **数据库**: SQLite + GORM
- **认证**: JWT
- **配置**: Viper
- **日志**: Zap

### 前端
- **框架**: React 18 + TypeScript
- **UI库**: Ant Design 5
- **构建工具**: Vite
- **状态管理**: Zustand
- **路由**: React Router 6
- **图表**: ECharts

## 功能特性

- 论文信息录入（手动录入 + 批量导入）
- 双审核流程（业务审核 + 政工审核）
- 多维度查询检索
- 自动归档管理
- 数据统计与导出
- 基于角色的权限控制
- 操作日志审计

## 项目结构

```
BioLitManager/
├── src/
│   ├── backend/           # Go 后端服务
│   │   ├── cmd/server/    # 入口程序
│   │   ├── internal/      # 业务逻辑
│   │   │   ├── handler/   # HTTP 处理器
│   │   │   ├── service/   # 业务服务层
│   │   │   ├── repository/# 数据访问层
│   │   │   ├── middleware/# 中间件
│   │   │   ├── security/  # 安全模块
│   │   │   └── config/    # 配置管理
│   │   ├── pkg/           # 公共包
│   │   └── data/          # 数据库文件
│   └── frontend/          # React 前端
│       ├── src/           # 源代码
│       └── dist/          # 构建产物
└── 论文管理数据库系统需求文档V2.md  # 需求文档
```

## 快速开始

### 前置要求

- Go 1.25+
- Node.js 18+
- npm 或 yarn

### 后端编译

```bash
cd src/backend
go build -o server.exe ./cmd/server
```

编译产物为 `server.exe`，位于 `src/backend` 目录下。

### 后端启动

```bash
cd src/backend
# 方式一：直接运行源码
go run cmd/server/main.go

# 方式二：运行编译后的可执行文件
./server.exe
```

服务默认运行在 `http://localhost:8080`

### 前端启动

```bash
cd src/frontend
npm install
npm run dev
```

前端默认运行在 `http://localhost:5173`

## 配置

后端配置文件位于 `src/backend/config.yaml`，可修改服务端口、数据库路径等配置。

## 角色说明

| 角色 | 权限 |
|------|------|
| 科研人员 | 论文录入、查询自己提交的论文 |
| 业务审核人员 | 业务审核、查询待审核论文 |
| 政工审核人员 | 政工审核、查询待审核论文 |
| 科研管理人员 | 查询所有论文、统计分析、数据导出 |
| 系统管理员 | 全权限、用户管理、系统配置 |

## 许可证

MIT License
