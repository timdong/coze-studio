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

// DatabaseConnection 数据库连接模型 - 从 Ent Schema 迁移到 GORM Model
type DatabaseConnection struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:varchar(500)" json:"description,omitempty"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // postgresql, mysql, sqlite, mongodb
	Host        string `gorm:"type:varchar(255);not null" json:"host"`
	Port        int    `gorm:"not null" json:"port"`
	Database    string `gorm:"type:varchar(100);not null" json:"database"`
	Username    string `gorm:"type:varchar(100);not null" json:"username"`
	Password    string `gorm:"type:varchar(500);not null" json:"-"` // 加密存储，不返回
	SSLMode     string `gorm:"type:varchar(20);default:'disable'" json:"ssl_mode"`
	Status      string `gorm:"type:varchar(50);default:'inactive'" json:"status"` // inactive, active, error
	
	// 连接配置
	Config      string `gorm:"type:jsonb" json:"config,omitempty"`
	
	// 测试信息
	LastTestAt  *time.Time `gorm:"" json:"last_test_at,omitempty"`
	LastError   string     `gorm:"type:text" json:"last_error,omitempty"`
	
	// 外键
	WorkspaceID uint `gorm:"not null;index" json:"workspace_id"`
	CreatorID   uint `gorm:"not null;index" json:"creator_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (DatabaseConnection) TableName() string {
	return "database_connections"
}

