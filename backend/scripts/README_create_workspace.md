# 创建工作空间脚本使用说明

## 简介

`create_workspace.sh` 是一个用于创建 Coze Studio 工作空间（Space/Workspace）的命令行脚本。它通过 OpenAPI 端点 `/v1/workspaces` 来创建工作空间。

## 前置条件

1. **服务器运行**: 确保 Coze Studio 后端服务器正在运行
   ```bash
   cd backend
   ./start_server.sh start
   ```

2. **用户认证**: 创建工作空间需要用户认证（session cookie）

## 使用方法

### 1. 交互式模式（推荐）

运行脚本不带参数，会进入交互式模式：

```bash
cd backend/scripts
./create_workspace.sh
```

脚本会提示您输入：
- 工作空间名称（必需）
- 工作空间描述（可选）
- 图标文件 ID（可选）
- Coze Account ID（可选）
- 所有者 UID（可选）
- Session Cookie（可选）

### 2. 命令行参数模式

直接通过命令行参数指定：

```bash
./create_workspace.sh --name "我的工作空间" --description "这是一个测试工作空间"
```

### 3. 完整参数示例

```bash
./create_workspace.sh \
  --name "测试工作空间" \
  --description "用于测试的工作空间" \
  --icon-file-id "icon123" \
  --account-id "account456" \
  --owner-uid "user789" \
  --url "http://localhost:8888" \
  --cookie "session_key=your_session_key_here"
```

## 参数说明

| 参数 | 短参数 | 说明 | 必需 |
|------|--------|------|------|
| `--name` | `-n` | 工作空间名称 | 是 |
| `--description` | `-d` | 工作空间描述 | 否 |
| `--icon-file-id` | `-i` | 图标文件 ID | 否 |
| `--account-id` | `-a` | Coze Account ID | 否 |
| `--owner-uid` | `-o` | 所有者 UID | 否 |
| `--url` | `-u` | API 服务器地址（默认: http://localhost:8888） | 否 |
| `--cookie` | `-c` | Session Cookie（格式: session_key=xxx） | 否* |
| `--help` | `-h` | 显示帮助信息 | - |

* 如果没有提供 cookie，脚本会尝试从浏览器或环境变量中获取。

## 获取 Session Cookie

创建工作空间需要用户认证。您可以通过以下方式获取 session cookie：

### 方法 1: 从浏览器获取

1. 在浏览器中登录 Coze Studio
2. 打开开发者工具 (F12)
3. 进入 Application/Storage -> Cookies
4. 找到 `session_key` 的值
5. 使用 `--cookie "session_key=你的session_key值"` 参数

### 方法 2: 通过登录 API 获取

```bash
# 登录获取 session
curl -X POST http://localhost:8888/api/passport/web/email/login \
  -H "Content-Type: application/json" \
  -d '{"email":"your_email@example.com","password":"your_password"}' \
  -c cookies.txt

# 从 cookies.txt 中提取 session_key
SESSION_KEY=$(grep session_key cookies.txt | awk '{print $7}')
```

## 环境变量

您也可以通过环境变量设置默认值：

```bash
export BASE_URL="http://localhost:8888"
export SESSION_KEY="your_session_key"
```

## 示例输出

成功创建工作空间时：

```
✓ 工作空间创建成功！

响应数据:
{
    "code": 0,
    "data": {
        "id": "123456789",
        "name": "我的工作空间",
        "description": "这是一个测试工作空间",
        "icon_url": ""
    }
}

工作空间 ID: 123456789
```

## 错误处理

### 认证错误

如果收到 401 错误，说明需要提供有效的 session cookie：

```
✗ 工作空间创建失败 (HTTP 401)

错误响应:
{
    "code": 40100,
    "msg": "authentication required"
}

提示: 可能需要先登录获取 session cookie
```

### 服务器连接错误

如果无法连接到服务器：

```
错误: 无法连接到服务器 http://localhost:8888
请确保服务器正在运行: cd backend && ./start_server.sh start
```

## 故障排除

1. **服务器未运行**: 检查服务器是否在运行
   ```bash
   cd backend
   ./start_server.sh status
   ```

2. **Session 过期**: 重新登录获取新的 session cookie

3. **端口错误**: 确认服务器运行在正确的端口（默认 8888）

4. **权限问题**: 确保脚本有执行权限
   ```bash
   chmod +x create_workspace.sh
   ```

## API 参考

脚本调用的 API 端点：

- **URL**: `POST /v1/workspaces`
- **认证**: 需要 Session Cookie
- **请求体**:
  ```json
  {
    "name": "工作空间名称",
    "description": "工作空间描述（可选）",
    "icon_file_id": "图标文件 ID（可选）",
    "coze_account_id": "Coze Account ID（可选）",
    "owner_uid": "所有者 UID（可选）"
  }
  ```

## 相关文档

- [API 文档](../README.md)
- [服务器启动脚本](../start_server.sh)

