# APP中台管理系统

**技术栈**: Vue 3 + Element Plus + Go + Gin + MySQL

统一的企业级应用管理平台，提供APP管理、用户管理、配置管理、监控告警等完整功能。

---

## 项目结构

```
app-platform/
├── frontend/          # Vue 3 前端
│   ├── src/
│   │   ├── views/    # 页面组件
│   │   ├── components/ # 公共组件
│   │   ├── router/   # 路由配置
│   │   ├── store/    # 状态管理
│   │   ├── api/      # API接口
│   │   └── utils/    # 工具函数
│   └── package.json
├── backend/           # Go 后端
│   ├── cmd/          # 主程序入口
│   ├── core/         # 核心模块
│   ├── models/       # 数据模型
│   ├── modules/      # 业务模块
│   ├── middleware/   # 中间件
│   └── configs/      # 配置文件
├── docs/             # 文档
├── deploy/           # 部署配置
└── tests/            # 测试文件
```

---

## 功能模块

### 核心功能
- ✅ 认证授权 - JWT登录、权限管理
- ✅ APP管理 - 项目管理、密钥管理
- ✅ 用户管理 - 用户列表、状态管理
- ✅ 配置中心 - 多环境配置管理
- ✅ 版本管理 - 版本发布、灰度发布

### 11个业务模块
1. **配置管理** - 配置列表、创建/更新/发布配置、配置历史
2. **埋点服务** - 事件上报、批量上报、事件统计、漏斗分析
3. **文件存储** - 文件上传/下载、文件列表、存储统计
4. **日志服务** - 日志查询、上报、统计、导出、清理
5. **消息中心** - 发送消息、消息列表、消息模板、批量发送
6. **监控服务** - 上报指标、监控指标、告警管理、健康检查
7. **推送服务** - 创建推送、发送推送、推送统计、推送模板
8. **用户管理** - 用户列表、用户详情、状态管理、用户统计
9. **版本管理** - 版本列表、创建/发布/下线版本、更新检查
10. **WebSocket服务** - 实时连接、数据推送、告警推送
11. **审计日志** - 日志列表、审计统计、导出审计日志

---

## 快速开始

### 前端开发

```bash
cd frontend
pnpm install
pnpm dev
```

前端将运行在 `http://localhost:5173`

### 后端开发

```bash
cd backend
go mod download
go run cmd/server/main.go
```

后端将运行在 `http://localhost:8080`

### 数据库配置

数据库配置文件：`backend/configs/config.yaml`

```yaml
database:
  driver: mysql
  host: your-database-host
  port: 3306
  username: your-username
  password: your-password
  database: app_platform
  charset: utf8mb4
```

详细配置信息请查看 `DATABASE_INFO.md`

---

## 部署

### 生产构建

**前端**:
```bash
cd frontend
pnpm build
```

**后端**:
```bash
cd backend
go build -o server cmd/server/main.go
```

### 部署方案

项目支持多种部署方式：
- 阿里云SAE（Serverless应用引擎）
- 阿里云ACK（容器服务Kubernetes）
- 阿里云ECS（云服务器）
- Docker容器部署

详细部署文档请查看 `deploy/` 目录。

---

## 环境变量

### 前端环境变量

创建 `frontend/.env` 文件：

```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 后端环境变量

后端配置通过 `backend/configs/config.yaml` 管理，包括：
- 数据库配置
- Redis配置
- JWT配置
- CORS配置
- 文件上传配置

---

## 测试

### 前端测试

```bash
cd frontend
pnpm test
```

### 后端测试

```bash
cd backend
go test ./...
```

测试报告请查看 `tests/` 目录。

---

## 文档

- [备份说明](BACKUP_README.md) - 完整的备份和恢复指南
- [数据库信息](DATABASE_INFO.md) - 数据库配置和表结构
- [安全审计报告](SECURITY_AUDIT_REPORT.md) - 安全审计结果
- [性能分析报告](BENCHMARK_ANALYSIS_REPORT.md) - 性能测试结果
- [系统健壮性报告](SYSTEM_ROBUSTNESS_REPORT.md) - 系统稳定性分析

---

## 技术栈详情

### 前端
- **框架**: Vue 3.5+ (Composition API)
- **UI组件**: Element Plus 2.13+
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP客户端**: Axios
- **图表**: ECharts 5
- **构建工具**: Vite 5

### 后端
- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis
- **认证**: JWT

---

## 默认账号

- **用户名**: admin
- **密码**: admin123

⚠️ 生产环境请务必修改默认密码！

---

## 许可证

MIT License

---

## 联系方式

- **项目负责人**: [待填写]
- **技术支持**: [待填写]
- **GitHub**: https://github.com/qidianbox/app-platform-backup

---

**最后更新**: 2026年1月11日
