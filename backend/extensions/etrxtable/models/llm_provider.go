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

// LLMProvider LLM 提供商模型 - 从 Ent Schema 迁移到 GORM Model
type LLMProvider struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	DisplayName string `gorm:"type:varchar(100)" json:"display_name,omitempty"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // openai, claude, ollama, qwen, etc.
	BaseURL     string `gorm:"type:varchar(500)" json:"base_url,omitempty"`
	Status      string `gorm:"type:varchar(50);default:'inactive'" json:"status"` // inactive, active, error
	IsEnabled   bool   `gorm:"default:false" json:"is_enabled"`
	
	// 配置
	Config     string `gorm:"type:jsonb" json:"config,omitempty"`
	Features   string `gorm:"type:jsonb;default:'[]'" json:"features,omitempty"` // chat, embedding, multimodal
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Credentials []LLMProviderCredential `gorm:"foreignKey:ProviderID" json:"credentials,omitempty"`
	Models      []LLMModel              `gorm:"foreignKey:ProviderID" json:"models,omitempty"`
}

// TableName 指定表名
func (LLMProvider) TableName() string {
	return "llm_providers"
}

// LLMProviderCredential LLM 提供商凭证模型
type LLMProviderCredential struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	APIKey      string `gorm:"type:varchar(500);not null" json:"-"` // 加密存储，不返回
	APISecret   string `gorm:"type:varchar(500)" json:"-"`
	IsDefault   bool   `gorm:"default:false" json:"is_default"`
	Status      string `gorm:"type:varchar(50);default:'inactive'" json:"status"`
	
	// 外键
	ProviderID  uint `gorm:"not null;index" json:"provider_id"`
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Provider  *LLMProvider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	Workspace *Workspace   `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}

// TableName 指定表名
func (LLMProviderCredential) TableName() string {
	return "llm_provider_credentials"
}

// LLMModel LLM 模型定义
type LLMModel struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ModelID     string `gorm:"type:varchar(100);not null;uniqueIndex" json:"model_id"`
	DisplayName string `gorm:"type:varchar(100)" json:"display_name,omitempty"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Type        string `gorm:"type:varchar(50)" json:"type"` // chat, embedding, multimodal
	Capabilities string `gorm:"type:jsonb;default:'[]'" json:"capabilities,omitempty"`
	MaxTokens   int    `gorm:"default:4096" json:"max_tokens"`
	IsEnabled   bool   `gorm:"default:true" json:"is_enabled"`
	
	// 外键
	ProviderID uint `gorm:"not null;index" json:"provider_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Provider *LLMProvider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
}

// TableName 指定表名
func (LLMModel) TableName() string {
	return "llm_models"
}

