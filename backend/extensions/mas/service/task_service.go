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

package service

import (
	"context"
	"fmt"
	"time"
	"gorm.io/gorm"
	"github.com/coze-dev/coze-studio/backend/extensions/mas/models"
	"github.com/coze-dev/coze-studio/backend/extensions/mas/repository"
)

// TaskService MAS Task 服务层
type TaskService struct {
	taskRepo    *repository.TaskRepository
	agentRepo   *repository.AgentRepository
	sessionRepo *repository.SessionRepository
	db          *gorm.DB
}

// NewTaskService 创建 Task 服务
func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{
		taskRepo:    repository.NewTaskRepository(db),
		agentRepo:   repository.NewAgentRepository(db),
		sessionRepo: repository.NewSessionRepository(db),
		db:          db,
	}
}

// CreateTask 创建 Task
func (s *TaskService) CreateTask(ctx context.Context, task *models.MASTask) error {
	// 验证必填字段
	if task.Name == "" {
		return fmt.Errorf("task name is required")
	}
	if task.Type == "" {
		return fmt.Errorf("task type is required")
	}
	if task.AgentID == 0 {
		return fmt.Errorf("agent_id is required")
	}
	if task.WorkspaceID == 0 {
		return fmt.Errorf("workspace_id is required")
	}
	if task.CreatorID == 0 {
		return fmt.Errorf("creator_id is required")
	}
	
	// 验证 Agent 是否存在
	_, err := s.agentRepo.GetByID(ctx, task.AgentID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}
	
	// 如果指定了 Session，验证 Session 是否存在
	if task.SessionID != nil && *task.SessionID > 0 {
		_, err := s.sessionRepo.GetByID(ctx, *task.SessionID)
		if err != nil {
			return fmt.Errorf("session not found: %w", err)
		}
	}
	
	// 设置默认值
	if task.Status == "" {
		task.Status = "pending"
	}
	
	// JSONB 字段的验证和默认值设置由模型的 BeforeCreate 钩子处理
	// 这里只需要确保字段不为空
	if len(task.Input) == 0 {
		task.Input = []byte("{}")
	}
	if len(task.Output) == 0 {
		task.Output = []byte("{}")
	}
	if len(task.Dependencies) == 0 {
		task.Dependencies = []byte("[]")
	}
	
	return s.taskRepo.Create(ctx, task)
}

// GetTask 获取 Task
func (s *TaskService) GetTask(ctx context.Context, id uint) (*models.MASTask, error) {
	return s.taskRepo.GetByID(ctx, id)
}

// ListTasks 获取 Task 列表
func (s *TaskService) ListTasks(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]models.MASTask, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	return s.taskRepo.List(ctx, filters, page, pageSize)
}

// UpdateTask 更新 Task
func (s *TaskService) UpdateTask(ctx context.Context, task *models.MASTask) error {
	// 验证 Task 是否存在
	existing, err := s.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	
	// 保留一些不可修改的字段
	task.CreatedAt = existing.CreatedAt
	task.WorkspaceID = existing.WorkspaceID
	task.CreatorID = existing.CreatorID
	
	// JSONB 字段的验证和默认值设置由模型的 BeforeUpdate 钩子处理
	// 如果新值为空，保留旧值
	if len(task.Input) == 0 {
		task.Input = existing.Input
	}
	if len(task.Output) == 0 {
		task.Output = existing.Output
	}
	if len(task.Dependencies) == 0 {
		task.Dependencies = existing.Dependencies
	}
	
	return s.taskRepo.Update(ctx, task)
}

// DeleteTask 删除 Task
func (s *TaskService) DeleteTask(ctx context.Context, id uint) error {
	// 验证 Task 是否存在
	_, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	
	return s.taskRepo.Delete(ctx, id)
}

// UpdateTaskStatus 更新 Task 状态
func (s *TaskService) UpdateTaskStatus(ctx context.Context, id uint, status string) error {
	validStatuses := map[string]bool{
		"pending":   true,
		"running":   true,
		"completed": true,
		"failed":    true,
		"cancelled": true,
	}
	
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}
	
	// 根据状态更新相关时间戳
	updates := map[string]interface{}{
		"status": status,
	}
	
	now := time.Now()
	if status == "running" {
		updates["started_at"] = &now
	} else if status == "completed" || status == "failed" || status == "cancelled" {
		updates["completed_at"] = &now
		
		// 计算持续时间
		task, err := s.taskRepo.GetByID(ctx, id)
		if err == nil && task.StartedAt != nil {
			duration := now.Sub(*task.StartedAt)
			updates["duration"] = int(duration.Milliseconds())
		}
	}
	
	return s.db.WithContext(context.Background()).
		Model(&models.MASTask{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// StartTask 开始执行 Task
func (s *TaskService) StartTask(ctx context.Context, id uint) error {
	return s.UpdateTaskStatus(ctx, id, "running")
}

// CompleteTask 完成 Task
func (s *TaskService) CompleteTask(ctx context.Context, id uint, output interface{}) error {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	
	// 更新输出和状态
	// TODO: 序列化 output 为 JSON
	task.Status = "completed"
	now := time.Now()
	task.CompletedAt = &now
	if task.StartedAt != nil {
		duration := now.Sub(*task.StartedAt)
		task.Duration = int(duration.Milliseconds())
	}
	
	return s.taskRepo.Update(ctx, task)
}

// FailTask 标记 Task 为失败
func (s *TaskService) FailTask(ctx context.Context, id uint, errorMsg string) error {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	
	task.Status = "failed"
	task.Error = errorMsg
	now := time.Now()
	task.CompletedAt = &now
	if task.StartedAt != nil {
		duration := now.Sub(*task.StartedAt)
		task.Duration = int(duration.Milliseconds())
	}
	
	return s.taskRepo.Update(ctx, task)
}

// GetPendingTasks 获取待处理的任务
func (s *TaskService) GetPendingTasks(ctx context.Context, limit int) ([]models.MASTask, error) {
	return s.taskRepo.GetPendingTasks(ctx, limit)
}

