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

// DocumentCategory 文档分类模型 - 从 Ent Schema 迁移到 GORM Model
type DocumentCategory struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:varchar(500)" json:"description,omitempty"`
	Icon        string `gorm:"type:varchar(50)" json:"icon,omitempty"`
	Color       string `gorm:"type:varchar(20);default:'#3b82f6'" json:"color"`
	ParentID    *uint  `gorm:"index" json:"parent_id,omitempty"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	
	// 外键
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	CreatorID   uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Workspace *Workspace         `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator   *User              `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Parent    *DocumentCategory  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []DocumentCategory `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Documents []Document         `gorm:"foreignKey:CategoryID" json:"documents,omitempty"`
}

// TableName 指定表名
func (DocumentCategory) TableName() string {
	return "document_categories"
}

// DocumentQuality 文档质量模型
type DocumentQuality struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	OverallScore    float64 `gorm:"not null" json:"overall_score"` // 0-1
	ReadabilityScore float64 `gorm:"default:0" json:"readability_score"`
	CompletenessScore float64 `gorm:"default:0" json:"completeness_score"`
	StructureScore  float64 `gorm:"default:0" json:"structure_score"`
	
	// 问题和建议
	Issues       string `gorm:"type:jsonb" json:"issues,omitempty"`
	Suggestions  string `gorm:"type:jsonb" json:"suggestions,omitempty"`
	
	// 外键
	DocumentID uint `gorm:"not null;uniqueIndex" json:"document_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Document *Document `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
}

// TableName 指定表名
func (DocumentQuality) TableName() string {
	return "document_qualities"
}

// DocumentVersion 文档版本模型
type DocumentVersion struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	VersionNumber string `gorm:"type:varchar(50);not null" json:"version_number"`
	Content       string `gorm:"type:text" json:"content,omitempty"`
	ChangeLog     string `gorm:"type:text" json:"change_log,omitempty"`
	FileURI       string `gorm:"type:varchar(500)" json:"file_uri,omitempty"`
	FileSize      int64  `gorm:"default:0" json:"file_size"`
	
	// 外键
	DocumentID uint `gorm:"not null;index" json:"document_id"`
	CreatorID  uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Document *Document `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
	Creator  *User     `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (DocumentVersion) TableName() string {
	return "document_versions"
}

