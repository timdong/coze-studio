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

package service

import (
	"math"
)

// PaginationRequest 分页请求
type PaginationRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// NormalizePagination 规范化分页参数
// 返回规范化后的page、pageSize和offset
func NormalizePagination(page, pageSize int) (normalizedPage, normalizedPageSize, offset int) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // 最大100条
	}
	offset = (page - 1) * pageSize
	return page, pageSize, offset
}

// CalculateTotalPages 计算总页数
func CalculateTotalPages(total, pageSize int) int {
	if total == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse(page, pageSize, total int) *PaginationResponse {
	return &PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: CalculateTotalPages(total, pageSize),
	}
}

// NormalizePaginationRequest 规范化分页请求
func NormalizePaginationRequest(req *PaginationRequest) (page, pageSize, offset int) {
	return NormalizePagination(req.Page, req.PageSize)
}

