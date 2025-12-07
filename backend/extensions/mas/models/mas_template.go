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

// MASTemplate MAS 模板模型 - 从 Ent Schema 迁移到 GORM Model
type MASTemplate struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // agent, workflow, session
	Config      string `gorm:"type:jsonb;not null" json:"config"`
	Tags        string `gorm:"type:jsonb;default:'[]'" json:"tags,omitempty"`
	
	// 使用统计
	UseCount int `gorm:"default:0" json:"use_count"`
	
	// 状态
	IsPublic   bool   `gorm:"default:false" json:"is_public"`
	IsOfficial bool   `gorm:"default:false" json:"is_official"`
	Status     string `gorm:"type:varchar(50);default:'draft'" json:"status"`
	
	// 外键
	CreatorID uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Creator *User `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (MASTemplate) TableName() string {
	return "mas_templates"
}

