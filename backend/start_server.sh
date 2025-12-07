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


# Coze Studio Super - Gin 服务器启动脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 编译产物
BIN_DIR="./bin"
BIN_NAME="coze-studio-gin"
BIN_PATH="$BIN_DIR/$BIN_NAME"

# 日志目录
LOG_DIR="./logs"
mkdir -p "$LOG_DIR"

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

# 函数：检查 PostgreSQL
check_postgres() {
    print_info "Checking PostgreSQL connection..."
    
    # 从配置文件读取数据库配置（优先使用环境变量，然后使用配置文件默认值）
    PGHOST=${COZE_DATABASE_HOST:-localhost}
    PGPORT=${COZE_DATABASE_PORT:-15432}
    PGUSER=${COZE_DATABASE_USER:-postgres}
    PGPASSWORD=${COZE_DATABASE_PASSWORD:-password}
    PGDATABASE=${COZE_DATABASE_NAME:-coze_studio_db}
    
    if PGPASSWORD=$PGPASSWORD psql -h $PGHOST -p $PGPORT -U $PGUSER -d postgres -c "SELECT 1" > /dev/null 2>&1; then
        print_info "PostgreSQL is running"
    else
        print_error "PostgreSQL is not running or connection failed"
        print_info "Please start PostgreSQL first:"
        print_info "  docker run -d -p $PGPORT:5432 -e POSTGRES_PASSWORD=$PGPASSWORD postgres:17"
        exit 1
    fi
    
    # 检查数据库是否存在
    if PGPASSWORD=$PGPASSWORD psql -h $PGHOST -p $PGPORT -U $PGUSER -d postgres -c "SELECT 1 FROM pg_database WHERE datname='$PGDATABASE'" 2>/dev/null | grep -q 1; then
        print_info "Database '$PGDATABASE' exists"
    else
        print_warning "Database '$PGDATABASE' does not exist, creating..."
        PGPASSWORD=$PGPASSWORD psql -h $PGHOST -p $PGPORT -U $PGUSER -d postgres -c "CREATE DATABASE $PGDATABASE;" 2>/dev/null
        print_info "Database '$PGDATABASE' created"
    fi
}

# 函数：编译
compile() {
    print_info "Compiling Gin server..."
    
    mkdir -p "$BIN_DIR"
    
    go build -a -trimpath -o "$BIN_PATH" ./cmd/gin_server/main.go
    
    if [ $? -eq 0 ]; then
        print_info "Compilation successful: $BIN_PATH"
    else
        print_error "Compilation failed"
        exit 1
    fi
}

# 函数：启动
start() {
    print_info "Starting Coze Studio Super (Gin version)..."
    
    # 检查是否已在运行
    if pgrep -f "$BIN_NAME" > /dev/null; then
        print_warning "Server is already running"
        print_info "Use './start_server.sh restart' to restart"
        exit 1
    fi
    
    # 启动服务器
    nohup "$BIN_PATH" > "$LOG_DIR/gin_server.log" 2>&1 &
    
    sleep 2
    
    if pgrep -f "$BIN_NAME" > /dev/null; then
        print_info "Server started successfully"
        print_info "  PID: $(pgrep -f "$BIN_NAME")"
        print_info "  Log: $LOG_DIR/gin_server.log"
        print_info "  URL: http://localhost:8888"
    else
        print_error "Server failed to start"
        print_info "Check log: tail -f $LOG_DIR/gin_server.log"
        exit 1
    fi
}

# 函数：停止
stop() {
    print_info "Stopping Gin server..."
    
    if pgrep -f "$BIN_NAME" > /dev/null; then
        pkill -f "$BIN_NAME"
        sleep 1
        print_info "Server stopped"
    else
        print_warning "Server is not running"
    fi
}

# 函数：重启
restart() {
    print_info "Restarting Gin server..."
    stop
    compile
    check_postgres
    start
}

# 函数：查看状态
status() {
    if pgrep -f "$BIN_NAME" > /dev/null; then
        print_info "Server is running"
        print_info "  PID: $(pgrep -f "$BIN_NAME")"
        print_info "  URL: http://localhost:8888"
    else
        print_warning "Server is not running"
    fi
}

# 函数：查看日志
logs() {
    if [ -f "$LOG_DIR/gin_server.log" ]; then
        tail -f "$LOG_DIR/gin_server.log"
    else
        print_error "Log file not found: $LOG_DIR/gin_server.log"
    fi
}

# 主逻辑
case "${1:-start}" in
    start)
        compile
        check_postgres
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    logs)
        logs
        ;;
    compile)
        compile
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|logs|compile}"
        echo ""
        echo "Commands:"
        echo "  start    - Compile and start the server"
        echo "  stop     - Stop the server"
        echo "  restart  - Restart the server"
        echo "  status   - Check server status"
        echo "  logs     - View server logs"
        echo "  compile  - Compile only"
        exit 1
        ;;
esac

