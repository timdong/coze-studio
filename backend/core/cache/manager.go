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
	"fmt"
	"sync"
	"time"

	"etrxlite/pkg/database/redis"

	"go.uber.org/zap"
)

// Manager 缓存管理器
type Manager struct {
	redis  *redis.Client
	memory *MemoryCache
	config *Config
	logger *zap.Logger
	mu     sync.RWMutex
}

// Config 缓存配置
type Config struct {
	Redis  *redis.Config `yaml:"redis" json:"redis"`
	Memory *MemoryConfig `yaml:"memory" json:"memory"`
}

// MemoryConfig 内存缓存配置
type MemoryConfig struct {
	MaxSize         int           `yaml:"max_size" json:"max_size"`
	DefaultTTL      time.Duration `yaml:"default_ttl" json:"default_ttl"`
	CleanupInterval time.Duration `yaml:"cleanup_interval" json:"cleanup_interval"`
}

// NewManager 创建缓存管理器
func NewManager(config *Config, logger *zap.Logger) (*Manager, error) {
	manager := &Manager{
		config: config,
		logger: logger,
	}

	// 初始化Redis
	if config.Redis != nil {
		redisClient, err := redis.NewClient(config.Redis)
		if err != nil {
			return nil, fmt.Errorf("failed to create redis client: %w", err)
		}
		manager.redis = redisClient
	}

	// 初始化内存缓存
	if config.Memory != nil {
		manager.memory = NewMemoryCache(config.Memory)
	}

	return manager, nil
}

// Set 设置缓存
func (m *Manager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// 优先使用Redis
	if m.redis != nil {
		return m.redis.Set(ctx, key, value, ttl)
	}

	// 回退到内存缓存
	if m.memory != nil {
		return m.memory.Set(key, value, ttl)
	}

	return fmt.Errorf("no cache backend available")
}

// Get 获取缓存
func (m *Manager) Get(ctx context.Context, key string) (string, error) {
	// 优先使用Redis
	if m.redis != nil {
		return m.redis.Get(ctx, key)
	}

	// 回退到内存缓存
	if m.memory != nil {
		return m.memory.Get(key)
	}

	return "", fmt.Errorf("no cache backend available")
}

// Delete 删除缓存
func (m *Manager) Delete(ctx context.Context, keys ...string) error {
	// 优先使用Redis
	if m.redis != nil {
		return m.redis.Del(ctx, keys...)
	}

	// 回退到内存缓存
	if m.memory != nil {
		for _, key := range keys {
			m.memory.Delete(key)
		}
		return nil
	}

	return fmt.Errorf("no cache backend available")
}

// Exists 检查键是否存在
func (m *Manager) Exists(ctx context.Context, keys ...string) (int64, error) {
	// 优先使用Redis
	if m.redis != nil {
		return m.redis.Exists(ctx, keys...)
	}

	// 回退到内存缓存
	if m.memory != nil {
		count := int64(0)
		for _, key := range keys {
			if m.memory.Exists(key) {
				count++
			}
		}
		return count, nil
	}

	return 0, fmt.Errorf("no cache backend available")
}

// SetWithTags 设置带标签的缓存
func (m *Manager) SetWithTags(ctx context.Context, key string, value interface{}, ttl time.Duration, tags ...string) error {
	// 设置缓存
	if err := m.Set(ctx, key, value, ttl); err != nil {
		return err
	}

	// 设置标签索引
	for _, tag := range tags {
		tagKey := fmt.Sprintf("tag:%s", tag)
		if err := m.Set(ctx, tagKey, key, ttl); err != nil {
			m.logger.Warn("Failed to set cache tag", zap.String("tag", tag), zap.Error(err))
		}
	}

	return nil
}

// InvalidateByTag 根据标签失效缓存
func (m *Manager) InvalidateByTag(ctx context.Context, tag string) error {
	tagKey := fmt.Sprintf("tag:%s", tag)

	// 获取标签下的所有键
	keys, err := m.Get(ctx, tagKey)
	if err != nil {
		return err
	}

	if keys != "" {
		// 删除标签下的所有键
		return m.Delete(ctx, keys)
	}

	return nil
}

// HealthCheck 健康检查
func (m *Manager) HealthCheck(ctx context.Context) error {
	if m.redis != nil {
		return m.redis.Ping(ctx)
	}
	return nil
}

// GetStats 获取统计信息
func (m *Manager) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	if m.redis != nil {
		stats["redis"] = m.redis.PoolStats()
	}

	if m.memory != nil {
		stats["memory"] = m.memory.GetStats()
	}

	return stats
}

// ClearMemoryPattern 清除内存缓存中匹配模式的键
func (m *Manager) ClearMemoryPattern(pattern string) {
	if m.memory != nil {
		m.memory.ClearPattern(pattern)
	}
}

// Publish 发布消息到频道
func (m *Manager) Publish(ctx context.Context, channel string, message interface{}) error {
	if m.redis != nil {
		return m.redis.Publish(ctx, channel, message)
	}
	return fmt.Errorf("no redis backend available for pub/sub")
}

// Subscribe 订阅频道，返回消息通道
func (m *Manager) Subscribe(ctx context.Context, channel string) (<-chan string, error) {
	if m.redis != nil {
		pubsub := m.redis.Subscribe(ctx, channel)
		ch := pubsub.Channel()
		
		// 创建字符串消息通道
		msgChan := make(chan string, 100)
		
		// 启动goroutine转换消息
		go func() {
			defer close(msgChan)
			for {
				select {
				case <-ctx.Done():
					pubsub.Close()
					return
				case msg, ok := <-ch:
					if !ok {
						return
					}
					msgChan <- msg.Payload
				}
			}
		}()
		
		return msgChan, nil
	}
	return nil, fmt.Errorf("no redis backend available for pub/sub")
}

// Unsubscribe 取消订阅频道
func (m *Manager) Unsubscribe(ctx context.Context, channels ...string) error {
	// Redis PubSub的取消订阅由关闭PubSub对象完成
	// 这里返回nil，实际的取消订阅在Subscribe方法的context取消时完成
	return nil
}

// Close 关闭缓存管理器
func (m *Manager) Close() error {
	if m.redis != nil {
		return m.redis.Close()
	}
	return nil
}
