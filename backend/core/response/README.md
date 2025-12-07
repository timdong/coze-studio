# Response Framework

## 概述

响应处理框架提供了统一的API响应格式处理，使所有API响应保持一致的结构和格式。

## 功能特性

- ✅ 统一的响应格式
- ✅ 类型安全的错误处理
- ✅ 支持分页响应
- ✅ 与Gin框架完美集成
- ✅ 自动错误处理
- ✅ 预定义的错误类型

## 核心组件

### 1. 错误类型 (`errors.go`)

定义了统一的错误类型和错误代码：

- `Error` 结构体 - 响应错误
- `ErrorCode` 类型 - 错误代码枚举
- 预定义的错误常量
- 错误包装和转换工具

**错误代码**:
- `ErrorCodeBadRequest` - 请求参数错误
- `ErrorCodeUnauthorized` - 未授权
- `ErrorCodeForbidden` - 禁止访问
- `ErrorCodeNotFound` - 资源不存在
- `ErrorCodeValidationFailed` - 验证失败
- `ErrorCodeInternalError` - 内部服务器错误
- 等等...

### 2. 响应格式化 (`formatter.go`)

提供统一的响应格式化函数：

**成功响应**:
- `Success()` - 成功响应
- `Created()` - 创建成功响应
- `Updated()` - 更新成功响应
- `Deleted()` - 删除成功响应
- `SuccessWithPagination()` - 带分页的成功响应

**错误响应**:
- `BadRequest()` - 请求参数错误
- `Unauthorized()` - 未授权
- `Forbidden()` - 禁止访问
- `NotFound()` - 资源不存在
- `Conflict()` - 资源冲突
- `ValidationFailed()` - 验证失败
- `InternalError()` - 内部服务器错误
- `ServiceUnavailable()` - 服务不可用
- `TooManyRequests()` - 请求过于频繁

**错误处理**:
- `ErrorResponse()` - 使用Error对象响应
- `HandleError()` - 自动处理错误

### 3. 分页响应 (`pagination.go`)

提供分页响应的便利方法：

- `SuccessPaginated()` - 带分页的成功响应
- `SuccessPaginatedFromSearch()` - 从搜索结果创建分页响应
- `SuccessPaginatedFromList()` - 从列表和总数创建分页响应
- `SuccessPaginatedWithOffset()` - 使用offset和limit创建分页响应

## 使用示例

### 基本用法

```go
import "github.com/coze-dev/coze-studio/backend/core/response"

// 成功响应
response.Success(c, data, "操作成功")

// 创建成功
response.Created(c, newItem, "创建成功")

// 更新成功
response.Updated(c, updatedItem, "更新成功")

// 删除成功
response.Deleted(c, "删除成功")
```

### 错误响应

```go
// 请求参数错误
response.BadRequest(c, "请求参数错误", map[string]interface{}{"field": "error"})

// 未授权
response.Unauthorized(c, "未授权访问")

// 资源不存在
response.NotFound(c, "资源不存在")

// 验证失败
response.ValidationFailed(c, "验证失败", validationErrors)

// 内部错误
response.InternalError(c, "内部服务器错误", err)
```

### 分页响应

```go
// 使用PaginationResponse
items := []*TableBody{...}
pagination := service.NewPaginationResponse(page, pageSize, total)
response.SuccessPaginated(c, items, pagination, "获取列表成功")

// 从搜索结果创建
result, _ := searchEngine.FullTextSearch(ctx, req)
response.SuccessPaginatedFromSearch(c, result, items, "搜索成功")

// 从列表和总数创建
response.SuccessPaginatedFromList(c, items, page, pageSize, total, "获取列表成功")
```

### 错误处理

```go
// 使用Error对象
err := response.NewError(response.ErrorCodeNotFound, "资源不存在", http.StatusNotFound)
response.ErrorResponse(c, err)

// 自动处理错误
if err != nil {
    response.HandleError(c, err)
    return
}
```

### 重构前 vs 重构后

**重构前**:
```go
c.JSON(http.StatusOK, gin.H{
    "code":    20000,
    "success": true,
    "data":    data,
    "message": "操作成功",
})
```

**重构后**:
```go
response.Success(c, data, "操作成功")
```

## API参考

### 成功响应

- `Success(c, data, message)` - 成功响应
- `Created(c, data, message)` - 创建成功（HTTP 201）
- `Updated(c, data, message)` - 更新成功
- `Deleted(c, message)` - 删除成功
- `SuccessWithPagination(c, data, pagination, message)` - 带分页的成功响应

### 错误响应

- `BadRequest(c, message, details)` - 请求参数错误（HTTP 400）
- `Unauthorized(c, message)` - 未授权（HTTP 401）
- `Forbidden(c, message)` - 禁止访问（HTTP 403）
- `NotFound(c, message)` - 资源不存在（HTTP 404）
- `Conflict(c, message, details)` - 资源冲突（HTTP 409）
- `ValidationFailed(c, message, details)` - 验证失败（HTTP 422）
- `InternalError(c, message, err)` - 内部服务器错误（HTTP 500）
- `ServiceUnavailable(c, message)` - 服务不可用（HTTP 503）
- `TooManyRequests(c, message)` - 请求过于频繁（HTTP 429）

### 错误处理

- `ErrorResponse(c, err)` - 使用Error对象响应
- `HandleError(c, err)` - 自动处理错误

### 分页响应

- `SuccessPaginated(c, items, pagination, message)` - 带分页的成功响应
- `SuccessPaginatedFromSearch(c, result, items, message)` - 从搜索结果创建
- `SuccessPaginatedFromList(c, items, page, pageSize, total, message)` - 从列表创建
- `SuccessPaginatedWithOffset(c, items, limit, offset, total, message)` - 使用offset创建

## 与现有代码集成

响应框架完全兼容现有的`types.ApiResponse`，可以逐步迁移现有代码。

## 测试

运行测试：
```bash
go test ./core/response/... -v
```

所有测试都已通过 ✅

