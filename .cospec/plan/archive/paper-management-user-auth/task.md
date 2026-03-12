## 实施

### 前端基础架构搭建

  - [x] 1.1 初始化前端项目（Vite + React + TypeScript）
      【目标对象】`src/frontend/package.json`
      【修改目的】创建前端项目并配置核心依赖
      【修改方式】新建文件
      【相关依赖】无
      【修改内容】
         - 初始化 Vite + React + TypeScript 项目
         - 安装依赖：antd@5.x、zustand@4.x、react-router-dom@6.x、axios@1.x、@ant-design/icons、dayjs
         - 配置 package.json 的 scripts（dev、build、preview）

  - [x] 1.2 配置 Vite 构建工具
      【目标对象】`src/frontend/vite.config.ts`
      【修改目的】配置 Vite 构建工具和开发服务器代理
      【修改方式】新建文件
      【相关依赖】`vite@5.x`
      【修改内容】
         - 配置开发服务器代理（API 请求代理到后端 8080 端口）
         - 配置路径别名（@ 指向 src 目录）
         - 配置构建输出目录

  - [x] 1.3 配置 TypeScript
      【目标对象】`src/frontend/tsconfig.json`
      【修改目的】配置 TypeScript 编译选项和路径别名
      【修改方式】新建文件
      【相关依赖】`typescript@5.x`
      【修改内容】
         - 配置编译选项（strict: true、jsx: react-jsx）
         - 配置路径别名（@/*）
         - 配置模块解析（node）

  - [x] 1.4 创建前端目录结构
      【目标对象】`src/frontend/src/`
      【修改目的】创建前端项目的基础目录结构
      【修改方式】创建目录
      【相关依赖】无
      【修改内容】
         - 创建 `components/` 目录（通用组件、业务组件）
         - 创建 `pages/` 目录（页面组件）
         - 创建 `stores/` 目录（状态管理）
         - 创建 `services/` 目录（API 服务）
         - 创建 `hooks/` 目录（自定义 Hooks）
         - 创建 `utils/` 目录（工具函数）
         - 创建 `types/` 目录（TypeScript 类型定义）

  - [x] 1.5 创建 API 客户端封装
      【目标对象】`src/frontend/src/services/api.ts`
      【修改目的】封装 Axios HTTP 客户端，统一处理请求拦截和响应拦截
      【修改方式】新建文件
      【相关依赖】`axios@1.x`
      【修改内容】
         - 创建 axios 实例，配置 baseURL 和超时时间
         - 实现请求拦截器（自动添加 Authorization Header）
         - 实现响应拦截器（统一处理错误、Token 过期跳转登录）
         - 导出 GET、POST、PUT、DELETE 方法

  - [x] 1.6 创建认证服务
      【目标对象】`src/frontend/src/services/authService.ts`
      【修改目的】封装用户认证相关的 API 调用
      【修改方式】新建文件
      【相关依赖】`src/frontend/src/services/api.ts`
      【修改内容】
         - 创建 LoginRequest 接口（username、password）
         - 创建 LoginResponse 接口（token、user）
         - 实现 login 方法（POST /api/auth/login）
         - 实现 logout 方法（POST /api/auth/logout）
         - 实现 getProfile 方法（GET /api/auth/profile）

  - [x] 1.7 创建用户状态管理（Zustand Store）
      【目标对象】`src/frontend/src/stores/userStore.ts`
      【修改目的】使用 Zustand 管理用户登录状态和权限信息
      【修改方式】新建文件
      【相关依赖】`zustand@4.x`
      【修改内容】
         - 定义 UserState 接口（user、token、permissions、isAuthenticated）
         - 实现 userStore，包含 state 和 actions
         - 实现 login action（调用 authService.login，存储 token 和用户信息）
         - 实现 logout action（清除 token 和用户信息）
         - 实现 updatePermissions action（更新权限列表）
         - 实现持久化（token 存储到 localStorage）

  - [x] 1.8 创建 TypeScript 类型定义
      【目标对象】`src/frontend/src/types/user.ts`
      【修改目的】定义用户相关的 TypeScript 类型
      【修改方式】新建文件
      【相关依赖】无
      【修改内容】
         - 定义 UserInfo 接口（id、username、name、role、department、email）
         - 定义 Role 类型（7种角色：super_admin、admin、dept_head、project_leader、business_reviewer、political_reviewer、user）
         - 定义 Permission 类型
         - 导出相关类型供其他模块使用

  - [x] 1.9 创建 API 响应类型定义
      【目标对象】`src/frontend/src/types/api.ts`
      【修改目的】定义 API 统一响应格式类型
      【修改方式】新建文件
      【相关依赖】无
      【修改内容】
         - 定义 ApiResponse<T> 接口（code、msg、data）
         - 定义 PageResult<T> 接口（list、total、page、size）
         - 导出相关类型供其他模块使用

  - [x] 1.10 创建登录页面组件
      【目标对象】`src/frontend/src/pages/login/LoginPage.tsx`
      【修改目的】实现用户登录界面，包含用户名密码输入和登录按钮
      【修改方式】新建文件，使用 Ant Design Form 组件创建登录表单
      【相关依赖】`antd@5.x`、`src/frontend/src/stores/userStore.ts`、`src/frontend/src/services/authService.ts`
      【修改内容】
         - 创建登录表单组件，包含用户名和密码输入框
         - 实现密码复杂度校验提示（至少8位，含大小写字母、数字、特殊字符）
         - 实现登录按钮和提交逻辑（调用 authService.login）
         - 实现错误提示（登录失败、账户锁定、账户禁用）
         - 登录成功后跳转到首页
         - 添加页面样式（居中布局、卡片样式）

  - [x] 1.11 创建自定义认证 Hook
      【目标对象】`src/frontend/src/hooks/useAuth.ts`
      【修改目的】封装认证相关的逻辑，供组件复用
      【修改方式】新建文件，创建自定义 Hook
      【相关依赖】`src/frontend/src/stores/userStore.ts`
      【修改内容】
         - 实现 useAuth Hook，返回 user、isAuthenticated、login、logout 等方法
         - 实现权限校验方法 hasPermission
         - 导出供其他组件使用

  - [x] 1.12 创建路由配置和路由保护
      【目标对象】`src/frontend/src/router/index.tsx`
      【修改目的】配置 React Router，实现路由保护（未登录跳转登录页）
      【修改方式】新建文件，创建路由配置组件
      【相关依赖】`react-router-dom@6.x`、`src/frontend/src/hooks/useAuth.ts`
      【修改内容】
         - 创建 BrowserRouter 和 Routes 配置
         - 创建 PrivateRoute 组件（检查登录状态，未登录跳转登录页）
         - 配置登录路由（/login）
         - 配置受保护路由（/、/papers、/reviews 等）
         - 实现自动重定向（登录成功后跳转原目标页面）

  - [x] 1.13 创建常量配置文件
      【目标对象】`src/frontend/src/utils/constants.ts`
      【修改目的】定义系统常量（API 地址、Token 键名、权限列表等）
      【修改方式】新建文件，定义常量
      【相关依赖】无
      【修改内容】
         - 定义 API_BASE_URL
         - 定义 TOKEN_KEY
         - 定义 USER_INFO_KEY
         - 定义 ROLE 常量（7种角色）
         - 定义 PERMISSION 常量

  - [x] 1.14 创建 Ant Design 全局配置
      【目标对象】`src/frontend/src/main.tsx`
      【修改目的】配置 Ant Design 主题和全局样式
      【修改方式】修改文件，添加 Ant Design 配置和 Router 渲染
      【相关依赖】`antd@5.x`
      【修改内容】
         - 配置 Ant Design ConfigProvider（中文语言包、主题颜色）
         - 引入 Ant Design 样式文件
         - 渲染 Router 组件

### 后端基础架构搭建

  - [x] 2.1 初始化后端项目（Go mod）
      【目标对象】`src/backend/go.mod`
      【修改目的】创建 Go 项目并配置核心依赖
      【修改方式】新建文件，初始化 Go module
      【相关依赖】无
      【修改内容】
         - 初始化 Go module（go mod init）
         - 配置依赖：gin-gonic/gin@1.9+、gorm.io/gorm@2.5+、gorm.io/driver/sqlite、golang.org/x/crypto/bcrypt、github.com/golang-jwt/jwt/v5、go.uber.org/zap、spf13/viper、github.com/google/uuid

  - [x] 2.2 创建后端目录结构
      【目标对象】`src/backend/`
      【修改目的】创建后端项目的基础目录结构
      【修改方式】创建目录
      【相关依赖】无
      【修改内容】
         - 创建 `cmd/server/` 目录（应用入口）
         - 创建 `internal/handler/` 目录（Handler 层）
         - 创建 `internal/service/` 目录（Service 层）
         - 创建 `internal/repository/` 目录（Repository 层）
         - 创建 `internal/model/` 目录（数据模型）
         - 创建 `internal/middleware/` 目录（中间件）
         - 创建 `internal/config/` 目录（配置）
         - 创建 `internal/database/` 目录（数据库）
         - 创建 `internal/cache/` 目录（缓存）
         - 创建 `internal/security/` 目录（安全）
         - 创建 `internal/utils/` 目录（工具）
         - 创建 `pkg/logger/` 目录（日志）
         - 创建 `pkg/response/` 目录（响应封装）

  - [x] 2.3 创建统一响应格式封装
      【目标对象】`src/backend/pkg/response/response.go`
      【修改目的】封装统一的 API 响应格式
      【修改方式】新建文件，定义响应结构体和方法
      【相关依赖】无
      【修改内容】
         - 定义 ApiResponse 结构体（code、msg、data）
         - 实现 Success 方法
         - 实现 Error 方法
         - 定义 PageResult 结构体（list、total、page、size）

  - [x] 2.4 创建错误码定义
      【目标对象】`src/backend/pkg/errors/errors.go`
      【修改目的】定义系统错误码和错误消息
      【修改方式】新建文件，定义错误码常量
      【相关依赖】无
      【修改内容】
         - 定义 ErrorCode 结构体（code、msg）
         - 定义通用错误（000000 成功、100001 系统异常、100002 参数错误、100003 未授权、100004 无权限、100005 资源不存在）
         - 定义认证错误（101001 登录失败、101002 账户锁定、101003 账户禁用、101004 会话过期）

  - [x] 2.5 创建配置管理模块
      【目标对象】`src/backend/internal/config/config.go`
      【修改目的】使用 Viper 加载和管理系统配置
      【修改方式】新建文件，实现配置加载和管理
      【相关依赖】`spf13/viper`
      【修改内容】
         - 定义 Config 结构体（server、database、jwt、cache）
         - 定义 ServerConfig 结构体（port、mode）
         - 定义 DatabaseConfig 结构体（path）
         - 定义 JWTConfig 结构体（secret、expire_hours）
         - 定义 CacheConfig 结构体（session_ttl、user_info_ttl）
         - 实现 LoadConfig 方法（从 config.yaml 加载配置）
         - 提供 GetConfig 方法

  - [x] 2.6 创建配置文件
      【目标对象】`src/backend/config.yaml`
      【修改目的】定义系统配置参数
      【修改方式】新建文件，定义配置项
      【相关依赖】无
      【修改内容】
         - server 配置（port: 8080、mode: debug）
         - database 配置（path: ./data/biolit.db）
         - jwt 配置（secret: biolitmanager、expire_hours: 2）
         - cache 配置（session_ttl: 2h、user_info_ttl: 30m）

  - [x] 2.7 创建数据库配置（SQLite + WAL 模式）
      【目标对象】`src/backend/internal/database/sqlite.go`
      【修改目的】配置 SQLite 数据库连接，启用 WAL 模式和连接池
      【修改方式】新建文件，实现数据库初始化和连接池配置
      【相关依赖】`gorm.io/gorm`、`gorm.io/driver/sqlite`
      【修改内容】
         - 实现 InitDB 方法（连接 SQLite 数据库）
         - 配置 DSN（启用 WAL 模式、设置 busy_timeout、设置 cache_size）
         - 配置连接池（SetMaxIdleConns=10、SetMaxOpenConns=100、SetConnMaxLifetime=1h）
         - 返回 *gorm.DB 实例

  - [x] 2.8 创建数据库迁移模块
      【目标对象】`src/backend/internal/database/migration.go`
      【修改目的】实现数据库表自动迁移（AutoMigrate）
      【修改方式】新建文件，实现数据库迁移
      【相关依赖】`src/backend/internal/database/sqlite.go`、`src/backend/internal/model/entity/`
      【修改内容】
         - 实现 AutoMigrate 方法
         - 迁移 User、OperationLog 等表

  - [x] 2.9 创建内存缓存实现
      【目标对象】`src/backend/internal/cache/memory_cache.go`
      【修改目的】实现基于 Go Map + RWMutex 的内存缓存
      【修改方式】新建文件，实现缓存结构体和方法
      【相关依赖】无
      【修改内容】
         - 定义 Item 结构体（value、expiration）
         - 定义 MemoryCache 结构体（items、mu）
         - 实现 NewMemoryCache 方法
         - 实现 Set 方法（设置缓存，支持过期时间）
         - 实现 Get 方法（获取缓存，自动检查过期）
         - 实现 Delete 方法（删除缓存）
         - 实现启动清理过期缓存的 goroutine

  - [x] 2.10 创建缓存键定义
      【目标对象】`src/backend/internal/cache/cache_keys.go`
      【修改目的】定义缓存键命名规范
      【修改方式】新建文件，定义缓存键生成方法
      【相关依赖】无
      【修改内容】
         - 定义 SessionKey 方法（session:{token}）
         - 定义 UserKey 方法（user:{userId}）
         - 定义 ConfigKey 方法（config:{key}）

### 数据模型和仓储层

  - [x] 3.1 创建用户实体模型
      【目标对象】`src/backend/internal/model/entity/user.go`
      【修改目的】定义用户表的数据模型
      【修改方式】新建文件，定义 User 结构体
      【相关依赖】`gorm.io/gorm`
      【修改内容】
         - 定义 User 结构体（id、username、password_hash、name、role、department、id_card、phone、email、is_locked、lock_until、is_disabled、login_fail_count、last_login_at、last_login_ip、created_at、updated_at）
         - 定义 TableName 方法（返回 "users"）
         - 添加 GORM 标签（主键、唯一索引、默认值）

  - [x] 3.2 创建用户 DTO
      【目标对象】`src/backend/internal/model/dto/response/user_dto.go`
      【修改目的】定义用户相关的数据传输对象
      【修改方式】新建文件，定义 DTO 结构体
      【相关依赖】无
      【修改内容】
         - 定义 UserDTO 结构体（id、username、name、role、department、email）
         - 定义 LoginRequest 结构体（username、password）
         - 定义 LoginResponse 结构体（token、user）

  - [x] 3.3 创建操作日志实体模型
      【目标对象】`src/backend/internal/model/entity/operation_log.go`
      【修改目的】定义操作日志表的数据模型
      【修改方式】新建文件，定义 OperationLog 结构体
      【相关依赖】`gorm.io/gorm`
      【修改内容】
         - 定义 OperationLog 结构体（id、user_id、operation_type、module、target_id、operation_content、operation_result、ip_address、created_at）
         - 定义 TableName 方法（返回 "operation_logs"）

  - [x] 3.4 创建用户仓储层
      【目标对象】`src/backend/internal/repository/user_repository.go`
      【修改目的】实现用户数据的 CRUD 操作
      【修改方式】新建文件，实现 UserRepository 结构体和方法
      【相关依赖】`src/backend/internal/model/entity/user.go`、`gorm.io/gorm`
      【修改内容】
         - 定义 UserRepository 结构体（db）
         - 实现 NewUserRepository 方法
         - 实现 FindByUsername 方法（根据用户名查询用户）
         - 实现 FindByID 方法（根据 ID 查询用户）
         - 实现 Create 方法（创建用户）
         - 实现 Update 方法（更新用户）
         - 实现 Delete 方法（删除用户）
         - 实现 List 方法（分页查询用户列表）

  - [x] 3.4.1 创建操作日志仓储层
      【目标对象】`src/backend/internal/repository/operation_log_repository.go`
      【修改目的】实现操作日志数据的 CRUD 操作
      【修改方式】新建文件，实现 OperationLogRepository 结构体和方法
      【相关依赖】`src/backend/internal/model/entity/operation_log.go`、`gorm.io/gorm`
      【修改内容】
         - 定义 OperationLogRepository 结构体（db）
         - 实现 NewOperationLogRepository 方法
         - 实现 Create 方法（创建操作日志）
         - 实现 List 方法（分页查询操作日志列表）

### 安全工具和权限模型

- [x] 4.1 创建 JWT 工具
     【目标对象】`src/backend/internal/security/jwt.go`
     【修改目的】实现 JWT Token 生成和解析
     【相关依赖】`github.com/golang-jwt/jwt/v5`
     【修改内容】
        - 定义 Claims 结构体（user_id、username、role、permissions、jwt.RegisteredClaims）
        - 定义 TokenExpireTime 常量（2小时）
        - 定义 jwtSecret 变量
        - 实现 GenerateToken 方法（生成 JWT Token）
        - 实现 ParseToken 方法（解析 JWT Token）

  - [x] 4.2 创建密码哈希工具
      【目标对象】`src/backend/internal/security/password.go`
      【修改目的】实现密码哈希和验证（使用 bcrypt）
      【修改方式】新建文件，实现密码哈希和验证方法
      【相关依赖】`golang.org/x/crypto/bcrypt`
      【修改内容】
         - 定义 Cost 常量（10）
         - 实现 HashPassword 方法（对密码进行哈希）
         - 实现 CheckPassword 方法（验证密码）

  - [x] 4.2.1 创建密码复杂度校验
      【目标对象】`src/backend/internal/security/password.go`
      【修改目的】实现密码复杂度校验逻辑
      【修改方式】在 password.go 文件中添加校验方法
      【相关依赖】无
      【修改内容】
         - 定义 PasswordComplexityError 错误类型
         - 实现 ValidatePasswordComplexity 方法
           - 校验密码长度至少8位
           - 校验包含大写字母
           - 校验包含小写字母
           - 校验包含数字
           - 校验包含特殊字符
         - 返回错误提示，指明不符合哪些条件

  - [x] 4.3 定义角色和权限常量
      【目标对象】`src/backend/internal/security/permission.go`
      【修改目的】定义 RBAC 角色和权限常量
      【修改方式】新建文件，定义角色和权限常量
      【相关依赖】无
      【修改内容】
         - 定义 Role 类型（super_admin、admin、dept_head、project_leader、business_reviewer、political_reviewer、user）
         - 定义 Permission 类型（paper:create、paper:edit、paper:view、paper:delete、paper:export、review:business、review:political、system:user:manage、system:project:manage、system:journal:manage、system:config:manage、stats:view、stats:export）
         - 定义 RolePermissions 映射（角色到权限列表的映射）
         - 实现 GetPermissionsByRole 方法（根据角色获取权限列表）

### 中间件实现

  - [x] 5.1 创建认证中间件
      【目标对象】`src/backend/internal/middleware/auth_middleware.go`
      【修改目的】实现 JWT 认证中间件，验证 Token 并设置用户信息到上下文
      【修改方式】新建文件，实现认证中间件
      【相关依赖】`src/backend/internal/security/jwt.go`、`src/backend/internal/cache/memory_cache.go`、`github.com/gin-gonic/gin`
      【修改内容】
         - 实现 AuthMiddleware 方法
         - 从请求头提取 Authorization Token
         - 解析 Token 验证有效性
         - 检查缓存中的会话是否存在
         - 将用户信息设置到 gin.Context
         - Token 无效或会话过期返回 401 错误

  - [x] 5.2 创建权限中间件
      【目标对象】`src/backend/internal/middleware/permission_middleware.go`
      【修改目的】实现权限校验中间件，检查用户是否拥有指定权限
      【修改方式】新建文件，实现权限中间件
      【相关依赖】`src/backend/internal/security/permission.go`、`github.com/gin-gonic/gin`
      【修改内容】
         - 实现 PermissionMiddleware 方法（接收 requiredPermission 参数）
         - 从上下文获取用户信息
         - 超级管理员拥有所有权限
         - 检查用户权限列表是否包含所需权限
         - 无权限返回 403 错误

  - [x] 5.3 创建日志中间件
      【目标对象】`src/backend/internal/middleware/logger_middleware.go`
      【修改目的】实现请求日志记录中间件，记录所有 HTTP 请求
      【修改方式】新建文件，实现日志中间件
      【相关依赖】`go.uber.org/zap`
      【修改内容】
         - 实现 LoggerMiddleware 方法
         - 记录请求方法、路径、状态码、响应时间、客户端 IP
         - 使用 zap 结构化日志

  - [x] 5.4 创建异常恢复中间件
      【目标对象】`src/backend/internal/middleware/recovery_middleware.go`
      【修改目的】实现全局异常捕获和恢复，防止 panic 导致服务崩溃
      【修改方式】新建文件，实现恢复中间件
      【相关依赖】`github.com/gin-gonic/gin`、`src/backend/pkg/response/response.go`
      【修改内容】
         - 实现 RecoveryMiddleware 方法
         - 使用 defer + recover 捕获 panic
         - 记录错误日志
         - 返回统一的错误响应（500 系统异常）

### 服务层实现

  - [x] 6.1 创建认证服务
      【目标对象】`src/backend/internal/service/auth_service.go`
      【修改目的】实现用户登录认证核心业务逻辑
      【修改方式】新建文件，实现 AuthService 结构体和方法
      【相关依赖】`src/backend/internal/repository/user_repository.go`、`src/backend/internal/security/password.go`、`src/backend/internal/security/jwt.go`、`src/backend/internal/cache/memory_cache.go`、`src/backend/internal/security/permission.go`、`src/backend/internal/service/operation_log_service.go`
      【修改内容】
         - 定义 AuthService 结构体（userRepo、cache、operationLogService）
         - 实现 NewAuthService 方法
         - 实现 Login 方法（登录逻辑）
           - 查询用户，验证密码
           - 检查账户状态（是否禁用、是否锁定）
           - 密码错误则增加失败计数，5次失败锁定1小时
           - 登录成功则重置失败计数
           - 获取用户权限列表
           - 生成 JWT Token
           - 存储会话到内存缓存（有效期2小时）
           - 更新最后登录时间和 IP
           - 记录操作日志（登录成功/失败）
           - 返回 Token 和用户信息
         - 实现 Logout 方法（登出逻辑）
           - 从缓存中删除会话
           - 记录操作日志（登出）
         - 实现 GetProfile 方法（获取当前用户信息）
         - 实现 incrementLoginFailCount 方法（增加登录失败计数）
         - 实现 resetLoginFailCount 方法（重置登录失败计数）

  - [x] 6.2 创建用户服务
      【目标对象】`src/backend/internal/service/user_service.go`
      【修改目的】实现用户管理业务逻辑（创建、修改、禁用/启用用户）
      【修改方式】新建文件，实现 UserService 结构体和方法
      【相关依赖】`src/backend/internal/repository/user_repository.go`、`src/backend/internal/security/password.go`
      【修改内容】
         - 定义 UserService 结构体（userRepo）
         - 实现 NewUserService 方法
         - 实现 CreateUser 方法（创建用户）
           - 校验用户名是否已存在
           - 校验密码复杂度（调用 ValidatePasswordComplexity）
           - 密码哈希
           - 创建用户记录
           - 记录操作日志
         - 实现 UpdateUser 方法（更新用户信息）
         - 实现 DisableUser 方法（禁用用户账户）
         - 实现 EnableUser 方法（启用用户账户）
         - 实现 ListUsers 方法（分页查询用户列表）

  - [x] 6.2.1 创建操作日志服务
      【目标对象】`src/backend/internal/service/operation_log_service.go`
      【修改目的】实现操作日志记录业务逻辑
      【修改方式】新建文件，实现 OperationLogService 结构体和方法
      【相关依赖】`src/backend/internal/repository/operation_log_repository.go`
      【修改内容】
         - 定义 OperationLogService 结构体（operationLogRepo）
         - 实现 NewOperationLogService 方法
         - 实现 LogOperation 方法
           - 创建操作日志记录
           - 记录操作类型、模块、目标 ID、操作内容、结果、IP 地址
           - 异步保存到数据库

  - [x] 6.3 创建 UserInfo 结构体（用于缓存）
      【目标对象】`src/backend/internal/model/user_info.go`
      【修改目的】定义用户信息结构体，用于缓存存储
      【修改方式】新建文件，定义 UserInfo 结构体
      【相关依赖】无
      【修改内容】
         - 定义 UserInfo 结构体（UserID、Username、Role、Permissions）
         - 用于 JWT Claims 和缓存存储

### Handler 层实现

  - [x] 7.1 创建认证处理器
      【目标对象】`src/backend/internal/handler/auth_handler.go`
      【修改目的】实现认证相关的 HTTP 接口（登录、登出、获取用户信息）
      【修改方式】新建文件，实现 AuthHandler 结构体和方法
      【相关依赖】`src/backend/internal/service/auth_service.go`、`src/backend/pkg/response/response.go`、`github.com/gin-gonic/gin`
      【修改内容】
         - 定义 AuthHandler 结构体（authService）
         - 实现 NewAuthHandler 方法
         - 实现 Login 方法（POST /api/auth/login）
           - 绑定 LoginRequest
           - 调用 authService.Login
           - 返回 LoginResponse
         - 实现 Logout 方法（POST /api/auth/logout）
           - 调用 authService.Logout
           - 返回成功响应
         - 实现 GetProfile 方法（GET /api/auth/profile）
           - 调用 authService.GetProfile
           - 返回用户信息

  - [x] 7.2 创建用户处理器
      【目标对象】`src/backend/internal/handler/user_handler.go`
      【修改目的】实现用户管理的 HTTP 接口（创建、更新、禁用/启用、查询）
      【修改方式】新建文件，实现 UserHandler 结构体和方法
      【相关依赖】`src/backend/internal/service/user_service.go`、`src/backend/pkg/response/response.go`、`src/backend/internal/middleware/auth_middleware.go`
      【修改内容】
         - 定义 UserHandler 结构体（userService）
         - 实现 NewUserHandler 方法
         - 实现 CreateUser 方法（POST /api/system/users）
           - 权限校验（system:user:manage）
           - 绑定请求参数
           - 调用 userService.CreateUser
           - 返回用户信息
         - 实现 UpdateUser 方法（PUT /api/system/users/:id）
           - 权限校验（system:user:manage）
           - 绑定请求参数
           - 调用 userService.UpdateUser
           - 返回成功响应
         - 实现 DisableUser 方法（PUT /api/system/users/:id/disable）
           - 权限校验（system:user:manage）
           - 调用 userService.DisableUser
           - 返回成功响应
         - 实现 EnableUser 方法（PUT /api/system/users/:id/enable）
           - 权限校验（system:user:manage）
           - 调用 userService.EnableUser
           - 返回成功响应
         - 实现 ListUsers 方法（GET /api/system/users）
           - 权限校验（system:user:manage）
           - 调用 userService.ListUsers
           - 返回分页结果

### 应用入口和路由配置

  - [x] 8.1 创建日志工具
      【目标对象】`src/backend/pkg/logger/logger.go`
      【修改目的】封装 zap 日志库，提供日志记录方法
      【修改方式】新建文件，实现日志工具
      【相关依赖】`go.uber.org/zap`
      【修改内容】
         - 实现 InitLogger 方法（初始化 zap logger）
         - 提供 Info、Error、Debug、Warn 等方法

  - [x] 8.2 创建应用入口
      【目标对象】`src/backend/cmd/server/main.go`
      【修改目的】实现后端应用启动入口，初始化所有组件和路由
      【修改方式】新建文件，实现应用启动逻辑
      【相关依赖】`src/backend/internal/config/config.go`、`src/backend/internal/database/sqlite.go`、`src/backend/internal/cache/memory_cache.go`、`src/backend/internal/middleware/`、`src/backend/internal/handler/`
      【修改内容】
         - 加载配置
         - 初始化日志
         - 初始化数据库
         - 执行数据库迁移
         - 初始化内存缓存
         - 初始化 repositories、services、handlers
         - 创建 Gin 实例
         - 注册中间件（Recovery、Logger、CORS）
         - 注册路由
           - /api/auth/login（POST）
           - /api/auth/logout（POST）
           - /api/auth/profile（GET，需要认证）
           - /api/system/users/*（CRUD 接口，需要认证和权限）
         - 启动 HTTP 服务器（监听 8080 端口）

### 数据库初始化

  - [x] 9.1 初始化数据库并创建表结构
      【目标对象】`src/backend/data/`
      【修改目的】执行数据库迁移，创建所有必需的表
      【修改方式】创建目录并初始化数据库文件
      【相关依赖】`src/backend/internal/database/migration.go`
      【修改内容】
         - 确保数据库目录存在
         - 运行 AutoMigrate 创建 users、operation_logs 等表
         - 验证表结构是否正确创建

  - [x] 9.2 创建初始管理员账户
      【目标对象】`src/backend/cmd/server/main.go`
      【修改目的】在应用启动时检查并创建初始超级管理员账户
      【修改方式】在 main.go 中添加初始账户创建逻辑
      【相关依赖】`src/backend/internal/repository/user_repository.go`、`src/backend/internal/security/password.go`
      【修改内容】
         - 检查数据库中是否存在超级管理员账户
         - 如果不存在，创建默认超级管理员（username: admin、password: Admin@123、role: super_admin）
         - 记录日志

### 验证和测试

  - [x] 10.1 启动后端服务并验证
      【目标对象】`src/backend/`
      【修改目的】验证后端服务能否正常启动
      【修改方式】执行验证操作
      【相关依赖】无
      【修改内容】
         - 运行 `go run cmd/server/main.go`
         - 检查日志输出，确认服务启动成功
         - 验证数据库文件是否创建
         - 验证初始管理员账户是否创建

  - [x] 10.2 启动前端服务并验证
      【目标对象】`src/frontend/`
      【修改目的】验证前端服务能否正常启动
      【修改方式】执行验证操作
      【相关依赖】无
      【修改内容】
         - 运行 `npm run dev`
         - 检查浏览器是否能访问 http://localhost:5173
         - 验证登录页面是否正常显示

  - [x] 10.3 测试用户登录功能
      【目标对象】`src/frontend/` 和 `src/backend/`
      【修改目的】验证用户登录流程是否正常工作
      【修改方式】执行功能测试
      【相关依赖】无
      【修改内容】
         - 使用初始管理员账户登录（username: admin、password: Admin@123）
         - 验证登录成功，Token 是否返回
         - 验证用户信息是否正确显示
         - 验证 Token 是否存储到 localStorage

  - [x] 10.4 测试登录失败场景
      【目标对象】`src/frontend/` 和 `src/backend/`
      【修改目的】验证登录失败和账户锁定机制
      【修改方式】执行功能测试
      【相关依赖】无
      【修改内容】
         - 测试密码错误场景（验证错误提示）
         - 连续输入错误密码5次，验证账户是否锁定
         - 验证锁定期间无法登录
         - 等待1小时后验证账户自动解锁

  - [x] 10.5 测试权限校验
      【目标对象】`src/backend/`
      【修改目的】验证权限中间件是否正常工作
      【修改方式】执行功能测试
      【相关依赖】无
      【修改内容】
         - 使用超级管理员账户访问所有接口，验证全部通过
         - 使用普通用户账户访问需要管理员权限的接口，验证返回403错误
         - 未登录访问受保护接口，验证返回401错误

  - [x] 10.6 测试会话管理
      【目标对象】`src/backend/`
      【修改目的】验证会话过期和登出功能
      【修改方式】执行功能测试
      【相关依赖】无
      【修改内容】
         - 登录后等待2小时，验证会话是否过期
         - 测试登出功能，验证 Token 是否失效
         - 测试多端登录场景（验证会话互不干扰）

  - [x] 10.7 测试用户管理功能
      【目标对象】`src/backend/`
      【修改目的】验证用户创建、禁用/启用功能
      【修改方式】执行功能测试
      【相关依赖】无
      【修改内容】
         - 使用超级管理员创建新用户
         - 验证新用户账户是否可用
         - 禁用用户账户
         - 使用被禁用账户登录，验证返回"账户已禁用"错误
         - 启用用户账户
         - 验证账户可正常登录
