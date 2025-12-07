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

// Plugin 插件模型 - 从 Ent Schema 迁移到 GORM Model
type Plugin struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	DisplayName string `gorm:"type:varchar(100)" json:"display_name,omitempty"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Version     string `gorm:"type:varchar(50);default:'1.0.0'" json:"version"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // http, grpc, python, nodejs
	Status      string `gorm:"type:varchar(50);default:'inactive'" json:"status"` // inactive, active, error
	
	// 配置
	Config   string `gorm:"type:jsonb" json:"config,omitempty"`
	Schema   string `gorm:"type:jsonb" json:"schema,omitempty"` // OpenAPI schema
	Manifest string `gorm:"type:jsonb" json:"manifest,omitempty"`
	
	// 图标和文档
	Icon    string `gorm:"type:varchar(500)" json:"icon,omitempty"`
	IconURL string `gorm:"type:varchar(500)" json:"icon_url,omitempty"`
	DocURL  string `gorm:"type:varchar(500)" json:"doc_url,omitempty"`
	
	// 分类和标签
	Category string `gorm:"type:varchar(50)" json:"category,omitempty"`
	Tags     string `gorm:"type:jsonb;default:'[]'" json:"tags,omitempty"`
	
	// 外键
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	CreatorID   uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (Plugin) TableName() string {
	return "plugins"
}

