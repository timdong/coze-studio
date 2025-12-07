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
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
	"github.com/coze-dev/coze-studio/backend/application/singleagent"
)

// DraftBotCreate 创建草稿 Bot
// @router /api/draftbot/create [POST]
func DraftBotCreate(c *gin.Context) {
	var req developer_api.DraftBotCreateRequest
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
			"msg":  "space id is not set",
		})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "name is nil",
		})
		return
	}

	if req.IconURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "icon uri is nil",
		})
		return
	}

	if utf8.RuneCountInString(req.Name) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "name is too long",
		})
		return
	}

	if utf8.RuneCountInString(req.Description) > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "description is too long",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.CreateSingleAgentDraft(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteDraftBot 删除草稿 Bot
// @router /api/draftbot/delete [POST]
func DeleteDraftBot(c *gin.Context) {
	var req developer_api.DeleteDraftBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.DeleteAgentDraft(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateDraftBotDisplayInfo 更新草稿 Bot 显示信息
// @router /api/draftbot/update_display_info [POST]
func UpdateDraftBotDisplayInfo(c *gin.Context) {
	var req developer_api.UpdateDraftBotDisplayInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.UpdateAgentDraftDisplayInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DuplicateDraftBot 复制草稿 Bot
// @router /api/draftbot/duplicate [POST]
func DuplicateDraftBot(c *gin.Context) {
	var req developer_api.DuplicateDraftBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.DuplicateDraftBot(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDraftBotDisplayInfo 获取草稿 Bot 显示信息
// @router /api/draftbot/get_display_info [POST]
func GetDraftBotDisplayInfo(c *gin.Context) {
	var req developer_api.GetDraftBotDisplayInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.GetAgentDraftDisplayInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PublishDraftBot 发布草稿 Bot
// @router /api/draftbot/publish [POST]
func PublishDraftBot(c *gin.Context) {
	var req developer_api.PublishDraftBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if len(req.Connectors) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "connectors is nil",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.PublishAgent(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListDraftBotHistory 列出草稿 Bot 历史
// @router /api/draftbot/list_draft_history [POST]
func ListDraftBotHistory(c *gin.Context) {
	var req developer_api.ListDraftBotHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if req.BotID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "bot id is not set",
		})
		return
	}

	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 30
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.ListAgentPublishHistory(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PublishConnectorList 获取发布连接器列表
// @router /api/draftbot/publish/connector/list [POST]
func PublishConnectorList(c *gin.Context) {
	var req developer_api.PublishConnectorListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if req.BotID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "bot id is not set",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.GetPublishConnectorList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CheckDraftBotCommit 检查草稿 Bot 提交
// @router /api/draftbot/commit_check [POST]
func CheckDraftBotCommit(c *gin.Context) {
	var req developer_api.CheckDraftBotCommitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := new(developer_api.CheckDraftBotCommitResponse)
	c.JSON(http.StatusOK, resp)
}

