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

// WorkflowTemplate 工作流模板模型 - 从 Ent Schema 迁移到 GORM Model
type WorkflowTemplate struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"type:varchar(255);not null" json:"name"`
	Description  string `gorm:"type:text" json:"description,omitempty"`
	Icon         string `gorm:"type:varchar(100)" json:"icon,omitempty"`
	CategoryName string `gorm:"type:varchar(50)" json:"category_name,omitempty"`
	Tags         string `gorm:"type:jsonb;default:'[]'" json:"tags,omitempty"`
	CanvasSchema string `gorm:"type:jsonb;not null" json:"canvas_schema"` // 工作流定义
	InputParams  string `gorm:"type:jsonb" json:"input_params,omitempty"`
	OutputParams string `gorm:"type:jsonb" json:"output_params,omitempty"`
	
	// 使用统计
	UseCount     int     `gorm:"default:0" json:"use_count"`
	AverageRating float64 `gorm:"default:0" json:"average_rating"`
	RatingCount  int     `gorm:"default:0" json:"rating_count"`
	
	// 状态
	IsPublic     bool   `gorm:"default:false" json:"is_public"`
	IsOfficial   bool   `gorm:"default:false" json:"is_official"`
	Status       string `gorm:"type:varchar(50);default:'draft'" json:"status"` // draft, published, deprecated
	
	// 外键
	WorkspaceID uint  `gorm:"not null;index" json:"workspace_id"`
	CreatorID   uint  `gorm:"not null;index" json:"creator_id"`
	CategoryID  *uint `gorm:"index" json:"category_id,omitempty"`
	
	// 时间戳
	CreatedAt   time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	PublishedAt *time.Time     `gorm:"" json:"published_at,omitempty"`
	
	// 关系
	Workspace *Workspace        `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator   *User             `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Category  *TemplateCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// TableName 指定表名
func (WorkflowTemplate) TableName() string {
	return "workflow_templates"
}

// TemplateCategory 模板分类模型
type TemplateCategory struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description string `gorm:"type:varchar(500)" json:"description,omitempty"`
	Icon        string `gorm:"type:varchar(50)" json:"icon,omitempty"`
	ParentID    *uint  `gorm:"index" json:"parent_id,omitempty"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	
	// 外键
	CreatorID uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Creator   *User               `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Parent    *TemplateCategory   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []TemplateCategory  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Templates []WorkflowTemplate  `gorm:"foreignKey:CategoryID" json:"templates,omitempty"`
}

// TableName 指定表名
func (TemplateCategory) TableName() string {
	return "template_categories"
}

// TemplateRating 模板评分模型
type TemplateRating struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Rating   int    `gorm:"not null" json:"rating"` // 1-5
	Comment  string `gorm:"type:text" json:"comment,omitempty"`
	
	// 外键
	TemplateID uint `gorm:"not null;index" json:"template_id"`
	UserID     uint `gorm:"not null;index" json:"user_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Template *WorkflowTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	User     *User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (TemplateRating) TableName() string {
	return "template_ratings"
}

// WorkflowPermission 工作流权限模型
type WorkflowPermission struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Permission string `gorm:"type:varchar(50);not null" json:"permission"` // read, write, execute, manage
	
	// 外键
	WorkflowID uint `gorm:"not null;index" json:"workflow_id"`
	UserID     *uint `gorm:"index" json:"user_id,omitempty"`
	RoleID     *uint `gorm:"index" json:"role_id,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Workflow *Workflow `gorm:"foreignKey:WorkflowID" json:"workflow,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role     *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName 指定表名
func (WorkflowPermission) TableName() string {
	return "workflow_permissions"
}

// WorkflowVersion 工作流版本模型
type WorkflowVersion struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	VersionNumber string `gorm:"type:varchar(50);not null" json:"version_number"`
	ChangeLog     string `gorm:"type:text" json:"change_log,omitempty"`
	CanvasSchema  string `gorm:"type:jsonb;not null" json:"canvas_schema"`
	
	// 外键
	WorkflowID uint `gorm:"not null;index" json:"workflow_id"`
	CreatorID  uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Workflow *Workflow `gorm:"foreignKey:WorkflowID" json:"workflow,omitempty"`
	Creator  *User     `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (WorkflowVersion) TableName() string {
	return "workflow_versions"
}

// WorkflowExecutionLog 工作流执行日志模型
type WorkflowExecutionLog struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	NodeID      string `gorm:"type:varchar(100);not null" json:"node_id"`
	NodeType    string `gorm:"type:varchar(50)" json:"node_type,omitempty"`
	Status      string `gorm:"type:varchar(50);not null" json:"status"` // running, completed, failed
	Input       string `gorm:"type:jsonb" json:"input,omitempty"`
	Output      string `gorm:"type:jsonb" json:"output,omitempty"`
	Error       string `gorm:"type:text" json:"error,omitempty"`
	StartedAt   time.Time  `gorm:"not null" json:"started_at"`
	CompletedAt *time.Time `gorm:"" json:"completed_at,omitempty"`
	Duration    int        `gorm:"default:0" json:"duration"` // 毫秒
	
	// 外键
	ExecutionID uint `gorm:"not null;index" json:"execution_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Execution *WorkflowExecution `gorm:"foreignKey:ExecutionID" json:"execution,omitempty"`
}

// TableName 指定表名
func (WorkflowExecutionLog) TableName() string {
	return "workflow_execution_logs"
}

