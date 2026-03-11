# 生命科学论文管理数据库系统 - 开发任务规划

**文档版本**：v1.0  
**创建时间**：2026-03-11  
**最后更新**：2026-03-11  

---

## 开发阶段划分

### 阶段 1：基础设施搭建（预计 2 周）
- 前端项目初始化
- 后端项目初始化
- 数据库设计与建表
- 基础框架配置

### 阶段 2：核心功能开发（预计 6 周）
- 用户认证与权限管理
- 论文信息管理
- 审核流程管理
- 论文归档管理

### 阶段 3：辅助功能开发（预计 4 周）
- 查询检索功能
- 统计分析功能
- 数据导出功能
- 消息通知功能
- 系统配置管理

### 阶段 4：集成测试与优化（预计 2 周）
- 集成测试
- 性能优化
- 安全加固
- 部署上线

---

## 开发任务清单

- [x] 1. 前端项目初始化与基础架构搭建
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：创建 React + TypeScript + Vite 前端项目，配置 Ant Design UI 框架、Zustand 状态管理、React Router 路由、Axios HTTP 客户端，建立项目目录结构
  - 复杂度：简单
  - 对应需求：spec.md 中的 NFR-024、NFR-025
  - 请参考需求文档和设计文档规划功能实现提案

- [x] 2. 后端项目初始化与基础架构搭建
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：创建 Spring Boot 3.x 后端项目，配置 MyBatis-Plus、HikariCP、SpringDoc OpenAPI，建立 Controller-Service-Repository 分层架构
  - 复杂度：简单
  - 对应需求：spec.md 中的 NFR-001、NFR-002、NFR-003
  - 请参考需求文档和设计文档规划功能实现提案

- [x] 3. PostgreSQL 数据库设计与建表
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：根据 ER 图创建 users、papers、authors、projects、journals、review_logs、archives、notifications、operation_logs 等核心表，建立索引和约束
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-001 到 FR-090、TC-003
  - 请参考需求文档和设计文档规划功能实现提案

- [x] 4. Redis 缓存配置与集成
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：配置 Redis 连接池，实现用户会话缓存、热点数据缓存、系统配置缓存，设计缓存 Key 规范和 TTL 策略
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-005、NFR-001、NFR-002
  - 请参考需求文档和设计文档规划功能实现提案

- [x] 5. 统一响应格式与异常处理机制
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现 ApiResponse 统一响应类、PageResult 分页结果类、ErrorCode 错误码枚举、GlobalExceptionHandler 全局异常处理器
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-001、FR-003、NFR-016
  - 请参考需求文档和设计文档规划功能实现提案

- [x] 6. 用户登录认证功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现用户名密码登录、密码加密验证、JWT Token 生成与验证、登录失败次数统计、账户锁定机制
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-001、FR-002、FR-003、FR-006、NFR-009、NFR-010
  - 请参考需求文档和设计文档规划功能实现提案

- [x] 7. RBAC 权限控制体系
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现 7 类角色定义、权限矩阵配置、AuthFilter 认证过滤器、PermissionEvaluator 权限评估器、权限注解与 AOP 切面
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-004、NFR-011、NFR-018、NFR-019
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 8. 用户会话管理
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现 Redis 存储用户会话、会话有效期管理（2 小时）、会话过期检测、多端登录支持
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-005、NFR-013
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 9. 用户管理功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现用户 CRUD 操作、用户角色分配、账户禁用/启用、密码重置、用户信息查询
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-007、FR-074、FR-075
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 10. 论文信息录入功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现论文表单组件、字段格式校验（DOI、日期、ISSN）、草稿保存、论文 ID 自动生成、重复论文检测
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-008、FR-009、FR-010、FR-011、FR-014、FR-015、FR-016
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 11. 论文附件上传功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现文件上传组件、文件大小限制（100MB）、文件类型校验、本地存储/MinIO 存储、文件命名规则、断点续传
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-012、NFR-004、NFR-007
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 12. 期刊信息管理与自动填充
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现期刊 CRUD 操作、期刊模糊搜索、期刊信息自动填充、JCR 数据库对接接口预留
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-013、FR-077、FR-078、ER-001
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 13. 批量导入论文功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现 Excel 模板下载、EasyExcel 数据解析、批量数据校验、错误标注与修正、导入进度显示、断点续传
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-017、FR-018、FR-019、NFR-008
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 14. 作者信息管理功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现作者列表组件、作者信息 CRUD、人员库选择自动填充、作者类型互斥选择、作者排名自动调整
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-020、FR-021、FR-022、FR-023、FR-024、FR-025
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 15. 课题信息管理功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现课题 CRUD 操作、课题库选择自动填充、课题审核流程、已关联论文课题保护
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-026、FR-027、FR-028、FR-029、FR-030
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 16. 论文提交审核功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现提交审核按钮、状态变更（待业务审核）、提交人和时间记录、邮件通知配置
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-031、FR-032
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 17. 业务审核功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现待审核论文列表、审核页面展示、审核通过/驳回操作、驳回原因必填、审核记录保存、审核时限管理（3 个工作日）
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-033、FR-034、FR-035、FR-036、FR-037、FR-038
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 18. 政工审核功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现待政工审核论文列表、政工审核页面、审核通过/驳回操作、业务审核意见展示、审核权限隔离
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-039、FR-040、FR-041、FR-042、FR-043、FR-044、NFR-019
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 19. 审核状态机与流程控制
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现 PaperStatus 状态枚举、状态流转控制、状态权限校验、审核流程可视化
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-031、FR-035、FR-041
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 20. 论文归档功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现自动归档触发、归档编号生成（年份 + 论文 ID+ 随机 3 位）、归档记录保存、归档状态管理
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-056、FR-057
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 21. 归档论文分类管理
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现按年份分类、按收录类型分类、按课题分类、按作者分类、隐藏状态设置
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-058、FR-059、FR-060、FR-061、FR-063
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 22. 归档论文修改与二次审核
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现修改申请提交、二次审核流程、修改痕迹记录、修改前后对比
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-062
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 23. 多维度查询检索功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现组合查询表单、动态 SQL 构建、模糊查询、精确查询、查询条件保存
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-045、FR-049、FR-050、FR-051、FR-052
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 24. 查询结果展示与分页
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现数据表格组件、结果排序、分页浏览（每页 20 条）、无结果提示、论文详情查看
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-046、FR-047、FR-048、FR-054、FR-055
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 25. 查询权限控制
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现审核通过论文公开查询、我的论文功能、权限过滤、数据范围控制
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-053、NFR-011
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 26. 统计分析功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现基础指标统计（论文总数、年份分布、收录类型、平均影响因子、引用次数）、按作者/课题/部门统计
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-064、FR-065、FR-066、FR-067
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 27. 统计图表展示
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：集成 ECharts 图表库、实现柱状图/折线图/饼图、图表切换、图表导出（PNG/PDF）
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-068、FR-069
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 28. 数据导出功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现 Excel 导出（自定义字段）、Word/PDF导出（单篇论文）、统计结果导出、导出权限控制
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-070、FR-071、FR-072、FR-073
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 29. 消息通知功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现系统内消息、邮件通知（SMTP）、通知模板、审核状态变更通知、时限提醒通知、密码重置通知
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-083、FR-084、FR-085、FR-086、FR-087、ER-003
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 30. 系统配置管理功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现系统参数配置（附件大小、审核时限、影响因子规则等）、配置 CRUD、配置缓存刷新
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-079、FR-080
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 31. 操作日志管理功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现操作日志自动记录（登录、CRUD、审核、导出）、日志查询、日志不可修改保护、日志保留策略
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-081、FR-082、NFR-020、CR-004
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 32. 历史数据迁移功能
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现历史数据 Excel 导入、数据格式校验、错误修正、归档状态标记、数据完整性校验
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-088、FR-089、FR-090
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 33. Web of Science 对接接口
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现 WoS API 对接、引用次数自动更新、他引次数自动更新、定时任务调度
  - 复杂度：复杂
  - 对应需求：spec.md 中的 ER-002、SC-010
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 34. 前端登录与认证页面
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现登录页面、密码复杂度校验提示、登录失败提示、账户锁定提示、会话过期跳转
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-001、FR-002、FR-003、FR-005
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 35. 前端权限路由与菜单
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现动态路由配置、基于角色的菜单显示、权限路由守卫、无权限提示页面
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-004、NFR-018
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 36. 前端论文录入与编辑页面
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现论文表单页面、字段校验、草稿保存、附件上传、作者管理、课题关联、期刊搜索
  - 复杂度：复杂
  - 对应需求：spec.md 中的 FR-008、FR-014、FR-020、FR-026
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 37. 前端审核页面
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现待审核列表页面、审核详情页面、审核通过/驳回操作、驳回原因富文本编辑
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-033、FR-034、FR-035、FR-036、FR-039、FR-040
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 38. 前端查询检索页面
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现查询表单页面、组合条件输入、查询结果表格、分页组件、排序功能、详情查看
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-045、FR-046、FR-047、FR-048
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 39. 前端统计分析页面
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现统计 Dashboard 页面、图表展示组件、统计维度切换、图表导出功能
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-064、FR-068、FR-069
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 40. 前端系统管理页面
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现用户管理页面、课题管理页面、期刊管理页面、系统配置页面、日志查询页面
  - 复杂度：简单
  - 对应需求：spec.md 中的 FR-074、FR-076、FR-077、FR-079、FR-082
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 41. 安全加固与防护
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现 SQL 注入防护、XSS 攻击防护、CSRF 防护、敏感数据加密存储、HTTPS 配置
  - 复杂度：复杂
  - 对应需求：spec.md 中的 NFR-012、NFR-013、NFR-016、CR-003、CR-005
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 42. 性能优化
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：实现数据库查询优化、索引优化、缓存优化、前端懒加载、代码分割、响应时间优化
  - 复杂度：复杂
  - 对应需求：spec.md 中的 NFR-001、NFR-002、NFR-003、NFR-004、NFR-005
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 43. 集成测试
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：编写单元测试、接口测试、端到端测试、验收测试场景执行、Bug 修复
  - 复杂度：复杂
  - 对应需求：spec.md 中的 SC-001 到 SC-015
  - 请参考需求文档和设计文档规划功能实现提案

- [ ] 44. 部署与上线
  - 整体需求文档位置：`.cospec/spec/paper-management/spec.md`
  - 整体设计文档路径：`.cospec/spec/paper-management/project.md`
  - 实现功能：配置生产环境、数据库初始化、数据备份策略配置、内网部署、系统监控配置
  - 复杂度：简单
  - 对应需求：spec.md 中的 TC-001、TC-002、TC-005、NFR-030、NFR-031
  - 请参考需求文档和设计文档规划功能实现提案

---

## 任务依赖关系与执行顺序

### 关键路径任务

```
阶段 1：基础设施搭建
  任务 1（前端初始化）→ 任务 34、35、36、37、38、39、40
  任务 2（后端初始化）→ 任务 3、4、5
  任务 3（数据库建表）→ 任务 6、7、10、11、12、13、14、15...
  任务 4（Redis 配置）→ 任务 8
  任务 5（统一响应）→ 所有后端任务

阶段 2：核心功能开发
  任务 6（登录认证）→ 任务 7、8、9
  任务 7（RBAC 权限）→ 任务 16、17、18、25、29
  任务 10（论文录入）→ 任务 11、12、13、14、15
  任务 16（提交审核）→ 任务 17、18
  任务 17（业务审核）→ 任务 18、19
  任务 18（政工审核）→ 任务 20
  任务 20（归档）→ 任务 21、22

阶段 3：辅助功能开发
  任务 23（查询检索）→ 任务 24、25
  任务 26（统计分析）→ 任务 27
  任务 28（数据导出）
  任务 29（消息通知）
  任务 30（系统配置）→ 任务 31
  任务 32（历史迁移）
  任务 33（WoS 对接）

阶段 4：集成测试与优化
  任务 41（安全加固）
  任务 42（性能优化）
  任务 43（集成测试）
  任务 44（部署上线）
```

### 可并行执行的任务

**阶段 1 并行任务**：
- 任务 1（前端初始化）与 任务 2（后端初始化）与 任务 3（数据库建表）可并行

**阶段 2 并行任务**：
- 任务 11（附件上传）与 任务 12（期刊管理）与 任务 13（批量导入）可并行
- 任务 14（作者管理）与 任务 15（课题管理）可并行

**阶段 3 并行任务**：
- 任务 23（查询检索）与 任务 26（统计分析）与 任务 28（数据导出）可并行
- 任务 29（消息通知）与 任务 30（系统配置）可并行
- 任务 34-40（前端页面开发）可与后端任务并行

---

## 关键里程碑

| 里程碑 | 完成标志 | 预计时间 |
|--------|----------|----------|
| M1：基础设施完成 | 任务 1-5 完成，前后端项目可运行，数据库可连接 | 第 2 周末 |
| M2：认证权限完成 | 任务 6-9 完成，用户可登录，权限控制生效 | 第 3 周末 |
| M3：核心功能完成 | 任务 10-22 完成，论文录入、审核、归档全流程贯通 | 第 8 周末 |
| M4：辅助功能完成 | 任务 23-33 完成，查询、统计、导出、通知功能可用 | 第 12 周末 |
| M5：前端页面完成 | 任务 34-40 完成，所有页面可访问可操作 | 第 12 周末 |
| M6：测试优化完成 | 任务 41-43 完成，通过验收测试，性能达标 | 第 13 周末 |
| M7：系统上线 | 任务 44 完成，系统部署到生产环境 | 第 14 周末 |

---

## 验收标准（DoD）

### 通用完成定义

每个任务完成后必须满足以下条件：

1. **代码完成**：功能代码编写完成，符合架构设计规范
2. **单元测试**：核心逻辑单元测试覆盖率≥80%
3. **代码审查**：代码通过团队审查，无重大质量问题
4. **文档更新**：API 文档、数据库文档同步更新
5. **无阻塞 Bug**：无 P0、P1 级别的缺陷

### 关键里程碑验收标准

**M1 验收标准**：
- 前端项目可启动，访问 localhost 显示欢迎页面
- 后端项目可启动，Swagger UI 可访问
- 数据库所有表创建成功，可执行 CRUD 操作
- Redis 连接成功，可执行缓存操作

**M2 验收标准**：
- 用户可使用正确账号密码登录
- 登录失败 5 次后账户锁定
- 不同角色登录后看到的菜单不同
- 无权限访问时自动跳转或提示

**M3 验收标准**：
- 用户可录入论文，字段校验生效
- 论文可提交审核，状态正确变更
- 业务审核员可审核论文，通过/驳回生效
- 政工审核员可审核论文，权限隔离生效
- 审核通过后论文自动归档

**M4 验收标准**：
- 多维度查询检索返回正确结果
- 统计图表展示正确数据
- 数据导出文件格式正确
- 通知邮件可发送成功

**M5 验收标准**：
- 所有页面可访问，无 404 错误
- 页面加载时间≤1 秒
- 表单提交、数据展示正常

**M6 验收标准**：
- 所有验收场景测试通过（SC-001 到 SC-015）
- 简单查询响应≤1 秒，复杂查询≤2 秒
- 50 用户并发无明显卡顿
- 安全扫描无高危漏洞

**M7 验收标准**：
- 系统部署到内网服务器
- 数据库备份策略生效
- 系统监控告警配置完成
- 用户手册、运维手册交付

---

## 风险与应对

| 风险 | 影响 | 概率 | 应对措施 |
|------|------|------|----------|
| JCR/WoS 接口对接困难 | 影响 FR-078、ER-002 | 中 | 预留手动导入功能作为备选方案 |
| 大文件上传性能问题 | 影响 NFR-004 | 中 | 采用分片上传、CDN 加速 |
| 复杂查询性能不达标 | 影响 NFR-003 | 中 | 提前进行 SQL 优化、索引优化 |
| 历史数据格式不规范 | 影响 FR-088 | 高 | 提供数据清洗工具，支持批量修正 |
| 内网部署环境限制 | 影响 TC-001 | 低 | 提前调研内网环境，准备离线包 |

---

**文档结束**
