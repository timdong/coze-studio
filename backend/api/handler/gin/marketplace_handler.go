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
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	product_public_api "github.com/coze-dev/coze-studio/backend/api/model/marketplace/product_public_api"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
	"github.com/coze-dev/coze-studio/backend/api/model/app/bot_common"
	"github.com/coze-dev/coze-studio/backend/api/model/marketplace/product_common"
	"github.com/coze-dev/coze-studio/backend/api/model/playground"
	appApplication "github.com/coze-dev/coze-studio/backend/application/app"
	"github.com/coze-dev/coze-studio/backend/application/modelmgr"
	"github.com/coze-dev/coze-studio/backend/application/plugin"
	"github.com/coze-dev/coze-studio/backend/application/search"
	"github.com/coze-dev/coze-studio/backend/application/singleagent"
	"github.com/coze-dev/coze-studio/backend/application/template"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// PublicGetProductList 获取产品列表
// @router /api/marketplace/product/list [GET]
func PublicGetProductList(c *gin.Context) {
	var req product_public_api.GetProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	var resp *product_public_api.GetProductListResponse
	var err error

	switch req.GetEntityType() {
	case product_common.ProductEntityType_Plugin:
		// 检查插件服务是否已初始化
		if plugin.PluginApplicationSVC == nil || plugin.PluginApplicationSVC.DomainSVC == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    50000,
				"message": "Plugin service not initialized",
				"success": false,
				"data":    nil,
			})
			return
		}
		resp, err = plugin.PluginApplicationSVC.PublicGetProductList(ctx, &req)
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}

	case product_common.ProductEntityType_TemplateCommon:
		resp, err = template.ApplicationSVC.PublicGetProductList(ctx, &req)
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
	case product_common.ProductEntityType_SaasPlugin:
		// 检查插件服务是否已初始化
		if plugin.PluginApplicationSVC == nil || plugin.PluginApplicationSVC.DomainSVC == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    50000,
				"message": "Plugin service not initialized",
				"success": false,
				"data":    nil,
			})
			return
		}
		resp, err = plugin.PluginApplicationSVC.GetCozeSaasPluginList(ctx, &req)
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
	default:
		// 如果 entity_type 不匹配任何已知类型，返回空列表
		resp = &product_public_api.GetProductListResponse{
			Data: &product_public_api.GetProductListData{
				Products: []*product_public_api.ProductInfo{},
				HasMore:  false,
				Total:    0,
			},
		}
	}

	if resp == nil {
		// 防止空指针解引用
		resp = &product_public_api.GetProductListResponse{
			Data: &product_public_api.GetProductListData{
				Products: []*product_public_api.ProductInfo{},
				HasMore:  false,
				Total:    0,
			},
		}
	}

	c.JSON(http.StatusOK, resp)
}

// PublicGetProductDetail 获取产品详情
// @router /api/marketplace/product/detail [GET]
func PublicGetProductDetail(c *gin.Context) {
	var req product_public_api.GetProductDetailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.GetProductID() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "productID is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.PublicGetProductDetail(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PublicFavoriteProduct 收藏产品
// @router /api/marketplace/product/favorite [POST]
func PublicFavoriteProduct(c *gin.Context) {
	var req product_public_api.FavoriteProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.GetEntityID() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "entityID is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	// check entity id is valid
	if req.GetEntityType() == product_common.ProductEntityType_Bot {
		_, err := singleagent.SingleAgentSVC.ValidateAgentDraftAccess(ctx, req.GetEntityID())
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
	} else if req.GetEntityType() == product_common.ProductEntityType_Project {
		_, err := appApplication.APPApplicationSVC.ValidateDraftAPPAccess(ctx, req.GetEntityID())
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
	}

	resp, err := search.SearchSVC.PublicFavoriteProduct(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PublicGetUserFavoriteListV2 获取用户收藏列表 V2
// @router /api/marketplace/product/favorite/list.v2 [GET]
func PublicGetUserFavoriteListV2(c *gin.Context) {
	var req product_public_api.GetUserFavoriteListV2Request
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if req.GetPageSize() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "pageSize is invalid",
		})
		return
	}
	if req.GetEntityType() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "entityType is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := search.SearchSVC.PublicGetUserFavoriteList(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PublicDuplicateProduct 复制产品
// @router /api/marketplace/product/duplicate [POST]
func PublicDuplicateProduct(c *gin.Context) {
	var req product_public_api.DuplicateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp := new(product_public_api.DuplicateProductResponse)
	resp.Data = new(product_public_api.DuplicateProductData)

	switch req.GetEntityType() {
	case product_common.ProductEntityType_BotTemplate:
		modelListResp, err := modelmgr.ModelmgrApplicationSVC.GetModelList(ctx, &developer_api.GetTypeListRequest{})
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
		if modelListResp == nil || modelListResp.Data == nil || len(modelListResp.Data.ModelList) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  "no model found",
			})
			return
		}

		bot, err := singleagent.SingleAgentSVC.DuplicateDraftBot(ctx, &developer_api.DuplicateDraftBotRequest{
			BotID:   req.GetProductID(),
			SpaceID: req.GetSpaceID(),
		})
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}

		botInfo, err := singleagent.SingleAgentSVC.GetAgentBotInfo(ctx, &playground.GetDraftBotInfoAgwRequest{
			BotID: bot.Data.BotID,
		})
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
		if botInfo.Data == nil || botInfo.Data.BotInfo == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  "no bot info found",
			})
			return
		}

		modelInfo := botInfo.GetData().GetBotInfo().ModelInfo
		if modelInfo == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  "no model info found in agent",
			})
			return
		}
		modelInfo.ModelId = &modelListResp.Data.ModelList[0].ModelType

		if req.Name != nil {
			_, err = singleagent.SingleAgentSVC.UpdateSingleAgentDraft(ctx, &playground.UpdateDraftBotInfoAgwRequest{
				BotInfo: &bot_common.BotInfoForUpdate{
					BotId:     &bot.Data.BotID,
					Name:      req.Name,
					ModelInfo: modelInfo,
				},
			})
			if err != nil {
				httputil.InternalErrorGin(ctx, c, err)
				return
			}
		}

		resp.Data.NewEntityID = bot.Data.BotID

	case product_common.ProductEntityType_WorkflowTemplateV2:
		workflowResp, err := appworkflow.SVC.CopyWorkflow(ctx, &workflow.CopyWorkflowRequest{
			WorkflowID: strconv.FormatInt(req.GetProductID(), 10),
			SpaceID:    strconv.FormatInt(req.GetSpaceID(), 10),
		})
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}

		newWorkflowID, err := strconv.ParseInt(workflowResp.Data.WorkflowID, 10, 64)
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
		resp.Data.NewEntityID = newWorkflowID
		resp.Data.NewPluginID = &newWorkflowID

		if req.Name != nil {
			_, err = appworkflow.SVC.UpdateWorkflowMeta(ctx, &workflow.UpdateWorkflowMetaRequest{
				WorkflowID: workflowResp.Data.WorkflowID,
				SpaceID:    strconv.FormatInt(req.GetSpaceID(), 10),
				Name:       req.Name,
			})
			if err != nil {
				httputil.InternalErrorGin(ctx, c, err)
				return
			}
		}
	}

	c.JSON(http.StatusOK, resp)
}

// PublicSearchProduct 搜索产品
// @router /api/marketplace/product/search [GET]
func PublicSearchProduct(c *gin.Context) {
	var req product_public_api.SearchProductRequest

	// Handle category_ids query parameter
	var categoryIDs []int64
	if categoryIDsStr := c.Query("category_ids"); categoryIDsStr != "" {
		var err error
		categoryIDs, err = handlerCategoryIDs(categoryIDsStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  err.Error(),
			})
			return
		}
		// Remove category_ids from query to avoid binding conflict
		c.Request.URL.RawQuery = strings.ReplaceAll(c.Request.URL.RawQuery, "category_ids="+categoryIDsStr, "")
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	if len(categoryIDs) > 0 {
		req.CategoryIDs = categoryIDs
	}

	ctx := c.Request.Context()
	// Call plugin application service
	resp, err := plugin.PluginApplicationSVC.PublicSearchProduct(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "PublicSearchProduct failed: %v", err)
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func handlerCategoryIDs(categoryIDsStr string) ([]int64, error) {
	var categoryIDs []int64
	if categoryIDsStr != "" {
		categoryIDStrs := strings.Split(categoryIDsStr, ",")
		categoryIDs = make([]int64, 0, len(categoryIDStrs))
		for _, idStr := range categoryIDStrs {
			idStr = strings.TrimSpace(idStr)
			if idStr != "" {
				// Validate that it's a valid integer
				if categoryID, parseErr := strconv.ParseInt(idStr, 10, 64); parseErr == nil {
					categoryIDs = append(categoryIDs, categoryID)
				} else {
					return nil, fmt.Errorf("invalid category_id: %s", idStr)
				}
			}
		}
	}
	return categoryIDs, nil
}

// PublicSearchSuggest 搜索建议
// @router /api/marketplace/product/search/suggest [GET]
func PublicSearchSuggest(c *gin.Context) {
	var req product_public_api.SearchSuggestRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	// Call plugin application service
	resp, err := plugin.PluginApplicationSVC.PublicSearchSuggest(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "PublicSearchSuggest failed: %v", err)
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PublicGetProductCategoryList 获取产品分类列表
// @router /api/marketplace/product/category/list [GET]
func PublicGetProductCategoryList(c *gin.Context) {
	var req product_public_api.GetProductCategoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	var resp *product_public_api.GetProductCategoryListResponse
	req.EntityType = product_common.ProductEntityType_SaasPlugin
	switch req.GetEntityType() {
	case product_common.ProductEntityType_SaasPlugin:
		var err error
		resp, err = plugin.PluginApplicationSVC.GetSaasProductCategoryList(ctx, &req)
		if err != nil {
			httputil.InternalErrorGin(ctx, c, err)
			return
		}
	}

	c.JSON(http.StatusOK, resp)
}

// PublicGetProductCallInfo 获取产品调用信息
// @router /api/marketplace/product/call_info [GET]
func PublicGetProductCallInfo(c *gin.Context) {
	var req product_public_api.GetProductCallInfoRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetProductCallInfo(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PublicGetMarketPluginConfig 获取市场插件配置
// @router /api/marketplace/product/config [GET]
func PublicGetMarketPluginConfig(c *gin.Context) {
	var req product_public_api.GetMarketPluginConfigRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := plugin.PluginApplicationSVC.GetMarketPluginConfig(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

