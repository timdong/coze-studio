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
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/core/logging"

	"go.uber.org/zap"
)

// Mem0CacheItem Mem0缓存项
type Mem0CacheItem struct {
	Memories  interface{} `json:"memories"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// Mem0CacheService Mem0缓存服务
type Mem0CacheService struct {
	manager *Manager
	prefix  string
}

// NewMem0CacheService 创建Mem0缓存服务
func NewMem0CacheService(manager *Manager) *Mem0CacheService {
	return &Mem0CacheService{
		manager: manager,
		prefix:  "mem0:",
	}
}

// GetSearchResult 获取搜索记忆结果缓存
func (m *Mem0CacheService) GetSearchResult(query, userID, agentID string, limit int) (interface{}, bool) {
	key := m.getSearchKey(query, userID, agentID, limit)

	// 从缓存获取
	value, err := m.manager.Get(context.Background(), key)
	if err == nil {
		var item Mem0CacheItem
		if err := json.Unmarshal([]byte(value), &item); err == nil {
			if time.Now().Before(item.ExpiresAt) {
				logging.Log.Debug("从缓存获取Mem0搜索结果",
					zap.String("query", query),
					zap.String("user_id", userID),
				)
				return item.Memories, true
			}
		}
	}

	return nil, false
}

// SetSearchResult 设置搜索记忆结果缓存
func (m *Mem0CacheService) SetSearchResult(query, userID, agentID string, limit int, memories interface{}, ttl time.Duration) error {
	key := m.getSearchKey(query, userID, agentID, limit)

	item := Mem0CacheItem{
		Memories:  memories,
		ExpiresAt: time.Now().Add(ttl),
	}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal cache item: %w", err)
	}

	// 设置缓存
	err = m.manager.Set(context.Background(), key, string(data), ttl)
	if err != nil {
		logging.Log.Warn("设置Mem0搜索结果缓存失败",
			zap.String("query", query),
			zap.Error(err),
		)
		return err
	}

	logging.Log.Debug("设置Mem0搜索结果缓存",
		zap.String("query", query),
		zap.String("user_id", userID),
		zap.Duration("ttl", ttl),
	)
	return nil
}

// GetAllMemoriesResult 获取所有记忆结果缓存
func (m *Mem0CacheService) GetAllMemoriesResult(userID, agentID, runID string, limit, offset int) (interface{}, bool) {
	key := m.getAllMemoriesKey(userID, agentID, runID, limit, offset)

	// 从缓存获取
	value, err := m.manager.Get(context.Background(), key)
	if err == nil {
		var item Mem0CacheItem
		if err := json.Unmarshal([]byte(value), &item); err == nil {
			if time.Now().Before(item.ExpiresAt) {
				logging.Log.Debug("从缓存获取Mem0所有记忆",
					zap.String("user_id", userID),
					zap.String("agent_id", agentID),
				)
				return item.Memories, true
			}
		}
	}

	return nil, false
}

// SetAllMemoriesResult 设置所有记忆结果缓存
func (m *Mem0CacheService) SetAllMemoriesResult(userID, agentID, runID string, limit, offset int, memories interface{}, ttl time.Duration) error {
	key := m.getAllMemoriesKey(userID, agentID, runID, limit, offset)

	item := Mem0CacheItem{
		Memories:  memories,
		ExpiresAt: time.Now().Add(ttl),
	}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal cache item: %w", err)
	}

	// 设置缓存
	err = m.manager.Set(context.Background(), key, string(data), ttl)
	if err != nil {
		logging.Log.Warn("设置Mem0所有记忆缓存失败",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return err
	}

	logging.Log.Debug("设置Mem0所有记忆缓存",
		zap.String("user_id", userID),
		zap.String("agent_id", agentID),
		zap.Duration("ttl", ttl),
	)
	return nil
}

// InvalidateUserCache 失效用户相关缓存
func (m *Mem0CacheService) InvalidateUserCache(userID string) error {
	pattern := m.prefix + "user:" + userID + ":*"
	m.manager.ClearMemoryPattern(pattern)
	logging.Log.Info("清除Mem0用户缓存", zap.String("user_id", userID))
	return nil
}

// InvalidateAgentCache 失效Agent相关缓存
func (m *Mem0CacheService) InvalidateAgentCache(agentID string) error {
	pattern := m.prefix + "agent:" + agentID + ":*"
	m.manager.ClearMemoryPattern(pattern)
	logging.Log.Info("清除Mem0 Agent缓存", zap.String("agent_id", agentID))
	return nil
}

// InvalidateAllCache 清除所有Mem0缓存
func (m *Mem0CacheService) InvalidateAllCache() error {
	pattern := m.prefix + "*"
	m.manager.ClearMemoryPattern(pattern)
	logging.Log.Info("清除所有Mem0缓存")
	return nil
}

// getSearchKey 获取搜索缓存键
func (m *Mem0CacheService) getSearchKey(query, userID, agentID string, limit int) string {
	// 使用hash避免key过长
	hash := m.hashKey(fmt.Sprintf("%s:%s:%s:%d", query, userID, agentID, limit))
	return fmt.Sprintf("%ssearch:%s", m.prefix, hash)
}

// getAllMemoriesKey 获取所有记忆缓存键
func (m *Mem0CacheService) getAllMemoriesKey(userID, agentID, runID string, limit, offset int) string {
	// 使用hash避免key过长
	hash := m.hashKey(fmt.Sprintf("%s:%s:%s:%d:%d", userID, agentID, runID, limit, offset))
	return fmt.Sprintf("%sall:%s", m.prefix, hash)
}

// hashKey 生成hash键
func (m *Mem0CacheService) hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])[:16] // 使用前16个字符
}

// GetStats 获取缓存统计
func (m *Mem0CacheService) GetStats() map[string]interface{} {
	return m.manager.GetStats()
}

