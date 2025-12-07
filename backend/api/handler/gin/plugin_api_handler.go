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

package gin

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/plugin_develop"
	common "github.com/coze-dev/coze-studio/backend/api/model/plugin_develop/common"
	"github.com/coze-dev/coze-studio/backend/api/model/resource"
	appApplication "github.com/coze-dev/coze-studio/backend/application/app"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	"github.com/coze-dev/coze-studio/backend/application/plugin"
	"github.com/coze-dev/coze-studio/backend/application/search"
)

// GetPlaygroundPluginList 获取 Playground 插件列表
// @router /api/plugin_api/get_playground_plugin_list [POST]
func GetPlaygroundPluginList(c *gin.Context) {
	var req plugin_develop.GetPlaygroundPluginListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.GetSpaceID() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "spaceID is invalid",
		})
		return
	}
	if req.GetPage() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "page is invalid",
		})
		return
	}
	if req.GetSize() >= 30 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "size is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	// when there is only one element in the types list, and the element type is workflow, use workflow service
	if len(req.GetPluginTypes()) == 1 && req.GetPluginTypes()[0] == int32(common.PluginType_WORKFLOW) {
		resp, err := appworkflow.SVC.GetPlaygroundPluginList(ctx, &req)
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp, err := plugin.PluginApplicationSVC.GetPlaygroundPluginList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RegisterPluginMeta 注册插件元数据
// @router /api/plugin_api/register_plugin_meta [POST]
func RegisterPluginMeta(c *gin.Context) {
	var req plugin_develop.RegisterPluginMetaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.GetName() == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin name is invalid",
		})
		return
	}
	if req.GetDesc() == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin desc is invalid",
		})
		return
	}
	if req.URL != nil && (*req.URL == "" || len(*req.URL) > 512) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin url is invalid",
		})
		return
	}
	if req.Icon == nil || req.Icon.URI == "" || len(req.Icon.URI) > 512 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin icon is invalid",
		})
		return
	}
	if req.AuthType == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin auth type is invalid",
		})
		return
	}
	if req.SpaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "spaceID is invalid",
		})
		return
	}
	if req.ProjectID != nil {
		if *req.ProjectID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  "projectID is invalid",
			})
			return
		}
	}
	if req.GetPluginType() != common.PluginType_PLUGIN {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin type is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.RegisterPluginMeta(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPluginAPIs 获取插件 APIs
// @router /api/plugin_api/get_plugin_apis [POST]
func GetPluginAPIs(c *gin.Context) {
	var req plugin_develop.GetPluginAPIsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}
	if len(req.APIIds) == 0 {
		if req.Page <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  "page is invalid",
			})
			return
		}
		if req.Size >= 30 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  "size is invalid",
			})
			return
		}
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetPluginAPIs(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPluginInfo 获取插件信息
// @router /api/plugin_api/get_plugin_info [POST]
func GetPluginInfo(c *gin.Context) {
	var req plugin_develop.GetPluginInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetPluginInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUpdatedAPIs 获取更新的 APIs
// @router /api/plugin_api/get_updated_apis [POST]
func GetUpdatedAPIs(c *gin.Context) {
	var req plugin_develop.GetUpdatedAPIsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetUpdatedAPIs(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetOAuthStatus 获取 OAuth 状态
// @router /api/plugin_api/get_oauth_status [POST]
func GetOAuthStatus(c *gin.Context) {
	var req plugin_develop.GetOAuthStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetOAuthStatus(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CheckAndLockPluginEdit 检查并锁定插件编辑
// @router /api/plugin_api/check_and_lock_plugin_edit [POST]
func CheckAndLockPluginEdit(c *gin.Context) {
	var req plugin_develop.CheckAndLockPluginEditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.CheckAndLockPluginEdit(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdatePlugin 更新插件
// @router /api/plugin_api/update [POST]
func UpdatePlugin(c *gin.Context) {
	var req plugin_develop.UpdatePluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}
	if req.AiPlugin == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin manifest is invalid",
		})
		return
	}
	if req.Openapi == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin openapi doc is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.UpdatePlugin(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteAPI 删除 API
// @router /api/plugin_api/delete_api [POST]
func DeleteAPI(c *gin.Context) {
	var req plugin_develop.DeleteAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}
	if req.APIID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "apiID is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.DeleteAPI(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DelPlugin 删除插件
// @router /api/plugin_api/del_plugin [POST]
func DelPlugin(c *gin.Context) {
	var req plugin_develop.DelPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.DelPlugin(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PublishPlugin 发布插件
// @router /api/plugin_api/publish_plugin [POST]
func PublishPlugin(c *gin.Context) {
	var req plugin_develop.PublishPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}
	if req.VersionName == "" || len(req.VersionName) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "version name is invalid",
		})
		return
	}

	match, _ := regexp.MatchString(`^v\d+\.\d+\.\d+$`, req.VersionName)
	if !match {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "version name is invalid",
		})
		return
	}

	if req.VersionDesc == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "version desc is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.PublishPlugin(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdatePluginMeta 更新插件元数据
// @router /api/plugin_api/update_plugin_meta [POST]
func UpdatePluginMeta(c *gin.Context) {
	var req plugin_develop.UpdatePluginMetaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}
	if req.Name != nil && *req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin name is invalid",
		})
		return
	}
	if req.Desc != nil && *req.Desc == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin desc is invalid",
		})
		return
	}
	if req.URL != nil && (*req.URL == "" || len(*req.URL) > 512) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin server url is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.UpdatePluginMeta(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetBotDefaultParams 获取 Bot 默认参数
// @router /api/plugin_api/get_bot_default_params [POST]
func GetBotDefaultParams(c *gin.Context) {
	var req plugin_develop.GetBotDefaultParamsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.SpaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "spaceID is invalid",
		})
		return
	}
	if req.BotID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "botID is invalid",
		})
		return
	}
	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}
	if req.APIName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "apiName is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetBotDefaultParams(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateBotDefaultParams 更新 Bot 默认参数
// @router /api/plugin_api/update_bot_default_params [POST]
func UpdateBotDefaultParams(c *gin.Context) {
	var req plugin_develop.UpdateBotDefaultParamsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.SpaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "spaceID is invalid",
		})
		return
	}
	if req.BotID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "botID is invalid",
		})
		return
	}
	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}
	if req.APIName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "apiName is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.UpdateBotDefaultParams(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateAPI 创建 API
// @router /api/plugin_api/create_api [POST]
func CreateAPI(c *gin.Context) {
	var req plugin_develop.CreateAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}
	if req.Name == "" || len(req.Name) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "api name is invalid",
		})
		return
	}
	if req.Desc == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "api desc is invalid",
		})
		return
	}
	if req.Path != nil && (*req.Path == "" || len(*req.Path) > 512) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "api path is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.CreateAPI(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateAPI 更新 API
// @router /api/plugin_api/update_api [POST]
func UpdateAPI(c *gin.Context) {
	var req plugin_develop.UpdateAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}
	if req.APIID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "apiID is invalid",
		})
		return
	}
	if req.Name != nil && (*req.Name == "" || len(*req.Name) > 255) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "api name is invalid",
		})
		return
	}
	if req.Desc != nil && (*req.Desc == "" || len(*req.Desc) > 255) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "api desc is invalid",
		})
		return
	}
	if req.Path != nil && (*req.Path == "" || len(*req.Path) > 512) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "api path is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.UpdateAPI(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserAuthority 获取用户权限
// @router /api/plugin_api/get_user_authority [POST]
func GetUserAuthority(c *gin.Context) {
	var req plugin_develop.GetUserAuthorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetUserAuthority(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DebugAPI 调试 API
// @router /api/plugin_api/debug_api [POST]
func DebugAPI(c *gin.Context) {
	var req plugin_develop.DebugAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}
	if req.APIID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "apiID is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.DebugAPI(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UnlockPluginEdit 解锁插件编辑
// @router /api/plugin_api/unlock_plugin_edit [POST]
func UnlockPluginEdit(c *gin.Context) {
	var req plugin_develop.UnlockPluginEditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.UnlockPluginEdit(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPluginNextVersion 获取插件下一个版本
// @router /api/plugin_api/get_plugin_next_version [POST]
func GetPluginNextVersion(c *gin.Context) {
	var req plugin_develop.GetPluginNextVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetPluginNextVersion(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDevPluginList 获取开发插件列表
// @router /api/plugin_api/get_dev_plugin_list [POST]
func GetDevPluginList(c *gin.Context) {
	var req plugin_develop.GetDevPluginListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.SpaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "spaceID is invalid",
		})
		return
	}
	if req.ProjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "projectID is invalid",
		})
		return
	}
	if req.GetPage() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "page is invalid",
		})
		return
	}
	if req.GetSize() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "size is invalid",
		})
		return
	}
	if req.GetSize() > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "size is too large",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetDevPluginList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Convert2OpenAPI 转换为 OpenAPI
// @router /api/plugin_api/convert_to_openapi [POST]
func Convert2OpenAPI(c *gin.Context) {
	var req plugin_develop.Convert2OpenAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.SpaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "spaceID is invalid",
		})
		return
	}
	if req.Data == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "data is invalid",
		})
		return
	}
	if req.PluginURL != nil && *req.PluginURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginURL is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.Convert2OpenAPI(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetOAuthSchemaAPI 获取 OAuth Schema API
// @router /api/plugin_api/get_oauth_schema [POST]
func GetOAuthSchemaAPI(c *gin.Context) {
	var req plugin_develop.GetOAuthSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetOAuthSchema(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetOAuthSchema 获取 OAuth Schema（非 API 路由）
// @router /api/plugin/get_oauth_schema [POST]
func GetOAuthSchema(c *gin.Context) {
	var req plugin_develop.GetOAuthSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetOAuthSchema(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// BatchCreateAPI 批量创建 API
// @router /api/plugin_api/batch_create_api [POST]
func BatchCreateAPI(c *gin.Context) {
	var req plugin_develop.BatchCreateAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.SpaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "spaceID is invalid",
		})
		return
	}
	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}
	if req.AiPlugin == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin manifest is invalid",
		})
		return
	}
	if req.Openapi == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "plugin openapi doc is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.BatchCreateAPI(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RevokeAuthToken 撤销授权令牌
// @router /api/plugin_api/revoke_auth_token [POST]
func RevokeAuthToken(c *gin.Context) {
	var req plugin_develop.RevokeAuthTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.PluginID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pluginID is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.RevokeAuthToken(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetQueriedOAuthPluginList 获取查询的 OAuth 插件列表
// @router /api/plugin_api/get_queried_oauth_plugins [POST]
func GetQueriedOAuthPluginList(c *gin.Context) {
	var req plugin_develop.GetQueriedOAuthPluginListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.BotID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "entityID is required",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetQueriedOAuthPluginList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// LibraryResourceList 库资源列表
// @router /api/plugin_api/library_resource_list [POST]
func LibraryResourceList(c *gin.Context) {
	var req resource.LibraryResourceListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.SpaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "space_id is invalid",
		})
		return
	}
	if req.GetSize() > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "size is too large",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := search.SearchSVC.LibraryResourceList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ProjectResourceList 项目资源列表
// @router /api/plugin_api/project_resource_list [POST]
func ProjectResourceList(c *gin.Context) {
	var req resource.ProjectResourceListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.SpaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "space_id is invalid",
		})
		return
	}
	if req.ProjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "project_id is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := search.SearchSVC.ProjectResourceList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ResourceCopyDispatch 资源复制分发
// @router /api/plugin_api/resource_copy_dispatch [POST]
func ResourceCopyDispatch(c *gin.Context) {
	var req resource.ResourceCopyDispatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ResID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "res_id is invalid",
		})
		return
	}
	if req.ResType <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "res_type is invalid",
		})
		return
	}
	if req.GetProjectID() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "project_id is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.ResourceCopyDispatch(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ResourceCopyDetail 资源复制详情
// @router /api/plugin_api/resource_copy_detail [POST]
func ResourceCopyDetail(c *gin.Context) {
	var req resource.ResourceCopyDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.TaskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "task_id is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.ResourceCopyDetail(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ResourceCopyRetry 资源复制重试
// @router /api/plugin_api/resource_copy_retry [POST]
func ResourceCopyRetry(c *gin.Context) {
	var req resource.ResourceCopyRetryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	resp := new(resource.ResourceCopyRetryResponse)
	c.JSON(http.StatusOK, resp)
}

// ResourceCopyCancel 资源复制取消
// @router /api/plugin_api/resource_copy_cancel [POST]
func ResourceCopyCancel(c *gin.Context) {
	var req resource.ResourceCopyCancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := new(resource.ResourceCopyCancelResponse)
	c.JSON(http.StatusOK, resp)
}

