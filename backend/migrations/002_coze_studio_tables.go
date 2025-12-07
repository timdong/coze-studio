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

import "gorm.io/gorm"

// createCozeStudioAgentTables 创建 Agent 相关表
func createCozeStudioAgentTables(db *gorm.DB) error {
	sql := `
	-- 创建 single_agent_draft 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'single_agent_draft') THEN
			CREATE TABLE single_agent_draft (
				id BIGSERIAL PRIMARY KEY,
				agent_id BIGINT NOT NULL DEFAULT 0,
				creator_id BIGINT NOT NULL DEFAULT 0,
				space_id BIGINT NOT NULL DEFAULT 0,
				name VARCHAR(255) NOT NULL DEFAULT '',
				description TEXT,
				icon_uri VARCHAR(255) NOT NULL DEFAULT '',
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT,
				variables_meta_id BIGINT,
				model_info JSONB,
				onboarding_info JSONB,
				prompt JSONB,
				plugin JSONB,
				knowledge JSONB,
				workflow JSONB,
				suggest_reply JSONB,
				jump_config JSONB,
				background_image_info_list JSONB,
				database_config JSONB,
				bot_mode SMALLINT NOT NULL DEFAULT 0,
				layout_info TEXT,
				shortcut_command JSONB
			);
			CREATE INDEX IF NOT EXISTS idx_creator_id ON single_agent_draft(creator_id);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_agent_id ON single_agent_draft(agent_id);
		END IF;
	END $$;

	-- 创建 single_agent_version 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'single_agent_version') THEN
			CREATE TABLE single_agent_version (
				id BIGSERIAL PRIMARY KEY,
				agent_id BIGINT NOT NULL DEFAULT 0,
				creator_id BIGINT NOT NULL DEFAULT 0,
				space_id BIGINT NOT NULL DEFAULT 0,
				name VARCHAR(255) NOT NULL DEFAULT '',
				description TEXT,
				icon_uri VARCHAR(255) NOT NULL DEFAULT '',
				created_at BIGINT NOT NULL DEFAULT 0,
				bot_mode SMALLINT NOT NULL DEFAULT 0,
				layout_info TEXT,
				updated_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT,
				variables_meta_id BIGINT,
				model_info JSONB,
				onboarding_info JSONB,
				prompt JSONB,
				plugin JSONB,
				knowledge JSONB,
				workflow JSONB,
				suggest_reply JSONB,
				jump_config JSONB,
				connector_id BIGINT NOT NULL DEFAULT 0,
				version VARCHAR(255) NOT NULL DEFAULT '',
				background_image_info_list JSONB,
				database_config JSONB,
				shortcut_command JSONB
			);
			CREATE INDEX IF NOT EXISTS idx_creator_id ON single_agent_version(creator_id);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_agent_id_and_version_connector_id ON single_agent_version(agent_id, version, connector_id);
		END IF;
	END $$;

	-- 创建 single_agent_publish 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'single_agent_publish') THEN
			CREATE TABLE single_agent_publish (
				id BIGSERIAL PRIMARY KEY,
				agent_id BIGINT NOT NULL DEFAULT 0,
				publish_id VARCHAR(50) NOT NULL DEFAULT '',
				connector_ids JSONB,
				version VARCHAR(255) NOT NULL DEFAULT '',
				publish_info TEXT,
				publish_time BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				creator_id BIGINT NOT NULL DEFAULT 0,
				status SMALLINT NOT NULL DEFAULT 0,
				extra JSONB
			);
			CREATE INDEX IF NOT EXISTS idx_agent_id_version ON single_agent_publish(agent_id, version);
			CREATE INDEX IF NOT EXISTS idx_creator_id ON single_agent_publish(creator_id);
			CREATE INDEX IF NOT EXISTS idx_publish_id ON single_agent_publish(publish_id);
		END IF;
	END $$;

	-- 创建 agent_tool_draft 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'agent_tool_draft') THEN
			CREATE TABLE agent_tool_draft (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				agent_id BIGINT NOT NULL DEFAULT 0,
				plugin_id BIGINT NOT NULL DEFAULT 0,
				tool_id BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				sub_url VARCHAR(512) NOT NULL DEFAULT '',
				method VARCHAR(64) NOT NULL DEFAULT '',
				tool_name VARCHAR(255) NOT NULL DEFAULT '',
				tool_version VARCHAR(255) NOT NULL DEFAULT '',
				operation JSONB,
				source SMALLINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_agent_plugin_tool ON agent_tool_draft(agent_id, plugin_id, tool_id);
			CREATE INDEX IF NOT EXISTS idx_agent_tool_bind ON agent_tool_draft(agent_id, created_at);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_agent_tool_id ON agent_tool_draft(agent_id, tool_id);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_agent_tool_name ON agent_tool_draft(agent_id, tool_name);
		END IF;
	END $$;

	-- 创建 agent_tool_version 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'agent_tool_version') THEN
			CREATE TABLE agent_tool_version (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				agent_id BIGINT NOT NULL DEFAULT 0,
				plugin_id BIGINT NOT NULL DEFAULT 0,
				tool_id BIGINT NOT NULL DEFAULT 0,
				agent_version VARCHAR(255) NOT NULL DEFAULT '',
				tool_name VARCHAR(255) NOT NULL DEFAULT '',
				tool_version VARCHAR(255) NOT NULL DEFAULT '',
				sub_url VARCHAR(512) NOT NULL DEFAULT '',
				method VARCHAR(64) NOT NULL DEFAULT '',
				operation JSONB,
				created_at BIGINT NOT NULL DEFAULT 0,
				source SMALLINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_agent_tool_id_created_at ON agent_tool_version(agent_id, tool_id, created_at);
			CREATE INDEX IF NOT EXISTS idx_agent_tool_name_created_at ON agent_tool_version(agent_id, tool_name, created_at);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_agent_tool_id_agent_version ON agent_tool_version(agent_id, tool_id, agent_version);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_agent_tool_name_agent_version ON agent_tool_version(agent_id, tool_name, agent_version);
		END IF;
	END $$;
	`
	return db.Exec(sql).Error
}

// createCozeStudioAppTables 创建 App 相关表
func createCozeStudioAppTables(db *gorm.DB) error {
	sql := `
	-- 创建 app_draft 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'app_draft') THEN
			CREATE TABLE app_draft (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				space_id BIGINT NOT NULL DEFAULT 0,
				owner_id BIGINT NOT NULL DEFAULT 0,
				icon_uri VARCHAR(512) NOT NULL DEFAULT '',
				name VARCHAR(255) NOT NULL DEFAULT '',
				description TEXT,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT
			);
		END IF;
	END $$;

	-- 创建 app_release_record 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'app_release_record') THEN
			CREATE TABLE app_release_record (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				app_id BIGINT NOT NULL DEFAULT 0,
				space_id BIGINT NOT NULL DEFAULT 0,
				owner_id BIGINT NOT NULL DEFAULT 0,
				icon_uri VARCHAR(512) NOT NULL DEFAULT '',
				name VARCHAR(255) NOT NULL DEFAULT '',
				description TEXT,
				connector_ids JSONB,
				extra_info JSONB,
				version VARCHAR(255) NOT NULL DEFAULT '',
				version_desc TEXT,
				publish_status SMALLINT NOT NULL DEFAULT 0,
				publish_at BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_app_publish_at ON app_release_record(app_id, publish_at);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_app_version_connector ON app_release_record(app_id, version);
		END IF;
	END $$;

	-- 创建 app_connector_release_ref 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'app_connector_release_ref') THEN
			CREATE TABLE app_connector_release_ref (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				record_id BIGINT NOT NULL DEFAULT 0,
				connector_id BIGINT,
				publish_config JSONB,
				publish_status SMALLINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0
			);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_record_connector ON app_connector_release_ref(record_id, connector_id);
		END IF;
	END $$;

	-- 创建 app_conversation_template_draft 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'app_conversation_template_draft') THEN
			CREATE TABLE app_conversation_template_draft (
				id BIGINT NOT NULL PRIMARY KEY,
				app_id BIGINT NOT NULL,
				space_id BIGINT NOT NULL,
				name VARCHAR(256) NOT NULL,
				template_id BIGINT NOT NULL,
				creator_id BIGINT NOT NULL,
				created_at BIGINT NOT NULL,
				updated_at BIGINT,
				deleted_at BIGINT
			);
			CREATE INDEX IF NOT EXISTS idx_space_id_app_id_template_id ON app_conversation_template_draft(space_id, app_id, template_id);
		END IF;
	END $$;

	-- 创建 app_conversation_template_online 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'app_conversation_template_online') THEN
			CREATE TABLE app_conversation_template_online (
				id BIGINT NOT NULL PRIMARY KEY,
				app_id BIGINT NOT NULL,
				space_id BIGINT NOT NULL,
				name VARCHAR(256) NOT NULL,
				template_id BIGINT NOT NULL,
				version VARCHAR(256) NOT NULL,
				creator_id BIGINT NOT NULL,
				created_at BIGINT NOT NULL
			);
			CREATE INDEX IF NOT EXISTS idx_space_id_app_id_template_id_version ON app_conversation_template_online(space_id, app_id, template_id, version);
		END IF;
	END $$;

	-- 创建 app_static_conversation_draft 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'app_static_conversation_draft') THEN
			CREATE TABLE app_static_conversation_draft (
				id BIGINT NOT NULL PRIMARY KEY,
				template_id BIGINT NOT NULL,
				user_id BIGINT NOT NULL,
				connector_id BIGINT NOT NULL,
				conversation_id BIGINT NOT NULL,
				created_at BIGINT NOT NULL,
				deleted_at BIGINT
			);
			CREATE INDEX IF NOT EXISTS idx_connector_id_user_id_template_id ON app_static_conversation_draft(connector_id, user_id, template_id);
		END IF;
	END $$;

	-- 创建 app_static_conversation_online 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'app_static_conversation_online') THEN
			CREATE TABLE app_static_conversation_online (
				id BIGINT NOT NULL PRIMARY KEY,
				template_id BIGINT NOT NULL,
				user_id BIGINT NOT NULL,
				connector_id BIGINT NOT NULL,
				conversation_id BIGINT NOT NULL,
				created_at BIGINT NOT NULL
			);
			CREATE INDEX IF NOT EXISTS idx_connector_id_user_id_template_id ON app_static_conversation_online(connector_id, user_id, template_id);
		END IF;
	END $$;

	-- 创建 app_dynamic_conversation_draft 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'app_dynamic_conversation_draft') THEN
			CREATE TABLE app_dynamic_conversation_draft (
				id BIGINT NOT NULL PRIMARY KEY,
				app_id BIGINT NOT NULL,
				name VARCHAR(256) NOT NULL,
				user_id BIGINT NOT NULL,
				connector_id BIGINT NOT NULL,
				conversation_id BIGINT NOT NULL,
				created_at BIGINT NOT NULL,
				deleted_at BIGINT
			);
			CREATE INDEX IF NOT EXISTS idx_app_id_connector_id_user_id ON app_dynamic_conversation_draft(app_id, connector_id, user_id);
			CREATE INDEX IF NOT EXISTS idx_connector_id_user_id_name ON app_dynamic_conversation_draft(connector_id, user_id, name);
		END IF;
	END $$;

	-- 创建 app_dynamic_conversation_online 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'app_dynamic_conversation_online') THEN
			CREATE TABLE app_dynamic_conversation_online (
				id BIGINT NOT NULL PRIMARY KEY,
				app_id BIGINT NOT NULL,
				name VARCHAR(256) NOT NULL,
				user_id BIGINT NOT NULL,
				connector_id BIGINT NOT NULL,
				conversation_id BIGINT NOT NULL,
				created_at BIGINT NOT NULL,
				deleted_at BIGINT
			);
			CREATE INDEX IF NOT EXISTS idx_app_id_connector_id_user_id ON app_dynamic_conversation_online(app_id, connector_id, user_id);
			CREATE INDEX IF NOT EXISTS idx_connector_id_user_id_name ON app_dynamic_conversation_online(connector_id, user_id, name);
		END IF;
	END $$;
	`
	return db.Exec(sql).Error
}

// createCozeStudioWorkflowTables 创建 Workflow 相关表
func createCozeStudioWorkflowTables(db *gorm.DB) error {
	sql := `
	-- 创建 workflow_draft 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'workflow_draft') THEN
			CREATE TABLE workflow_draft (
				id BIGINT NOT NULL PRIMARY KEY,
				workflow_id BIGINT NOT NULL,
				canvas TEXT,
				input_params TEXT,
				output_params TEXT,
				commit_id VARCHAR(255) NOT NULL,
				created_at BIGINT NOT NULL,
				updated_at BIGINT
			);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_workflow_id ON workflow_draft(workflow_id);
		END IF;
	END $$;

	-- 创建 workflow_meta 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'workflow_meta') THEN
			CREATE TABLE workflow_meta (
				id BIGINT NOT NULL PRIMARY KEY,
				workflow_id BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				description TEXT,
				icon_uri VARCHAR(512),
				status SMALLINT NOT NULL DEFAULT 0,
				content_type SMALLINT NOT NULL DEFAULT 0,
				author_id BIGINT NOT NULL,
				space_id BIGINT NOT NULL,
				updater_id BIGINT,
				source_id BIGINT,
				app_id BIGINT,
				created_at BIGINT NOT NULL,
				updated_at BIGINT
			);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_workflow_id ON workflow_meta(workflow_id);
		END IF;
	END $$;

	-- 创建 workflow_reference 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'workflow_reference') THEN
			CREATE TABLE workflow_reference (
				id BIGINT NOT NULL PRIMARY KEY,
				source_workflow_id BIGINT NOT NULL,
				target_workflow_id BIGINT NOT NULL,
				created_at BIGINT NOT NULL,
				deleted_at BIGINT
			);
		END IF;
	END $$;

	-- 创建 workflow_snapshot 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'workflow_snapshot') THEN
			CREATE TABLE workflow_snapshot (
				id BIGSERIAL PRIMARY KEY,
				workflow_id BIGINT NOT NULL,
				commit_id VARCHAR(255) NOT NULL,
				canvas TEXT,
				input_params TEXT,
				output_params TEXT,
				created_at BIGINT NOT NULL
			);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_workflow_id_commit_id ON workflow_snapshot(workflow_id, commit_id);
		END IF;
	END $$;

	-- 创建 connector_workflow_version 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'connector_workflow_version') THEN
			CREATE TABLE connector_workflow_version (
				id BIGSERIAL PRIMARY KEY,
				app_id BIGINT NOT NULL,
				connector_id BIGINT NOT NULL,
				workflow_id BIGINT NOT NULL,
				version VARCHAR(256) NOT NULL,
				created_at BIGINT NOT NULL
			);
			CREATE INDEX IF NOT EXISTS idx_connector_id_workflow_id_create_at ON connector_workflow_version(connector_id, workflow_id, created_at);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_connector_id_workflow_id_version ON connector_workflow_version(connector_id, workflow_id, version);
		END IF;
	END $$;

	-- 创建 node_execution 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'node_execution') THEN
			CREATE TABLE node_execution (
				id BIGINT NOT NULL PRIMARY KEY,
				execute_id BIGINT NOT NULL,
				node_id VARCHAR(128) NOT NULL,
				node_name VARCHAR(128) NOT NULL,
				node_type VARCHAR(128) NOT NULL,
				created_at BIGINT NOT NULL,
				status SMALLINT NOT NULL,
				duration BIGINT,
				input TEXT,
				output TEXT,
				raw_output TEXT,
				error_info TEXT,
				error_level VARCHAR(32),
				input_tokens BIGINT,
				output_tokens BIGINT,
				updated_at BIGINT,
				composite_node_index BIGINT,
				composite_node_items TEXT,
				parent_node_id VARCHAR(128),
				sub_execute_id BIGINT,
				extra TEXT
			);
			CREATE INDEX IF NOT EXISTS idx_execute_id_node_id ON node_execution(execute_id, node_id);
			CREATE INDEX IF NOT EXISTS idx_execute_id_parent_node_id ON node_execution(execute_id, parent_node_id);
		END IF;
	END $$;

	-- 创建 chat_flow_role_config 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'chat_flow_role_config') THEN
			CREATE TABLE chat_flow_role_config (
				id BIGINT NOT NULL PRIMARY KEY,
				workflow_id BIGINT NOT NULL,
				connector_id BIGINT,
				name VARCHAR(256) NOT NULL,
				description TEXT,
				version VARCHAR(256),
				avatar VARCHAR(256) NOT NULL,
				background_image_info TEXT,
				onboarding_info TEXT,
				suggest_reply_info TEXT,
				audio_config TEXT,
				user_input_config VARCHAR(256) NOT NULL,
				creator_id BIGINT NOT NULL,
				created_at BIGINT NOT NULL,
				updated_at BIGINT,
				deleted_at BIGINT
			);
			CREATE INDEX IF NOT EXISTS idx_connector_id_version ON chat_flow_role_config(connector_id, version);
			CREATE INDEX IF NOT EXISTS idx_workflow_id_version ON chat_flow_role_config(workflow_id, version);
		END IF;
	END $$;
	`
	return db.Exec(sql).Error
}

// createCozeStudioKnowledgeTables 创建 Knowledge 相关表
func createCozeStudioKnowledgeTables(db *gorm.DB) error {
	sql := `
	-- 创建 knowledge 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'knowledge') THEN
			CREATE TABLE knowledge (
				id BIGINT NOT NULL PRIMARY KEY,
				name VARCHAR(150) NOT NULL DEFAULT '',
				app_id BIGINT NOT NULL DEFAULT 0,
				creator_id BIGINT NOT NULL DEFAULT 0,
				space_id BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT,
				status SMALLINT NOT NULL DEFAULT 1,
				description TEXT,
				icon_uri VARCHAR(150),
				format_type SMALLINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_app_id ON knowledge(app_id);
			CREATE INDEX IF NOT EXISTS idx_creator_id ON knowledge(creator_id);
			CREATE INDEX IF NOT EXISTS idx_space_id_deleted_at_updated_at ON knowledge(space_id, deleted_at, updated_at);
		END IF;
	END $$;

	-- 创建 knowledge_document 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'knowledge_document') THEN
			CREATE TABLE knowledge_document (
				id BIGINT NOT NULL PRIMARY KEY,
				knowledge_id BIGINT NOT NULL DEFAULT 0,
				name VARCHAR(150) NOT NULL DEFAULT '',
				file_extension VARCHAR(20) NOT NULL DEFAULT '0',
				document_type INTEGER NOT NULL DEFAULT 0,
				uri TEXT,
				size BIGINT NOT NULL DEFAULT 0,
				slice_count BIGINT NOT NULL DEFAULT 0,
				char_count BIGINT NOT NULL DEFAULT 0,
				creator_id BIGINT NOT NULL DEFAULT 0,
				space_id BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT,
				source_type INTEGER DEFAULT 0,
				status INTEGER NOT NULL DEFAULT 0,
				fail_reason TEXT,
				parse_rule JSONB,
				table_info JSONB
			);
			CREATE INDEX IF NOT EXISTS idx_creator_id ON knowledge_document(creator_id);
			CREATE INDEX IF NOT EXISTS idx_knowledge_id_deleted_at_updated_at ON knowledge_document(knowledge_id, deleted_at, updated_at);
		END IF;
	END $$;

	-- 创建 knowledge_document_slice 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'knowledge_document_slice') THEN
			CREATE TABLE knowledge_document_slice (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				knowledge_id BIGINT NOT NULL DEFAULT 0,
				document_id BIGINT NOT NULL DEFAULT 0,
				content TEXT,
				sequence NUMERIC(20,5) NOT NULL,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT,
				creator_id BIGINT NOT NULL DEFAULT 0,
				space_id BIGINT NOT NULL DEFAULT 0,
				status INTEGER NOT NULL DEFAULT 0,
				fail_reason TEXT,
				hit BIGINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_document_id_deleted_at_sequence ON knowledge_document_slice(document_id, deleted_at, sequence);
			CREATE INDEX IF NOT EXISTS idx_knowledge_id_document_id ON knowledge_document_slice(knowledge_id, document_id);
			CREATE INDEX IF NOT EXISTS idx_sequence ON knowledge_document_slice(sequence);
		END IF;
	END $$;

	-- 创建 knowledge_document_review 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'knowledge_document_review') THEN
			CREATE TABLE knowledge_document_review (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				knowledge_id BIGINT NOT NULL DEFAULT 0,
				space_id BIGINT NOT NULL DEFAULT 0,
				name VARCHAR(150) NOT NULL DEFAULT '',
				type VARCHAR(10) NOT NULL DEFAULT '0',
				uri TEXT,
				format_type SMALLINT NOT NULL DEFAULT 0,
				status SMALLINT NOT NULL DEFAULT 0,
				chunk_resp_uri TEXT,
				deleted_at BIGINT,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				creator_id BIGINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_dataset_id ON knowledge_document_review(knowledge_id, status, updated_at);
		END IF;
	END $$;
	`
	return db.Exec(sql).Error
}

// createCozeStudioPluginTables 创建 Plugin 相关表
func createCozeStudioPluginTables(db *gorm.DB) error {
	sql := `
	-- 创建 plugin_draft 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'plugin_draft') THEN
			CREATE TABLE plugin_draft (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				space_id BIGINT NOT NULL DEFAULT 0,
				developer_id BIGINT NOT NULL DEFAULT 0,
				app_id BIGINT NOT NULL DEFAULT 0,
				icon_uri VARCHAR(512) NOT NULL DEFAULT '',
				server_url VARCHAR(512) NOT NULL DEFAULT '',
				plugin_type SMALLINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT,
				manifest JSONB,
				openapi_doc JSONB
			);
			CREATE INDEX IF NOT EXISTS idx_app_id ON plugin_draft(app_id, id);
			CREATE INDEX IF NOT EXISTS idx_space_app_created_at ON plugin_draft(space_id, app_id, created_at);
			CREATE INDEX IF NOT EXISTS idx_space_app_updated_at ON plugin_draft(space_id, app_id, updated_at);
		END IF;
	END $$;

	-- 创建 plugin_version 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'plugin_version') THEN
			CREATE TABLE plugin_version (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				space_id BIGINT NOT NULL DEFAULT 0,
				developer_id BIGINT NOT NULL DEFAULT 0,
				plugin_id BIGINT NOT NULL DEFAULT 0,
				app_id BIGINT NOT NULL DEFAULT 0,
				icon_uri VARCHAR(512) NOT NULL DEFAULT '',
				server_url VARCHAR(512) NOT NULL DEFAULT '',
				plugin_type SMALLINT NOT NULL DEFAULT 0,
				version VARCHAR(255) NOT NULL DEFAULT '',
				version_desc TEXT,
				manifest JSONB,
				openapi_doc JSONB,
				created_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT
			);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_plugin_version ON plugin_version(plugin_id, version);
		END IF;
	END $$;

	-- 创建 plugin_oauth_auth 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'plugin_oauth_auth') THEN
			CREATE TABLE plugin_oauth_auth (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				user_id VARCHAR(255) NOT NULL DEFAULT '',
				plugin_id BIGINT NOT NULL DEFAULT 0,
				is_draft BOOLEAN NOT NULL DEFAULT false,
				oauth_config JSONB,
				access_token TEXT,
				refresh_token TEXT,
				token_expired_at BIGINT,
				next_token_refresh_at BIGINT,
				last_active_at BIGINT,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_last_active_at ON plugin_oauth_auth(last_active_at);
			CREATE INDEX IF NOT EXISTS idx_last_token_expired_at ON plugin_oauth_auth(token_expired_at);
			CREATE INDEX IF NOT EXISTS idx_next_token_refresh_at ON plugin_oauth_auth(next_token_refresh_at);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_user_plugin_is_draft ON plugin_oauth_auth(user_id, plugin_id, is_draft);
		END IF;
	END $$;

	-- 创建 tool 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tool') THEN
			CREATE TABLE tool (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				plugin_id BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				version VARCHAR(255) NOT NULL DEFAULT '',
				sub_url VARCHAR(512) NOT NULL DEFAULT '',
				method VARCHAR(64) NOT NULL DEFAULT '',
				operation JSONB,
				activated_status SMALLINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_plugin_activated_status ON tool(plugin_id, activated_status);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_plugin_sub_url_method ON tool(plugin_id, sub_url, method);
		END IF;
	END $$;

	-- 创建 tool_draft 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tool_draft') THEN
			CREATE TABLE tool_draft (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				plugin_id BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				sub_url VARCHAR(512) NOT NULL DEFAULT '',
				method VARCHAR(64) NOT NULL DEFAULT '',
				operation JSONB,
				debug_status SMALLINT NOT NULL DEFAULT 0,
				activated_status SMALLINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_plugin_created_at_id ON tool_draft(plugin_id, created_at, id);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_plugin_sub_url_method ON tool_draft(plugin_id, sub_url, method);
		END IF;
	END $$;

	-- 创建 tool_version 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tool_version') THEN
			CREATE TABLE tool_version (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				tool_id BIGINT NOT NULL DEFAULT 0,
				plugin_id BIGINT NOT NULL DEFAULT 0,
				version VARCHAR(255) NOT NULL DEFAULT '',
				sub_url VARCHAR(512) NOT NULL DEFAULT '',
				method VARCHAR(64) NOT NULL DEFAULT '',
				operation JSONB,
				created_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT
			);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_tool_version ON tool_version(tool_id, version);
		END IF;
	END $$;
	`
	return db.Exec(sql).Error
}

// createCozeStudioConversationTables 创建 Conversation 相关表
func createCozeStudioConversationTables(db *gorm.DB) error {
	sql := `
	-- 创建 message 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'message') THEN
			CREATE TABLE message (
				id BIGSERIAL PRIMARY KEY,
				run_id BIGINT NOT NULL DEFAULT 0,
				conversation_id BIGINT NOT NULL DEFAULT 0,
				user_id VARCHAR(60) NOT NULL DEFAULT '',
				agent_id BIGINT NOT NULL DEFAULT 0,
				role VARCHAR(100) NOT NULL DEFAULT '',
				content_type VARCHAR(100) NOT NULL DEFAULT '',
				content TEXT,
				message_type VARCHAR(100) NOT NULL DEFAULT '',
				display_content TEXT,
				ext TEXT,
				section_id BIGINT,
				broken_position INTEGER DEFAULT -1,
				status SMALLINT NOT NULL DEFAULT 0,
				model_content TEXT,
				meta_info TEXT,
				reasoning_content TEXT,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_conversation_id ON message(conversation_id);
			CREATE INDEX IF NOT EXISTS idx_run_id ON message(run_id);
		END IF;
	END $$;

	-- 创建 run_record 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'run_record') THEN
			CREATE TABLE run_record (
				id BIGINT NOT NULL PRIMARY KEY,
				conversation_id BIGINT NOT NULL DEFAULT 0,
				section_id BIGINT NOT NULL DEFAULT 0,
				agent_id BIGINT NOT NULL DEFAULT 0,
				user_id VARCHAR(255) NOT NULL DEFAULT '',
				source SMALLINT NOT NULL DEFAULT 0,
				status VARCHAR(255) NOT NULL DEFAULT '',
				creator_id BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				failed_at BIGINT NOT NULL DEFAULT 0,
				last_error TEXT,
				completed_at BIGINT NOT NULL DEFAULT 0,
				chat_request TEXT,
				ext TEXT,
				usage JSONB
			);
			CREATE INDEX IF NOT EXISTS idx_c_s ON run_record(conversation_id, section_id);
		END IF;
	END $$;
	`
	return db.Exec(sql).Error
}

// createCozeStudioMemoryTables 创建 Memory 相关表
func createCozeStudioMemoryTables(db *gorm.DB) error {
	sql := `
	-- 创建 variables_meta 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'variables_meta') THEN
			CREATE TABLE variables_meta (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				creator_id BIGINT NOT NULL,
				biz_type SMALLINT NOT NULL,
				biz_id VARCHAR(128) NOT NULL DEFAULT '',
				variable_list JSONB,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				version VARCHAR(255) NOT NULL
			);
			CREATE INDEX IF NOT EXISTS idx_user_key ON variables_meta(creator_id);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_project_key ON variables_meta(biz_id, biz_type, version);
		END IF;
	END $$;

	-- 创建 variable_instance 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'variable_instance') THEN
			CREATE TABLE variable_instance (
				id BIGINT NOT NULL DEFAULT 0 PRIMARY KEY,
				biz_type SMALLINT NOT NULL,
				biz_id VARCHAR(128) NOT NULL DEFAULT '',
				version VARCHAR(255) NOT NULL,
				keyword VARCHAR(255) NOT NULL,
				type SMALLINT NOT NULL,
				content TEXT,
				connector_uid VARCHAR(255) NOT NULL,
				connector_id BIGINT NOT NULL,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_connector_key ON variable_instance(biz_id, biz_type, version, connector_uid, connector_id);
		END IF;
	END $$;

	-- 创建 draft_database_info 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'draft_database_info') THEN
			CREATE TABLE draft_database_info (
				id BIGINT NOT NULL PRIMARY KEY,
				app_id BIGINT,
				space_id BIGINT NOT NULL,
				related_online_id BIGINT NOT NULL,
				is_visible SMALLINT NOT NULL DEFAULT 1,
				prompt_disabled SMALLINT NOT NULL DEFAULT 0,
				table_name VARCHAR(255) NOT NULL,
				table_desc VARCHAR(256),
				table_field TEXT,
				creator_id BIGINT NOT NULL DEFAULT 0,
				icon_uri VARCHAR(255) NOT NULL,
				physical_table_name VARCHAR(255),
				rw_mode BIGINT NOT NULL DEFAULT 1,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT
			);
			CREATE INDEX IF NOT EXISTS idx_space_app_creator_deleted ON draft_database_info(space_id, app_id, creator_id, deleted_at);
		END IF;
	END $$;

	-- 创建 online_database_info 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'online_database_info') THEN
			CREATE TABLE online_database_info (
				id BIGINT NOT NULL PRIMARY KEY,
				app_id BIGINT,
				space_id BIGINT NOT NULL,
				related_draft_id BIGINT NOT NULL,
				is_visible SMALLINT NOT NULL DEFAULT 1,
				prompt_disabled SMALLINT NOT NULL DEFAULT 0,
				table_name VARCHAR(255) NOT NULL,
				table_desc VARCHAR(256),
				table_field TEXT,
				creator_id BIGINT NOT NULL DEFAULT 0,
				icon_uri VARCHAR(255) NOT NULL,
				physical_table_name VARCHAR(255),
				rw_mode BIGINT NOT NULL DEFAULT 1,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT
			);
			CREATE INDEX IF NOT EXISTS idx_space_app_creator_deleted ON online_database_info(space_id, app_id, creator_id, deleted_at);
		END IF;
	END $$;

	-- 创建 agent_to_database 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'agent_to_database') THEN
			CREATE TABLE agent_to_database (
				id BIGINT NOT NULL PRIMARY KEY,
				agent_id BIGINT NOT NULL,
				database_id BIGINT NOT NULL,
				is_draft BOOLEAN NOT NULL,
				prompt_disable BOOLEAN NOT NULL DEFAULT false
			);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_agent_db_draft ON agent_to_database(agent_id, database_id, is_draft);
		END IF;
	END $$;
	`
	return db.Exec(sql).Error
}

// createCozeStudioOtherTables 创建其他表
func createCozeStudioOtherTables(db *gorm.DB) error {
	sql := `
	-- 创建 files 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'files') THEN
			CREATE TABLE files (
				id BIGINT NOT NULL PRIMARY KEY,
				name VARCHAR(255) NOT NULL DEFAULT '',
				file_size BIGINT NOT NULL DEFAULT 0,
				tos_uri VARCHAR(1024) NOT NULL DEFAULT '',
				status SMALLINT NOT NULL DEFAULT 0,
				comment VARCHAR(1024) NOT NULL DEFAULT '',
				source SMALLINT NOT NULL DEFAULT 0,
				creator_id VARCHAR(512) NOT NULL DEFAULT '',
				content_type VARCHAR(255) NOT NULL DEFAULT '',
				coze_account_id BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				deleted_at BIGINT
			);
			CREATE INDEX IF NOT EXISTS idx_creator_id ON files(creator_id);
		END IF;
	END $$;

	-- 创建 prompt_resource 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'prompt_resource') THEN
			CREATE TABLE prompt_resource (
				id BIGSERIAL PRIMARY KEY,
				space_id BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				description VARCHAR(255) NOT NULL,
				prompt_text TEXT,
				status INTEGER NOT NULL,
				creator_id BIGINT NOT NULL,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_creator_id ON prompt_resource(creator_id);
		END IF;
	END $$;

	-- 创建 shortcut_command 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'shortcut_command') THEN
			CREATE TABLE shortcut_command (
				id BIGSERIAL PRIMARY KEY,
				object_id BIGINT NOT NULL DEFAULT 0,
				command_id BIGINT NOT NULL DEFAULT 0,
				command_name VARCHAR(255) NOT NULL DEFAULT '',
				shortcut_command VARCHAR(255) NOT NULL DEFAULT '',
				description VARCHAR(2000) NOT NULL DEFAULT '',
				send_type SMALLINT NOT NULL DEFAULT 0,
				tool_type SMALLINT NOT NULL DEFAULT 0,
				work_flow_id BIGINT NOT NULL DEFAULT 0,
				plugin_id BIGINT NOT NULL DEFAULT 0,
				plugin_tool_name VARCHAR(255) NOT NULL DEFAULT '',
				template_query TEXT,
				components JSONB,
				card_schema TEXT,
				tool_info JSONB,
				status SMALLINT NOT NULL DEFAULT 0,
				creator_id BIGINT DEFAULT 0,
				is_online SMALLINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				agent_id BIGINT NOT NULL DEFAULT 0,
				shortcut_icon JSONB,
				plugin_tool_id BIGINT NOT NULL DEFAULT 0,
				source SMALLINT DEFAULT 0
			);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_object_command_id_type ON shortcut_command(object_id, command_id, is_online);
		END IF;
	END $$;

	-- 创建 data_copy_task 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'data_copy_task') THEN
			CREATE TABLE data_copy_task (
				id BIGSERIAL PRIMARY KEY,
				master_task_id VARCHAR(128) DEFAULT '',
				origin_data_id BIGINT NOT NULL DEFAULT 0,
				target_data_id BIGINT NOT NULL DEFAULT 0,
				origin_space_id BIGINT NOT NULL DEFAULT 0,
				target_space_id BIGINT NOT NULL DEFAULT 0,
				origin_user_id BIGINT NOT NULL DEFAULT 0,
				target_user_id BIGINT DEFAULT 0,
				origin_app_id BIGINT NOT NULL DEFAULT 0,
				target_app_id BIGINT NOT NULL DEFAULT 0,
				data_type SMALLINT NOT NULL DEFAULT 0,
				ext_info VARCHAR(255) NOT NULL DEFAULT '',
				start_time BIGINT DEFAULT 0,
				finish_time BIGINT,
				status SMALLINT NOT NULL DEFAULT 1,
				error_msg VARCHAR(128)
			);
			CREATE UNIQUE INDEX IF NOT EXISTS uniq_master_task_id_origin_data_id_data_type ON data_copy_task(master_task_id, origin_data_id, data_type);
		END IF;
	END $$;

	-- 创建 api_key 表
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'api_key') THEN
			CREATE TABLE api_key (
				id BIGSERIAL PRIMARY KEY,
				api_key VARCHAR(255) NOT NULL DEFAULT '',
				name VARCHAR(255) NOT NULL DEFAULT '',
				status SMALLINT NOT NULL DEFAULT 0,
				user_id BIGINT NOT NULL DEFAULT 0,
				expired_at BIGINT NOT NULL DEFAULT 0,
				created_at BIGINT NOT NULL DEFAULT 0,
				updated_at BIGINT NOT NULL DEFAULT 0,
				last_used_at BIGINT NOT NULL DEFAULT 0,
				ak_type SMALLINT NOT NULL DEFAULT 0
			);
		END IF;
	END $$;

	-- 创建 kv_entries 表（用于存储键值对配置）
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'kv_entries') THEN
			CREATE TABLE kv_entries (
				id BIGSERIAL PRIMARY KEY,
				namespace VARCHAR(255) NOT NULL,
				key_data VARCHAR(255) NOT NULL,
				value_data BYTEA NOT NULL,
				CONSTRAINT uniq_namespace_key UNIQUE (namespace, key_data)
			);
			CREATE INDEX IF NOT EXISTS idx_kv_namespace_key ON kv_entries(namespace, key_data);
		END IF;
	END $$;
	`
	return db.Exec(sql).Error
}

