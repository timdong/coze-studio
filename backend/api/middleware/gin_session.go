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
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coze-dev/coze-studio/backend/application/user"
	"github.com/coze-dev/coze-studio/backend/bizpkg/config"
	"github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

// 不需要 session 检查的路径（Gin 版本）
var ginNoNeedSessionCheckPath = map[string]bool{
	"/api/passport/web/email/login":        true,
	"/api/passport/web/email/register/v2":  true,
	"/api/passport/web/email/login/":       true,
	"/api/passport/web/email/register/v2/": true,
}

// 可选认证的路径（没有 session 也可以访问）
var ginOptionalAuthPath = map[string]bool{
	"/api/passport/account/info/v2": true,
}

// GinSessionAuth Gin Session 认证中间件
func GinSessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 检查是否需要 session 验证
		if ginNoNeedSessionCheckPath[path] {
			c.Next()
			return
		}

		// 从 cookie 获取 session key
		sessionKey, err := c.Cookie(entity.SessionKey)
		if err != nil || sessionKey == "" {
			// 如果是可选认证的路径，允许继续
			if ginOptionalAuthPath[path] {
				logs.Debugf("[GinSessionAuth] session key not found for optional auth path: %s", path)
				c.Next()
				return
			}
			// 否则返回 401 错误
			logs.Warnf("[GinSessionAuth] session key not found for path: %s", path)
			c.JSON(401, gin.H{
				"code":    40100,
				"success": false,
				"msg":     "authentication required: session key not found",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 验证 session
		ctx := c.Request.Context()
		// 初始化 context cache（如果还没有初始化）
		ctx = ctxcache.Init(ctx)

		session, err := user.UserApplicationSVC.ValidateSession(ctx, sessionKey)
		if err != nil {
			// 如果是可选认证的路径，允许继续
			if ginOptionalAuthPath[path] {
				logs.Debugf("[GinSessionAuth] session validation failed for optional auth path: %s, err: %v", path, err)
				c.Next()
				return
			}
			// 否则返回 401 错误
			logs.Warnf("[GinSessionAuth] validate session failed, err: %v", err)
			c.JSON(401, gin.H{
				"code":    40100,
				"success": false,
				"msg":     "authentication failed: " + err.Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 将 session 存储到 context
		if session != nil {
			ctxcache.Store(ctx, consts.SessionDataKeyInCtx, session)
			// 更新 request context
			c.Request = c.Request.WithContext(ctx)
		} else {
			// Session 验证返回 nil，说明 session 无效
			if !ginOptionalAuthPath[path] {
				logs.Warnf("[GinSessionAuth] session is nil for path: %s", path)
				c.JSON(401, gin.H{
					"code":    40100,
					"success": false,
					"msg":     "authentication failed: session is invalid",
					"data":    nil,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// GinAdminAuth Gin Admin 认证中间件
func GinAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 从 context 获取 session
		session, ok := ctxcache.Get[*entity.Session](ctx, consts.SessionDataKeyInCtx)
		if !ok {
			logs.Errorf("[GinAdminAuth] session data is nil")
			c.JSON(401, gin.H{
				"code": 40100,
				"msg":  "session data is nil",
			})
			c.Abort()
			return
		}

		// 获取基础配置
		baseConf, err := config.Base().GetBaseConfig(ctx)
		if err != nil {
			logs.Errorf("[GinAdminAuth] get base config failed, err: %v", err)
			c.JSON(500, gin.H{
				"code": 50000,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}

		// 如果未配置管理员邮箱，允许访问
		if baseConf.AdminEmails == "" {
			logs.CtxWarnf(ctx, "[GinAdminAuth] admin emails is empty")
			c.Next()
			return
		}

		// 检查用户邮箱是否在管理员列表中
		adminEmails := strings.Split(baseConf.AdminEmails, ",")
		for _, adminEmail := range adminEmails {
			if strings.EqualFold(strings.TrimSpace(adminEmail), session.UserEmail) {
				c.Next()
				return
			}
		}

		// 用户不在管理员列表中
		c.JSON(403, gin.H{
			"code": 40300,
			"msg":  "the account does not have permission to access",
		})
		c.Abort()
	}
}

