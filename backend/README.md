# Coze Studio 后端服务

这是 Coze Studio 的后端服务，基于 Gin 框架构建，使用 PostgreSQL 数据库。

## 🚀 快速开始

### 前置要求

- Go 1.24+
- PostgreSQL 17+ (推荐使用 pgvector 扩展)
- Docker (可选)

### 启动服务

```bash
# 进入后端目录
cd backend

# 启动服务器（自动编译）
./start_server.sh start

# 查看日志
./start_server.sh logs

# 停止服务器
./start_server.sh stop

# 重启服务器
./start_server.sh restart

# 查看状态
./start_server.sh status
```

### 配置

服务器配置文件位于 `config.yaml`，主要配置项：

```yaml
# 数据库配置
database:
  host: localhost
  port: 15432
  user: postgres
  password: password
  name: coze_studio_db

# 服务器配置
server:
  host: 0.0.0.0
  port: 8888

# 日志配置
logging:
  level: info
  format: json
  output: file
  file_path: logs/api.log
```

支持环境变量覆盖，格式：`COZE_DATABASE_HOST=localhost`

## 📁 项目结构

```
backend/
├── api/                    # API 层
│   ├── handler/           # 请求处理器
│   ├── middleware/        # 中间件
│   └── router/           # 路由注册
├── application/           # 应用层 (DDD)
├── domain/               # 领域层 (DDD)
├── infra/                # 基础设施层
├── core/                 # 核心基础设施（日志/缓存/认证/配置）
├── extensions/           # 扩展层 (EtrxLite 功能)
│   ├── etrxtable/       # 多维表格系统
│   ├── er-diagram/      # ER 图编辑器
│   ├── code-generation/ # 代码生成系统
│   └── mas/             # MAS 多智能体系统
├── migrations/          # 数据库迁移
├── cmd/                 # 命令行工具
│   └── gin_server/     # Gin 服务器入口
├── config.yaml         # 配置文件
└── start_server.sh     # 启动脚本
```

## 🔧 开发

### 编译

```bash
# 编译到 bin/ 目录
go build -o ./bin/coze-studio ./cmd/gin_server/main.go
```

### 运行

```bash
# 开发模式
go run ./cmd/gin_server/main.go

# 或使用启动脚本
./start_server.sh start
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./extensions/etrxtable/...

# 带覆盖率
go test -cover ./...
```

## 🌐 API 文档

服务器启动后，访问：
- API 端点：`http://localhost:8888/api/v1`
- Swagger 文档：`http://localhost:8888/docs`（如果启用）

### 主要 API 端点

#### EtrxTable 多维表格系统

```bash
# 获取表格列表
GET /api/v1/extensions/etrxtable/tables

# 创建表格
POST /api/v1/extensions/etrxtable/tables

# 获取表格详情
GET /api/v1/extensions/etrxtable/tables/:id

# 更新表格
PUT /api/v1/extensions/etrxtable/tables/:id

# 删除表格
DELETE /api/v1/extensions/etrxtable/tables/:id
```

## 🗄️ 数据库

### PostgreSQL 安装

```bash
# 使用 Docker 启动 PostgreSQL 17 + pgvector
docker run -d \
  --name coze-postgres \
  -p 15432:5432 \
  -e POSTGRES_PASSWORD=password \
  pgvector/pgvector:pg17

# 创建数据库（首次使用）
PGPASSWORD=password psql -h localhost -p 15432 -U postgres -c "CREATE DATABASE coze_studio_db;"

# 安装 pgvector 扩展
PGPASSWORD=password psql -h localhost -p 15432 -U postgres -d coze_studio_db -c "CREATE EXTENSION IF NOT EXISTS vector;"
```

### 数据库迁移

数据库迁移会在服务器启动时自动运行（使用 GORM AutoMigrate）。

手动运行迁移：
```bash
go run ./cmd/gin_server/main.go
```

## 📝 日志

日志文件位于 `logs/` 目录：
- `logs/gin_server.log` - 服务器主日志
- `logs/api.log` - API 请求日志（如果配置）

查看实时日志：
```bash
./start_server.sh logs
# 或
tail -f logs/gin_server.log
```

## 🔐 认证

目前服务器使用基础认证，未来将集成 JWT。

## 🌍 环境变量

支持以下环境变量（会覆盖 config.yaml）：

```bash
# 数据库配置
export COZE_DATABASE_HOST=localhost
export COZE_DATABASE_PORT=15432
export COZE_DATABASE_USER=postgres
export COZE_DATABASE_PASSWORD=password
export COZE_DATABASE_NAME=coze_studio_db

# 服务器配置
export COZE_SERVER_HOST=0.0.0.0
export COZE_SERVER_PORT=8888

# 日志配置
export COZE_LOGGING_LEVEL=debug
```

## 🐛 故障排除

### 端口已被占用

```bash
# 查找占用端口的进程
lsof -i :8888

# 或使用脚本停止
./start_server.sh stop
```

### 数据库连接失败

1. 检查 PostgreSQL 是否运行：
   ```bash
   docker ps | grep postgres
   ```

2. 检查端口配置是否正确（15432）

3. 测试数据库连接：
   ```bash
   PGPASSWORD=password psql -h localhost -p 15432 -U postgres -c "SELECT 1"
   ```

### 编译失败

1. 确保 Go 版本 >= 1.24
2. 清理缓存：
   ```bash
   go clean -modcache
   go mod tidy
   ```

## 📚 参考文档

- [项目架构](../INTEGRATION_DEVELOPMENT_PLAN.md)
- [开发规范](../.cursor/rules/coze-super.mdc)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)

## 🤝 贡献

请阅读 [开发规范](../.cursor/rules/coze-super.mdc) 了解代码规范和最佳实践。

## 📄 许可证

Apache License 2.0
