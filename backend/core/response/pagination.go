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

import (
	"net/http"

	"github.com/coze-dev/coze-studio/backend/core/service"
	"etrxlite/internal/types"
	"etrxlite/tools/search"

	"github.com/gin-gonic/gin"
)

// PaginatedData 分页数据包装
type PaginatedData struct {
	Items      interface{}                `json:"items"`
	Pagination *service.PaginationResponse `json:"pagination"`
}

// SuccessPaginated 带分页的成功响应（使用PaginatedData）
func SuccessPaginated(c *gin.Context, items interface{}, pagination *service.PaginationResponse, message string) {
	if message == "" {
		message = types.MessageSuccess
	}

	data := PaginatedData{
		Items:      items,
		Pagination: pagination,
	}

	c.JSON(http.StatusOK, types.NewSuccessResponse(data, message))
}

// SuccessPaginatedFromSearch 从搜索结果创建分页响应
func SuccessPaginatedFromSearch(c *gin.Context, result *search.SearchResult, items interface{}, message string) {
	if message == "" {
		message = types.MessageSuccess
	}

	pagination := &service.PaginationResponse{
		Page:       result.Page,
		PageSize:   result.PageSize,
		Total:      int(result.Total),
		TotalPages: result.TotalPages,
	}

	SuccessPaginated(c, items, pagination, message)
}

// SuccessPaginatedFromList 从列表和总数创建分页响应
func SuccessPaginatedFromList(c *gin.Context, items interface{}, page, pageSize, total int, message string) {
	if message == "" {
		message = types.MessageSuccess
	}

	pagination := service.NewPaginationResponse(page, pageSize, total)
	SuccessPaginated(c, items, pagination, message)
}

// SuccessPaginatedWithOffset 使用offset和limit创建分页响应
func SuccessPaginatedWithOffset(c *gin.Context, items interface{}, limit, offset, total int, message string) {
	if message == "" {
		message = types.MessageSuccess
	}

	// 计算页码
	page := 1
	if limit > 0 {
		page = offset/limit + 1
	}

	pagination := service.NewPaginationResponse(page, limit, total)
	SuccessPaginated(c, items, pagination, message)
}

