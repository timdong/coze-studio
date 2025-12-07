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
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/playground"
	appApplication "github.com/coze-dev/coze-studio/backend/application/app"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	"github.com/coze-dev/coze-studio/backend/application/prompt"
	"github.com/coze-dev/coze-studio/backend/application/shortcutcmd"
	"github.com/coze-dev/coze-studio/backend/application/singleagent"
	"github.com/coze-dev/coze-studio/backend/application/upload"
	"github.com/coze-dev/coze-studio/backend/application/user"
	userService "github.com/coze-dev/coze-studio/backend/domain/user/service"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// MGetUserBasicInfo 批量获取用户基本信息
// @router /api/playground_api/mget_user_info [POST]
func MGetUserBasicInfo(c *gin.Context) {
	var req playground.MGetUserBasicInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := user.UserApplicationSVC.MGetUserBasicInfo(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "MGetUserBasicInfo failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetSpaceListV2 获取空间列表
// @router /api/playground_api/space/list [POST]
func GetSpaceListV2(c *gin.Context) {
	var req playground.GetSpaceListV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := user.UserApplicationSVC.GetSpaceListV2(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "GetSpaceListV2 failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// OpenCreateSpace 创建空间（OpenAPI）
// @router /v1/workspaces [POST]
func OpenCreateSpace(c *gin.Context) {
	var req user.OpenCreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := user.UserApplicationSVC.OpenCreateSpace(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "OpenCreateSpace failed: %v", err)
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SaveSpaceV2 保存/创建空间（Playground API）
// @router /api/playground_api/space/save [POST]
func SaveSpaceV2(c *gin.Context) {
	var req struct {
		SpaceID        *string `json:"space_id,omitempty"`
		Name           string  `json:"name" binding:"required"`
		Description    string  `json:"description"`
		IconURI        string  `json:"icon_uri"`
		SpaceType      int32   `json:"space_type" binding:"required"`
		SpaceMode      *int32  `json:"space_mode,omitempty"`
		EnterpriseID   *int64  `json:"enterprise_id,omitempty"`
		OrganizationID *int64  `json:"organization_id,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 40100,
			"msg":  "authentication required",
		})
		return
	}

	// 如果有 space_id，说明是更新；否则是创建
	if req.SpaceID != nil && *req.SpaceID != "" && *req.SpaceID != "0" {
		// TODO: 实现更新空间的逻辑
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": 50100,
			"msg":  "update space not implemented yet",
		})
		return
	}

	// 创建新空间
	space, err := user.UserApplicationSVC.DomainSVC.CreateSpace(ctx, &userService.CreateSpaceRequest{
		Name:        req.Name,
		Description: req.Description,
		IconURI:     req.IconURI,
		OwnerID:     *uid,
		CreatorID:   *uid,
	})
	if err != nil {
		logs.CtxErrorf(ctx, "SaveSpaceV2 failed: %v", err)
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":          strconv.FormatInt(space.ID, 10), // 转换为字符串以匹配前端期望
			"name":        space.Name,
			"description": space.Description,
			"icon_url":    space.IconURL,
		},
		"code": 0,
		"msg":  "success",
	})
}

// GetDraftBotInfoAgw 获取草稿 Bot 信息
// @router /api/playground_api/draftbot/get_draft_bot_info [POST]
func GetDraftBotInfoAgw(c *gin.Context) {
	var req playground.GetDraftBotInfoAgwRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.BotID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "bot id is nil",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.GetAgentBotInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateDraftBotInfoAgw 更新草稿 Bot 信息
// @router /api/playground_api/draftbot/update_draft_bot_info [POST]
func UpdateDraftBotInfoAgw(c *gin.Context) {
	var req playground.UpdateDraftBotInfoAgwRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.BotInfo == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "bot info is nil",
		})
		return
	}

	if req.BotInfo.BotId == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "bot id is nil",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.UpdateSingleAgentDraft(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetOfficialPromptResourceList 获取官方提示资源列表
// @router /api/playground_api/get_official_prompt_list [POST]
func GetOfficialPromptResourceList(c *gin.Context) {
	var req playground.GetOfficialPromptResourceListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := prompt.PromptSVC.GetOfficialPromptResourceList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPromptResourceInfo 获取提示资源信息
// @router /api/playground_api/get_prompt_resource_info [GET]
func GetPromptResourceInfo(c *gin.Context) {
	var req playground.GetPromptResourceInfoRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := prompt.PromptSVC.GetPromptResourceInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpsertPromptResource 更新或插入提示资源
// @router /api/playground_api/upsert_prompt_resource [POST]
func UpsertPromptResource(c *gin.Context) {
	var req playground.UpsertPromptResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.Prompt == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "prompt is nil",
		})
		return
	}

	if req.Prompt.GetSpaceID() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "space id is invalid",
		})
		return
	}

	if len(req.Prompt.GetName()) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "name is empty",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := prompt.PromptSVC.UpsertPromptResource(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeletePromptResource 删除提示资源
// @router /api/playground_api/delete_prompt_resource [POST]
func DeletePromptResource(c *gin.Context) {
	var req playground.DeletePromptResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := prompt.PromptSVC.DeletePromptResource(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetImagexShortUrl 获取图片短链接
// @router /api/playground_api/get_imagex_url [POST]
func GetImagexShortUrl(c *gin.Context) {
	var req playground.GetImagexShortUrlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if len(req.Uris) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "uris is empty",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.GetImagexShortUrl(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetBotPopupInfo 获取 Bot 弹窗信息
// @router /api/playground_api/operate/get_bot_popup_info [POST]
func GetBotPopupInfo(c *gin.Context) {
	var req playground.GetBotPopupInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if len(req.BotPopupTypes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "bot popup types is empty",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.GetAgentPopupInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateBotPopupInfo 更新 Bot 弹窗信息
// @router /api/playground_api/operate/update_bot_popup_info [POST]
func UpdateBotPopupInfo(c *gin.Context) {
	var req playground.UpdateBotPopupInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.UpdateAgentPopupInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateUpdateShortcutCommand 创建或更新快捷命令
// @router /api/playground_api/create_update_shortcut_command [POST]
func CreateUpdateShortcutCommand(c *gin.Context) {
	var req playground.CreateUpdateShortcutCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	shortCuts, err := shortcutcmd.ShortcutCmdSVC.Handler(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	resp := new(playground.CreateUpdateShortcutCommandResponse)
	resp.Shortcuts = shortCuts
	resp.Code = 0
	resp.Msg = ""
	c.JSON(http.StatusOK, resp)
}

// ReportUserBehavior 报告用户行为
// @router /api/playground_api/report_user_behavior [POST]
func ReportUserBehavior(c *gin.Context) {
	var req playground.ReportUserBehaviorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ResourceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "resource id is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp := new(playground.ReportUserBehaviorResponse)

	if req.ResourceType == playground.SpaceResourceType_DraftBot {
		var err error
		resp, err = singleagent.SingleAgentSVC.ReportUserBehavior(ctx, &req)
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
	} else if req.ResourceType == playground.SpaceResourceType_Project {
		var err error
		resp, err = appApplication.APPApplicationSVC.ReportUserBehavior(ctx, &req)
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
	}

	c.JSON(http.StatusOK, resp)
}

// GetFileUrls 获取文件 URL 列表
// @router /api/playground_api/get_file_list [POST]
func GetFileUrls(c *gin.Context) {
	var req playground.GetFileUrlsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	iconList, err := upload.SVC.GetShortcutIcons(ctx)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	resp := new(playground.GetFileUrlsResponse)
	resp.FileList = iconList
	resp.Code = 0

	c.JSON(http.StatusOK, resp)
}
