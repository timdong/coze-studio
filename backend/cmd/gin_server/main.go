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

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/api/middleware"
	"github.com/coze-dev/coze-studio/backend/api/router"
	"github.com/coze-dev/coze-studio/backend/application"
	"github.com/coze-dev/coze-studio/backend/core/logging"
	"github.com/coze-dev/coze-studio/backend/migrations"
)

func main() {
	// 1. 加载配置
	if err := loadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	if err := initLogging(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logging: %v\n", err)
		os.Exit(1)
	}

	// 记录配置加载成功（此时日志系统已初始化）
	logging.Log.Info("Config loaded successfully",
		zap.String("file", viper.ConfigFileUsed()),
	)

	logging.Log.Info("Starting Coze Studio Super (Gin version)...")

	// 3. 连接数据库
	db, err := connectDatabase()
	if err != nil {
		logging.Log.Fatal("Failed to connect database", zap.Error(err))
	}

	// 4. 运行数据库迁移
	if err := migrations.AutoMigrate(db); err != nil {
		logging.Log.Fatal("Failed to migrate database", zap.Error(err))
	}
	logging.Log.Info("Database migration completed")

	// 5. 设置存储服务环境变量（必须在 application.Init 之前）
	setupStorageEnvironment()

	// 6. 初始化所有应用服务（包括 SearchSVC、UserSVC 等）
	// 传入已连接的 PostgreSQL 数据库连接
	ctx := context.Background()
	if err := application.Init(ctx, db); err != nil {
		logging.Log.Fatal("Failed to initialize application services", zap.Error(err))
	}
	logging.Log.Info("Application services initialized successfully")

	// 7. 初始化 Gin 引擎
	ginEngine := initGinEngine(db)

	// 8. 启动服务器
	startServer(ginEngine)
}

// setupStorageEnvironment 从配置文件读取存储配置并设置环境变量
func setupStorageEnvironment() {
	// 从 Viper 读取配置并设置环境变量（存储服务目前从环境变量读取配置）
	if viper.GetBool("minio.enabled") {
		os.Setenv("STORAGE_TYPE", "minio")
		os.Setenv("MINIO_ENDPOINT", viper.GetString("minio.endpoint"))
		os.Setenv("MINIO_AK", viper.GetString("minio.access_key"))
		os.Setenv("MINIO_SK", viper.GetString("minio.secret_key"))
		os.Setenv("STORAGE_BUCKET", viper.GetString("minio.bucket_name"))
		if viper.GetBool("minio.use_ssl") {
			os.Setenv("MINIO_USE_SSL", "true")
		} else {
			os.Setenv("MINIO_USE_SSL", "false")
		}
		logging.Log.Info("Storage environment variables set for MinIO",
			zap.String("endpoint", viper.GetString("minio.endpoint")),
			zap.String("bucket", viper.GetString("minio.bucket_name")),
		)
	} else {
		// 如果没有启用 MinIO，设置默认值或从其他配置读取
		// 这里可以根据需要添加其他存储类型的配置
		logging.Log.Warn("MinIO is not enabled, storage type may not be set")
	}
}

// loadConfig 加载配置
func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend")

	// 支持环境变量覆盖
	viper.AutomaticEnv()
	viper.SetEnvPrefix("COZE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		// 检查是否是配置文件不存在的错误
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，使用默认值
			// 不返回错误，继续使用默认配置
		} else {
			// 其他错误（如格式错误）
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	// 注意：这里不使用 logging.Log，因为日志系统还未初始化
	// 配置加载成功的日志会在 initLogging() 之后记录

	return nil
}

// initLogging 初始化日志系统
func initLogging() error {
	cfg := &logging.Config{
		Level:      viper.GetString("logging.level"),
		Format:     viper.GetString("logging.format"),
		Output:     viper.GetString("logging.output"),
		FilePath:   viper.GetString("logging.file_path"),
		MaxSize:    viper.GetInt("logging.max_size"),
		MaxBackups: viper.GetInt("logging.max_backups"),
		MaxAge:     viper.GetInt("logging.max_age"),
		Compress:   viper.GetBool("logging.compress"),
	}

	logger, err := logging.NewLogger(cfg)
	if err != nil {
		return err
	}

	// 设置全局日志
	logging.SetGlobalLogger(logger)

	return nil
}

// connectDatabase 连接数据库
func connectDatabase() (*gorm.DB, error) {
	// 构建 DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		viper.GetString("database.host"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.name"),
		viper.GetInt("database.port"),
		viper.GetString("database.sslmode"),
		viper.GetString("database.timezone"),
	)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// DisableForeignKeyConstraintWhenMigrating: true, // 如果需要禁用外键
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(viper.GetInt("database.max_open_conns"))
	sqlDB.SetMaxIdleConns(viper.GetInt("database.max_idle_conns"))
	sqlDB.SetConnMaxLifetime(viper.GetDuration("database.conn_max_lifetime"))

	logging.Log.Info("Database connected successfully",
		zap.String("host", viper.GetString("database.host")),
		zap.Int("port", viper.GetInt("database.port")),
		zap.String("database", viper.GetString("database.name")),
	)

	return db, nil
}

// initGinEngine 初始化 Gin 引擎
func initGinEngine(db *gorm.DB) *gin.Engine {
	// 根据环境设置模式
	if !viper.GetBool("app.debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建引擎（不使用默认中间件）
	engine := gin.New()

	// 注册中间件（顺序重要）
	engine.Use(middleware.GinLogger())     // 日志
	engine.Use(middleware.GinRecovery())   // 恢复
	engine.Use(middleware.GinCORS())       // CORS
	engine.Use(middleware.GinRequestID())  // 请求 ID
	// engine.Use(middleware.GinAuth())    // 认证（可选，根据路由需要）

	// 注册路由
	router.GinRegister(engine, db)

	return engine
}

// startServer 启动服务器
func startServer(engine *gin.Engine) {
	// 获取服务器配置
	host := viper.GetString("server.host")
	port := viper.GetInt("server.port")
	addr := fmt.Sprintf("%s:%d", host, port)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:           addr,
		Handler:        engine,
		ReadTimeout:    viper.GetDuration("server.read_timeout"),
		WriteTimeout:   viper.GetDuration("server.write_timeout"),
		IdleTimeout:    viper.GetDuration("server.idle_timeout"),
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// 启动服务器（异步）
	go func() {
		logging.Log.Info("Starting HTTP server",
			zap.String("addr", addr),
			zap.Bool("debug", viper.GetBool("app.debug")),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logging.Log.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logging.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logging.Log.Info("Server exited")
}

