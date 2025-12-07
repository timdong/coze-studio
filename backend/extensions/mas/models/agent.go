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
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// BeforeCreate GORM 钩子：在创建前验证 JSONB 字段
func (a *Agent) BeforeCreate(tx *gorm.DB) error {
	return a.validateJSONBFields()
}

// BeforeUpdate GORM 钩子：在更新前验证 JSONB 字段
func (a *Agent) BeforeUpdate(tx *gorm.DB) error {
	return a.validateJSONBFields()
}

// validateJSONBFields 验证并规范化所有 JSONB 字段
func (a *Agent) validateJSONBFields() error {
	// 使用 datatypes.JSON，GORM 会自动处理 JSON 序列化
	// 如果字段为空，设置为默认值
	if len(a.Config) == 0 {
		a.Config = []byte("{}")
	}
	if len(a.Capabilities) == 0 {
		a.Capabilities = []byte("[]")
	}
	if len(a.Constraints) == 0 {
		a.Constraints = []byte("[]")
	}
	if len(a.Tags) == 0 {
		a.Tags = []byte("[]")
	}
	// WorkflowConfig 和 MemoryConfig 是可选的，如果为空则保持为空
	return nil
}

// validateJSONBString 验证并规范化 JSON 字符串（用于模型层）
func validateJSONBString(jsonStr string, defaultValue string) string {
	if jsonStr == "" {
		return defaultValue
	}
	
	// 尝试解析 JSON
	var js interface{}
	if err := json.Unmarshal([]byte(jsonStr), &js); err != nil {
		return defaultValue
	}
	
	// 重新序列化
	validated, err := json.Marshal(js)
	if err != nil {
		return defaultValue
	}
	
	return string(validated)
}

// Agent MAS Agent 模型 - 从 Ent Schema 迁移到 GORM Model
type Agent struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:varchar(1000)" json:"description,omitempty"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // general, expert, coordinator, tool
	Status      string `gorm:"type:varchar(50);default:'inactive'" json:"status"` // inactive, active, busy, error
	
	// 图标
	IconURI string `gorm:"type:varchar(500)" json:"icon_uri,omitempty"` // 存储路径
	IconURL string `gorm:"type:varchar(500)" json:"icon_url,omitempty"` // 访问地址
	
	// 配置
	Config       datatypes.JSON `gorm:"type:jsonb" json:"config,omitempty"`       // LLM配置、工具配置等
	Capabilities datatypes.JSON `gorm:"type:jsonb" json:"capabilities,omitempty"` // 能力列表
	Constraints  datatypes.JSON `gorm:"type:jsonb" json:"constraints,omitempty"`  // 约束条件
	
	// 工作流相关
	WorkflowID     *int          `gorm:"" json:"workflow_id,omitempty"`
	WorkflowConfig datatypes.JSON `gorm:"type:jsonb" json:"workflow_config,omitempty"`
	
	// 提示词
	SystemPrompt string `gorm:"type:text" json:"system_prompt,omitempty"`
	UserPrompt   string `gorm:"type:text" json:"user_prompt,omitempty"`
	
	// 记忆相关
	MemoryEnabled bool          `gorm:"default:false" json:"memory_enabled"`
	MemoryConfig  datatypes.JSON `gorm:"type:jsonb" json:"memory_config,omitempty"`
	
	// 性能指标
	TotalTasks      int     `gorm:"default:0" json:"total_tasks"`
	CompletedTasks  int     `gorm:"default:0" json:"completed_tasks"`
	FailedTasks     int     `gorm:"default:0" json:"failed_tasks"`
	AvgResponseTime float64 `gorm:"default:0" json:"avg_response_time"`
	SuccessRate     float64 `gorm:"default:0" json:"success_rate"`
	LastActiveAt    *time.Time `gorm:"" json:"last_active_at,omitempty"`
	
	// 版本管理
	Version       string `gorm:"type:varchar(50);default:'1.0.0'" json:"version"`
	LatestVersion string `gorm:"type:varchar(50)" json:"latest_version,omitempty"`
	
	// 标签
	Tags datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"tags,omitempty"` // 标签列表
	
	// 外键
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	CreatorID   uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (Agent) TableName() string {
	return "agents"
}

// User 用户模型（简化版）
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type:varchar(50)" json:"username"`
}

func (User) TableName() string {
	return "users"
}

// Workspace 工作空间模型（简化版）
type Workspace struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(100)" json:"name"`
}

func (Workspace) TableName() string {
	return "workspaces"
}

