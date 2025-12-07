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

// AgentTag Agent 标签模型 - 从 Ent Schema 迁移到 GORM Model
type AgentTag struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	Description string `gorm:"type:varchar(500)" json:"description,omitempty"`
	Color       string `gorm:"type:varchar(20);default:'#3b82f6'" json:"color"`
	Icon        string `gorm:"type:varchar(50)" json:"icon,omitempty"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	
	// 外键
	CreatorID uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Creator *User   `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Agents  []Agent `gorm:"many2many:agent_tag_relations;" json:"agents,omitempty"`
}

// TableName 指定表名
func (AgentTag) TableName() string {
	return "agent_tags"
}

