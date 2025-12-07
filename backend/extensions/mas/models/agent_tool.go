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

package models

import (
	"time"

	"gorm.io/gorm"
)

// AgentTool Agent 工具模型 - 从 Ent Schema 迁移到 GORM Model
type AgentTool struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // function, plugin, api, workflow
	Config      string `gorm:"type:jsonb" json:"config,omitempty"`
	Schema      string `gorm:"type:jsonb" json:"schema,omitempty"` // OpenAPI schema
	IsEnabled   bool   `gorm:"default:true" json:"is_enabled"`
	
	// 外键
	AgentID uint `gorm:"not null;index" json:"agent_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Agent *Agent `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
}

// TableName 指定表名
func (AgentTool) TableName() string {
	return "agent_tools"
}

// AgentMessage Agent 消息模型
type AgentMessage struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Role      string `gorm:"type:varchar(50);not null" json:"role"` // user, assistant, system
	Content   string `gorm:"type:text;not null" json:"content"`
	Metadata  string `gorm:"type:jsonb" json:"metadata,omitempty"`
	
	// 外键
	AgentID   uint `gorm:"not null;index" json:"agent_id"`
	SessionID uint `gorm:"not null;index" json:"session_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Agent   *Agent      `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Session *MASSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}

// TableName 指定表名
func (AgentMessage) TableName() string {
	return "agent_messages"
}

// AgentExecutionLog Agent 执行日志模型
type AgentExecutionLog struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Action      string `gorm:"type:varchar(100);not null" json:"action"`
	Input       string `gorm:"type:jsonb" json:"input,omitempty"`
	Output      string `gorm:"type:jsonb" json:"output,omitempty"`
	Error       string `gorm:"type:text" json:"error,omitempty"`
	Status      string `gorm:"type:varchar(50);not null" json:"status"` // success, failed
	Duration    int    `gorm:"default:0" json:"duration"` // 毫秒
	
	// 外键
	AgentID   uint  `gorm:"not null;index" json:"agent_id"`
	TaskID    *uint `gorm:"index" json:"task_id,omitempty"`
	SessionID *uint `gorm:"index" json:"session_id,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Agent   *Agent      `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Task    *MASTask    `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Session *MASSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}

// TableName 指定表名
func (AgentExecutionLog) TableName() string {
	return "agent_execution_logs"
}

// AgentFavorite Agent 收藏模型
type AgentFavorite struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	
	// 外键
	AgentID uint `gorm:"not null;index" json:"agent_id"`
	UserID  uint `gorm:"not null;index" json:"user_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Agent *Agent `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (AgentFavorite) TableName() string {
	return "agent_favorites"
}

