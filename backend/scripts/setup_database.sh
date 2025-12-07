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


# Coze Studio Super - 数据库设置脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 数据库配置（从环境变量或使用默认值）
DB_HOST=${COZE_DATABASE_HOST:-localhost}
DB_PORT=${COZE_DATABASE_PORT:-5432}
DB_USER=${COZE_DATABASE_USER:-postgres}
DB_PASSWORD=${COZE_DATABASE_PASSWORD:-password}
DB_NAME=${COZE_DATABASE_NAME:-coze_studio_db}

# 函数：打印信息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 函数：检查 PostgreSQL 是否运行
check_postgres() {
    print_step "检查 PostgreSQL 连接..."
    
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "SELECT 1" > /dev/null 2>&1; then
        print_info "PostgreSQL 连接成功"
        return 0
    else
        print_error "无法连接到 PostgreSQL"
        print_info "请确保 PostgreSQL 正在运行："
        print_info "  方式 1（Docker Compose）："
        print_info "    cd /code/timdong/coze-studio"
        print_info "    docker-compose -f docker-compose.postgres.yml up -d"
        print_info ""
        print_info "  方式 2（Docker 手动）："
        print_info "    docker run -d -p $DB_PORT:5432 \\"
        print_info "      -e POSTGRES_PASSWORD=$DB_PASSWORD \\"
        print_info "      -e POSTGRES_DB=$DB_NAME \\"
        print_info "      --name coze-postgres \\"
        print_info "      pgvector/pgvector:pg17"
        return 1
    fi
}

# 函数：创建数据库
create_database() {
    print_step "检查数据库 '$DB_NAME'..."
    
    # 检查数据库是否存在
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
        print_info "数据库 '$DB_NAME' 已存在"
    else
        print_warning "数据库 '$DB_NAME' 不存在，正在创建..."
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME WITH ENCODING 'UTF8';"
        print_info "数据库 '$DB_NAME' 创建成功"
    fi
}

# 函数：安装扩展
install_extensions() {
    print_step "安装 PostgreSQL 扩展..."
    
    # 安装 pgvector 扩展
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "CREATE EXTENSION IF NOT EXISTS vector;" || true
    
    # 安装其他扩展
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";" || true
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "CREATE EXTENSION IF NOT EXISTS \"pg_trgm\";" || true
    
    print_info "扩展安装完成"
}

# 函数：运行初始化脚本
run_init_script() {
    print_step "运行初始化脚本..."
    
    SCRIPT_PATH="$(dirname "$0")/init_postgres.sql"
    
    if [ -f "$SCRIPT_PATH" ]; then
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$SCRIPT_PATH"
        print_info "初始化脚本执行成功"
    else
        print_warning "初始化脚本不存在: $SCRIPT_PATH"
    fi
}

# 函数：运行数据库迁移
run_migration() {
    print_step "运行数据库迁移（GORM AutoMigrate）..."
    
    print_info "迁移将在服务器启动时自动执行"
    print_info "或者运行: go run cmd/gin_server/main.go"
}

# 函数：验证数据库
verify_database() {
    print_step "验证数据库设置..."
    
    # 检查扩展
    echo -e "\n${BLUE}已安装的扩展:${NC}"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "\dx" | cat
    
    # 检查数据库大小
    echo -e "\n${BLUE}数据库信息:${NC}"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT pg_size_pretty(pg_database_size('$DB_NAME')) AS size;" | cat
    
    print_info "数据库验证完成"
}

# 主函数
main() {
    echo -e "${GREEN}"
    echo "╔════════════════════════════════════════════╗"
    echo "║   Coze Studio Super Database Setup        ║"
    echo "║   PostgreSQL 17 + pgvector                 ║"
    echo "╚════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    echo "数据库配置:"
    echo "  Host: $DB_HOST"
    echo "  Port: $DB_PORT"
    echo "  User: $DB_USER"
    echo "  Database: $DB_NAME"
    echo ""
    
    # 执行步骤
    if ! check_postgres; then
        exit 1
    fi
    
    create_database
    install_extensions
    run_init_script
    verify_database
    
    echo ""
    print_info "✅ 数据库设置完成！"
    echo ""
    print_info "下一步："
    print_info "  1. 启动服务器："
    print_info "     cd backend && ./start_server.sh start"
    print_info ""
    print_info "  2. 查看数据库表："
    print_info "     PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c \"\\dt\" | cat"
    echo ""
}

# 执行主函数
main

