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
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/admin/config"
	bizConf "github.com/coze-dev/coze-studio/backend/bizpkg/config"
	"github.com/coze-dev/coze-studio/backend/bizpkg/config/modelmgr"
	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/coze-studio/backend/bizpkg/llm/modelbuilder"
	"github.com/coze-dev/coze-studio/backend/infra/embedding/impl"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/conv"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// GetBasicConfiguration 获取基础配置
// @router /api/admin/config/basic/get [GET]
func GetBasicConfiguration(c *gin.Context) {
	ctx := c.Request.Context()
	baseConfig, err := bizConf.Base().GetBaseConfig(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(config.GetBasicConfigurationResp)
	resp.Configuration = baseConfig

	c.JSON(http.StatusOK, resp)
}

// SaveBasicConfiguration 保存基础配置
// @router /api/admin/config/basic/save [POST]
func SaveBasicConfiguration(c *gin.Context) {
	var req config.SaveBasicConfigurationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if req.Configuration == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "Configuration is nil",
		})
		return
	}

	// TODO: check coze api token

	// Validate ServerHost: allow http/https URLs, or hostname:port
	if req.Configuration.ServerHost == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "ServerHost is empty",
		})
		return
	}

	host := req.Configuration.ServerHost
	if strings.Contains(host, "://") {
		u, parseErr := url.Parse(host)
		if parseErr != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  "ServerHost is invalid URL, require http/https",
			})
			return
		}
	} else {
		// Expect hostname:port format
		h, p, splitErr := net.SplitHostPort(host)
		if splitErr != nil || h == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  "ServerHost must be hostname:port or http(s) URL",
			})
			return
		}
		port, portErr := strconv.Atoi(p)
		if portErr != nil || port <= 0 || port > 65535 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  "ServerHost port is invalid",
			})
			return
		}
	}

	logs.Infof("server host is valid %s", req.Configuration.ServerHost)

	if req.Configuration.CodeRunnerType.String() == "<UNSET>" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "CodeRunnerType is invalid",
		})
		return
	}

	ctx := c.Request.Context()
	err := bizConf.Base().SaveBaseConfig(ctx, req.Configuration)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, fmt.Errorf("save basic config failed: %w", err))
		return
	}

	resp := new(config.SaveBasicConfigurationResp)
	c.JSON(http.StatusOK, resp)
}

// GetKnowledgeConfig 获取知识库配置
// @router /api/admin/config/knowledge/get [GET]
func GetKnowledgeConfig(c *gin.Context) {
	ctx := c.Request.Context()
	knowledgeConfig, err := bizConf.Knowledge().GetKnowledgeConfig(ctx)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, fmt.Errorf("get knowledge config failed: %w", err))
		return
	}

	resp := new(config.GetKnowledgeConfigResp)
	resp.KnowledgeConfig = knowledgeConfig

	c.JSON(http.StatusOK, resp)
}

// UpdateKnowledgeConfig 更新知识库配置
// @router /api/admin/config/knowledge/save [POST]
func UpdateKnowledgeConfig(c *gin.Context) {
	var req config.UpdateKnowledgeConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if req.KnowledgeConfig == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "KnowledgeConfig is nil",
		})
		return
	}

	if req.KnowledgeConfig.EmbeddingConfig == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "EmbeddingConfig is nil",
		})
		return
	}

	if req.KnowledgeConfig.EmbeddingConfig.Connection == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "Connection is nil",
		})
		return
	}

	if req.KnowledgeConfig.EmbeddingConfig.Connection.EmbeddingInfo == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  "EmbeddingInfo is nil",
		})
		return
	}

	ctx := c.Request.Context()
	embedding, err := impl.GetEmbedding(ctx, req.KnowledgeConfig.EmbeddingConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  fmt.Sprintf("get embedding failed: %v", err),
		})
		return
	}

	if req.KnowledgeConfig.EmbeddingConfig.Connection.EmbeddingInfo.Dims <= 0 {
		req.KnowledgeConfig.EmbeddingConfig.Connection.EmbeddingInfo.Dims = int32(embedding.Dimensions())

		embedding, err = impl.GetEmbedding(ctx, req.KnowledgeConfig.EmbeddingConfig)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 40000,
				"msg":  fmt.Sprintf("get embedding failed: %v", err),
			})
			return
		}
	}

	denseEmbeddings, err := embedding.EmbedStrings(ctx, []string{"test"})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  fmt.Sprintf("embed test string failed: %v", err),
		})
		return
	}

	if len(denseEmbeddings) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  fmt.Sprintf("embed test string failed: %v", err),
		})
		return
	}

	logs.CtxDebugf(ctx, "embed test string result: %d, expect %d",
		len(denseEmbeddings[0]), req.KnowledgeConfig.EmbeddingConfig.Connection.EmbeddingInfo.Dims)
	if len(denseEmbeddings[0]) != int(req.KnowledgeConfig.EmbeddingConfig.Connection.EmbeddingInfo.Dims) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg": fmt.Sprintf("embed test string failed: dims not match, expect %d, got %d",
				req.KnowledgeConfig.EmbeddingConfig.Connection.EmbeddingInfo.Dims, len(denseEmbeddings[0])),
		})
		return
	}

	err = bizConf.Knowledge().SaveKnowledgeConfig(ctx, req.KnowledgeConfig)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, fmt.Errorf("save knowledge config failed: %w", err))
		return
	}

	resp := new(config.UpdateKnowledgeConfigResp)
	c.JSON(http.StatusOK, resp)
}

// GetModelList 获取模型列表
// @router /api/admin/config/model/list [GET]
func GetModelList(c *gin.Context) {
	var req config.GetModelListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	modelList, err := bizConf.ModelConf().GetProviderModelList(ctx)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, fmt.Errorf("get model list failed: %w", err))
		return
	}

	resp := new(config.GetModelListResp)
	resp.ProviderModelList = modelList

	c.JSON(http.StatusOK, resp)
}

// CreateModel 创建模型
// @router /api/admin/config/model/create [POST]
func CreateModel(c *gin.Context) {
	var req config.CreateModelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	modelBuilder, err := modelbuilder.NewModelBuilder(req.ModelClass, &config.Model{
		EnableBase64URL: req.EnableBase64URL,
		Connection:      req.Connection,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  fmt.Sprintf("create model builder failed: %v", err),
		})
		return
	}

	ctx := c.Request.Context()
	logs.CtxDebugf(ctx, "create model req: %s, conn: %s", conv.DebugJsonToStr(req), conv.DebugJsonToStr(req.Connection.BaseConnInfo))

	chatModel, err := modelBuilder.Build(ctx, &modelbuilder.LLMParams{EnableThinking: ptr.Of(false)})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  fmt.Sprintf("build model failed: %v", err),
		})
		return
	}

	respMsgs, err := chatModel.Generate(ctx, []*schema.Message{
		schema.SystemMessage("1+1=?,Just answer with a number, no explanation.")})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  fmt.Sprintf("generate model failed: %v", err),
		})
		return
	}

	logs.CtxDebugf(ctx, "chatModel.Generate resp : %s", conv.DebugJsonToStr(respMsgs))

	id, err := bizConf.ModelConf().CreateModel(ctx, req.ModelClass, req.ModelName, req.Connection, &modelmgr.ModelExtra{
		EnableBase64URL: req.EnableBase64URL,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  fmt.Sprintf("create model failed: %v", err),
		})
		return
	}

	resp := new(config.CreateModelResp)
	resp.ID = id

	c.JSON(http.StatusOK, resp)
}

// DeleteModel 删除模型
// @router /api/admin/config/model/delete [POST]
func DeleteModel(c *gin.Context) {
	var req config.DeleteModelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	err := bizConf.ModelConf().DeleteModel(ctx, req.ID)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, fmt.Errorf("delete model failed: %w", err))
		return
	}

	resp := new(config.DeleteModelResp)
	c.JSON(http.StatusOK, resp)
}

