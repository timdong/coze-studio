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

package cache

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// MemoryCache 内存缓存
type MemoryCache struct {
	items  map[string]*CacheItem
	mu     sync.RWMutex
	config *MemoryConfig
}

// CacheItem 缓存项
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(config *MemoryConfig) *MemoryCache {
	cache := &MemoryCache{
		items:  make(map[string]*CacheItem),
		config: config,
	}

	// 启动清理协程
	go cache.cleanup()

	return cache
}

// Set 设置缓存
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiresAt := time.Now().Add(ttl)
	c.items[key] = &CacheItem{
		Value:     value,
		ExpiresAt: expiresAt,
	}

	return nil
}

// Get 获取缓存
func (c *MemoryCache) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return "", fmt.Errorf("key not found")
	}

	// 检查是否过期
	if time.Now().After(item.ExpiresAt) {
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		c.mu.RLock()
		return "", fmt.Errorf("key expired")
	}

	// 转换为字符串
	if str, ok := item.Value.(string); ok {
		return str, nil
	}

	return "", fmt.Errorf("value is not a string")
}

// Delete 删除缓存
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Exists 检查键是否存在
func (c *MemoryCache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return false
	}

	// 检查是否过期
	if time.Now().After(item.ExpiresAt) {
		return false
	}

	return true
}

// GetStats 获取统计信息
func (c *MemoryCache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"size":     len(c.items),
		"max_size": c.config.MaxSize,
	}
}

// cleanup 清理过期项
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(c.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// ClearPattern 清除匹配模式的键
func (c *MemoryCache) ClearPattern(pattern string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.items {
		if strings.Contains(key, pattern) {
			delete(c.items, key)
		}
	}
}
