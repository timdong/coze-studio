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

	"github.com/coze-dev/coze-studio/backend/core/cache"
	"github.com/coze-dev/coze-studio/backend/core/database"
	"github.com/coze-dev/coze-studio/backend/core/logging"
	"etrxlite/platform/ent"
	"etrxlite/tools/crud"
	"etrxlite/tools/search"
)

// ServiceManager 服务管理器接口
// 这个接口定义了服务管理器需要提供的能力
// 现有的 platform/services.Manager 应该实现这个接口
type ServiceManager interface {
	GetDatabase() *database.Manager
	GetCache() *cache.Manager
	GetLogger() *logging.Logger
	GetClient() *ent.Client
	GetCRUD() *crud.EntCRUD
	GetSearch() *search.Engine
}

// BaseService 基础服务接口
type BaseService interface {
	// HealthCheck 健康检查
	HealthCheck(ctx context.Context) error

	// GetLogger 获取日志器
	GetLogger() *logging.Logger

	// GetDatabase 获取数据库管理器
	GetDatabase() *database.Manager

	// GetCache 获取缓存管理器
	GetCache() *cache.Manager
}

// BaseServiceImpl 基础服务实现
type BaseServiceImpl struct {
	manager ServiceManager
	logger  *logging.Logger
}

// NewBaseService 创建基础服务
func NewBaseService(manager ServiceManager) *BaseServiceImpl {
	return &BaseServiceImpl{
		manager: manager,
		logger:  manager.GetLogger(),
	}
}

// HealthCheck 默认健康检查实现
func (s *BaseServiceImpl) HealthCheck(ctx context.Context) error {
	// 检查数据库连接
	if err := s.manager.GetDatabase().HealthCheck(ctx); err != nil {
		return err
	}

	// 检查缓存连接
	if err := s.manager.GetCache().HealthCheck(ctx); err != nil {
		return err
	}

	return nil
}

// GetLogger 获取日志器
func (s *BaseServiceImpl) GetLogger() *logging.Logger {
	return s.logger
}

// GetDatabase 获取数据库管理器
func (s *BaseServiceImpl) GetDatabase() *database.Manager {
	return s.manager.GetDatabase()
}

// GetCache 获取缓存管理器
func (s *BaseServiceImpl) GetCache() *cache.Manager {
	return s.manager.GetCache()
}

// GetManager 获取服务管理器
func (s *BaseServiceImpl) GetManager() ServiceManager {
	return s.manager
}

// GetClient 获取Ent客户端
func (s *BaseServiceImpl) GetClient() *ent.Client {
	return s.manager.GetClient()
}

// GetCRUD 获取CRUD工具
func (s *BaseServiceImpl) GetCRUD() *crud.EntCRUD {
	return s.manager.GetCRUD()
}

// GetSearch 获取搜索工具
func (s *BaseServiceImpl) GetSearch() *search.Engine {
	return s.manager.GetSearch()
}

