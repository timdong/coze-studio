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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/extensions/mas/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate
	err = db.AutoMigrate(&models.Agent{})
	require.NoError(t, err)

	return db
}

func TestAgentService_CreateAgent(t *testing.T) {
	db := setupTestDB(t)
	service := NewAgentService(db)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		agent := &models.Agent{
			Name:        "Test Agent",
			Type:        "simple",
			WorkspaceID: 1,
			CreatorID:   1,
			Status:      "inactive",
			Version:     "1.0.0",
		}

		err := service.CreateAgent(ctx, agent)
		assert.NoError(t, err)
		assert.NotZero(t, agent.ID)
	})

	t.Run("missing name", func(t *testing.T) {
		agent := &models.Agent{
			Type:        "simple",
			WorkspaceID: 1,
			CreatorID:   1,
		}

		err := service.CreateAgent(ctx, agent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "agent name is required")
	})

	t.Run("missing type", func(t *testing.T) {
		agent := &models.Agent{
			Name:        "Test Agent",
			WorkspaceID: 1,
			CreatorID:   1,
		}

		err := service.CreateAgent(ctx, agent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "agent type is required")
	})

	t.Run("missing workspace_id", func(t *testing.T) {
		agent := &models.Agent{
			Name:      "Test Agent",
			Type:      "simple",
			CreatorID: 1,
		}

		err := service.CreateAgent(ctx, agent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "workspace_id is required")
	})

	t.Run("default values", func(t *testing.T) {
		agent := &models.Agent{
			Name:        "Test Agent",
			Type:        "simple",
			WorkspaceID: 1,
			CreatorID:   1,
		}

		err := service.CreateAgent(ctx, agent)
		assert.NoError(t, err)
		assert.Equal(t, "inactive", agent.Status)
		assert.Equal(t, "1.0.0", agent.Version)
	})
}

func TestAgentService_GetAgent(t *testing.T) {
	db := setupTestDB(t)
	service := NewAgentService(db)
	ctx := context.Background()

	// Create an agent first
	agent := &models.Agent{
		Name:        "Test Agent",
		Type:        "simple",
		WorkspaceID: 1,
		CreatorID:   1,
	}
	err := service.CreateAgent(ctx, agent)
	require.NoError(t, err)

	t.Run("get existing agent", func(t *testing.T) {
		retrieved, err := service.GetAgent(ctx, agent.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, agent.ID, retrieved.ID)
		assert.Equal(t, "Test Agent", retrieved.Name)
	})

	t.Run("get non-existent agent", func(t *testing.T) {
		_, err := service.GetAgent(ctx, 999)
		assert.Error(t, err)
	})
}

func TestAgentService_ListAgents(t *testing.T) {
	db := setupTestDB(t)
	service := NewAgentService(db)
	ctx := context.Background()

	// Create multiple agents
	for i := 0; i < 5; i++ {
		agent := &models.Agent{
			Name:        "Test Agent",
			Type:        "simple",
			WorkspaceID: 1,
			CreatorID:   1,
		}
		err := service.CreateAgent(ctx, agent)
		require.NoError(t, err)
	}

	t.Run("list all agents", func(t *testing.T) {
		agents, total, err := service.ListAgents(ctx, 1, 1, 10, nil)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, agents, 5)
	})

	t.Run("pagination", func(t *testing.T) {
		agents, total, err := service.ListAgents(ctx, 1, 1, 2, nil)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, agents, 2)
	})
}

func TestAgentService_UpdateAgentStatus(t *testing.T) {
	db := setupTestDB(t)
	service := NewAgentService(db)
	ctx := context.Background()

	// Create an agent
	agent := &models.Agent{
		Name:        "Test Agent",
		Type:        "simple",
		WorkspaceID: 1,
		CreatorID:   1,
		Status:      "inactive",
	}
	err := service.CreateAgent(ctx, agent)
	require.NoError(t, err)

	t.Run("update status to active", func(t *testing.T) {
		err := service.UpdateAgentStatus(ctx, agent.ID, "active")
		assert.NoError(t, err)

		updated, err := service.GetAgent(ctx, agent.ID)
		assert.NoError(t, err)
		assert.Equal(t, "active", updated.Status)
	})

	t.Run("invalid status", func(t *testing.T) {
		err := service.UpdateAgentStatus(ctx, agent.ID, "invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
	})
}

