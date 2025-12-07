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
	"errors"
	"testing"
)

func TestNewServiceError(t *testing.T) {
	cause := errors.New("original error")
	serviceErr := NewServiceError("test-service", "test-operation", "test message", cause)

	if serviceErr.Service != "test-service" {
		t.Errorf("expected service 'test-service', got '%s'", serviceErr.Service)
	}
	if serviceErr.Operation != "test-operation" {
		t.Errorf("expected operation 'test-operation', got '%s'", serviceErr.Operation)
	}
	if serviceErr.Message != "test message" {
		t.Errorf("expected message 'test message', got '%s'", serviceErr.Message)
	}
	if serviceErr.Cause != cause {
		t.Errorf("expected cause to be original error")
	}
	if serviceErr.Stack == "" {
		t.Error("expected stack trace to be captured")
	}
}

func TestServiceError_Error(t *testing.T) {
	t.Run("with cause", func(t *testing.T) {
		cause := errors.New("original error")
		serviceErr := NewServiceError("test-service", "test-operation", "test message", cause)
		errMsg := serviceErr.Error()
		if errMsg == "" {
			t.Error("expected error message")
		}
	})

	t.Run("without cause", func(t *testing.T) {
		serviceErr := NewServiceError("test-service", "test-operation", "test message", nil)
		errMsg := serviceErr.Error()
		if errMsg == "" {
			t.Error("expected error message")
		}
	})
}

func TestServiceError_Unwrap(t *testing.T) {
	cause := errors.New("original error")
	serviceErr := NewServiceError("test-service", "test-operation", "test message", cause)

	unwrapped := serviceErr.Unwrap()
	if unwrapped != cause {
		t.Errorf("expected unwrapped error to be original cause")
	}
}

func TestWrapError(t *testing.T) {
	t.Run("wrap nil error", func(t *testing.T) {
		err := WrapError("test-service", "test-operation", "test message", nil)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	t.Run("wrap standard error", func(t *testing.T) {
		originalErr := errors.New("original error")
		wrappedErr := WrapError("test-service", "test-operation", "test message", originalErr)

		if wrappedErr == nil {
			t.Error("expected wrapped error")
		}

		serviceErr, ok := wrappedErr.(*ServiceError)
		if !ok {
			t.Error("expected ServiceError type")
		}
		if serviceErr.Cause != originalErr {
			t.Errorf("expected cause to be original error")
		}
	})

	t.Run("wrap ServiceError", func(t *testing.T) {
		originalServiceErr := NewServiceError("original-service", "original-operation", "original message", nil)
		wrappedErr := WrapError("test-service", "test-operation", "test message", originalServiceErr)

		// 应该直接返回原始的ServiceError
		if wrappedErr != originalServiceErr {
			t.Errorf("expected original ServiceError to be returned")
		}
	})
}

func TestWrapErrorf(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapErrorf("test-service", "test-operation", originalErr, "formatted message: %s", "value")

	if wrappedErr == nil {
		t.Error("expected wrapped error")
	}

	serviceErr, ok := wrappedErr.(*ServiceError)
	if !ok {
		t.Error("expected ServiceError type")
	}
	if serviceErr.Message != "formatted message: value" {
		t.Errorf("expected formatted message, got '%s'", serviceErr.Message)
	}
}

func TestIsServiceError(t *testing.T) {
	t.Run("is ServiceError", func(t *testing.T) {
		serviceErr := NewServiceError("test-service", "test-operation", "test message", nil)
		if !IsServiceError(serviceErr) {
			t.Error("expected IsServiceError to return true")
		}
	})

	t.Run("is not ServiceError", func(t *testing.T) {
		stdErr := errors.New("standard error")
		if IsServiceError(stdErr) {
			t.Error("expected IsServiceError to return false")
		}
	})
}

func TestGetServiceError(t *testing.T) {
	t.Run("get ServiceError", func(t *testing.T) {
		serviceErr := NewServiceError("test-service", "test-operation", "test message", nil)
		retrieved := GetServiceError(serviceErr)
		if retrieved != serviceErr {
			t.Errorf("expected retrieved ServiceError to be same instance")
		}
	})

	t.Run("get nil for standard error", func(t *testing.T) {
		stdErr := errors.New("standard error")
		retrieved := GetServiceError(stdErr)
		if retrieved != nil {
			t.Errorf("expected nil, got %v", retrieved)
		}
	})
}

func TestGetService(t *testing.T) {
	t.Run("get service from ServiceError", func(t *testing.T) {
		serviceErr := NewServiceError("test-service", "test-operation", "test message", nil)
		service := GetService(serviceErr)
		if service != "test-service" {
			t.Errorf("expected 'test-service', got '%s'", service)
		}
	})

	t.Run("get empty string from standard error", func(t *testing.T) {
		stdErr := errors.New("standard error")
		service := GetService(stdErr)
		if service != "" {
			t.Errorf("expected empty string, got '%s'", service)
		}
	})
}

func TestGetOperation(t *testing.T) {
	t.Run("get operation from ServiceError", func(t *testing.T) {
		serviceErr := NewServiceError("test-service", "test-operation", "test message", nil)
		operation := GetOperation(serviceErr)
		if operation != "test-operation" {
			t.Errorf("expected 'test-operation', got '%s'", operation)
		}
	})

	t.Run("get empty string from standard error", func(t *testing.T) {
		stdErr := errors.New("standard error")
		operation := GetOperation(stdErr)
		if operation != "" {
			t.Errorf("expected empty string, got '%s'", operation)
		}
	})
}

