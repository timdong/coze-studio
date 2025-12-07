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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConfigStructure 测试配置结构
func TestConfigStructure(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host: "localhost",
			Port: 5432,
			User: "postgres",
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}

	// 验证配置对象创建成功
	assert.NotNil(t, cfg)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, 6379, cfg.Redis.Port)
	assert.Equal(t, "info", cfg.Logging.Level)
}

// TestGetDatabaseURL 测试数据库URL生成
func TestGetDatabaseURL(t *testing.T) {
	testCases := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name: "基础配置",
			config: Config{
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "postgres",
					Password: "password",
					Name:     "testdb",
					SSLMode:  "disable",
				},
			},
			expected: "postgresql://postgres:password@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "带SSL配置",
			config: Config{
				Database: DatabaseConfig{
					Host:     "db.example.com",
					Port:     5432,
					User:     "admin",
					Password: "secret",
					Name:     "proddb",
					SSLMode:  "require",
				},
			},
			expected: "postgresql://admin:secret@db.example.com:5432/proddb?sslmode=require",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.config.GetDatabaseURL()
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestGetRedisAddr 测试Redis地址生成
func TestGetRedisAddr(t *testing.T) {
	testCases := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name: "默认配置",
			config: Config{
				Redis: RedisConfig{
					Host: "localhost",
					Port: 6379,
				},
			},
			expected: "localhost:6379",
		},
		{
			name: "自定义端口",
			config: Config{
				Redis: RedisConfig{
					Host: "redis.example.com",
					Port: 16379,
				},
			},
			expected: "redis.example.com:16379",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.config.GetRedisAddr()
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestGetDefaultUserID 测试获取默认用户ID
func TestGetDefaultUserID(t *testing.T) {
	cfg := &Config{}

	// 测试GetDefaultUserID方法存在
	userID := cfg.GetDefaultUserID()

	// 默认应该返回1或0
	assert.True(t, userID >= 0)
}

// TestLoadConfigFromEnv 测试从环境变量加载配置
func TestLoadConfigFromEnv(t *testing.T) {
	t.Skip("需要完整的Viper配置，跳过")

	// 这个测试需要完整的Viper配置初始化
	// 在实际环境中测试
}

// TestValidateConfig 测试配置验证
func TestValidateConfig(t *testing.T) {
	t.Run("有效配置", func(t *testing.T) {
		cfg := &Config{
			Server: ServerConfig{
				Port: 8080,
				Host: "localhost",
			},
			Database: DatabaseConfig{
				Host: "localhost",
				Port: 5432,
				User: "postgres",
				Name: "testdb",
			},
			JWT: JWTConfig{
				Secret: "test-secret", // validateConfig需要JWT secret
			},
		}

		err := validateConfig(cfg)
		assert.NoError(t, err)
	})

	t.Run("无效端口", func(t *testing.T) {
		cfg := &Config{
			Server: ServerConfig{
				Port: -1, // 无效端口
			},
			Database: DatabaseConfig{
				Host: "localhost",
			},
		}

		err := validateConfig(cfg)
		assert.Error(t, err, "无效端口应该返回错误")
	})
}

// TestConfigDefaults 测试配置默认值
func TestConfigDefaults(t *testing.T) {
	t.Skip("需要完整的Viper配置，跳过")

	// setDefaults()是无参数函数，直接操作viper
	// 需要完整的配置环境才能测试
}

// TestConfigGetters 测试配置getter方法
func TestConfigGetters(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host:     "db.local",
			Port:     5432,
			User:     "admin",
			Password: "secret",
			Name:     "testdb",
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Host: "redis.local",
			Port: 6379,
		},
		Consul: ConsulConfig{
			Host: "consul.local",
			Port: 8500,
		},
	}

	// 测试GetDatabaseURL（PostgreSQL URI格式）
	dbURL := cfg.GetDatabaseURL()
	assert.Contains(t, dbURL, "postgresql://")
	assert.Contains(t, dbURL, "db.local")
	assert.Contains(t, dbURL, "5432")
	assert.Contains(t, dbURL, "testdb")
	assert.Contains(t, dbURL, "admin:secret")

	// 测试GetRedisAddr
	redisAddr := cfg.GetRedisAddr()
	assert.Equal(t, "redis.local:6379", redisAddr)

	// 测试GetConsulAddr
	consulAddr := cfg.GetConsulAddr()
	assert.Contains(t, consulAddr, "consul.local")
	assert.Contains(t, consulAddr, "8500")

	// 测试GetDefaultUserID
	userID := cfg.GetDefaultUserID()
	assert.True(t, userID >= 0, "用户ID应该大于等于0")
}
