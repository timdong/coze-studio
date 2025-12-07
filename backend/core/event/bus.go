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

package event

import (
	"context"
	"fmt"
	"sync"
)

// EventHandler 事件处理函数
type EventHandler func(ctx context.Context, event *Event) error

// EventBus 事件总线
type EventBus struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Unsubscribe 取消订阅事件（移除所有该类型的处理器）
func (eb *EventBus) Unsubscribe(eventType EventType) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	delete(eb.handlers, eventType)
}

// Publish 发布事件
// 如果任何处理器返回错误，会收集所有错误并返回
func (eb *EventBus) Publish(ctx context.Context, event *Event) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}

	eb.mu.RLock()
	handlers, exists := eb.handlers[event.Type]
	eb.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		return nil // 没有订阅者，不算错误
	}

	var errs []error
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("some event handlers failed: %v", errs)
	}

	return nil
}

// PublishAsync 异步发布事件
// 不等待处理器完成，不返回错误
func (eb *EventBus) PublishAsync(ctx context.Context, event *Event) {
	go func() {
		_ = eb.Publish(ctx, event)
	}()
}

// GetHandlerCount 获取指定事件类型的处理器数量
func (eb *EventBus) GetHandlerCount(eventType EventType) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.handlers[eventType])
}

// ListSubscribedEvents 列出所有已订阅的事件类型
func (eb *EventBus) ListSubscribedEvents() []EventType {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	events := make([]EventType, 0, len(eb.handlers))
	for eventType := range eb.handlers {
		events = append(events, eventType)
	}
	return events
}

