# APP中台管理系统 - 功能清单

## 已完成功能 ✅

### 核心功能
- [x] 认证授权 - JWT登录、权限管理
- [x] APP管理 - 项目管理、密钥管理、重置AppSecret
- [x] 用户管理 - 用户列表、状态管理、搜索筛选
- [x] 配置中心 - 8个API接口（创建/编辑/删除/发布/历史/多环境支持）
- [x] 版本管理 - 8个API接口（创建/编辑/删除/发布/下线/灰度发布/强制更新）

### 核心架构
- [x] 模块接口定义 (core/module/module.go)
- [x] 模块注册中心 (core/module/registry.go)
- [x] 模块同步器 (core/module/sync.go)
- [x] 模块加载器 (modules/loader.go)
- [x] 主程序入口 (cmd/server/main.go)

### 数据库
- [x] 16张核心数据表已创建
- [x] 数据模型定义完成

## 后端API (已完成) ✅

### 模块1：存储服务 (file_storage) - 5个功能
- [x] 文件上传API (POST /api/v1/files)
- [x] 文件列表API (GET /api/v1/files)
- [x] 文件下载API (GET /api/v1/files/download/:id)
- [x] 文件删除API (DELETE /api/v1/files/:id)
- [x] 文件统计API (GET /api/v1/files/stats)

### 模块2：消息中心 (message_center) - 6个功能
- [x] 消息列表API (GET /api/v1/messages)
- [x] 发送消息API (POST /api/v1/messages)
- [x] 消息详情API (GET /api/v1/messages/:id)
- [x] 标记已读API (PUT /api/v1/messages/:id/read)
- [x] 批量发送API (POST /api/v1/messages/batch)
- [x] 消息统计API (GET /api/v1/messages/stats)

### 模块3：日志服务 (log_service) - 5个功能
- [x] 日志查询API (GET /api/v1/logs)
- [x] 日志上报API (POST /api/v1/logs)
- [x] 日志统计API (GET /api/v1/logs/stats)
- [x] 日志导出API (GET /api/v1/logs/export)
- [x] 日志清理API (DELETE /api/v1/logs/clean)

### 模块4：Push推送 (push_service) - 6个功能
- [x] 推送列表API (GET /api/v1/push)
- [x] 创建推送API (POST /api/v1/push)
- [x] 推送详情API (GET /api/v1/push/:id)
- [x] 发送推送API (POST /api/v1/push/:id/send)
- [x] 取消推送API (POST /api/v1/push/:id/cancel)
- [x] 推送统计API (GET /api/v1/push/stats)

### 模块5：数据埋点 (event_tracking) - 6个功能
- [x] 事件上报API (POST /api/v1/events)
- [x] 批量上报API (POST /api/v1/events/batch)
- [x] 事件列表API (GET /api/v1/events)
- [x] 事件统计API (GET /api/v1/events/stats)
- [x] 漏斗分析API (GET /api/v1/events/funnel)
- [x] 事件定义管理API (CRUD /api/v1/events/definitions)

### 模块6：监控告警 (monitor_service) - 5个功能
- [x] 监控指标API (GET /api/v1/monitor/metrics)
- [x] 上报指标API (POST /api/v1/monitor/metrics)
- [x] 告警管理API (CRUD /api/v1/monitor/alerts)
- [x] 监控统计API (GET /api/v1/monitor/stats)
- [x] 健康检查API (GET /api/v1/monitor/health)

## 前端页面 (待开发) 🔄

### 模块管理页面
- [ ] 存储服务管理页面
- [ ] 消息中心管理页面
- [ ] 日志服务查询页面
- [ ] Push推送管理页面
- [ ] 数据埋点分析页面
- [ ] 监控告警看板页面

### 系统集成
- [ ] 所有模块前端页面集成到侧边栏
- [ ] 移动端响应式适配
- [ ] 创建checkpoint推送到Git


## Bug修复

- [x] 修复前端登录“请求失败”错误

- [x] 优化APP管理页面UI，实现卡片式APP列表展示
- [x] 修复创建APP对话框中模块列表为空的问题
- [x] 将项目代码推送到GitHub仓库永久保存
- [ ] 完善APP详情页面（进入配置后的页面）
  - [ ] 顶部Tab：工作台 + 配置中心
  - [ ] 工作台：左侧边栏功能菜单 + 主内容区差异化功能
  - [ ] 配置中心：显示该APP选择的模块配置项


## 数据库迁移到Manus平台
- [x] 获取Manus数据库连接信息
- [x] 在Manus数据库中创建所需表结构（apps、app_modules、module_templates、admins等）
- [x] 修改后端配置使用Manus数据库
- [x] 测试后端连接和功能
- [x] 推送代码到GitHub


## 根据视频优化APP配置中心
- [x] 修复前端构建和服务器问题（使用SPA服务器）
- [x] 修复CORS跨域问题
- [x] 修复数据库连接问题（使用本地MySQL）
- [x] 登录功能正常
- [ ] 重新设计APP详情页面左侧边栏导航（按模块分组展开/收起）
- [ ] 实现APP概览页面（统计卡片、APP信息、已启用模块）
- [ ] 实现用户管理配置（登录配置、用户信息管理、实名认证、账号注销）
- [ ] 实现支付配置（安全验证、限额控制、回调配置）
- [ ] 实现短信配置（服务商配置、验证码配置、通知配置）
- [ ] 实现日志服务配置（基础配置、上报配置）
- [ ] 添加配置数据的后端API
- [ ] 添加配置数据的数据库表
