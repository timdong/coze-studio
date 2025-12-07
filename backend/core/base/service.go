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

package base

import (
	"github.com/coze-dev/coze-studio/backend/core/logging"

	"go.uber.org/zap"
)

// Service 服务基类，提供统一的日志功能
type Service struct {
	logger *logging.Logger
}

// NewService 创建服务基类
func NewService() *Service {
	return &Service{
		logger: logging.GetGlobalLogger(),
	}
}

// SetLogger 设置日志器（可选，用于测试或特殊场景）
func (s *Service) SetLogger(logger *logging.Logger) {
	s.logger = logger
}

// Logger 获取日志器
func (s *Service) Logger() *logging.Logger {
	return s.logger
}

// Debug 调试日志
func (s *Service) Debug(msg string, fields ...zap.Field) {
	s.logger.Debug(msg, fields...)
}

// Info 信息日志
func (s *Service) Info(msg string, fields ...zap.Field) {
	s.logger.Info(msg, fields...)
}

// Warn 警告日志
func (s *Service) Warn(msg string, fields ...zap.Field) {
	s.logger.Warn(msg, fields...)
}

// Error 错误日志
func (s *Service) Error(msg string, fields ...zap.Field) {
	s.logger.Error(msg, fields...)
}

// Fatal 致命错误日志
func (s *Service) Fatal(msg string, fields ...zap.Field) {
	s.logger.Fatal(msg, fields...)
}
