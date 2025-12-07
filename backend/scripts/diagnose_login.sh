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


# Coze Studio - 登录问题诊断脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

BASE_URL="${BASE_URL:-http://localhost:8888}"

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  登录问题诊断工具${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# 1. 检查服务器
echo -e "${BLUE}1. 检查服务器状态...${NC}"
if curl -s -f "${BASE_URL}/api/v1/extensions/etrxtable/tables" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 服务器运行正常${NC}"
else
    echo -e "${RED}✗ 服务器未运行或无法连接${NC}"
    exit 1
fi
echo ""

# 2. 检查数据库中的用户
echo -e "${BLUE}2. 检查数据库中的用户...${NC}"
USERS=$(PGPASSWORD=password psql -h localhost -p 15432 -U postgres -d coze_studio_db -t -c "SELECT email, username FROM users;" 2>&1 | cat)
if [ -n "$USERS" ]; then
    echo -e "${GREEN}✓ 找到以下用户:${NC}"
    echo "$USERS" | while read line; do
        if [ -n "$line" ]; then
            echo -e "  ${CYAN}-${NC} $line"
        fi
    done
else
    echo -e "${YELLOW}⚠ 未找到用户${NC}"
fi
echo ""

# 3. 测试登录
echo -e "${BLUE}3. 测试登录...${NC}"
echo -e "${CYAN}尝试使用默认账户登录:${NC}"
echo -e "  邮箱: ${YELLOW}admin@coze.studio${NC}"
echo -e "  密码: ${YELLOW}admin123${NC}"
echo ""

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/api/passport/web/email/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@coze.studio","password":"admin123"}')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    CODE=$(echo "$BODY" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('code', ''))" 2>/dev/null || echo "")
    if [ "$CODE" = "0" ]; then
        echo -e "${GREEN}✓ 登录成功！${NC}"
    else
        echo -e "${RED}✗ 登录失败 (code: ${CODE})${NC}"
        echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
    fi
else
    echo -e "${RED}✗ 登录失败 (HTTP ${HTTP_CODE})${NC}"
    echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
fi
echo ""

# 4. 提供建议
echo -e "${BLUE}4. 建议:${NC}"
echo -e "  - 如果登录失败，可以重置密码:"
echo -e "    ${CYAN}cd backend && go run scripts/reset_admin_password.go [新密码]${NC}"
echo -e ""
echo -e "  - 或者使用完整脚本创建工作空间:"
echo -e "    ${CYAN}cd backend/scripts && ./create_workspace_with_login.sh${NC}"

