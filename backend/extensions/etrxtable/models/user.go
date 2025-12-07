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

// User 用户模型 - 从 Ent Schema 迁移到 GORM Model
type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"` // 不返回密码
	FirstName    string `gorm:"type:varchar(50)" json:"first_name,omitempty"`
	LastName     string `gorm:"type:varchar(50)" json:"last_name,omitempty"`
	Avatar       string `gorm:"type:varchar(255)" json:"avatar,omitempty"`
	IsActive     bool   `gorm:"default:true;index" json:"is_active"`
	IsSuperuser  bool   `gorm:"default:false;index" json:"is_superuser"`
	LastLoginAt  *time.Time `gorm:"" json:"last_login_at,omitempty"`
	Phone        string `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Address      string `gorm:"type:varchar(255)" json:"address,omitempty"`
	Timezone     string `gorm:"type:varchar(50);default:'UTC'" json:"timezone"`
	Language     string `gorm:"type:varchar(10);default:'zh-CN'" json:"language"`
	Metadata     string `gorm:"type:jsonb" json:"metadata,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系（使用 Preload 加载）
	Roles          []Role          `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	Workspaces     []Workspace     `gorm:"foreignKey:OwnerID" json:"workspaces,omitempty"`
	CreatedTables  []TableBody     `gorm:"foreignKey:OwnerID" json:"created_tables,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeSave GORM 钩子 - 保存前
func (u *User) BeforeSave(tx *gorm.DB) error {
	// 可以在这里添加保存前的逻辑，如密码加密
	return nil
}

