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
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestMemoryCache(t *testing.T) {
	// 创建内存缓存配置
	config := &MemoryConfig{
		MaxSize:         100,
		DefaultTTL:      5 * time.Minute,
		CleanupInterval: 1 * time.Minute,
	}

	// 创建内存缓存
	cache := NewMemoryCache(config)
	if cache == nil {
		t.Fatal("Memory cache should not be nil")
	}

	// 测试设置和获取
	key := "test-key"
	value := "test-value"
	ttl := time.Minute

	err := cache.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	retrievedValue, err := cache.Get(key)
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}

	if retrievedValue != value {
		t.Fatalf("Expected %s, got %s", value, retrievedValue)
	}

	// 测试删除
	cache.Delete(key)

	// 验证删除
	_, err = cache.Get(key)
	if err == nil {
		t.Fatal("Deleted key should not exist")
	}
}

func TestMemoryCacheTTL(t *testing.T) {
	// 创建内存缓存配置
	config := &MemoryConfig{
		MaxSize:         100,
		DefaultTTL:      100 * time.Millisecond,
		CleanupInterval: 50 * time.Millisecond,
	}

	// 创建内存缓存
	cache := NewMemoryCache(config)
	if cache == nil {
		t.Fatal("Memory cache should not be nil")
	}

	// 设置短期缓存
	key := "ttl-test-key"
	value := "ttl-test-value"
	ttl := 50 * time.Millisecond

	err := cache.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// 立即获取应该成功
	retrievedValue, err := cache.Get(key)
	if err != nil {
		t.Fatalf("Failed to get cache immediately: %v", err)
	}

	if retrievedValue != value {
		t.Fatalf("Expected %s, got %s", value, retrievedValue)
	}

	// 等待过期
	time.Sleep(100 * time.Millisecond)

	// 过期后获取应该失败
	_, err = cache.Get(key)
	if err == nil {
		t.Fatal("Expired cache should not be retrievable")
	}
}

func TestMemoryCacheClearPattern(t *testing.T) {
	// 创建内存缓存配置
	config := &MemoryConfig{
		MaxSize:         100,
		DefaultTTL:      5 * time.Minute,
		CleanupInterval: 1 * time.Minute,
	}

	// 创建内存缓存
	cache := NewMemoryCache(config)
	if cache == nil {
		t.Fatal("Memory cache should not be nil")
	}

	// 设置多个键
	keys := []string{
		"user:1",
		"user:2",
		"session:1",
		"session:2",
		"data:1",
	}

	for i, key := range keys {
		err := cache.Set(key, "value", time.Minute)
		if err != nil {
			t.Fatalf("Failed to set cache %d: %v", i, err)
		}
	}

	// 清除user:模式的键
	cache.ClearPattern("user:")

	// 验证user:键被清除
	for _, key := range []string{"user:1", "user:2"} {
		_, err := cache.Get(key)
		if err == nil {
			t.Fatalf("Key %s should be cleared", key)
		}
	}

	// 验证其他键仍然存在
	for _, key := range []string{"session:1", "session:2", "data:1"} {
		_, err := cache.Get(key)
		if err != nil {
			t.Fatalf("Key %s should still exist", key)
		}
	}
}

func TestCacheManager(t *testing.T) {
	// 创建测试logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 创建缓存配置
	config := &Config{
		Memory: &MemoryConfig{
			MaxSize:         100,
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	// 创建缓存管理器
	manager, err := NewManager(config, logger)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	// 测试设置和获取
	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	ttl := time.Minute

	err = manager.Set(ctx, key, value, ttl)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	retrievedValue, err := manager.Get(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}

	if retrievedValue != value {
		t.Fatalf("Expected %s, got %s", value, retrievedValue)
	}

	// 测试删除
	err = manager.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete cache: %v", err)
	}

	// 验证删除
	_, err = manager.Get(ctx, key)
	if err == nil {
		t.Fatal("Deleted key should not exist")
	}
}

func TestRAGCacheService(t *testing.T) {
	// 创建测试logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 创建缓存配置
	config := &Config{
		Memory: &MemoryConfig{
			MaxSize:         100,
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	// 创建缓存管理器
	manager, err := NewManager(config, logger)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	// 创建RAG缓存服务
	ragCache := NewRAGCacheService(manager)
	if ragCache == nil {
		t.Fatal("RAG cache service should not be nil")
	}

	// 测试查询结果缓存
	query := "test query"
	results := []string{"result1", "result2", "result3"}
	ttl := time.Minute

	err = ragCache.SetQueryResult(query, results, ttl)
	if err != nil {
		t.Fatalf("Failed to set query result: %v", err)
	}

	retrievedResults, exists := ragCache.GetQueryResult(query)
	if !exists {
		t.Fatal("Query result should exist")
	}

	if retrievedResults == nil {
		t.Fatal("Retrieved results should not be nil")
	}

	// 测试嵌入向量缓存
	text := "test text"
	embedding := []float32{0.1, 0.2, 0.3, 0.4, 0.5}

	err = ragCache.SetEmbeddingCache(text, embedding, ttl)
	if err != nil {
		t.Fatalf("Failed to set embedding cache: %v", err)
	}

	retrievedEmbedding, exists := ragCache.GetEmbeddingCache(text)
	if !exists {
		t.Fatal("Embedding should exist")
	}

	if len(retrievedEmbedding) != len(embedding) {
		t.Fatalf("Embedding length mismatch: expected %d, got %d", len(embedding), len(retrievedEmbedding))
	}

	// 测试清除缓存
	err = ragCache.ClearQueryCache()
	if err != nil {
		t.Fatalf("Failed to clear query cache: %v", err)
	}

	err = ragCache.ClearEmbeddingCache()
	if err != nil {
		t.Fatalf("Failed to clear embedding cache: %v", err)
	}
}
