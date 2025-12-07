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

// TableLink 表格关联模型 - 从 Ent Schema 迁移到 GORM Model
type TableLink struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	Name             string `gorm:"type:varchar(100);not null" json:"name"`
	Description      string `gorm:"type:varchar(500)" json:"description,omitempty"`
	Type             string `gorm:"type:varchar(50);default:'one_to_many'" json:"type"` // one_to_one, one_to_many, many_to_many
	SourceColumnKey  string `gorm:"type:varchar(100);not null" json:"source_column_key"`
	TargetColumnKey  string `gorm:"type:varchar(100);not null" json:"target_column_key"`
	CascadeDelete    bool   `gorm:"default:false" json:"cascade_delete"`
	CascadeUpdate    bool   `gorm:"default:false" json:"cascade_update"`
	
	// 外键
	SourceTableID uint `gorm:"not null;index" json:"source_table_id"`
	TargetTableID uint `gorm:"not null;index" json:"target_table_id"`
	WorkspaceID   uint `gorm:"not null;index" json:"workspace_id"`
	CreatorID     uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	SourceTable *TableBody `gorm:"foreignKey:SourceTableID" json:"source_table,omitempty"`
	TargetTable *TableBody `gorm:"foreignKey:TargetTableID" json:"target_table,omitempty"`
	Workspace   *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator     *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (TableLink) TableName() string {
	return "table_links"
}

