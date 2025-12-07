# 数据库迁移

本目录包含数据库迁移脚本和工具。

## 迁移文件

- `001_initial_schema.go` - 初始化数据库表结构（GORM AutoMigrate）

## 使用方法

### 自动迁移（推荐）

使用 GORM AutoMigrate 自动创建和更新表结构：

```go
import (
    "github.com/coze-dev/coze-studio/backend/migrations"
)

func initDatabase(db *gorm.DB) error {
    return migrations.AutoMigrate(db)
}
```

### 手动执行

```bash
# 在 main.go 中调用
go run main.go migrate
```

## 迁移内容

### 已包含的表

#### 用户和权限（6 个）
- `users` - 用户
- `roles` - 角色
- `workspaces` - 工作空间
- `menus` - 菜单
- `companies` - 公司
- `departments` - 部门

#### 多维表格（8 个）
- `table_bodies` - 表格主体
- `table_columns` - 表格列
- `table_rows` - 表格行
- `table_views` - 表格视图
- `table_sorts` - 表格排序
- `table_links` - 表格关联
- `dynamic_table_metadata` - 动态表元数据
- `views` - 自定义视图

#### ER 图（2 个）
- `er_diagrams` - ER 图
- `er_table_sync_histories` - ER 表同步历史

#### MAS 系统（14 个）
- `agents` - Agent
- `mas_tasks` - MAS 任务
- `mas_sessions` - MAS 会话
- `agent_versions` - Agent 版本
- `agent_memories` - Agent 记忆
- `agent_metrics` - Agent 指标
- `agent_tags` - Agent 标签
- `agent_tools` - Agent 工具
- `agent_messages` - Agent 消息
- `agent_execution_logs` - Agent 执行日志
- `agent_favorites` - Agent 收藏
- `agent_collaborations` - Agent 协作
- `agent_knowledges` - Agent 知识库
- `mas_templates` - MAS 模板

#### 工作流（8 个）
- `workflows` - 工作流
- `workflow_executions` - 工作流执行
- `workflow_templates` - 工作流模板
- `workflow_permissions` - 工作流权限
- `workflow_versions` - 工作流版本
- `workflow_execution_logs` - 工作流执行日志
- `template_categories` - 模板分类
- `template_ratings` - 模板评分

#### RAG 系统（8 个）
- `documents` - 文档
- `chunks` - 文档分块
- `embeddings` - 向量嵌入
- `conversations` - 对话
- `qa_histories` - 问答历史
- `document_categories` - 文档分类
- `document_qualities` - 文档质量
- `document_versions` - 文档版本

#### LLM 系统（6 个）
- `llm_providers` - LLM 提供商
- `llm_provider_credentials` - LLM 凭证
- `llm_models` - LLM 模型
- `llm_model_audit_logs` - LLM 审计日志
- `llm_model_instances` - LLM 模型实例

#### 其他（2 个）
- `plugins` - 插件
- `database_connections` - 数据库连接

**总计**: **54 个表**

## PostgreSQL 特性

### pgvector 扩展

需要在 PostgreSQL 中安装 pgvector 扩展：

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

### JSON 支持

使用 PostgreSQL 的 JSONB 类型存储 JSON 数据，提供更好的性能和查询能力。

## 注意事项

1. **首次运行**: AutoMigrate 会自动创建所有表
2. **字段变更**: 修改模型后重新运行 AutoMigrate
3. **外键约束**: GORM 会自动创建外键约束
4. **索引**: GORM 会自动创建索引
5. **软删除**: 使用 `DeletedAt` 字段实现软删除

## 回滚

GORM AutoMigrate 不支持自动回滚，如需回滚：

```sql
-- 删除所有表（谨慎操作）
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;
```

## 手动迁移SQL

如果需要手动执行 SQL，可以使用：

```bash
# 导出当前表结构
pg_dump -h localhost -U postgres -d coze_studio_db --schema-only > schema.sql

# 导入表结构
psql -h localhost -U postgres -d coze_studio_db < schema.sql
```

