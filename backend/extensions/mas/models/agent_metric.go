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

// AgentMetric Agent 性能指标模型 - 从 Ent Schema 迁移到 GORM Model
type AgentMetric struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	MetricType      string  `gorm:"type:varchar(50);not null;index" json:"metric_type"` // response_time, success_rate, task_count
	MetricValue     float64 `gorm:"not null" json:"metric_value"`
	Unit            string  `gorm:"type:varchar(20)" json:"unit,omitempty"` // ms, %, count
	Timestamp       time.Time `gorm:"not null;index" json:"timestamp"`
	
	// 聚合信息
	PeriodType      string  `gorm:"type:varchar(20);default:'realtime'" json:"period_type"` // realtime, hourly, daily, weekly, monthly
	AggregateValue  float64 `gorm:"default:0" json:"aggregate_value"`
	
	// 外键
	AgentID     uint `gorm:"not null;index" json:"agent_id"`
	TaskID      *uint `gorm:"index" json:"task_id,omitempty"`
	SessionID   *uint `gorm:"index" json:"session_id,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Agent   *Agent      `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Task    *MASTask    `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Session *MASSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}

// TableName 指定表名
func (AgentMetric) TableName() string {
	return "agent_metrics"
}

