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

package router

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	codegenHandler "github.com/coze-dev/coze-studio/backend/extensions/code-generation/handler"
	erdiagramHandler "github.com/coze-dev/coze-studio/backend/extensions/er-diagram/handler"
	etrxtableHandler "github.com/coze-dev/coze-studio/backend/extensions/etrxtable/handler"
	masHandler "github.com/coze-dev/coze-studio/backend/extensions/mas/handler"
	ginHandler "github.com/coze-dev/coze-studio/backend/api/handler/gin"
	"github.com/coze-dev/coze-studio/backend/api/middleware"
	"github.com/coze-dev/coze-studio/backend/internal/openapi"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// GinRegister 注册 Gin 路由
// ✅ 所有路由已从 Hertz 迁移到 Gin
func GinRegister(router *gin.Engine, db *gorm.DB) {
	// 注册 Passport 认证路由（核心路由）
	registerPassportRoutes(router)

	// 注册 Coze Studio 核心路由（✅ 迁移完成）
	registerCozeRoutes(router, db)

	// 注册扩展功能路由
	registerExtensions(router, db)

	// 注册静态文件路由
	ginStaticFileRegister(router)

	// 注册 OpenAPI 文档路由
	registerOpenAPIRoutes(router)
}

// registerPassportRoutes 注册 Passport 认证路由
func registerPassportRoutes(router *gin.Engine) {
	api := router.Group("/api")
	api.Use(middleware.GinSessionAuth()) // 添加 session 认证中间件
	{
		// Passport 路由组
		passport := api.Group("/passport")
		{
			// 账户信息
			account := passport.Group("/account")
			{
				info := account.Group("/info")
				{
					v2 := info.Group("/v2")
					v2.POST("/", ginHandler.PassportAccountInfoV2)
				}
			}

			// Web 路由组
			web := passport.Group("/web")
			{
				// 登出
				web.GET("/logout", ginHandler.PassportWebLogoutGet)

				// 邮箱相关
				email := web.Group("/email")
				{
					// 登录
					email.POST("/login", ginHandler.PassportWebEmailLoginPost)

					// 注册
					register := email.Group("/register")
					{
						register.POST("/v2", ginHandler.PassportWebEmailRegisterV2Post)
					}

					// 重置密码
					password := email.Group("/password")
					{
						password.GET("/reset", ginHandler.PassportWebEmailPasswordResetGet)
					}
				}
			}
		}

		// 用户相关路由
		{
			// 更新用户资料
			api.POST("/user/update_profile", ginHandler.UserUpdateProfile)

			// 更新头像
			web := api.Group("/web")
			{
				user := web.Group("/user")
				{
					update := user.Group("/update")
					{
						update.POST("/upload_avatar", ginHandler.UserUpdateAvatar)
					}
				}
			}
		}

	}
}

// registerExtensions 注册扩展功能路由
func registerExtensions(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api/v1")

	// 公开路由（不需要认证）
	public := api.Group("")
	{
		// 健康检查
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"code":    20000,
				"success": true,
				"message": "服务运行正常",
				"data":    gin.H{"status": "ok"},
			})
		})
	}

	// 需要认证的路由
	authenticated := api.Group("")
	authenticated.Use(middleware.GinOptionalAuth()) // 可选认证，允许未认证访问但会设置用户信息
	{
		// 多维表格系统路由
		tableHandler := etrxtableHandler.NewTableHandler(db)
		tableHandler.RegisterRoutes(authenticated)

		// MAS 多智能体系统路由
		agentHandler := masHandler.NewAgentHandler(db)
		agentHandler.RegisterRoutes(authenticated)

		sessionHandler := masHandler.NewSessionHandler(db)
		sessionHandler.RegisterRoutes(authenticated)

		taskHandler := masHandler.NewTaskHandler(db)
		taskHandler.RegisterRoutes(authenticated)

		// ER 图编辑器路由
		erDiagramHandler := erdiagramHandler.NewERDiagramHandler(db)
		erDiagramHandler.RegisterRoutes(authenticated)

		// 代码生成路由
		codeGenHandler := codegenHandler.NewCodeGenHandler(db)
		codeGenHandler.RegisterRoutes(authenticated)
	}
}

// ginStaticFileRegister 注册静态文件路由（Gin 版本）
func ginStaticFileRegister(router *gin.Engine) {
	cwd, err := os.Getwd()
	if err != nil {
		logs.Warnf("[ginStaticFileRegister] Failed to get current working directory: %v", err)
		cwd = os.Getenv("PWD")
	}

	staticFile := path.Join(cwd, "resources/static/index.html")

	// 静态文件服务
	router.Static("/static", path.Join(cwd, "/resources/static"))
	router.StaticFile("/favicon.png", "./resources/static/favicon.png")
	router.StaticFile("/", staticFile)
	router.StaticFile("/sign", staticFile)

	// Admin 静态文件
	router.StaticFS("/admin", http.Dir(path.Join(cwd, "/resources/conf")))

	// 404 处理
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") ||
			strings.HasPrefix(path, "/v1/") ||
			strings.HasPrefix(path, "/v3/") {
			c.JSON(404, gin.H{
				"code": 404,
				"msg":  "not found",
			})
			return
		}

		// 其他路径返回 index.html（SPA 路由）
		c.File(staticFile)
	})
}

// registerOpenAPIRoutes 注册 OpenAPI 文档路由
func registerOpenAPIRoutes(router *gin.Engine) {
	// 动态生成 OpenAPI JSON
	router.GET("/api/v1/openapi.json", openapi.DynamicHandler("http://localhost:8888", "1.0.0", router))

	// Swagger UI 页面
	router.GET("/docs", openapi.ServeSwaggerUI())
}

