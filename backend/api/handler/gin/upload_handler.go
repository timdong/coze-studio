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
	upload "github.com/coze-dev/coze-studio/backend/api/model/file/upload"
	uploadSVC "github.com/coze-dev/coze-studio/backend/application/upload"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
)

// CommonUpload 通用文件上传
// @router /api/common/upload/:tos_uri [POST]
func CommonUpload(c *gin.Context) {
	var req upload.CommonUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	// 获取完整 URL，包括路径参数
	tosURI := c.Param("tos_uri")
	fullURL := c.Request.URL.Scheme + "://" + c.Request.Host + c.Request.URL.Path
	if tosURI != "" {
		fullURL = fullURL + "/" + tosURI
	}
	if c.Request.URL.RawQuery != "" {
		fullURL = fullURL + "?" + c.Request.URL.RawQuery
	}

	resp, err := uploadSVC.SVC.UploadFileCommon(ctx, &req, fullURL)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ApplyUploadAction 申请上传操作
// @router /api/common/upload/apply_upload_action [GET, POST]
func ApplyUploadAction(c *gin.Context) {
	var req upload.ApplyUploadActionRequest
	var err error

	// 根据请求方法绑定参数
	if c.Request.Method == http.MethodGet {
		err = c.ShouldBindQuery(&req)
	} else {
		err = c.ShouldBindJSON(&req)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	host := c.Request.Host
	resp := new(upload.ApplyUploadActionResponse)

	if ptr.From(req.Action) == "ApplyImageUpload" {
		resp, err = uploadSVC.SVC.ApplyImageUpload(ctx, &req, host)
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
	} else if ptr.From(req.Action) == "CommitImageUpload" {
		resp, err = uploadSVC.SVC.CommitImageUpload(ctx, &req, host)
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
	}

	c.JSON(http.StatusOK, resp)
}
