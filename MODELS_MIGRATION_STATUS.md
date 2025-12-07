# Ent 到 GORM 模型迁移状态

## 📊 迁移进度

**总计**: 55 个 Ent Schema
**已转换**: 15 个 (27%)
**待转换**: 40 个 (73%)

---

## ✅ 已转换模型（15 个）

### 多维表格系统（7 个）
- ✅ `table_body.go` - 表格主体
- ✅ `table_column.go` - 表格列
- ✅ `table_row.go` - 表格行
- ✅ `table_view.go` - 表格视图
- ✅ `table_sort.go` - 表格排序
- ✅ `table_link.go` - 表格关联
- ✅ `user.go` - 用户
- ✅ `workspace.go` - 工作空间
- ✅ `role.go` - 角色

### ER 图系统（1 个）
- ✅ `er_diagram.go` - ER 图

### MAS 系统（2 个）
- ✅ `agent.go` - Agent
- ✅ `mas_task.go` - MAS 任务
- ✅ `mas_session.go` - MAS 会话

### 工作流系统（1 个）
- ✅ `workflow.go` - 工作流（包含 WorkflowExecution）

### RAG 系统（3 个）
- ✅ `document.go` - 文档（包含 Chunk, Embedding）

---

## 🚧 待转换模型（40 个）

### Agent 相关（9 个）
- [ ] `agent_collaboration.go` - Agent 协作
- [ ] `agent_execution_log.go` - Agent 执行日志
- [ ] `agent_favorite.go` - Agent 收藏
- [ ] `agent_knowledge.go` - Agent 知识库
- [ ] `agent_memory.go` - Agent 记忆
- [ ] `agent_message.go` - Agent 消息
- [ ] `agent_metric.go` - Agent 指标
- [ ] `agent_tag.go` - Agent 标签
- [ ] `agent_tag_relation.go` - Agent 标签关系
- [ ] `agent_tool.go` - Agent 工具
- [ ] `agent_version.go` - Agent 版本

### 文档相关（3 个）
- [ ] `conversation.go` - 对话
- [ ] `document_category.go` - 文档分类
- [ ] `document_quality.go` - 文档质量
- [ ] `document_version.go` - 文档版本
- [ ] `qa_history.go` - 问答历史

### 工作流相关（4 个）
- [ ] `workflow_execution_log.go` - 工作流执行日志
- [ ] `workflow_permission.go` - 工作流权限
- [ ] `workflow_template.go` - 工作流模板
- [ ] `workflow_version.go` - 工作流版本

### 模板相关（3 个）
- [ ] `template_category.go` - 模板分类
- [ ] `template_rating.go` - 模板评分
- [ ] `template_version.go` - 模板版本
- [ ] `mas_template.go` - MAS 模板

### LLM 相关（5 个）
- [ ] `llm_model.go` - LLM 模型
- [ ] `llm_model_audit_log.go` - LLM 模型审计日志
- [ ] `llm_model_instance.go` - LLM 模型实例
- [ ] `llm_provider.go` - LLM 提供商
- [ ] `llm_provider_credential.go` - LLM 提供商凭证

### 其他（16 个）
- [ ] `company.go` - 公司
- [ ] `department.go` - 部门
- [ ] `database_connection.go` - 数据库连接
- [ ] `dynamic_table_metadata.go` - 动态表元数据
- [ ] `er_table_sync_history.go` - ER 表同步历史
- [ ] `menu.go` - 菜单
- [ ] `plugin.go` - 插件
- [ ] `view.go` - 自定义视图

---

## 📝 迁移策略

### 优先级分组

#### 第一优先级（核心功能）- 已完成 ✅
- 多维表格系统
- 用户和权限系统
- 工作空间管理

#### 第二优先级（扩展功能）- 部分完成 🚧
- ER 图系统 ✅
- MAS 系统 ✅
- 工作流系统 ✅
- RAG 系统 ✅

#### 第三优先级（辅助功能）- 待开始 ⏳
- Agent 相关功能
- 文档管理功能
- 模板系统
- LLM 配置

#### 第四优先级（可选功能）- 待开始 ⏳
- 公司部门管理
- 插件系统
- 其他辅助功能

---

## 🎯 下一步行动

### 立即开始

1. **转换 Agent 相关模型**（11 个）
   - 这些模型对 MAS 系统很重要
   - 预计时间：2-3 小时

2. **转换文档相关模型**（5 个）
   - 对 RAG 系统很重要
   - 预计时间：1-2 小时

3. **转换工作流相关模型**（4 个）
   - 对工作流引擎很重要
   - 预计时间：1-2 小时

---

## 📁 已创建的模型文件

### backend/extensions/etrxtable/models/
1. `table_body.go` - 表格主体 ✅
2. `table_column.go` - 表格列 ✅
3. `table_row.go` - 表格行 ✅
4. `table_view.go` - 表格视图 ✅
5. `table_sort.go` - 表格排序 ✅
6. `table_link.go` - 表格关联 ✅
7. `user.go` - 用户 ✅
8. `workspace.go` - 工作空间 ✅
9. `role.go` - 角色 ✅
10. `workflow.go` - 工作流 ✅
11. `document.go` - 文档（包含 Chunk, Embedding）✅

### backend/extensions/er-diagram/models/
1. `er_diagram.go` - ER 图 ✅

### backend/extensions/mas/models/
1. `agent.go` - Agent ✅
2. `mas_task.go` - MAS 任务 ✅
3. `mas_session.go` - MAS 会话 ✅

---

**注意**：本文档会随着迁移进度持续更新。

