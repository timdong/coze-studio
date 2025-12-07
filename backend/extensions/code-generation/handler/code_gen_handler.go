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

package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/coze-dev/coze-studio/backend/extensions/code-generation/service"
)

// CodeGenHandler 代码生成处理器
type CodeGenHandler struct {
	codeGenService *service.CodeGenService
}

// NewCodeGenHandler 创建代码生成处理器
func NewCodeGenHandler(db *gorm.DB) *CodeGenHandler {
	return &CodeGenHandler{
		codeGenService: service.NewCodeGenService(db),
	}
}

// RegisterRoutes 注册路由
func (h *CodeGenHandler) RegisterRoutes(router *gin.RouterGroup) {
	codegen := router.Group("/extensions/code-generation")
	{
		codegen.POST("/generate", h.GenerateCode)
	}
}

// GenerateCode 生成代码
func (h *CodeGenHandler) GenerateCode(c *gin.Context) {
	var req service.GenerateCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	// 验证语言
	validLanguages := map[string]bool{
		"go":     true,
		"python": true,
		"java":   true,
	}
	if !validLanguages[req.Language] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "不支持的语言: " + req.Language,
			"data":    nil,
		})
		return
	}
	
	// 生成代码
	result, err := h.codeGenService.GenerateCode(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "生成代码失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "生成代码成功",
		"data":    result,
	})
}

