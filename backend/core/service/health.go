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
	"fmt"
)

// HealthChecker 健康检查接口
type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}

// HealthCheckResult 健康检查结果
type HealthCheckResult struct {
	Service string `json:"service"`
	Status  string `json:"status"` // "healthy", "unhealthy"
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// CompositeHealthChecker 组合健康检查器
// 可以同时检查多个服务的健康状态
type CompositeHealthChecker struct {
	checkers map[string]HealthChecker
}

// NewCompositeHealthChecker 创建组合健康检查器
func NewCompositeHealthChecker() *CompositeHealthChecker {
	return &CompositeHealthChecker{
		checkers: make(map[string]HealthChecker),
	}
}

// Register 注册健康检查器
func (c *CompositeHealthChecker) Register(name string, checker HealthChecker) {
	c.checkers[name] = checker
}

// HealthCheck 执行所有健康检查
func (c *CompositeHealthChecker) HealthCheck(ctx context.Context) error {
	var errors []error

	for name, checker := range c.checkers {
		if err := checker.HealthCheck(ctx); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("health check failed: %v", errors)
	}

	return nil
}

// HealthCheckDetailed 执行详细健康检查，返回每个服务的状态
func (c *CompositeHealthChecker) HealthCheckDetailed(ctx context.Context) ([]HealthCheckResult, error) {
	results := make([]HealthCheckResult, 0, len(c.checkers))
	var hasError bool

	for name, checker := range c.checkers {
		result := HealthCheckResult{
			Service: name,
			Status:  "healthy",
		}

		if err := checker.HealthCheck(ctx); err != nil {
			result.Status = "unhealthy"
			result.Error = err.Error()
			result.Message = fmt.Sprintf("Health check failed: %v", err)
			hasError = true
		}

		results = append(results, result)
	}

	if hasError {
		return results, fmt.Errorf("some services are unhealthy")
	}

	return results, nil
}

