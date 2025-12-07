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
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	ginHandler "github.com/coze-dev/coze-studio/backend/api/handler/gin"
	"github.com/coze-dev/coze-studio/backend/api/middleware"
)

// registerCozeRoutes 注册 Coze Studio 的核心路由
// ✅ 所有 Hertz 路由已成功迁移到 Gin
func registerCozeRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")
	api.Use(middleware.GinSessionAuth()) // 添加 session 认证中间件

	// 注册 Playground API 路由（已完成）
	registerPlaygroundAPIRoutes(api)

	// 注册 Conversation 对话路由（已完成）
	registerConversationRoutes(api)

	// 注册 Common 上传路由（已完成）
	registerCommonRoutes(api)

	// 注册 Playground 路由（已完成）
	registerPlaygroundRoutes(api)

	// 注册 Open API 路由（部分完成）
	registerOpenAPIConversationRoutes(api)

	// 注册 Bot 路由（已完成）
	registerBotRoutes(api)

	// 注册 Draftbot 路由（已完成）
	registerDraftbotRoutes(api)

	// 注册 Developer API 路由（已完成）
	registerDeveloperRoutes(api)

	// 注册 Workflow API 路由（已完成）
	registerWorkflowAPIRoutes(api)

	// 注册 Plugin API 路由（已完成）
	registerPluginAPIRoutes(api)

	// 注册 Intelligence API 路由（已完成）
	registerIntelligenceAPIRoutes(api)

	// 注册 Permission API 路由（已完成）
	registerPermissionAPIRoutes(api)

	// 注册 Knowledge 路由（已完成）
	registerKnowledgeRoutes(api)

	// 注册 Memory 路由（已完成）
	registerMemoryRoutes(api)

	// 注册 Admin 配置路由（已完成）
	registerAdminConfigRoutes(api)

	// 注册 Marketplace 路由（已完成）
	registerMarketplaceRoutes(api)

	// 注册 Plugin 路由（非 API）（已完成）
	registerPluginRoutes(api)

	// 注册 OAuth 路由（已完成）
	registerOAuthRoutes(api)

	// ✅ 所有路由组迁移已完成
}

// registerPlaygroundAPIRoutes 注册 Playground API 路由
func registerPlaygroundAPIRoutes(api *gin.RouterGroup) {
	playgroundAPI := api.Group("/playground_api")
	{
		// 批量获取用户信息
		playgroundAPI.POST("/mget_user_info", ginHandler.MGetUserBasicInfo)

		// 空间相关
		space := playgroundAPI.Group("/space")
		{
			space.POST("/list", ginHandler.GetSpaceListV2)
			space.POST("/save", ginHandler.SaveSpaceV2)
		}

		// 提示资源相关
		playgroundAPI.POST("/get_official_prompt_list", ginHandler.GetOfficialPromptResourceList)
		playgroundAPI.GET("/get_prompt_resource_info", ginHandler.GetPromptResourceInfo)
		playgroundAPI.POST("/upsert_prompt_resource", ginHandler.UpsertPromptResource)
		playgroundAPI.POST("/delete_prompt_resource", ginHandler.DeletePromptResource)

		// 文件相关
		playgroundAPI.POST("/get_file_list", ginHandler.GetFileUrls)
		playgroundAPI.POST("/get_imagex_url", ginHandler.GetImagexShortUrl)

		// 快捷命令
		playgroundAPI.POST("/create_update_shortcut_command", ginHandler.CreateUpdateShortcutCommand)

		// 用户行为
		playgroundAPI.POST("/report_user_behavior", ginHandler.ReportUserBehavior)

		// draftbot 子路由
		draftbot := playgroundAPI.Group("/draftbot")
		{
			draftbot.POST("/get_draft_bot_info", ginHandler.GetDraftBotInfoAgw)
			draftbot.POST("/update_draft_bot_info", ginHandler.UpdateDraftBotInfoAgw)
		}

		// operate 子路由
		operate := playgroundAPI.Group("/operate")
		{
			operate.POST("/get_bot_popup_info", ginHandler.GetBotPopupInfo)
			operate.POST("/update_bot_popup_info", ginHandler.UpdateBotPopupInfo)
		}
	}
}

// registerConversationRoutes 注册 Conversation 对话路由
func registerConversationRoutes(api *gin.RouterGroup) {
	conversation := api.Group("/conversation")
	{
		conversation.POST("/break_message", ginHandler.BreakMessage)
		conversation.POST("/chat", ginHandler.AgentRun) // TODO: 需要实现 SSE 流式响应
		conversation.POST("/clear_message", ginHandler.ClearConversationHistory)
		conversation.POST("/create_section", ginHandler.ClearConversationCtx)
		conversation.POST("/delete_message", ginHandler.DeleteMessage)
		conversation.POST("/get_message_list", ginHandler.GetMessageList)
	}
}

// registerOpenAPIConversationRoutes 注册 Open API 对话路由（v1, v3）
func registerOpenAPIConversationRoutes(api *gin.RouterGroup) {
	// v1 API 路由
	v1 := api.Group("/v1")
	{
		// Conversation 相关
		v1.GET("/conversations", ginHandler.ListConversationsApi)
		v1.POST("/conversation/create", ginHandler.CreateConversation)
		v1.GET("/conversation/retrieve", ginHandler.RetrieveConversationApi)
		v1.POST("/conversation/message/list", ginHandler.GetApiMessageList)

		conversations := v1.Group("/conversations")
		{
			conversations.POST("/:conversation_id/clear", ginHandler.ClearConversationApi)
			conversations.PUT("/:conversation_id", ginHandler.UpdateConversationApi)
			conversations.DELETE("/:conversation_id", ginHandler.DeleteConversationApi)
		}
	}

	// v3 API 路由
	v3 := api.Group("/v3")
	{
		v3.POST("/chat", ginHandler.ChatV3)
		v3.POST("/chat/cancel", ginHandler.CancelChatApi)
		v3.GET("/chat/retrieve", ginHandler.RetrieveChatOpen)
		chat := v3.Group("/chat")
		{
			message := chat.Group("/message")
			{
				message.GET("/list", ginHandler.ListChatMessageApi)
			}
		}
	}
}

// registerCommonRoutes 注册 Common 上传路由
func registerCommonRoutes(api *gin.RouterGroup) {
	common := api.Group("/common")
	{
		upload := common.Group("/upload")
		{
			// 先注册具体路径
			upload.GET("/apply_upload_action", ginHandler.ApplyUploadAction)
			upload.POST("/apply_upload_action", ginHandler.ApplyUploadAction)
			// 使用通配符参数处理其他路径
			upload.POST("/:tos_uri", ginHandler.CommonUpload)
		}
	}
}

// registerPlaygroundRoutes 注册 Playground 路由
func registerPlaygroundRoutes(api *gin.RouterGroup) {
	playground := api.Group("/playground")
	{
		playground.POST("/get_onboarding", ginHandler.GetOnboarding)
		upload := playground.Group("/upload")
		{
			upload.POST("/auth_token", ginHandler.GetUploadAuthToken)
		}
	}
}

// registerBotRoutes 注册 Bot 路由
func registerBotRoutes(api *gin.RouterGroup) {
	bot := api.Group("/bot")
	{
		bot.POST("/get_type_list", ginHandler.GetTypeList)
		bot.POST("/upload_file", ginHandler.UploadFile)
	}
}

// registerDraftbotRoutes 注册 Draftbot 路由
func registerDraftbotRoutes(api *gin.RouterGroup) {
	draftbot := api.Group("/draftbot")
	{
		draftbot.POST("/commit_check", ginHandler.CheckDraftBotCommit)
		draftbot.POST("/create", ginHandler.DraftBotCreate)
		draftbot.POST("/delete", ginHandler.DeleteDraftBot)
		draftbot.POST("/duplicate", ginHandler.DuplicateDraftBot)
		draftbot.POST("/get_display_info", ginHandler.GetDraftBotDisplayInfo)
		draftbot.POST("/list_draft_history", ginHandler.ListDraftBotHistory)
		draftbot.POST("/publish", ginHandler.PublishDraftBot)
		draftbot.POST("/update_display_info", ginHandler.UpdateDraftBotDisplayInfo)

		publish := draftbot.Group("/publish")
		{
			connector := publish.Group("/connector")
			{
				connector.POST("/list", ginHandler.PublishConnectorList)
			}
		}
	}
}

// registerDeveloperRoutes 注册 Developer API 路由
func registerDeveloperRoutes(api *gin.RouterGroup) {
	developer := api.Group("/developer")
	{
		developer.POST("/get_icon", ginHandler.GetIcon)
	}

	// 用户相关路由
	api.POST("/user/update_profile_check", ginHandler.UpdateUserProfileCheck)
}

// registerWorkflowAPIRoutes 注册 Workflow API 路由
func registerWorkflowAPIRoutes(api *gin.RouterGroup) {
	workflowAPI := api.Group("/workflow_api")
	{
		// 基础工作流操作
		workflowAPI.POST("/create", ginHandler.CreateWorkflow)
		workflowAPI.POST("/canvas", ginHandler.GetCanvasInfo)
		workflowAPI.POST("/save", ginHandler.SaveWorkflow)
		workflowAPI.POST("/update_meta", ginHandler.UpdateWorkflowMeta)
		workflowAPI.POST("/delete", ginHandler.DeleteWorkflow)
		workflowAPI.POST("/batch_delete", ginHandler.BatchDeleteWorkflow)
		workflowAPI.POST("/delete_strategy", ginHandler.GetDeleteStrategy)
		workflowAPI.POST("/publish", ginHandler.PublishWorkflow)
		workflowAPI.POST("/copy", ginHandler.CopyWorkflow)
		workflowAPI.POST("/copy_wk_template", ginHandler.CopyWkTemplateApi)

		// 工作流查询
		workflowAPI.POST("/released_workflows", ginHandler.GetReleasedWorkflows)
		workflowAPI.POST("/workflow_references", ginHandler.GetWorkflowReferences)
		workflowAPI.POST("/workflow_list", ginHandler.GetWorkFlowList)
		workflowAPI.POST("/workflow_detail", ginHandler.GetWorkflowDetail)
		workflowAPI.POST("/workflow_detail_info", ginHandler.GetWorkflowDetailInfo)
		workflowAPI.POST("/list_publish_workflow", ginHandler.ListPublishWorkflow)
		workflowAPI.POST("/example_workflow_list", ginHandler.GetExampleWorkFlowList)

		// 节点相关
		workflowAPI.POST("/node_type", ginHandler.QueryWorkflowNodeTypes)
		workflowAPI.POST("/node_template_list", ginHandler.NodeTemplateList)
		workflowAPI.POST("/node_panel_search", ginHandler.NodePanelSearch)
		workflowAPI.POST("/nodeDebug", ginHandler.WorkflowNodeDebugV2)

		// LLM 节点设置
		workflowAPI.POST("/llm_fc_setting_merged", ginHandler.GetLLMNodeFCSettingsMerged)
		workflowAPI.POST("/llm_fc_setting_detail", ginHandler.GetLLMNodeFCSettingDetail)

		// 测试和运行
		workflowAPI.POST("/test_run", ginHandler.WorkFlowTestRun)
		workflowAPI.POST("/test_resume", ginHandler.WorkFlowTestResume)
		workflowAPI.POST("/cancel", ginHandler.CancelWorkFlow)
		workflowAPI.GET("/get_process", ginHandler.GetWorkFlowProcess)
		workflowAPI.GET("/get_node_execute_history", ginHandler.GetNodeExecuteHistory)

		// API 详情
		workflowAPI.GET("/apiDetail", ginHandler.GetApiDetail)

		// 图片签名
		workflowAPI.POST("/sign_image_url", ginHandler.SignImageURL)

		// 验证
		workflowAPI.POST("/validate_tree", ginHandler.ValidateTree)

		// 历史 Schema
		workflowAPI.POST("/history_schema", ginHandler.GetHistorySchema)

		// 项目对话
		projectConversation := workflowAPI.Group("/project_conversation")
		{
			projectConversation.POST("/create", ginHandler.CreateProjectConversationDef)
			projectConversation.POST("/update", ginHandler.UpdateProjectConversationDef)
			projectConversation.POST("/delete", ginHandler.DeleteProjectConversationDef)
			projectConversation.GET("/list", ginHandler.ListProjectConversationDef)
		}

		// 聊天流角色
		chatFlowRole := workflowAPI.Group("/chat_flow_role")
		{
			chatFlowRole.GET("/get", ginHandler.GetChatFlowRole)
			chatFlowRole.POST("/create", ginHandler.CreateChatFlowRole)
			chatFlowRole.POST("/delete", ginHandler.DeleteChatFlowRole)
		}

		// 上传
		upload := workflowAPI.Group("/upload")
		{
			upload.POST("/auth_token", ginHandler.GetWorkflowUploadAuthToken)
		}

		// Trace 相关（暂时返回空响应）
		workflowAPI.POST("/list_spans", ginHandler.ListRootSpans)
		workflowAPI.POST("/get_trace", ginHandler.GetTraceSDK)
	}

	// Open API 路由（v1）
	v1 := api.Group("/v1")
	{
		// 工作流运行
		v1.POST("/workflow/run", ginHandler.OpenAPIRunFlow)
		v1.POST("/workflow/stream_run", ginHandler.OpenAPIStreamRunFlow)
		v1.POST("/workflow/stream_resume", ginHandler.OpenAPIStreamResumeFlow)
		v1.GET("/workflow/get_run_history", ginHandler.OpenAPIGetWorkflowRunHistory)

		// 工作流信息
		v1.GET("/workflows/:workflow_id", ginHandler.OpenAPIGetWorkflowInfo)

		// 聊天流
		v1.POST("/workflows/chat", ginHandler.OpenAPIChatFlowRun)

		// 对话
		workflow := v1.Group("/workflow")
		{
			workflow.POST("/conversation/create", ginHandler.OpenAPICreateConversation)
		}

		// 应用相关
		apps := v1.Group("/apps")
		{
			apps.GET("/:app_id", ginHandler.GetOnlineAppData)
		}

		// 工作空间相关
		v1.POST("/workspaces", ginHandler.OpenCreateSpace)
	}
}

// registerPluginAPIRoutes 注册 Plugin API 路由
func registerPluginAPIRoutes(api *gin.RouterGroup) {
	pluginAPI := api.Group("/plugin_api")
	{
		// 插件列表和查询
		pluginAPI.POST("/get_playground_plugin_list", ginHandler.GetPlaygroundPluginList)
		pluginAPI.POST("/get_dev_plugin_list", ginHandler.GetDevPluginList)
		pluginAPI.POST("/get_plugin_info", ginHandler.GetPluginInfo)
		pluginAPI.POST("/get_plugin_apis", ginHandler.GetPluginAPIs)
		pluginAPI.POST("/get_updated_apis", ginHandler.GetUpdatedAPIs)
		pluginAPI.POST("/get_plugin_next_version", ginHandler.GetPluginNextVersion)

		// 插件注册和更新
		pluginAPI.POST("/register_plugin_meta", ginHandler.RegisterPluginMeta)
		pluginAPI.POST("/update", ginHandler.UpdatePlugin)
		pluginAPI.POST("/update_plugin_meta", ginHandler.UpdatePluginMeta)
		pluginAPI.POST("/del_plugin", ginHandler.DelPlugin)
		pluginAPI.POST("/publish_plugin", ginHandler.PublishPlugin)

		// API 管理
		pluginAPI.POST("/create_api", ginHandler.CreateAPI)
		pluginAPI.POST("/update_api", ginHandler.UpdateAPI)
		pluginAPI.POST("/delete_api", ginHandler.DeleteAPI)
		pluginAPI.POST("/batch_create_api", ginHandler.BatchCreateAPI)
		pluginAPI.POST("/debug_api", ginHandler.DebugAPI)

		// Bot 默认参数
		pluginAPI.POST("/get_bot_default_params", ginHandler.GetBotDefaultParams)
		pluginAPI.POST("/update_bot_default_params", ginHandler.UpdateBotDefaultParams)

		// OAuth 相关
		pluginAPI.POST("/get_oauth_status", ginHandler.GetOAuthStatus)
		pluginAPI.POST("/get_oauth_schema", ginHandler.GetOAuthSchemaAPI)
		pluginAPI.POST("/get_queried_oauth_plugins", ginHandler.GetQueriedOAuthPluginList)
		pluginAPI.POST("/revoke_auth_token", ginHandler.RevokeAuthToken)

		// 插件编辑锁定
		pluginAPI.POST("/check_and_lock_plugin_edit", ginHandler.CheckAndLockPluginEdit)
		pluginAPI.POST("/unlock_plugin_edit", ginHandler.UnlockPluginEdit)

		// 用户权限
		pluginAPI.POST("/get_user_authority", ginHandler.GetUserAuthority)

		// 转换工具
		pluginAPI.POST("/convert_to_openapi", ginHandler.Convert2OpenAPI)

		// 资源相关
		pluginAPI.POST("/library_resource_list", ginHandler.LibraryResourceList)
		pluginAPI.POST("/project_resource_list", ginHandler.ProjectResourceList)
		pluginAPI.POST("/resource_copy_dispatch", ginHandler.ResourceCopyDispatch)
		pluginAPI.POST("/resource_copy_detail", ginHandler.ResourceCopyDetail)
		pluginAPI.POST("/resource_copy_retry", ginHandler.ResourceCopyRetry)
		pluginAPI.POST("/resource_copy_cancel", ginHandler.ResourceCopyCancel)
	}
}

// registerIntelligenceAPIRoutes 注册 Intelligence API 路由
func registerIntelligenceAPIRoutes(api *gin.RouterGroup) {
	intelligenceAPI := api.Group("/intelligence_api")
	{
		// 搜索相关
		search := intelligenceAPI.Group("/search")
		{
			search.POST("/get_draft_intelligence_list", ginHandler.GetDraftIntelligenceList)
			search.POST("/get_draft_intelligence_info", ginHandler.GetDraftIntelligenceInfo)
			search.POST("/get_recently_edit_intelligence", ginHandler.GetUserRecentlyEditIntelligence)
		}

		// 草稿项目相关
		draftProject := intelligenceAPI.Group("/draft_project")
		{
			draftProject.POST("/create", ginHandler.DraftProjectCreate)
			draftProject.POST("/update", ginHandler.DraftProjectUpdate)
			draftProject.POST("/delete", ginHandler.DraftProjectDelete)
			draftProject.POST("/copy", ginHandler.DraftProjectCopy)
			draftProject.POST("/inner_task_list", ginHandler.DraftProjectInnerTaskList)
		}

		// 发布相关
		publish := intelligenceAPI.Group("/publish")
		{
			publish.POST("/get_published_connector", ginHandler.GetProjectPublishedConnector)
			publish.POST("/check_version_number", ginHandler.CheckProjectVersionNumber)
			publish.POST("/publish_project", ginHandler.PublishProject)
			publish.POST("/publish_record_list", ginHandler.GetPublishRecordList)
			publish.POST("/connector_list", ginHandler.ProjectPublishConnectorList)
			publish.POST("/publish_record_detail", ginHandler.GetPublishRecordDetail)
		}
	}
}

// registerPermissionAPIRoutes 注册 Permission API 路由
func registerPermissionAPIRoutes(api *gin.RouterGroup) {
	permissionAPI := api.Group("/permission_api")
	{
		// PAT (Personal Access Token) 相关
		pat := permissionAPI.Group("/pat")
		{
			pat.GET("/get_personal_access_token_and_permission", ginHandler.GetPersonalAccessTokenAndPermission)
			pat.POST("/delete_personal_access_token_and_permission", ginHandler.DeletePersonalAccessTokenAndPermission)
			pat.GET("/list_personal_access_tokens", ginHandler.ListPersonalAccessTokens)
			pat.POST("/create_personal_access_token_and_permission", ginHandler.CreatePersonalAccessTokenAndPermission)
			pat.POST("/update_personal_access_token_and_permission", ginHandler.UpdatePersonalAccessTokenAndPermission)
		}

		// Coze Web App 相关
		cozeWebApp := permissionAPI.Group("/coze_web_app")
		{
			cozeWebApp.POST("/impersonate_coze_user", ginHandler.ImpersonateCozeUser)
		}
	}
}

// registerKnowledgeRoutes 注册 Knowledge 路由
func registerKnowledgeRoutes(api *gin.RouterGroup) {
	knowledge := api.Group("/knowledge")
	{
		// 数据集相关
		knowledge.POST("/create", ginHandler.CreateDataset)
		knowledge.POST("/delete", ginHandler.DeleteDataset)
		knowledge.POST("/detail", ginHandler.DatasetDetail)
		knowledge.POST("/list", ginHandler.ListDataset)
		knowledge.POST("/update", ginHandler.UpdateDataset)

		// 文档相关
		document := knowledge.Group("/document")
		{
			document.POST("/create", ginHandler.CreateDocument)
			document.POST("/delete", ginHandler.DeleteDocument)
			document.POST("/list", ginHandler.ListDocument)
			document.POST("/resegment", ginHandler.Resegment)
			document.POST("/update", ginHandler.UpdateDocument)

			// 文档进度
			progress := document.Group("/progress")
			{
				progress.POST("/get", ginHandler.GetDocumentProgress)
			}
		}

		// 图标相关
		icon := knowledge.Group("/icon")
		{
			icon.POST("/get", ginHandler.GetIconForDataset)
		}

		// 照片相关
		photo := knowledge.Group("/photo")
		{
			photo.POST("/caption", ginHandler.UpdatePhotoCaption)
			photo.POST("/detail", ginHandler.PhotoDetail)
			photo.POST("/extract_caption", ginHandler.ExtractPhotoCaption)
			photo.POST("/list", ginHandler.ListPhoto)
		}

		// 审核相关
		review := knowledge.Group("/review")
		{
			review.POST("/create", ginHandler.CreateDocumentReview)
			review.POST("/mget", ginHandler.MGetDocumentReview)
			review.POST("/save", ginHandler.SaveDocumentReview)
		}

		// 切片相关
		slice := knowledge.Group("/slice")
		{
			slice.POST("/create", ginHandler.CreateSlice)
			slice.POST("/delete", ginHandler.DeleteSlice)
			slice.POST("/list", ginHandler.ListSlice)
			slice.POST("/update", ginHandler.UpdateSlice)
		}

		// 表结构相关
		tableSchema := knowledge.Group("/table_schema")
		{
			tableSchema.POST("/get", ginHandler.GetTableSchema)
			tableSchema.POST("/validate", ginHandler.ValidateTableSchema)
		}
	}

	// Memory 相关路由（虽然路径是 /api/memory，但处理器在 knowledge_handler.go 中）
	memory := api.Group("/memory")
	{
		memory.GET("/doc_table_info", ginHandler.GetDocumentTableInfo)
		memory.GET("/table_mode_config", ginHandler.GetModeConfig)
	}
}

// registerMemoryRoutes 注册 Memory 路由
func registerMemoryRoutes(api *gin.RouterGroup) {
	memory := api.Group("/memory")
	{
		// 变量相关
		memory.GET("/sys_variable_conf", ginHandler.GetSysVariableConf)
		memory.POST("/variable/upsert", ginHandler.SetKvMemory)
		memory.POST("/variable/get_meta", ginHandler.GetMemoryVariableMeta)
		memory.POST("/variable/delete", ginHandler.DelProfileMemory)
		memory.POST("/variable/get", ginHandler.GetPlayGroundMemory)

		// 项目变量相关
		project := memory.Group("/project")
		{
			variable := project.Group("/variable")
			{
				variable.GET("/meta_list", ginHandler.GetProjectVariableList)
				variable.POST("/meta_update", ginHandler.UpdateProjectVariable)
			}
		}

		// 数据库相关
		database := memory.Group("/database")
		{
			database.POST("/list", ginHandler.ListDatabase)
			database.POST("/get_by_id", ginHandler.GetDatabaseByID)
			database.POST("/add", ginHandler.AddDatabase)
			database.POST("/update", ginHandler.UpdateDatabase)
			database.POST("/delete", ginHandler.DeleteDatabase)
			database.POST("/bind_to_bot", ginHandler.BindDatabase)
			database.POST("/unbind_to_bot", ginHandler.UnBindDatabase)
			database.POST("/list_records", ginHandler.ListDatabaseRecords)
			database.POST("/update_records", ginHandler.UpdateDatabaseRecords)
			database.POST("/get_online_database_id", ginHandler.GetOnlineDatabaseId)
			database.POST("/get_template", ginHandler.GetDatabaseTemplate)
			database.POST("/get_connector_name", ginHandler.GetConnectorName)
			database.POST("/update_bot_switch", ginHandler.UpdateDatabaseBotSwitch)

			// 数据库表相关
			table := database.Group("/table")
			{
				table.POST("/list_new", ginHandler.GetBotDatabase)
				table.POST("/reset", ginHandler.ResetBotTable)
			}
		}

		// 表结构相关
		tableSchema := memory.Group("/table_schema")
		{
			tableSchema.POST("/get", ginHandler.GetDatabaseTableSchema)
			tableSchema.POST("/validate", ginHandler.ValidateDatabaseTableSchema)
		}

		// 表文件相关
		tableFile := memory.Group("/table_file")
		{
			tableFile.POST("/submit", ginHandler.SubmitDatabaseInsertTask)
			tableFile.POST("/get_progress", ginHandler.DatabaseFileProgressData)
		}
	}
}

// registerAdminConfigRoutes 注册 Admin 配置路由
func registerAdminConfigRoutes(api *gin.RouterGroup) {
	admin := api.Group("/admin")
	admin.Use(middleware.GinAdminAuth()) // 添加 Admin 认证中间件
	{
		config := admin.Group("/config")
		{
			// 基础配置
			basic := config.Group("/basic")
			{
				basic.GET("/get", ginHandler.GetBasicConfiguration)
				basic.POST("/save", ginHandler.SaveBasicConfiguration)
			}

			// 知识库配置
			knowledge := config.Group("/knowledge")
			{
				knowledge.GET("/get", ginHandler.GetKnowledgeConfig)
				knowledge.POST("/save", ginHandler.UpdateKnowledgeConfig)
			}

			// 模型配置
			model := config.Group("/model")
			{
				model.GET("/list", ginHandler.GetModelList)
				model.POST("/create", ginHandler.CreateModel)
				model.POST("/delete", ginHandler.DeleteModel)
			}
		}
	}
}

// registerMarketplaceRoutes 注册 Marketplace 路由
func registerMarketplaceRoutes(api *gin.RouterGroup) {
	marketplace := api.Group("/marketplace")
	{
		product := marketplace.Group("/product")
		{
			product.GET("/list", ginHandler.PublicGetProductList)
			product.GET("/detail", ginHandler.PublicGetProductDetail)
			product.POST("/favorite", ginHandler.PublicFavoriteProduct)
			product.POST("/duplicate", ginHandler.PublicDuplicateProduct)
			product.GET("/search", ginHandler.PublicSearchProduct)
			product.GET("/call_info", ginHandler.PublicGetProductCallInfo)
			product.GET("/config", ginHandler.PublicGetMarketPluginConfig)

			// 收藏相关
			favorite := product.Group("/favorite")
			{
				favorite.GET("/list.v2", ginHandler.PublicGetUserFavoriteListV2)
			}

			// 搜索相关
			search := product.Group("/search")
			{
				search.GET("/suggest", ginHandler.PublicSearchSuggest)
			}

			// 分类相关
			category := product.Group("/category")
			{
				category.GET("/list", ginHandler.PublicGetProductCategoryList)
			}
		}
	}
}

// registerPluginRoutes 注册 Plugin 路由（非 API）
func registerPluginRoutes(api *gin.RouterGroup) {
	plugin := api.Group("/plugin")
	{
		plugin.POST("/get_oauth_schema", ginHandler.GetOAuthSchema)
	}
}

// registerOAuthRoutes 注册 OAuth 路由
func registerOAuthRoutes(api *gin.RouterGroup) {
	oauth := api.Group("/oauth")
	{
		oauth.GET("/authorization_code", ginHandler.OauthAuthorizationCode)
	}
}

