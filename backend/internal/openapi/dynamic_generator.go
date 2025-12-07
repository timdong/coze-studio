/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package openapi

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// OpenAPI3 OpenAPI 3.0规范结构
type OpenAPI3 struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Servers    []Server            `json:"servers"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
	Tags       []Tag               `json:"tags"`
}

// Info API信息
type Info struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Version        string  `json:"version"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
	TermsOfService   string  `json:"termsOfService"`
}

// Contact 联系信息
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// License 许可证信息
type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Server 服务器信息
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// PathItem 路径项
type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
}

// Operation 操作
type Operation struct {
	Tags        []string              `json:"tags"`
	Summary     string                `json:"summary"`
	Description string                `json:"description"`
	OperationID string                `json:"operationId"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []map[string][]string  `json:"security,omitempty"`
}

// Parameter 参数
type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required"`
	Schema      *Schema     `json:"schema"`
	Example     interface{} `json:"example,omitempty"`
}

// RequestBody 请求体
type RequestBody struct {
	Description string               `json:"description"`
	Required    bool                 `json:"required"`
	Content     map[string]MediaType `json:"content"`
}

// MediaType 媒体类型
type MediaType struct {
	Schema *Schema `json:"schema"`
}

// Schema 模式
type Schema struct {
	Type                 string             `json:"type"`
	Format               string             `json:"format,omitempty"`
	Description          string             `json:"description,omitempty"`
	Example              interface{}        `json:"example,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	Required             []string           `json:"required,omitempty"`
	Items                *Schema            `json:"items,omitempty"`
	AdditionalProperties *Schema            `json:"additionalProperties,omitempty"`
}

// Response 响应
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// Components 组件
type Components struct {
	Schemas         map[string]*Schema        `json:"schemas,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

// SecurityScheme 安全方案
type SecurityScheme struct {
	Type         string `json:"type"`
	Description  string `json:"description"`
	Name         string `json:"name,omitempty"`
	In           string `json:"in,omitempty"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
}

// Tag 标签
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DynamicGenerator 动态OpenAPI生成器
type DynamicGenerator struct {
	baseURL string
	version string
	router  *gin.Engine
}

// NewDynamicGenerator 创建新的动态生成器
func NewDynamicGenerator(baseURL, version string, router *gin.Engine) *DynamicGenerator {
	return &DynamicGenerator{
		baseURL: baseURL,
		version: version,
		router:  router,
	}
}

// Generate 生成OpenAPI文档
func (g *DynamicGenerator) Generate() *OpenAPI3 {
	return &OpenAPI3{
		OpenAPI:    "3.0.0",
		Info:       g.generateInfo(),
		Servers:    g.generateServers(),
		Paths:      g.generatePathsFromRouter(),
		Components: g.generateComponents(),
		Tags:       g.generateTags(),
	}
}

// generateInfo 生成API信息
func (g *DynamicGenerator) generateInfo() Info {
	return Info{
		Title:          "Coze Studio API",
		Description:    "Coze Studio 融合平台API文档",
		Version:        g.version,
		TermsOfService: "http://swagger.io/terms/",
		Contact: Contact{
			Name:  "Coze Studio Team",
			Email: "support@coze-studio.com",
			URL:   "https://coze-studio.com",
		},
		License: License{
			Name: "Apache 2.0",
			URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
		},
	}
}

// generateServers 生成服务器信息
func (g *DynamicGenerator) generateServers() []Server {
	return []Server{
		{
			URL:         g.baseURL,
			Description: "开发环境",
		},
	}
}

// generatePathsFromRouter 从路由器动态生成路径
func (g *DynamicGenerator) generatePathsFromRouter() map[string]PathItem {
	paths := make(map[string]PathItem)

	// 获取路由器的所有路由
	routes := g.router.Routes()

	for _, route := range routes {
		path := route.Path
		method := strings.ToLower(route.Method)

		// 跳过一些不需要的路由
		if path == "/" || path == "/health" || path == "/docs" ||
			path == "/swagger/*any" || path == "/api/v1/openapi.json" ||
			strings.HasPrefix(path, "/static") || strings.HasPrefix(path, "/admin") {
			continue
		}

		// 创建路径项
		if _, exists := paths[path]; !exists {
			paths[path] = PathItem{}
		}

		pathItem := paths[path]

		// 根据方法设置操作
		operation := g.createOperation(path, method)

		switch method {
		case "get":
			pathItem.Get = operation
		case "post":
			pathItem.Post = operation
		case "put":
			pathItem.Put = operation
		case "delete":
			pathItem.Delete = operation
		case "patch":
			pathItem.Patch = operation
		}

		paths[path] = pathItem
	}

	return paths
}

// createOperation 创建操作
func (g *DynamicGenerator) createOperation(path, method string) *Operation {
	// 根据路径和方法推断操作信息
	tags := g.inferTags(path)
	summary := g.inferSummary(path, method)
	operationID := g.inferOperationID(path, method)

	operation := &Operation{
		Tags:        tags,
		Summary:     summary,
		Description: g.inferDescription(path, method),
		OperationID: operationID,
		Parameters:  g.inferParameters(path),
		Responses:   g.inferResponses(path, method),
	}

	// 为POST和PUT请求添加请求体
	if method == "post" || method == "put" {
		operation.RequestBody = g.inferRequestBody(path, method)
	}

	return operation
}

// inferTags 推断标签
func (g *DynamicGenerator) inferTags(path string) []string {
	// 根据路径推断标签
	if strings.Contains(path, "/mas/") {
		return []string{"Multiple Agents System (MAS)"}
	} else if strings.Contains(path, "/extensions/etrxtable") {
		return []string{"EtrxTable"}
	} else if strings.Contains(path, "/extensions/er-diagram") {
		return []string{"ER Diagram"}
	} else if strings.Contains(path, "/extensions/code-generation") {
		return []string{"Code Generation"}
	} else if strings.Contains(path, "/tables") {
		return []string{"表格管理"}
	} else if strings.Contains(path, "/users") {
		return []string{"用户管理"}
	} else if strings.Contains(path, "/roles") {
		return []string{"角色管理"}
	} else if strings.Contains(path, "/workspaces") {
		return []string{"工作空间"}
	} else if strings.Contains(path, "/documents") {
		return []string{"文档管理"}
	} else if strings.Contains(path, "/conversations") {
		return []string{"对话管理"}
	} else if strings.Contains(path, "/auth") {
		return []string{"认证"}
	} else if path == "/health" || path == "/api/v1/health" {
		return []string{"系统"}
	}
	return []string{"其他"}
}

// inferSummary 推断摘要
func (g *DynamicGenerator) inferSummary(path, method string) string {
	entity := g.getEntityFromPath(path)

	switch method {
	case "get":
		if strings.Contains(path, "/:id") || strings.Contains(path, "/{id}") {
			return "获取" + entity + "详情"
		}
		return "获取" + entity + "列表"
	case "post":
		return "创建" + entity
	case "put":
		return "更新" + entity
	case "delete":
		return "删除" + entity
	}
	return "操作" + entity
}

// inferDescription 推断描述
func (g *DynamicGenerator) inferDescription(path, method string) string {
	entity := g.getEntityFromPath(path)

	switch method {
	case "get":
		if strings.Contains(path, "/:id") || strings.Contains(path, "/{id}") {
			return "根据ID获取" + entity + "详情"
		}
		return "分页获取" + entity + "列表"
	case "post":
		return "创建新的" + entity
	case "put":
		return "更新" + entity + "信息"
	case "delete":
		return "删除指定的" + entity
	}
	return "对" + entity + "进行操作"
}

// inferOperationID 推断操作ID
func (g *DynamicGenerator) inferOperationID(path, method string) string {
	entity := g.getEntityFromPath(path)
	return strings.ToLower(method) + entity
}

// inferParameters 推断参数
func (g *DynamicGenerator) inferParameters(path string) []Parameter {
	var parameters []Parameter

	// 检查路径参数（Gin使用:param格式）
	if strings.Contains(path, "/:id") {
		parameters = append(parameters, Parameter{
			Name:        "id",
			In:          "path",
			Description: "资源ID",
			Required:    true,
			Schema: &Schema{
				Type:   "integer",
				Format: "int64",
			},
		})
	}

	// 为列表接口添加分页参数
	if !strings.Contains(path, "/:id") && !strings.Contains(path, "/{id}") {
		// 检查是否是列表接口（以复数形式结尾）
		parts := strings.Split(path, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			if strings.HasSuffix(lastPart, "s") && lastPart != "docs" {
				parameters = append(parameters, Parameter{
					Name:        "page",
					In:          "query",
					Description: "页码",
					Required:    false,
					Schema: &Schema{
						Type: "integer",
					},
				})
				parameters = append(parameters, Parameter{
					Name:        "page_size",
					In:          "query",
					Description: "每页数量",
					Required:    false,
					Schema: &Schema{
						Type: "integer",
					},
				})
			}
		}
	}

	return parameters
}

// inferRequestBody 推断请求体
func (g *DynamicGenerator) inferRequestBody(path, method string) *RequestBody {
	entity := g.getEntityFromPath(path)

	return &RequestBody{
		Description: entity + "数据",
		Required:    true,
		Content: map[string]MediaType{
			"application/json": {
				Schema: &Schema{
					Type: "object",
					Properties: map[string]*Schema{
						"name": {
							Type:        "string",
							Description: entity + "名称",
						},
						"description": {
							Type:        "string",
							Description: entity + "描述",
						},
					},
				},
			},
		},
	}
}

// inferResponses 推断响应
func (g *DynamicGenerator) inferResponses(path, method string) map[string]Response {
	responses := make(map[string]Response)

	// 通用响应
	responses["200"] = Response{
		Description: "成功",
		Content: map[string]MediaType{
			"application/json": {
				Schema: &Schema{
					Type: "object",
					Properties: map[string]*Schema{
						"success": {
							Type:    "boolean",
							Example: true,
						},
						"data": {
							Type: "object",
						},
						"code": {
							Type:    "integer",
							Example: 20000,
						},
						"message": {
							Type:    "string",
							Example: "操作成功",
						},
					},
				},
			},
		},
	}

	responses["400"] = Response{
		Description: "请求参数错误",
	}

	responses["401"] = Response{
		Description: "未授权",
	}

	responses["404"] = Response{
		Description: "资源不存在",
	}

	responses["500"] = Response{
		Description: "服务器内部错误",
	}

	return responses
}

// generateComponents 生成组件
func (g *DynamicGenerator) generateComponents() Components {
	return Components{
		SecuritySchemes: map[string]SecurityScheme{
			"BearerAuth": {
				Type:        "http",
				Scheme:      "bearer",
				BearerFormat: "JWT",
				Description: "JWT认证令牌",
			},
		},
	}
}

// generateTags 生成标签
func (g *DynamicGenerator) generateTags() []Tag {
	return []Tag{
		{Name: "Multiple Agents System (MAS)", Description: "多智能体系统相关API，包括Agent管理、任务调度、会话管理、记忆系统、性能指标等"},
		{Name: "EtrxTable", Description: "多维表格系统相关API"},
		{Name: "ER Diagram", Description: "ER图编辑器相关API"},
		{Name: "Code Generation", Description: "代码生成系统相关API"},
		{Name: "表格管理", Description: "动态表格管理相关API，包括表格、列、行的增删改查及数据迁移"},
		{Name: "用户管理", Description: "用户相关的API接口"},
		{Name: "角色管理", Description: "角色相关的API接口"},
		{Name: "工作空间", Description: "工作空间相关的API接口"},
		{Name: "文档管理", Description: "文档相关的API接口"},
		{Name: "对话管理", Description: "对话相关的API接口"},
		{Name: "认证", Description: "认证相关的API接口"},
		{Name: "系统", Description: "系统相关的API接口"},
	}
}

// getEntityFromPath 从路径中提取实体名称
func (g *DynamicGenerator) getEntityFromPath(path string) string {
	// 移除API前缀
	path = strings.TrimPrefix(path, "/api/v1/")
	path = strings.TrimPrefix(path, "/extensions/")

	// 分割路径
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return "资源"
	}

	// 获取第一个部分作为实体名称
	entity := parts[0]

	// 移除复数形式
	if strings.HasSuffix(entity, "s") {
		entity = strings.TrimSuffix(entity, "s")
	}

	// 处理特殊实体名称
	switch entity {
	case "user":
		return "用户"
	case "role":
		return "角色"
	case "workspace":
		return "工作空间"
	case "table":
		return "表格"
	case "document":
		return "文档"
	case "conversation":
		return "对话"
	case "etrxtable":
		return "EtrxTable"
	case "er-diagram":
		return "ER图"
	case "code-generation":
		return "代码生成"
	case "mas":
		return "多智能体"
	default:
		return entity
	}
}

