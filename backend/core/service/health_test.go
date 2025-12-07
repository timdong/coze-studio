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
	"context"
	"errors"
	"testing"
)

// mockHealthChecker 模拟健康检查器
type mockHealthChecker struct {
	shouldFail bool
	name       string
}

func (m *mockHealthChecker) HealthCheck(ctx context.Context) error {
	if m.shouldFail {
		return errors.New("health check failed")
	}
	return nil
}

func TestCompositeHealthChecker_Register(t *testing.T) {
	checker := NewCompositeHealthChecker()
	mockChecker := &mockHealthChecker{name: "test"}

	checker.Register("test", mockChecker)

	if len(checker.checkers) != 1 {
		t.Errorf("expected 1 checker, got %d", len(checker.checkers))
	}
}

func TestCompositeHealthChecker_HealthCheck_Success(t *testing.T) {
	checker := NewCompositeHealthChecker()
	checker.Register("service1", &mockHealthChecker{shouldFail: false})
	checker.Register("service2", &mockHealthChecker{shouldFail: false})

	ctx := context.Background()
	err := checker.HealthCheck(ctx)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCompositeHealthChecker_HealthCheck_Failure(t *testing.T) {
	checker := NewCompositeHealthChecker()
	checker.Register("service1", &mockHealthChecker{shouldFail: false})
	checker.Register("service2", &mockHealthChecker{shouldFail: true})

	ctx := context.Background()
	err := checker.HealthCheck(ctx)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCompositeHealthChecker_HealthCheckDetailed(t *testing.T) {
	checker := NewCompositeHealthChecker()
	checker.Register("service1", &mockHealthChecker{shouldFail: false})
	checker.Register("service2", &mockHealthChecker{shouldFail: true})

	ctx := context.Background()
	results, err := checker.HealthCheckDetailed(ctx)

	if err == nil {
		t.Error("expected error, got nil")
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if results[0].Status != "healthy" {
		t.Errorf("expected service1 to be healthy, got %s", results[0].Status)
	}

	if results[1].Status != "unhealthy" {
		t.Errorf("expected service2 to be unhealthy, got %s", results[1].Status)
	}
}

