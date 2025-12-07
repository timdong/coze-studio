# Coze Studio Super

**Coze Studio + EtrxLite 融合项目**

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.24.0+-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-336791?style=flat&logo=postgresql)](https://postgresql.org)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)](https://react.dev)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

**现代化 AI Agent 开发平台 + 企业级数据管理系统**

[English](README.md) | 中文

</div>

---

## 🎯 项目概述

**Coze Studio Super** 是 Coze Studio 和 EtrxLite 的融合项目，在 Coze Studio 基础上集成 EtrxLite 的优秀特性。

### 核心优势

| Coze Studio | EtrxLite |
|-------------|----------|
| ✅ AI Agent 开发平台 | ✅ 统一基础设施架构 |
| ✅ 工作流引擎 | ✅ 多维表格系统 |
| ✅ 知识库系统 | ✅ ER 图编辑器 |
| ✅ 插件系统 | ✅ 代码生成系统 |
| ✅ 活跃社区 | ✅ MAS 多智能体系统 |

### 融合特性

- 🔄 **统一技术栈**: Gin + GORM + PostgreSQL 17 + Viper + Zap
- 🔌 **松耦合集成**: 核心功能和扩展功能分离，可独立升级
- 📦 **模块化设计**: 清晰的分层架构，易于扩展
- 🚀 **企业级特性**: 完整的基础设施支持

---

## 🛠️ 技术栈

### 后端技术
- **语言**: Go 1.24.0+
- **Web 框架**: Gin v1.10.1
- **ORM**: GORM v1.25.11
- **数据库**: PostgreSQL 17 (with pgvector)
- **缓存**: Redis 6.0+
- **配置**: Viper (YAML)
- **日志**: Zap + Lumberjack
- **认证**: JWT

### 前端技术
- **框架**: React 18 + TypeScript
- **构建**: Rsbuild
- **包管理**: Rush + PNPM
- **状态管理**: Zustand
- **UI**: @coze-arch/coze-design

---

## 🚀 快速开始

### 环境要求

- Go 1.24.0+
- Node.js 21+
- PostgreSQL 17+
- Redis 6.0+

### 安装依赖

```bash
# 后端依赖
cd backend
go mod tidy

# 前端依赖
cd frontend
rush install
```

### 配置环境

1. 复制配置文件
```bash
cd backend
cp config.yaml config.local.yaml
```

2. 修改数据库配置（config.local.yaml）
```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  name: coze_studio_db
```

### 启动服务

#### 方式 1: 使用启动脚本（推荐）

```bash
cd backend
./start_server.sh start
```

#### 方式 2: 手动启动

```bash
cd backend
go run cmd/gin_server/main.go
```

### 访问应用

- **前端**: http://localhost:8888
- **API**: http://localhost:8888/api/v1
- **健康检查**: http://localhost:8888/health

---

## 📦 核心功能

### Coze Studio 核心功能

- 🤖 **AI Agent 开发**: 创建、配置、发布 AI Agent
- 🔄 **工作流引擎**: 可视化工作流编辑和执行
- 📚 **知识库系统**: 文档管理、向量检索、智能问答
- 🔌 **插件系统**: 第三方插件集成
- 💬 **对话管理**: 多轮对话、上下文管理

### EtrxLite 扩展功能 🆕

#### 1. 多维表格系统
- 📋 动态表结构创建
- 🔗 表间关联
- 👥 实时协作编辑
- 🎨 多视图支持（Grid/Kanban/Gantt）
- 🔐 细粒度权限控制

**API 端点**:
```
GET    /api/v1/extensions/etrxtable/tables
POST   /api/v1/extensions/etrxtable/tables
GET    /api/v1/extensions/etrxtable/tables/:id/rows
```

#### 2. ER 图编辑器
- 🎨 可视化数据库设计
- 📤 SQL/DBML 导入导出
- 🔄 数据库同步
- 📊 同步到表格

#### 3. 代码生成系统
- 🔧 表格到代码自动生成
- 🌐 多语言支持（Go/GORM, Python/FastAPI, Java/Spring）
- 📝 基于工作流的代码生成

#### 4. MAS 多智能体系统
- 🤖 Agent 管理和配置
- 📋 任务调度和执行
- 💬 会话管理和协调
- 🧠 记忆系统（短期/长期）
- 📊 性能指标监控

---

## 📁 项目结构

```
coze-studio/
├── backend/
│   ├── core/                   # EtrxLite 统一基础设施
│   │   ├── logging/           # Zap + Lumberjack
│   │   ├── cache/             # 内存 + Redis
│   │   ├── auth/              # JWT + BCrypt
│   │   ├── config/            # Viper
│   │   └── ...                # 其他 9 个模块
│   │
│   ├── extensions/            # 扩展功能
│   │   ├── etrxtable/        # 多维表格系统
│   │   ├── er-diagram/       # ER 图编辑器
│   │   ├── code-generation/  # 代码生成系统
│   │   └── mas/              # MAS 系统
│   │
│   ├── cmd/gin_server/        # Gin 服务器入口
│   ├── migrations/            # 数据库迁移
│   ├── config.yaml            # Viper 配置
│   └── start_server.sh        # 启动脚本
│
├── frontend/                   # 前端应用
│
└── docs/                       # 项目文档
    ├── INTEGRATION_DEVELOPMENT_PLAN.md
    ├── TECHNICAL_DIFFERENCES.md
    └── ...
```

---

## 📖 文档

### 核心文档
- [融合开发计划](INTEGRATION_DEVELOPMENT_PLAN.md) - 详细的开发计划
- [技术差异对比](TECHNICAL_DIFFERENCES.md) - 技术栈差异分析
- [Ent 迁移指南](ENT_TO_GORM_MIGRATION_GUIDE.md) - Ent 到 GORM 迁移
- [Hertz 迁移指南](HERTZ_TO_GIN_MIGRATION_GUIDE.md) - Hertz 到 Gin 迁移
- [后端开发文档](backend/README.md) - 后端开发指南
- [项目开发规范](.cursor/rules/coze-super.mdc) - 开发规范

### 进度文档
- [项目完成报告](PROJECT_COMPLETION_REPORT.md) - 完整的项目报告
- [今日成果总结](TODAY_ACHIEVEMENTS.md) - 每日成果
- [当前状态](CURRENT_STATUS.md) - 实时状态

---

## 🔧 开发指南

### 添加新的扩展功能

```bash
# 1. 创建扩展目录
mkdir -p backend/extensions/my-extension/{models,repository,service,handler}

# 2. 定义 GORM 模型
# backend/extensions/my-extension/models/my_model.go

# 3. 实现 Repository 层
# backend/extensions/my-extension/repository/repository.go

# 4. 实现 Service 层
# backend/extensions/my-extension/service/service.go

# 5. 实现 Handler 层并注册路由
# backend/extensions/my-extension/handler/handler.go
```

### 运行测试

```bash
cd backend
go test ./...
```

---

## 📊 项目状态

### 完成度

```
整体进度:                [██████████░░░░░░░░░░] 50%

阶段 0: 基础设施准备     [████████████████████] 100% ✅
阶段 0.1: Ent→GORM       [████████████████████] 100% ✅
阶段 0.2: Hertz→Gin      [████████████████████] 100% ✅
阶段 0.3: MySQL→PG       [░░░░░░░░░░░░░░░░░░░░] 0% ⏳
阶段 0.4: Env→Viper      [████████████████████] 100% ✅
```

### 下一步

1. **配置 PostgreSQL 17** - 运行数据库迁移
2. **实现扩展功能** - Repository/Service/Handler
3. **前端开发** - Vue → React 转换
4. **测试优化** - 单元测试和性能优化

---

## 🤝 贡献

欢迎贡献代码！请查看 [CONTRIBUTING.md](CONTRIBUTING.md)（待创建）。

---

## 📄 许可证

Apache License 2.0

---

## 📞 联系方式

- **项目地址**: [GitHub](https://github.com/coze-dev/coze-studio)
- **问题反馈**: [Issues](https://github.com/coze-dev/coze-studio/issues)

---

**项目状态**: 🚀 **阶段 0 完成，框架搭建完毕，进入实施阶段！**

