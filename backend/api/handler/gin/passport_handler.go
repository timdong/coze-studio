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
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
	"github.com/coze-dev/coze-studio/backend/api/model/passport"
	"github.com/coze-dev/coze-studio/backend/application/user"
	"github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/i18n"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

// PassportWebEmailRegisterV2Post 用户注册
// @router /api/passport/web/email/register/v2/ [POST]
func PassportWebEmailRegisterV2Post(c *gin.Context) {
	var req passport.PassportWebEmailRegisterV2PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	locale := string(i18n.GetLocale(ctx))

	resp, sessionKey, err := user.UserApplicationSVC.PassportWebEmailRegisterV2(ctx, locale, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "PassportWebEmailRegisterV2 failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  err.Error(),
		})
		return
	}

	// 设置 Cookie
	c.SetCookie(
		entity.SessionKey,
		sessionKey,
		consts.SessionMaxAgeSecond,
		"/",
		c.Request.Host,
		false,
		true,
	)

	c.JSON(http.StatusOK, resp)
}

// PassportWebEmailLoginPost 用户登录
// @router /api/passport/web/email/login/ [POST]
func PassportWebEmailLoginPost(c *gin.Context) {
	var req passport.PassportWebEmailLoginPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// 检查 UserApplicationSVC 是否已初始化
	if user.UserApplicationSVC == nil || user.UserApplicationSVC.DomainSVC == nil {
		logs.CtxErrorf(ctx, "UserApplicationSVC is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  "UserApplicationSVC is not initialized",
		})
		return
	}

	resp, sessionKey, err := user.UserApplicationSVC.PassportWebEmailLoginPost(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "PassportWebEmailLoginPost failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  err.Error(),
		})
		return
	}

	logs.Infof("[PassportWebEmailLoginPost] sessionKey: %s", sessionKey)

	// 设置 Cookie
	c.SetCookie(
		entity.SessionKey,
		sessionKey,
		consts.SessionMaxAgeSecond,
		"/",
		c.Request.Host,
		false,
		true,
	)

	c.JSON(http.StatusOK, resp)
}

// PassportWebLogoutGet 用户登出
// @router /api/passport/web/logout/ [GET]
func PassportWebLogoutGet(c *gin.Context) {
	var req passport.PassportWebLogoutGetRequest
	ctx := c.Request.Context()

	resp, err := user.UserApplicationSVC.PassportWebLogoutGet(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "PassportWebLogoutGet failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PassportWebEmailPasswordResetGet 重置密码
// @router /api/passport/web/email/password/reset/ [GET]
func PassportWebEmailPasswordResetGet(c *gin.Context) {
	var req passport.PassportWebEmailPasswordResetGetRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := user.UserApplicationSVC.PassportWebEmailPasswordResetGet(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "PassportWebEmailPasswordResetGet failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PassportAccountInfoV2 获取账户信息
// @router /api/passport/account/info/v2/ [POST]
func PassportAccountInfoV2(c *gin.Context) {
	var req passport.PassportAccountInfoV2Request
	ctx := c.Request.Context()

	// 检查是否有 session（这个接口用于检查登录状态，允许未登录）
	session := ctxutil.GetUserSessionFromCtx(ctx)
	if session == nil {
		// 未登录时返回成功响应，但 data 为空，前端可以据此判断用户未登录
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  "",
		})
		return
	}

	resp, err := user.UserApplicationSVC.PassportAccountInfoV2(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "PassportAccountInfoV2 failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UserUpdateAvatar 更新头像
// @router /api/web/user/update/upload_avatar/ [POST]
func UserUpdateAvatar(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		logs.CtxErrorf(c.Request.Context(), "Get Avatar Fail failed, err=%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "missing avatar file",
		})
		return
	}

	// 检查文件类型
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "invalid file type, only image allowed",
		})
		return
	}

	// 读取文件内容
	src, err := file.Open()
	if err != nil {
		logs.CtxErrorf(c.Request.Context(), "Open file failed, err=%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  err.Error(),
		})
		return
	}
	defer src.Close()

	fileContent, err := io.ReadAll(src)
	if err != nil {
		logs.CtxErrorf(c.Request.Context(), "Read file failed, err=%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  err.Error(),
		})
		return
	}

	req := passport.UserUpdateAvatarRequest{
		Avatar: fileContent,
	}
	mimeType := file.Header.Get("Content-Type")

	ctx := c.Request.Context()
	resp, err := user.UserApplicationSVC.UserUpdateAvatar(ctx, mimeType, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "UserUpdateAvatar failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UserUpdateProfile 更新用户资料
// @router /api/user/update_profile [POST]
func UserUpdateProfile(c *gin.Context) {
	var req passport.UserUpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := user.UserApplicationSVC.UserUpdateProfile(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "UserUpdateProfile failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateUserProfileCheck 更新用户资料检查
// @router /api/user/update_profile_check [POST]
func UpdateUserProfileCheck(c *gin.Context) {
	var req developer_api.UpdateUserProfileCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := user.UserApplicationSVC.UpdateUserProfileCheck(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

