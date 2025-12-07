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

// BeforeCreate GORM 钩子：在创建前验证 JSONB 字段
func (t *MASTask) BeforeCreate(tx *gorm.DB) error {
	return t.validateJSONBFields()
}

// BeforeUpdate GORM 钩子：在更新前验证 JSONB 字段
func (t *MASTask) BeforeUpdate(tx *gorm.DB) error {
	return t.validateJSONBFields()
}

// validateJSONBFields 验证并规范化所有 JSONB 字段
func (t *MASTask) validateJSONBFields() error {
	if len(t.Input) == 0 {
		t.Input = []byte("{}")
	}
	if len(t.Output) == 0 {
		t.Output = []byte("{}")
	}
	if len(t.Dependencies) == 0 {
		t.Dependencies = []byte("[]")
	}
	return nil
}

// MASTask MAS 任务模型 - 从 Ent Schema 迁移到 GORM Model
type MASTask struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // query, analysis, generation, execution
	Priority    int    `gorm:"default:0" json:"priority"`
	Status      string `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, running, completed, failed, cancelled
	
	// 输入输出
	Input  datatypes.JSON `gorm:"type:jsonb" json:"input,omitempty"`
	Output datatypes.JSON `gorm:"type:jsonb" json:"output,omitempty"`
	Error  string         `gorm:"type:text" json:"error,omitempty"`
	
	// 执行信息
	StartedAt   *time.Time `gorm:"" json:"started_at,omitempty"`
	CompletedAt *time.Time `gorm:"" json:"completed_at,omitempty"`
	Duration    int        `gorm:"default:0" json:"duration"` // 毫秒
	
	// 依赖关系
	Dependencies datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"dependencies,omitempty"` // 依赖的任务 ID 列表
	
	// 外键
	AgentID     uint  `gorm:"not null;index" json:"agent_id"`
	SessionID   *uint `gorm:"index" json:"session_id,omitempty"`
	WorkspaceID uint  `gorm:"not null;index" json:"workspace_id"`
	CreatorID   uint  `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Agent     *Agent     `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	// Session   *MASSession `gorm:"foreignKey:SessionID" json:"session,omitempty"` // 在 mas_session.go 中定义
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (MASTask) TableName() string {
	return "mas_tasks"
}

