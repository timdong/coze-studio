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

// Conversation 对话模型 - 从 Ent Schema 迁移到 GORM Model
type Conversation struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Type        string `gorm:"type:varchar(50);default:'chat'" json:"type"` // chat, qa, workflow
	Status      string `gorm:"type:varchar(50);default:'active'" json:"status"` // active, archived, deleted
	
	// 消息统计
	MessageCount int `gorm:"default:0" json:"message_count"`
	TokenCount   int `gorm:"default:0" json:"token_count"`
	
	// 上下文
	Context  string `gorm:"type:jsonb" json:"context,omitempty"`
	Metadata string `gorm:"type:jsonb" json:"metadata,omitempty"`
	
	// 外键
	UserID      uint  `gorm:"not null;index" json:"user_id"`
	WorkspaceID uint  `gorm:"not null;index" json:"workspace_id"`
	AgentID     *uint `gorm:"index" json:"agent_id,omitempty"`
	
	// 时间戳
	CreatedAt      time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	LastMessageAt  *time.Time     `gorm:"" json:"last_message_at,omitempty"`
	
	// 关系
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	// Agent     *Agent     `gorm:"foreignKey:AgentID" json:"agent,omitempty"` // TODO: Agent 在 mas 包中
}

// TableName 指定表名
func (Conversation) TableName() string {
	return "conversations"
}

