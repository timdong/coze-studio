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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/app/intelligence"
	"github.com/coze-dev/coze-studio/backend/api/model/app/intelligence/common"
	project "github.com/coze-dev/coze-studio/backend/api/model/app/intelligence/project"
	publish "github.com/coze-dev/coze-studio/backend/api/model/app/intelligence/publish"
	task "github.com/coze-dev/coze-studio/backend/api/model/app/intelligence/task"
	appApplication "github.com/coze-dev/coze-studio/backend/application/app"
	"github.com/coze-dev/coze-studio/backend/application/search"
)

// GetDraftIntelligenceList 获取草稿智能体列表
// @router /api/intelligence_api/search/get_draft_intelligence_list [POST]
func GetDraftIntelligenceList(c *gin.Context) {
	var req intelligence.GetDraftIntelligenceListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := search.SearchSVC.GetDraftIntelligenceList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDraftIntelligenceInfo 获取草稿智能体信息
// @router /api/intelligence_api/search/get_draft_intelligence_info [POST]
func GetDraftIntelligenceInfo(c *gin.Context) {
	var req intelligence.GetDraftIntelligenceInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.IntelligenceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid intelligence id",
		})
		return
	}
	if req.IntelligenceType != common.IntelligenceType_Project {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  fmt.Sprintf("invalid intelligence type '%d'", req.IntelligenceType),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.GetDraftIntelligenceInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserRecentlyEditIntelligence 获取用户最近编辑的智能体
// @router /api/intelligence_api/search/get_recently_edit_intelligence [POST]
func GetUserRecentlyEditIntelligence(c *gin.Context) {
	var req intelligence.GetUserRecentlyEditIntelligenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := new(intelligence.GetUserRecentlyEditIntelligenceResponse)
	c.JSON(http.StatusOK, resp)
}

// DraftProjectCreate 创建草稿项目
// @router /api/intelligence_api/draft_project/create [POST]
func DraftProjectCreate(c *gin.Context) {
	var req project.DraftProjectCreateRequest
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
			"msg":  "invalid space id",
		})
		return
	}
	if req.Name == "" || len(req.Name) > 256 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid name",
		})
		return
	}
	if req.IconURI == "" || len(req.IconURI) > 512 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid icon uri",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.DraftProjectCreate(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DraftProjectUpdate 更新草稿项目
// @router /api/intelligence_api/draft_project/update [POST]
func DraftProjectUpdate(c *gin.Context) {
	var req project.DraftProjectUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ProjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid project id",
		})
		return
	}
	if req.Name != nil && (len(*req.Name) == 0 || len(*req.Name) > 256) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid name",
		})
		return
	}
	if req.IconURI != nil && (len(*req.IconURI) == 0 || len(*req.IconURI) > 512) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid icon uri",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.DraftProjectUpdate(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DraftProjectDelete 删除草稿项目
// @router /api/intelligence_api/draft_project/delete [POST]
func DraftProjectDelete(c *gin.Context) {
	var req project.DraftProjectDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ProjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid project id",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.DraftProjectDelete(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DraftProjectCopy 复制草稿项目
// @router /api/intelligence_api/draft_project/copy [POST]
func DraftProjectCopy(c *gin.Context) {
	var req project.DraftProjectCopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ProjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid project id",
		})
		return
	}
	if req.ToSpaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid to space id",
		})
		return
	}
	if req.Name == "" || len(req.Name) > 256 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid name",
		})
		return
	}
	if req.IconURI == "" || len(req.IconURI) > 512 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid icon uri",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.DraftProjectCopy(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DraftProjectInnerTaskList 获取草稿项目内部任务列表
// @router /api/intelligence_api/draft_project/inner_task_list [POST]
func DraftProjectInnerTaskList(c *gin.Context) {
	var req task.DraftProjectInnerTaskListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ProjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid project id",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.DraftProjectInnerTaskList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProjectPublishedConnector 获取项目已发布的连接器
// @router /api/intelligence_api/publish/get_published_connector [POST]
func GetProjectPublishedConnector(c *gin.Context) {
	var req publish.GetProjectPublishedConnectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := new(publish.GetProjectPublishedConnectorResponse)
	c.JSON(http.StatusOK, resp)
}

// CheckProjectVersionNumber 检查项目版本号
// @router /api/intelligence_api/publish/check_version_number [POST]
func CheckProjectVersionNumber(c *gin.Context) {
	var req publish.CheckProjectVersionNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ProjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid project id",
		})
		return
	}
	if req.VersionNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid version number",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.CheckProjectVersionNumber(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PublishProject 发布项目
// @router /api/intelligence_api/publish/publish_project [POST]
func PublishProject(c *gin.Context) {
	var req publish.PublishProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ProjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid project id",
		})
		return
	}
	if req.VersionNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid version number",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.PublishAPP(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPublishRecordList 获取发布记录列表
// @router /api/intelligence_api/publish/publish_record_list [POST]
func GetPublishRecordList(c *gin.Context) {
	var req publish.GetPublishRecordListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ProjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid project id",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.GetPublishRecordList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ProjectPublishConnectorList 获取项目发布连接器列表
// @router /api/intelligence_api/publish/connector_list [POST]
func ProjectPublishConnectorList(c *gin.Context) {
	var req publish.PublishConnectorListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ProjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid project id",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.ProjectPublishConnectorList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPublishRecordDetail 获取发布记录详情
// @router /api/intelligence_api/publish/publish_record_detail [POST]
func GetPublishRecordDetail(c *gin.Context) {
	var req publish.GetPublishRecordDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ProjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid project id",
		})
		return
	}
	if req.PublishRecordID != nil && *req.PublishRecordID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid publish record id",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.GetPublishRecordDetail(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetOnlineAppData 获取在线应用数据
// @router /v1/apps/:app_id [GET]
func GetOnlineAppData(c *gin.Context) {
	var req project.GetOnlineAppDataRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := appApplication.APPApplicationSVC.GetOnlineAppData(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

