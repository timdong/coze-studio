# Core Service Framework

## 概述

这是EtrxLite Framework重构计划的基础服务框架，提供了统一的服务基础结构、CRUD操作、分页、健康检查、事务、事件和审计日志能力。

## 文件结构

```
core/service/
├── base.go          # 基础服务接口和实现
├── crud.go          # CRUD操作基类
├── pagination.go    # 分页工具
├── health.go        # 健康检查基类
├── transaction.go   # 事务管理器
├── event.go         # 事件发布能力
├── audit.go         # 审计日志能力
├── example.go       # 使用示例
├── README.md        # 本文档
├── PHASE6_SUMMARY.md # Phase 6总结
├── PHASE7_SUMMARY.md # Phase 7总结
└── *_test.go        # 单元测试
```

## 核心组件

### 1. BaseService (base.go)

基础服务接口和实现，提供：
- `ServiceManager` 接口：定义服务管理器需要提供的能力
- `BaseService` 接口：基础服务接口
- `BaseServiceImpl` 实现：基础服务实现类

**使用示例**:
```go
import "github.com/coze-dev/coze-studio/backend/core/service"

// 创建基础服务
baseService := service.NewBaseService(manager)

// 使用基础服务的方法
logger := baseService.GetLogger()
db := baseService.GetDatabase()
cache := baseService.GetCache()
```

### 2. CRUDService (crud.go)

CRUD操作基类，提供统一的CRUD操作：
- `Create()` - 创建资源
- `GetByID()` - 根据ID获取资源
- `Update()` - 更新资源
- `Delete()` - 删除资源
- `CreateAndGet()` - 创建并获取资源（常用模式）

**使用示例**:
```go
// 创建CRUD服务
base := service.NewBaseService(manager)
crudService := service.NewCRUDService(base, manager.GetCRUD(), "table_bodies")

// 使用CRUD操作
id, err := crudService.Create(ctx, data)
var table ent.TableBody
err = crudService.GetByID(ctx, id, &table)
```

### 3. Pagination (pagination.go)

分页工具，提供：
- `NormalizePagination()` - 规范化分页参数
- `CalculateTotalPages()` - 计算总页数
- `NewPaginationResponse()` - 创建分页响应

**使用示例**:
```go
page, pageSize, offset := service.NormalizePagination(1, 20)
pagination := service.NewPaginationResponse(page, pageSize, total)
```

### 4. HealthChecker (health.go)

健康检查基类，提供：
- `HealthChecker` 接口
- `CompositeHealthChecker` - 组合健康检查器，可以同时检查多个服务

**使用示例**:
```go
checker := service.NewCompositeHealthChecker()
checker.Register("database", dbChecker)
checker.Register("cache", cacheChecker)
err := checker.HealthCheck(ctx)
```

### 5. TransactionManager (transaction.go)

事务管理器，提供：
- `WithTransaction()` - 在事务中执行函数（自动回滚和panic恢复）
- `WithTransactionOptions()` - 使用选项在事务中执行函数

**使用示例**:
```go
txManager := service.NewTransactionManager(client)
err := txManager.WithTransaction(ctx, func(tx *ent.Tx) error {
    // 执行数据库操作
    return nil
})
```

### 6. Event Publishing (event.go)

事件发布能力，提供：
- `PublishEvent()` - 发布事件
- `PublishEventAsync()` - 异步发布事件
- `PublishResourceCreated()` - 发布资源创建事件
- `PublishResourceUpdated()` - 发布资源更新事件
- `PublishResourceDeleted()` - 发布资源删除事件

**使用示例**:
```go
// 发布资源创建事件
err := s.PublishResourceCreated(ctx, "table", tableID, tableData)
```

### 7. Audit Logging (audit.go)

审计日志能力，提供：
- `LogAction()` - 记录操作审计日志
- `LogCreate()` - 记录创建操作
- `LogUpdate()` - 记录更新操作
- `LogDelete()` - 记录删除操作
- `LogSecurityEvent()` - 记录安全事件
- `LogPerformanceMetric()` - 记录性能指标

**使用示例**:
```go
// 记录创建操作
err := s.LogCreate(ctx, "table", fmt.Sprintf("%d", tableID), data)
```

## 集成说明

### 让Manager实现ServiceManager接口

`platform/services.Manager` 已经实现了 `ServiceManager` 接口所需的所有方法：
- `GetDatabase()`
- `GetCache()`
- `GetLogger()`
- `GetClient()`
- `GetCRUD()`
- `GetSearch()`

### Manager实现EventBusProvider和AuditServiceProvider

`platform/services.Manager` 也实现了：
- `GetEventBus()` - 获取事件总线
- `GetAuditService()` - 获取审计服务

### 重构现有服务

参考 `example.go` 中的示例，可以这样重构现有服务：

```go
type EtrxtableService struct {
    *service.CRUDService
}

func NewEtrxtableService(manager *Manager) *EtrxtableService {
    base := service.NewBaseService(manager)
    crudService := service.NewCRUDService(base, manager.GetCRUD(), "table_bodies")
    return &EtrxtableService{
        CRUDService: crudService,
    }
}

// 自动继承Create、GetByID、Update、Delete等方法
func (s *EtrxtableService) CreateTable(ctx context.Context, data map[string]interface{}) (*ent.TableBody, error) {
    var table ent.TableBody
    if err := s.CreateAndGet(ctx, data, &table); err != nil {
        return nil, err
    }
    
    // 发布事件
    _ = s.PublishResourceCreated(ctx, "table", table.ID, table)
    
    // 记录审计日志
    _ = s.LogCreate(ctx, "table", fmt.Sprintf("%d", table.ID), data)
    
    return &table, nil
}
```

## 测试

运行测试：
```bash
go test ./core/service/... -v
```

所有测试都已通过 ✅

## 相关文档

- [Phase 6 总结](PHASE6_SUMMARY.md) - 事务和上下文管理
- [Phase 7 总结](PHASE7_SUMMARY.md) - 事件和审计系统
- [Framework使用指南](../../docs/framework/FRAMEWORK_USAGE_GUIDE.md)
- [迁移指南](../../docs/framework/MIGRATION_GUIDE.md)

## 下一步

所有核心框架已完成，接下来可以：

1. 开始服务重构：将现有服务迁移到新框架
2. 添加集成测试：测试框架与现有服务的集成
3. 性能优化：根据实际使用情况进行性能优化
