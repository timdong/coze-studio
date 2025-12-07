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

// DynamicTableMetadata 动态表元数据模型 - 从 Ent Schema 迁移到 GORM Model
type DynamicTableMetadata struct {
	ID                uint   `gorm:"primaryKey" json:"id"`
	PhysicalTableName string `gorm:"type:varchar(255);not null;uniqueIndex" json:"physical_table_name"`
	TableType         string `gorm:"type:varchar(50);default:'user_defined'" json:"table_type"`
	Schema            string `gorm:"type:jsonb" json:"schema,omitempty"` // 表结构定义
	Indexes           string `gorm:"type:jsonb" json:"indexes,omitempty"` // 索引定义
	Constraints       string `gorm:"type:jsonb" json:"constraints,omitempty"` // 约束定义
	
	// 统计信息
	RowCount      int64      `gorm:"default:0" json:"row_count"`
	LastSyncAt    *time.Time `gorm:"" json:"last_sync_at,omitempty"`
	
	// 外键
	TableBodyID uint `gorm:"not null;uniqueIndex" json:"table_body_id"`
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	TableBody *TableBody `gorm:"foreignKey:TableBodyID" json:"table_body,omitempty"`
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}

// TableName 指定表名
func (DynamicTableMetadata) TableName() string {
	return "dynamic_table_metadata"
}

// ERTableSyncHistory ER 表同步历史模型
type ERTableSyncHistory struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	SyncType    string `gorm:"type:varchar(50);not null" json:"sync_type"` // to_database, to_table
	Status      string `gorm:"type:varchar(50);not null" json:"status"` // success, failed, partial
	ChangesApplied string `gorm:"type:jsonb" json:"changes_applied,omitempty"`
	ErrorMessage   string `gorm:"type:text" json:"error_message,omitempty"`
	Duration       int    `gorm:"default:0" json:"duration"` // 毫秒
	
	// 外键
	ERDiagramID uint  `gorm:"not null;index" json:"er_diagram_id"`
	TableBodyID *uint `gorm:"index" json:"table_body_id,omitempty"`
	CreatorID   uint  `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	// ERDiagram *ERDiagram `gorm:"foreignKey:ERDiagramID" json:"er_diagram,omitempty"` // TODO: ERDiagram 在 er-diagram 包中
	TableBody *TableBody `gorm:"foreignKey:TableBodyID" json:"table_body,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (ERTableSyncHistory) TableName() string {
	return "er_table_sync_histories"
}

// View 自定义视图模型
type View struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Type        string `gorm:"type:varchar(50);default:'grid'" json:"type"`
	Config      string `gorm:"type:jsonb" json:"config,omitempty"`
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
func (View) TableName() string {
	return "views"
}

