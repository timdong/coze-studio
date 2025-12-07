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

package repository

import (
	"context"
	"gorm.io/gorm"
	"github.com/coze-dev/coze-studio/backend/extensions/mas/models"
)

// SessionRepository MAS Session 数据访问层
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository 创建 Session 仓库
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create 创建 Session
func (r *SessionRepository) Create(ctx context.Context, session *models.MASSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByID 根据 ID 获取 Session
func (r *SessionRepository) GetByID(ctx context.Context, id uint) (*models.MASSession, error) {
	var session models.MASSession
	err := r.db.WithContext(ctx).
		Preload("Workspace").
		Preload("Creator").
		Preload("Tasks").
		Preload("Tasks.Agent").
		First(&session, id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// List 获取 Session 列表
func (r *SessionRepository) List(ctx context.Context, workspaceID uint, page, pageSize int, filters map[string]interface{}) ([]models.MASSession, int64, error) {
	var sessions []models.MASSession
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.MASSession{})
	
	// 过滤工作空间
	if workspaceID > 0 {
		query = query.Where("workspace_id = ?", workspaceID)
	}
	
	// 应用其他过滤条件
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if sessionType, ok := filters["type"].(string); ok && sessionType != "" {
		query = query.Where("type = ?", sessionType)
	}
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.
		Preload("Workspace").
		Preload("Creator").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&sessions).Error
	
	return sessions, total, err
}

// Update 更新 Session
func (r *SessionRepository) Update(ctx context.Context, session *models.MASSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// Delete 删除 Session（软删除）
func (r *SessionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.MASSession{}, id).Error
}

// UpdateStatus 更新 Session 状态
func (r *SessionRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.MASSession{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetByWorkspace 获取工作空间的所有 Session
func (r *SessionRepository) GetByWorkspace(ctx context.Context, workspaceID uint) ([]models.MASSession, error) {
	var sessions []models.MASSession
	err := r.db.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

