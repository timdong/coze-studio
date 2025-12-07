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
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Environment string            `mapstructure:"environment"`
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	Consul      ConsulConfig      `mapstructure:"consul"`
	Services    ServicesConfig    `mapstructure:"services"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	Logging     LoggingConfig     `mapstructure:"logging"`
	Monitoring  MonitoringConfig  `mapstructure:"monitoring"`
	Security    SecurityConfig    `mapstructure:"security"`
	Vector      VectorConfig      `mapstructure:"vector"`
	AI          AIConfig          `mapstructure:"ai"`
	MinIO       MinIOConfig       `mapstructure:"minio"`
	Mem0        Mem0Config        `mapstructure:"mem0"`
	Development DevelopmentConfig `mapstructure:"development"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	DialTimeout  string `mapstructure:"dial_timeout"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
}

// ConsulConfig Consul配置
type ConsulConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Scheme     string `mapstructure:"scheme"`
	Datacenter string `mapstructure:"datacenter"`
	Token      string `mapstructure:"token"`
}

// ServicesConfig 服务配置
type ServicesConfig struct {
	Etrxtable   ServiceInfo `mapstructure:"etrxtable"`
	EntSMBMain  ServiceInfo `mapstructure:"ent_smb_main"`
	EntSMBMicro ServiceInfo `mapstructure:"ent_smb_micro"`
	Entrag      ServiceInfo `mapstructure:"entrag"`
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name        string   `mapstructure:"name"`
	Port        int      `mapstructure:"port"`
	HealthCheck string   `mapstructure:"health_check"`
	Tags        []string `mapstructure:"tags"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string `mapstructure:"secret"`
	ExpireHours        int    `mapstructure:"expire_hours"`
	RefreshExpireHours int    `mapstructure:"refresh_expire_hours"`
	Issuer             string `mapstructure:"issuer"`
	Disabled           bool   `mapstructure:"disabled"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	MetricsPort         int           `mapstructure:"metrics_port"`
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval"`
	PrometheusEnabled   bool          `mapstructure:"prometheus_enabled"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	EncryptionKey string          `mapstructure:"encryption_key"`
	CORS          CORSConfig      `mapstructure:"cors"`
	RateLimit     RateLimitConfig `mapstructure:"rate_limit"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Global  struct {
		RequestsPerSecond int `mapstructure:"requests_per_second"`
		Burst             int `mapstructure:"burst"`
	} `mapstructure:"global"`
	API struct {
		RequestsPerSecond int `mapstructure:"requests_per_second"`
		Burst             int `mapstructure:"burst"`
	} `mapstructure:"api"`
	PerIP struct {
		RequestsPerMinute int `mapstructure:"requests_per_minute"`
		Burst             int `mapstructure:"burst"`
	} `mapstructure:"per_ip"`
}

// VectorConfig 向量数据库配置
type VectorConfig struct {
	PGVectorEnabled     bool    `mapstructure:"pgvector_enabled"`
	EmbeddingDimension  int     `mapstructure:"embedding_dimension"`
	SimilarityThreshold float64 `mapstructure:"similarity_threshold"`
	MaxResults          int     `mapstructure:"max_results"`
}

// AIConfig AI配置
type AIConfig struct {
	Ollama      OllamaConfig      `mapstructure:"ollama"`
	SiliconFlow SiliconFlowConfig `mapstructure:"siliconflow"`
	RAG         RAGConfig         `mapstructure:"rag"`
}

// MinIOConfig MinIO配置
type MinIOConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Endpoint    string `mapstructure:"endpoint"`
	AccessKey   string `mapstructure:"access_key"`
	SecretKey   string `mapstructure:"secret_key"`
	UseSSL      bool   `mapstructure:"use_ssl"`
	BucketName  string `mapstructure:"bucket_name"`
	Region      string `mapstructure:"region"`
	MaxFileSize int64  `mapstructure:"max_file_size"`
}

// OllamaConfig Ollama配置
type OllamaConfig struct {
	Host       string `mapstructure:"host"`
	Model      string `mapstructure:"model"`
	Timeout    string `mapstructure:"timeout"`
	EmbedModel string `mapstructure:"embed_model"`
}

// SiliconFlowConfig 硅基流动配置
type SiliconFlowConfig struct {
	BaseURL      string `mapstructure:"base_url"`
	APIKey       string `mapstructure:"api_key"`
	DefaultModel string `mapstructure:"default_model"`
	Timeout      string `mapstructure:"timeout"`
}

// RAGConfig RAG配置
type RAGConfig struct {
	ChunkSize        int     `mapstructure:"chunk_size"`
	ChunkOverlap     int     `mapstructure:"chunk_overlap"`
	MaxTokens        int     `mapstructure:"max_tokens"`
	Temperature      float64 `mapstructure:"temperature"`
	TokenEncoding    string  `mapstructure:"token_encoding"`
	MaxSimilarChunks int     `mapstructure:"max_similar_chunks"`
	MinChunkSize     int     `mapstructure:"min_chunk_size"`
}

// Mem0Config Mem0配置
type Mem0Config struct {
	Enabled   bool   `mapstructure:"enabled"`
	BaseURL   string `mapstructure:"base_url"`
	APIKey    string `mapstructure:"api_key"`
	Timeout   string `mapstructure:"timeout"`
	RetryTimes int  `mapstructure:"retry_times"`
	Infer     bool   `mapstructure:"infer"` // 是否进行事实提取和冲突检测（true=自动提取事实并处理冲突，false=直接保存对话内容）
}

// DevelopmentConfig 开发配置
type DevelopmentConfig struct {
	DefaultUserID string `mapstructure:"default_user_id"`
}

// GetDefaultUserID 获取开发环境的默认用户ID
func (c *Config) GetDefaultUserID() uint {
	if c.Environment == "development" || c.Environment == "dev" {
		// 在开发环境中，返回配置的默认用户ID
		if c.Development.DefaultUserID != "" {
			// 这里可以根据需要映射用户名到ID
			switch c.Development.DefaultUserID {
			case "admin":
				return 2
			case "timdong":
				return 4
			case "testuser2":
				return 9
			case "test111":
				return 10
			default:
				return 2 // 默认使用admin用户
			}
		}
		return 2 // 默认使用admin用户
	}
	// 在生产环境中，应该从JWT token获取
	return 0 // 表示需要从认证信息获取
}

// Load 加载配置
func Load() (*Config, error) {
	// 设置默认值
	setDefaults()

	// 从配置文件加载
	if err := loadFromFile(); err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// 从环境变量加载
	loadFromEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

var (
	globalConfig *Config
)

// SetGlobalConfig 设置全局配置
func SetGlobalConfig(cfg *Config) {
	globalConfig = cfg
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if globalConfig == nil {
		// 如果全局配置未设置，尝试加载
		if cfg, err := Load(); err == nil {
			globalConfig = cfg
		} else {
			// 如果加载失败，返回一个默认配置
			globalConfig = &Config{
				JWT: JWTConfig{
					Secret: "your-super-secret-jwt-key-change-in-production",
					Issuer: "coze-studio",
				},
			}
		}
	}
	return globalConfig
}

// setDefaults 设置默认值
func setDefaults() {
	viper.SetDefault("environment", "development")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")

	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "erdcloud")
	viper.SetDefault("database.password", "Pw!123456")
	viper.SetDefault("database.name", "etrxlite_db")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 200)
	viper.SetDefault("database.max_idle_conns", 50)
	viper.SetDefault("database.conn_max_lifetime", "600s")

	viper.SetDefault("redis.enabled", true)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)
	viper.SetDefault("redis.dial_timeout", "5s")
	viper.SetDefault("redis.read_timeout", "3s")
	viper.SetDefault("redis.write_timeout", "3s")

	viper.SetDefault("consul.host", "localhost")
	viper.SetDefault("consul.port", 8500)
	viper.SetDefault("consul.scheme", "http")
	viper.SetDefault("consul.datacenter", "dc1")
	viper.SetDefault("consul.token", "")

	viper.SetDefault("jwt.secret", "your-super-secret-jwt-key-change-in-production")
	viper.SetDefault("jwt.expire_hours", 24)
	viper.SetDefault("jwt.refresh_expire_hours", 168)
	viper.SetDefault("jwt.issuer", "etrxlite")
	viper.SetDefault("jwt.disabled", false)

	viper.SetDefault("logging.level", "debug")
	viper.SetDefault("logging.format", "text")
	viper.SetDefault("logging.output", "file")
	viper.SetDefault("logging.file_path", "logs/etrxlite.log")
	viper.SetDefault("logging.max_size", 100)
	viper.SetDefault("logging.max_age", 30)
	viper.SetDefault("logging.max_backups", 10)
	viper.SetDefault("logging.compress", true)

	viper.SetDefault("monitoring.metrics_port", 9090)
	viper.SetDefault("monitoring.health_check_interval", "30s")
	viper.SetDefault("monitoring.prometheus_enabled", true)

	viper.SetDefault("security.cors.allowed_origins", []string{"http://localhost:3000", "http://localhost:3001", "*"})
	viper.SetDefault("security.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"})
	viper.SetDefault("security.cors.allowed_headers", []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"})
	viper.SetDefault("security.cors.allow_credentials", true)
	viper.SetDefault("security.cors.expose_headers", []string{"Content-Length", "Content-Type"})

	viper.SetDefault("security.rate_limit.enabled", true)
	viper.SetDefault("security.rate_limit.global.requests_per_second", 5000)
	viper.SetDefault("security.rate_limit.global.burst", 3000)
	viper.SetDefault("security.rate_limit.api.requests_per_second", 5000)
	viper.SetDefault("security.rate_limit.api.burst", 3000)
	viper.SetDefault("security.rate_limit.per_ip.requests_per_minute", 500000)
	viper.SetDefault("security.rate_limit.per_ip.burst", 200000)
	viper.SetDefault("security.encryption_key", "")

	viper.SetDefault("vector.pgvector_enabled", true)
	viper.SetDefault("vector.embedding_dimension", 768)
	viper.SetDefault("vector.similarity_threshold", 0.7)
	viper.SetDefault("vector.max_results", 100)

	viper.SetDefault("ai.ollama.host", "http://localhost:11434")
	viper.SetDefault("ai.ollama.model", "llama3.2:3b")
	viper.SetDefault("ai.ollama.timeout", "600s") // 大模型（如30B）需要更长时间，默认10分钟
	viper.SetDefault("ai.ollama.embed_model", "nomic-embed-text")

	viper.SetDefault("ai.rag.chunk_size", 600)
	viper.SetDefault("ai.rag.chunk_overlap", 80)
	viper.SetDefault("ai.rag.max_tokens", 4000)
	viper.SetDefault("ai.rag.temperature", 0.7)
	viper.SetDefault("ai.rag.token_encoding", "cl100k_base")
	viper.SetDefault("ai.rag.max_similar_chunks", 4)
	viper.SetDefault("ai.rag.min_chunk_size", 50)

	// MinIO默认配置
	viper.SetDefault("minio.enabled", true)
	viper.SetDefault("minio.endpoint", "localhost:9000")
	viper.SetDefault("minio.access_key", "minioadmin")
	viper.SetDefault("minio.secret_key", "minioadmin")
	viper.SetDefault("minio.use_ssl", false)
	viper.SetDefault("minio.bucket_name", "documents")
	viper.SetDefault("minio.region", "us-east-1")
	viper.SetDefault("minio.max_file_size", int64(104857600)) // 100MB

	// Mem0默认配置
	viper.SetDefault("mem0.enabled", false)
	viper.SetDefault("mem0.base_url", "http://localhost:8000")
	viper.SetDefault("mem0.api_key", "")
	viper.SetDefault("mem0.timeout", "30s")
	viper.SetDefault("mem0.retry_times", 3)
	viper.SetDefault("mem0.infer", true) // 默认启用事实提取和冲突检测

	viper.SetDefault("development.default_user_id", "1")
}

// loadFromFile 从配置文件加载
func loadFromFile() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		// 配置文件不存在，使用默认值
		return nil
	}

	return nil
}

// loadFromEnv 从环境变量加载
func loadFromEnv() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", config.Database.Port)
	}

	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	return nil
}

// GetDatabaseURL 获取数据库连接URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetRedisAddr 获取Redis地址
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// GetConsulAddr 获取Consul地址
func (c *Config) GetConsulAddr() string {
	return fmt.Sprintf("%s://%s:%d", c.Consul.Scheme, c.Consul.Host, c.Consul.Port)
}
