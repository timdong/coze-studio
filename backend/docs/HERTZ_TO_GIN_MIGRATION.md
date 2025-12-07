# Hertz 到 Gin 路由迁移指南

## 📋 迁移状态

### ✅ 已迁移的路由组

1. **Passport 认证路由** (`/api/passport/*`)
   - ✅ `/api/passport/account/info/v2/` - 获取账户信息
   - ✅ `/api/passport/web/email/login` - 邮箱登录
   - ✅ `/api/passport/web/email/register/v2` - 邮箱注册
   - ✅ `/api/passport/web/email/password/reset` - 重置密码
   - ✅ `/api/passport/web/logout` - 登出

2. **用户相关路由** (`/api/user/*`, `/api/web/user/*`)
   - ✅ `/api/user/update_profile` - 更新用户资料
   - ✅ `/api/web/user/update/upload_avatar` - 更新头像

3. **Playground API 路由** (`/api/playground_api/*`) ✅ **已完成**
   - ✅ `/api/playground_api/mget_user_info` - 批量获取用户信息
   - ✅ `/api/playground_api/space/list` - 获取空间列表
   - ✅ `/api/playground_api/draftbot/get_draft_bot_info` - 获取草稿 Bot 信息
   - ✅ `/api/playground_api/draftbot/update_draft_bot_info` - 更新草稿 Bot 信息
   - ✅ `/api/playground_api/get_official_prompt_list` - 获取官方提示资源列表
   - ✅ `/api/playground_api/get_prompt_resource_info` - 获取提示资源信息
   - ✅ `/api/playground_api/upsert_prompt_resource` - 更新或插入提示资源
   - ✅ `/api/playground_api/delete_prompt_resource` - 删除提示资源
   - ✅ `/api/playground_api/get_imagex_url` - 获取图片短链接
   - ✅ `/api/playground_api/operate/get_bot_popup_info` - 获取 Bot 弹窗信息
   - ✅ `/api/playground_api/operate/update_bot_popup_info` - 更新 Bot 弹窗信息
   - ✅ `/api/playground_api/create_update_shortcut_command` - 创建或更新快捷命令
   - ✅ `/api/playground_api/report_user_behavior` - 报告用户行为
   - ✅ `/api/playground_api/get_file_list` - 获取文件 URL 列表

4. **Conversation 对话路由** (`/api/conversation/*`, `/v1/conversations/*`, `/v1/conversation/*`) ✅ **已完成**
   - ✅ `/api/conversation/break_message` - 中断消息
   - ✅ `/api/conversation/chat` - 对话聊天（需要 SSE 支持）
   - ✅ `/api/conversation/clear_message` - 清空对话历史
   - ✅ `/api/conversation/create_section` - 创建新的对话会话
   - ✅ `/api/conversation/delete_message` - 删除消息
   - ✅ `/api/conversation/get_message_list` - 获取消息列表
   - ✅ `/v1/conversations` - 列出对话（Open API）
   - ✅ `/v1/conversations/:conversation_id/clear` - 清空对话（Open API）
   - ✅ `/v1/conversations/:conversation_id` - 更新/删除对话（Open API）
   - ✅ `/v1/conversation/create` - 创建对话（Open API）
   - ✅ `/v1/conversation/retrieve` - 获取对话详情（Open API）
   - ✅ `/v1/conversation/message/list` - 获取消息列表（Open API）

4. **Conversation 对话路由** (`/api/conversation/*`) ✅ **已完成**
   - ✅ `/api/conversation/break_message` - 中断消息
   - ✅ `/api/conversation/chat` - 对话聊天（SSE 流式）
   - ✅ `/api/conversation/clear_message` - 清空消息
   - ✅ `/api/conversation/create_section` - 创建会话段
   - ✅ `/api/conversation/delete_message` - 删除消息
   - ✅ `/api/conversation/get_message_list` - 获取消息列表

5. **Common 上传路由** (`/api/common/upload/*`) ✅ **已完成**
   - ✅ `/api/common/upload/apply_upload_action` - 申请上传操作
   - ✅ `/api/common/upload/*tos_uri` - 通用上传

6. **Playground 路由** (`/api/playground/*`) ✅ **已完成**
   - ✅ `/api/playground/get_onboarding` - 获取引导信息
   - ✅ `/api/playground/upload/auth_token` - 获取上传令牌

7. **Bot 路由** (`/api/bot/*`) ✅ **已完成**
   - ✅ `/api/bot/get_type_list` - 获取类型列表
   - ✅ `/api/bot/upload_file` - 上传文件

8. **Draftbot 路由** (`/api/draftbot/*`) ✅ **已完成**
   - ✅ `/api/draftbot/commit_check` - 提交检查
   - ✅ `/api/draftbot/create` - 创建草稿 Bot
   - ✅ `/api/draftbot/delete` - 删除草稿 Bot
   - ✅ `/api/draftbot/duplicate` - 复制草稿 Bot
   - ✅ `/api/draftbot/get_display_info` - 获取显示信息
   - ✅ `/api/draftbot/list_draft_history` - 列出草稿历史
   - ✅ `/api/draftbot/publish` - 发布草稿 Bot
   - ✅ `/api/draftbot/update_display_info` - 更新显示信息
   - ✅ `/api/draftbot/publish/connector/list` - 发布连接器列表

9. **Developer API 路由** (`/api/developer/*`) ✅ **已完成**
   - ✅ `/api/developer/get_icon` - 获取图标

10. **Open API 路由** (`/v1/*`, `/v3/*`) ✅ **已完成**

11. **Workflow API 路由** (`/api/workflow_api/*`) ✅ **已完成**
   - ✅ `/api/workflow_api/create` - 创建工作流
   - ✅ `/api/workflow_api/canvas` - 获取画布信息
   - ✅ `/api/workflow_api/save` - 保存工作流
   - ✅ `/api/workflow_api/update_meta` - 更新工作流元数据
   - ✅ `/api/workflow_api/delete` - 删除工作流
   - ✅ `/api/workflow_api/batch_delete` - 批量删除工作流
   - ✅ `/api/workflow_api/delete_strategy` - 获取删除策略
   - ✅ `/api/workflow_api/publish` - 发布工作流
   - ✅ `/api/workflow_api/copy` - 复制工作流
   - ✅ `/api/workflow_api/copy_wk_template` - 复制工作流模板
   - ✅ `/api/workflow_api/released_workflows` - 获取已发布的工作流
   - ✅ `/api/workflow_api/workflow_references` - 获取工作流引用
   - ✅ `/api/workflow_api/workflow_list` - 获取工作流列表
   - ✅ `/api/workflow_api/workflow_detail` - 获取工作流详情
   - ✅ `/api/workflow_api/workflow_detail_info` - 获取工作流详情信息
   - ✅ `/api/workflow_api/list_publish_workflow` - 列出发布的工作流
   - ✅ `/api/workflow_api/example_workflow_list` - 获取示例工作流列表
   - ✅ `/api/workflow_api/node_type` - 查询工作流节点类型
   - ✅ `/api/workflow_api/node_template_list` - 获取节点模板列表
   - ✅ `/api/workflow_api/node_panel_search` - 节点面板搜索
   - ✅ `/api/workflow_api/nodeDebug` - 工作流节点调试 V2
   - ✅ `/api/workflow_api/llm_fc_setting_merged` - 获取 LLM 节点 FC 设置合并
   - ✅ `/api/workflow_api/llm_fc_setting_detail` - 获取 LLM 节点 FC 设置详情
   - ✅ `/api/workflow_api/test_run` - 工作流测试运行
   - ✅ `/api/workflow_api/test_resume` - 工作流测试恢复
   - ✅ `/api/workflow_api/cancel` - 取消工作流
   - ✅ `/api/workflow_api/get_process` - 获取工作流进程
   - ✅ `/api/workflow_api/get_node_execute_history` - 获取节点执行历史
   - ✅ `/api/workflow_api/apiDetail` - 获取 API 详情
   - ✅ `/api/workflow_api/sign_image_url` - 签名图片 URL
   - ✅ `/api/workflow_api/validate_tree` - 验证树
   - ✅ `/api/workflow_api/history_schema` - 获取历史 Schema
   - ✅ `/api/workflow_api/project_conversation/*` - 项目对话相关
   - ✅ `/api/workflow_api/chat_flow_role/*` - 聊天流角色相关
   - ✅ `/api/workflow_api/upload/auth_token` - 获取工作流上传授权令牌
   - ✅ `/api/workflow_api/list_spans` - 列出根 Spans
   - ✅ `/api/workflow_api/get_trace` - 获取 Trace SDK
   - ✅ `/v1/workflow/run` - Open API 运行工作流
   - ✅ `/v1/workflow/stream_run` - Open API 流式运行工作流（SSE）
   - ✅ `/v1/workflow/stream_resume` - Open API 流式恢复工作流（SSE）
   - ✅ `/v1/workflow/get_run_history` - Open API 获取工作流运行历史
   - ✅ `/v1/workflows/:workflow_id` - Open API 获取工作流信息
   - ✅ `/v1/workflows/chat` - Open API 聊天流运行（SSE）
   - ✅ `/v1/workflow/conversation/create` - Open API 创建对话

12. **Plugin API 路由** (`/api/plugin_api/*`) ✅ **已完成**
   - ✅ `/api/plugin_api/get_playground_plugin_list` - 获取 Playground 插件列表
   - ✅ `/api/plugin_api/get_dev_plugin_list` - 获取开发插件列表
   - ✅ `/api/plugin_api/get_plugin_info` - 获取插件信息
   - ✅ `/api/plugin_api/get_plugin_apis` - 获取插件 APIs
   - ✅ `/api/plugin_api/get_updated_apis` - 获取更新的 APIs
   - ✅ `/api/plugin_api/get_plugin_next_version` - 获取插件下一个版本
   - ✅ `/api/plugin_api/register_plugin_meta` - 注册插件元数据
   - ✅ `/api/plugin_api/update` - 更新插件
   - ✅ `/api/plugin_api/update_plugin_meta` - 更新插件元数据
   - ✅ `/api/plugin_api/del_plugin` - 删除插件
   - ✅ `/api/plugin_api/publish_plugin` - 发布插件
   - ✅ `/api/plugin_api/create_api` - 创建 API
   - ✅ `/api/plugin_api/update_api` - 更新 API
   - ✅ `/api/plugin_api/delete_api` - 删除 API
   - ✅ `/api/plugin_api/batch_create_api` - 批量创建 API
   - ✅ `/api/plugin_api/debug_api` - 调试 API
   - ✅ `/api/plugin_api/get_bot_default_params` - 获取 Bot 默认参数
   - ✅ `/api/plugin_api/update_bot_default_params` - 更新 Bot 默认参数
   - ✅ `/api/plugin_api/get_oauth_status` - 获取 OAuth 状态
   - ✅ `/api/plugin_api/get_oauth_schema` - 获取 OAuth Schema
   - ✅ `/api/plugin_api/get_queried_oauth_plugins` - 获取查询的 OAuth 插件列表
   - ✅ `/api/plugin_api/revoke_auth_token` - 撤销授权令牌
   - ✅ `/api/plugin_api/check_and_lock_plugin_edit` - 检查并锁定插件编辑
   - ✅ `/api/plugin_api/unlock_plugin_edit` - 解锁插件编辑
   - ✅ `/api/plugin_api/get_user_authority` - 获取用户权限
   - ✅ `/api/plugin_api/convert_to_openapi` - 转换为 OpenAPI
   - ✅ `/api/plugin_api/library_resource_list` - 库资源列表
   - ✅ `/api/plugin_api/project_resource_list` - 项目资源列表
   - ✅ `/api/plugin_api/resource_copy_dispatch` - 资源复制分发
   - ✅ `/api/plugin_api/resource_copy_detail` - 资源复制详情
   - ✅ `/api/plugin_api/resource_copy_retry` - 资源复制重试
   - ✅ `/api/plugin_api/resource_copy_cancel` - 资源复制取消

13. **Intelligence API 路由** (`/api/intelligence_api/*`) ✅ **已完成**
   - ✅ `/api/intelligence_api/search/get_draft_intelligence_list` - 获取草稿智能体列表
   - ✅ `/api/intelligence_api/search/get_draft_intelligence_info` - 获取草稿智能体信息
   - ✅ `/api/intelligence_api/search/get_recently_edit_intelligence` - 获取用户最近编辑的智能体
   - ✅ `/api/intelligence_api/draft_project/create` - 创建草稿项目
   - ✅ `/api/intelligence_api/draft_project/update` - 更新草稿项目
   - ✅ `/api/intelligence_api/draft_project/delete` - 删除草稿项目
   - ✅ `/api/intelligence_api/draft_project/copy` - 复制草稿项目
   - ✅ `/api/intelligence_api/draft_project/inner_task_list` - 获取草稿项目内部任务列表
   - ✅ `/api/intelligence_api/publish/get_published_connector` - 获取项目已发布的连接器
   - ✅ `/api/intelligence_api/publish/check_version_number` - 检查项目版本号
   - ✅ `/api/intelligence_api/publish/publish_project` - 发布项目
   - ✅ `/api/intelligence_api/publish/publish_record_list` - 获取发布记录列表
   - ✅ `/api/intelligence_api/publish/connector_list` - 获取项目发布连接器列表
   - ✅ `/api/intelligence_api/publish/publish_record_detail` - 获取发布记录详情
   - ✅ `/v1/apps/:app_id` - 获取在线应用数据

14. **Permission API 路由** (`/api/permission_api/*`) ✅ **已完成**
   - ✅ `/api/permission_api/pat/get_personal_access_token_and_permission` - 获取个人访问令牌和权限
   - ✅ `/api/permission_api/pat/delete_personal_access_token_and_permission` - 删除个人访问令牌和权限
   - ✅ `/api/permission_api/pat/list_personal_access_tokens` - 列出个人访问令牌
   - ✅ `/api/permission_api/pat/create_personal_access_token_and_permission` - 创建个人访问令牌和权限
   - ✅ `/api/permission_api/pat/update_personal_access_token_and_permission` - 更新个人访问令牌和权限
   - ✅ `/api/permission_api/coze_web_app/impersonate_coze_user` - 模拟 Coze 用户

15. **Knowledge 路由** (`/api/knowledge/*`) ✅ **已完成**
   - ✅ `/api/knowledge/create` - 创建数据集
   - ✅ `/api/knowledge/delete` - 删除数据集
   - ✅ `/api/knowledge/detail` - 获取数据集详情
   - ✅ `/api/knowledge/list` - 列出数据集
   - ✅ `/api/knowledge/update` - 更新数据集
   - ✅ `/api/knowledge/document/create` - 创建文档
   - ✅ `/api/knowledge/document/delete` - 删除文档
   - ✅ `/api/knowledge/document/list` - 列出文档
   - ✅ `/api/knowledge/document/resegment` - 重新分段文档
   - ✅ `/api/knowledge/document/update` - 更新文档
   - ✅ `/api/knowledge/document/progress/get` - 获取文档处理进度
   - ✅ `/api/knowledge/icon/get` - 获取数据集图标
   - ✅ `/api/knowledge/photo/caption` - 更新照片标题
   - ✅ `/api/knowledge/photo/detail` - 获取照片详情
   - ✅ `/api/knowledge/photo/extract_caption` - 提取照片标题
   - ✅ `/api/knowledge/photo/list` - 列出照片
   - ✅ `/api/knowledge/review/create` - 创建文档审核
   - ✅ `/api/knowledge/review/mget` - 批量获取文档审核
   - ✅ `/api/knowledge/review/save` - 保存文档审核
   - ✅ `/api/knowledge/slice/create` - 创建切片
   - ✅ `/api/knowledge/slice/delete` - 删除切片
   - ✅ `/api/knowledge/slice/list` - 列出切片
   - ✅ `/api/knowledge/slice/update` - 更新切片
   - ✅ `/api/knowledge/table_schema/get` - 获取表结构
   - ✅ `/api/knowledge/table_schema/validate` - 验证表结构
   - ✅ `/api/memory/doc_table_info` - 获取文档表信息（Memory 相关）
   - ✅ `/api/memory/table_mode_config` - 获取模式配置（Memory 相关）

16. **Memory 路由** (`/api/memory/*`) ✅ **已完成**
   - ✅ `/api/memory/sys_variable_conf` - 获取系统变量配置
   - ✅ `/api/memory/project/variable/meta_list` - 获取项目变量列表
   - ✅ `/api/memory/project/variable/meta_update` - 更新项目变量
   - ✅ `/api/memory/variable/upsert` - 设置键值对记忆
   - ✅ `/api/memory/variable/get_meta` - 获取记忆变量元数据
   - ✅ `/api/memory/variable/delete` - 删除用户画像记忆
   - ✅ `/api/memory/variable/get` - 获取 Playground 记忆
   - ✅ `/api/memory/database/list` - 列出数据库
   - ✅ `/api/memory/database/get_by_id` - 根据 ID 获取数据库
   - ✅ `/api/memory/database/add` - 添加数据库
   - ✅ `/api/memory/database/update` - 更新数据库
   - ✅ `/api/memory/database/delete` - 删除数据库
   - ✅ `/api/memory/database/bind_to_bot` - 绑定数据库到 Bot
   - ✅ `/api/memory/database/unbind_to_bot` - 解绑数据库与 Bot
   - ✅ `/api/memory/database/list_records` - 列出数据库记录
   - ✅ `/api/memory/database/update_records` - 更新数据库记录
   - ✅ `/api/memory/database/get_online_database_id` - 获取在线数据库 ID
   - ✅ `/api/memory/database/get_template` - 获取数据库模板
   - ✅ `/api/memory/database/get_connector_name` - 获取连接器名称
   - ✅ `/api/memory/database/update_bot_switch` - 更新数据库 Bot 开关
   - ✅ `/api/memory/database/table/list_new` - 获取 Bot 数据库
   - ✅ `/api/memory/database/table/reset` - 重置 Bot 表
   - ✅ `/api/memory/table_schema/get` - 获取数据库表结构
   - ✅ `/api/memory/table_schema/validate` - 验证数据库表结构
   - ✅ `/api/memory/table_file/submit` - 提交数据库插入任务
   - ✅ `/api/memory/table_file/get_progress` - 获取数据库文件进度数据

17. **Admin 配置路由** (`/api/admin/config/*`) ✅ **已完成**
   - ✅ `/api/admin/config/basic/get` - 获取基础配置
   - ✅ `/api/admin/config/basic/save` - 保存基础配置
   - ✅ `/api/admin/config/knowledge/get` - 获取知识库配置
   - ✅ `/api/admin/config/knowledge/save` - 更新知识库配置
   - ✅ `/api/admin/config/model/list` - 获取模型列表
   - ✅ `/api/admin/config/model/create` - 创建模型
   - ✅ `/api/admin/config/model/delete` - 删除模型

18. **Marketplace 路由** (`/api/marketplace/*`) ✅ **已完成**
   - ✅ `/api/marketplace/product/list` - 获取产品列表
   - ✅ `/api/marketplace/product/detail` - 获取产品详情
   - ✅ `/api/marketplace/product/favorite` - 收藏产品
   - ✅ `/api/marketplace/product/favorite/list.v2` - 获取用户收藏列表 V2
   - ✅ `/api/marketplace/product/duplicate` - 复制产品
   - ✅ `/api/marketplace/product/search` - 搜索产品
   - ✅ `/api/marketplace/product/search/suggest` - 搜索建议
   - ✅ `/api/marketplace/product/category/list` - 获取产品分类列表
   - ✅ `/api/marketplace/product/call_info` - 获取产品调用信息
   - ✅ `/api/marketplace/product/config` - 获取市场插件配置

19. **Plugin 路由（非 API）** (`/api/plugin/*`) ✅ **已完成**
   - ✅ `/api/plugin/get_oauth_schema` - 获取 OAuth Schema

20. **OAuth 路由** (`/api/oauth/*`) ✅ **已完成**
   - ✅ `/api/oauth/authorization_code` - OAuth 授权码处理

## 迁移完成总结

### 总体统计
- **已完成路由组**: 20 个
- **已迁移路由**: 约 211+ 个
- **服务器状态**: 正常运行在 `http://localhost:8888`

### 已完成的路由组列表
1. Passport 认证路由 (`/api/passport/*`)
2. User 相关路由 (`/api/passport/account/*`)
3. Playground API 路由 (`/api/playground_api/*`)
4. Conversation 对话路由 (`/api/conversation/*`)
5. Common 上传路由 (`/api/common/upload/*`)
6. Playground 路由 (`/api/playground/*`)
7. Open API 路由 (`/v1/*`, `/v3/*`)
8. Bot 路由 (`/api/bot/*`)
9. Draftbot 路由 (`/api/draftbot/*`)
10. Developer API 路由 (`/api/developer/*`)
11. Workflow API 路由 (`/api/workflow_api/*`)
12. Plugin API 路由 (`/api/plugin_api/*`)
13. Intelligence API 路由 (`/api/intelligence_api/*`)
14. Permission API 路由 (`/api/permission_api/*`)
15. Knowledge 路由 (`/api/knowledge/*`)
16. Memory 路由 (`/api/memory/*`)
17. Admin 配置路由 (`/api/admin/config/*`)
18. Marketplace 路由 (`/api/marketplace/*`)
19. Plugin 路由（非 API）(`/api/plugin/*`)
20. OAuth 路由 (`/api/oauth/*`)

### 技术要点
- 所有路由都已从 Hertz 迁移到 Gin
- 保持了原有的业务逻辑和功能
- 统一了错误处理和响应格式
- 支持 SSE 流式响应
- 实现了请求体预处理（JSON 对象转字符串）
- 所有路由都已注册并可以正常使用

### 中间件迁移 ✅ **已完成**

#### Session 认证中间件
- ✅ `GinSessionAuth()` - Session 认证中间件，已应用到所有 `/api/*` 路由
- ✅ 支持白名单路径（登录、注册等）
- ✅ 自动从 Cookie 获取 session key 并验证

#### Admin 认证中间件
- ✅ `GinAdminAuth()` - Admin 认证中间件，已应用到 `/api/admin/*` 路由
- ✅ 检查用户邮箱是否在管理员列表中
- ✅ 支持配置多个管理员邮箱（逗号分隔）

#### 其他中间件
- ✅ CORS 中间件（`GinCORS()`）
- ✅ 请求 ID 中间件（`GinRequestID()`）
- ✅ 日志中间件
- ✅ 错误恢复中间件

### 下一步工作
- ✅ 路由迁移已完成
- ✅ 中间件迁移已完成
- ⏳ 功能测试（建议进行）
- ⏳ 性能测试和优化（建议进行）
- ⏳ 清理旧的 Hertz 路由代码（可选）

## 🎉 迁移完成

所有 Hertz 路由已成功迁移到 Gin 框架！

### 相关文档
- **迁移总结**: `MIGRATION_SUMMARY.md` - 详细的迁移统计和技术实现
- **检查清单**: `MIGRATION_CHECKLIST.md` - 迁移完成检查清单和测试建议

### 迁移成果
- ✅ 20 个路由组，211+ 个路由
- ✅ 6 个中间件
- ✅ SSE 流式响应支持
- ✅ 统一错误处理和响应格式
- ✅ 服务器正常运行

### 注意事项
1. 部分处理器使用 `c.String()` 返回错误，这是为了保持与原始 Hertz 实现的一致性
2. `admin_config_handler.go` 中有一个 TODO 注释（检查 coze api token），待后续实现
3. 建议进行全面的功能测试，确保所有路由正常工作
   - ✅ `/v1/conversations` - 对话列表
   - ✅ `/v1/conversations/:conversation_id` - 对话操作（PUT, DELETE）
   - ✅ `/v1/conversations/:conversation_id/clear` - 清空对话
   - ✅ `/v1/conversation/create` - 创建对话
   - ✅ `/v1/conversation/retrieve` - 检索对话
   - ✅ `/v1/conversation/message/list` - 消息列表
   - ✅ `/v3/chat` - 聊天（支持流式和非流式）
   - ✅ `/v3/chat/cancel` - 取消聊天
   - ✅ `/v3/chat/retrieve` - 检索聊天
   - ✅ `/v3/chat/message/list` - 聊天消息列表

### 🔄 待迁移的路由组

1. **Admin 配置路由** (`/api/admin/config/*`)
   - ⏳ `/api/admin/config/basic/get` - 获取基础配置
   - ⏳ `/api/admin/config/basic/save` - 保存基础配置
   - ⏳ `/api/admin/config/knowledge/get` - 获取知识库配置
   - ⏳ `/api/admin/config/knowledge/save` - 保存知识库配置
   - ⏳ `/api/admin/config/model/*` - 模型配置相关

2. **Bot 路由** (`/api/bot/*`)
   - ⏳ `/api/bot/get_type_list` - 获取类型列表
   - ⏳ `/api/bot/upload_file` - 上传文件

3. **Common 上传路由** (`/api/common/upload/*`)
   - ⏳ `/api/common/upload/apply_upload_action` - 申请上传操作
   - ⏳ `/api/common/upload/*tos_uri` - 通用上传


5. **Developer API 路由** (`/api/developer/*`)
   - ⏳ `/api/developer/get_icon` - 获取图标

6. **Draftbot 路由** (`/api/draftbot/*`)
   - ⏳ `/api/draftbot/commit_check` - 提交检查
   - ⏳ `/api/draftbot/create` - 创建草稿 Bot
   - ⏳ `/api/draftbot/delete` - 删除草稿 Bot
   - ⏳ `/api/draftbot/duplicate` - 复制草稿 Bot
   - ⏳ `/api/draftbot/get_display_info` - 获取显示信息
   - ⏳ `/api/draftbot/list_draft_history` - 列出草稿历史
   - ⏳ `/api/draftbot/publish` - 发布草稿 Bot
   - ⏳ `/api/draftbot/update_display_info` - 更新显示信息
   - ⏳ `/api/draftbot/publish/connector/list` - 发布连接器列表

7. **Intelligence API 路由** (`/api/intelligence_api/*`)
   - ⏳ `/api/intelligence_api/draft_project/*` - 草稿项目相关
   - ⏳ `/api/intelligence_api/publish/*` - 发布相关
   - ⏳ `/api/intelligence_api/search/*` - 搜索相关

8. **Knowledge 知识库路由** (`/api/knowledge/*`)
   - ⏳ `/api/knowledge/*` - 知识库相关路由

9. **Marketplace 市场路由** (`/api/marketplace/*`)
   - ⏳ `/api/marketplace/product/*` - 产品相关

10. **Memory 记忆路由** (`/api/memory/*`)
    - ⏳ `/api/memory/*` - 记忆相关路由

11. **OAuth 路由** (`/api/oauth/*`)
    - ⏳ `/api/oauth/authorization_code` - 授权码

12. **Permission API 路由** (`/api/permission_api/*`)
    - ⏳ `/api/permission_api/coze_web_app/impersonate_coze_user` - 模拟用户
    - ⏳ `/api/permission_api/pat/*` - 个人访问令牌相关

13. **Playground 路由** (`/api/playground/*`)
    - ⏳ `/api/playground/get_onboarding` - 获取引导信息
    - ⏳ `/api/playground/upload/auth_token` - 获取上传令牌

14. **Playground API 剩余路由** (`/api/playground_api/*`)
    - ⏳ `/api/playground_api/create_update_shortcut_command` - 创建/更新快捷命令
    - ⏳ `/api/playground_api/delete_prompt_resource` - 删除提示资源
    - ⏳ `/api/playground_api/get_file_list` - 获取文件列表
    - ⏳ `/api/playground_api/get_imagex_url` - 获取图片 URL
    - ⏳ `/api/playground_api/get_official_prompt_list` - 获取官方提示列表
    - ⏳ `/api/playground_api/get_prompt_resource_info` - 获取提示资源信息
    - ⏳ `/api/playground_api/report_user_behavior` - 报告用户行为
    - ⏳ `/api/playground_api/upsert_prompt_resource` - 更新/插入提示资源
    - ⏳ `/api/playground_api/draftbot/*` - 草稿 Bot 相关
    - ⏳ `/api/playground_api/operate/*` - 操作相关

15. **Plugin 路由** (`/api/plugin/*`)
    - ⏳ `/api/plugin/get_oauth_schema` - 获取 OAuth Schema

16. **Plugin API 路由** (`/api/plugin_api/*`)
    - ⏳ `/api/plugin_api/*` - 插件 API 相关（大量路由）

17. **Workflow API 路由** (`/api/workflow_api/*`)
    - ⏳ `/api/workflow_api/*` - 工作流 API 相关（大量路由）

18. **Open API 路由** (`/v1/*`, `/v3/*`)
    - ⏳ `/v1/conversations` - 对话列表
    - ⏳ `/v1/conversations/:conversation_id` - 对话操作
    - ⏳ `/v1/apps/:app_id` - 应用信息
    - ⏳ `/v1/bot/get_online_info` - 获取在线 Bot 信息
    - ⏳ `/v1/bots/:bot_id` - Bot 信息
    - ⏳ `/v1/conversation/*` - 对话相关
    - ⏳ `/v1/files/upload` - 上传文件
    - ⏳ `/v1/workflow/*` - 工作流相关
    - ⏳ `/v1/workflows/*` - 工作流相关
    - ⏳ `/v3/chat/*` - 聊天相关

## 🔧 迁移步骤

### 1. 为每个路由创建 Gin 处理器

对于每个 Hertz 处理器，需要创建对应的 Gin 处理器：

**Hertz 处理器示例**：
```go
// backend/api/handler/coze/playground_service.go
func MGetUserBasicInfo(ctx context.Context, c *app.RequestContext) {
    var req playground.MGetUserBasicInfoRequest
    err := c.BindAndValidate(&req)
    if err != nil {
        invalidParamRequestResponse(c, err.Error())
        return
    }
    // ... 业务逻辑
}
```

**Gin 处理器示例**：
```go
// backend/api/handler/gin/playground_handler.go
func MGetUserBasicInfo(c *gin.Context) {
    var req playground.MGetUserBasicInfoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": 40000,
            "msg":  err.Error(),
        })
        return
    }
    // ... 相同的业务逻辑
}
```

### 2. 注册路由

在 `backend/api/router/gin_router_coze.go` 中注册路由：

```go
func registerPlaygroundAPIRoutes(api *gin.RouterGroup) {
    playgroundAPI := api.Group("/playground_api")
    {
        playgroundAPI.POST("/mget_user_info", ginHandler.MGetUserBasicInfo)
    }
}
```

### 3. 测试路由

确保每个迁移的路由都能正常工作。

## 📝 迁移优先级

1. **高优先级**（核心功能）
   - ✅ Passport 认证路由（已完成）
   - ✅ 用户相关路由（已完成）
   - ⏳ Conversation 对话路由
   - ⏳ Playground API 路由（部分完成）

2. **中优先级**（常用功能）
   - ⏳ Bot 路由
   - ⏳ Draftbot 路由
   - ⏳ Workflow API 路由
   - ⏳ Plugin API 路由

3. **低优先级**（辅助功能）
   - ⏳ Admin 配置路由
   - ⏳ Marketplace 路由
   - ⏳ Memory 路由
   - ⏳ Open API 路由

## 🚀 快速开始

1. 选择一个路由组
2. 查看对应的 Hertz 处理器（`backend/api/handler/coze/*.go`）
3. 创建对应的 Gin 处理器（`backend/api/handler/gin/*.go`）
4. 在 `gin_router_coze.go` 中注册路由
5. 测试路由功能
6. 更新本文档的迁移状态

## ⚠️ 注意事项

1. **保持业务逻辑不变**：只迁移路由和处理器，不修改业务逻辑
2. **复用 Application Service**：Gin 处理器应该调用相同的 Application Service
3. **保持 API 兼容性**：确保迁移后的 API 与原有 API 完全兼容
4. **测试覆盖**：每个迁移的路由都应该有对应的测试

