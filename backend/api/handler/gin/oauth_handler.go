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
	"github.com/coze-dev/coze-studio/backend/api/model/app/bot_open_api"
	"github.com/coze-dev/coze-studio/backend/application/plugin"
	"github.com/coze-dev/coze-studio/backend/bizpkg/config"
)

// OauthAuthorizationCode OAuth 授权码处理
// @router /api/oauth/authorization_code [GET]
func OauthAuthorizationCode(c *gin.Context) {
	var req bot_open_api.OauthAuthorizationCodeReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "authorization failed, code is required",
		})
		return
	}
	if req.State == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "state is required",
		})
		return
	}

	ctx := c.Request.Context()
	_, err := plugin.PluginApplicationSVC.OauthAuthorizationCode(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	host, err := config.Base().GetServerHost(ctx)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	redirectURL := fmt.Sprintf("%s/information/auth/success", host)
	c.Redirect(http.StatusFound, redirectURL)
}

