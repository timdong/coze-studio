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
	"encoding/json"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/core/logging"

	"go.uber.org/zap"
)

// RAGCacheItem RAG缓存项
type RAGCacheItem struct {
	Query     string      `json:"query"`
	Results   interface{} `json:"results"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// RAGCacheService RAG缓存服务
type RAGCacheService struct {
	manager *Manager
	prefix  string
}

// NewRAGCacheService 创建RAG缓存服务
func NewRAGCacheService(manager *Manager) *RAGCacheService {
	return &RAGCacheService{
		manager: manager,
		prefix:  "rag:",
	}
}

// GetQueryResult 获取查询结果缓存
func (r *RAGCacheService) GetQueryResult(query string) (interface{}, bool) {
	key := r.getQueryKey(query)

	// 先尝试从统一缓存获取
	value, err := r.manager.Get(context.Background(), key)
	if err == nil {
		var item RAGCacheItem
		if err := json.Unmarshal([]byte(value), &item); err == nil {
			if time.Now().Before(item.ExpiresAt) {
				logging.Log.Debug("从缓存获取RAG查询结果", zap.String("query", query))
				return item.Results, true
			}
		}
	}

	return nil, false
}

// SetQueryResult 设置查询结果缓存
func (r *RAGCacheService) SetQueryResult(query string, results interface{}, ttl time.Duration) error {
	key := r.getQueryKey(query)

	item := RAGCacheItem{
		Query:     query,
		Results:   results,
		ExpiresAt: time.Now().Add(ttl),
	}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal cache item: %w", err)
	}

	// 设置统一缓存
	err = r.manager.Set(context.Background(), key, string(data), ttl)
	if err != nil {
		logging.Log.Warn("设置缓存失败", zap.Error(err))
	}

	logging.Log.Debug("设置RAG查询结果缓存", zap.String("query", query), zap.Duration("ttl", ttl))
	return nil
}

// GetEmbeddingCache 获取嵌入向量缓存
func (r *RAGCacheService) GetEmbeddingCache(text string) ([]float32, bool) {
	key := r.getEmbeddingKey(text)

	// 从统一缓存获取
	value, err := r.manager.Get(context.Background(), key)
	if err == nil {
		var embedding []float32
		if err := json.Unmarshal([]byte(value), &embedding); err == nil {
			logging.Log.Debug("从缓存获取嵌入向量", zap.String("text", text))
			return embedding, true
		}
	}

	return nil, false
}

// SetEmbeddingCache 设置嵌入向量缓存
func (r *RAGCacheService) SetEmbeddingCache(text string, embedding []float32, ttl time.Duration) error {
	key := r.getEmbeddingKey(text)

	data, err := json.Marshal(embedding)
	if err != nil {
		return fmt.Errorf("failed to marshal embedding: %w", err)
	}

	// 设置统一缓存
	err = r.manager.Set(context.Background(), key, string(data), ttl)
	if err != nil {
		logging.Log.Warn("设置嵌入向量缓存失败", zap.Error(err))
	}

	logging.Log.Debug("设置嵌入向量缓存", zap.String("text", text), zap.Duration("ttl", ttl))
	return nil
}

// ClearQueryCache 清除查询缓存
func (r *RAGCacheService) ClearQueryCache() error {
	pattern := r.prefix + "query:*"

	// 清除内存缓存中的查询缓存
	r.manager.ClearMemoryPattern(pattern)

	logging.Log.Info("清除RAG查询缓存")
	return nil
}

// ClearEmbeddingCache 清除嵌入向量缓存
func (r *RAGCacheService) ClearEmbeddingCache() error {
	pattern := r.prefix + "embedding:*"

	// 清除内存缓存中的嵌入向量缓存
	r.manager.ClearMemoryPattern(pattern)

	logging.Log.Info("清除RAG嵌入向量缓存")
	return nil
}

// getQueryKey 获取查询缓存键
func (r *RAGCacheService) getQueryKey(query string) string {
	return r.prefix + "query:" + query
}

// getEmbeddingKey 获取嵌入向量缓存键
func (r *RAGCacheService) getEmbeddingKey(text string) string {
	return r.prefix + "embedding:" + text
}

// GetEmbedding 获取嵌入向量（兼容旧接口）
func (r *RAGCacheService) GetEmbedding(text string) ([]float32, bool) {
	return r.GetEmbeddingCache(text)
}

// SetEmbedding 设置嵌入向量（兼容旧接口）
func (r *RAGCacheService) SetEmbedding(text string, embedding []float32, ttl time.Duration) error {
	return r.SetEmbeddingCache(text, embedding, ttl)
}

// GetChunk 获取分块缓存
func (r *RAGCacheService) GetChunk(key string) (interface{}, bool) {
	cacheKey := r.prefix + "chunk:" + key
	value, err := r.manager.Get(context.Background(), cacheKey)
	if err != nil {
		return nil, false
	}

	var chunk interface{}
	if err := json.Unmarshal([]byte(value), &chunk); err == nil {
		return chunk, true
	}
	return nil, false
}

// SetChunk 设置分块缓存
func (r *RAGCacheService) SetChunk(key string, chunk interface{}, ttl time.Duration) error {
	cacheKey := r.prefix + "chunk:" + key
	data, err := json.Marshal(chunk)
	if err != nil {
		return fmt.Errorf("failed to marshal chunk: %w", err)
	}

	return r.manager.Set(context.Background(), cacheKey, string(data), ttl)
}

// GetStats 获取缓存统计
func (r *RAGCacheService) GetStats() map[string]interface{} {
	return r.manager.GetStats()
}

// Clear 清除所有缓存
func (r *RAGCacheService) Clear() error {
	// 清除查询缓存
	r.ClearQueryCache()
	// 清除嵌入向量缓存
	r.ClearEmbeddingCache()
	return nil
}
