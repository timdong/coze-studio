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
	"testing"
	"time"

	"etrxlite/pkg/auth/jwt"
	"etrxlite/platform/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// setupTestService 创建测试用的认证服务
func setupTestService(t *testing.T) (*Service, context.Context) {
	// 创建内存数据库客户端
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")

	// 创建测试配置
	cfg := &Config{
		JWT: &jwt.Config{
			SecretKey: "test-secret-key-for-testing",
			Method:    "HS256",
		},
	}

	// 创建测试logger
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "创建logger失败")

	// 创建认证管理器
	manager, err := NewManager(cfg, logger)
	require.NoError(t, err, "创建认证管理器失败")

	// 设置全局认证管理器
	SetGlobalAuthManager(manager)

	// 创建认证服务
	service := NewService(manager, client)

	ctx := context.Background()

	return service, ctx
}

// TestServiceRegister 测试用户注册功能
func TestServiceRegister(t *testing.T) {
	service, ctx := setupTestService(t)

	tests := []struct {
		name    string
		req     *RegisterRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功注册",
			req: &RegisterRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: false,
		},
		{
			name: "用户名已存在",
			req: &RegisterRequest{
				Username:  "testuser",
				Email:     "test2@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User2",
			},
			wantErr: true,
			errMsg:  "用户名已存在",
		},
		{
			name: "邮箱已存在",
			req: &RegisterRequest{
				Username:  "testuser2",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User3",
			},
			wantErr: true,
			errMsg:  "邮箱已存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Register(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)
				assert.NotEmpty(t, resp.Token)
				assert.Equal(t, tt.req.Username, resp.User.Username)
				assert.Equal(t, tt.req.Email, resp.User.Email)
				assert.Equal(t, tt.req.FirstName, resp.User.FirstName)
				assert.Equal(t, tt.req.LastName, resp.User.LastName)

				// 验证密码已加密
				assert.NotEqual(t, tt.req.Password, resp.User.PasswordHash)

				// 验证token可以被解析
				userInfo, err := service.ValidateToken(resp.Token)
				assert.NoError(t, err)
				assert.Equal(t, resp.User.ID, userInfo.ID)
				assert.Equal(t, resp.User.Username, userInfo.Username)
			}
		})
	}
}

// TestServiceLogin 测试用户登录功能
func TestServiceLogin(t *testing.T) {
	service, ctx := setupTestService(t)

	// 先注册一个用户
	registerReq := &RegisterRequest{
		Username:  "logintest",
		Email:     "logintest@example.com",
		Password:  "password123",
		FirstName: "Login",
		LastName:  "Test",
	}
	_, err := service.Register(ctx, registerReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		req     *LoginRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功登录",
			req: &LoginRequest{
				Username: "logintest",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "用户不存在",
			req: &LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "用户名或密码错误",
		},
		{
			name: "密码错误",
			req: &LoginRequest{
				Username: "logintest",
				Password: "wrongpassword",
			},
			wantErr: true,
			errMsg:  "用户名或密码错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Login(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)
				assert.NotEmpty(t, resp.Token)
				assert.Equal(t, tt.req.Username, resp.User.Username)
				assert.Greater(t, resp.ExpiresAt, time.Now().Unix())

				// 验证token可以被解析
				userInfo, err := service.ValidateToken(resp.Token)
				assert.NoError(t, err)
				assert.Equal(t, resp.User.ID, userInfo.ID)
				assert.Equal(t, resp.User.Username, userInfo.Username)
			}
		})
	}
}

// TestServiceRefreshToken 测试token刷新功能
func TestServiceRefreshToken(t *testing.T) {
	service, ctx := setupTestService(t)

	// 先注册并登录一个用户
	registerReq := &RegisterRequest{
		Username:  "refreshtest",
		Email:     "refreshtest@example.com",
		Password:  "password123",
		FirstName: "Refresh",
		LastName:  "Test",
	}
	registerResp, err := service.Register(ctx, registerReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "成功刷新token",
			token:   registerResp.Token,
			wantErr: false,
		},
		{
			name:    "无效的token",
			token:   "invalid.token.here",
			wantErr: true,
			errMsg:  "无效的令牌",
		},
		{
			name:    "空token",
			token:   "",
			wantErr: true,
			errMsg:  "无效的令牌",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.RefreshToken(ctx, tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)
				assert.NotEmpty(t, resp.Token)
				// 注意：新token可能与旧token相同（如果在同一秒内生成），这是正常的
				assert.Greater(t, resp.ExpiresAt, time.Now().Unix())

				// 验证新token可以被解析
				userInfo, err := service.ValidateToken(resp.Token)
				assert.NoError(t, err)
				assert.Equal(t, registerResp.User.ID, userInfo.ID)
				assert.Equal(t, registerResp.User.Username, userInfo.Username)
			}
		})
	}
}

// TestServiceChangePassword 测试修改密码功能
func TestServiceChangePassword(t *testing.T) {
	service, ctx := setupTestService(t)

	// 先注册一个用户
	registerReq := &RegisterRequest{
		Username:  "changepasstest",
		Email:     "changepasstest@example.com",
		Password:  "oldpassword123",
		FirstName: "Change",
		LastName:  "Password",
	}
	registerResp, err := service.Register(ctx, registerReq)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      uint
		oldPassword string
		newPassword string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "成功修改密码",
			userID:      registerResp.User.ID,
			oldPassword: "oldpassword123",
			newPassword: "newpassword456",
			wantErr:     false,
		},
		{
			name:        "旧密码错误",
			userID:      registerResp.User.ID,
			oldPassword: "wrongoldpassword",
			newPassword: "newpassword789",
			wantErr:     true,
			errMsg:      "旧密码错误",
		},
		{
			name:        "用户不存在",
			userID:      99999,
			oldPassword: "oldpassword123",
			newPassword: "newpassword000",
			wantErr:     true,
			errMsg:      "修改密码失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ChangePassword(ctx, tt.userID, tt.oldPassword, tt.newPassword)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)

				// 验证新密码可以登录
				loginReq := &LoginRequest{
					Username: registerResp.User.Username,
					Password: tt.newPassword,
				}
				loginResp, err := service.Login(ctx, loginReq)
				assert.NoError(t, err)
				assert.NotNil(t, loginResp)

				// 验证旧密码不能再登录
				oldLoginReq := &LoginRequest{
					Username: registerResp.User.Username,
					Password: tt.oldPassword,
				}
				_, err = service.Login(ctx, oldLoginReq)
				assert.Error(t, err)
			}
		})
	}
}

// TestServiceValidateToken 测试token验证功能
func TestServiceValidateToken(t *testing.T) {
	service, ctx := setupTestService(t)

	// 先注册一个用户
	registerReq := &RegisterRequest{
		Username:  "validatetest",
		Email:     "validatetest@example.com",
		Password:  "password123",
		FirstName: "Validate",
		LastName:  "Test",
	}
	registerResp, err := service.Register(ctx, registerReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "有效的token",
			token:   registerResp.Token,
			wantErr: false,
		},
		{
			name:    "无效的token格式",
			token:   "invalid.token",
			wantErr: true,
		},
		{
			name:    "空token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userInfo, err := service.ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, userInfo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userInfo)
				assert.Equal(t, registerResp.User.ID, userInfo.ID)
				assert.Equal(t, registerResp.User.Username, userInfo.Username)
				assert.Equal(t, registerResp.User.Email, userInfo.Email)
			}
		})
	}
}

// TestServiceIntegration 测试完整的用户认证流程
func TestServiceIntegration(t *testing.T) {
	service, ctx := setupTestService(t)

	// 1. 注册用户
	registerReq := &RegisterRequest{
		Username:  "integrationtest",
		Email:     "integration@example.com",
		Password:  "password123",
		FirstName: "Integration",
		LastName:  "Test",
	}
	registerResp, err := service.Register(ctx, registerReq)
	require.NoError(t, err)
	require.NotNil(t, registerResp)

	// 2. 验证注册返回的token
	userInfo, err := service.ValidateToken(registerResp.Token)
	require.NoError(t, err)
	assert.Equal(t, registerResp.User.ID, userInfo.ID)

	// 3. 使用用户名和密码登录
	loginReq := &LoginRequest{
		Username: "integrationtest",
		Password: "password123",
	}
	loginResp, err := service.Login(ctx, loginReq)
	require.NoError(t, err)
	require.NotNil(t, loginResp)

	// 4. 验证登录返回的token
	userInfo2, err := service.ValidateToken(loginResp.Token)
	require.NoError(t, err)
	assert.Equal(t, loginResp.User.ID, userInfo2.ID)

	// 5. 刷新token
	refreshResp, err := service.RefreshToken(ctx, loginResp.Token)
	require.NoError(t, err)
	require.NotNil(t, refreshResp)
	// 注意：新token可能与旧token相同（如果在同一秒内生成），这是正常的

	// 6. 修改密码
	err = service.ChangePassword(ctx, loginResp.User.ID, "password123", "newpassword456")
	require.NoError(t, err)

	// 7. 使用旧密码登录应该失败
	_, err = service.Login(ctx, &LoginRequest{
		Username: "integrationtest",
		Password: "password123",
	})
	assert.Error(t, err)

	// 8. 使用新密码登录应该成功
	newLoginResp, err := service.Login(ctx, &LoginRequest{
		Username: "integrationtest",
		Password: "newpassword456",
	})
	require.NoError(t, err)
	require.NotNil(t, newLoginResp)
}
