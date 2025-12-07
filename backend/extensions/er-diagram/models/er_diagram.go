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

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ERDiagram ER 图模型 - 从 Ent Schema 迁移到 GORM Model
type ERDiagram struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"type:varchar(100);not null" json:"name"`
	Description    string `gorm:"type:varchar(500)" json:"description,omitempty"`
	DiagramData    datatypes.JSON `gorm:"type:jsonb;not null" json:"diagram_data"` // ER 图数据（DrawDB 格式）
	DatabaseType   string         `gorm:"type:varchar(50);default:'postgresql'" json:"database_type"`
	LinkedTableIDs datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"linked_table_ids"` // 关联的 ETRX-Table ID 列表
	SyncConfig     datatypes.JSON `gorm:"type:jsonb" json:"sync_config,omitempty"`         // 同步策略配置
	
	// 外键
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	CreatedBy   uint `gorm:"not null;index" json:"created_by"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName 指定表名
func (ERDiagram) TableName() string {
	return "er_diagrams"
}

// User 用户模型（简化版，用于关联）
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type:varchar(50)" json:"username"`
	Email    string `gorm:"type:varchar(100)" json:"email"`
}

func (User) TableName() string {
	return "users"
}

// Workspace 工作空间模型（简化版，用于关联）
type Workspace struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(100)" json:"name"`
}

func (Workspace) TableName() string {
	return "workspaces"
}

