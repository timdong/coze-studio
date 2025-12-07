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

// Menu 菜单模型 - 从 Ent Schema 迁移到 GORM Model
type Menu struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Title       string `gorm:"type:varchar(100);not null" json:"title"`
	Icon        string `gorm:"type:varchar(50)" json:"icon,omitempty"`
	Path        string `gorm:"type:varchar(255);not null" json:"path"`
	Component   string `gorm:"type:varchar(255)" json:"component,omitempty"`
	Permission  string `gorm:"type:varchar(100)" json:"permission,omitempty"`
	Type        string `gorm:"type:varchar(20);default:'menu'" json:"type"` // menu, button, tab
	ParentID    *uint  `gorm:"index" json:"parent_id,omitempty"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	IsVisible   bool   `gorm:"default:true" json:"is_visible"`
	IsEnabled   bool   `gorm:"default:true" json:"is_enabled"`
	
	// 元数据
	Metadata    string `gorm:"type:jsonb" json:"metadata,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Parent   *Menu  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Menu `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Roles    []Role `gorm:"many2many:role_menus;" json:"roles,omitempty"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

