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

// TableBody 表格主体 - 从 Ent Schema 迁移到 GORM Model
type TableBody struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:varchar(500)" json:"description,omitempty"`
	Icon        string         `gorm:"type:varchar(100)" json:"icon,omitempty"`
	Type        string         `gorm:"type:varchar(20);default:'grid'" json:"type"` // grid, kanban, gantt, calendar, form
	Settings    string         `gorm:"type:jsonb;default:'{}'" json:"settings,omitempty"`        // JSON 字段
	IsPublic    bool           `gorm:"default:false" json:"is_public"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active, archived, deleted
	TableType   string         `gorm:"type:varchar(20);default:'single'" json:"table_type"` // single, master, detail, hierarchical, network, temporal, spatial, document, custom
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Metadata    string         `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"` // JSON 字段
	
	// 外键
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	OwnerID     uint `gorm:"not null;index" json:"owner_id"`
	ERDiagramID *uint `gorm:"index" json:"er_diagram_id,omitempty"`
	
	// 同步配置
	SyncConfig string `gorm:"type:jsonb;default:'{\"auto_sync_enabled\":false,\"sync_mode\":\"manual\",\"sync_on_update\":false}'" json:"sync_config,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系（使用 Preload 加载）
	Workspace  *Workspace     `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator    *User          `gorm:"foreignKey:OwnerID" json:"creator,omitempty"`
	Columns    []TableColumn  `gorm:"foreignKey:TableBodyID" json:"columns,omitempty"`
	Rows       []TableRow     `gorm:"foreignKey:TableBodyID" json:"rows,omitempty"`
	Views      []TableView    `gorm:"foreignKey:TableBodyID" json:"views,omitempty"`
	Sorts      []TableSort    `gorm:"foreignKey:TableBodyID" json:"sorts,omitempty"`
}

// TableName 指定表名
func (TableBody) TableName() string {
	return "table_bodies"
}

// BeforeCreate GORM 钩子 - 创建前
func (t *TableBody) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// AfterCreate GORM 钩子 - 创建后
func (t *TableBody) AfterCreate(tx *gorm.DB) error {
	// 可以在这里添加创建后的逻辑
	return nil
}

