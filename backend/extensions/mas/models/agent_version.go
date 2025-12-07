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

// AgentVersion Agent 版本模型 - 从 Ent Schema 迁移到 GORM Model
type AgentVersion struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	VersionNumber  string `gorm:"type:varchar(50);not null" json:"version_number"`
	VersionName    string `gorm:"type:varchar(100)" json:"version_name,omitempty"`
	Description    string `gorm:"type:text" json:"description,omitempty"`
	Config         string `gorm:"type:jsonb" json:"config,omitempty"`
	SystemPrompt   string `gorm:"type:text" json:"system_prompt,omitempty"`
	UserPrompt     string `gorm:"type:text" json:"user_prompt,omitempty"`
	Status         string `gorm:"type:varchar(50);default:'draft'" json:"status"` // draft, published, deprecated
	IsActive       bool   `gorm:"default:false" json:"is_active"`
	ChangeLog      string `gorm:"type:text" json:"change_log,omitempty"`
	
	// 外键
	AgentID   uint `gorm:"not null;index" json:"agent_id"`
	CreatorID uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt   time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	PublishedAt *time.Time     `gorm:"" json:"published_at,omitempty"`
	
	// 关系
	Agent   *Agent `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Creator *User  `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (AgentVersion) TableName() string {
	return "agent_versions"
}

