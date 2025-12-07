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

package errors

import (
	"fmt"
	"runtime"
)

// ServiceError 服务错误
// 提供统一的错误包装和追踪能力
type ServiceError struct {
	Code      string
	Message   string
	Cause     error
	Stack     string
	Service   string
	Operation string
}

// Error 实现error接口
func (e *ServiceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", e.Service, e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s:%s] %s", e.Service, e.Operation, e.Message)
}

// Unwrap 返回原始错误
// 实现errors.Unwrap接口，支持错误链
func (e *ServiceError) Unwrap() error {
	return e.Cause
}

// NewServiceError 创建服务错误
func NewServiceError(service, operation, message string, cause error) *ServiceError {
	stack := captureStack()
	return &ServiceError{
		Code:      "SERVICE_ERROR",
		Message:   message,
		Cause:     cause,
		Stack:     stack,
		Service:   service,
		Operation: operation,
	}
}

// WrapError 包装错误
// 如果错误已经是ServiceError，直接返回
// 否则创建新的ServiceError包装原始错误
func WrapError(service, operation, message string, err error) error {
	if err == nil {
		return nil
	}

	// 如果已经是ServiceError，直接返回
	if se, ok := err.(*ServiceError); ok {
		return se
	}

	return NewServiceError(service, operation, message, err)
}

// WrapErrorf 使用格式化字符串包装错误
func WrapErrorf(service, operation string, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	message := fmt.Sprintf(format, args...)
	return WrapError(service, operation, message, err)
}

// IsServiceError 检查错误是否是ServiceError
func IsServiceError(err error) bool {
	_, ok := err.(*ServiceError)
	return ok
}

// GetServiceError 获取ServiceError，如果不是则返回nil
func GetServiceError(err error) *ServiceError {
	if se, ok := err.(*ServiceError); ok {
		return se
	}
	return nil
}

// GetService 从错误中提取服务名
func GetService(err error) string {
	if se := GetServiceError(err); se != nil {
		return se.Service
	}
	return ""
}

// GetOperation 从错误中提取操作名
func GetOperation(err error) string {
	if se := GetServiceError(err); se != nil {
		return se.Operation
	}
	return ""
}

// captureStack 捕获调用栈
func captureStack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

