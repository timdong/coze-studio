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

// AgentMemory Agent 记忆模型 - 从 Ent Schema 迁移到 GORM Model
type AgentMemory struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // short_term, long_term, episodic, semantic
	Content     string `gorm:"type:text;not null" json:"content"`
	Context     string `gorm:"type:jsonb" json:"context,omitempty"`
	Importance  float64 `gorm:"default:0" json:"importance"`  // 0-1
	AccessCount int     `gorm:"default:0" json:"access_count"`
	LastAccess  *time.Time `gorm:"" json:"last_access,omitempty"`
	ExpiresAt   *time.Time `gorm:"" json:"expires_at,omitempty"`
	
	// 外键
	AgentID     uint  `gorm:"not null;index" json:"agent_id"`
	SessionID   *uint `gorm:"index" json:"session_id,omitempty"`
	TaskID      *uint `gorm:"index" json:"task_id,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Agent   *Agent      `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Session *MASSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
	Task    *MASTask    `gorm:"foreignKey:TaskID" json:"task,omitempty"`
}

// TableName 指定表名
func (AgentMemory) TableName() string {
	return "agent_memories"
}

