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
	"testing"
	"time"

	"etrxlite/pkg/auth/jwt"

	"go.uber.org/zap"
)

func TestHashPassword(t *testing.T) {
	// 创建测试logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 创建认证管理器
	config := &Config{
		JWT: &jwt.Config{
			SecretKey: "test-secret-key",
			Method:    "HS256",
		},
	}
	manager, err := NewManager(config, logger)
	if err != nil {
		t.Fatalf("Failed to create auth manager: %v", err)
	}

	// 测试密码加密
	password := "testpassword123"
	hashedPassword, err := manager.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashedPassword == "" {
		t.Fatal("Hashed password should not be empty")
	}

	if hashedPassword == password {
		t.Fatal("Hashed password should not be the same as original password")
	}

	// 测试密码验证
	err = manager.VerifyPassword(hashedPassword, password)
	if err != nil {
		t.Fatalf("Password verification failed: %v", err)
	}

	// 测试错误密码
	wrongPassword := "wrongpassword"
	err = manager.VerifyPassword(hashedPassword, wrongPassword)
	if err == nil {
		t.Fatal("Wrong password should fail verification")
	}
}

func TestGenerateRandomToken(t *testing.T) {
	// 创建测试logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 创建认证管理器
	config := &Config{
		JWT: &jwt.Config{
			SecretKey: "test-secret-key",
			Method:    "HS256",
		},
	}
	manager, err := NewManager(config, logger)
	if err != nil {
		t.Fatalf("Failed to create auth manager: %v", err)
	}

	// 测试随机令牌生成
	token, err := manager.GenerateRandomToken(32)
	if err != nil {
		t.Fatalf("Failed to generate random token: %v", err)
	}

	if len(token) != 64 { // 32 bytes = 64 hex characters
		t.Fatalf("Token length should be 64, got %d", len(token))
	}

	// 测试不同长度的令牌
	shortToken, err := manager.GenerateRandomToken(16)
	if err != nil {
		t.Fatalf("Failed to generate short token: %v", err)
	}

	if len(shortToken) != 32 { // 16 bytes = 32 hex characters
		t.Fatalf("Short token length should be 32, got %d", len(shortToken))
	}
}

func TestGenerateToken(t *testing.T) {
	// 创建测试logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 创建认证管理器
	config := &Config{
		JWT: &jwt.Config{
			SecretKey: "test-secret-key",
			Method:    "HS256",
		},
	}
	manager, err := NewManager(config, logger)
	if err != nil {
		t.Fatalf("Failed to create auth manager: %v", err)
	}

	// 测试JWT令牌生成
	userInfo := &UserInfo{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"user"},
	}

	token, err := manager.GenerateToken(userInfo, time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Token should not be empty")
	}

	// 测试令牌验证
	validatedUserInfo, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if validatedUserInfo.ID != userInfo.ID {
		t.Fatalf("User ID mismatch: expected %d, got %d", userInfo.ID, validatedUserInfo.ID)
	}

	if validatedUserInfo.Username != userInfo.Username {
		t.Fatalf("Username mismatch: expected %s, got %s", userInfo.Username, validatedUserInfo.Username)
	}
}

func TestGlobalAuth(t *testing.T) {
	// 创建测试logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 创建认证管理器
	config := &Config{
		JWT: &jwt.Config{
			SecretKey: "test-secret-key",
			Method:    "HS256",
		},
	}
	manager, err := NewManager(config, logger)
	if err != nil {
		t.Fatalf("Failed to create auth manager: %v", err)
	}

	// 设置全局认证管理器
	SetGlobalAuthManager(manager)

	// 测试全局认证方法
	password := "testpassword123"
	hashedPassword, err := Auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Global HashPassword failed: %v", err)
	}

	err = Auth.VerifyPassword(hashedPassword, password)
	if err != nil {
		t.Fatalf("Global VerifyPassword failed: %v", err)
	}

	token, err := Auth.GenerateRandomToken(32)
	if err != nil {
		t.Fatalf("Global GenerateRandomToken failed: %v", err)
	}

	if len(token) != 64 {
		t.Fatalf("Global token length should be 64, got %d", len(token))
	}
}

func TestTokenExpiration(t *testing.T) {
	// 创建测试logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 创建认证管理器
	config := &Config{
		JWT: &jwt.Config{
			SecretKey: "test-secret-key",
			Method:    "HS256",
		},
	}
	manager, err := NewManager(config, logger)
	if err != nil {
		t.Fatalf("Failed to create auth manager: %v", err)
	}

	userInfo := &UserInfo{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"user"},
	}

	// 生成短期令牌
	token, err := manager.GenerateToken(userInfo, time.Second*1)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 立即验证应该成功
	_, err = manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Token should be valid immediately: %v", err)
	}

	// 等待令牌过期
	time.Sleep(time.Second * 2)

	// 过期后验证应该失败
	_, err = manager.ValidateToken(token)
	if err == nil {
		t.Fatal("Expired token should fail validation")
	}
}
