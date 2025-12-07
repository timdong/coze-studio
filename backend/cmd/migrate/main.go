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
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

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

	logging.Log.Info("Starting database migration...")

	// 3. 连接数据库
	db, err := connectDatabase()
	if err != nil {
		logging.Log.Fatal("Failed to connect database", zap.Error(err))
	}

	// 4. 运行数据库迁移
	if err := migrations.AutoMigrate(db); err != nil {
		logging.Log.Fatal("Failed to migrate database", zap.Error(err))
	}

	logging.Log.Info("✅ Database migration completed successfully!")
	fmt.Println("\n✅ 数据库迁移完成！所有 Coze Studio 核心表已创建。")
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")

	// 支持环境变量覆盖
	viper.AutomaticEnv()
	viper.SetEnvPrefix("COZE")

	return viper.ReadInConfig()
}

func initLogging() error {
	logger, err := logging.NewLogger(&logging.Config{
		Level:  "info",
		Format: "text",
		Output: "stdout",
	})
	if err != nil {
		return err
	}
	logging.SetGlobalLogger(logger)
	return nil
}

func connectDatabase() (*gorm.DB, error) {
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	name := viper.GetString("database.name")
	sslmode := viper.GetString("database.sslmode")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		host, user, password, name, port, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}

