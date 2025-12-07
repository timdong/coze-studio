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

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// BeforeCreate GORM 钩子：在创建前验证 JSONB 字段
func (s *MASSession) BeforeCreate(tx *gorm.DB) error {
	return s.validateJSONBFields()
}

// BeforeUpdate GORM 钩子：在更新前验证 JSONB 字段
func (s *MASSession) BeforeUpdate(tx *gorm.DB) error {
	return s.validateJSONBFields()
}

// validateJSONBFields 验证并规范化所有 JSONB 字段
func (s *MASSession) validateJSONBFields() error {
	if len(s.Config) == 0 {
		s.Config = []byte("{}")
	}
	if len(s.Context) == 0 {
		s.Context = []byte("{}")
	}
	if len(s.AgentIDs) == 0 {
		s.AgentIDs = []byte("[]")
	}
	return nil
}

// MASSession MAS 会话模型 - 从 Ent Schema 迁移到 GORM Model
type MASSession struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // conversation, collaboration, workflow
	Status      string `gorm:"type:varchar(50);default:'active'" json:"status"` // active, completed, failed, cancelled
	
	// 配置
	Config datatypes.JSON `gorm:"type:jsonb" json:"config,omitempty"`
	
	// 上下文
	Context datatypes.JSON `gorm:"type:jsonb" json:"context,omitempty"`
	
	// 参与的 Agent
	AgentIDs datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"agent_ids,omitempty"` // Agent ID 列表
	
	// 执行信息
	StartedAt   *time.Time `gorm:"" json:"started_at,omitempty"`
	CompletedAt *time.Time `gorm:"" json:"completed_at,omitempty"`
	Duration    int        `gorm:"default:0" json:"duration"` // 毫秒
	
	// 外键
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	CreatorID   uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Tasks     []MASTask  `gorm:"foreignKey:SessionID" json:"tasks,omitempty"`
}

// TableName 指定表名
func (MASSession) TableName() string {
	return "mas_sessions"
}

