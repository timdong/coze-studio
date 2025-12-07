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
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	"github.com/coze-dev/coze-studio/backend/application/modelmgr"
	"github.com/coze-dev/coze-studio/backend/application/upload"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

const baseWord = "1Aa2Bb3Cc4Dd5Ee6Ff7Gg8Hh9Ii0JjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"

// GetTypeList 获取类型列表
// @router /api/bot/get_type_list [POST]
func GetTypeList(c *gin.Context) {
	var req developer_api.GetTypeListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := modelmgr.ModelmgrApplicationSVC.GetModelList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UploadFile Bot 文件上传
// @router /api/bot/upload_file [POST]
func UploadFile(c *gin.Context) {
	var req developer_api.UploadFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp := new(developer_api.UploadFileResponse)

	fileContent, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	userID := ctxutil.GetUIDFromCtx(ctx)
	if userID == nil {
		httputil.InternalErrorGin(ctx, c, errorx.New(errno.ErrUploadPermissionCode, errorx.KV("msg", "session required")))
		return
	}

	secret := createSecret(ptr.From(userID), req.FileHead.FileType)
	fileName := fmt.Sprintf("%d_%d_%s.%s", ptr.From(userID), time.Now().UnixNano(), secret, req.FileHead.FileType)
	objectName := fmt.Sprintf("%s/%s", req.FileHead.BizType.String(), fileName)

	resp, err = upload.SVC.UploadFile(ctx, fileContent, objectName)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// createSecret 创建密钥
func createSecret(uid int64, fileType string) string {
	num := 10
	input := fmt.Sprintf("upload_%d_Ma*9)fhi_%d_gou_%s_rand_%d", uid, time.Now().Unix(), fileType, rand.Intn(100000))
	// Do sha256, take the first 20, map the number to Base62
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s", input)))
	hashString := base64.StdEncoding.EncodeToString(hash[:])
	if len(hashString) > num {
		hashString = hashString[:num]
	}

	result := ""
	for _, char := range hashString {
		index := int(char) % 62
		result += string(baseWord[index])
	}
	return result
}

