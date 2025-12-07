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

import "time"

// EventType 事件类型
type EventType string

const (
	// 资源事件
	EventResourceCreated EventType = "resource.created"
	EventResourceUpdated EventType = "resource.updated"
	EventResourceDeleted EventType = "resource.deleted"

	// 用户事件
	EventUserCreated    EventType = "user.created"
	EventUserUpdated    EventType = "user.updated"
	EventUserDeleted    EventType = "user.deleted"
	EventUserLoggedIn   EventType = "user.logged_in"
	EventUserLoggedOut  EventType = "user.logged_out"

	// 工作空间事件
	EventWorkspaceCreated EventType = "workspace.created"
	EventWorkspaceUpdated EventType = "workspace.updated"
	EventWorkspaceDeleted EventType = "workspace.deleted"

	// 文档事件
	EventDocumentCreated EventType = "document.created"
	EventDocumentUpdated EventType = "document.updated"
	EventDocumentDeleted EventType = "document.deleted"

	// 表格事件
	EventTableCreated EventType = "table.created"
	EventTableUpdated EventType = "table.updated"
	EventTableDeleted EventType = "table.deleted"

	// 工作流事件
	EventWorkflowCreated   EventType = "workflow.created"
	EventWorkflowUpdated   EventType = "workflow.updated"
	EventWorkflowDeleted   EventType = "workflow.deleted"
	EventWorkflowExecuted  EventType = "workflow.executed"
	EventWorkflowCompleted EventType = "workflow.completed"
	EventWorkflowFailed    EventType = "workflow.failed"
)

// Event 事件
type Event struct {
	ID          string                 `json:"id"`
	Type        EventType              `json:"type"`
	Source      string                 `json:"source"`
	Timestamp   time.Time              `json:"timestamp"`
	UserID      uint                   `json:"user_id,omitempty"`
	Username    string                 `json:"username,omitempty"`
	WorkspaceID uint                   `json:"workspace_id,omitempty"`
	Data        interface{}            `json:"data"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewEvent 创建新事件
func NewEvent(eventType EventType, source string, data interface{}) *Event {
	return &Event{
		ID:        generateEventID(),
		Type:      eventType,
		Source:    source,
		Timestamp: time.Now(),
		Data:      data,
		Metadata:  make(map[string]interface{}),
	}
}

// WithUser 设置用户信息
func (e *Event) WithUser(userID uint, username string) *Event {
	e.UserID = userID
	e.Username = username
	return e
}

// WithWorkspace 设置工作空间信息
func (e *Event) WithWorkspace(workspaceID uint) *Event {
	e.WorkspaceID = workspaceID
	return e
}

// WithMetadata 添加元数据
func (e *Event) WithMetadata(key string, value interface{}) *Event {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// WithMetadataMap 批量添加元数据
func (e *Event) WithMetadataMap(metadata map[string]interface{}) *Event {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	for k, v := range metadata {
		e.Metadata[k] = v
	}
	return e
}

// generateEventID 生成事件ID
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	// 使用时间戳的纳秒部分作为随机种子
	seed := uint64(time.Now().UnixNano())
	for i := range b {
		seed = seed*1103515245 + 12345 // 线性同余生成器
		b[i] = charset[seed%uint64(len(charset))]
	}
	return string(b)
}

