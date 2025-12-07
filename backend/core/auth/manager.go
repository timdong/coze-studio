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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/auth/jwt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Manager 认证管理器
type Manager struct {
	jwtClient *jwt.Client
	config    *Config
	logger    *zap.Logger
}

// Config 认证配置
type Config struct {
	JWT *jwt.Config `yaml:"jwt" json:"jwt"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// NewManager 创建认证管理器
func NewManager(config *Config, logger *zap.Logger) (*Manager, error) {
	jwtClient, err := jwt.NewClient(config.JWT)
	if err != nil {
		return nil, fmt.Errorf("failed to create jwt client: %w", err)
	}

	return &Manager{
		jwtClient: jwtClient,
		config:    config,
		logger:    logger,
	}, nil
}

// GenerateToken 生成访问令牌
func (m *Manager) GenerateToken(userInfo *UserInfo, expiration time.Duration) (string, error) {
	claims := &jwt.Claims{
		UserID:   userInfo.ID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Roles:    userInfo.Roles,
	}

	return m.jwtClient.GenerateToken(claims, expiration)
}

// GenerateTokenPair 生成令牌对
func (m *Manager) GenerateTokenPair(userInfo *UserInfo) (*TokenPair, error) {
	// 生成访问令牌
	accessToken, err := m.GenerateToken(userInfo, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成刷新令牌
	refreshToken, err := m.GenerateToken(userInfo, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

// ValidateToken 验证令牌
func (m *Manager) ValidateToken(tokenString string) (*UserInfo, error) {
	claims, err := m.jwtClient.ParseToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return &UserInfo{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Roles:    claims.Roles,
	}, nil
}

// RefreshToken 刷新令牌
func (m *Manager) RefreshToken(refreshToken string) (*TokenPair, error) {
	// 验证刷新令牌
	userInfo, err := m.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 生成新的令牌对
	return m.GenerateTokenPair(userInfo)
}

// ExtractUserInfo 从上下文中提取用户信息
func (m *Manager) ExtractUserInfo(ctx context.Context) (*UserInfo, error) {
	// 从上下文中获取令牌
	token, ok := ctx.Value("token").(string)
	if !ok {
		return nil, fmt.Errorf("no token in context")
	}

	return m.ValidateToken(token)
}

// HasRole 检查用户是否有指定角色
func (m *Manager) HasRole(userInfo *UserInfo, role string) bool {
	for _, r := range userInfo.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否有任意一个指定角色
func (m *Manager) HasAnyRole(userInfo *UserInfo, roles ...string) bool {
	for _, role := range roles {
		if m.HasRole(userInfo, role) {
			return true
		}
	}
	return false
}

// HasAllRoles 检查用户是否有所有指定角色
func (m *Manager) HasAllRoles(userInfo *UserInfo, roles ...string) bool {
	for _, role := range roles {
		if !m.HasRole(userInfo, role) {
			return false
		}
	}
	return true
}

// GetUserID 从上下文中获取用户ID
func (m *Manager) GetUserID(ctx context.Context) (uint, error) {
	userInfo, err := m.ExtractUserInfo(ctx)
	if err != nil {
		return 0, err
	}
	return userInfo.ID, nil
}

// GetUsername 从上下文中获取用户名
func (m *Manager) GetUsername(ctx context.Context) (string, error) {
	userInfo, err := m.ExtractUserInfo(ctx)
	if err != nil {
		return "", err
	}
	return userInfo.Username, nil
}

// GetUserRoles 从上下文中获取用户角色
func (m *Manager) GetUserRoles(ctx context.Context) ([]string, error) {
	userInfo, err := m.ExtractUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return userInfo.Roles, nil
}

// 全局认证管理器
var (
	globalAuthManager *Manager
)

// SetGlobalAuthManager 设置全局认证管理器
func SetGlobalAuthManager(manager *Manager) {
	globalAuthManager = manager
}

// GetGlobalAuthManager 获取全局认证管理器
func GetGlobalAuthManager() *Manager {
	return globalAuthManager
}

// 全局认证方法
var Auth = &GlobalAuth{}

// GlobalAuth 全局认证方法
type GlobalAuth struct{}

// HashPassword 全局密码加密
func (g *GlobalAuth) HashPassword(password string) (string, error) {
	if globalAuthManager == nil {
		return "", fmt.Errorf("global auth manager not initialized")
	}
	return globalAuthManager.HashPassword(password)
}

// VerifyPassword 全局密码验证
func (g *GlobalAuth) VerifyPassword(hashedPassword, password string) error {
	if globalAuthManager == nil {
		return fmt.Errorf("global auth manager not initialized")
	}
	return globalAuthManager.VerifyPassword(hashedPassword, password)
}

// GenerateRandomToken 全局随机令牌生成
func (g *GlobalAuth) GenerateRandomToken(length int) (string, error) {
	if globalAuthManager == nil {
		return "", fmt.Errorf("global auth manager not initialized")
	}
	return globalAuthManager.GenerateRandomToken(length)
}

// ValidateToken 全局令牌验证
func (g *GlobalAuth) ValidateToken(tokenString string) (*UserInfo, error) {
	if globalAuthManager == nil {
		return nil, fmt.Errorf("global auth manager not initialized")
	}
	return globalAuthManager.ValidateToken(tokenString)
}

// HashPassword 加密密码
func (m *Manager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// VerifyPassword 验证密码
func (m *Manager) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateRandomToken 生成随机令牌
func (m *Manager) GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
