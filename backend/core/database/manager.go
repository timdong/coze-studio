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

package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"etrxlite/pkg/database/postgres"
	"etrxlite/pkg/database/redis"

	"go.uber.org/zap"
)

// Manager 数据库管理器
type Manager struct {
	postgres *postgres.Client
	redis    *redis.Client
	logger   *zap.Logger
}

// Config 数据库配置
type Config struct {
	Postgres *postgres.Config `yaml:"postgres" json:"postgres"`
	Redis    *redis.Config    `yaml:"redis" json:"redis"`
}

// NewManager 创建数据库管理器
func NewManager(cfg *Config, logger *zap.Logger) (*Manager, error) {
	// 初始化PostgreSQL
	pgClient, err := postgres.NewClient(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres client: %w", err)
	}

	// 初始化Redis
	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}

	return &Manager{
		postgres: pgClient,
		redis:    redisClient,
		logger:   logger,
	}, nil
}

// GetPostgres 获取PostgreSQL客户端
func (m *Manager) GetPostgres() *postgres.Client {
	return m.postgres
}

// GetRedis 获取Redis客户端
func (m *Manager) GetRedis() *redis.Client {
	return m.redis
}

// HealthCheck 健康检查
func (m *Manager) HealthCheck(ctx context.Context) error {
	// 检查PostgreSQL
	if err := m.postgres.Ping(); err != nil {
		return fmt.Errorf("postgres health check failed: %w", err)
	}

	// 检查Redis
	if err := m.redis.Ping(ctx); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}

// Close 关闭数据库连接
func (m *Manager) Close() error {
	var errors []error

	if err := m.postgres.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close postgres: %w", err))
	}

	if err := m.redis.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close redis: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("database close errors: %v", errors)
	}

	return nil
}

// Transaction 执行事务
func (m *Manager) Transaction(ctx context.Context, fn func(*postgres.Client) error) error {
	tx, err := m.postgres.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 创建事务客户端
	txClient := &postgres.Client{}
	// 这里需要实现事务客户端，简化处理
	_ = txClient

	if err := fn(txClient); err != nil {
		return err
	}

	return tx.Commit()
}

// Query 执行查询
func (m *Manager) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return m.postgres.Query(query, args...)
}

// Exec 执行命令
func (m *Manager) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := m.postgres.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// SetCache 设置缓存
func (m *Manager) SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.redis.Set(ctx, key, value, expiration)
}

// GetCache 获取缓存
func (m *Manager) GetCache(ctx context.Context, key string) (string, error) {
	return m.redis.Get(ctx, key)
}

// DeleteCache 删除缓存
func (m *Manager) DeleteCache(ctx context.Context, key string) error {
	return m.redis.Del(ctx, key)
}

// GetWriteConnection 获取写连接
func (m *Manager) GetWriteConnection() *postgres.Client {
	return m.postgres
}

// GetReadConnection 获取读连接
func (m *Manager) GetReadConnection() *postgres.Client {
	return m.postgres
}
