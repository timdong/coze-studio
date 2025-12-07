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


# Coze Studio Super - 登录并创建工作空间脚本
# 自动登录然后创建工作空间

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
LOGIN_URL="${BASE_URL}/api/passport/web/email/login"
WORKSPACE_URL="${BASE_URL}/api/v1/workspaces"
COOKIE_FILE="/tmp/coze_session_$$.txt"

# 默认登录信息（可以修改或通过环境变量设置）
EMAIL="${COZE_EMAIL:-admin@coze.studio}"
PASSWORD="${COZE_PASSWORD:-admin123}"

# 工作空间信息
WORKSPACE_NAME="${WORKSPACE_NAME:-我的新工作空间}"
WORKSPACE_DESC="${WORKSPACE_DESC:-通过脚本创建的工作空间}"

# 函数：打印使用说明
usage() {
    echo -e "${BLUE}使用方法:${NC}"
    echo "  $0 [选项]"
    echo ""
    echo -e "${BLUE}选项:${NC}"
    echo "  -n, --name NAME           工作空间名称（默认: 我的新工作空间）"
    echo "  -d, --description DESC    工作空间描述"
    echo "  -e, --email EMAIL         登录邮箱（默认: admin@coze.studio）"
    echo "  -p, --password PASSWORD   登录密码（默认: admin123）"
    echo "  -u, --url URL             API 服务器地址（默认: http://localhost:8888）"
    echo "  -h, --help                显示帮助信息"
    echo ""
    echo -e "${BLUE}环境变量:${NC}"
    echo "  COZE_EMAIL              登录邮箱"
    echo "  COZE_PASSWORD           登录密码"
    echo "  WORKSPACE_NAME          工作空间名称"
    echo "  WORKSPACE_DESC          工作空间描述"
    echo "  BASE_URL                API 服务器地址"
    echo ""
    echo -e "${BLUE}示例:${NC}"
    echo "  # 使用默认参数创建工作空间"
    echo "  $0"
    echo ""
    echo "  # 指定工作空间名称和描述"
    echo "  $0 --name \"测试工作空间\" --description \"用于测试\""
    echo ""
    echo "  # 使用自定义登录信息"
    echo "  $0 --email \"user@example.com\" --password \"mypassword\" --name \"我的工作空间\""
    exit 0
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            WORKSPACE_NAME="$2"
            shift 2
            ;;
        -d|--description)
            WORKSPACE_DESC="$2"
            shift 2
            ;;
        -e|--email)
            EMAIL="$2"
            shift 2
            ;;
        -p|--password)
            PASSWORD="$2"
            shift 2
            ;;
        -u|--url)
            BASE_URL="$2"
            LOGIN_URL="${BASE_URL}/api/passport/web/email/login"
            WORKSPACE_URL="${BASE_URL}/api/v1/workspaces"
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

# 清理函数
cleanup() {
    rm -f "$COOKIE_FILE"
}
trap cleanup EXIT

# 函数：检查服务器是否运行
check_server() {
    echo -e "${BLUE}检查服务器连接...${NC}"
    if ! curl -s -f "${BASE_URL}/api/v1/extensions/etrxtable/tables" > /dev/null 2>&1; then
        echo -e "${RED}错误: 无法连接到服务器 ${BASE_URL}${NC}"
        echo -e "${YELLOW}请确保服务器正在运行: cd backend && ./start_server.sh start${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ 服务器连接成功${NC}"
}

# 函数：登录
login() {
    echo -e "${BLUE}正在登录...${NC}"
    echo -e "${CYAN}邮箱: ${EMAIL}${NC}"

    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$LOGIN_URL" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\"}" \
        -c "$COOKIE_FILE")

    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')

    if [ "$HTTP_CODE" != "200" ]; then
        echo -e "${RED}✗ 登录失败 (HTTP ${HTTP_CODE})${NC}"
        echo -e "${RED}错误响应:${NC}"
        echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
        exit 1
    fi

    # 检查响应中的 code 字段
    CODE=$(echo "$BODY" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('code', ''))" 2>/dev/null || echo "")
    if [ "$CODE" != "0" ] && [ -n "$CODE" ]; then
        echo -e "${RED}✗ 登录失败${NC}"
        echo -e "${RED}错误响应:${NC}"
        echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
        exit 1
    fi

    echo -e "${GREEN}✓ 登录成功${NC}"

    # 提取 session key
    SESSION_KEY=$(grep session_key "$COOKIE_FILE" 2>/dev/null | awk '{print $7}')
    if [ -z "$SESSION_KEY" ]; then
        echo -e "${YELLOW}警告: 未能从 cookie 中提取 session key${NC}"
    fi
}

# 函数：创建工作空间
create_workspace() {
    echo ""
    echo -e "${BLUE}创建工作空间...${NC}"
    echo -e "${CYAN}名称: ${WORKSPACE_NAME}${NC}"
    if [ -n "$WORKSPACE_DESC" ]; then
        echo -e "${CYAN}描述: ${WORKSPACE_DESC}${NC}"
    fi

    # 构建请求 JSON
    JSON_DATA="{\"name\":\"${WORKSPACE_NAME}\""
    if [ -n "$WORKSPACE_DESC" ]; then
        JSON_DATA+=",\"description\":\"${WORKSPACE_DESC}\""
    fi
    JSON_DATA+="}"

    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$WORKSPACE_URL" \
        -H "Content-Type: application/json" \
        -b "$COOKIE_FILE" \
        -d "$JSON_DATA")

    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')

    if [ "$HTTP_CODE" != "200" ]; then
        echo -e "${RED}✗ 工作空间创建失败 (HTTP ${HTTP_CODE})${NC}"
        echo -e "${RED}错误响应:${NC}"
        echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
        exit 1
    fi

    # 检查响应
    CODE=$(echo "$BODY" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('code', ''))" 2>/dev/null || echo "")

    if [ "$CODE" != "0" ] && [ -n "$CODE" ]; then
        echo -e "${RED}✗ 工作空间创建失败 (code: ${CODE})${NC}"
        echo -e "${RED}错误响应:${NC}"
        echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
        exit 1
    fi

    echo -e "${GREEN}✓ 工作空间创建成功！${NC}"
    echo ""
    echo -e "${BLUE}工作空间信息:${NC}"
    echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"

    # 提取工作空间 ID
    WORKSPACE_ID=$(echo "$BODY" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('data', {}).get('id', ''))" 2>/dev/null || echo "")
    if [ -n "$WORKSPACE_ID" ]; then
        echo ""
        echo -e "${GREEN}工作空间 ID: ${WORKSPACE_ID}${NC}"
    fi
}

# 主流程
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  Coze Studio - 创建新工作空间${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# 检查服务器
check_server
echo ""

# 登录
login
echo ""

# 创建工作空间
create_workspace

echo ""
echo -e "${GREEN}完成！${NC}"

