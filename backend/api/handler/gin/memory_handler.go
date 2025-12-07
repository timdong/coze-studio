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

	"github.com/gin-gonic/gin"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/data/variable/kvmemory"
	"github.com/coze-dev/coze-studio/backend/api/model/data/variable/project_memory"
	appApplication "github.com/coze-dev/coze-studio/backend/application/app"
	"github.com/coze-dev/coze-studio/backend/application/memory"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/conv"
)

// GetSysVariableConf 获取系统变量配置
// @router /api/memory/sys_variable_conf [GET]
func GetSysVariableConf(c *gin.Context) {
	var req kvmemory.GetSysVariableConfRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.VariableApplicationSVC.GetSysVariableConf(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProjectVariableList 获取项目变量列表
// @router /api/memory/project/variable/meta_list [GET]
func GetProjectVariableList(c *gin.Context) {
	var req project_memory.GetProjectVariableListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if req.ProjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "project_id is empty",
		})
		return
	}

	ctx := c.Request.Context()
	pID, err := conv.StrToInt64(req.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "project_id is not int",
		})
		return
	}

	pInfo, err := appApplication.APPApplicationSVC.DomainSVC.GetDraftAPP(ctx, pID)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	resp, err := memory.VariableApplicationSVC.GetProjectVariablesMeta(ctx, pInfo.OwnerID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateProjectVariable 更新项目变量
// @router /api/memory/project/variable/meta_update [POST]
func UpdateProjectVariable(c *gin.Context) {
	var req project_memory.UpdateProjectVariableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if req.ProjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "project_id is empty",
		})
		return
	}

	key2Var := make(map[string]*project_memory.Variable)
	for _, v := range req.VariableList {
		if v.Keyword == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  "variable name is empty",
			})
			return
		}

		if key2Var[v.Keyword] != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  "variable keyword is duplicate",
			})
			return
		}

		key2Var[v.Keyword] = v
	}

	ctx := c.Request.Context()
	resp, err := memory.VariableApplicationSVC.UpdateProjectVariable(ctx, req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SetKvMemory 设置键值对记忆
// @router /api/memory/variable/upsert [POST]
func SetKvMemory(c *gin.Context) {
	var req kvmemory.SetKvMemoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if req.BotID == 0 && req.GetProjectID() == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "bot_id and project_id are both empty",
		})
		return
	}

	if len(req.Data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "data is empty",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.VariableApplicationSVC.SetVariableInstance(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMemoryVariableMeta 获取记忆变量元数据
// @router /api/memory/variable/get_meta [POST]
func GetMemoryVariableMeta(c *gin.Context) {
	var req project_memory.GetMemoryVariableMetaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.VariableApplicationSVC.GetVariableMeta(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DelProfileMemory 删除用户画像记忆
// @router /api/memory/variable/delete [POST]
func DelProfileMemory(c *gin.Context) {
	var req kvmemory.DelProfileMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if req.BotID == 0 && req.GetProjectID() == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "bot_id and project_id are both empty",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.VariableApplicationSVC.DeleteVariableInstance(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPlayGroundMemory 获取 Playground 记忆
// @router /api/memory/variable/get [POST]
func GetPlayGroundMemory(c *gin.Context) {
	var req kvmemory.GetProfileMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if req.BotID == 0 && req.GetProjectID() == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "bot_id and project_id are both empty",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.VariableApplicationSVC.GetPlayGroundMemory(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

