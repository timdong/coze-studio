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

// TableColumn 表格列 - 从 Ent Schema 迁移到 GORM Model
type TableColumn struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Key         string `gorm:"type:varchar(100);not null" json:"key"`
	Label       string `gorm:"type:varchar(100)" json:"label,omitempty"`
	Type        string `gorm:"type:varchar(20);default:'text'" json:"type"` // text, number, date, datetime, boolean, select, multiselect, email, url, phone, file, image, json, formula
	Required    bool   `gorm:"default:false" json:"required"`
	Unique      bool   `gorm:"default:false" json:"unique"`
	Searchable  bool   `gorm:"default:true" json:"searchable"`
	Sortable    bool   `gorm:"default:true" json:"sortable"`
	Visible     bool   `gorm:"default:true" json:"visible"`
	Width       *int   `gorm:"" json:"width,omitempty"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	
	// JSON 字段
	Options      string `gorm:"type:jsonb" json:"options,omitempty"`      // 列选项配置
	Validation   string `gorm:"type:jsonb" json:"validation,omitempty"`   // 验证规则
	DefaultValue string `gorm:"type:varchar(500)" json:"default_value,omitempty"`
	Formula      string `gorm:"type:varchar(1000)" json:"formula,omitempty"`
	Description  string `gorm:"type:varchar(500)" json:"description,omitempty"`
	
	// 外键
	TableBodyID uint `gorm:"not null;index" json:"table_body_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	TableBody *TableBody `gorm:"foreignKey:TableBodyID" json:"table_body,omitempty"`
}

// TableName 指定表名
func (TableColumn) TableName() string {
	return "table_columns"
}

