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

// Workspace 工作空间模型 - 从 Ent Schema 迁移到 GORM Model
type Workspace struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:varchar(500)" json:"description,omitempty"`
	Icon        string `gorm:"type:varchar(100)" json:"icon,omitempty"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"`
	IsActive    bool   `gorm:"default:true;index" json:"is_active"`
	Settings    string `gorm:"type:jsonb" json:"settings,omitempty"`
	
	// 外键
	OwnerID uint `gorm:"not null;index" json:"owner_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Owner  *User        `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Tables []TableBody  `gorm:"foreignKey:WorkspaceID" json:"tables,omitempty"`
}

// TableName 指定表名
func (Workspace) TableName() string {
	return "workspaces"
}

