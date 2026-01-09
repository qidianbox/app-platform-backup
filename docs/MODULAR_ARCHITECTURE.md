# 编译时模块化架构方案

**版本**: v9.0 (已实现)  
**作者**: Manus AI  
**日期**: 2026-01-08

---

## 实现状态

✅ **已完成实现** - 本方案已在APP中台系统中成功落地。

### 已实现的核心组件

| 组件 | 文件路径 | 状态 |
|:---|:---|:---|
| 模块接口定义 | `/backend/core/module/module.go` | ✅ 完成 |
| 模块注册中心 | `/backend/core/module/registry.go` | ✅ 完成 |
| 模块同步器 | `/backend/core/module/sync.go` | ✅ 完成 |
| 模块加载器 | `/backend/modules/loader.go` | ✅ 完成 |
| 主程序入口 | `/backend/cmd/server/main.go` | ✅ 完成 |

### 已实现的功能模块

| 模块代码 | 模块名称 | 功能数量 | 状态 |
|:---|:---|:---|:---|
| `user_management` | 用户管理 | 4 | ✅ 完成 |
| `message_center` | 消息中心 | 4 | ✅ 完成 |
| `push_service` | 推送服务 | 4 | ✅ 完成 |
| `event_tracking` | 埋点服务 | 5 | ✅ 完成 |
| `log_service` | 日志服务 | 4 | ✅ 完成 |
| `monitor_service` | 监控服务 | 4 | ✅ 完成 |
| `file_storage` | 文件存储 | 4 | ✅ 完成 |
| `config_management` | 配置管理 | 5 | ✅ 完成 |
| `version_management` | 版本管理 | 5 | ✅ 完成 |

**总计**: 9个模块，39个功能点

---

## 1. 核心思想：编译时组装，单一应用

本方案旨在为内部开发团队提供一个高效、安全、易于协作的模块化开发模式。我们放弃在运行时动态加载插件的复杂模式，转而采用**编译时模块化**的策略。

- **开发时分离**：不同团队在各自独立的Go包（或Git仓库）中开发功能模块。
- **编译时整合**：在应用编译阶段，通过Go的包管理机制将所有需要的模块整合进来，最终编译成一个单一的、完整的二进制文件进行部署。

这种方式兼具**模块化开发**的灵活性和**单体应用**的性能与部署简便性，是内部团队协作的最佳实践。

## 2. 模块接口规范 (Interface)

为了实现模块的"可插拔"，我们定义一个所有模块都必须实现的统一接口。这个接口是模块与主框架之间的契约。

**文件路径**: `/backend/core/module/module.go`

```go
package module

import "github.com/gin-gonic/gin"

// Meta 包含了模块的基本元数据
type Meta struct {
    Code        string // 模块唯一标识，例如 "user_management"
    Name        string // 人类可读的名称，例如 "用户管理"
    Description string // 模块功能描述
    Icon        string // 模块图标
    SortOrder   int    // 排序顺序
}

// Function 定义了一个具体的功能点，对应数据库中的一条记录
type Function struct {
    Code         string                 // 功能的唯一标识
    Name         string                 // 功能名称
    Description  string                 // 功能描述
    Type         string                 // 功能类型: "active" 或 "passive"
    ConfigSchema map[string]interface{} // 功能的JSON Schema配置
    Dependencies []string               // 依赖的其他功能Code列表
    SortOrder    int                    // 排序顺序
}

// Module 是所有功能模块必须实现的接口
type Module interface {
    Meta() Meta
    RegisterRoutes(router *gin.RouterGroup)
    GetFunctions() []Function
    Init() error
}
```

## 3. 模块注册机制

我们利用Go语言的`init()`函数特性，实现模块的自动注册。

**文件路径**: `/backend/core/module/registry.go`

```go
package module

var (
    modules   = make(map[string]Module)
    initOrder []string
)

// Register 用于注册一个模块实例
func Register(m Module) {
    meta := m.Meta()
    if _, exists := modules[meta.Code]; exists {
        panic("module already registered: " + meta.Code)
    }
    modules[meta.Code] = m
    initOrder = append(initOrder, meta.Code)
}

// GetAllModules 返回所有已注册的模块
func GetAllModules() []Module {
    all := make([]Module, 0, len(modules))
    for _, code := range initOrder {
        if m, exists := modules[code]; exists {
            all = append(all, m)
        }
    }
    return all
}
```

### 模块实现示例

**文件路径**: `/backend/modules/user/module.go`

```go
package user

import (
    "app-platform-backend/core/module"
    "app-platform-backend/internal/api/v1/user"
    "github.com/gin-gonic/gin"
)

type userModule struct {
    *module.BaseModule
}

func init() {
    module.Register(&userModule{
        BaseModule: module.NewBaseModule(
            module.Meta{
                Code:        "user_management",
                Name:        "用户管理",
                Description: "提供用户信息管理、状态管理、统计分析等功能",
                Icon:        "user",
                SortOrder:   100,
            },
            []module.Function{
                {
                    Code:        "user_list",
                    Name:        "用户列表",
                    Description: "查看和搜索所有用户",
                    Type:        "active",
                    SortOrder:   1,
                },
                // ... 更多功能
            },
        ),
    })
}

func (m *userModule) RegisterRoutes(router *gin.RouterGroup) {
    userGroup := router.Group("/users")
    {
        userGroup.GET("", user.List)
        userGroup.GET("/stats", user.Stats)
        userGroup.GET("/:id", user.Detail)
        userGroup.PUT("/:id/status", user.UpdateStatus)
    }
}
```

## 4. 目录结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # 应用主入口
├── core/                     # 核心框架代码
│   └── module/
│       ├── module.go         # 模块接口定义
│       ├── registry.go       # 模块注册中心
│       └── sync.go           # 模块同步器
├── modules/                  # 功能模块
│   ├── loader.go             # 模块加载器
│   ├── user/                 # 用户管理模块
│   ├── message/              # 消息中心模块
│   ├── push/                 # 推送服务模块
│   ├── event/                # 埋点服务模块
│   ├── log/                  # 日志服务模块
│   ├── monitor/              # 监控服务模块
│   ├── file/                 # 文件存储模块
│   ├── config/               # 配置管理模块
│   └── version/              # 版本管理模块
└── internal/                 # 内部实现
    ├── api/                  # API处理器
    ├── model/                # 数据模型
    └── ...
```

## 5. 应用启动流程

主程序启动时的模块化初始化流程：

```go
// main.go

import (
    "app-platform-backend/core/module"
    _ "app-platform-backend/modules" // 导入触发所有模块的init()注册
)

func main() {
    // 1. 初始化所有模块
    module.InitAllModules()
    
    // 2. 同步模块功能到数据库
    syncer := module.NewSyncer(database.GetDB())
    syncer.SyncModulesToDB()
    
    // 3. 注册所有模块的路由
    modules := module.GetAllModules()
    for _, m := range modules {
        m.RegisterRoutes(authGroup)
    }
    
    // 4. 启动服务器
    r.Run(":8080")
}
```

## 6. 数据库同步

模块同步器会在应用启动时，自动将所有模块的功能同步到 `module_templates` 表：

| 字段名 | 数据类型 | 描述 |
|:---|:---|:---|
| `id` | INT | 主键 |
| `module_code` | VARCHAR | 功能唯一标识 |
| `module_name` | VARCHAR | 功能名称 |
| `description` | TEXT | 功能描述 |
| `config_schema` | JSON | 配置Schema |
| `source_module` | VARCHAR | 来源模块Code |
| `function_type` | VARCHAR | 功能类型 (active/passive) |
| `is_active` | BOOLEAN | 是否启用 |

## 7. API端点

### 健康检查
```
GET /health
```
返回：
```json
{
    "status": "ok",
    "modules": 9,
    "architecture": "modular"
}
```

### 模块信息
```
GET /api/v1/system/modules
```
返回所有已注册模块及其功能的详细信息。

## 8. 开发新模块指南

### 步骤1：创建模块目录
```bash
mkdir -p backend/modules/your_module
```

### 步骤2：实现模块接口
创建 `backend/modules/your_module/module.go`：

```go
package your_module

import (
    "app-platform-backend/core/module"
    "github.com/gin-gonic/gin"
)

type yourModule struct {
    *module.BaseModule
}

func init() {
    module.Register(&yourModule{
        BaseModule: module.NewBaseModule(
            module.Meta{
                Code:        "your_module",
                Name:        "你的模块",
                Description: "模块描述",
                Icon:        "icon-name",
                SortOrder:   1000,
            },
            []module.Function{
                // 定义功能点
            },
        ),
    })
}

func (m *yourModule) RegisterRoutes(router *gin.RouterGroup) {
    // 注册路由
}
```

### 步骤3：在loader.go中导入
编辑 `backend/modules/loader.go`，添加导入：
```go
import (
    _ "app-platform-backend/modules/your_module"
)
```

### 步骤4：重新编译
```bash
go build -o app-platform-server ./cmd/server/main.go
```

启动后，新模块会自动注册并同步到数据库。

## 9. Git协作模式

### Git Submodule (推荐)
```bash
# 添加外部模块仓库
git submodule add <URL> plugins/external-module

# 更新子模块
git submodule update --init --recursive
```

### Go Workspace (本地开发)
```bash
# 创建工作区
go work init ./app-platform ./plugins/external-module
```

---

## 总结

本模块化架构方案已成功实现并投入使用，主要特点：

1. **编译时整合**：所有模块在编译时组装成单一二进制文件
2. **自动注册**：利用Go的`init()`机制实现模块自动注册
3. **自动同步**：启动时自动将模块功能同步到数据库
4. **标准接口**：所有模块实现统一的`Module`接口
5. **团队协作**：支持Git Submodule和Go Workspace协作模式
