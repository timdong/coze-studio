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

package migrations

import (
	"gorm.io/gorm"

	etrxtableModels "github.com/coze-dev/coze-studio/backend/extensions/etrxtable/models"
	erDiagramModels "github.com/coze-dev/coze-studio/backend/extensions/er-diagram/models"
	masModels "github.com/coze-dev/coze-studio/backend/extensions/mas/models"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	// 为 users 表添加 Coze Studio 需要的字段
	if err := addCozeStudioFieldsToUsersTable(db); err != nil {
		return err
	}

	// Coze Studio 核心域模型（Space 和 SpaceUser）
	if err := createCozeStudioCoreTables(db); err != nil {
		return err
	}

	// 用户和权限系统（扩展功能）
	// 注意：Coze Studio 核心用户系统使用 users 表（与扩展功能共享）
	if err := db.AutoMigrate(
		&etrxtableModels.User{},
		&etrxtableModels.Role{},
		&etrxtableModels.Workspace{},
		&etrxtableModels.Menu{},
		&etrxtableModels.Company{},
		&etrxtableModels.Department{},
	); err != nil {
		return err
	}

	// 多维表格系统
	if err := db.AutoMigrate(
		&etrxtableModels.TableBody{},
		&etrxtableModels.TableColumn{},
		&etrxtableModels.TableRow{},
		&etrxtableModels.TableView{},
		&etrxtableModels.TableSort{},
		&etrxtableModels.TableLink{},
		&etrxtableModels.DynamicTableMetadata{},
		&etrxtableModels.View{},
	); err != nil {
		return err
	}

	// ER 图系统
	if err := db.AutoMigrate(
		&erDiagramModels.ERDiagram{},
		&etrxtableModels.ERTableSyncHistory{},
	); err != nil {
		return err
	}

	// MAS 多智能体系统
	if err := db.AutoMigrate(
		&masModels.Agent{},
		&masModels.MASSession{},
		&masModels.MASTask{},
		&masModels.AgentVersion{},
		&masModels.AgentMemory{},
		&masModels.AgentMetric{},
		&masModels.AgentTag{},
		&masModels.AgentTool{},
		&masModels.AgentMessage{},
		&masModels.AgentExecutionLog{},
		&masModels.AgentFavorite{},
		&masModels.AgentCollaboration{},
		&masModels.AgentKnowledge{},
		&masModels.MASTemplate{},
	); err != nil {
		return err
	}

	// 工作流系统
	if err := db.AutoMigrate(
		&etrxtableModels.Workflow{},
		&etrxtableModels.WorkflowExecution{},
		&etrxtableModels.WorkflowTemplate{},
		&etrxtableModels.WorkflowPermission{},
		&etrxtableModels.WorkflowVersion{},
		&etrxtableModels.WorkflowExecutionLog{},
		&etrxtableModels.TemplateCategory{},
		&etrxtableModels.TemplateRating{},
	); err != nil {
		return err
	}

	// RAG 系统
	if err := db.AutoMigrate(
		&etrxtableModels.DocumentCategory{},
		&etrxtableModels.Document{},
		&etrxtableModels.DocumentQuality{},
		&etrxtableModels.DocumentVersion{},
		&etrxtableModels.Chunk{},
		&etrxtableModels.Embedding{},
		&etrxtableModels.Conversation{},
		&etrxtableModels.QAHistory{},
	); err != nil {
		return err
	}

	// 插件和 LLM 系统
	if err := db.AutoMigrate(
		&etrxtableModels.Plugin{},
		&etrxtableModels.LLMProvider{},
		&etrxtableModels.LLMProviderCredential{},
		&etrxtableModels.LLMModelInstance{},
		&etrxtableModels.LLMModel{},
		&etrxtableModels.LLMModelAuditLog{},
	); err != nil {
		return err
	}

	// 其他系统
	if err := db.AutoMigrate(
		&etrxtableModels.DatabaseConnection{},
	); err != nil {
		return err
	}

	return nil
}

// addCozeStudioFieldsToUsersTable 为 users 表添加 Coze Studio 需要的字段
func addCozeStudioFieldsToUsersTable(db *gorm.DB) error {
	sql := `
	-- 添加 Coze Studio 需要的字段（如果不存在）
	DO $$
	BEGIN
		-- name 字段（如果不存在，使用 username 作为默认值）
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='name') THEN
			ALTER TABLE users ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '';
			UPDATE users SET name = username WHERE name = '';
		END IF;

		-- unique_name 字段
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='unique_name') THEN
			ALTER TABLE users ADD COLUMN unique_name VARCHAR(255) NOT NULL DEFAULT '';
			UPDATE users SET unique_name = username WHERE unique_name = '';
		END IF;

		-- description 字段
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='description') THEN
			ALTER TABLE users ADD COLUMN description TEXT NOT NULL DEFAULT '';
		END IF;

		-- icon_uri 字段（如果不存在，使用 avatar 作为默认值）
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='icon_uri') THEN
			ALTER TABLE users ADD COLUMN icon_uri VARCHAR(255) NOT NULL DEFAULT '';
			UPDATE users SET icon_uri = COALESCE(avatar, '') WHERE icon_uri = '';
		END IF;

		-- user_verified 字段（如果不存在，使用 is_active 作为默认值）
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='user_verified') THEN
			ALTER TABLE users ADD COLUMN user_verified BOOLEAN NOT NULL DEFAULT false;
			UPDATE users SET user_verified = is_active WHERE user_verified = false;
		END IF;

		-- locale 字段（如果不存在，使用 language 作为默认值）
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='locale') THEN
			ALTER TABLE users ADD COLUMN locale VARCHAR(50) NOT NULL DEFAULT 'zh-CN';
			UPDATE users SET locale = COALESCE(language, 'zh-CN') WHERE locale = 'zh-CN';
		END IF;

		-- session_key 字段
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='session_key') THEN
			ALTER TABLE users ADD COLUMN session_key VARCHAR(255) NOT NULL DEFAULT '';
		END IF;

		-- created_at_ms 字段（存储毫秒时间戳）
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='created_at_ms') THEN
			ALTER TABLE users ADD COLUMN created_at_ms BIGINT NOT NULL DEFAULT 0;
			UPDATE users SET created_at_ms = EXTRACT(EPOCH FROM created_at)::BIGINT * 1000 WHERE created_at_ms = 0 AND created_at IS NOT NULL;
		END IF;

		-- updated_at_ms 字段（存储毫秒时间戳）
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='updated_at_ms') THEN
			ALTER TABLE users ADD COLUMN updated_at_ms BIGINT NOT NULL DEFAULT 0;
			UPDATE users SET updated_at_ms = EXTRACT(EPOCH FROM updated_at)::BIGINT * 1000 WHERE updated_at_ms = 0 AND updated_at IS NOT NULL;
		END IF;

		-- 创建索引
		CREATE INDEX IF NOT EXISTS idx_users_unique_name ON users(unique_name);
		CREATE INDEX IF NOT EXISTS idx_users_session_key ON users(session_key);
	END $$;
	`

	return db.Exec(sql).Error
}

// createCozeStudioCoreTables 创建 Coze Studio 核心表
// 包括：Space, SpaceUser, 以及所有其他 Coze Studio 核心域模型表
func createCozeStudioCoreTables(db *gorm.DB) error {
	// 由于表很多，分批创建以提高可读性和维护性
	if err := createCozeStudioBasicTables(db); err != nil {
		return err
	}
	if err := createCozeStudioAgentTables(db); err != nil {
		return err
	}
	if err := createCozeStudioAppTables(db); err != nil {
		return err
	}
	if err := createCozeStudioWorkflowTables(db); err != nil {
		return err
	}
	if err := createCozeStudioKnowledgeTables(db); err != nil {
		return err
	}
	if err := createCozeStudioPluginTables(db); err != nil {
		return err
	}
	if err := createCozeStudioConversationTables(db); err != nil {
		return err
	}
	if err := createCozeStudioMemoryTables(db); err != nil {
		return err
	}
	if err := createCozeStudioOtherTables(db); err != nil {
		return err
	}
	return nil
}

// createCozeStudioBasicTables 创建基础表（Space, SpaceUser, Template）
func createCozeStudioBasicTables(db *gorm.DB) error {
	sql := `
	-- 创建 space 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'space') THEN
			CREATE TABLE space (
				id BIGSERIAL PRIMARY KEY,
				owner_id BIGINT NOT NULL DEFAULT 0,
				name VARCHAR(200) NOT NULL DEFAULT '',
				description VARCHAR(2000) NOT NULL DEFAULT '',
				icon_uri VARCHAR(200) NOT NULL DEFAULT '',
				creator_id BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT
			);
			CREATE INDEX IF NOT EXISTS idx_space_creator_id ON space(creator_id);
			CREATE INDEX IF NOT EXISTS idx_space_owner_id ON space(owner_id);
		END IF;
	END $$;

	-- 创建 space_user 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'space_user') THEN
			CREATE TABLE space_user (
				id BIGSERIAL PRIMARY KEY,
				space_id BIGINT NOT NULL DEFAULT 0,
				user_id BIGINT NOT NULL DEFAULT 0,
				role_type INTEGER NOT NULL DEFAULT 3,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_space_user_user_id ON space_user(user_id);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_space_user ON space_user(space_id, user_id);
		END IF;
	END $$;

	-- 创建 template 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'template') THEN
			CREATE TABLE template (
				id BIGSERIAL PRIMARY KEY,
				agent_id BIGINT NOT NULL DEFAULT 0,
				workflow_id BIGINT NOT NULL DEFAULT 0,
				space_id BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				heat BIGINT NOT NULL DEFAULT 0,
				product_entity_type BIGINT NOT NULL DEFAULT 0,
				meta_info JSONB,
				agent_extra JSONB,
				workflow_extra JSONB,
				project_extra JSONB
			);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_agent_id ON template(agent_id);
		END IF;
	END $$;
	`
	return db.Exec(sql).Error
}

