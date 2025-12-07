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

// TableRow 表格行 - 从 Ent Schema 迁移到 GORM Model
type TableRow struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Data        string `gorm:"type:jsonb;not null" json:"data"` // JSON 字段存储行数据
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	IsDeleted   bool   `gorm:"default:false" json:"is_deleted"`
	
	// 外键
	TableBodyID uint `gorm:"not null;index" json:"table_body_id"`
	CreatorID   uint `gorm:"not null;index" json:"creator_id"`
	UpdaterID   *uint `gorm:"index" json:"updater_id,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	TableBody *TableBody `gorm:"foreignKey:TableBodyID" json:"table_body,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Updater   *User      `gorm:"foreignKey:UpdaterID" json:"updater,omitempty"`
}

// TableName 指定表名
func (TableRow) TableName() string {
	return "table_rows"
}

