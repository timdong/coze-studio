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

package container

import (
	"github.com/coze-dev/coze-studio/backend/core/config"
	"github.com/coze-dev/coze-studio/backend/core/logging"
	"etrxlite/platform/ent"
	"sync"
)

// Container 依赖注入容器
type Container struct {
	mu        sync.RWMutex
	logger    *logging.Logger
	config    *config.Config
	entClient *ent.Client
	services  map[string]interface{}
}

var (
	globalContainer *Container
	once            sync.Once
)

// GetContainer 获取全局容器
func GetContainer() *Container {
	once.Do(func() {
		globalContainer = &Container{
			services: make(map[string]interface{}),
		}
	})
	return globalContainer
}

// SetLogger 设置日志器
func (c *Container) SetLogger(logger *logging.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
}

// GetLogger 获取日志器
func (c *Container) GetLogger() *logging.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.logger == nil {
		// 如果未设置，创建默认日志器
		defaultLogger, _ := logging.NewDefaultLogger()
		c.logger = defaultLogger
	}
	return c.logger
}

// SetConfig 设置配置
func (c *Container) SetConfig(cfg *config.Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config = cfg
}

// GetConfig 获取配置
func (c *Container) GetConfig() *config.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// SetEntClient 设置Ent客户端
func (c *Container) SetEntClient(client *ent.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entClient = client
}

// GetEntClient 获取Ent客户端
func (c *Container) GetEntClient() *ent.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.entClient
}

// Register 注册服务
func (c *Container) Register(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// Get 获取服务
func (c *Container) Get(name string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	service, exists := c.services[name]
	return service, exists
}

// MustGet 获取服务（如果不存在则panic）
func (c *Container) MustGet(name string) interface{} {
	service, exists := c.Get(name)
	if !exists {
		panic("service not found: " + name)
	}
	return service
}
