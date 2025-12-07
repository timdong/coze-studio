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


# Coze Studio Super - 改进的 API 测试脚本
# 使用有效的 JSON 数据和正确的依赖关系

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

# 存储创建的 ID
CREATED_IDS=()

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
    else
        print_test "$name" "FAIL" "Expected HTTP $expected_code, got $http_code"
        echo -e "${RED}  Response:${NC} $(echo "$body" | head -c 200)"
    fi
    echo ""
}

# 函数：测试 POST 请求并提取 ID
test_post_with_id() {
    local name=$1
    local url=$2
    local data=$3
    local expected_code=${4:-200}
    local id_var=$5  # 变量名，用于存储提取的 ID
    
    response=$(curl -s -w "\n%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -d "$data" \
        "${url}")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_code" ]; then
        print_test "$name" "PASS" "HTTP $http_code"
        # 尝试提取 ID
        if [ -n "$id_var" ]; then
            id=$(echo "$body" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
            if [ -n "$id" ]; then
                eval "$id_var=$id"
                CREATED_IDS+=("$id")
            fi
        fi
    else
        print_test "$name" "FAIL" "Expected HTTP $expected_code, got $http_code"
        echo -e "${RED}  Response:${NC} $(echo "$body" | head -c 200)"
    fi
    echo ""
}

# 函数：测试 POST 请求
test_post() {
    test_post_with_id "$1" "$2" "$3" "$4" ""
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
    else
        print_test "$name" "FAIL" "Expected HTTP $expected_code, got $http_code"
        echo -e "${RED}  Response:${NC} $(echo "$body" | head -c 200)"
    fi
    echo ""
}

# 检查服务器是否运行
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Coze Studio Super API 测试（改进版）${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

if ! curl -s "${BASE_URL}/api/v1/extensions/etrxtable/tables" > /dev/null 2>&1; then
    echo -e "${RED}错误: 无法连接到服务器 ${BASE_URL}${NC}"
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

# 使用已存在的表格 ID（从数据库获取）
TABLE_ID=1

# 1.2 获取表格详情
test_get "获取表格详情" \
    "${API_BASE}/extensions/etrxtable/tables/${TABLE_ID}"

# 1.3 获取列列表
test_get "获取列列表" \
    "${API_BASE}/extensions/etrxtable/tables/${TABLE_ID}/columns"

# 1.4 获取行列表
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

# 2.2 创建 Agent（使用有效的 JSON）
AGENT_DATA='{
  "name": "test_agent",
  "description": "测试 Agent",
  "type": "general",
  "status": "inactive",
  "workspace_id": 1,
  "creator_id": 1,
  "version": "1.0.0",
  "config": "{}",
  "capabilities": "[]",
  "constraints": "[]",
  "tags": "[]"
}'
AGENT_ID=""
test_post_with_id "创建 Agent" \
    "${API_BASE}/extensions/mas/agents" \
    "$AGENT_DATA" \
    "200" \
    "AGENT_ID"

# 2.3 获取 Agent 详情（如果创建成功）
if [ -n "$AGENT_ID" ]; then
    test_get "获取 Agent 详情" \
        "${API_BASE}/extensions/mas/agents/${AGENT_ID}"
fi

# 2.4 获取 Session 列表
test_get "获取 Session 列表" \
    "${API_BASE}/extensions/mas/sessions?page=1&page_size=10"

# 2.5 创建 Session（使用有效的 JSON）
SESSION_DATA='{
  "name": "test_session",
  "description": "测试 Session",
  "type": "conversation",
  "status": "active",
  "workspace_id": 1,
  "creator_id": 1,
  "config": "{}",
  "context": "{}",
  "agent_ids": "[]"
}'
SESSION_ID=""
test_post_with_id "创建 Session" \
    "${API_BASE}/extensions/mas/sessions" \
    "$SESSION_DATA" \
    "200" \
    "SESSION_ID"

# 2.6 获取 Session 详情（如果创建成功）
if [ -n "$SESSION_ID" ]; then
    test_get "获取 Session 详情" \
        "${API_BASE}/extensions/mas/sessions/${SESSION_ID}"
fi

# 2.7 获取 Task 列表
test_get "获取 Task 列表" \
    "${API_BASE}/extensions/mas/tasks?page=1&page_size=10"

# 2.8 创建 Task（需要有效的 Agent ID）
if [ -n "$AGENT_ID" ]; then
    TASK_DATA="{
      \"name\": \"test_task\",
      \"description\": \"测试 Task\",
      \"type\": \"query\",
      \"status\": \"pending\",
      \"agent_id\": ${AGENT_ID},
      \"workspace_id\": 1,
      \"creator_id\": 1,
      \"priority\": 0,
      \"input\": \"{}\",
      \"output\": \"{}\",
      \"dependencies\": \"[]\"
    }"
    TASK_ID=""
    test_post_with_id "创建 Task" \
        "${API_BASE}/extensions/mas/tasks" \
        "$TASK_DATA" \
        "200" \
        "TASK_ID"
    
    # 2.9 获取 Task 详情（如果创建成功）
    if [ -n "$TASK_ID" ]; then
        test_get "获取 Task 详情" \
            "${API_BASE}/extensions/mas/tasks/${TASK_ID}"
    fi
fi

# 2.10 获取待处理任务
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

# 3.2 创建 ER 图（使用有效的 JSON）
ER_DIAGRAM_DATA='{
  "name": "test_diagram",
  "description": "测试 ER 图",
  "diagram_data": "{\"tables\": [], \"relationships\": []}",
  "database_type": "postgresql",
  "workspace_id": 1,
  "created_by": 1,
  "linked_table_ids": "[]",
  "sync_config": "{}"
}'
DIAGRAM_ID=""
test_post_with_id "创建 ER 图" \
    "${API_BASE}/extensions/er-diagram/diagrams" \
    "$ER_DIAGRAM_DATA" \
    "200" \
    "DIAGRAM_ID"

# 3.3 获取 ER 图详情（如果创建成功）
if [ -n "$DIAGRAM_ID" ]; then
    test_get "获取 ER 图详情" \
        "${API_BASE}/extensions/er-diagram/diagrams/${DIAGRAM_ID}"
fi

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

