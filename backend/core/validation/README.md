# Validation Framework

## 概述

验证框架提供了统一的验证能力，包括字段验证、结构体验证和工作空间验证。

## 功能特性

- ✅ 类型安全的验证规则
- ✅ 支持多种内置验证规则
- ✅ 字段级和结构体级验证
- ✅ 工作空间验证
- ✅ 可组合的验证器
- ✅ 自定义验证规则支持

## 核心组件

### 1. 验证器接口 (`validator.go`)

- `Validator` - 基础验证器接口
- `FieldValidator` - 字段验证器
- `StructValidator` - 结构体验证器
- `CompositeValidator` - 组合验证器

### 2. 验证规则 (`rules.go`)

支持以下内置规则：

- `RequiredRule` - 必填规则
- `MinLengthRule` / `MaxLengthRule` / `LengthRule` - 长度规则
- `EmailRule` - 邮箱规则
- `URLRule` - URL规则
- `MinRule` / `MaxRule` / `RangeRule` - 数值范围规则
- `InRule` / `NotInRule` - 列表规则
- `PatternRule` - 正则表达式规则
- `CustomRule` - 自定义规则
- `TrimmedStringRule` - 去除空白字符后验证

### 3. 工作空间验证 (`workspace.go`)

- `ValidateWorkspaceExists` - 验证工作空间是否存在
- `ValidateWorkspaceActive` - 验证工作空间是否激活
- `ValidateWorkspaceAccess` - 验证用户是否有权限访问工作空间
- `BatchValidateWorkspaces` - 批量验证工作空间

## 使用示例

### 字段验证

```go
import "github.com/coze-dev/coze-studio/backend/core/validation"

// 创建字段验证器
validator := validation.NewFieldValidator("email").
    AddRule(validation.NewRequiredRule()).
    AddRule(validation.NewEmailRule())

// 验证
err := validator.Validate("user@example.com")
```

### 结构体验证

```go
type UserInput struct {
    Name  string
    Email string
}

validator := validation.NewStructValidator()
validator.Field("Name").
    AddRule(validation.NewRequiredRule()).
    AddRule(validation.NewMinLengthRule(3))
validator.Field("Email").
    AddRule(validation.NewRequiredRule()).
    AddRule(validation.NewEmailRule())

user := UserInput{
    Name:  "John",
    Email: "john@example.com",
}

err := validator.Validate(user)
```

### 工作空间验证

```go
import "github.com/coze-dev/coze-studio/backend/core/validation"

validator := validation.NewWorkspaceValidator(client)

// 验证工作空间是否存在
err := validator.ValidateWorkspaceExists(ctx, workspaceID)

// 验证工作空间是否激活
err := validator.ValidateWorkspaceActive(ctx, workspaceID)

// 验证用户是否有权限访问
err := validator.ValidateWorkspaceAccess(ctx, workspaceID, userID)
```

### 组合验证器

```go
composite := validation.NewCompositeValidator().
    Add(validation.NewRequiredRule()).
    Add(validation.NewMinLengthRule(5))

err := composite.Validate("test")
```

## API参考

### 验证器

- `NewFieldValidator(fieldName)` - 创建字段验证器
- `NewStructValidator()` - 创建结构体验证器
- `NewCompositeValidator()` - 创建组合验证器

### 验证规则

- `NewRequiredRule()` - 必填规则
- `NewMinLengthRule(min)` - 最小长度规则
- `NewMaxLengthRule(max)` - 最大长度规则
- `NewLengthRule(min, max)` - 长度范围规则
- `NewEmailRule()` - 邮箱规则
- `NewURLRule()` - URL规则
- `NewMinRule(min)` - 最小值规则
- `NewMaxRule(max)` - 最大值规则
- `NewRangeRule(min, max)` - 范围规则
- `NewInRule(values...)` - 在列表中规则
- `NewNotInRule(values...)` - 不在列表中规则
- `NewPatternRule(pattern, message)` - 正则表达式规则
- `NewCustomRule(func)` - 自定义规则

### 工作空间验证

- `NewWorkspaceValidator(client)` - 创建工作空间验证器
- `ValidateWorkspaceExists(ctx, workspaceID)` - 验证存在
- `ValidateWorkspaceActive(ctx, workspaceID)` - 验证激活
- `ValidateWorkspaceAccess(ctx, workspaceID, userID)` - 验证访问权限

## 测试

运行测试：
```bash
go test ./core/validation/... -v
```

所有测试都已通过 ✅

