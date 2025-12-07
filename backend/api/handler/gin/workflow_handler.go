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
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
	"github.com/hertz-contrib/sse"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	gin_sse "github.com/coze-dev/coze-studio/backend/infra/sse/impl/gin"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
)

// CreateWorkflow 创建工作流
// @router /api/workflow_api/create [POST]
func CreateWorkflow(c *gin.Context) {
	var req workflow.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.CreateWorkflow(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetCanvasInfo 获取画布信息
// @router /api/workflow_api/canvas [POST]
func GetCanvasInfo(c *gin.Context) {
	var req workflow.GetCanvasInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.GetCanvasInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SaveWorkflow 保存工作流
// @router /api/workflow_api/save [POST]
func SaveWorkflow(c *gin.Context) {
	var req workflow.SaveWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.SaveWorkflow(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateWorkflowMeta 更新工作流元数据
// @router /api/workflow_api/update_meta [POST]
func UpdateWorkflowMeta(c *gin.Context) {
	var req workflow.UpdateWorkflowMetaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.UpdateWorkflowMeta(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteWorkflow 删除工作流
// @router /api/workflow_api/delete [POST]
func DeleteWorkflow(c *gin.Context) {
	var req workflow.DeleteWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.DeleteWorkflow(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// BatchDeleteWorkflow 批量删除工作流
// @router /api/workflow_api/batch_delete [POST]
func BatchDeleteWorkflow(c *gin.Context) {
	var req workflow.BatchDeleteWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.BatchDeleteWorkflow(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDeleteStrategy 获取删除策略
// @router /api/workflow_api/delete_strategy [POST]
func GetDeleteStrategy(c *gin.Context) {
	var req workflow.GetDeleteStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	resp := new(workflow.GetDeleteStrategyResponse)
	c.JSON(http.StatusOK, resp)
}

// PublishWorkflow 发布工作流
// @router /api/workflow_api/publish [POST]
func PublishWorkflow(c *gin.Context) {
	var req workflow.PublishWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.PublishWorkflow(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CopyWorkflow 复制工作流
// @router /api/workflow_api/copy [POST]
func CopyWorkflow(c *gin.Context) {
	var req workflow.CopyWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.CopyWorkflow(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CopyWkTemplateApi 复制工作流模板
// @router /api/workflow_api/copy_wk_template [POST]
func CopyWkTemplateApi(c *gin.Context) {
	var req workflow.CopyWkTemplateApiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.CopyWkTemplateApi(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetReleasedWorkflows 获取已发布的工作流
// @router /api/workflow_api/released_workflows [POST]
func GetReleasedWorkflows(c *gin.Context) {
	var req workflow.GetReleasedWorkflowsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	resp := new(workflow.GetReleasedWorkflowsResponse)
	c.JSON(http.StatusOK, resp)
}

// GetWorkflowReferences 获取工作流引用
// @router /api/workflow_api/workflow_references [POST]
func GetWorkflowReferences(c *gin.Context) {
	var req workflow.GetWorkflowReferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.GetWorkflowReferences(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetWorkFlowList 获取工作流列表
// @router /api/workflow_api/workflow_list [POST]
func GetWorkFlowList(c *gin.Context) {
	var req workflow.GetWorkFlowListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.ListWorkflow(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// QueryWorkflowNodeTypes 查询工作流节点类型
// @router /api/workflow_api/node_type [POST]
func QueryWorkflowNodeTypes(c *gin.Context) {
	var req workflow.QueryWorkflowNodeTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.QueryWorkflowNodeTypes(ctx, &req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// NodeTemplateList 获取节点模板列表
// @router /api/workflow_api/node_template_list [POST]
func NodeTemplateList(c *gin.Context) {
	var req workflow.NodeTemplateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.GetNodeTemplateList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// NodePanelSearch 节点面板搜索
// @router /api/workflow_api/node_panel_search [POST]
func NodePanelSearch(c *gin.Context) {
	var req workflow.NodePanelSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	resp := new(workflow.NodePanelSearchResponse)
	c.JSON(http.StatusOK, resp)
}

// GetLLMNodeFCSettingsMerged 获取 LLM 节点 FC 设置合并
// @router /api/workflow_api/llm_fc_setting_merged [POST]
func GetLLMNodeFCSettingsMerged(c *gin.Context) {
	var req workflow.GetLLMNodeFCSettingsMergedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.GetLLMNodeFCSettingsMerged(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetLLMNodeFCSettingDetail 获取 LLM 节点 FC 设置详情
// @router /api/workflow_api/llm_fc_setting_detail [POST]
func GetLLMNodeFCSettingDetail(c *gin.Context) {
	var req workflow.GetLLMNodeFCSettingDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.GetLLMNodeFCSettingDetail(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// WorkFlowTestRun 工作流测试运行
// @router /api/workflow_api/test_run [POST]
func WorkFlowTestRun(c *gin.Context) {
	var req workflow.WorkFlowTestRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.TestRun(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// WorkFlowTestResume 工作流测试恢复
// @router /api/workflow_api/test_resume [POST]
func WorkFlowTestResume(c *gin.Context) {
	var req workflow.WorkflowTestResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.TestResume(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CancelWorkFlow 取消工作流
// @router /api/workflow_api/cancel [POST]
func CancelWorkFlow(c *gin.Context) {
	var req workflow.CancelWorkFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.Cancel(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetWorkFlowProcess 获取工作流进程
// @router /api/workflow_api/get_process [GET]
func GetWorkFlowProcess(c *gin.Context) {
	var req workflow.GetWorkflowProcessRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.GetProcess(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetNodeExecuteHistory 获取节点执行历史
// @router /api/workflow_api/get_node_execute_history [GET]
func GetNodeExecuteHistory(c *gin.Context) {
	var req workflow.GetNodeExecuteHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.GetNodeExecuteHistory(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetApiDetail 获取 API 详情
// @router /api/workflow_api/apiDetail [GET]
func GetApiDetail(c *gin.Context) {
	var req workflow.GetApiDetailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	toolDetailInfo, err := appworkflow.SVC.GetApiDetail(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	response := map[string]interface{}{
		"data": toolDetailInfo,
		"code": 0,
		"msg":  "",
	}

	c.JSON(http.StatusOK, response)
}

// WorkflowNodeDebugV2 工作流节点调试 V2
// @router /api/workflow_api/nodeDebug [POST]
func WorkflowNodeDebugV2(c *gin.Context) {
	var req workflow.WorkflowNodeDebugV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.NodeDebug(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SignImageURL 签名图片 URL
// @router /api/workflow_api/sign_image_url [POST]
func SignImageURL(c *gin.Context) {
	var req workflow.SignImageURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.SignImageURL(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateProjectConversationDef 创建项目对话定义
// @router /api/workflow_api/project_conversation/create [POST]
func CreateProjectConversationDef(c *gin.Context) {
	var req workflow.CreateProjectConversationDefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.CreateApplicationConversationDef(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateProjectConversationDef 更新项目对话定义
// @router /api/workflow_api/project_conversation/update [POST]
func UpdateProjectConversationDef(c *gin.Context) {
	var req workflow.UpdateProjectConversationDefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.UpdateApplicationConversationDef(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteProjectConversationDef 删除项目对话定义
// @router /api/workflow_api/project_conversation/delete [POST]
func DeleteProjectConversationDef(c *gin.Context) {
	var req workflow.DeleteProjectConversationDefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.DeleteApplicationConversationDef(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListProjectConversationDef 列出项目对话定义
// @router /api/workflow_api/project_conversation/list [GET]
func ListProjectConversationDef(c *gin.Context) {
	var req workflow.ListProjectConversationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.ListApplicationConversationDef(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListRootSpans 列出根 Spans
// @router /api/workflow_api/list_spans [POST]
func ListRootSpans(c *gin.Context) {
	var req workflow.ListRootSpansRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	resp := new(workflow.ListRootSpansResponse)
	c.JSON(http.StatusOK, resp)
}

// GetTraceSDK 获取 Trace SDK
// @router /api/workflow_api/get_trace [POST]
func GetTraceSDK(c *gin.Context) {
	var req workflow.GetTraceSDKRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	resp := new(workflow.GetTraceSDKResponse)
	c.JSON(http.StatusOK, resp)
}

// GetWorkflowDetail 获取工作流详情
// @router /api/workflow_api/workflow_detail [POST]
func GetWorkflowDetail(c *gin.Context) {
	var req workflow.GetWorkflowDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	workflowDetailDataList, err := appworkflow.SVC.GetWorkflowDetail(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	response := map[string]any{
		"data":    workflowDetailDataList,
		"code":    0,
		"message": "",
	}

	c.JSON(http.StatusOK, response)
}

// GetWorkflowDetailInfo 获取工作流详情信息
// @router /api/workflow_api/workflow_detail_info [POST]
func GetWorkflowDetailInfo(c *gin.Context) {
	var req workflow.GetWorkflowDetailInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	workflowDetailInfoDataList, err := appworkflow.SVC.GetWorkflowDetailInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	response := map[string]any{
		"data":    workflowDetailInfoDataList,
		"code":    0,
		"message": "",
	}

	c.JSON(http.StatusOK, response)
}

// ValidateTree 验证树
// @router /api/workflow_api/validate_tree [POST]
func ValidateTree(c *gin.Context) {
	var req workflow.ValidateTreeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.ValidateTree(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetChatFlowRole 获取聊天流角色
// @router /api/workflow_api/chat_flow_role/get [GET]
func GetChatFlowRole(c *gin.Context) {
	var req workflow.GetChatFlowRoleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.GetChatFlowRole(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateChatFlowRole 创建聊天流角色
// @router /api/workflow_api/chat_flow_role/create [POST]
func CreateChatFlowRole(c *gin.Context) {
	var req workflow.CreateChatFlowRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.CreateChatFlowRole(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteChatFlowRole 删除聊天流角色
// @router /api/workflow_api/chat_flow_role/delete [POST]
func DeleteChatFlowRole(c *gin.Context) {
	var req workflow.DeleteChatFlowRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.DeleteChatFlowRole(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListPublishWorkflow 列出发布的工作流
// @router /api/workflow_api/list_publish_workflow [POST]
func ListPublishWorkflow(c *gin.Context) {
	var req workflow.ListPublishWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	resp := new(workflow.ListPublishWorkflowResponse)
	c.JSON(http.StatusOK, resp)
}

// GetWorkflowUploadAuthToken 获取工作流上传授权令牌
// @router /api/workflow_api/upload/auth_token [POST]
func GetWorkflowUploadAuthToken(c *gin.Context) {
	var req workflow.GetUploadAuthTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.GetWorkflowUploadAuthToken(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetHistorySchema 获取历史 Schema
// @router /api/workflow_api/history_schema [POST]
func GetHistorySchema(c *gin.Context) {
	var req workflow.GetHistorySchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.GetHistorySchema(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetExampleWorkFlowList 获取示例工作流列表
// @router /api/workflow_api/example_workflow_list [POST]
func GetExampleWorkFlowList(c *gin.Context) {
	var req workflow.GetExampleWorkFlowListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.GetExampleWorkFlowList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// preprocessWorkflowRequestBody 预处理工作流请求体，将 parameters 字段从对象转换为字符串
func preprocessWorkflowRequestBody(c *gin.Context) error {
	// 读取原始请求体
	rawData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	// 解析为临时 map
	var bodyData map[string]interface{}
	if err = sonic.Unmarshal(rawData, &bodyData); err != nil {
		// 如果解析失败，恢复原始 body 并返回错误
		c.Request.Body = io.NopCloser(bytes.NewReader(rawData))
		return fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	// 处理 'parameters' 字段
	if parameters, ok := bodyData["parameters"]; ok {
		if _, isString := parameters.(string); !isString {
			// 不是字符串，需要转换
			paramsBytes, marshalErr := sonic.Marshal(parameters)
			if marshalErr != nil {
				c.Request.Body = io.NopCloser(bytes.NewReader(rawData))
				return fmt.Errorf("failed to marshal parameters: %w", marshalErr)
			}
			bodyData["parameters"] = string(paramsBytes)

			newRawData, err := sonic.Marshal(bodyData)
			if err != nil {
				c.Request.Body = io.NopCloser(bytes.NewReader(rawData))
				return fmt.Errorf("failed to marshal modified body: %w", err)
			}
			// 重新设置请求体
			c.Request.Body = io.NopCloser(bytes.NewReader(newRawData))
			return nil
		}
	}

	// 如果没有修改，恢复原始 body
	c.Request.Body = io.NopCloser(bytes.NewReader(rawData))
	return nil
}

// processOpenAPIGetWorkflowInfoRequest 处理 Open API 获取工作流信息请求
func processOpenAPIGetWorkflowInfoRequest(c *gin.Context) error {
	isDebug := c.Query("is_debug")
	if isDebug == "" {
		values := c.Request.URL.Query()
		values.Set("is_debug", "false")
		c.Request.URL.RawQuery = values.Encode()
	}
	return nil
}

// streamRunData 流式运行数据结构
type streamRunData struct {
	Content       *string        `json:"content,omitempty"`
	ContentType   *string        `json:"content_type,omitempty"`
	NodeSeqID     *string        `json:"node_seq_id,omitempty"`
	NodeID        *string        `json:"node_id,omitempty"`
	NodeIsFinish  *bool          `json:"node_is_finish,omitempty"`
	NodeType      *string        `json:"node_type,omitempty"`
	NodeTitle     *string        `json:"node_title,omitempty"`
	Token         *int64         `json:"token,omitempty"`
	DebugURL      *string        `json:"debug_url,omitempty"`
	ErrorCode     *int64         `json:"error_code,omitempty"`
	ErrorMessage  *string        `json:"error_message,omitempty"`
	InterruptData *interruptData `json:"interrupt_data,omitempty"`
}

// interruptData 中断数据结构
type interruptData struct {
	EventID string `json:"event_id"`
	Type    int64  `json:"type"`
	Data    string `json:"data"`
}

// convertStreamRunData 转换流式运行数据
func convertStreamRunData(msg *workflow.OpenAPIStreamRunFlowResponse) *streamRunData {
	var ie *interruptData
	if msg.InterruptData != nil {
		ie = &interruptData{
			EventID: msg.InterruptData.EventID,
			Type:    int64(msg.InterruptData.Type),
			Data:    msg.InterruptData.InData,
		}
	}

	return &streamRunData{
		Content:       msg.Content,
		ContentType:   msg.ContentType,
		NodeSeqID:     msg.NodeSeqID,
		NodeID:        msg.NodeID,
		NodeIsFinish:  msg.NodeIsFinish,
		NodeType:      msg.NodeType,
		NodeTitle:     msg.NodeTitle,
		Token:         msg.Token,
		DebugURL:      msg.DebugUrl,
		ErrorCode:     msg.ErrorCode,
		ErrorMessage:  msg.ErrorMessage,
		InterruptData: ie,
	}
}

// sendStreamRunSSE 发送流式运行 SSE 事件
func sendStreamRunSSE(ctx context.Context, c *gin.Context, sr *schema.StreamReader[*workflow.OpenAPIStreamRunFlowResponse]) {
	sender := gin_sse.NewGinSSESender(c)
	defer func() {
		sr.Close()
	}()

	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// 完成
				break
			}

			event := &sse.Event{
				Event: "error",
				Data:  []byte(err.Error()),
			}

			if err = sender.Send(ctx, event); err != nil {
				logs.CtxErrorf(ctx, "publish stream event failed, err:%v", err)
			}
			return
		}

		converted := convertStreamRunData(msg)
		msgBytes, err := sonic.Marshal(converted)
		if err != nil {
			event := &sse.Event{
				Event: "error",
				Data:  []byte(err.Error()),
			}
			if err = sender.Send(ctx, event); err != nil {
				logs.CtxErrorf(ctx, "publish stream event failed, err:%v", err)
			}
			return
		}

		event := &sse.Event{
			ID:    msg.ID,
			Event: msg.Event,
			Data:  msgBytes,
		}

		if err = sender.Send(ctx, event); err != nil {
			logs.CtxErrorf(ctx, "publish stream event failed, err:%v", err)
			return
		}
	}
}

// sendChatFlowStreamRunSSE 发送聊天流流式运行 SSE 事件
func sendChatFlowStreamRunSSE(ctx context.Context, c *gin.Context, sr *schema.StreamReader[[]*workflow.ChatFlowRunResponse]) {
	sender := gin_sse.NewGinSSESender(c)
	defer func() {
		sr.Close()
	}()

	seq := int64(1)
	for {
		respList, err := sr.Recv()

		if err != nil {
			if errors.Is(err, io.EOF) {
				// 完成
				break
			}

			event := &sse.Event{
				Event: "error",
				Data:  []byte(err.Error()),
			}

			if err = sender.Send(ctx, event); err != nil {
				logs.CtxErrorf(ctx, "publish stream event failed, err:%v", err)
			}
			return
		}

		for _, resp := range respList {
			event := &sse.Event{
				ID:    strconv.FormatInt(seq, 10),
				Event: resp.Event,
				Data:  []byte(resp.Data),
			}

			if err = sender.Send(ctx, event); err != nil {
				logs.CtxErrorf(ctx, "publish stream event failed, err:%v", err)
				return
			}
			seq++
		}
	}
}

// OpenAPIRunFlow Open API 运行工作流
// @router /v1/workflow/run [POST]
func OpenAPIRunFlow(c *gin.Context) {
	if err := preprocessWorkflowRequestBody(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	var req workflow.OpenAPIRunFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.OpenAPIRun(ctx, &req)
	if err != nil {
		var se vo.WorkflowError
		if errors.As(err, &se) {
			resp = new(workflow.OpenAPIRunFlowResponse)
			resp.Code = int64(se.OpenAPICode())
			resp.Msg = ptr.Of(se.Msg())
			debugURL := se.DebugURL()
			if debugURL != "" {
				resp.DebugUrl = ptr.Of(debugURL)
			}
			c.JSON(http.StatusOK, resp)
			return
		}

		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// OpenAPIStreamRunFlow Open API 流式运行工作流
// @router /v1/workflow/stream_run [POST]
func OpenAPIStreamRunFlow(c *gin.Context) {
	if err := preprocessWorkflowRequestBody(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	var req workflow.OpenAPIRunFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	ctx := c.Request.Context()
	sr, err := appworkflow.SVC.OpenAPIStreamRun(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	sendStreamRunSSE(ctx, c, sr)
}

// OpenAPIStreamResumeFlow Open API 流式恢复工作流
// @router /v1/workflow/stream_resume [POST]
func OpenAPIStreamResumeFlow(c *gin.Context) {
	var req workflow.OpenAPIStreamResumeFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	ctx := c.Request.Context()
	sr, err := appworkflow.SVC.OpenAPIStreamResume(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	sendStreamRunSSE(ctx, c, sr)
}

// OpenAPIGetWorkflowRunHistory Open API 获取工作流运行历史
// @router /v1/workflow/get_run_history [GET]
func OpenAPIGetWorkflowRunHistory(c *gin.Context) {
	var req workflow.GetWorkflowRunHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.OpenAPIGetWorkflowRunHistory(ctx, &req)
	if err != nil {
		var se vo.WorkflowError
		if errors.As(err, &se) {
			resp = new(workflow.GetWorkflowRunHistoryResponse)
			resp.Code = ptr.Of(int64(se.OpenAPICode()))
			resp.Msg = ptr.Of(se.Msg())
			c.JSON(http.StatusOK, resp)
			return
		}

		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// OpenAPIChatFlowRun Open API 聊天流运行
// @router /v1/workflows/chat [POST]
func OpenAPIChatFlowRun(c *gin.Context) {
	if err := preprocessWorkflowRequestBody(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	var req workflow.ChatFlowRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	ctx := c.Request.Context()
	sr, err := appworkflow.SVC.OpenAPIChatFlowRun(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	sendChatFlowStreamRunSSE(ctx, c, sr)
}

// OpenAPIGetWorkflowInfo Open API 获取工作流信息
// @router /v1/workflows/:workflow_id [GET]
func OpenAPIGetWorkflowInfo(c *gin.Context) {
	if err := processOpenAPIGetWorkflowInfoRequest(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	var req workflow.OpenAPIGetWorkflowInfoRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.OpenAPIGetWorkflowInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// OpenAPICreateConversation Open API 创建对话
// @router /v1/workflow/conversation/create [POST]
func OpenAPICreateConversation(c *gin.Context) {
	var req workflow.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := appworkflow.SVC.OpenAPICreateConversation(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

