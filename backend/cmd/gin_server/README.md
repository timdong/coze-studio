# Gin Server

这是 Coze Studio Super 的 Gin 版本服务器入口。

## 功能

- 使用 Gin 框架替代 Hertz
- 集成统一基础设施（Viper, Zap, GORM）
- 支持扩展功能路由
- 优雅关闭
- 配置热加载

## 运行

### 开发模式

```bash
cd backend
go run cmd/gin_server/main.go
```

### 编译运行

```bash
cd backend
go build -o bin/coze-studio-gin cmd/gin_server/main.go
./bin/coze-studio-gin
```

## 配置

使用 `backend/config.yaml` 配置文件。

支持环境变量覆盖：
```bash
export COZE_SERVER_PORT=8888
export COZE_DATABASE_HOST=localhost
```

## 端口

默认端口：`8888`

可以通过配置文件或环境变量修改：
```yaml
server:
  port: 8888
```

或

```bash
export COZE_SERVER_PORT=8888
```

## API 文档

启动后访问：
- http://localhost:8888/swagger/index.html (TODO)
- http://localhost:8888/docs (TODO)

## 健康检查

```bash
curl http://localhost:8888/health
```

## 与原服务器的区别

| 特性 | Hertz 版本 (main.go) | Gin 版本 (cmd/gin_server/main.go) |
|------|---------------------|-----------------------------------|
| Web 框架 | Hertz | Gin |
| 配置 | 环境变量 + godotenv | Viper (YAML) |
| 日志 | 自定义日志 | Zap + Lumberjack |
| 数据库 | MySQL | PostgreSQL 17 |
| 路由 | IDL 生成 | 手动注册 |
| 扩展功能 | 无 | 支持（etrxtable, er-diagram, mas 等）|

## 迁移计划

当前状态：**并行运行**
- Hertz 版本：`go run main.go`
- Gin 版本：`go run cmd/gin_server/main.go`

最终目标：**完全替代 Hertz 版本**

