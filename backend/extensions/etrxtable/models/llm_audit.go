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

// LLMModelAuditLog LLM 模型审计日志 - 从 Ent Schema 迁移到 GORM Model
type LLMModelAuditLog struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Action     string `gorm:"type:varchar(50);not null;index" json:"action"` // create, update, delete, test, use
	EntityType string `gorm:"type:varchar(50);not null" json:"entity_type"` // provider, model, instance, credential
	EntityID   uint   `gorm:"not null;index" json:"entity_id"`
	Changes    string `gorm:"type:jsonb" json:"changes,omitempty"` // 变更内容
	Result     string `gorm:"type:varchar(50)" json:"result,omitempty"` // success, failed
	Error      string `gorm:"type:text" json:"error,omitempty"`
	IPAddress  string `gorm:"type:varchar(50)" json:"ip_address,omitempty"`
	UserAgent  string `gorm:"type:varchar(255)" json:"user_agent,omitempty"`
	
	// 外键
	OperatorID  uint `gorm:"not null;index" json:"operator_id"`
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Operator  *User      `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}

// TableName 指定表名
func (LLMModelAuditLog) TableName() string {
	return "llm_model_audit_logs"
}

// LLMModelInstance LLM 模型实例
type LLMModelInstance struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Config      string `gorm:"type:jsonb" json:"config,omitempty"` // 实例配置（温度、max_tokens等）
	IsDefault   bool   `gorm:"default:false" json:"is_default"`
	IsEnabled   bool   `gorm:"default:true" json:"is_enabled"`
	
	// 限流配置
	RateLimitConfig string `gorm:"type:jsonb" json:"rate_limit_config,omitempty"`
	
	// 外键
	ModelID     uint `gorm:"not null;index" json:"model_id"`
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	CreatorID   uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Model     *LLMModel  `gorm:"foreignKey:ModelID" json:"model,omitempty"`
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (LLMModelInstance) TableName() string {
	return "llm_model_instances"
}

