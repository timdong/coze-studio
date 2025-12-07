# 技术栈差异对比文档

本文档详细对比 Coze Studio 和 EtrxLite 的技术栈差异，以及统一方案。

## 📊 技术栈对比表

| 技术维度 | Coze Studio | EtrxLite | 统一方案 | 优先级 |
|---------|------------|----------|---------|--------|
| **后端 Web 框架** | Hertz | Gin | **Gin** | 🔴 高 |
| **ORM** | GORM | Ent | **GORM** | 🔴 高 |
| **数据库** | MySQL | PostgreSQL | **PostgreSQL 17** | 🔴 高 |
| **前端框架** | React 18 | Vue 3 | **React 18** | 🔴 高 |
| **前端包管理** | Rush + PNPM | npm/pnpm | **Rush + PNPM** | 🟡 中 |
| **配置管理** | 环境变量 + godotenv | Viper (YAML) | **Viper (YAML)** | 🔴 高 |
| **日志系统** | 自定义日志 | Zap + Lumberjack | **Zap + Lumberjack** | 🟢 低 |
| **认证系统** | JWT (自定义) | JWT (统一实现) | **统一 JWT 实现** | 🟡 中 |
| **缓存系统** | Redis | Redis + 内存缓存 | **Redis + 内存缓存** | 🟢 低 |

---

## 🔴 高优先级差异（必须统一）

### 1. Web 框架：Hertz vs Gin

#### 差异说明

**Coze Studio (Hertz)**：
- 字节跳动开源的高性能 HTTP 框架
- 基于 Netpoll，性能优于标准库
- 支持 IDL 代码生成
- 完整的中间件系统

**EtrxLite (Gin)**：
- 流行的 Go Web 框架
- 基于标准库 net/http
- 简单易用，生态成熟

#### 统一方案：使用 Gin

**原因**：
1. **生态成熟**：Gin 是 Go 生态中最流行的 Web 框架
2. **简单易用**：API 更直观，学习成本更低
3. **兼容性好**：基于标准库，兼容性更好
4. **迁移成本低**：EtrxLite 已使用 Gin

**迁移步骤**：

1. **路由迁移**：
```go
// Hertz 路由 (原)
h := server.Default()
h.GET("/api/v1/tables", handler.ListTables)

// Gin 路由 (迁移后)
router := gin.New()
router.GET("/api/v1/tables", handler.ListTables)
```

2. **中间件迁移**：
```go
// Hertz 中间件 (原)
h.Use(middleware.AccessLogMW())
h.Use(middleware.RecoveryMW())

// Gin 中间件 (迁移后)
router.Use(gin.Logger())
router.Use(gin.Recovery())
```

3. **处理器迁移**：
```go
// Hertz 处理器 (原)
func ListTables(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, map[string]interface{}{"data": tables})
}

// Gin 处理器 (迁移后)
func ListTables(c *gin.Context) {
    c.JSON(200, gin.H{"data": tables})
}
```

**迁移工具**：可以编写脚本自动转换 Hertz 路由到 Gin

---

### 2. ORM：GORM vs Ent

详见 `ENT_TO_GORM_MIGRATION_GUIDE.md`

---

### 3. 前端框架：React vs Vue

#### 差异说明

**Coze Studio (React)**：
- React 18 + TypeScript
- Rush Monorepo 管理
- 企业级组件库

**EtrxLite (Vue)**：
- Vue 3 + TypeScript
- Vite 构建
- Element Plus UI 库

#### 统一方案：使用 React

**迁移步骤**：

1. **组件迁移**：
```vue
<!-- Vue 组件 (原) -->
<template>
  <div>
    <el-table :data="tables">
      <el-table-column prop="name" label="名称" />
    </el-table>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const tables = ref([])
</script>
```

```tsx
// React 组件 (迁移后)
import { useState } from 'react'
import { Table } from '@coze-arch/coze-design'

function TableList() {
  const [tables, setTables] = useState([])
  return (
    <Table dataSource={tables}>
      <Table.Column dataIndex="name" title="名称" />
    </Table>
  )
}
```

2. **状态管理迁移**：
```typescript
// Vue Pinia (原)
import { defineStore } from 'pinia'
export const useTableStore = defineStore('table', {
  state: () => ({ tables: [] }),
  actions: { fetchTables() {} }
})
```

```typescript
// React Zustand (迁移后)
import { create } from 'zustand'
export const useTableStore = create((set) => ({
  tables: [],
  fetchTables: async () => {}
}))
```

3. **路由迁移**：
```typescript
// Vue Router (原)
import { createRouter } from 'vue-router'
const router = createRouter({
  routes: [{ path: '/tables', component: TableList }]
})
```

```typescript
// React Router (迁移后)
import { createBrowserRouter } from 'react-router-dom'
const router = createBrowserRouter([
  { path: '/tables', element: <TableList /> }
])
```

**迁移工具**：
- 可以使用工具如 `vue-to-react` 辅助转换
- 手动转换确保代码质量

---

## 🟡 中优先级差异（建议统一）

### 4. 数据库：MySQL vs PostgreSQL

#### 差异说明

**Coze Studio (MySQL)**：
- 使用 MySQL 作为主数据库
- GORM MySQL driver

**EtrxLite (PostgreSQL)**：
- 使用 PostgreSQL 作为主数据库
- 需要 pgvector 扩展支持向量检索

#### 统一方案：PostgreSQL 17

**原因**：
1. **PostgreSQL 17 优势**：最新版本，性能优化和功能增强
2. **向量支持**：原生支持 pgvector，适合 RAG 等 AI 功能
3. **JSON 支持**：更好的 JSON 数据类型支持
4. **扩展性**：PostgreSQL 的扩展机制更灵活
5. **迁移成本低**：EtrxLite 已使用 PostgreSQL

**迁移步骤**：

1. **数据模型迁移**：
```go
// MySQL GORM Model (原)
type User struct {
    ID        uint   `gorm:"primaryKey;autoIncrement"`
    Name      string `gorm:"type:varchar(255)"`
    CreatedAt time.Time
}

// PostgreSQL GORM Model (迁移后)
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"type:varchar(255)"`
    CreatedAt time.Time
}
```

2. **数据库连接迁移**：
```go
// MySQL 连接 (原)
dsn := "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// PostgreSQL 连接 (迁移后)
dsn := "host=localhost user=postgres password=password dbname=dbname port=5432 sslmode=disable TimeZone=Asia/Shanghai"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

3. **SQL 语法迁移**：
```sql
-- MySQL (原)
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255)
);

-- PostgreSQL (迁移后)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);
```

4. **pgvector 扩展**：
```sql
-- 安装 pgvector 扩展
CREATE EXTENSION IF NOT EXISTS vector;

-- 创建向量列
ALTER TABLE documents ADD COLUMN embedding vector(1536);
```

---

### 5. 配置管理：环境变量 vs YAML

#### 差异说明

**Coze Studio (环境变量)**：
- 使用 `.env` 文件
- 通过 `godotenv` 加载
- 符合 12-Factor App

**EtrxLite (Viper)**：
- 使用 YAML 配置文件
- 支持环境变量覆盖
- 配置集中管理

#### 统一方案：Viper (YAML 配置)

**原因**：
1. **配置集中管理**：YAML 配置文件更直观，便于管理复杂配置
2. **环境变量支持**：Viper 支持环境变量覆盖，兼顾灵活性
3. **类型安全**：Viper 提供更好的配置类型检查和验证
4. **多格式支持**：Viper 支持 YAML、JSON、TOML 等多种格式
5. **迁移成本低**：EtrxLite 已使用 Viper

**迁移步骤**：

1. **环境变量转换为 YAML 配置**：
```bash
# .env (原)
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=password
```

```yaml
# config.yaml (迁移后)
database:
  host: localhost
  port: 5432
  user: postgres
  password: password
```

2. **配置加载代码**：
```go
// godotenv (原)
godotenv.Load()
host := os.Getenv("DATABASE_HOST")
```

```go
// Viper (迁移后)
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.ReadInConfig()

// 支持环境变量覆盖
viper.AutomaticEnv()
viper.SetEnvPrefix("APP")
viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

host := viper.GetString("database.host")
```

3. **配置结构体**：
```go
// 定义配置结构体
type Config struct {
    Database struct {
        Host     string `yaml:"host" mapstructure:"host"`
        Port     int    `yaml:"port" mapstructure:"port"`
        User     string `yaml:"user" mapstructure:"user"`
        Password string `yaml:"password" mapstructure:"password"`
    } `yaml:"database" mapstructure:"database"`
}

// 加载配置到结构体
var cfg Config
viper.Unmarshal(&cfg)
```

---

### 6. 前端包管理：Rush vs 单仓库

#### 差异说明

**Coze Studio (Rush Monorepo)**：
- 使用 Rush 管理 monorepo
- PNPM 作为包管理器
- 更好的依赖管理

**EtrxLite (单仓库)**：
- 单仓库结构
- npm/pnpm 管理依赖

#### 统一方案：Rush Monorepo

**迁移步骤**：

1. **创建扩展包**：
```json
// frontend/packages/extensions/etrxtable/package.json
{
  "name": "@coze-extensions/etrxtable",
  "version": "0.0.1",
  "dependencies": {
    "react": "workspace:*",
    "@coze-arch/coze-design": "workspace:*"
  }
}
```

2. **配置 Rush**：
```json
// rush.json
{
  "projects": [
    {
      "packageName": "@coze-extensions/etrxtable",
      "projectFolder": "packages/extensions/etrxtable"
    }
  ]
}
```

---

## 🟢 低优先级差异（可选统一）

### 7. 日志系统：自定义 vs Zap

#### 差异说明

**Coze Studio (自定义日志)**：
- 自定义日志接口
- 简单易用

**EtrxLite (Zap + Lumberjack)**：
- 高性能结构化日志
- 日志轮转支持

#### 统一方案：使用 Zap + Lumberjack

**迁移步骤**：

1. **创建日志适配器**：
```go
// 适配 Coze Studio 日志接口到 Zap
type LoggingAdapter struct {
    logger *zap.Logger
}

func (a *LoggingAdapter) Info(msg string, args ...interface{}) {
    a.logger.Info(msg, args...)
}
```

2. **逐步替换**：
- 保持 Coze Studio 日志接口不变
- 底层实现切换到 Zap
- 逐步迁移调用代码

---

### 8. 认证系统：JWT 实现差异

#### 差异说明

两者都使用 JWT，但实现细节可能不同。

#### 统一方案：统一 JWT 实现

**建议**：
- 使用 EtrxLite 的统一认证系统
- 创建适配器适配 Coze Studio 接口

---

## 📋 迁移优先级和时间表

### 阶段 1: 核心框架统一 (4-6 周)

- [ ] **Week 1-2**: Ent → GORM 迁移
- [ ] **Week 3**: Hertz → Gin 迁移
- [ ] **Week 4**: MySQL → PostgreSQL 17 迁移
- [ ] **Week 5**: 环境变量 → Viper 迁移
- [ ] **Week 6**: Vue → React 迁移

### 阶段 2: 基础设施统一 (2-3 周)

- [ ] **Week 7**: 配置管理统一
- [ ] **Week 8**: 日志系统统一
- [ ] **Week 9**: 认证系统统一

### 阶段 3: 数据库和包管理 (2-3 周)

- [ ] **Week 10**: 数据库配置统一
- [ ] **Week 11**: 前端包管理统一
- [ ] **Week 12**: 测试和优化

---

## 🛠️ 迁移工具和脚本

### 1. Gin → Hertz 转换工具

```bash
# 自动转换 Gin 路由到 Hertz
tools/gin_to_hertz_converter/main.go
```

### 2. Vue → React 转换工具

```bash
# 辅助转换 Vue 组件到 React
tools/vue_to_react_converter/
```

### 3. Ent → GORM 转换工具

详见 `ENT_TO_GORM_MIGRATION_GUIDE.md`

---

## ✅ 迁移检查清单

### Web 框架迁移

- [ ] 所有 Hertz 路由已转换为 Gin
- [ ] 所有中间件已迁移
- [ ] 所有处理器已迁移
- [ ] 测试通过

### 数据库迁移

- [ ] 所有 MySQL 数据模型已迁移到 PostgreSQL
- [ ] 数据库迁移脚本正常
- [ ] pgvector 扩展已安装
- [ ] 数据完整性验证通过
- [ ] 性能测试通过

### 配置管理迁移

- [ ] 所有环境变量已转换为 YAML 配置
- [ ] Viper 配置加载正常
- [ ] 环境变量覆盖功能正常
- [ ] 配置验证通过
- [ ] 测试通过

### ORM 迁移

- [ ] 所有 Ent Schema 已转换为 GORM Model
- [ ] 所有 CRUD 操作已迁移
- [ ] 数据库迁移脚本正常
- [ ] 测试通过

### 前端框架迁移

- [ ] 所有 Vue 组件已转换为 React
- [ ] 状态管理已迁移
- [ ] 路由已迁移
- [ ] UI 组件已替换
- [ ] 测试通过

### 配置管理迁移

- [ ] YAML 配置已转换为环境变量
- [ ] 配置加载代码已更新
- [ ] 环境变量文档已更新
- [ ] 测试通过

---

## 📚 参考资料

- [Hertz 官方文档](https://www.cloudwego.io/zh/docs/hertz/)
- [GORM 官方文档](https://gorm.io/docs/)
- [React 官方文档](https://react.dev/)
- [Rush 官方文档](https://rushjs.io/)

---

**注意**：本文档会根据实际迁移情况持续更新。

