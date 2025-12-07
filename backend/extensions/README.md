# Extensions 扩展功能模块

本目录包含从 EtrxLite 集成的扩展功能模块。

---

## 📦 扩展模块列表

### 1. EtrxTable - 多维表格系统

**目录**: `etrxtable/`

**功能**:
- 动态表结构创建和管理
- 多种列类型支持（13 种）
- 表间关联（TableLink）
- 多视图支持（Grid/Kanban/Gantt/Calendar）
- 实时协作编辑
- 权限管理

**已实现**:
- ✅ 15 个 GORM 模型
- ✅ Repository 层（3 个仓库类）
- ✅ Service 层（完整业务逻辑）
- ✅ Handler 层（15+ API 端点）

**API 端点**:
```
GET    /api/v1/extensions/etrxtable/tables
POST   /api/v1/extensions/etrxtable/tables
GET    /api/v1/extensions/etrxtable/tables/:id
PUT    /api/v1/extensions/etrxtable/tables/:id
DELETE /api/v1/extensions/etrxtable/tables/:id
GET    /api/v1/extensions/etrxtable/tables/:id/columns
POST   /api/v1/extensions/etrxtable/tables/:id/columns
GET    /api/v1/extensions/etrxtable/tables/:id/rows
POST   /api/v1/extensions/etrxtable/tables/:id/rows
...
```

---

### 2. ER-Diagram - ER 图编辑器

**目录**: `er-diagram/`

**功能**:
- 可视化数据库设计
- SQL/DBML 导入导出
- 数据库同步
- 同步到 EtrxTable

**已实现**:
- ✅ 1 个 GORM 模型
- ⏳ Repository 层（待实现）
- ⏳ Service 层（待实现）
- ⏳ Handler 层（待实现）

**计划 API 端点**:
```
GET    /api/v1/extensions/er-diagram/diagrams
POST   /api/v1/extensions/er-diagram/diagrams
GET    /api/v1/extensions/er-diagram/diagrams/:id
PUT    /api/v1/extensions/er-diagram/diagrams/:id
DELETE /api/v1/extensions/er-diagram/diagrams/:id
GET    /api/v1/extensions/er-diagram/diagrams/:id/export/sql
POST   /api/v1/extensions/er-diagram/diagrams/:id/sync
```

---

### 3. Code-Generation - 代码生成系统

**目录**: `code-generation/`

**功能**:
- 从表格定义生成代码
- 支持多种语言（Go/GORM, Python/FastAPI, Java/Spring Boot）
- 代码模板管理
- 生成任务管理

**已实现**:
- ⏳ 模型层（待实现）
- ⏳ Service 层（待实现）
- ⏳ Handler 层（待实现）

**计划 API 端点**:
```
POST   /api/v1/extensions/code-generation/generate
GET    /api/v1/extensions/code-generation/tasks/:id
GET    /api/v1/extensions/code-generation/tasks/:id/files
```

---

### 4. MAS - 多智能体系统

**目录**: `mas/`

**功能**:
- Agent 管理和配置
- 任务调度和执行
- 会话管理和协调
- 记忆系统（短期/长期）
- 性能指标监控
- Agent 协作

**已实现**:
- ✅ 14 个 GORM 模型
- ⏳ Repository 层（待实现）
- ⏳ Service 层（待实现）
- ⏳ Handler 层（待实现）

**计划 API 端点**:
```
GET    /api/v1/extensions/mas/agents
POST   /api/v1/extensions/mas/agents
GET    /api/v1/extensions/mas/agents/:id
PUT    /api/v1/extensions/mas/agents/:id
DELETE /api/v1/extensions/mas/agents/:id
POST   /api/v1/extensions/mas/sessions
POST   /api/v1/extensions/mas/tasks
GET    /api/v1/extensions/mas/agents/:id/memories
GET    /api/v1/extensions/mas/metrics
```

---

## 🏗️ 架构设计

所有扩展模块遵循统一的分层架构：

```
extension-module/
├── models/           # GORM 数据模型
├── repository/       # 数据访问层
├── service/          # 业务逻辑层
├── handler/          # HTTP 处理器层
└── README.md         # 模块文档
```

### 分层职责

#### Model 层
- 定义 GORM 数据模型
- 数据库字段映射
- 关系定义

#### Repository 层
- 数据库 CRUD 操作
- 查询构建
- 事务处理

#### Service 层
- 业务逻辑实现
- 数据验证
- 业务规则

#### Handler 层
- HTTP 请求处理
- 参数验证
- 响应格式化

---

## 🔧 使用示例

### 注册扩展路由

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/coze-dev/coze-studio/backend/extensions/etrxtable/handler"
)

func registerExtensions(router *gin.Engine, db *gorm.DB) {
    api := router.Group("/api/v1")
    
    // 注册多维表格路由
    tableHandler := handler.NewTableHandler(db)
    tableHandler.RegisterRoutes(api)
    
    // 注册其他扩展...
}
```

### 使用扩展服务

```go
import (
    "github.com/coze-dev/coze-studio/backend/extensions/etrxtable/service"
)

func example(db *gorm.DB) {
    tableService := service.NewTableService(db)
    
    // 创建表格
    table := &models.TableBody{
        Name:        "项目任务表",
        WorkspaceID: 1,
        OwnerID:     1,
    }
    err := tableService.CreateTable(ctx, table)
}
```

---

## 📋 开发规范

### 统一返回格式

所有 API 使用统一的返回格式：

```go
c.JSON(200, gin.H{
    "code":    20000,
    "success": true,
    "message": "操作成功",
    "data": gin.H{
        "items": items,
        "pagination": gin.H{
            "page": 1,
            "page_size": 10,
            "total": 100,
            "total_pages": 10,
        },
    },
})
```

### 错误处理

```go
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "code":    50000,
        "success": false,
        "message": "操作失败: " + err.Error(),
        "data":    nil,
    })
    return
}
```

---

## 🔗 依赖关系

```
Handler → Service → Repository → Model
   ↓         ↓          ↓
  Gin     Business    GORM
         Logic
```

---

## 📚 参考文档

- [INTEGRATION_DEVELOPMENT_PLAN.md](../../INTEGRATION_DEVELOPMENT_PLAN.md) - 融合开发计划
- [.cursor/rules/coze-super.mdc](../../.cursor/rules/coze-super.mdc) - 开发规范
- [backend/README.md](../README.md) - 后端开发文档

---

**注意**: 扩展功能遵循松耦合原则，可以独立开发、测试和部署。

