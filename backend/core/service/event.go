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

	"github.com/coze-dev/coze-studio/backend/core/event"
	ctxutil "github.com/coze-dev/coze-studio/backend/core/context"
)

// EventPublisher 事件发布接口
type EventPublisher interface {
	PublishEvent(ctx context.Context, eventType event.EventType, data interface{}, metadata map[string]interface{}) error
}

// EventBusProvider 事件总线提供者接口
type EventBusProvider interface {
	GetEventBus() *event.EventBus
}

// PublishEvent 发布事件
// 从context中提取用户ID和工作空间ID，并发布事件
func (s *BaseServiceImpl) PublishEvent(ctx context.Context, eventType event.EventType, data interface{}, metadata map[string]interface{}) error {
	// 从context中提取用户ID和工作空间ID
	userID, _ := ctxutil.GetUserID(ctx)
	username, _ := ctxutil.GetUsername(ctx)
	workspaceID, _ := ctxutil.GetWorkspaceID(ctx)

	// 创建事件
	e := event.NewEvent(eventType, s.getServiceName(), data).
		WithUser(userID, username).
		WithWorkspace(workspaceID).
		WithMetadataMap(metadata)

	// 通过Manager获取EventBus并发布
	if provider, ok := s.manager.(EventBusProvider); ok {
		eventBus := provider.GetEventBus()
		if eventBus != nil {
			return eventBus.Publish(ctx, e)
		}
	}

	return nil // 如果没有EventBus，静默失败
}

// PublishEventAsync 异步发布事件
func (s *BaseServiceImpl) PublishEventAsync(ctx context.Context, eventType event.EventType, data interface{}, metadata map[string]interface{}) {
	go func() {
		_ = s.PublishEvent(ctx, eventType, data, metadata)
	}()
}

// PublishResourceCreated 发布资源创建事件
func (s *BaseServiceImpl) PublishResourceCreated(ctx context.Context, resourceType string, resourceID uint, data interface{}) error {
	metadata := map[string]interface{}{
		"resource_type": resourceType,
		"resource_id":   resourceID,
	}
	return s.PublishEvent(ctx, event.EventResourceCreated, data, metadata)
}

// PublishResourceUpdated 发布资源更新事件
func (s *BaseServiceImpl) PublishResourceUpdated(ctx context.Context, resourceType string, resourceID uint, data interface{}) error {
	metadata := map[string]interface{}{
		"resource_type": resourceType,
		"resource_id":   resourceID,
	}
	return s.PublishEvent(ctx, event.EventResourceUpdated, data, metadata)
}

// PublishResourceDeleted 发布资源删除事件
func (s *BaseServiceImpl) PublishResourceDeleted(ctx context.Context, resourceType string, resourceID uint, data interface{}) error {
	metadata := map[string]interface{}{
		"resource_type": resourceType,
		"resource_id":   resourceID,
	}
	return s.PublishEvent(ctx, event.EventResourceDeleted, data, metadata)
}

// getServiceName 获取服务名称
// 通过反射获取服务类型名称
func (s *BaseServiceImpl) getServiceName() string {
	// 简化实现，返回通用服务名
	// 实际可以通过反射获取
	return "base-service"
}

