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

	"github.com/gin-gonic/gin"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}, message string) {
	if message == "" {
		message = types.MessageSuccess
	}
	c.JSON(http.StatusOK, types.NewSuccessResponse(data, message))
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, data interface{}, message string) {
	Success(c, data, message)
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}, message string) {
	if message == "" {
		message = types.MessageCreated
	}
	c.JSON(http.StatusCreated, types.NewSuccessResponse(data, message))
}

// Updated 更新成功响应
func Updated(c *gin.Context, data interface{}, message string) {
	if message == "" {
		message = types.MessageUpdated
	}
	c.JSON(http.StatusOK, types.NewSuccessResponse(data, message))
}

// Deleted 删除成功响应
func Deleted(c *gin.Context, message string) {
	if message == "" {
		message = types.MessageDeleted
	}
	c.JSON(http.StatusOK, types.NewSuccessResponse(nil, message))
}

// SuccessWithPagination 带分页的成功响应
func SuccessWithPagination(c *gin.Context, data interface{}, pagination *service.PaginationResponse, message string) {
	if message == "" {
		message = types.MessageSuccess
	}

	responseData := map[string]interface{}{
		"items": data,
		"pagination": pagination,
	}

	c.JSON(http.StatusOK, types.NewSuccessResponse(responseData, message))
}

// ErrorResponse 错误响应（内部函数）
func errorResponse(c *gin.Context, status int, code int, message string, details map[string]interface{}) {
	c.JSON(status, types.NewErrorResponse(code, message, string(GetErrorCode(status)), details))
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string, details map[string]interface{}) {
	if message == "" {
		message = "请求参数错误"
	}
	errorResponse(c, http.StatusBadRequest, types.ResponseCodeBadRequest, message, details)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = types.MessageUnauthorized
	}
	errorResponse(c, http.StatusUnauthorized, types.ResponseCodeUnauthorized, message, nil)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = types.MessageForbidden
	}
	errorResponse(c, http.StatusForbidden, types.ResponseCodeForbidden, message, nil)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = types.MessageNotFound
	}
	errorResponse(c, http.StatusNotFound, types.ResponseCodeNotFound, message, nil)
}

// Conflict 资源冲突
func Conflict(c *gin.Context, message string, details map[string]interface{}) {
	if message == "" {
		message = "资源冲突"
	}
	errorResponse(c, http.StatusConflict, types.ResponseCodeConflict, message, details)
}

// ValidationFailed 验证失败
func ValidationFailed(c *gin.Context, message string, details map[string]interface{}) {
	if message == "" {
		message = types.MessageValidationFailed
	}
	errorResponse(c, http.StatusUnprocessableEntity, types.ResponseCodeValidationFailed, message, details)
}

// InternalError 内部服务器错误
func InternalError(c *gin.Context, message string, err error) {
	if message == "" {
		message = types.MessageInternalError
	}
	details := make(map[string]interface{})
	if err != nil {
		details["error"] = err.Error()
	}
	errorResponse(c, http.StatusInternalServerError, types.ResponseCodeInternalError, message, details)
}

// ServiceUnavailable 服务不可用
func ServiceUnavailable(c *gin.Context, message string) {
	if message == "" {
		message = "服务不可用"
	}
	errorResponse(c, http.StatusServiceUnavailable, types.ResponseCodeServiceUnavailable, message, nil)
}

// TooManyRequests 请求过于频繁
func TooManyRequests(c *gin.Context, message string) {
	if message == "" {
		message = "请求过于频繁"
	}
	errorResponse(c, http.StatusTooManyRequests, types.ResponseCodeTooManyRequests, message, nil)
}

// ErrorResponse 使用Error对象响应
func ErrorResponse(c *gin.Context, err *Error) {
	if err == nil {
		InternalError(c, "未知错误", nil)
		return
	}

	code := GetResponseCode(err.Status)
	errorResponse(c, err.Status, code, err.Message, err.Details)
}

// HandleError 处理错误（自动判断错误类型）
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 如果是响应错误，直接使用
	if respErr, ok := err.(*Error); ok {
		ErrorResponse(c, respErr)
		return
	}

	// 否则包装为内部错误
	InternalError(c, err.Error(), err)
}

