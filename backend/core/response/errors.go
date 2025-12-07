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
	"fmt"
	"net/http"

	"etrxlite/internal/types"
)

// ErrorCode 错误代码类型
type ErrorCode string

const (
	// 客户端错误
	ErrorCodeBadRequest       ErrorCode = "BAD_REQUEST"
	ErrorCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrorCodeMethodNotAllowed ErrorCode = "METHOD_NOT_ALLOWED"
	ErrorCodeConflict         ErrorCode = "CONFLICT"
	ErrorCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrorCodeTooManyRequests  ErrorCode = "TOO_MANY_REQUESTS"

	// 服务器错误
	ErrorCodeInternalError        ErrorCode = "INTERNAL_ERROR"
	ErrorCodeNotImplemented       ErrorCode = "NOT_IMPLEMENTED"
	ErrorCodeServiceUnavailable   ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeDatabaseError        ErrorCode = "DATABASE_ERROR"
	ErrorCodeExternalServiceError ErrorCode = "EXTERNAL_SERVICE_ERROR"
)

// Error 响应错误
type Error struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Status  int                    `json:"-"` // HTTP状态码，不序列化到JSON
}

// Error 实现error接口
func (e *Error) Error() string {
	return e.Message
}

// NewError 创建错误
func NewError(code ErrorCode, message string, status int) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Status:  status,
		Details: make(map[string]interface{}),
	}
}

// WithDetails 添加错误详情
func (e *Error) WithDetails(key string, value interface{}) *Error {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithDetailsMap 批量添加错误详情
func (e *Error) WithDetailsMap(details map[string]interface{}) *Error {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// 预定义的错误
var (
	ErrBadRequest       = NewError(ErrorCodeBadRequest, "请求参数错误", http.StatusBadRequest)
	ErrUnauthorized     = NewError(ErrorCodeUnauthorized, "未授权访问", http.StatusUnauthorized)
	ErrForbidden        = NewError(ErrorCodeForbidden, "禁止访问", http.StatusForbidden)
	ErrNotFound         = NewError(ErrorCodeNotFound, "资源不存在", http.StatusNotFound)
	ErrMethodNotAllowed = NewError(ErrorCodeMethodNotAllowed, "方法不允许", http.StatusMethodNotAllowed)
	ErrConflict         = NewError(ErrorCodeConflict, "资源冲突", http.StatusConflict)
	ErrValidationFailed = NewError(ErrorCodeValidationFailed, "验证失败", http.StatusUnprocessableEntity)
	ErrTooManyRequests  = NewError(ErrorCodeTooManyRequests, "请求过于频繁", http.StatusTooManyRequests)

	ErrInternalError        = NewError(ErrorCodeInternalError, "内部服务器错误", http.StatusInternalServerError)
	ErrNotImplemented       = NewError(ErrorCodeNotImplemented, "功能未实现", http.StatusNotImplemented)
	ErrServiceUnavailable   = NewError(ErrorCodeServiceUnavailable, "服务不可用", http.StatusServiceUnavailable)
	ErrDatabaseError        = NewError(ErrorCodeDatabaseError, "数据库错误", http.StatusInternalServerError)
	ErrExternalServiceError = NewError(ErrorCodeExternalServiceError, "外部服务错误", http.StatusBadGateway)
)

// GetErrorCode 根据HTTP状态码获取错误代码
func GetErrorCode(status int) ErrorCode {
	switch status {
	case http.StatusBadRequest:
		return ErrorCodeBadRequest
	case http.StatusUnauthorized:
		return ErrorCodeUnauthorized
	case http.StatusForbidden:
		return ErrorCodeForbidden
	case http.StatusNotFound:
		return ErrorCodeNotFound
	case http.StatusMethodNotAllowed:
		return ErrorCodeMethodNotAllowed
	case http.StatusConflict:
		return ErrorCodeConflict
	case http.StatusUnprocessableEntity:
		return ErrorCodeValidationFailed
	case http.StatusTooManyRequests:
		return ErrorCodeTooManyRequests
	case http.StatusInternalServerError:
		return ErrorCodeInternalError
	case http.StatusNotImplemented:
		return ErrorCodeNotImplemented
	case http.StatusServiceUnavailable:
		return ErrorCodeServiceUnavailable
	case http.StatusBadGateway:
		return ErrorCodeExternalServiceError
	default:
		return ErrorCodeInternalError
	}
}

// GetResponseCode 根据HTTP状态码获取响应代码（兼容现有代码）
func GetResponseCode(status int) int {
	switch status {
	case http.StatusOK:
		return types.ResponseCodeSuccess
	case http.StatusBadRequest:
		return types.ResponseCodeBadRequest
	case http.StatusUnauthorized:
		return types.ResponseCodeUnauthorized
	case http.StatusForbidden:
		return types.ResponseCodeForbidden
	case http.StatusNotFound:
		return types.ResponseCodeNotFound
	case http.StatusMethodNotAllowed:
		return types.ResponseCodeMethodNotAllowed
	case http.StatusConflict:
		return types.ResponseCodeConflict
	case http.StatusUnprocessableEntity:
		return types.ResponseCodeValidationFailed
	case http.StatusTooManyRequests:
		return types.ResponseCodeTooManyRequests
	case http.StatusInternalServerError:
		return types.ResponseCodeInternalError
	case http.StatusNotImplemented:
		return types.ResponseCodeNotImplemented
	case http.StatusServiceUnavailable:
		return types.ResponseCodeServiceUnavailable
	case http.StatusBadGateway:
		return types.ResponseCodeExternalServiceError
	default:
		return types.ResponseCodeInternalError
	}
}

// WrapError 包装标准error为响应错误
func WrapError(err error, defaultError *Error) *Error {
	if err == nil {
		return nil
	}

	// 如果已经是Error类型，直接返回
	if respErr, ok := err.(*Error); ok {
		return respErr
	}

	// 创建新的错误
	wrapped := &Error{
		Code:    defaultError.Code,
		Message: err.Error(),
		Status:  defaultError.Status,
		Details: make(map[string]interface{}),
	}

	// 复制默认错误的详情
	if defaultError.Details != nil {
		for k, v := range defaultError.Details {
			wrapped.Details[k] = v
		}
	}

	return wrapped
}

// NewErrorFromString 从错误字符串创建错误
func NewErrorFromString(message string, status int) *Error {
	return &Error{
		Code:    GetErrorCode(status),
		Message: message,
		Status:  status,
		Details: make(map[string]interface{}),
	}
}

// NewErrorf 使用格式化字符串创建错误
func NewErrorf(status int, format string, args ...interface{}) *Error {
	return &Error{
		Code:    GetErrorCode(status),
		Message: fmt.Sprintf(format, args...),
		Status:  status,
		Details: make(map[string]interface{}),
	}
}

