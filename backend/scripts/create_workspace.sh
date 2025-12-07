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


# Coze Studio Super - 创建工作空间脚本
# 使用 OpenAPI 创建新的工作空间（Space/Workspace）

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 配置
BASE_URL="${BASE_URL:-http://localhost:8888}"
API_URL="${BASE_URL}/api/v1/workspaces"

# 默认值
DEFAULT_NAME="我的工作空间"
DEFAULT_DESCRIPTION=""

# 函数：打印使用说明
usage() {
    echo -e "${BLUE}使用方法:${NC}"
    echo "  $0 [选项]"
    echo ""
    echo -e "${BLUE}选项:${NC}"
    echo "  -n, --name NAME           工作空间名称（必需）"
    echo "  -d, --description DESC    工作空间描述"
    echo "  -i, --icon-file-id ID     图标文件 ID"
    echo "  -a, --account-id ID       Coze Account ID"
    echo "  -o, --owner-uid UID       所有者 UID"
    echo "  -u, --url URL             API 服务器地址（默认: http://localhost:8888）"
    echo "  -c, --cookie COOKIE       Session Cookie（格式: session_key=xxx）"
    echo "  -h, --help                显示帮助信息"
    echo ""
    echo -e "${BLUE}示例:${NC}"
    echo "  # 使用交互式模式创建工作空间"
    echo "  $0"
    echo ""
    echo "  # 直接指定名称创建工作空间"
    echo "  $0 --name \"测试工作空间\" --description \"这是一个测试工作空间\""
    echo ""
    echo "  # 指定服务器地址和 cookie"
    echo "  $0 --name \"我的工作空间\" --url http://localhost:8888 --cookie \"session_key=xxx\""
    exit 0
}

# 函数：检查服务器是否运行
check_server() {
    echo -e "${BLUE}检查服务器连接...${NC}"
    if ! curl -s -f "${BASE_URL}/api/v1/extensions/etrxtable/tables" > /dev/null 2>&1; then
        echo -e "${RED}错误: 无法连接到服务器 ${BASE_URL}${NC}"
        echo -e "${YELLOW}请确保服务器正在运行: cd backend && ./start_server.sh start${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ 服务器连接成功${NC}"
    echo ""
}

# 函数：交互式输入
interactive_input() {
    echo -e "${CYAN}=== 创建工作空间 ===${NC}"
    echo ""

    read -p "工作空间名称 (默认: ${DEFAULT_NAME}): " name
    name="${name:-$DEFAULT_NAME}"

    read -p "工作空间描述 (可选): " description
    description="${description:-$DEFAULT_DESCRIPTION}"

    read -p "图标文件 ID (可选): " icon_file_id
    icon_file_id="${icon_file_id:-}"

    read -p "Coze Account ID (可选): " coze_account_id
    coze_account_id="${coze_account_id:-}"

    read -p "所有者 UID (可选): " owner_uid
    owner_uid="${owner_uid:-}"

    read -p "Session Cookie (可选，格式: session_key=xxx): " cookie
    cookie="${cookie:-}"

    echo ""
}

# 解析命令行参数
NAME=""
DESCRIPTION=""
ICON_FILE_ID=""
COZE_ACCOUNT_ID=""
OWNER_UID=""
COOKIE=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            NAME="$2"
            shift 2
            ;;
        -d|--description)
            DESCRIPTION="$2"
            shift 2
            ;;
        -i|--icon-file-id)
            ICON_FILE_ID="$2"
            shift 2
            ;;
        -a|--account-id)
            COZE_ACCOUNT_ID="$2"
            shift 2
            ;;
        -o|--owner-uid)
            OWNER_UID="$2"
            shift 2
            ;;
        -u|--url)
            BASE_URL="$2"
            API_URL="${BASE_URL}/v1/workspaces"
            shift 2
            ;;
        -c|--cookie)
            COOKIE="$2"
            shift 2
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo -e "${RED}错误: 未知参数 $1${NC}"
            usage
            exit 1
            ;;
    esac
done

# 检查服务器
check_server

# 如果没有提供名称，使用交互式模式
if [ -z "$NAME" ]; then
    interactive_input
    NAME="${name:-$DEFAULT_NAME}"
    DESCRIPTION="${description:-$DEFAULT_DESCRIPTION}"
    ICON_FILE_ID="${icon_file_id:-}"
    COZE_ACCOUNT_ID="${coze_account_id:-}"
    OWNER_UID="${owner_uid:-}"
    COOKIE="${cookie:-}"
fi

# 构建请求 JSON
JSON_DATA="{"
JSON_DATA+="\"name\": \"${NAME}\""

if [ -n "$DESCRIPTION" ]; then
    JSON_DATA+=", \"description\": \"${DESCRIPTION}\""
fi

if [ -n "$ICON_FILE_ID" ]; then
    JSON_DATA+=", \"icon_file_id\": \"${ICON_FILE_ID}\""
fi

if [ -n "$COZE_ACCOUNT_ID" ]; then
    JSON_DATA+=", \"coze_account_id\": \"${COZE_ACCOUNT_ID}\""
fi

if [ -n "$OWNER_UID" ]; then
    JSON_DATA+=", \"owner_uid\": \"${OWNER_UID}\""
fi

JSON_DATA+="}"

echo -e "${BLUE}创建请求:${NC}"
echo -e "${CYAN}URL:${NC} ${API_URL}"
echo -e "${CYAN}数据:${NC} ${JSON_DATA}"
echo ""

# 构建 curl 命令
CURL_CMD="curl -s -w \"\\n%{http_code}\" -X POST"
CURL_CMD+=" -H \"Content-Type: application/json\""

# 如果有 cookie，添加到请求头
if [ -n "$COOKIE" ]; then
    CURL_CMD+=" -H \"Cookie: ${COOKIE}\""
fi

CURL_CMD+=" -d '${JSON_DATA}'"
CURL_CMD+=" \"${API_URL}\""

# 执行请求
echo -e "${YELLOW}发送创建请求...${NC}"
RESPONSE=$(eval $CURL_CMD)
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

echo ""

# 检查响应
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ 工作空间创建成功！${NC}"
    echo ""
    echo -e "${BLUE}响应数据:${NC}"
    echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"

    # 尝试提取工作空间 ID
    WORKSPACE_ID=$(echo "$BODY" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('data', {}).get('id', ''))" 2>/dev/null || echo "")
    if [ -n "$WORKSPACE_ID" ]; then
        echo ""
        echo -e "${GREEN}工作空间 ID: ${WORKSPACE_ID}${NC}"
    fi
else
    echo -e "${RED}✗ 工作空间创建失败 (HTTP ${HTTP_CODE})${NC}"
    echo ""
    echo -e "${RED}错误响应:${NC}"
    echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"

    # 检查是否是认证错误
    if [ "$HTTP_CODE" = "401" ] || echo "$BODY" | grep -q "authentication\|unauthorized\|session" 2>/dev/null; then
        echo ""
        echo -e "${YELLOW}提示: 可能需要先登录获取 session cookie${NC}"
        echo -e "${YELLOW}可以使用 --cookie 参数传递 session_key${NC}"
    fi

    exit 1
fi

