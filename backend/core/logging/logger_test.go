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

package logging

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewDefaultLogger(t *testing.T) {
	logger, err := NewDefaultLogger()
	if err != nil {
		t.Fatalf("Failed to create default logger: %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("Logger should not be nil")
	}
}

func TestNewLogger(t *testing.T) {
	config := &Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("Logger should not be nil")
	}
}

func TestGlobalLogger(t *testing.T) {
	// 创建测试logger
	logger, err := NewDefaultLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// 设置全局logger
	SetGlobalLogger(logger)

	// 测试全局logger访问
	globalLogger := Log
	// Global是结构体，不能与nil比较，我们通过调用方法来测试

	// 测试日志输出（这里只是确保不会panic）
	globalLogger.Info("Test message", zap.String("test", "value"))
	globalLogger.Debug("Debug message", zap.Int("count", 42))
	globalLogger.Warn("Warning message", zap.Bool("flag", true))
	globalLogger.Error("Error message", zap.Error(nil))
}

func TestLoggerLevels(t *testing.T) {
	logger, err := NewDefaultLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// 测试不同级别的日志
	logger.Info("Info level test")
	logger.Debug("Debug level test")
	logger.Warn("Warn level test")
	logger.Error("Error level test")
}

func TestLoggerWithFields(t *testing.T) {
	logger, err := NewDefaultLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// 测试带字段的日志
	logger.Info("User action",
		zap.String("user_id", "123"),
		zap.String("action", "login"),
		zap.Int("attempts", 1),
		zap.Bool("success", true),
	)
}

func TestLoggerConfig(t *testing.T) {
	configs := []*Config{
		{
			Level:  "debug",
			Format: "json",
			Output: "stdout",
		},
		{
			Level:  "info",
			Format: "console",
			Output: "stderr",
		},
		{
			Level:  "warn",
			Format: "json",
			Output: "stdout",
		},
	}

	for i, config := range configs {
		logger, err := NewLogger(config)
		if err != nil {
			t.Fatalf("Config %d failed: %v", i, err)
		}
		defer logger.Sync()

		if logger == nil {
			t.Fatalf("Config %d: Logger should not be nil", i)
		}
	}
}
