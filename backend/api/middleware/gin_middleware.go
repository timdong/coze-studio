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

package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	coreAuth "github.com/coze-dev/coze-studio/backend/core/auth"
	coreConfig "github.com/coze-dev/coze-studio/backend/core/config"
	"github.com/coze-dev/coze-studio/backend/core/logging"
	"go.uber.org/zap"
)

// GinLogger Gin 日志中间件（使用统一日志系统）
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		latency := time.Since(start)

		// 日志字段
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		// 使用统一日志系统
		logging.Log.Info("HTTP Request",
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
		)
	}
}

// GinRecovery Gin 恢复中间件
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录详细的错误信息，包括堆栈跟踪
				logging.Log.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("query", c.Request.URL.RawQuery),
				)

				// 记录错误消息
				var errMsg string
				if e, ok := err.(error); ok {
					errMsg = e.Error()
				} else {
					errMsg = fmt.Sprintf("%v", err)
				}

				c.JSON(500, gin.H{
					"code":    50000,
					"success": false,
					"message": "Internal server error: " + errMsg,
					"data":    nil,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// GinCORS Gin CORS 中间件
func GinCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// GinAuth Gin 认证中间件（JWT）
func GinAuth() gin.HandlerFunc {
	// 从配置获取 JWT 配置
	cfg := coreConfig.GetConfig()

	authMiddleware := coreAuth.NewMiddleware(&coreAuth.MiddlewareConfig{
		Secret: cfg.JWT.Secret,
		Issuer: cfg.JWT.Issuer,
	})

	return authMiddleware.AuthRequired()
}

// GinOptionalAuth Gin 可选认证中间件（JWT）
func GinOptionalAuth() gin.HandlerFunc {
	cfg := coreConfig.GetConfig()

	authMiddleware := coreAuth.NewMiddleware(&coreAuth.MiddlewareConfig{
		Secret: cfg.JWT.Secret,
		Issuer: cfg.JWT.Issuer,
	})

	return authMiddleware.OptionalAuth()
}

// GinRoleRequired Gin 角色权限中间件
func GinRoleRequired(roles ...string) gin.HandlerFunc {
	cfg := coreConfig.GetConfig()

	authMiddleware := coreAuth.NewMiddleware(&coreAuth.MiddlewareConfig{
		Secret: cfg.JWT.Secret,
		Issuer: cfg.JWT.Issuer,
	})

	return authMiddleware.RoleRequired(roles...)
}

// GinRequestID Gin 请求 ID 中间件
func GinRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成或获取请求 ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 设置到上下文和响应头
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID 生成请求 ID
func generateRequestID() string {
	// 简单实现，可以使用 UUID 或其他方式
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

