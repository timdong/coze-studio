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

package ctxutil

import (
	"context"
)

type contextKey string

const (
	// UserIDKey 用户ID的context key
	UserIDKey contextKey = "user_id"
	// UsernameKey 用户名的context key
	UsernameKey contextKey = "username"
	// WorkspaceIDKey 工作空间ID的context key
	WorkspaceIDKey contextKey = "workspace_id"
	// RequestIDKey 请求ID的context key
	RequestIDKey contextKey = "request_id"
)

// WithUserID 在context中添加用户ID
func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID 从context中获取用户ID
// 返回用户ID和是否存在的标志
func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

// MustGetUserID 从context中获取用户ID，如果不存在则返回0
func MustGetUserID(ctx context.Context) uint {
	userID, _ := GetUserID(ctx)
	return userID
}

// WithUsername 在context中添加用户名
func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, UsernameKey, username)
}

// GetUsername 从context中获取用户名
// 返回用户名和是否存在的标志
func GetUsername(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}

// MustGetUsername 从context中获取用户名，如果不存在则返回空字符串
func MustGetUsername(ctx context.Context) string {
	username, _ := GetUsername(ctx)
	return username
}

// WithWorkspaceID 在context中添加工作空间ID
func WithWorkspaceID(ctx context.Context, workspaceID uint) context.Context {
	return context.WithValue(ctx, WorkspaceIDKey, workspaceID)
}

// GetWorkspaceID 从context中获取工作空间ID
// 返回工作空间ID和是否存在的标志
func GetWorkspaceID(ctx context.Context) (uint, bool) {
	workspaceID, ok := ctx.Value(WorkspaceIDKey).(uint)
	return workspaceID, ok
}

// MustGetWorkspaceID 从context中获取工作空间ID，如果不存在则返回0
func MustGetWorkspaceID(ctx context.Context) uint {
	workspaceID, _ := GetWorkspaceID(ctx)
	return workspaceID
}

// WithRequestID 在context中添加请求ID
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID 从context中获取请求ID
// 返回请求ID和是否存在的标志
func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	return requestID, ok
}

// MustGetRequestID 从context中获取请求ID，如果不存在则返回空字符串
func MustGetRequestID(ctx context.Context) string {
	requestID, _ := GetRequestID(ctx)
	return requestID
}

// FromGinContext 从Gin的context中提取信息并设置到Go的context中
// 这是一个辅助函数，用于在中间件中将Gin的context转换为Go的context
func FromGinContext(ctx context.Context, userID interface{}, username interface{}, workspaceID interface{}) context.Context {
	if userID != nil {
		if uid, ok := userID.(int); ok {
			ctx = WithUserID(ctx, uint(uid))
		} else if uid, ok := userID.(uint); ok {
			ctx = WithUserID(ctx, uid)
		}
	}
	if username != nil {
		if uname, ok := username.(string); ok {
			ctx = WithUsername(ctx, uname)
		}
	}
	if workspaceID != nil {
		if wid, ok := workspaceID.(int); ok {
			ctx = WithWorkspaceID(ctx, uint(wid))
		} else if wid, ok := workspaceID.(uint); ok {
			ctx = WithWorkspaceID(ctx, wid)
		}
	}
	return ctx
}

