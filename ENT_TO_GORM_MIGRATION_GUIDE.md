# Ent 到 GORM 迁移指南

本文档提供将 EtrxLite 的 Ent Schema 迁移到 GORM Model 的详细指南。

## 📋 迁移概述

### 为什么迁移？

1. **统一技术栈**：与 Coze Studio 保持一致，统一使用 GORM
2. **降低复杂度**：避免维护两套 ORM 框架
3. **简化集成**：无需适配器层，直接使用 GORM
4. **成熟稳定**：GORM 是 Go 生态中最成熟的 ORM 之一

### 迁移范围

需要迁移的 EtrxLite 核心模型：
- ✅ 多维表格系统（TableBody, TableColumn, TableRow, TableView 等）
- ✅ ER 图编辑器（ERDiagram, Table, Field, Relationship 等）
- ✅ MAS 多智能体系统（Agent, Task, Session, Memory 等）
- ✅ 代码生成系统（CodeGenerationTask, GeneratedFile 等）
- ✅ 用户和权限系统（User, Role, Workspace 等）

---

## 🔄 迁移步骤

### 步骤 1: 分析 Ent Schema

首先，分析 EtrxLite 中的 Ent Schema 定义：

```bash
# 查看所有 Ent Schema
find etrxlite/backend/platform/ent/schema -name "*.go" -type f
```

### 步骤 2: 转换 Schema 到 GORM Model

#### 2.1 基本字段映射

| Ent | GORM |
|-----|------|
| `field.Int("id")` | `ID uint \`gorm:"primaryKey"\`` |
| `field.String("name")` | `Name string \`gorm:"type:varchar(255)"\`` |
| `field.Time("created_at")` | `CreatedAt time.Time \`gorm:"autoCreateTime"\`` |
| `field.Time("updated_at")` | `UpdatedAt time.Time \`gorm:"autoUpdateTime"\`` |
| `field.Bool("deleted")` | `DeletedAt gorm.DeletedAt \`gorm:"index"\`` |

#### 2.2 关系映射

**Ent 关系**：
```go
// Ent Schema
func (TableBody) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("columns", TableColumn.Type),
        edge.To("rows", TableRow.Type),
    }
}
```

**GORM 关系**：
```go
// GORM Model
type TableBody struct {
    ID        uint           `gorm:"primaryKey"`
    Columns   []TableColumn  `gorm:"foreignKey:TableBodyID"`
    Rows      []TableRow     `gorm:"foreignKey:TableBodyID"`
}
```

### 步骤 3: 转换示例

#### 示例 1: TableBody 模型

**Ent Schema (原)**：
```go
// etrxlite/platform/ent/schema/table_body.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/edge"
)

type TableBody struct {
    ent.Schema
}

func (TableBody) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id"),
        field.String("name").MaxLen(255),
        field.String("description").Optional(),
        field.Int("workspace_id"),
        field.Time("created_at"),
        field.Time("updated_at"),
    }
}

func (TableBody) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("columns", TableColumn.Type),
        edge.To("rows", TableRow.Type),
        edge.To("views", TableView.Type),
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
    ID          uint           `gorm:"primaryKey" json:"id"`
    Name        string         `gorm:"type:varchar(255);not null" json:"name"`
    Description string         `gorm:"type:text" json:"description,omitempty"`
    WorkspaceID uint           `gorm:"not null;index" json:"workspace_id"`
    CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
    
    // 关系
    Columns     []TableColumn  `gorm:"foreignKey:TableBodyID" json:"columns,omitempty"`
    Rows        []TableRow     `gorm:"foreignKey:TableBodyID" json:"rows,omitempty"`
    Views       []TableView    `gorm:"foreignKey:TableBodyID" json:"views,omitempty"`
}

func (TableBody) TableName() string {
    return "table_bodies"
}
```

#### 示例 2: TableColumn 模型

**Ent Schema (原)**：
```go
type TableColumn struct {
    ent.Schema
}

func (TableColumn) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id"),
        field.String("name"),
        field.String("type").Default("text"),
        field.Bool("required").Default(false),
        field.Int("table_body_id"),
        field.Int("order").Default(0),
    }
}

func (TableColumn) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("table_body", TableBody.Type).
            Ref("columns").
            Field("table_body_id").
            Required().
            Unique(),
    }
}
```

**GORM Model (迁移后)**：
```go
type TableColumn struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Name        string         `gorm:"type:varchar(255);not null" json:"name"`
    Type        string         `gorm:"type:varchar(50);default:'text'" json:"type"`
    Required    bool           `gorm:"default:false" json:"required"`
    TableBodyID uint           `gorm:"not null;index" json:"table_body_id"`
    Order       int            `gorm:"default:0" json:"order"`
    CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
    
    // 关系
    TableBody   TableBody      `gorm:"foreignKey:TableBodyID" json:"table_body,omitempty"`
}

func (TableColumn) TableName() string {
    return "table_columns"
}
```

### 步骤 4: 迁移 CRUD 操作

#### 4.1 创建操作

**Ent 方式 (原)**：
```go
// Ent 创建
tableBody, err := client.TableBody.Create().
    SetName("My Table").
    SetWorkspaceID(workspaceID).
    Save(ctx)
```

**GORM 方式 (迁移后)**：
```go
// GORM 创建
tableBody := &TableBody{
    Name:        "My Table",
    WorkspaceID: workspaceID,
}
if err := db.Create(tableBody).Error; err != nil {
    return nil, err
}
```

#### 4.2 查询操作

**Ent 方式 (原)**：
```go
// Ent 查询
tableBody, err := client.TableBody.Query().
    Where(tablebody.ID(id)).
    WithColumns().
    WithRows().
    Only(ctx)
```

**GORM 方式 (迁移后)**：
```go
// GORM 查询
var tableBody TableBody
if err := db.Preload("Columns").Preload("Rows").
    First(&tableBody, id).Error; err != nil {
    return nil, err
}
```

#### 4.3 更新操作

**Ent 方式 (原)**：
```go
// Ent 更新
tableBody, err := client.TableBody.UpdateOneID(id).
    SetName("Updated Name").
    Save(ctx)
```

**GORM 方式 (迁移后)**：
```go
// GORM 更新
if err := db.Model(&TableBody{}).
    Where("id = ?", id).
    Update("name", "Updated Name").Error; err != nil {
    return err
}
```

#### 4.4 删除操作

**Ent 方式 (原)**：
```go
// Ent 删除（硬删除）
err := client.TableBody.DeleteOneID(id).Exec(ctx)

// Ent 软删除
err := client.TableBody.UpdateOneID(id).
    SetDeletedAt(time.Now()).
    Save(ctx)
```

**GORM 方式 (迁移后)**：
```go
// GORM 软删除（推荐）
if err := db.Delete(&TableBody{}, id).Error; err != nil {
    return err
}

// GORM 硬删除
if err := db.Unscoped().Delete(&TableBody{}, id).Error; err != nil {
    return err
}
```

### 步骤 5: 数据库迁移

#### 5.1 使用 GORM AutoMigrate

```go
// backend/extensions/etrxtable/migrations/migrate.go
package migrations

import (
    "gorm.io/gorm"
    "github.com/coze-dev/coze-studio/backend/extensions/etrxtable/models"
)

func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.TableBody{},
        &models.TableColumn{},
        &models.TableRow{},
        &models.TableView{},
        // ... 其他模型
    )
}
```

#### 5.2 在主程序中调用

```go
// backend/main.go 或 backend/extensions/init.go
import (
    "github.com/coze-dev/coze-studio/backend/extensions/etrxtable/migrations"
)

func initExtensions(db *gorm.DB) error {
    // 迁移 EtrxTable 模型
    if err := migrations.AutoMigrate(db); err != nil {
        return err
    }
    return nil
}
```

---

## 🛠️ 迁移工具脚本

### 自动转换脚本（示例）

```go
// tools/ent_to_gorm_converter/main.go
package main

import (
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "os"
)

// 这是一个简化的示例，实际转换需要更复杂的逻辑
func convertEntSchemaToGORM(entSchemaPath string) error {
    // 1. 解析 Ent Schema 文件
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, entSchemaPath, nil, parser.ParseComments)
    if err != nil {
        return err
    }
    
    // 2. 提取 Schema 定义
    // 3. 转换为 GORM Model
    // 4. 生成 GORM Model 文件
    
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: ent_to_gorm_converter <ent_schema_file>")
        os.Exit(1)
    }
    
    if err := convertEntSchemaToGORM(os.Args[1]); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
```

---

## ✅ 迁移检查清单

### 模型迁移检查

- [ ] 所有 Ent Schema 已转换为 GORM Model
- [ ] 字段类型映射正确
- [ ] 关系映射正确（外键、关联）
- [ ] 索引定义正确
- [ ] 表名定义正确（TableName() 方法）

### 操作迁移检查

- [ ] 创建操作已迁移
- [ ] 查询操作已迁移（包括预加载）
- [ ] 更新操作已迁移
- [ ] 删除操作已迁移（软删除/硬删除）
- [ ] 事务处理已迁移

### 数据库迁移检查

- [ ] AutoMigrate 可以正常运行
- [ ] 数据库表结构正确
- [ ] 索引创建正确
- [ ] 外键约束正确

### 测试检查

- [ ] 单元测试已更新
- [ ] 集成测试已更新
- [ ] 所有测试通过
- [ ] 性能测试通过

---

## 📚 参考资料

### GORM 文档
- [GORM 官方文档](https://gorm.io/docs/)
- [GORM 模型定义](https://gorm.io/docs/models.html)
- [GORM 关联](https://gorm.io/docs/belongs_to.html)
- [GORM 迁移](https://gorm.io/docs/migration.html)

### Ent 文档
- [Ent 官方文档](https://entgo.io/docs/getting-started)
- [Ent Schema 定义](https://entgo.io/docs/schema-def)

---

## 🔍 常见问题

### Q1: 如何处理 Ent 的复杂查询？

**A**: GORM 提供了强大的查询构建器，可以替代 Ent 的查询：

```go
// Ent 复杂查询
users, err := client.User.Query().
    Where(user.AgeGT(18)).
    Where(user.NameContains("John")).
    Order(user.ByAge()).
    Limit(10).
    All(ctx)

// GORM 等价查询
var users []User
db.Where("age > ?", 18).
    Where("name LIKE ?", "%John%").
    Order("age ASC").
    Limit(10).
    Find(&users)
```

### Q2: 如何处理 Ent 的预加载（With）？

**A**: GORM 使用 Preload 实现：

```go
// Ent 预加载
user, err := client.User.Query().
    WithPosts().
    WithComments().
    Only(ctx)

// GORM 预加载
var user User
db.Preload("Posts").Preload("Comments").
    First(&user, id)
```

### Q3: 如何处理 Ent 的事务？

**A**: GORM 的事务使用方式：

```go
// Ent 事务
err := client.WithTx(ctx, func(tx *ent.Tx) error {
    // ...
    return nil
})

// GORM 事务
err := db.Transaction(func(tx *gorm.DB) error {
    // ...
    return nil
})
```

---

## 📝 迁移进度跟踪

### 已完成迁移

- [ ] TableBody
- [ ] TableColumn
- [ ] TableRow
- [ ] TableView
- [ ] ERDiagram
- [ ] Agent
- [ ] Task
- [ ] Session
- [ ] Memory

### 待迁移

- [ ] 其他模型...

---

**注意**：本文档会根据实际迁移情况持续更新。

