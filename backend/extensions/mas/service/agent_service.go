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
	"gorm.io/gorm"
	"github.com/coze-dev/coze-studio/backend/extensions/mas/models"
	"github.com/coze-dev/coze-studio/backend/extensions/mas/repository"
)

// AgentService Agent 服务层
type AgentService struct {
	agentRepo *repository.AgentRepository
	db        *gorm.DB
}

// NewAgentService 创建 Agent 服务
func NewAgentService(db *gorm.DB) *AgentService {
	return &AgentService{
		agentRepo: repository.NewAgentRepository(db),
		db:         db,
	}
}

// CreateAgent 创建 Agent
func (s *AgentService) CreateAgent(ctx context.Context, agent *models.Agent) error {
	// 验证必填字段
	if agent.Name == "" {
		return fmt.Errorf("agent name is required")
	}
	if agent.Type == "" {
		return fmt.Errorf("agent type is required")
	}
	if agent.WorkspaceID == 0 {
		return fmt.Errorf("workspace_id is required")
	}
	if agent.CreatorID == 0 {
		return fmt.Errorf("creator_id is required")
	}
	
	// 设置默认值
	if agent.Status == "" {
		agent.Status = "inactive"
	}
	if agent.Version == "" {
		agent.Version = "1.0.0"
	}
	
	// JSONB 字段的验证和默认值设置由模型的 BeforeCreate 钩子处理
	// 这里只需要确保字段不为 nil
	if len(agent.Config) == 0 {
		agent.Config = []byte("{}")
	}
	if len(agent.Capabilities) == 0 {
		agent.Capabilities = []byte("[]")
	}
	if len(agent.Constraints) == 0 {
		agent.Constraints = []byte("[]")
	}
	if len(agent.Tags) == 0 {
		agent.Tags = []byte("[]")
	}
	
	return s.agentRepo.Create(ctx, agent)
}

// GetAgent 获取 Agent
func (s *AgentService) GetAgent(ctx context.Context, id uint) (*models.Agent, error) {
	return s.agentRepo.GetByID(ctx, id)
}

// ListAgents 获取 Agent 列表
func (s *AgentService) ListAgents(ctx context.Context, workspaceID uint, page, pageSize int, filters map[string]interface{}) ([]models.Agent, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	return s.agentRepo.List(ctx, workspaceID, page, pageSize, filters)
}

// UpdateAgent 更新 Agent
func (s *AgentService) UpdateAgent(ctx context.Context, agent *models.Agent) error {
	// 验证 Agent 是否存在
	existing, err := s.agentRepo.GetByID(ctx, agent.ID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}
	
	// 保留一些不可修改的字段
	agent.CreatedAt = existing.CreatedAt
	agent.WorkspaceID = existing.WorkspaceID
	agent.CreatorID = existing.CreatorID
	
	// JSONB 字段的验证和默认值设置由模型的 BeforeUpdate 钩子处理
	// 如果新值为空，保留旧值
	if len(agent.Config) == 0 {
		agent.Config = existing.Config
	}
	if len(agent.Capabilities) == 0 {
		agent.Capabilities = existing.Capabilities
	}
	if len(agent.Constraints) == 0 {
		agent.Constraints = existing.Constraints
	}
	if len(agent.Tags) == 0 {
		agent.Tags = existing.Tags
	}
	if len(agent.WorkflowConfig) == 0 {
		agent.WorkflowConfig = existing.WorkflowConfig
	}
	if len(agent.MemoryConfig) == 0 {
		agent.MemoryConfig = existing.MemoryConfig
	}
	
	return s.agentRepo.Update(ctx, agent)
}

// DeleteAgent 删除 Agent
func (s *AgentService) DeleteAgent(ctx context.Context, id uint) error {
	// 验证 Agent 是否存在
	_, err := s.agentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}
	
	// TODO: 检查是否有正在运行的任务，如果有则不允许删除
	
	return s.agentRepo.Delete(ctx, id)
}

// UpdateAgentStatus 更新 Agent 状态
func (s *AgentService) UpdateAgentStatus(ctx context.Context, id uint, status string) error {
	validStatuses := map[string]bool{
		"inactive": true,
		"active":   true,
		"busy":     true,
		"error":    true,
	}
	
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}
	
	return s.agentRepo.UpdateStatus(ctx, id, status)
}

// UpdateAgentMetrics 更新 Agent 性能指标
func (s *AgentService) UpdateAgentMetrics(ctx context.Context, id uint, metrics map[string]interface{}) error {
	return s.agentRepo.UpdateMetrics(ctx, id, metrics)
}

