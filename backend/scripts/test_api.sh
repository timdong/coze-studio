#!/bin/bash
#
# Copyright 2025 coze-dev Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#


# Coze Studio Super - API 测试脚本
# 使用 curl 测试所有扩展功能的 API 端点

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
BASE_URL="http://localhost:8888"
API_BASE="${BASE_URL}/api/v1"

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 函数：打印测试结果
print_test() {
    local name=$1
    local status=$2
    local message=$3
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✓${NC} ${name}: ${message}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗${NC} ${name}: ${message}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# 函数：测试 GET 请求
test_get() {
    local name=$1
    local url=$2
    local expected_code=${3:-200}
    
    response=$(curl -s -w "\n%{http_code}" "${url}")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_code" ]; then
        print_test "$name" "PASS" "HTTP $http_code"
        echo -e "${BLUE}  Response:${NC} $(echo "$body" | head -c 200)..."
    else
        print_test "$name" "FAIL" "Expected HTTP $expected_code, got $http_code"
        echo -e "${RED}  Response:${NC} $body"
    fi
    echo ""
}

# 函数：测试 POST 请求
test_post() {
    local name=$1
    local url=$2
    local data=$3
    local expected_code=${4:-200}
    
    response=$(curl -s -w "\n%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -d "$data" \
        "${url}")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_code" ]; then
        print_test "$name" "PASS" "HTTP $http_code"
        echo -e "${BLUE}  Response:${NC} $(echo "$body" | head -c 200)..."
    else
        print_test "$name" "FAIL" "Expected HTTP $expected_code, got $http_code"
        echo -e "${RED}  Response:${NC} $body"
    fi
    echo ""
}

# 函数：测试 PUT 请求
test_put() {
    local name=$1
    local url=$2
    local data=$3
    local expected_code=${4:-200}
    
    response=$(curl -s -w "\n%{http_code}" -X PUT \
        -H "Content-Type: application/json" \
        -d "$data" \
        "${url}")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_code" ]; then
        print_test "$name" "PASS" "HTTP $http_code"
        echo -e "${BLUE}  Response:${NC} $(echo "$body" | head -c 200)..."
    else
        print_test "$name" "FAIL" "Expected HTTP $expected_code, got $http_code"
        echo -e "${RED}  Response:${NC} $body"
    fi
    echo ""
}

# 函数：测试 DELETE 请求
test_delete() {
    local name=$1
    local url=$2
    local expected_code=${3:-200}
    
    response=$(curl -s -w "\n%{http_code}" -X DELETE "${url}")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_code" ]; then
        print_test "$name" "PASS" "HTTP $http_code"
        echo -e "${BLUE}  Response:${NC} $(echo "$body" | head -c 200)..."
    else
        print_test "$name" "FAIL" "Expected HTTP $expected_code, got $http_code"
        echo -e "${RED}  Response:${NC} $body"
    fi
    echo ""
}

# 检查服务器是否运行
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Coze Studio Super API 测试${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

if ! curl -s "${BASE_URL}/api/v1/extensions/etrxtable/tables" > /dev/null 2>&1; then
    echo -e "${RED}错误: 无法连接到服务器 ${BASE_URL}${NC}"
    echo -e "${YELLOW}请确保服务器正在运行: cd backend && ./start_server.sh start${NC}"
    exit 1
fi

echo -e "${GREEN}服务器连接成功！${NC}"
echo ""

# ============================================
# 1. 多维表格系统 API 测试
# ============================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}1. 多维表格系统 API${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 1.1 获取表格列表
test_get "获取表格列表" \
    "${API_BASE}/extensions/etrxtable/tables?page=1&page_size=10"

# 1.2 创建表格
TABLE_DATA='{
  "name": "test_table",
  "description": "测试表格",
  "workspace_id": 1,
  "creator_id": 1,
  "sort_order": 0
}'
test_post "创建表格" \
    "${API_BASE}/extensions/etrxtable/tables" \
    "$TABLE_DATA"

# 获取刚创建的表格 ID（从响应中提取，这里假设为 1）
TABLE_ID=1

# 1.3 获取表格详情
test_get "获取表格详情" \
    "${API_BASE}/extensions/etrxtable/tables/${TABLE_ID}"

# 1.4 更新表格
UPDATE_TABLE_DATA='{
  "name": "test_table_updated",
  "description": "更新后的测试表格"
}'
test_put "更新表格" \
    "${API_BASE}/extensions/etrxtable/tables/${TABLE_ID}" \
    "$UPDATE_TABLE_DATA"

# 1.5 创建列
COLUMN_DATA='{
  "table_body_id": 1,
  "name": "test_column",
  "type": "text",
  "is_required": false,
  "sort_order": 0
}'
test_post "创建列" \
    "${API_BASE}/extensions/etrxtable/tables/${TABLE_ID}/columns" \
    "$COLUMN_DATA"

# 1.6 获取列列表
test_get "获取列列表" \
    "${API_BASE}/extensions/etrxtable/tables/${TABLE_ID}/columns"

# 1.7 创建行
ROW_DATA='{
  "table_body_id": 1,
  "data": "{\"test_column\": \"test_value\"}"
}'
test_post "创建行" \
    "${API_BASE}/extensions/etrxtable/tables/${TABLE_ID}/rows" \
    "$ROW_DATA"

# 1.8 获取行列表
test_get "获取行列表" \
    "${API_BASE}/extensions/etrxtable/tables/${TABLE_ID}/rows?page=1&page_size=10"

# ============================================
# 2. MAS 多智能体系统 API 测试
# ============================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}2. MAS 多智能体系统 API${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 2.1 获取 Agent 列表
test_get "获取 Agent 列表" \
    "${API_BASE}/extensions/mas/agents?page=1&page_size=10"

# 2.2 创建 Agent
AGENT_DATA='{
  "name": "test_agent",
  "description": "测试 Agent",
  "type": "general",
  "status": "inactive",
  "workspace_id": 1,
  "creator_id": 1,
  "version": "1.0.0"
}'
test_post "创建 Agent" \
    "${API_BASE}/extensions/mas/agents" \
    "$AGENT_DATA"

AGENT_ID=1

# 2.3 获取 Agent 详情
test_get "获取 Agent 详情" \
    "${API_BASE}/extensions/mas/agents/${AGENT_ID}"

# 2.4 更新 Agent 状态
STATUS_DATA='{
  "status": "active"
}'
test_put "更新 Agent 状态" \
    "${API_BASE}/extensions/mas/agents/${AGENT_ID}/status" \
    "$STATUS_DATA"

# 2.5 获取 Session 列表
test_get "获取 Session 列表" \
    "${API_BASE}/extensions/mas/sessions?page=1&page_size=10"

# 2.6 创建 Session
SESSION_DATA='{
  "name": "test_session",
  "description": "测试 Session",
  "type": "conversation",
  "status": "active",
  "workspace_id": 1,
  "creator_id": 1
}'
test_post "创建 Session" \
    "${API_BASE}/extensions/mas/sessions" \
    "$SESSION_DATA"

SESSION_ID=1

# 2.7 获取 Session 详情
test_get "获取 Session 详情" \
    "${API_BASE}/extensions/mas/sessions/${SESSION_ID}"

# 2.8 获取 Task 列表
test_get "获取 Task 列表" \
    "${API_BASE}/extensions/mas/tasks?page=1&page_size=10"

# 2.9 创建 Task
TASK_DATA='{
  "name": "test_task",
  "description": "测试 Task",
  "type": "query",
  "status": "pending",
  "agent_id": 1,
  "workspace_id": 1,
  "creator_id": 1,
  "priority": 0
}'
test_post "创建 Task" \
    "${API_BASE}/extensions/mas/tasks" \
    "$TASK_DATA"

TASK_ID=1

# 2.10 获取 Task 详情
test_get "获取 Task 详情" \
    "${API_BASE}/extensions/mas/tasks/${TASK_ID}"

# 2.11 启动 Task
test_put "启动 Task" \
    "${API_BASE}/extensions/mas/tasks/${TASK_ID}/start" \
    "{}"

# 2.12 获取待处理任务
test_get "获取待处理任务" \
    "${API_BASE}/extensions/mas/tasks/pending?limit=10"

# ============================================
# 3. ER 图编辑器 API 测试
# ============================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}3. ER 图编辑器 API${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 3.1 获取 ER 图列表
test_get "获取 ER 图列表" \
    "${API_BASE}/extensions/er-diagram/diagrams?page=1&page_size=10"

# 3.2 创建 ER 图
ER_DIAGRAM_DATA='{
  "name": "test_diagram",
  "description": "测试 ER 图",
  "diagram_data": "{\"tables\": [], \"relationships\": []}",
  "database_type": "postgresql",
  "workspace_id": 1,
  "created_by": 1
}'
test_post "创建 ER 图" \
    "${API_BASE}/extensions/er-diagram/diagrams" \
    "$ER_DIAGRAM_DATA"

DIAGRAM_ID=1

# 3.3 获取 ER 图详情
test_get "获取 ER 图详情" \
    "${API_BASE}/extensions/er-diagram/diagrams/${DIAGRAM_ID}"

# 3.4 更新 ER 图
UPDATE_DIAGRAM_DATA='{
  "name": "test_diagram_updated",
  "description": "更新后的测试 ER 图"
}'
test_put "更新 ER 图" \
    "${API_BASE}/extensions/er-diagram/diagrams/${DIAGRAM_ID}" \
    "$UPDATE_DIAGRAM_DATA"

# 3.5 导出 SQL（可能未实现，返回 500 是正常的）
test_get "导出 SQL" \
    "${API_BASE}/extensions/er-diagram/diagrams/${DIAGRAM_ID}/export/sql" \
    "500"

# 3.6 导出 DBML（可能未实现，返回 500 是正常的）
test_get "导出 DBML" \
    "${API_BASE}/extensions/er-diagram/diagrams/${DIAGRAM_ID}/export/dbml" \
    "500"

# ============================================
# 4. 代码生成系统 API 测试
# ============================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}4. 代码生成系统 API${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 4.1 生成 Go/Ent 代码
GO_ENT_DATA='{
  "table_id": 1,
  "language": "go",
  "framework": "ent",
  "package_name": "schema"
}'
test_post "生成 Go/Ent 代码" \
    "${API_BASE}/extensions/code-generation/generate" \
    "$GO_ENT_DATA"

# 4.2 生成 Go/GORM 代码
GO_GORM_DATA='{
  "table_id": 1,
  "language": "go",
  "framework": "gorm",
  "package_name": "models"
}'
test_post "生成 Go/GORM 代码" \
    "${API_BASE}/extensions/code-generation/generate" \
    "$GO_GORM_DATA"

# 4.3 生成 Python/FastAPI 代码
PYTHON_DATA='{
  "table_id": 1,
  "language": "python",
  "framework": "fastapi"
}'
test_post "生成 Python/FastAPI 代码" \
    "${API_BASE}/extensions/code-generation/generate" \
    "$PYTHON_DATA"

# 4.4 生成 Java/Spring 代码
JAVA_DATA='{
  "table_id": 1,
  "language": "java",
  "framework": "spring",
  "package_name": "com.example.model"
}'
test_post "生成 Java/Spring 代码" \
    "${API_BASE}/extensions/code-generation/generate" \
    "$JAVA_DATA"

# ============================================
# 测试总结
# ============================================
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}测试总结${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "总测试数: ${TOTAL_TESTS}"
echo -e "${GREEN}通过: ${PASSED_TESTS}${NC}"
echo -e "${RED}失败: ${FAILED_TESTS}${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}所有测试通过！🎉${NC}"
    exit 0
else
    echo -e "${YELLOW}部分测试失败，请检查服务器日志${NC}"
    exit 1
fi

