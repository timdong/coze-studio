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
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/app/bot_open_api"
	"github.com/coze-dev/coze-studio/backend/api/model/permission/openapiauth"
	openapiauthApp "github.com/coze-dev/coze-studio/backend/application/openauth"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// GetPersonalAccessTokenAndPermission 获取个人访问令牌和权限
// @router /api/permission_api/pat/get_personal_access_token_and_permission [GET]
func GetPersonalAccessTokenAndPermission(c *gin.Context) {
	var req openapiauth.GetPersonalAccessTokenAndPermissionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "id is required",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := openapiauthApp.OpenAuthApplication.GetPersonalAccessTokenAndPermission(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "OpenAuthApplicationService.GetPersonalAccessTokenAndPermission failed, err=%v", err)
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeletePersonalAccessTokenAndPermission 删除个人访问令牌和权限
// @router /api/permission_api/pat/delete_personal_access_token_and_permission [POST]
func DeletePersonalAccessTokenAndPermission(c *gin.Context) {
	var req openapiauth.DeletePersonalAccessTokenAndPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "id is required",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := openapiauthApp.OpenAuthApplication.DeletePersonalAccessTokenAndPermission(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "OpenAuthApplication.DeletePersonalAccessTokenAndPermission failed, err=%v", err)
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListPersonalAccessTokens 列出个人访问令牌
// @router /api/permission_api/pat/list_personal_access_tokens [GET]
func ListPersonalAccessTokens(c *gin.Context) {
	var req openapiauth.ListPersonalAccessTokensRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.Page == nil || *req.Page <= 0 {
		req.Page = ptr.Of(int64(1))
	}
	if req.Size == nil || *req.Size <= 0 {
		req.Size = ptr.Of(int64(10))
	}

	ctx := c.Request.Context()
	resp, err := openapiauthApp.OpenAuthApplication.ListPersonalAccessTokens(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "OpenAuthApplication.ListPersonalAccessTokens failed, err=%v", err)
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreatePersonalAccessTokenAndPermission 创建个人访问令牌和权限
// @router /api/permission_api/pat/create_personal_access_token_and_permission [POST]
func CreatePersonalAccessTokenAndPermission(c *gin.Context) {
	var req openapiauth.CreatePersonalAccessTokenAndPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if err := checkCPATParams(ctx, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	resp, err := openapiauthApp.OpenAuthApplication.CreatePersonalAccessToken(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "OpenAuthApplicationService.CreatePersonalAccessToken failed, err=%v", err)
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// checkCPATParams 检查创建个人访问令牌的参数
func checkCPATParams(ctx context.Context, req *openapiauth.CreatePersonalAccessTokenAndPermissionRequest) error {
	if req.Name == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode, errorx.KV("msg", "name is required"))
	}
	return nil
}

// UpdatePersonalAccessTokenAndPermission 更新个人访问令牌和权限
// @router /api/permission_api/pat/update_personal_access_token_and_permission [POST]
func UpdatePersonalAccessTokenAndPermission(c *gin.Context) {
	var req openapiauth.UpdatePersonalAccessTokenAndPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := openapiauthApp.OpenAuthApplication.UpdatePersonalAccessTokenAndPermission(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "OpenAuthApplication.UpdatePersonalAccessTokenAndPermission failed, err=%v", err)
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ImpersonateCozeUser 模拟 Coze 用户
// @router /api/permission_api/coze_web_app/impersonate_coze_user [POST]
func ImpersonateCozeUser(c *gin.Context) {
	var req bot_open_api.ImpersonateCozeUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := openapiauthApp.OpenAuthApplication.ImpersonateCozeUserAccessToken(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

