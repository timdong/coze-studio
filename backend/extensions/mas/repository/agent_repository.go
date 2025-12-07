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

// AgentRepository Agent 数据访问层
type AgentRepository struct {
	db *gorm.DB
}

// NewAgentRepository 创建 Agent 仓库
func NewAgentRepository(db *gorm.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

// Create 创建 Agent
func (r *AgentRepository) Create(ctx context.Context, agent *models.Agent) error {
	return r.db.WithContext(ctx).Create(agent).Error
}

// GetByID 根据 ID 获取 Agent
func (r *AgentRepository) GetByID(ctx context.Context, id uint) (*models.Agent, error) {
	var agent models.Agent
	err := r.db.WithContext(ctx).
		Preload("Workspace").
		Preload("Creator").
		First(&agent, id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// List 获取 Agent 列表
func (r *AgentRepository) List(ctx context.Context, workspaceID uint, page, pageSize int, filters map[string]interface{}) ([]models.Agent, int64, error) {
	var agents []models.Agent
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.Agent{})
	
	// 过滤工作空间
	if workspaceID > 0 {
		query = query.Where("workspace_id = ?", workspaceID)
	}
	
	// 应用其他过滤条件
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if agentType, ok := filters["type"].(string); ok && agentType != "" {
		query = query.Where("type = ?", agentType)
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
		Find(&agents).Error
	
	return agents, total, err
}

// Update 更新 Agent
func (r *AgentRepository) Update(ctx context.Context, agent *models.Agent) error {
	return r.db.WithContext(ctx).Save(agent).Error
}

// Delete 删除 Agent（软删除）
func (r *AgentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Agent{}, id).Error
}

// GetByWorkspaceAndName 根据工作空间和名称获取 Agent
func (r *AgentRepository) GetByWorkspaceAndName(ctx context.Context, workspaceID uint, name string) (*models.Agent, error) {
	var agent models.Agent
	err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND name = ?", workspaceID, name).
		First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// CountByWorkspace 统计工作空间的 Agent 数量
func (r *AgentRepository) CountByWorkspace(ctx context.Context, workspaceID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Agent{}).
		Where("workspace_id = ?", workspaceID).
		Count(&count).Error
	return count, err
}

// UpdateStatus 更新 Agent 状态
func (r *AgentRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.Agent{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateMetrics 更新 Agent 性能指标
func (r *AgentRepository) UpdateMetrics(ctx context.Context, id uint, metrics map[string]interface{}) error {
	updates := make(map[string]interface{})
	if totalTasks, ok := metrics["total_tasks"].(int); ok {
		updates["total_tasks"] = totalTasks
	}
	if completedTasks, ok := metrics["completed_tasks"].(int); ok {
		updates["completed_tasks"] = completedTasks
	}
	if failedTasks, ok := metrics["failed_tasks"].(int); ok {
		updates["failed_tasks"] = failedTasks
	}
	if avgResponseTime, ok := metrics["avg_response_time"].(float64); ok {
		updates["avg_response_time"] = avgResponseTime
	}
	if successRate, ok := metrics["success_rate"].(float64); ok {
		updates["success_rate"] = successRate
	}
	
	return r.db.WithContext(ctx).
		Model(&models.Agent{}).
		Where("id = ?", id).
		Updates(updates).Error
}

