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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hertz-contrib/sse"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/conversation/conversation"
	"github.com/coze-dev/coze-studio/backend/api/model/conversation/message"
	"github.com/coze-dev/coze-studio/backend/api/model/conversation/run"
	application "github.com/coze-dev/coze-studio/backend/application/conversation"
	ginSSE "github.com/coze-dev/coze-studio/backend/infra/sse/impl/gin"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// ClearConversationHistory 清空对话历史
// @router /api/conversation/clear_message [POST]
func ClearConversationHistory(c *gin.Context) {
	var req conversation.ClearConversationHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if checkErr := checkCCHParams(ctx, &req); checkErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  checkErr.Error(),
		})
		return
	}

	resp, err := application.ConversationSVC.ClearHistory(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func checkCCHParams(_ interface{}, req *conversation.ClearConversationHistoryRequest) error {
	if req.ConversationID <= 0 {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "invalid conversation id"))
	}

	if req.Scene == nil {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "scene is required"))
	}

	return nil
}

// ClearConversationCtx 创建会话段（清空上下文）
// @router /api/conversation/create_section [POST]
func ClearConversationCtx(c *gin.Context) {
	var req conversation.ClearConversationCtxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if checkErr := checkCCCParams(ctx, &req); checkErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  checkErr.Error(),
		})
		return
	}

	newSectionID, err := application.ConversationSVC.CreateSection(ctx, req.ConversationID)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	resp := new(conversation.ClearConversationCtxResponse)
	resp.NewSectionID = newSectionID

	c.JSON(http.StatusOK, resp)
}

func checkCCCParams(_ interface{}, req *conversation.ClearConversationCtxRequest) error {
	if req.ConversationID <= 0 {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "invalid conversation id"))
	}

	if req.Scene == nil {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "scene is required"))
	}
	return nil
}

// CreateConversation 创建对话
// @router /api/conversation/create [POST]
func CreateConversation(c *gin.Context) {
	var req conversation.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := application.ConversationSVC.CreateConversation(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMessageList 获取消息列表
// @router /api/conversation/get_message_list [POST]
func GetMessageList(c *gin.Context) {
	var req message.GetMessageListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if checkErr := checkMLParams(ctx, &req); checkErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  checkErr.Error(),
		})
		return
	}

	resp, err := application.ConversationSVC.GetMessageList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func checkMLParams(_ interface{}, req *message.GetMessageListRequest) error {
	if req.BotID == "" {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "agent id is required"))
	}

	return nil
}

// DeleteMessage 删除消息
// @router /api/conversation/delete_message [POST]
func DeleteMessage(c *gin.Context) {
	var req message.DeleteMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if checkErr := checkDMParams(ctx, &req); checkErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  checkErr.Error(),
		})
		return
	}

	resp, err := application.ConversationSVC.DeleteMessage(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func checkDMParams(_ interface{}, req *message.DeleteMessageRequest) error {
	if req.MessageID <= 0 {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "message id is invalid"))
	}

	return nil
}

// BreakMessage 中断消息
// @router /api/conversation/break_message [POST]
func BreakMessage(c *gin.Context) {
	var req message.BreakMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if checkErr := checkBMParams(ctx, &req); checkErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  checkErr.Error(),
		})
		return
	}

	resp, err := application.ConversationSVC.BreakMessage(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func checkBMParams(_ interface{}, req *message.BreakMessageRequest) error {
	if req.AnswerMessageID == nil {
		return errors.New("answer message id is required")
	}
	if *req.AnswerMessageID <= 0 {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "answer message id is invalid"))
	}

	return nil
}

// AgentRun 对话聊天（SSE 流式）
// @router /api/conversation/chat [POST]
func AgentRun(c *gin.Context) {
	var req run.AgentRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if checkErr := checkAgentRunParams(ctx, &req); checkErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  checkErr.Error(),
		})
		return
	}

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	// 创建 Gin 版本的 SSE 发送器（兼容 sseImpl.SSenderImpl 接口）
	sseSender := ginSSE.NewGinSSESender(c)

	err := application.ConversationSVC.Run(ctx, sseSender, &req)
	if err != nil {
		errData := run.ErrorData{
			Code: errno.ErrConversationAgentRunError,
			Msg:  err.Error(),
		}
		ed, _ := json.Marshal(errData)
		_ = sseSender.Send(ctx, &sse.Event{
			Event: run.RunEventError,
			Data:  ed,
		})
	}
}

func checkAgentRunParams(_ interface{}, ar *run.AgentRunRequest) error {
	if ar.BotID == 0 {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "bot id is required"))
	}

	if ar.Scene == nil {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "scene is required"))
	}

	if ar.ContentType == nil {
		ar.ContentType = ptr.Of(run.ContentTypeText)
	}
	return nil
}

// ListConversationsApi 获取对话列表（Open API）
// @router /v1/conversations [GET]
func ListConversationsApi(c *gin.Context) {
	var req conversation.ListConversationsApiRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := application.ConversationSVC.ListConversation(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RetrieveConversationApi 检索对话（Open API）
// @router /v1/conversation/retrieve [GET]
func RetrieveConversationApi(c *gin.Context) {
	var req conversation.RetrieveConversationApiRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := application.ConversationSVC.RetrieveConversation(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetApiMessageList 获取消息列表（Open API）
// @router /v1/conversation/message/list [POST]
func GetApiMessageList(c *gin.Context) {
	var req message.ListMessageApiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := application.OpenapiMessageSVC.GetApiMessageList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ClearConversationApi 清空对话（Open API）
// @router /v1/conversations/:conversation_id/clear [POST]
func ClearConversationApi(c *gin.Context) {
	var req conversation.ClearConversationApiRequest
	conversationIDStr := c.Param("conversation_id")
	if conversationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "conversation_id is required",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	// 解析 conversation_id
	var conversationID int64
	if _, err := fmt.Sscanf(conversationIDStr, "%d", &conversationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid conversation_id",
		})
		return
	}
	req.ConversationID = conversationID

	ctx := c.Request.Context()
	sectionID, err := application.ConversationSVC.CreateSection(ctx, conversationID)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	resp := new(conversation.ClearConversationApiResponse)
	resp.Data = &conversation.Section{
		ID:             sectionID,
		ConversationID: req.ConversationID,
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateConversationApi 更新对话（Open API）
// @router /v1/conversations/:conversation_id [PUT]
func UpdateConversationApi(c *gin.Context) {
	var req conversation.UpdateConversationApiRequest
	conversationIDStr := c.Param("conversation_id")
	if conversationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "conversation_id is required",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	// 解析 conversation_id
	var conversationID int64
	if _, err := fmt.Sscanf(conversationIDStr, "%d", &conversationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid conversation_id",
		})
		return
	}
	req.ConversationID = &conversationID

	ctx := c.Request.Context()
	resp, err := application.ConversationSVC.UpdateConversation(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteConversationApi 删除对话（Open API）
// @router /v1/conversations/:conversation_id [DELETE]
func DeleteConversationApi(c *gin.Context) {
	var req conversation.DeleteConversationApiRequest
	conversationIDStr := c.Param("conversation_id")
	if conversationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "conversation_id is required",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	// 解析 conversation_id
	var conversationID int64
	if _, err := fmt.Sscanf(conversationIDStr, "%d", &conversationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid conversation_id",
		})
		return
	}
	req.ConversationID = &conversationID

	ctx := c.Request.Context()
	resp, err := application.ConversationSVC.DeleteConversation(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ChatV3 聊天（Open API v3）
// @router /v3/chat [POST]
func ChatV3(c *gin.Context) {
	var req run.ChatV3Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if checkErr := checkParamsV3(ctx, &req); checkErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  checkErr.Error(),
		})
		return
	}

	// 非流式模式
	if req.Stream != nil && !*req.Stream {
		resp, err := application.ConversationOpenAPISVC.OpenapiAgentRunSync(ctx, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	// 流式模式（默认）
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	sseSender := ginSSE.NewGinSSESender(c)
	err := application.ConversationOpenAPISVC.OpenapiAgentRun(ctx, sseSender, &req)
	if err != nil {
		errData := run.ErrorData{
			Code: errno.ErrConversationAgentRunError,
			Msg:  err.Error(),
		}
		ed, _ := json.Marshal(errData)
		_ = sseSender.Send(ctx, &sse.Event{
			Event: run.RunEventError,
			Data:  ed,
		})
	}
}

func checkParamsV3(_ interface{}, ar *run.ChatV3Request) error {
	if ar.BotID == 0 {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "bot id is required"))
	}
	return nil
}

// CancelChatApi 取消聊天（Open API v3）
// @router /v3/chat/cancel [POST]
func CancelChatApi(c *gin.Context) {
	var req run.CancelChatApiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := application.ConversationOpenAPISVC.CancelRun(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RetrieveChatOpen 检索聊天（Open API v3）
// @router /v3/chat/retrieve [GET]
func RetrieveChatOpen(c *gin.Context) {
	var req run.RetrieveChatOpenRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := application.ConversationOpenAPISVC.RetrieveRunRecord(ctx, &req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListChatMessageApi 获取聊天消息列表（Open API v3）
// @router /v3/chat/message/list [GET]
func ListChatMessageApi(c *gin.Context) {
	var req message.ListChatMessageApiRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := application.ConversationOpenAPISVC.ListChatMessageApi(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
