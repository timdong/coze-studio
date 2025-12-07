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

// AgentCollaboration Agent 协作模型 - 从 Ent Schema 迁移到 GORM Model
type AgentCollaboration struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // sequential, parallel, conditional
	Status      string `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, active, completed, failed
	Config      string `gorm:"type:jsonb" json:"config,omitempty"`
	Result      string `gorm:"type:jsonb" json:"result,omitempty"`
	
	// 参与的 Agent
	AgentIDs    string `gorm:"type:jsonb;not null" json:"agent_ids"` // Agent ID 列表
	
	// 执行信息
	StartedAt   *time.Time `gorm:"" json:"started_at,omitempty"`
	CompletedAt *time.Time `gorm:"" json:"completed_at,omitempty"`
	Duration    int        `gorm:"default:0" json:"duration"`
	
	// 外键
	SessionID   uint `gorm:"not null;index" json:"session_id"`
	InitiatorID uint `gorm:"not null;index" json:"initiator_id"` // 发起者 Agent ID
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Session   *MASSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
	Initiator *Agent      `gorm:"foreignKey:InitiatorID" json:"initiator,omitempty"`
}

// TableName 指定表名
func (AgentCollaboration) TableName() string {
	return "agent_collaborations"
}

// AgentKnowledge Agent 知识库模型
type AgentKnowledge struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Content     string `gorm:"type:text;not null" json:"content"`
	Type        string `gorm:"type:varchar(50);default:'text'" json:"type"` // text, document, code, url
	Source      string `gorm:"type:varchar(500)" json:"source,omitempty"`
	Tags        string `gorm:"type:jsonb;default:'[]'" json:"tags,omitempty"`
	
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
func (AgentKnowledge) TableName() string {
	return "agent_knowledges"
}

