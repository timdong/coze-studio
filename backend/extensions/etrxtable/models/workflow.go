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

// Workflow 工作流模型 - 从 Ent Schema 迁移到 GORM Model
type Workflow struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"type:varchar(255);not null" json:"name"`
	Description  string `gorm:"type:text" json:"description,omitempty"`
	CanvasSchema string `gorm:"type:jsonb" json:"canvas_schema,omitempty"` // 画布架构定义
	InputParams  string `gorm:"type:jsonb" json:"input_params,omitempty"`  // 输入参数定义
	OutputParams string `gorm:"type:jsonb" json:"output_params,omitempty"` // 输出参数定义
	Version      string `gorm:"type:varchar(50);default:'1.0.0'" json:"version"`
	Status       string `gorm:"type:varchar(50);default:'draft'" json:"status"` // draft, published, archived
	
	// 外键
	CreatorID   uint `gorm:"not null;index" json:"creator_id"`
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	
	// 时间戳
	CreatedAt   time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	PublishedAt *time.Time     `gorm:"" json:"published_at,omitempty"`
	
	// 关系
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}

// TableName 指定表名
func (Workflow) TableName() string {
	return "workflows"
}

// WorkflowExecution 工作流执行记录
type WorkflowExecution struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	WorkflowID uint   `gorm:"not null;index" json:"workflow_id"`
	Status     string `gorm:"type:varchar(50);default:'running'" json:"status"` // running, completed, failed, cancelled
	Input      string `gorm:"type:jsonb" json:"input,omitempty"`
	Output     string `gorm:"type:jsonb" json:"output,omitempty"`
	Error      string `gorm:"type:text" json:"error,omitempty"`
	
	// 执行信息
	StartedAt   time.Time  `gorm:"" json:"started_at"`
	CompletedAt *time.Time `gorm:"" json:"completed_at,omitempty"`
	Duration    int        `gorm:"default:0" json:"duration"` // 毫秒
	
	// 外键
	ExecutorID  uint `gorm:"not null;index" json:"executor_id"`
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Workflow  *Workflow  `gorm:"foreignKey:WorkflowID" json:"workflow,omitempty"`
	Executor  *User      `gorm:"foreignKey:ExecutorID" json:"executor,omitempty"`
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}

// TableName 指定表名
func (WorkflowExecution) TableName() string {
	return "workflow_executions"
}

