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
	"sync"

	"go.uber.org/zap"
)

var (
	globalLogger *Logger
	once         sync.Once
)

// SetGlobalLogger 设置全局日志器
func SetGlobalLogger(logger *Logger) {
	once.Do(func() {
		globalLogger = logger
	})
}

// GetGlobalLogger 获取全局日志器
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		// 如果全局日志器未设置，创建一个默认的
		defaultLogger, err := NewDefaultLogger()
		if err != nil {
			// 如果创建失败，使用zap的默认日志器
			zapLogger, _ := zap.NewProduction()
			globalLogger = &Logger{Logger: zapLogger}
		} else {
			globalLogger = defaultLogger
		}
	}
	return globalLogger
}

// Global 提供全局日志方法
type Global struct{}

// Debug 全局Debug日志
func (g Global) Debug(msg string, fields ...zap.Field) {
	GetGlobalLogger().Debug(msg, fields...)
}

// Info 全局Info日志
func (g Global) Info(msg string, fields ...zap.Field) {
	GetGlobalLogger().Info(msg, fields...)
}

// Warn 全局Warn日志
func (g Global) Warn(msg string, fields ...zap.Field) {
	GetGlobalLogger().Warn(msg, fields...)
}

// Error 全局Error日志
func (g Global) Error(msg string, fields ...zap.Field) {
	GetGlobalLogger().Error(msg, fields...)
}

// Fatal 全局Fatal日志
func (g Global) Fatal(msg string, fields ...zap.Field) {
	GetGlobalLogger().Fatal(msg, fields...)
}

// Sync 同步全局日志器
func (g Global) Sync() error {
	return GetGlobalLogger().Sync()
}

// 全局日志实例
var Log = Global{}
