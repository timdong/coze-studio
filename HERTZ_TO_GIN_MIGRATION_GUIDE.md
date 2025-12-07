# Hertz 到 Gin 迁移指南

本文档提供将 Coze Studio 的 Hertz 框架迁移到 Gin 的详细指南。

## 📋 迁移概述

### 为什么迁移？

1. **统一技术栈**：与 EtrxLite 保持一致，统一使用 Gin
2. **生态成熟**：Gin 是 Go 生态中最流行的 Web 框架
3. **简单易用**：Gin 的 API 更直观，学习成本更低
4. **兼容性好**：Gin 基于标准库，兼容性更好

### 迁移范围

需要迁移的 Coze Studio 组件：
- ✅ HTTP 路由定义
- ✅ 中间件实现
- ✅ 处理器函数
- ✅ 请求/响应处理
- ✅ 错误处理

---

## 🔄 迁移步骤

### 步骤 1: 依赖更新

#### 1.1 更新 go.mod

```go
// 移除 Hertz 依赖
// github.com/cloudwego/hertz v0.10.2

// 添加 Gin 依赖
require (
    github.com/gin-gonic/gin v1.10.1
)
```

#### 1.2 安装依赖

```bash
go mod tidy
```

---

### 步骤 2: 路由迁移

#### 2.1 基本路由

**Hertz (原)**：
```go
import (
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/app"
)

func main() {
    h := server.Default()
    h.GET("/api/v1/tables", ListTables)
    h.Spin()
}
```

**Gin (迁移后)**：
```go
import (
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    router.GET("/api/v1/tables", ListTables)
    router.Run(":8888")
}
```

#### 2.2 路由组

**Hertz (原)**：
```go
api := h.Group("/api")
v1 := api.Group("/v1")
v1.GET("/tables", ListTables)
v1.POST("/tables", CreateTable)
```

**Gin (迁移后)**：
```go
api := router.Group("/api")
v1 := api.Group("/v1")
{
    v1.GET("/tables", ListTables)
    v1.POST("/tables", CreateTable)
}
```

#### 2.3 路由参数

**Hertz (原)**：
```go
h.GET("/api/v1/tables/:id", GetTable)
```

**Gin (迁移后)**：
```go
router.GET("/api/v1/tables/:id", GetTable)
```

---

### 步骤 3: 处理器迁移

#### 3.1 基本处理器

**Hertz (原)**：
```go
import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

func ListTables(ctx context.Context, c *app.RequestContext) {
    tables := []Table{...}
    c.JSON(200, map[string]interface{}{
        "data": tables,
    })
}
```

**Gin (迁移后)**：
```go
import (
    "github.com/gin-gonic/gin"
)

func ListTables(c *gin.Context) {
    tables := []Table{...}
    c.JSON(200, gin.H{
        "data": tables,
    })
}
```

#### 3.2 请求参数获取

**Hertz (原)**：
```go
func GetTable(ctx context.Context, c *app.RequestContext) {
    id := c.Param("id")
    name := c.Query("name")
    
    var req CreateTableRequest
    c.BindJSON(&req)
}
```

**Gin (迁移后)**：
```go
func GetTable(c *gin.Context) {
    id := c.Param("id")
    name := c.Query("name")
    
    var req CreateTableRequest
    c.ShouldBindJSON(&req)
}
```

#### 3.3 响应处理

**Hertz (原)**：
```go
c.JSON(200, map[string]interface{}{"data": result})
c.String(200, "text")
c.Data(200, "application/json", []byte("data"))
```

**Gin (迁移后)**：
```go
c.JSON(200, gin.H{"data": result})
c.String(200, "text")
c.Data(200, "application/json", []byte("data"))
```

---

### 步骤 4: 中间件迁移

#### 4.1 基本中间件

**Hertz (原)**：
```go
import (
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/middleware/server"
)

func main() {
    h := server.Default()
    h.Use(middleware.AccessLogMW())
    h.Use(middleware.RecoveryMW())
}
```

**Gin (迁移后)**：
```go
import (
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default() // 默认包含 Logger 和 Recovery
    // 或
    router := gin.New()
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
}
```

#### 4.2 自定义中间件

**Hertz (原)**：
```go
func AuthMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatus(401)
            return
        }
        // 验证 token
        c.Next(ctx)
    }
}
```

**Gin (迁移后)**：
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatus(401)
            return
        }
        // 验证 token
        c.Next()
    }
}
```

#### 4.3 中间件注册

**Hertz (原)**：
```go
h.Use(middleware.AccessLogMW())
h.Use(middleware.SessionAuthMW())
```

**Gin (迁移后)**：
```go
router.Use(gin.Logger())
router.Use(AuthMiddleware())
```

---

### 步骤 5: 错误处理

#### 5.1 错误响应

**Hertz (原)**：
```go
func CreateTable(ctx context.Context, c *app.RequestContext) {
    if err != nil {
        c.JSON(400, map[string]interface{}{
            "error": err.Error(),
        })
        return
    }
}
```

**Gin (迁移后)**：
```go
func CreateTable(c *gin.Context) {
    if err != nil {
        c.JSON(400, gin.H{
            "error": err.Error(),
        })
        return
    }
}
```

#### 5.2 统一错误处理

**Hertz (原)**：
```go
h.NoRoute(func(ctx context.Context, c *app.RequestContext) {
    c.JSON(404, map[string]interface{}{
        "error": "Not Found",
    })
})
```

**Gin (迁移后)**：
```go
router.NoRoute(func(c *gin.Context) {
    c.JSON(404, gin.H{
        "error": "Not Found",
    })
})
```

---

### 步骤 6: 上下文处理

#### 6.1 上下文传递

**Hertz (原)**：
```go
func Handler(ctx context.Context, c *app.RequestContext) {
    // 使用 context.Context
    userID := ctx.Value("userID")
}
```

**Gin (迁移后)**：
```go
func Handler(c *gin.Context) {
    // 使用 gin.Context
    userID := c.Get("userID")
    
    // 或获取原始 context
    ctx := c.Request.Context()
}
```

#### 6.2 设置上下文值

**Hertz (原)**：
```go
ctx = context.WithValue(ctx, "userID", userID)
c.Next(ctx)
```

**Gin (迁移后)**：
```go
c.Set("userID", userID)
c.Next()
```

---

## 🛠️ 迁移工具脚本

### 自动转换脚本（示例）

```go
// tools/hertz_to_gin_converter/main.go
package main

import (
    "go/ast"
    "go/parser"
    "go/token"
    "os"
)

// 这是一个简化的示例，实际转换需要更复杂的逻辑
func convertHertzToGin(hertzFilePath string) error {
    // 1. 解析 Hertz 代码文件
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, hertzFilePath, nil, parser.ParseComments)
    if err != nil {
        return err
    }
    
    // 2. 提取路由定义
    // 3. 转换为 Gin 路由
    // 4. 生成 Gin 代码文件
    
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: hertz_to_gin_converter <hertz_file>")
        os.Exit(1)
    }
    
    if err := convertHertzToGin(os.Args[1]); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
```

---

## ✅ 迁移检查清单

### 路由迁移检查

- [ ] 所有 Hertz 路由已转换为 Gin
- [ ] 路由参数正确映射
- [ ] 路由组正确嵌套
- [ ] 测试通过

### 处理器迁移检查

- [ ] 所有处理器函数签名已更新
- [ ] 请求参数获取正确
- [ ] 响应处理正确
- [ ] 错误处理正确
- [ ] 测试通过

### 中间件迁移检查

- [ ] 所有中间件已迁移
- [ ] 中间件顺序正确
- [ ] 自定义中间件功能正常
- [ ] 测试通过

### 上下文处理检查

- [ ] 上下文值传递正确
- [ ] 上下文获取正确
- [ ] 测试通过

---

## 📚 参考资料

- [Gin 官方文档](https://gin-gonic.com/docs/)
- [Hertz 官方文档](https://www.cloudwego.io/zh/docs/hertz/)

---

## 🔍 常见问题

### Q1: 如何处理 Hertz 的 IDL 代码生成？

**A**: Gin 不支持 IDL 代码生成，需要手动编写路由和处理器。可以：
- 保留 IDL 定义作为文档
- 手动实现对应的 Gin 路由和处理器
- 编写脚本辅助生成

### Q2: 如何处理 Hertz 的特殊功能？

**A**: 对于 Hertz 的特殊功能，需要找到 Gin 的等价实现：
- **Hertz 的流式响应**：使用 Gin 的 `c.Stream()` 或 `c.SSEvent()`
- **Hertz 的文件上传**：使用 Gin 的 `c.FormFile()` 或 `c.MultipartForm()`

### Q3: 性能会有影响吗？

**A**: Gin 基于标准库，性能略低于 Hertz，但对于大多数应用场景，性能差异可以忽略。如果确实需要极致性能，可以考虑：
- 使用 Gin 的 Release 模式：`gin.SetMode(gin.ReleaseMode)`
- 优化中间件和处理器逻辑
- 使用连接池和缓存

---

**注意**：本文档会根据实际迁移情况持续更新。

