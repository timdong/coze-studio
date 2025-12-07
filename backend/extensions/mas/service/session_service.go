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

// SessionService MAS Session 服务层
type SessionService struct {
	sessionRepo *repository.SessionRepository
	taskRepo    *repository.TaskRepository
	db          *gorm.DB
}

// NewSessionService 创建 Session 服务
func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{
		sessionRepo: repository.NewSessionRepository(db),
		taskRepo:    repository.NewTaskRepository(db),
		db:          db,
	}
}

// CreateSession 创建 Session
func (s *SessionService) CreateSession(ctx context.Context, session *models.MASSession) error {
	// 验证必填字段
	if session.Name == "" {
		return fmt.Errorf("session name is required")
	}
	if session.Type == "" {
		return fmt.Errorf("session type is required")
	}
	if session.WorkspaceID == 0 {
		return fmt.Errorf("workspace_id is required")
	}
	if session.CreatorID == 0 {
		return fmt.Errorf("creator_id is required")
	}
	
	// 设置默认值
	if session.Status == "" {
		session.Status = "active"
	}
	now := time.Now()
	if session.StartedAt == nil {
		session.StartedAt = &now
	}
	
	// JSONB 字段的验证和默认值设置由模型的 BeforeCreate 钩子处理
	// 这里只需要确保字段不为空
	if len(session.Config) == 0 {
		session.Config = []byte("{}")
	}
	if len(session.Context) == 0 {
		session.Context = []byte("{}")
	}
	if len(session.AgentIDs) == 0 {
		session.AgentIDs = []byte("[]")
	}
	
	return s.sessionRepo.Create(ctx, session)
}

// GetSession 获取 Session
func (s *SessionService) GetSession(ctx context.Context, id uint) (*models.MASSession, error) {
	return s.sessionRepo.GetByID(ctx, id)
}

// ListSessions 获取 Session 列表
func (s *SessionService) ListSessions(ctx context.Context, workspaceID uint, page, pageSize int, filters map[string]interface{}) ([]models.MASSession, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	return s.sessionRepo.List(ctx, workspaceID, page, pageSize, filters)
}

// UpdateSession 更新 Session
func (s *SessionService) UpdateSession(ctx context.Context, session *models.MASSession) error {
	// 验证 Session 是否存在
	existing, err := s.sessionRepo.GetByID(ctx, session.ID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}
	
	// 保留一些不可修改的字段
	session.CreatedAt = existing.CreatedAt
	session.WorkspaceID = existing.WorkspaceID
	session.CreatorID = existing.CreatorID
	
	// JSONB 字段的验证和默认值设置由模型的 BeforeUpdate 钩子处理
	// 如果新值为空，保留旧值
	if len(session.Config) == 0 {
		session.Config = existing.Config
	}
	if len(session.Context) == 0 {
		session.Context = existing.Context
	}
	if len(session.AgentIDs) == 0 {
		session.AgentIDs = existing.AgentIDs
	}
	
	return s.sessionRepo.Update(ctx, session)
}

// DeleteSession 删除 Session
func (s *SessionService) DeleteSession(ctx context.Context, id uint) error {
	// 验证 Session 是否存在
	_, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}
	
	// 在事务中删除 Session 及其关联的 Task
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除 Session 的所有 Task
		if err := tx.Where("session_id = ?", id).Delete(&models.MASTask{}).Error; err != nil {
			return fmt.Errorf("failed to delete tasks: %w", err)
		}
		
		// 删除 Session
		if err := tx.Delete(&models.MASSession{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete session: %w", err)
		}
		
		return nil
	})
}

// UpdateSessionStatus 更新 Session 状态
func (s *SessionService) UpdateSessionStatus(ctx context.Context, id uint, status string) error {
	validStatuses := map[string]bool{
		"active":    true,
		"completed": true,
		"failed":    true,
		"cancelled": true,
	}
	
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}
	
	// 如果状态变为 completed，更新完成时间
	updates := map[string]interface{}{
		"status": status,
	}
	if status == "completed" {
		now := time.Now()
		updates["completed_at"] = &now
	}
	
	return s.db.WithContext(context.Background()).
		Model(&models.MASSession{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// CompleteSession 完成 Session
func (s *SessionService) CompleteSession(ctx context.Context, id uint) error {
	session, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}
	
	// 计算持续时间
	now := time.Now()
	if session.StartedAt != nil {
		duration := now.Sub(*session.StartedAt)
		session.Duration = int(duration.Milliseconds())
	}
	session.CompletedAt = &now
	session.Status = "completed"
	
	return s.sessionRepo.Update(ctx, session)
}

