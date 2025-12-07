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

package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/core/logging"
	"gorm.io/gorm"
	"github.com/coze-dev/coze-studio/backend/extensions/mas/models"

	"go.uber.org/zap"
)

// Service 统一认证服务
type Service struct {
	manager *Manager
	db      *gorm.DB
}

// NewService 创建统一认证服务
func NewService(manager *Manager, db *gorm.DB) *Service {
	return &Service{
		manager: manager,
		db:      db,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string        `json:"token"`
	User      *models.User  `json:"user"`
	ExpiresAt int64         `json:"expires_at"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// Login 用户登录
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	logging.Log.Info("用户登录尝试", zap.String("username", req.Username))

	// 查找用户（使用 GORM）
	var user models.User
	if err := s.db.WithContext(ctx).Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logging.Log.Warn("用户不存在", zap.String("username", req.Username))
			return nil, fmt.Errorf("用户名或密码错误")
		}
		logging.Log.Error("查询用户失败", zap.Error(err))
		return nil, fmt.Errorf("登录失败")
	}

	// 验证密码（假设 User 模型有 PasswordHash 字段）
	// TODO: 需要根据实际的 User 模型字段调整
	if err := s.manager.VerifyPassword(user.Username, req.Password); err != nil {
		logging.Log.Warn("密码验证失败", zap.String("username", req.Username))
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 生成JWT令牌
	userInfo := &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    "", // TODO: 根据实际模型字段调整
		Roles:    []string{"user"}, // 默认角色
	}

	token, err := s.manager.GenerateToken(userInfo, 24*time.Hour)
	if err != nil {
		logging.Log.Error("生成令牌失败", zap.Error(err))
		return nil, fmt.Errorf("登录失败")
	}

	logging.Log.Info("用户登录成功", zap.String("username", req.Username), zap.Uint("user_id", user.ID))

	return &LoginResponse{
		Token:     token,
		User:      &user,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

// Register 用户注册
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	logging.Log.Info("用户注册尝试", zap.String("username", req.Username), zap.String("email", req.Email))

	// 检查用户名是否已存在（使用 GORM）
	var count int64
	if err := s.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		logging.Log.Error("检查用户名失败", zap.Error(err))
		return nil, fmt.Errorf("注册失败")
	}
	if count > 0 {
		return nil, fmt.Errorf("用户名已存在")
	}

	// TODO: 实现邮箱检查和用户创建逻辑
	// 这里暂时返回错误，需要根据实际的 User 模型字段来实现
	return nil, fmt.Errorf("用户注册功能暂未实现，需要根据实际 User 模型调整")
}

// ValidateToken 验证令牌
func (s *Service) ValidateToken(tokenString string) (*UserInfo, error) {
	return s.manager.ValidateToken(tokenString)
}

// RefreshToken 刷新令牌
func (s *Service) RefreshToken(ctx context.Context, tokenString string) (*LoginResponse, error) {
	// 验证当前令牌
	userInfo, err := s.manager.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("无效的令牌")
	}

	// 查找用户（使用 GORM）
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userInfo.ID).Error; err != nil {
		logging.Log.Error("查找用户失败", zap.Error(err))
		return nil, fmt.Errorf("刷新令牌失败")
	}

	// 生成新令牌
	newToken, err := s.manager.GenerateToken(userInfo, 24*time.Hour)
	if err != nil {
		logging.Log.Error("生成新令牌失败", zap.Error(err))
		return nil, fmt.Errorf("刷新令牌失败")
	}

	return &LoginResponse{
		Token:     newToken,
		User:      &user,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

// ChangePassword 修改密码
func (s *Service) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	logging.Log.Info("用户修改密码", zap.Uint("user_id", userID))

	// 查找用户（使用 GORM）
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		logging.Log.Error("查找用户失败", zap.Error(err))
		return fmt.Errorf("修改密码失败")
	}

	// TODO: 实现密码验证和更新逻辑，需要根据实际的 User 模型字段调整
	// 这里暂时返回错误
	return fmt.Errorf("修改密码功能暂未实现，需要根据实际 User 模型调整")
}
