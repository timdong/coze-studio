# Coze Studio Super 快速开始指南

## 🚀 5 分钟快速启动

本指南帮助你在 5 分钟内启动 Coze Studio Super 项目。

---

## 📋 前置要求

### 必需
- ✅ Go 1.24.0+
- ✅ Docker 和 Docker Compose

### 可选
- Node.js 21+ (前端开发)
- PostgreSQL 客户端工具（psql）

---

## 🎯 快速启动步骤

### 步骤 1: 启动数据库和中间件（1 分钟）

```bash
cd /code/timdong/coze-studio

# 使用 Docker Compose 启动 PostgreSQL 17 + Redis + MinIO
docker-compose -f docker-compose.postgres.yml up -d

# 等待服务启动
sleep 10
```

**验证服务**：
```bash
# 检查容器状态
docker ps

# 应该看到 3 个容器：
# - coze-studio-postgres
# - coze-studio-redis
# - coze-studio-minio
```

### 步骤 2: 初始化数据库（1 分钟）

```bash
cd backend

# 运行数据库初始化脚本
./scripts/setup_database.sh

# 输出应该显示：
# ✅ PostgreSQL 连接成功
# ✅ 数据库创建成功
# ✅ 扩展安装完成
# ✅ 数据库设置完成！
```

### 步骤 3: 安装依赖（1 分钟）

```bash
# 在 backend 目录下
go mod tidy

# 等待依赖下载完成
```

### 步骤 4: 启动 Gin 服务器（2 分钟）

```bash
# 方式 1: 使用启动脚本（推荐）
./start_server.sh start

# 方式 2: 手动运行
go run cmd/gin_server/main.go

# 输出应该显示：
# [INFO] Config loaded successfully
# [INFO] Database connected successfully
# [INFO] Database migration completed
# [INFO] Starting HTTP server addr=0.0.0.0:8888
```

### 步骤 5: 验证服务（30 秒）

```bash
# 健康检查
curl http://localhost:8888/health

# 测试多维表格 API
curl http://localhost:8888/api/v1/extensions/etrxtable/tables

# 应该返回 JSON 响应
```

---

## ✅ 启动成功！

服务已经运行在：
- 🌐 **前端**: http://localhost:8888
- 🔌 **API**: http://localhost:8888/api/v1
- 📊 **健康检查**: http://localhost:8888/health

---

## 🔧 常用命令

### 服务器管理

```bash
# 启动服务器
./start_server.sh start

# 停止服务器
./start_server.sh stop

# 重启服务器
./start_server.sh restart

# 查看状态
./start_server.sh status

# 查看日志
./start_server.sh logs
```

### 数据库管理

```bash
# 连接数据库
PGPASSWORD=password psql -h localhost -p 5432 -U postgres -d coze_studio_db

# 查看所有表
PGPASSWORD=password psql -h localhost -p 5432 -U postgres -d coze_studio_db -c "\dt" | cat

# 查看表结构
PGPASSWORD=password psql -h localhost -p 5432 -U postgres -d coze_studio_db -c "\d table_bodies" | cat
```

### Docker 管理

```bash
# 查看容器状态
docker ps

# 查看 PostgreSQL 日志
docker logs coze-studio-postgres

# 停止所有服务
docker-compose -f docker-compose.postgres.yml down

# 停止并删除数据
docker-compose -f docker-compose.postgres.yml down -v
```

---

## 🐛 故障排查

### 问题 1: PostgreSQL 连接失败

**检查**：
```bash
docker ps | grep postgres
```

**解决**：
```bash
# 重启 PostgreSQL
docker restart coze-studio-postgres

# 或重新创建
docker-compose -f docker-compose.postgres.yml up -d postgres
```

### 问题 2: 端口被占用

**检查**：
```bash
lsof -i:5432  # PostgreSQL
lsof -i:8888  # Gin 服务器
lsof -i:6379  # Redis
```

**解决**：
- 修改 `docker-compose.postgres.yml` 中的端口映射
- 或停止占用端口的服务

### 问题 3: 数据库迁移失败

**检查日志**：
```bash
tail -f backend/logs/coze-studio.log | grep ERROR
```

**重新运行迁移**：
```bash
cd backend
./start_server.sh restart
```

### 问题 4: 编译失败

**清理并重新编译**：
```bash
cd backend
rm -rf bin/
go clean -cache
go mod tidy
go build -o bin/coze-studio-gin cmd/gin_server/main.go
```

---

## 📚 下一步

### 开发

1. **阅读开发文档**
   - [后端开发文档](backend/README.md)
   - [项目开发规范](.cursor/rules/coze-super.mdc)

2. **查看 API 文档**
   - http://localhost:8888/swagger/index.html (TODO)
   - http://localhost:8888/docs (TODO)

3. **开始开发**
   - 添加新的扩展功能
   - 实现业务逻辑
   - 编写测试用例

### 部署

1. **生产环境配置**
   - 修改 `config.yaml`
   - 设置环境变量
   - 配置安全选项

2. **编译生产版本**
   ```bash
   go build -ldflags="-s -w" -o bin/coze-studio-gin cmd/gin_server/main.go
   ```

3. **使用 Docker 部署**
   ```bash
   docker build -t coze-studio-super .
   docker run -d -p 8888:8888 coze-studio-super
   ```

---

## 💡 提示

- **日志位置**: `backend/logs/coze-studio.log`
- **配置文件**: `backend/config.yaml`
- **数据库**: PostgreSQL 17 with pgvector
- **默认端口**: 8888

---

## 📞 获取帮助

- 📖 查看 [完整文档](INTEGRATION_DEVELOPMENT_PLAN.md)
- 🐛 提交 [Issue](https://github.com/coze-dev/coze-studio/issues)
- 💬 加入社区讨论

---

**祝你使用愉快！** 🎉

