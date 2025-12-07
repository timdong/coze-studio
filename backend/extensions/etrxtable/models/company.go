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

// Company 公司模型 - 从 Ent Schema 迁移到 GORM Model
type Company struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description string `gorm:"type:varchar(500)" json:"description,omitempty"`
	Address     string `gorm:"type:varchar(500)" json:"address,omitempty"`
	Phone       string `gorm:"type:varchar(50)" json:"phone,omitempty"`
	Email       string `gorm:"type:varchar(100)" json:"email,omitempty"`
	Website     string `gorm:"type:varchar(255)" json:"website,omitempty"`
	Logo        string `gorm:"type:varchar(500)" json:"logo,omitempty"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Departments []Department `gorm:"foreignKey:CompanyID" json:"departments,omitempty"`
}

// TableName 指定表名
func (Company) TableName() string {
	return "companies"
}

// Department 部门模型
type Department struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:varchar(500)" json:"description,omitempty"`
	ParentID    *uint  `gorm:"index" json:"parent_id,omitempty"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	
	// 外键
	CompanyID uint `gorm:"not null;index" json:"company_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Company  *Company     `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Parent   *Department  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName 指定表名
func (Department) TableName() string {
	return "departments"
}

