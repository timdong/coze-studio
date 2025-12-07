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

// TableView 表格视图 - 从 Ent Schema 迁移到 GORM Model
type TableView struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Type        string `gorm:"type:varchar(20);default:'grid'" json:"type"` // grid, kanban, gantt, calendar, form
	Config      string `gorm:"type:jsonb" json:"config,omitempty"`          // JSON 字段存储视图配置
	IsDefault   bool   `gorm:"default:false" json:"is_default"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	
	// 外键
	TableBodyID uint `gorm:"not null;index" json:"table_body_id"`
	CreatorID   uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	TableBody *TableBody `gorm:"foreignKey:TableBodyID" json:"table_body,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (TableView) TableName() string {
	return "table_views"
}

