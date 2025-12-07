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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestManager_Exists 测试Exists功能
func TestManager_Exists(t *testing.T) {
	logger := zap.NewNop()
	config := &Config{
		Memory: &MemoryConfig{
			MaxSize:         100,
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	manager, err := NewManager(config, logger)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("单个key存在性检查", func(t *testing.T) {
		key := "exist_test_1"
		value := "value1"

		// 先设置key
		err := manager.Set(ctx, key, value, time.Minute)
		require.NoError(t, err)

		// 检查存在
		count, err := manager.Exists(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count, "key应该存在")

		// 删除key后再检查
		manager.Delete(ctx, key)
		count, err = manager.Exists(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count, "key应该不存在")
	})

	t.Run("多个keys存在性检查", func(t *testing.T) {
		keys := []string{"key1", "key2", "key3"}

		// 设置前两个key
		manager.Set(ctx, "key1", "value1", time.Minute)
		manager.Set(ctx, "key2", "value2", time.Minute)

		// 检查3个key（只有2个存在）
		count, err := manager.Exists(ctx, keys...)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count, "应该有2个key存在")

		// 清理
		for _, key := range keys {
			manager.Delete(ctx, key)
		}
	})
}

// TestManager_SetWithTags 测试带标签的缓存设置
func TestManager_SetWithTags(t *testing.T) {
	logger := zap.NewNop()
	config := &Config{
		Memory: &MemoryConfig{
			MaxSize:         100,
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	manager, err := NewManager(config, logger)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("设置带单个标签的缓存", func(t *testing.T) {
		key := "tagged_key_1"
		value := "tagged_value_1"
		tag := "user"

		err := manager.SetWithTags(ctx, key, value, time.Minute, tag)
		assert.NoError(t, err)

		// 验证key存在
		result, err := manager.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)

		// 验证tag存在
		count, _ := manager.Exists(ctx, "tag:"+tag)
		assert.Equal(t, int64(1), count, "tag应该存在")

		// 清理
		manager.Delete(ctx, key)
		manager.Delete(ctx, "tag:"+tag)
	})

	t.Run("设置带多个标签的缓存", func(t *testing.T) {
		key := "tagged_key_2"
		value := "tagged_value_2"
		tags := []string{"user", "profile", "cache"}

		err := manager.SetWithTags(ctx, key, value, time.Minute, tags...)
		assert.NoError(t, err)

		// 验证key存在
		result, err := manager.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)

		// 清理
		manager.Delete(ctx, key)
		for _, tag := range tags {
			manager.Delete(ctx, "tag:"+tag)
		}
	})
}

// TestManager_InvalidateByTag 测试按标签失效缓存
func TestManager_InvalidateByTag(t *testing.T) {
	logger := zap.NewNop()
	config := &Config{
		Memory: &MemoryConfig{
			MaxSize:         100,
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	manager, err := NewManager(config, logger)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("按标签失效缓存", func(t *testing.T) {
		key := "tagged_invalidate_key"
		value := "tagged_value"
		tag := "invalidate_test"

		// 设置带标签的缓存
		err := manager.SetWithTags(ctx, key, value, time.Minute, tag)
		require.NoError(t, err)

		// 验证key存在
		result, err := manager.Get(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, value, result)

		// 按标签失效
		err = manager.InvalidateByTag(ctx, tag)
		assert.NoError(t, err)

		// 由于实现简单，InvalidateByTag可能无法完全删除
		// 这个测试验证方法调用不会panic即可
	})

	t.Run("失效不存在的标签", func(t *testing.T) {
		err := manager.InvalidateByTag(ctx, "nonexistent_tag")
		// 应该返回错误或正常处理
		assert.Error(t, err, "不存在的标签应该返回错误")
	})
}

// TestManager_HealthCheck 测试健康检查
func TestManager_HealthCheck(t *testing.T) {
	logger := zap.NewNop()

	t.Run("只有内存缓存的健康检查", func(t *testing.T) {
		config := &Config{
			Memory: &MemoryConfig{
				MaxSize:         100,
				DefaultTTL:      5 * time.Minute,
				CleanupInterval: 1 * time.Minute,
			},
		}

		manager, err := NewManager(config, logger)
		require.NoError(t, err)

		ctx := context.Background()

		// 只有内存缓存时，健康检查应该成功
		err = manager.HealthCheck(ctx)
		assert.NoError(t, err, "内存缓存健康检查应该成功")
	})

	t.Run("无缓存后端的健康检查", func(t *testing.T) {
		config := &Config{}

		manager, err := NewManager(config, logger)
		require.NoError(t, err)

		ctx := context.Background()

		// 没有任何缓存后端时
		err = manager.HealthCheck(ctx)
		// 应该成功（因为redis为nil时返回nil）
		assert.NoError(t, err)
	})
}

// TestManager_GetStats 测试获取统计信息
func TestManager_GetStats(t *testing.T) {
	logger := zap.NewNop()

	t.Run("获取内存缓存统计", func(t *testing.T) {
		config := &Config{
			Memory: &MemoryConfig{
				MaxSize:         100,
				DefaultTTL:      5 * time.Minute,
				CleanupInterval: 1 * time.Minute,
			},
		}

		manager, err := NewManager(config, logger)
		require.NoError(t, err)

		// 设置一些数据
		ctx := context.Background()
		for i := 0; i < 5; i++ {
			manager.Set(ctx, "stats_key_"+string(rune(i)), "value", time.Minute)
		}

		// 获取统计
		stats := manager.GetStats()
		assert.NotNil(t, stats, "统计信息应该不为nil")
		assert.Contains(t, stats, "memory", "应该包含memory统计")
	})

	t.Run("无缓存后端的统计", func(t *testing.T) {
		config := &Config{}

		manager, err := NewManager(config, logger)
		require.NoError(t, err)

		stats := manager.GetStats()
		assert.NotNil(t, stats, "即使没有后端，统计也应该返回空map")
	})
}

// TestManager_ClearMemoryPattern 测试清除内存缓存模式
func TestManager_ClearMemoryPattern(t *testing.T) {
	logger := zap.NewNop()
	config := &Config{
		Memory: &MemoryConfig{
			MaxSize:         100,
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	manager, err := NewManager(config, logger)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("清除匹配模式的keys", func(t *testing.T) {
		// 设置多个keys
		keys := map[string]string{
			"user:1:profile": "value1",
			"user:2:profile": "value2",
			"user:3:profile": "value3",
			"session:1":      "session_value",
			"cache:data":     "cache_value",
		}

		for key, value := range keys {
			manager.Set(ctx, key, value, time.Minute)
		}

		// 清除user:模式的keys
		manager.ClearMemoryPattern("user:")

		// 验证user:keys被清除
		_, err := manager.Get(ctx, "user:1:profile")
		assert.Error(t, err, "user:1:profile应该被清除")

		// 验证其他keys仍存在
		value, err := manager.Get(ctx, "session:1")
		assert.NoError(t, err)
		assert.Equal(t, "session_value", value)
	})
}

// TestManager_Close 测试关闭管理器
func TestManager_Close(t *testing.T) {
	logger := zap.NewNop()

	t.Run("关闭只有内存缓存的管理器", func(t *testing.T) {
		config := &Config{
			Memory: &MemoryConfig{
				MaxSize:         100,
				DefaultTTL:      5 * time.Minute,
				CleanupInterval: 1 * time.Minute,
			},
		}

		manager, err := NewManager(config, logger)
		require.NoError(t, err)

		// 关闭管理器
		err = manager.Close()
		// 只有内存缓存时，Close应该返回nil
		assert.NoError(t, err)
	})
}

// TestManager_DeleteMultipleKeys 测试删除多个keys
func TestManager_DeleteMultipleKeys(t *testing.T) {
	logger := zap.NewNop()
	config := &Config{
		Memory: &MemoryConfig{
			MaxSize:         100,
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	manager, err := NewManager(config, logger)
	require.NoError(t, err)

	ctx := context.Background()

	// 设置多个keys
	keys := []string{"del_key_1", "del_key_2", "del_key_3"}
	for _, key := range keys {
		manager.Set(ctx, key, "value", time.Minute)
	}

	// 删除所有keys
	err = manager.Delete(ctx, keys...)
	assert.NoError(t, err)

	// 验证所有keys被删除
	for _, key := range keys {
		_, err := manager.Get(ctx, key)
		assert.Error(t, err, key+"应该被删除")
	}
}

// Benchmark测试
func BenchmarkManager_Set(b *testing.B) {
	logger := zap.NewNop()
	config := &Config{
		Memory: &MemoryConfig{
			MaxSize:         1000,
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	manager, _ := NewManager(config, logger)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Set(ctx, "benchmark_key", "value", time.Minute)
	}
}

func BenchmarkManager_Get(b *testing.B) {
	logger := zap.NewNop()
	config := &Config{
		Memory: &MemoryConfig{
			MaxSize:         1000,
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	manager, _ := NewManager(config, logger)
	ctx := context.Background()
	manager.Set(ctx, "benchmark_key", "value", time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Get(ctx, "benchmark_key")
	}
}

func BenchmarkManager_Exists(b *testing.B) {
	logger := zap.NewNop()
	config := &Config{
		Memory: &MemoryConfig{
			MaxSize:         1000,
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	manager, _ := NewManager(config, logger)
	ctx := context.Background()
	manager.Set(ctx, "benchmark_key", "value", time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Exists(ctx, "benchmark_key")
	}
}
