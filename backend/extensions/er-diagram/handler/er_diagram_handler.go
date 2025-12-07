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
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/coze-dev/coze-studio/backend/extensions/er-diagram/models"
	"github.com/coze-dev/coze-studio/backend/extensions/er-diagram/service"
)

// ERDiagramHandler ER 图处理器
type ERDiagramHandler struct {
	diagramService *service.ERDiagramService
}

// NewERDiagramHandler 创建 ER 图处理器
func NewERDiagramHandler(db *gorm.DB) *ERDiagramHandler {
	return &ERDiagramHandler{
		diagramService: service.NewERDiagramService(db),
	}
}

// RegisterRoutes 注册路由
func (h *ERDiagramHandler) RegisterRoutes(router *gin.RouterGroup) {
	diagrams := router.Group("/extensions/er-diagram/diagrams")
	{
		diagrams.GET("", h.ListDiagrams)
		diagrams.POST("", h.CreateDiagram)
		diagrams.GET("/:id", h.GetDiagram)
		diagrams.PUT("/:id", h.UpdateDiagram)
		diagrams.DELETE("/:id", h.DeleteDiagram)
		diagrams.GET("/:id/export/sql", h.ExportToSQL)
		diagrams.GET("/:id/export/dbml", h.ExportToDBML)
		diagrams.POST("/import/sql", h.ImportFromSQL)
		diagrams.POST("/import/dbml", h.ImportFromDBML)
	}
}

// ListDiagrams 获取 ER 图列表
func (h *ERDiagramHandler) ListDiagrams(c *gin.Context) {
	workspaceID, _ := strconv.ParseUint(c.Query("workspace_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	// 构建过滤条件
	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if dbType := c.Query("database_type"); dbType != "" {
		filters["database_type"] = dbType
	}
	
	diagrams, total, err := h.diagramService.ListDiagrams(c.Request.Context(), uint(workspaceID), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "获取 ER 图列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取 ER 图列表成功",
		"data": gin.H{
			"diagrams": diagrams,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateDiagram 创建 ER 图
func (h *ERDiagramHandler) CreateDiagram(c *gin.Context) {
	var diagram models.ERDiagram
	if err := c.ShouldBindJSON(&diagram); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	if err := h.diagramService.CreateDiagram(c.Request.Context(), &diagram); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "创建 ER 图失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "创建 ER 图成功",
		"data":    diagram,
	})
}

// GetDiagram 获取 ER 图详情
func (h *ERDiagramHandler) GetDiagram(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 ER 图 ID",
			"data":    nil,
		})
		return
	}
	
	diagram, err := h.diagramService.GetDiagram(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    40400,
			"success": false,
			"message": "ER 图不存在: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取 ER 图成功",
		"data":    diagram,
	})
}

// UpdateDiagram 更新 ER 图
func (h *ERDiagramHandler) UpdateDiagram(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 ER 图 ID",
			"data":    nil,
		})
		return
	}
	
	var diagram models.ERDiagram
	if err := c.ShouldBindJSON(&diagram); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	diagram.ID = uint(id)
	if err := h.diagramService.UpdateDiagram(c.Request.Context(), &diagram); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "更新 ER 图失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "更新 ER 图成功",
		"data":    diagram,
	})
}

// DeleteDiagram 删除 ER 图
func (h *ERDiagramHandler) DeleteDiagram(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 ER 图 ID",
			"data":    nil,
		})
		return
	}
	
	if err := h.diagramService.DeleteDiagram(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "删除 ER 图失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "删除 ER 图成功",
		"data":    nil,
	})
}

// ExportToSQL 导出为 SQL
func (h *ERDiagramHandler) ExportToSQL(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 ER 图 ID",
			"data":    nil,
		})
		return
	}
	
	sql, err := h.diagramService.ExportToSQL(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "导出 SQL 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "导出 SQL 成功",
		"data": gin.H{
			"sql": sql,
		},
	})
}

// ExportToDBML 导出为 DBML
func (h *ERDiagramHandler) ExportToDBML(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 ER 图 ID",
			"data":    nil,
		})
		return
	}
	
	dbml, err := h.diagramService.ExportToDBML(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "导出 DBML 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "导出 DBML 成功",
		"data": gin.H{
			"dbml": dbml,
		},
	})
}

// ImportFromSQL 从 SQL 导入
func (h *ERDiagramHandler) ImportFromSQL(c *gin.Context) {
	var req struct {
		SQL         string `json:"sql" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		WorkspaceID uint   `json:"workspace_id" binding:"required"`
		CreatedBy   uint   `json:"created_by" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	diagram, err := h.diagramService.ImportFromSQL(c.Request.Context(), req.SQL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "导入 SQL 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	// 设置导入的元数据
	diagram.Name = req.Name
	diagram.Description = req.Description
	diagram.WorkspaceID = req.WorkspaceID
	diagram.CreatedBy = req.CreatedBy
	
	if err := h.diagramService.CreateDiagram(c.Request.Context(), diagram); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "保存导入的 ER 图失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "导入 SQL 成功",
		"data":    diagram,
	})
}

// ImportFromDBML 从 DBML 导入
func (h *ERDiagramHandler) ImportFromDBML(c *gin.Context) {
	var req struct {
		DBML        string `json:"dbml" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		WorkspaceID uint   `json:"workspace_id" binding:"required"`
		CreatedBy   uint   `json:"created_by" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	diagram, err := h.diagramService.ImportFromDBML(c.Request.Context(), req.DBML)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "导入 DBML 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	// 设置导入的元数据
	diagram.Name = req.Name
	diagram.Description = req.Description
	diagram.WorkspaceID = req.WorkspaceID
	diagram.CreatedBy = req.CreatedBy
	
	if err := h.diagramService.CreateDiagram(c.Request.Context(), diagram); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "保存导入的 ER 图失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "导入 DBML 成功",
		"data":    diagram,
	})
}

