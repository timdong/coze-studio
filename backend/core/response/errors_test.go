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
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(ErrorCodeBadRequest, "test error", http.StatusBadRequest)
	if err.Code != ErrorCodeBadRequest {
		t.Errorf("expected code %s, got %s", ErrorCodeBadRequest, err.Code)
	}
	if err.Message != "test error" {
		t.Errorf("expected message 'test error', got '%s'", err.Message)
	}
	if err.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, err.Status)
	}
}

func TestError_WithDetails(t *testing.T) {
	err := NewError(ErrorCodeBadRequest, "test error", http.StatusBadRequest)
	err.WithDetails("field", "value")

	if err.Details["field"] != "value" {
		t.Errorf("expected detail 'value', got '%v'", err.Details["field"])
	}
}

func TestError_WithDetailsMap(t *testing.T) {
	err := NewError(ErrorCodeBadRequest, "test error", http.StatusBadRequest)
	err.WithDetailsMap(map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	})

	if len(err.Details) != 2 {
		t.Errorf("expected 2 details, got %d", len(err.Details))
	}
}

func TestGetErrorCode(t *testing.T) {
	tests := []struct {
		status   int
		expected ErrorCode
	}{
		{http.StatusBadRequest, ErrorCodeBadRequest},
		{http.StatusUnauthorized, ErrorCodeUnauthorized},
		{http.StatusForbidden, ErrorCodeForbidden},
		{http.StatusNotFound, ErrorCodeNotFound},
		{http.StatusInternalServerError, ErrorCodeInternalError},
		{999, ErrorCodeInternalError}, // 未知状态码
	}

	for _, tt := range tests {
		result := GetErrorCode(tt.status)
		if result != tt.expected {
			t.Errorf("status %d: expected %s, got %s", tt.status, tt.expected, result)
		}
	}
}

func TestGetResponseCode(t *testing.T) {
	tests := []struct {
		status   int
		expected int
	}{
		{http.StatusOK, 20000},
		{http.StatusBadRequest, 40000},
		{http.StatusUnauthorized, 40001},
		{http.StatusForbidden, 40003},
		{http.StatusNotFound, 40004},
		{http.StatusInternalServerError, 50000},
	}

	for _, tt := range tests {
		result := GetResponseCode(tt.status)
		if result != tt.expected {
			t.Errorf("status %d: expected %d, got %d", tt.status, tt.expected, result)
		}
	}
}

func TestWrapError(t *testing.T) {
	originalErr := NewError(ErrorCodeBadRequest, "original error", http.StatusBadRequest)
	wrapped := WrapError(originalErr, ErrInternalError)

	if wrapped.Code != ErrorCodeBadRequest {
		t.Errorf("expected code %s, got %s", ErrorCodeBadRequest, wrapped.Code)
	}

	// 测试nil错误
	if WrapError(nil, ErrInternalError) != nil {
		t.Error("expected nil for nil error")
	}
}

func TestNewErrorFromString(t *testing.T) {
	err := NewErrorFromString("test error", http.StatusBadRequest)
	if err.Message != "test error" {
		t.Errorf("expected message 'test error', got '%s'", err.Message)
	}
	if err.Code != ErrorCodeBadRequest {
		t.Errorf("expected code %s, got %s", ErrorCodeBadRequest, err.Code)
	}
}

func TestNewErrorf(t *testing.T) {
	err := NewErrorf(http.StatusBadRequest, "test %s", "error")
	expected := "test error"
	if err.Message != expected {
		t.Errorf("expected message '%s', got '%s'", expected, err.Message)
	}
}

