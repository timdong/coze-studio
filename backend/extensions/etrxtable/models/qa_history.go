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

// QAHistory 问答历史模型 - 从 Ent Schema 迁移到 GORM Model
type QAHistory struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Question string `gorm:"type:text;not null" json:"question"`
	Answer   string `gorm:"type:text;not null" json:"answer"`
	
	// 检索信息
	SimilarChunks string `gorm:"type:jsonb" json:"similar_chunks,omitempty"` // 相似文档块
	RetrievalTime int    `gorm:"default:0" json:"retrieval_time"` // 检索耗时（毫秒）
	GenerationTime int   `gorm:"default:0" json:"generation_time"` // 生成耗时（毫秒）
	
	// 质量评分
	Confidence  float64 `gorm:"default:0" json:"confidence"`  // 置信度 0-1
	Relevance   float64 `gorm:"default:0" json:"relevance"`   // 相关性 0-1
	UserRating  *int    `gorm:"" json:"user_rating,omitempty"` // 用户评分 1-5
	
	// 模型信息
	Model       string `gorm:"type:varchar(100)" json:"model,omitempty"`
	Temperature float64 `gorm:"default:0.7" json:"temperature"`
	MaxTokens   int     `gorm:"default:4000" json:"max_tokens"`
	
	// 外键
	ConversationID uint  `gorm:"not null;index" json:"conversation_id"`
	UserID         uint  `gorm:"not null;index" json:"user_id"`
	WorkspaceID    uint  `gorm:"not null;index" json:"workspace_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关系
	Conversation *Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
	User         *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Workspace    *Workspace    `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}

// TableName 指定表名
func (QAHistory) TableName() string {
	return "qa_histories"
}

