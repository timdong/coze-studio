## 📊 API 测试快速参考

### 运行测试
```bash
# 改进版测试脚本（推荐）
./scripts/test_api_improved.sh

# 基础测试脚本
./scripts/test_api.sh
```

### 测试结果
- ✅ **17/18 测试通过** (94.4%)
- ❌ **1 个失败**: Agent 创建（JSON 格式问题）

### 主要 API 端点

#### 多维表格系统
- GET `/api/v1/extensions/etrxtable/tables` - 获取表格列表
- GET `/api/v1/extensions/etrxtable/tables/:id` - 获取表格详情

#### MAS 系统
- GET `/api/v1/extensions/mas/agents` - 获取 Agent 列表
- GET `/api/v1/extensions/mas/sessions` - 获取 Session 列表
- GET `/api/v1/extensions/mas/tasks` - 获取 Task 列表

#### ER 图编辑器
- GET `/api/v1/extensions/er-diagram/diagrams` - 获取 ER 图列表
- POST `/api/v1/extensions/er-diagram/diagrams` - 创建 ER 图

#### 代码生成系统
- POST `/api/v1/extensions/code-generation/generate` - 生成代码

详细测试报告请查看: `backend/scripts/API_TEST_REPORT.md`

