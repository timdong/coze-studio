# Coze Studio + EtrxLite 融合开发计划

## 📋 项目概述

本文档旨在制定一个在 **Coze Studio** 基础上，参考 **EtrxLite** 优势功能的融合开发计划。核心原则是**松耦合集成**，确保不影响未来 Coze Studio 的升级和维护。

## 🎯 目标

1. **保留 Coze Studio 核心优势**：成熟的 AI Agent 开发平台、工作流引擎、知识库系统
2. **引入 EtrxLite 优秀特性**：统一基础设施架构、多维表格系统、ER图编辑器、代码生成等
3. **实现松耦合集成**：通过插件化、适配器模式、独立服务等方式实现
4. **保持升级兼容性**：确保 Coze Studio 核心代码可以独立升级

---

## 📊 两个项目对比分析

### Coze Studio 核心优势

| 维度 | 技术栈/特性 | 说明 |
|------|-----------|------|
| **后端框架** | Hertz (字节跳动开源) | 高性能 HTTP 框架，基于 Netpoll |
| **ORM** | GORM (v1.25+) | 成熟的 ORM 框架，支持 MySQL/SQLite |
| **架构模式** | DDD (领域驱动设计) | 清晰的领域边界，易于扩展 |
| **前端架构** | React 18 + TypeScript + Rush Monorepo | 企业级 monorepo 管理 |
| **工作流引擎** | FlowGram 集成 | 成熟的可视化工作流编辑器 |
| **AI 能力** | Eino 框架集成 | 统一的 AI 模型抽象和实现 |
| **知识库** | 完整的 RAG 实现 | 文档处理、向量检索、智能问答 |
| **插件系统** | 完善的插件机制 | 支持第三方插件扩展 |
| **社区生态** | 活跃的开源社区 | 持续更新和维护 |

### EtrxLite 核心优势

| 维度 | 技术栈/特性 | 说明 |
|------|-----------|------|
| **统一基础设施** | Core/Platform/Internal 三层架构 | 统一日志、认证、缓存、配置、数据库管理 |
| **多维表格系统** | EtrxTable | 动态表结构、多视图、实时协作 |
| **ER 图编辑器** | 可视化数据库设计 | SQL/DBML 导入导出、数据库同步 |
| **代码生成系统** | 表格到代码自动生成 | 支持 Go/Ent、Python/FastAPI、Java/Spring |
| **MAS 多智能体系统** | Agent 管理、任务调度、会话协调 | 完整的 Agent 生命周期管理 |
| **统一日志系统** | Zap + Lumberjack | 结构化日志，零侵入性 |
| **统一缓存系统** | 内存 + Redis 多级缓存 | 高性能缓存策略 |
| **WebSocket 实时协作** | 完整的实时通信框架 | 多用户协同编辑 |

---

## 🏗️ 融合架构设计

### 总体架构原则

```
┌─────────────────────────────────────────────────────────────┐
│                    Coze Studio 核心层                        │
│  (保持不变，可独立升级)                                      │
│  - Gin 框架 (统一后)                                         │
│  - DDD 领域模型                                              │
│  - 工作流引擎                                                │
│  - 知识库系统                                                │
└─────────────────────────────────────────────────────────────┘
                          ↕ (适配器层)
┌─────────────────────────────────────────────────────────────┐
│                  EtrxLite 增强层 (插件化)                    │
│  (独立模块，通过接口集成)                                      │
│  - 统一基础设施 Core (backend/core/)                        │
│  - 多维表格系统 (backend/extensions/etrxtable/)            │
│  - ER 图编辑器 (backend/extensions/er-diagram/)            │
│  - 代码生成系统 (backend/extensions/code-generation/)      │
│  - MAS 多智能体系统 (backend/extensions/mas/)              │
└─────────────────────────────────────────────────────────────┘
```

### 三层集成策略

#### 1. **基础设施层集成** (Core Layer)

**目标**：将 EtrxLite 的统一基础设施能力引入 Coze Studio

**集成方式**：**适配器模式 + 可选依赖**

```
coze-studio/
├── backend/
│   ├── core/                     # EtrxLite 统一基础设施（核心层）
│   │   ├── logging/              # 统一日志系统
│   │   ├── cache/                # 统一缓存系统
│   │   ├── auth/                 # 统一认证系统
│   │   ├── config/               # 统一配置管理
│   │   └── database/             # 统一数据库管理
│   │
│   ├── infra/                    # Coze Studio 原有基础设施（保留）
│   │   ├── database/
│   │   ├── cache/
│   │   └── logger/
│   │
│   └── extensions/               # 扩展层（EtrxLite 业务功能）
│       ├── etrxtable/            # 多维表格系统
│       ├── er-diagram/           # ER 图编辑器
│       ├── code-generation/      # 代码生成系统
│       └── mas/                  # MAS 多智能体系统
```

**实现策略**：
- 创建适配器接口，将 EtrxLite 的 Core 能力封装为 Coze Studio 可用的接口
- 通过配置开关控制是否启用 EtrxLite 基础设施
- 保持 Coze Studio 原有基础设施不变，作为 fallback

**代码示例**：
```go
// backend/core/logging/adapter.go
package logging

import (
    "github.com/coze-dev/coze-studio/backend/pkg/logs"
    "go.uber.org/zap"
)

// LoggingAdapter 统一日志适配器
type LoggingAdapter struct {
    logger *zap.Logger
}

func (a *LoggingAdapter) Info(msg string, fields ...interface{}) {
    a.logger.Info(msg, fields...)
}

// 初始化统一日志系统
func InitLogger() logs.Logger {
    logger, _ := NewLogger(&Config{
        Level:  "info",
        Format: "json",
    })
    return &LoggingAdapter{logger: logger}
}
```

#### 1.3 Ent 到 GORM 迁移 (重要)

**目标**：将 EtrxLite 的 Ent Schema 迁移到 GORM Model

**实施步骤**：
1. **分析 Ent Schema**：梳理 EtrxLite 中的所有 Ent Schema 定义
2. **转换为 GORM Model**：将 Ent Schema 转换为 GORM 结构体
3. **迁移数据库操作**：将 Ent 的 CRUD 操作转换为 GORM 操作
4. **保持 API 兼容**：确保对外 API 接口不变

**文件结构**：
```
backend/extensions/etrxtable/models/  # 多维表格模型
├── table_body.go
├── table_column.go
└── table_row.go

backend/extensions/er_diagram/models/  # ER 图模型
└── er_diagram.go

backend/extensions/mas/models/  # MAS 模型
├── agent.go
└── task.go
```

**迁移示例**：

**Ent Schema (原)**：
```go
// etrxlite/platform/ent/schema/table_body.go
type TableBody struct {
    ent.Schema
}

func (TableBody) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id"),
        field.String("name"),
        field.Time("created_at"),
    }
}
```

**GORM Model (迁移后)**：
```go
// backend/extensions/etrxtable/models/table_body.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type TableBody struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `gorm:"type:varchar(255);not null" json:"name"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (TableBody) TableName() string {
    return "table_bodies"
}
```

**CRUD 操作迁移示例**：

**Ent 操作 (原)**：
```go
// Ent 方式
client.TableBody.Create().
    SetName("My Table").
    Save(ctx)
```

**GORM 操作 (迁移后)**：
```go
// GORM 方式
tableBody := &TableBody{
    Name: "My Table",
}
db.Create(tableBody)
```

**迁移工具**：
- 可以编写脚本自动转换 Ent Schema 到 GORM Model
- 使用 GORM 的 AutoMigrate 功能管理数据库迁移

#### 2. **业务功能层集成** (Feature Layer)

**目标**：将 EtrxLite 的业务功能作为独立服务集成

**集成方式**：**微服务 + API Gateway**

```
coze-studio/
├── backend/
│   ├── api/
│   │   └── router/
│   │       └── extensions.go    # 扩展路由注册
│   │
│   └── extensions/
│       ├── etrxtable/            # 多维表格服务
│       │   ├── service.go
│       │   ├── handler.go
│       │   └── adapter.go        # 适配 Coze Studio 接口
│       │
│       ├── er-diagram/           # ER 图编辑器服务
│       │   ├── service.go
│       │   ├── handler.go
│       │   └── adapter.go
│       │
│       ├── code-generation/     # 代码生成服务
│       │   ├── service.go
│       │   ├── handler.go
│       │   └── adapter.go
│       │
│       └── mas/                  # MAS 多智能体系统
│           ├── service.go
│           ├── handler.go
│           └── adapter.go
```

**实现策略**：
- 每个功能模块独立实现，通过 HTTP/gRPC 接口提供服务
- 通过 Coze Studio 的路由系统注册扩展路由
- 使用适配器模式将 EtrxLite 的数据模型转换为 Coze Studio 的数据模型

**路由注册示例**：
```go
// backend/api/router/extensions.go
package router

import (
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/coze-dev/coze-studio/backend/extensions/etrxtable"
    "github.com/coze-dev/coze-studio/backend/extensions/er_diagram"
)

func RegisterExtensions(h *server.Hertz) {
    // 多维表格路由
    etrxtableGroup := h.Group("/api/v1/extensions/etrxtable")
    {
        etrxtableGroup.GET("/tables", etrxtable.ListTables)
        etrxtableGroup.POST("/tables", etrxtable.CreateTable)
        // ...
    }
    
    // ER 图编辑器路由
    erDiagramGroup := h.Group("/api/v1/extensions/er-diagram")
    {
        erDiagramGroup.GET("/diagrams", er_diagram.ListDiagrams)
        erDiagramGroup.POST("/diagrams", er_diagram.CreateDiagram)
        // ...
    }
}
```

#### 3. **前端集成** (Frontend Layer)

**目标**：在 Coze Studio 前端中集成 EtrxLite 的功能模块

**集成方式**：**独立包 + 动态加载**

```
coze-studio/
├── frontend/
│   ├── packages/
│   │   ├── studio/               # Coze Studio 原有包
│   │   │
│   │   └── extensions/          # 新增：扩展包
│   │       ├── etrxtable/       # 多维表格前端组件
│   │       │   ├── src/
│   │       │   │   ├── components/
│   │       │   │   ├── views/
│   │       │   │   └── index.tsx
│   │       │   └── package.json
│   │       │
│   │       ├── er-diagram/      # ER 图编辑器前端组件
│   │       │   ├── src/
│   │       │   └── package.json
│   │       │
│   │       └── code-generation/ # 代码生成前端组件
│   │           ├── src/
│   │           └── package.json
│   │
│   └── apps/
│       └── coze-studio/
│           └── src/
│               └── routes/
│                   └── extensions.tsx  # 扩展路由配置
```

**实现策略**：
- 将 EtrxLite 的 Vue 组件转换为 React 组件（或使用 Web Components）
- 通过 Rush monorepo 管理扩展包
- 使用动态路由和懒加载，按需加载扩展功能

**路由配置示例**：
```typescript
// frontend/apps/coze-studio/src/routes/extensions.tsx
import { lazy } from 'react';
import { Route, Routes } from 'react-router-dom';

// 懒加载扩展模块
const EtrxTable = lazy(() => import('@coze-extensions/etrxtable'));
const ERDiagram = lazy(() => import('@coze-extensions/er-diagram'));
const CodeGeneration = lazy(() => import('@coze-extensions/code-generation'));

export function ExtensionRoutes() {
  return (
    <Routes>
      <Route path="/extensions/etrxtable/*" element={<EtrxTable />} />
      <Route path="/extensions/er-diagram/*" element={<ERDiagram />} />
      <Route path="/extensions/code-generation/*" element={<CodeGeneration />} />
    </Routes>
  );
}
```

---

## 📦 具体集成方案

### 重要架构决策

#### 1. ORM 统一方案

**决策**：统一使用 **GORM** 作为 ORM 框架

**原因**：
1. **保持一致性**：与 Coze Studio 核心架构保持一致
2. **降低复杂度**：避免维护两套 ORM 框架
3. **简化集成**：无需适配器层，直接使用 GORM
4. **成熟稳定**：GORM 是 Go 生态中最成熟的 ORM 之一

**迁移策略**：
- 将 EtrxLite 的 Ent Schema 转换为 GORM Model
- 使用 GORM 的迁移功能管理数据库 Schema
- 保持数据模型结构不变，仅改变 ORM 实现方式

#### 2. Web 框架差异

**差异**：
- **Coze Studio**: 使用 **Hertz** (字节跳动开源的高性能 HTTP 框架)
- **EtrxLite**: 使用 **Gin** (流行的 Go Web 框架)

**决策**：**统一使用 Gin**

**原因**：
1. **生态成熟**：Gin 是 Go 生态中最流行的 Web 框架，生态更成熟
2. **简单易用**：Gin 的 API 更直观，学习成本更低
3. **兼容性好**：Gin 基于标准库，兼容性更好
4. **与 EtrxLite 一致**：EtrxLite 已使用 Gin，迁移成本更低

**迁移策略**：
- 将 Coze Studio 的 Hertz 路由和处理器迁移到 Gin
- 使用 Gin 的中间件系统
- 保持 API 接口不变，仅改变框架实现
- 创建适配器层，逐步迁移

#### 3. 数据库差异

**差异**：
- **Coze Studio**: 使用 **MySQL** (gorm.io/driver/mysql)
- **EtrxLite**: 使用 **PostgreSQL** (lib/pq, pgvector 支持)

**决策**：**统一使用 PostgreSQL 17**

**原因**：
1. **PostgreSQL 17 优势**：最新版本，性能优化和功能增强
2. **向量支持**：原生支持 pgvector，适合 RAG 等 AI 功能
3. **JSON 支持**：更好的 JSON 数据类型支持
4. **扩展性**：PostgreSQL 的扩展机制更灵活
5. **与 EtrxLite 一致**：EtrxLite 已使用 PostgreSQL，迁移成本更低

**迁移策略**：
- 将 Coze Studio 的 MySQL 数据模型迁移到 PostgreSQL
- 使用 GORM 的 PostgreSQL driver (gorm.io/driver/postgres)
- 数据库迁移脚本从 MySQL 语法转换为 PostgreSQL 语法
- 利用 PostgreSQL 17 的新特性（如更好的 JSON 支持、性能优化等）

#### 4. 前端框架差异

**差异**：
- **Coze Studio**: **React 18** + TypeScript + Rush Monorepo
- **EtrxLite**: **Vue 3** + TypeScript + Vite

**决策**：**统一使用 React**

**原因**：
1. **与 Coze Studio 一致**：保持前端技术栈统一
2. **Monorepo 优势**：Rush 管理的大型 monorepo 更适合企业级项目
3. **生态成熟**：React 生态更成熟，组件库更丰富

**迁移策略**：
- 将 EtrxLite 的 Vue 组件转换为 React 组件
- 使用 Coze Studio 的组件库和设计系统
- 集成到 Rush monorepo 中

#### 5. 配置管理差异

**差异**：
- **Coze Studio**: 使用 **环境变量** + godotenv (.env 文件)
- **EtrxLite**: 使用 **Viper** (YAML 配置文件)

**决策**：**统一使用 Viper (YAML 配置)**

**原因**：
1. **配置集中管理**：YAML 配置文件更直观，便于管理复杂配置
2. **环境变量支持**：Viper 支持环境变量覆盖，兼顾灵活性
3. **类型安全**：Viper 提供更好的配置类型检查和验证
4. **多格式支持**：Viper 支持 YAML、JSON、TOML 等多种格式
5. **与 EtrxLite 一致**：EtrxLite 已使用 Viper，迁移成本更低

**迁移策略**：
- 将 Coze Studio 的环境变量配置转换为 YAML 配置文件
- 使用 Viper 加载配置，支持环境变量覆盖
- 保持配置项不变，仅改变配置加载方式
- 创建统一的配置结构体，便于类型检查

#### 6. 日志系统差异

**差异**：
- **Coze Studio**: 自定义日志系统 (pkg/logs)
- **EtrxLite**: **Zap + Lumberjack** (结构化日志)

**决策**：**统一使用 EtrxLite 的日志系统（Zap + Lumberjack）**

**原因**：
1. **更优的实现**：Zap 是 Go 生态中性能最好的日志库
2. **结构化日志**：支持结构化日志，便于日志分析
3. **功能完整**：Lumberjack 提供日志轮转功能

**迁移策略**：
- 将 Coze Studio 的日志调用迁移到 Zap
- 创建适配器，保持 API 兼容
- 逐步替换，不影响现有功能

#### 7. 前端包管理差异

**差异**：
- **Coze Studio**: **Rush + PNPM** (monorepo)
- **EtrxLite**: npm/pnpm (单仓库)

**决策**：**统一使用 Rush + PNPM**

**原因**：
1. **Monorepo 优势**：更好的代码组织和依赖管理
2. **与 Coze Studio 一致**：保持前端架构统一
3. **企业级**：Rush 更适合大型项目

**迁移策略**：
- 将 EtrxLite 的前端代码集成到 Rush monorepo
- 创建独立的扩展包
- 使用 Rush 的依赖管理

### 方案 1: 统一基础设施集成 (优先级：高)

#### 1.1 统一日志系统

**目标**：引入 EtrxLite 的统一日志系统，提供更好的日志管理能力

**实施步骤**：
1. 创建日志适配器，将 EtrxLite 的日志接口适配到 Coze Studio
2. 通过配置开关控制是否启用统一日志
3. 保持 Coze Studio 原有日志系统作为 fallback

**文件结构**：
```
backend/core/logging/   # 统一日志系统
├── logger.go          # 日志实现
├── config.go          # 日志配置
└── middleware.go      # 日志中间件
```

**配置示例**：
```yaml
# backend/conf/config.yaml
extensions:
  etrxlite:
    enabled: true
    logging:
      enabled: true
      level: info
      output: file
      file_path: logs/coze-studio.log
```

#### 1.2 统一缓存系统

**目标**：引入 EtrxLite 的多级缓存系统（内存 + Redis）

**实施步骤**：
1. 创建缓存适配器
2. 集成到 Coze Studio 的缓存接口
3. 支持配置切换缓存策略

**文件结构**：
```
backend/core/cache/     # 统一缓存系统
├── manager.go         # 缓存管理器
├── strategy.go        # 缓存策略
└── adapter.go         # 缓存适配器（如需要）
```

### 方案 2: 多维表格系统集成 (优先级：高)

#### 2.1 后端集成

**目标**：将 EtrxLite 的多维表格系统作为独立服务集成

**实施步骤**：
1. **迁移数据模型**：将 Ent Schema 转换为 GORM Model
2. **迁移服务层**：将 Ent 的 CRUD 操作转换为 GORM 操作
3. **创建服务适配器**：适配 Coze Studio 的接口规范
4. **通过 API Gateway 暴露服务**

**文件结构**：
```
backend/extensions/etrxtable/
├── models/             # GORM 模型（从 Ent 迁移）
│   ├── table_body.go
│   ├── table_column.go
│   └── table_row.go
├── service.go          # 表格服务（使用 GORM）
├── handler.go          # HTTP 处理器
└── repository.go       # 数据访问层（GORM）
```

**API 设计**：
```go
// 保持 EtrxLite 的 API 风格，但通过 /api/v1/extensions/etrxtable 前缀
GET    /api/v1/extensions/etrxtable/tables
POST   /api/v1/extensions/etrxtable/tables
GET    /api/v1/extensions/etrxtable/tables/{id}/rows
POST   /api/v1/extensions/etrxtable/tables/{id}/rows
```

#### 2.2 前端集成

**目标**：在 Coze Studio 前端中添加多维表格功能

**实施步骤**：
1. 将 EtrxLite 的 Vue 组件转换为 React 组件
2. 创建独立的扩展包
3. 集成到 Coze Studio 的路由系统

**文件结构**：
```
frontend/packages/extensions/etrxtable/
├── src/
│   ├── components/
│   │   ├── TableView.tsx
│   │   ├── TableEditor.tsx
│   │   └── TableList.tsx
│   ├── hooks/
│   │   └── useTable.ts
│   ├── api/
│   │   └── table.ts
│   └── index.tsx
└── package.json
```

### 方案 3: ER 图编辑器集成 (优先级：中)

#### 3.1 后端集成

**目标**：集成 ER 图编辑器功能，支持可视化数据库设计

**实施步骤**：
1. 创建 ER 图服务适配器
2. 实现 SQL/DBML 导入导出功能
3. 集成数据库同步能力

**文件结构**：
```
backend/extensions/er-diagram/
├── service.go
├── handler.go
├── adapter.go
└── exporters/
    ├── sql.go
    └── dbml.go
```

#### 3.2 前端集成

**目标**：在 Coze Studio 中添加 ER 图编辑器界面

**实施步骤**：
1. 使用 React Flow 或类似库实现可视化编辑器
2. 参考 EtrxLite 的 ER 图编辑器实现
3. 集成到 Coze Studio 的界面中

### 方案 4: 代码生成系统集成 (优先级：中)

#### 4.1 后端集成

**目标**：集成代码生成功能，支持从表格定义生成代码

**实施步骤**：
1. 创建代码生成服务适配器
2. 支持 Go/Ent、Python/FastAPI、Java/Spring Boot 代码生成
3. 集成到 Coze Studio 的工作流中

**文件结构**：
```
backend/extensions/code-generation/
├── service.go
├── handler.go
├── generators/
│   ├── go.go
│   ├── python.go
│   └── java.go
└── templates/
    ├── go/
    ├── python/
    └── java/
```

### 方案 5: MAS 多智能体系统集成 (优先级：低)

#### 5.1 后端集成

**目标**：集成 MAS 系统，增强 Coze Studio 的 Agent 能力

**实施步骤**：
1. 创建 MAS 服务适配器
2. 与 Coze Studio 的 Agent 系统集成
3. 支持 Agent 间协作和任务调度

**注意事项**：
- Coze Studio 已有 Agent 系统，需要评估是否需要 MAS
- 可以作为可选增强功能

---

## 🔌 接口设计规范

### 1. API 路由规范

所有扩展功能的 API 都使用统一前缀：

```
/api/v1/extensions/{module-name}/{resource}
```

示例：
- `/api/v1/extensions/etrxtable/tables`
- `/api/v1/extensions/er-diagram/diagrams`
- `/api/v1/extensions/code-generation/generate`

### 2. 数据模型转换规范

**原则**：保持 Coze Studio 的数据模型不变，通过适配器转换

```go
// 适配器接口
type Adapter interface {
    // 将 Coze Studio 模型转换为 EtrxLite 模型
    ToEtrxLiteModel(cozeModel interface{}) (interface{}, error)
    
    // 将 EtrxLite 模型转换为 Coze Studio 模型
    ToCozeStudioModel(etrxliteModel interface{}) (interface{}, error)
}
```

### 3. 配置管理规范

所有扩展功能通过统一的配置管理：

```yaml
# backend/conf/config.yaml
extensions:
  etrxlite:
    enabled: true
    modules:
      etrxtable:
        enabled: true
        database:
          host: localhost
          port: 5432
      er_diagram:
        enabled: true
      code_generation:
        enabled: false
```

---

## 🚀 实施计划

### 阶段 0: Ent 到 GORM 迁移 (2-3 周) - 前置工作

**目标**：将 EtrxLite 的 Ent Schema 迁移到 GORM Model

**任务清单**：
- [ ] 分析 EtrxLite 中的所有 Ent Schema
- [ ] 编写 Ent Schema 到 GORM Model 的转换脚本
- [ ] 转换核心模型（TableBody, TableColumn, TableRow 等）
- [ ] 转换 ER 图模型
- [ ] 转换 MAS 模型
- [ ] 迁移 CRUD 操作代码
- [ ] 编写数据库迁移脚本
- [ ] 编写单元测试验证迁移正确性
- [ ] 更新文档

**验收标准**：
- 所有 Ent Schema 已转换为 GORM Model
- 所有 CRUD 操作已迁移到 GORM
- 数据库迁移脚本可以正常运行
- 所有测试通过

### 阶段 1: 基础设施集成 (2-3 周)

**目标**：完成统一基础设施的集成

**任务清单**：
- [ ] 创建扩展层目录结构
- [ ] 实现日志系统适配器
- [ ] 实现缓存系统适配器
- [ ] 实现认证系统适配器
- [ ] 集成 GORM 数据库连接（统一使用 Coze Studio 的 GORM 实例）
- [ ] 编写单元测试
- [ ] 更新文档

**验收标准**：
- 可以通过配置开关启用/禁用 EtrxLite 基础设施
- 统一使用 GORM 进行数据库操作
- 不影响 Coze Studio 原有功能
- 所有测试通过

### 阶段 2: 多维表格系统集成 (3-4 周)

**目标**：完成多维表格系统的集成

**任务清单**：
- [ ] 创建 EtrxTable 服务适配器
- [ ] 实现数据模型转换
- [ ] 实现 API 路由注册
- [ ] 前端组件转换（Vue → React）
- [ ] 集成到 Coze Studio 界面
- [ ] 编写单元测试和集成测试
- [ ] 更新文档

**验收标准**：
- 可以在 Coze Studio 中创建和管理多维表格
- 支持实时协作编辑
- 与 Coze Studio 的数据模型兼容

### 阶段 3: ER 图编辑器集成 (2-3 周)

**目标**：完成 ER 图编辑器的集成

**任务清单**：
- [ ] 创建 ER 图服务适配器
- [ ] 实现 SQL/DBML 导入导出
- [ ] 实现数据库同步功能
- [ ] 前端可视化编辑器实现
- [ ] 集成到 Coze Studio 界面
- [ ] 编写测试
- [ ] 更新文档

**验收标准**：
- 可以在 Coze Studio 中设计 ER 图
- 支持 SQL/DBML 导入导出
- 支持数据库同步

### 阶段 4: 代码生成系统集成 (2-3 周)

**目标**：完成代码生成系统的集成

**任务清单**：
- [ ] 创建代码生成服务适配器
- [ ] 实现多语言代码生成器
- [ ] 集成到工作流系统
- [ ] 前端界面实现
- [ ] 编写测试
- [ ] 更新文档

**验收标准**：
- 可以从表格定义生成代码
- 支持多种编程语言
- 生成代码质量良好

### 阶段 5: 优化和测试 (2 周)

**目标**：系统优化和全面测试

**任务清单**：
- [ ] 性能优化
- [ ] 安全性检查
- [ ] 全面集成测试
- [ ] 文档完善
- [ ] 用户手册编写

---

## 🔒 兼容性保证

### 1. 代码隔离

- 所有扩展代码放在 `extensions/` 目录下
- 不修改 Coze Studio 核心代码
- 通过接口和适配器实现集成

### 2. 配置隔离

- 扩展功能通过独立配置管理
- 可以通过配置开关完全禁用扩展功能
- 不影响 Coze Studio 原有配置

### 3. 依赖隔离

- 扩展功能使用独立的依赖管理
- 不强制 Coze Studio 升级依赖
- 通过适配器解决版本兼容问题

### 4. 数据隔离

- 扩展功能使用独立的数据表或命名空间
- 不修改 Coze Studio 原有数据模型
- 统一使用 GORM，但通过表名前缀或数据库隔离实现数据隔离
- 使用 GORM 的 TableName() 方法管理表名

---

## 📝 开发规范

### 1. 代码规范

- 遵循 Coze Studio 的代码规范
- 使用 Go 1.24+ 和 React 18+
- 编写清晰的注释和文档

### 2. 测试规范

- 每个功能模块都要有单元测试
- 集成测试覆盖主要功能
- 测试覆盖率不低于 70%

### 3. 文档规范

- 每个模块都要有 README
- API 文档使用 Swagger/OpenAPI
- 更新主项目文档

### 4. 版本管理

- 扩展功能独立版本管理
- 遵循语义化版本规范
- 记录变更日志

---

## 🎯 成功指标

### 技术指标

- [ ] 所有扩展功能可以独立启用/禁用
- [ ] 不影响 Coze Studio 核心功能
- [ ] 测试覆盖率 ≥ 70%
- [ ] API 响应时间 < 200ms (P95)
- [ ] 系统可用性 ≥ 99.9%

### 功能指标

- [ ] 多维表格系统完整可用
- [ ] ER 图编辑器完整可用
- [ ] 代码生成系统完整可用
- [ ] 统一基础设施完整可用

### 兼容性指标

- [ ] Coze Studio 可以独立升级
- [ ] 扩展功能可以独立升级
- [ ] 向后兼容性保持

---

## 📚 参考资料

### Coze Studio 文档
- [Coze Studio GitHub](https://github.com/coze-dev/coze-studio)
- [Coze Studio Wiki](https://github.com/coze-dev/coze-studio/wiki)
- [Coze 开发平台文档](https://www.coze.cn/open/docs)

### EtrxLite 文档
- [EtrxLite README](../etrxlite/README.md)
- [EtrxLite 后端文档](../etrxlite/backend/README.md)
- [EtrxLite 开发规范](../etrxlite/.cursor/rules/etrxlite.mdc)

### 技术文档
- [Hertz 框架文档](https://www.cloudwego.io/zh/docs/hertz/)
- [Ent ORM 文档](https://entgo.io/)
- [React 文档](https://react.dev/)
- [DDD 设计模式](https://martinfowler.com/bliki/DomainDrivenDesign.html)

---

## 🔄 更新日志

### 2025-01-XX
- 初始版本创建
- 完成项目对比分析
- 制定融合架构设计
- 制定实施计划

---

## 📞 联系方式

如有问题或建议，请联系：
- 项目维护者：[Your Name]
- 邮箱：[Your Email]
- GitHub Issues: [Project Issues]

---

**注意**：本文档会根据实际开发情况持续更新和完善。

