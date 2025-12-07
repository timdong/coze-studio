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

	dataset "github.com/coze-dev/coze-studio/backend/api/model/data/knowledge"
	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/application/knowledge"
	application "github.com/coze-dev/coze-studio/backend/application/knowledge"
	"github.com/coze-dev/coze-studio/backend/application/memory"
	"github.com/coze-dev/coze-studio/backend/application/upload"
)

// CreateDataset 创建数据集
// @router /api/knowledge/create [POST]
func CreateDataset(c *gin.Context) {
	var req dataset.CreateDatasetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.CreateKnowledge(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DatasetDetail 获取数据集详情
// @router /api/knowledge/detail [POST]
func DatasetDetail(c *gin.Context) {
	var req dataset.DatasetDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.DatasetDetail(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListDataset 列出数据集
// @router /api/knowledge/list [POST]
func ListDataset(c *gin.Context) {
	var req dataset.ListDatasetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.ListKnowledge(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteDataset 删除数据集
// @router /api/knowledge/delete [POST]
func DeleteDataset(c *gin.Context) {
	var req dataset.DeleteDatasetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.DeleteKnowledge(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateDataset 更新数据集
// @router /api/knowledge/update [POST]
func UpdateDataset(c *gin.Context) {
	var req dataset.UpdateDatasetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.UpdateKnowledge(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateDocument 创建文档
// @router /api/knowledge/document/create [POST]
func CreateDocument(c *gin.Context) {
	var req dataset.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.CreateDocument(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListDocument 列出文档
// @router /api/knowledge/document/list [POST]
func ListDocument(c *gin.Context) {
	var req dataset.ListDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.ListDocument(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteDocument 删除文档
// @router /api/knowledge/document/delete [POST]
func DeleteDocument(c *gin.Context) {
	var req dataset.DeleteDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.DeleteDocument(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateDocument 更新文档
// @router /api/knowledge/document/update [POST]
func UpdateDocument(c *gin.Context) {
	var req dataset.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.UpdateDocument(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDocumentProgress 获取文档处理进度
// @router /api/knowledge/document/progress/get [POST]
func GetDocumentProgress(c *gin.Context) {
	var req dataset.GetDocumentProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.GetDocumentProgress(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Resegment 重新分段文档
// @router /api/knowledge/document/resegment [POST]
func Resegment(c *gin.Context) {
	var req dataset.ResegmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.Resegment(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdatePhotoCaption 更新照片标题
// @router /api/knowledge/photo/caption [POST]
func UpdatePhotoCaption(c *gin.Context) {
	var req dataset.UpdatePhotoCaptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.UpdatePhotoCaption(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListPhoto 列出照片
// @router /api/knowledge/photo/list [POST]
func ListPhoto(c *gin.Context) {
	var req dataset.ListPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.ListPhoto(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PhotoDetail 获取照片详情
// @router /api/knowledge/photo/detail [POST]
func PhotoDetail(c *gin.Context) {
	var req dataset.PhotoDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.PhotoDetail(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTableSchema 获取表结构
// @router /api/knowledge/table_schema/get [POST]
func GetTableSchema(c *gin.Context) {
	var req dataset.GetTableSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.GetTableSchema(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ValidateTableSchema 验证表结构
// @router /api/knowledge/table_schema/validate [POST]
func ValidateTableSchema(c *gin.Context) {
	var req dataset.ValidateTableSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.ValidateTableSchema(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteSlice 删除切片
// @router /api/knowledge/slice/delete [POST]
func DeleteSlice(c *gin.Context) {
	var req dataset.DeleteSliceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.DeleteSlice(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateSlice 创建切片
// @router /api/knowledge/slice/create [POST]
func CreateSlice(c *gin.Context) {
	var req dataset.CreateSliceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.CreateSlice(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateSlice 更新切片
// @router /api/knowledge/slice/update [POST]
func UpdateSlice(c *gin.Context) {
	var req dataset.UpdateSliceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.UpdateSlice(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListSlice 列出切片
// @router /api/knowledge/slice/list [POST]
func ListSlice(c *gin.Context) {
	var req dataset.ListSliceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.ListSlice(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateDocumentReview 创建文档审核
// @router /api/knowledge/review/create [POST]
func CreateDocumentReview(c *gin.Context) {
	var req dataset.CreateDocumentReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.CreateDocumentReview(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// MGetDocumentReview 批量获取文档审核
// @router /api/knowledge/review/mget [POST]
func MGetDocumentReview(c *gin.Context) {
	var req dataset.MGetDocumentReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.MGetDocumentReview(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SaveDocumentReview 保存文档审核
// @router /api/knowledge/review/save [POST]
func SaveDocumentReview(c *gin.Context) {
	var req dataset.SaveDocumentReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.SaveDocumentReview(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetIconForDataset 获取数据集图标
// @router /api/knowledge/icon/get [POST]
func GetIconForDataset(c *gin.Context) {
	var req dataset.GetIconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := upload.SVC.GetIconForDataset(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ExtractPhotoCaption 提取照片标题
// @router /api/knowledge/photo/extract_caption [POST]
func ExtractPhotoCaption(c *gin.Context) {
	var req dataset.ExtractPhotoCaptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := application.KnowledgeSVC.ExtractPhotoCaption(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDocumentTableInfo 获取文档表信息（Memory 相关，但放在 Knowledge 路由中）
// @router /api/memory/doc_table_info [GET]
func GetDocumentTableInfo(c *gin.Context) {
	var req dataset.GetDocumentTableInfoRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	resp, err := knowledge.KnowledgeSVC.GetDocumentTableInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetModeConfig 获取模式配置（Memory 相关，但放在 Knowledge 路由中）
// @router /api/memory/table_mode_config [GET]
func GetModeConfig(c *gin.Context) {
	var req dataset.GetModeConfigRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if req.BotID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "bot_id is zero",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.GetModeConfig(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

