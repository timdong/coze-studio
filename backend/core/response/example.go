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

package response

// 这个文件展示了如何使用响应框架重构现有API Handler
// 这是一个示例，不会在实际代码中使用

/*
示例：重构后的API Handler

package api

import (
	"net/http"
	"github.com/coze-dev/coze-studio/backend/core/response"
	"github.com/coze-dev/coze-studio/backend/core/service"
	"github.com/coze-dev/coze-studio/backend/core/validation"
)

// TableHandler 使用响应框架
type TableHandler struct {
	tableService *services.EtrxtableService
	validator    *validation.WorkspaceValidator
}

func (h *TableHandler) GetTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID参数", nil)
		return
	}

	// 验证工作空间
	workspaceID := getWorkspaceIDFromContext(c)
	if err := h.validator.ValidateWorkspaceAccess(c.Request.Context(), workspaceID, getUserID(c)); err != nil {
		response.Forbidden(c, "无权访问该工作空间")
		return
	}

	table, err := h.tableService.GetTable(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "表格不存在")
		return
	}

	response.Success(c, table, "获取表格成功")
}

func (h *TableHandler) ListTables(c *gin.Context) {
	workspaceID := getWorkspaceIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	tables, pagination, err := h.tableService.GetTables(c.Request.Context(), workspaceID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取表格列表失败", err)
		return
	}

	response.SuccessPaginated(c, tables, pagination, "获取表格列表成功")
}

func (h *TableHandler) CreateTable(c *gin.Context) {
	var req CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", map[string]interface{}{"error": err.Error()})
		return
	}

	// 验证请求数据
	validator := validation.NewStructValidator()
	validator.Field("Name").
		AddRule(validation.NewRequiredRule()).
		AddRule(validation.NewMinLengthRule(3))
	
	if err := validator.Validate(req); err != nil {
		response.ValidationFailed(c, "验证失败", map[string]interface{}{"errors": err.Error()})
		return
	}

	table, err := h.tableService.CreateTable(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, "创建表格失败", err)
		return
	}

	response.Created(c, table, "创建表格成功")
}

func (h *TableHandler) UpdateTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID参数", nil)
		return
	}

	var req UpdateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", map[string]interface{}{"error": err.Error()})
		return
	}

	err = h.tableService.UpdateTable(c.Request.Context(), uint(id), req)
	if err != nil {
		if err == services.ErrTableNotFound {
			response.NotFound(c, "表格不存在")
			return
		}
		response.InternalError(c, "更新表格失败", err)
		return
	}

	response.Updated(c, nil, "更新表格成功")
}

func (h *TableHandler) DeleteTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID参数", nil)
		return
	}

	err = h.tableService.DeleteTable(c.Request.Context(), uint(id))
	if err != nil {
		if err == services.ErrTableNotFound {
			response.NotFound(c, "表格不存在")
			return
		}
		response.InternalError(c, "删除表格失败", err)
		return
	}

	response.Deleted(c, "删除表格成功")
}
*/

