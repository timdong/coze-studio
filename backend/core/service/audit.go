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

	"etrxlite/platform/services/audit"
	ctxutil "github.com/coze-dev/coze-studio/backend/core/context"
)

// AuditLogger 审计日志接口
type AuditLogger interface {
	LogAction(ctx context.Context, action, resource, resourceID string, details interface{}, result string) error
}

// AuditServiceProvider 审计服务提供者接口
type AuditServiceProvider interface {
	GetAuditService() *audit.Logger
}

// LogAction 记录操作审计日志
func (s *BaseServiceImpl) LogAction(ctx context.Context, action, resource, resourceID string, details interface{}, result string) error {
	userID, _ := ctxutil.GetUserID(ctx)
	username, _ := ctxutil.GetUsername(ctx)

	// 通过Manager获取AuditService
	if provider, ok := s.manager.(AuditServiceProvider); ok {
		auditService := provider.GetAuditService()
		if auditService != nil {
			return auditService.LogAction(ctx, userID, username, action, resource, resourceID, details, result)
		}
	}

	return nil // 如果没有AuditService，静默失败
}

// LogCreate 记录创建操作
func (s *BaseServiceImpl) LogCreate(ctx context.Context, resource, resourceID string, details interface{}) error {
	return s.LogAction(ctx, "create", resource, resourceID, details, "success")
}

// LogUpdate 记录更新操作
func (s *BaseServiceImpl) LogUpdate(ctx context.Context, resource, resourceID string, details interface{}) error {
	return s.LogAction(ctx, "update", resource, resourceID, details, "success")
}

// LogDelete 记录删除操作
func (s *BaseServiceImpl) LogDelete(ctx context.Context, resource, resourceID string, details interface{}) error {
	return s.LogAction(ctx, "delete", resource, resourceID, details, "success")
}

// LogQuery 记录查询操作
func (s *BaseServiceImpl) LogQuery(ctx context.Context, resource string, details interface{}) error {
	return s.LogAction(ctx, "query", resource, "", details, "success")
}

// LogActionWithError 记录操作审计日志（带错误信息）
func (s *BaseServiceImpl) LogActionWithError(ctx context.Context, action, resource, resourceID string, details interface{}, err error) error {
	result := "success"
	if err != nil {
		result = "failure"
		// 将错误信息添加到details中
		if detailsMap, ok := details.(map[string]interface{}); ok {
			detailsMap["error"] = err.Error()
			details = detailsMap
		} else {
			details = map[string]interface{}{
				"original_details": details,
				"error":            err.Error(),
			}
		}
	}
	return s.LogAction(ctx, action, resource, resourceID, details, result)
}

// LogSecurityEvent 记录安全事件
func (s *BaseServiceImpl) LogSecurityEvent(ctx context.Context, eventType, severity, description string, details interface{}) error {
	userID, _ := ctxutil.GetUserID(ctx)

	// 通过Manager获取AuditService
	if provider, ok := s.manager.(AuditServiceProvider); ok {
		auditService := provider.GetAuditService()
		if auditService != nil {
			return auditService.LogSecurityEvent(ctx, userID, eventType, severity, description, details)
		}
	}

	return nil
}

// LogPerformanceMetric 记录性能指标
func (s *BaseServiceImpl) LogPerformanceMetric(ctx context.Context, operation string, duration int64, memoryUsage int64, cpuUsage float64, details interface{}) error {
	serviceName := s.getServiceName()

	// 通过Manager获取AuditService
	if provider, ok := s.manager.(AuditServiceProvider); ok {
		auditService := provider.GetAuditService()
		if auditService != nil {
			return auditService.LogPerformanceMetric(ctx, serviceName, operation, duration, memoryUsage, cpuUsage, details)
		}
	}

	return nil
}

