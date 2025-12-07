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

// TaskRepository MAS Task 数据访问层
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository 创建 Task 仓库
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create 创建 Task
func (r *TaskRepository) Create(ctx context.Context, task *models.MASTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// GetByID 根据 ID 获取 Task
func (r *TaskRepository) GetByID(ctx context.Context, id uint) (*models.MASTask, error) {
	var task models.MASTask
	err := r.db.WithContext(ctx).
		Preload("Agent").
		Preload("Workspace").
		Preload("Creator").
		First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// List 获取 Task 列表
func (r *TaskRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]models.MASTask, int64, error) {
	var tasks []models.MASTask
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.MASTask{})
	
	// 应用过滤条件
	if agentID, ok := filters["agent_id"].(uint); ok && agentID > 0 {
		query = query.Where("agent_id = ?", agentID)
	}
	if sessionID, ok := filters["session_id"].(uint); ok && sessionID > 0 {
		query = query.Where("session_id = ?", sessionID)
	}
	if workspaceID, ok := filters["workspace_id"].(uint); ok && workspaceID > 0 {
		query = query.Where("workspace_id = ?", workspaceID)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if taskType, ok := filters["type"].(string); ok && taskType != "" {
		query = query.Where("type = ?", taskType)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.
		Preload("Agent").
		Preload("Workspace").
		Preload("Creator").
		Order("priority DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&tasks).Error
	
	return tasks, total, err
}

// Update 更新 Task
func (r *TaskRepository) Update(ctx context.Context, task *models.MASTask) error {
	return r.db.WithContext(ctx).Save(task).Error
}

// Delete 删除 Task（软删除）
func (r *TaskRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.MASTask{}, id).Error
}

// UpdateStatus 更新 Task 状态
func (r *TaskRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.MASTask{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetByAgent 获取 Agent 的所有 Task
func (r *TaskRepository) GetByAgent(ctx context.Context, agentID uint) ([]models.MASTask, error) {
	var tasks []models.MASTask
	err := r.db.WithContext(ctx).
		Where("agent_id = ?", agentID).
		Order("priority DESC, created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

// GetBySession 获取 Session 的所有 Task
func (r *TaskRepository) GetBySession(ctx context.Context, sessionID uint) ([]models.MASTask, error) {
	var tasks []models.MASTask
	err := r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("priority DESC, created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

// GetPendingTasks 获取待处理的任务
func (r *TaskRepository) GetPendingTasks(ctx context.Context, limit int) ([]models.MASTask, error) {
	var tasks []models.MASTask
	query := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Order("priority DESC, created_at ASC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.
		Preload("Agent").
		Find(&tasks).Error
	return tasks, err
}

