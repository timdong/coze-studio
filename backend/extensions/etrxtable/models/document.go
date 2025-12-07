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

// Document 文档模型 - 从 Ent Schema 迁移到 GORM Model（用于 RAG）
type Document struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Content     string `gorm:"type:text" json:"content,omitempty"`
	Type        string `gorm:"type:varchar(50);default:'text'" json:"type"` // text, pdf, word, markdown
	FileURI     string `gorm:"type:varchar(500)" json:"file_uri,omitempty"` // 文件存储路径
	FileURL     string `gorm:"type:varchar(500)" json:"file_url,omitempty"` // 文件访问地址
	FileSize    int64  `gorm:"default:0" json:"file_size"`
	Status      string `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, processing, completed, failed
	
	// 处理信息
	ProcessedAt *time.Time `gorm:"" json:"processed_at,omitempty"`
	ChunkCount  int        `gorm:"default:0" json:"chunk_count"`
	
	// 质量评分
	QualityScore float64 `gorm:"default:0" json:"quality_score"`
	
	// 元数据
	Metadata string `gorm:"type:jsonb" json:"metadata,omitempty"`
	
	// 外键
	WorkspaceID uint  `gorm:"not null;index" json:"workspace_id"`
	CreatorID   uint  `gorm:"not null;index" json:"creator_id"`
	CategoryID  *uint `gorm:"index" json:"category_id,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Chunks    []Chunk    `gorm:"foreignKey:DocumentID" json:"chunks,omitempty"`
}

// TableName 指定表名
func (Document) TableName() string {
	return "documents"
}

// Chunk 文档分块模型
type Chunk struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Content    string `gorm:"type:text;not null" json:"content"`
	Position   int    `gorm:"not null" json:"position"`
	TokenCount int    `gorm:"default:0" json:"token_count"`
	
	// 外键
	DocumentID uint `gorm:"not null;index" json:"document_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Document  *Document  `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
	Embedding *Embedding `gorm:"foreignKey:ChunkID" json:"embedding,omitempty"`
}

// TableName 指定表名
func (Chunk) TableName() string {
	return "chunks"
}

// Embedding 向量嵌入模型
type Embedding struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Vector    string `gorm:"type:vector(1536)" json:"vector,omitempty"` // pgvector 类型
	Model     string `gorm:"type:varchar(100)" json:"model,omitempty"`
	
	// 外键
	ChunkID uint `gorm:"not null;uniqueIndex" json:"chunk_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Chunk *Chunk `gorm:"foreignKey:ChunkID" json:"chunk,omitempty"`
}

// TableName 指定表名
func (Embedding) TableName() string {
	return "embeddings"
}

