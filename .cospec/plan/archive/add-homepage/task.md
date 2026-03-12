## 实施

- [x] 1.1 创建统计类型定义 `statistics.ts`
     【目标对象】`src/types/statistics.ts`
     【修改目的】定义统计相关类型
     【修改方式】新建文件
     【相关依赖】无（独立类型定义）
     【修改内容】
        - 定义 Statistics 接口（总数、待审核、已通过、草稿等）
        - 定义 RecentPaper 接口
        - 定义 PendingReview 接口
        - 定义 HomeData 接口，聚合所有首页数据

- [x] 1.2 创建统计服务 `statisticsService.ts`
     【目标对象】`src/services/statisticsService.ts`
     【修改目的】提供首页统计数据获取能力
     【修改方式】新建文件
     【相关依赖】
        - `src/services/api.ts` - HTTP 请求封装
        - `src/services/paperService.ts` - 论文服务（getPapers, getMyPapers）
        - `src/services/reviewService.ts` - 审核服务（getPendingBusinessReviews, getPendingPoliticalReviews）
        - `src/types/statistics.ts` - 统计类型定义
     【修改内容】
        - 创建 getStatistics 方法，调用 getPapers 获取论文列表并按状态统计
        - 创建 getPendingReviews 方法，获取待审核数据
        - 创建 getRecentPapers 方法，获取最近录入的论文（按时间倒序）
        - 创建 getHomeData 方法，聚合所有首页数据

- [x] 1.3 实现首页统计卡片区域
     【目标对象】`src/pages/home/HomePage.tsx`
     【修改目的】展示关键数据指标
     【修改方式】修改现有文件
     【相关依赖】
        - `src/services/statisticsService.ts` - 统计服务
        - `src/types/statistics.ts` - 统计类型定义
        - Ant Design: Card, Row, Col, Statistic, Spin
        - `src/stores/userStore.ts` - 获取用户角色信息
     【修改内容】
        - 导入 Ant Design 组件（Card, Row, Col, Statistic, Spin）
        - 创建统计卡片组件，展示论文总数、待审核、已通过、我的草稿
        - 使用不同颜色区分不同类型的统计卡片（总数: blue, 待审核: orange, 已通过: green, 草稿: gray）
        - 添加数据加载状态（Spin 组件包裹）
        - 处理空数据情况（显示 0 而非空白）

- [x] 1.4 实现快捷操作入口区域
     【目标对象】`src/pages/home/HomePage.tsx`
     【修改目的】提供快速操作入口
     【修改方式】修改现有文件
     【相关依赖】
        - React Router: useNavigate 路由跳转
        - Ant Design: Card, List, Button, Space, Icons
        - `src/stores/userStore.ts` - 获取用户角色（判断是否显示审核入口）
     【修改内容】
        - 使用 Card + List 组件展示快捷操作
        - 基础操作：录入论文、批量导入、我的论文（所有用户可见）
        - 审核操作：根据用户角色显示业务审核/政工审核入口
        - 使用图标增强视觉效果（PlusOutlined, UploadOutlined, FileTextOutlined, CheckCircleOutlined）
        - 使用 useAuth 或 userStore 获取当前用户角色

- [x] 1.5 实现最近动态区域
     【目标对象】`src/pages/home/HomePage.tsx`
     【修改目的】展示最近活动和待办事项
     【修改方式】修改现有文件
     【相关依赖】
        - React Router: useNavigate 路由跳转
        - Ant Design: Card, List, Tag, Avatar, Typography
        - `src/services/statisticsService.ts` - 获取最近论文数据
        - Day.js: 格式化时间显示
     【修改内容】
        - 使用 Card + List 组件展示最近论文（最多显示 5 条）
        - 显示论文标题、状态标签、提交时间
        - 点击跳转到论文详情页（/papers/:id）
        - 根据用户角色显示待审核任务列表
        - 处理空数据情况（显示"暂无最近动态"）

- [x] 1.6 添加欢迎区域和整体布局优化
     【目标对象】`src/pages/home/HomePage.tsx`
     【修改目的】提升用户体验
     【修改方式】修改现有文件
     【相关依赖】
        - Ant Design: Row, Col, Typography, Avatar, Space
        - `src/stores/userStore.ts` - 获取当前用户信息
     【修改内容】
        - 添加欢迎区域，显示"欢迎回来，{用户名}"
        - 显示当前日期
        - 优化整体布局，使用 Row/Col 实现响应式栅格
        - 桌面端：统计卡片 4 列，快捷操作和最近动态各占 12 列
        - 平板端：统计卡片 2 列，快捷操作和最近动态各占 12 列
        - 手机端：统计卡片 1 列，快捷操作和最近动态各占 24 列
        - 添加适当的间距（gutter: [16, 16]）和内边距

- [x] 1.7 添加首页加载状态和错误处理
     【目标对象】`src/pages/home/HomePage.tsx`
     【修改目的】保证用户体验
     【修改方式】修改现有文件
     【相关依赖】
        - Ant Design: Alert, Button, Result
        - React: useState, useEffect
     【修改内容】
        - 添加数据加载中的 loading 状态（全局 loading）
        - 添加数据获取失败的错误提示（Alert 组件）
        - 错误时显示重试按钮
        - 优化首次加载体验（骨架屏或占位）
