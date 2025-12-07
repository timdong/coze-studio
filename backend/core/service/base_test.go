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
	"testing"

	"github.com/coze-dev/coze-studio/backend/core/cache"
	"github.com/coze-dev/coze-studio/backend/core/database"
	"github.com/coze-dev/coze-studio/backend/core/logging"
	"etrxlite/platform/ent"
	"etrxlite/tools/crud"
	"etrxlite/tools/search"
)

// mockServiceManager 模拟服务管理器
type mockServiceManager struct {
	db     *database.Manager
	cache  *cache.Manager
	logger *logging.Logger
	client *ent.Client
	crud   *crud.EntCRUD
	search *search.Engine
}

func (m *mockServiceManager) GetDatabase() *database.Manager {
	return m.db
}

func (m *mockServiceManager) GetCache() *cache.Manager {
	return m.cache
}

func (m *mockServiceManager) GetLogger() *logging.Logger {
	return m.logger
}

func (m *mockServiceManager) GetClient() *ent.Client {
	return m.client
}

func (m *mockServiceManager) GetCRUD() *crud.EntCRUD {
	return m.crud
}

func (m *mockServiceManager) GetSearch() *search.Engine {
	return m.search
}

func TestNewBaseService(t *testing.T) {
	logger, _ := logging.NewDefaultLogger()
	manager := &mockServiceManager{
		logger: logger,
	}

	service := NewBaseService(manager)

	if service == nil {
		t.Fatal("NewBaseService returned nil")
	}

	if service.GetLogger() != logger {
		t.Error("GetLogger() returned incorrect logger")
	}

	if service.GetManager() != manager {
		t.Error("GetManager() returned incorrect manager")
	}
}

func TestBaseServiceImpl_GetLogger(t *testing.T) {
	logger, _ := logging.NewDefaultLogger()
	manager := &mockServiceManager{
		logger: logger,
	}

	service := NewBaseService(manager)

	if service.GetLogger() != logger {
		t.Error("GetLogger() returned incorrect logger")
	}
}

