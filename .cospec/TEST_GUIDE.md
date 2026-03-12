# Test Guide

## 可运行性验证命令

### 前端
```bash
cd src/frontend
npm run build
```
说明：执行 TypeScript 编译和 Vite 构建验证

### 后端
```bash
cd src/backend
go build
```
说明：执行 Go 语言构建验证

## 测试用例管理方法

当前项目未建立测试用例管理体系，未找到现有测试文件。

建议：
- 前端测试文件位置：`src/frontend/src/**/*.{test,spec}.{ts,tsx}`
- 后端测试文件位置：`src/backend/**/*_test.go`
- 测试文件应与源码文件位于同一目录下，便于维护

## 测试执行方法

### 前端
```bash
cd src/frontend
npm test
```
说明：需先配置 Jest 或 Vitest 测试框架

### 后端
```bash
cd src/backend
go test ./...
```
说明：运行所有 Go 测试文件（依赖 github.com/stretchr/testify）
