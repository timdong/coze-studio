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

	"github.com/coze-dev/coze-studio/backend/extensions/etrxtable/models"
	"github.com/coze-dev/coze-studio/backend/extensions/etrxtable/service"
)

// TableHandler 表格处理器
type TableHandler struct {
	tableService *service.TableService
}

// NewTableHandler 创建表格处理器
func NewTableHandler(db *gorm.DB) *TableHandler {
	return &TableHandler{
		tableService: service.NewTableService(db),
	}
}

// RegisterRoutes 注册路由
func (h *TableHandler) RegisterRoutes(router *gin.RouterGroup) {
	tables := router.Group("/extensions/etrxtable/tables")
	{
		tables.GET("", h.ListTables)
		tables.POST("", h.CreateTable)
		tables.GET("/:id", h.GetTable)
		tables.PUT("/:id", h.UpdateTable)
		tables.DELETE("/:id", h.DeleteTable)
		
		// 列管理
		tables.POST("/:id/columns", h.CreateColumn)
		tables.GET("/:id/columns", h.GetColumns)
		tables.PUT("/:id/columns/:columnId", h.UpdateColumn)
		tables.DELETE("/:id/columns/:columnId", h.DeleteColumn)
		
		// 行管理
		tables.POST("/:id/rows", h.CreateRow)
		tables.GET("/:id/rows", h.GetRows)
		tables.PUT("/:id/rows/:rowId", h.UpdateRow)
		tables.DELETE("/:id/rows/:rowId", h.DeleteRow)
		tables.POST("/:id/rows/batch", h.BatchCreateRows)
		tables.PUT("/:id/rows/batch", h.BatchUpdateRows)
		tables.DELETE("/:id/rows/batch", h.BatchDeleteRows)
	}
}

// ListTables 获取表格列表
// @Summary 获取表格列表
// @Tags 多维表格
// @Param workspace_id query int false "工作空间ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/extensions/etrxtable/tables [get]
func (h *TableHandler) ListTables(c *gin.Context) {
	workspaceID, _ := strconv.ParseUint(c.Query("workspace_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	tables, total, err := h.tableService.ListTables(c.Request.Context(), uint(workspaceID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "获取表格列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取表格列表成功",
		"data": gin.H{
			"tables": tables,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateTable 创建表格
// @Summary 创建表格
// @Tags 多维表格
// @Param table body models.TableBody true "表格信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/extensions/etrxtable/tables [post]
func (h *TableHandler) CreateTable(c *gin.Context) {
	var table models.TableBody
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	if err := h.tableService.CreateTable(c.Request.Context(), &table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "创建表格失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "创建表格成功",
		"data":    table,
	})
}

// GetTable 获取表格详情
// @Summary 获取表格详情
// @Tags 多维表格
// @Param id path int true "表格ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/extensions/etrxtable/tables/{id} [get]
func (h *TableHandler) GetTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的表格ID",
			"data":    nil,
		})
		return
	}
	
	table, err := h.tableService.GetTable(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    40400,
			"success": false,
			"message": "表格不存在: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取表格成功",
		"data":    table,
	})
}

// UpdateTable 更新表格
// @Summary 更新表格
// @Tags 多维表格
// @Param id path int true "表格ID"
// @Param table body models.TableBody true "表格信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/extensions/etrxtable/tables/{id} [put]
func (h *TableHandler) UpdateTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的表格ID",
			"data":    nil,
		})
		return
	}
	
	var table models.TableBody
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	table.ID = uint(id)
	if err := h.tableService.UpdateTable(c.Request.Context(), &table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "更新表格失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "更新表格成功",
		"data":    table,
	})
}

// DeleteTable 删除表格
// @Summary 删除表格
// @Tags 多维表格
// @Param id path int true "表格ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/extensions/etrxtable/tables/{id} [delete]
func (h *TableHandler) DeleteTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的表格ID",
			"data":    nil,
		})
		return
	}
	
	if err := h.tableService.DeleteTable(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "删除表格失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "删除表格成功",
		"data":    nil,
	})
}

// CreateColumn 创建列
func (h *TableHandler) CreateColumn(c *gin.Context) {
	tableID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var column models.TableColumn
	if err := c.ShouldBindJSON(&column); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	column.TableBodyID = uint(tableID)
	if err := h.tableService.CreateColumn(c.Request.Context(), &column); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "创建列失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "创建列成功",
		"data":    column,
	})
}

// GetColumns 获取列列表
func (h *TableHandler) GetColumns(c *gin.Context) {
	tableID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	columns, err := h.tableService.GetColumns(c.Request.Context(), uint(tableID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "获取列列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取列列表成功",
		"data":    columns,
	})
}

// UpdateColumn 更新列
func (h *TableHandler) UpdateColumn(c *gin.Context) {
	columnID, _ := strconv.ParseUint(c.Param("columnId"), 10, 32)
	
	var column models.TableColumn
	if err := c.ShouldBindJSON(&column); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	column.ID = uint(columnID)
	if err := h.tableService.UpdateColumn(c.Request.Context(), &column); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "更新列失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "更新列成功",
		"data":    column,
	})
}

// DeleteColumn 删除列
func (h *TableHandler) DeleteColumn(c *gin.Context) {
	columnID, _ := strconv.ParseUint(c.Param("columnId"), 10, 32)
	
	if err := h.tableService.DeleteColumn(c.Request.Context(), uint(columnID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "删除列失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "删除列成功",
		"data":    nil,
	})
}

// CreateRow 创建行
func (h *TableHandler) CreateRow(c *gin.Context) {
	tableID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var row models.TableRow
	if err := c.ShouldBindJSON(&row); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	row.TableBodyID = uint(tableID)
	if err := h.tableService.CreateRow(c.Request.Context(), &row); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "创建行失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "创建行成功",
		"data":    row,
	})
}

// GetRows 获取行列表
func (h *TableHandler) GetRows(c *gin.Context) {
	tableID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	rows, total, err := h.tableService.GetRows(c.Request.Context(), uint(tableID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "获取行列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取行列表成功",
		"data": gin.H{
			"rows": rows,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// UpdateRow 更新行
func (h *TableHandler) UpdateRow(c *gin.Context) {
	rowID, _ := strconv.ParseUint(c.Param("rowId"), 10, 32)
	
	var row models.TableRow
	if err := c.ShouldBindJSON(&row); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	row.ID = uint(rowID)
	if err := h.tableService.UpdateRow(c.Request.Context(), &row); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "更新行失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "更新行成功",
		"data":    row,
	})
}

// DeleteRow 删除行
func (h *TableHandler) DeleteRow(c *gin.Context) {
	rowID, _ := strconv.ParseUint(c.Param("rowId"), 10, 32)
	
	if err := h.tableService.DeleteRow(c.Request.Context(), uint(rowID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "删除行失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "删除行成功",
		"data":    nil,
	})
}

// BatchCreateRows 批量创建行
func (h *TableHandler) BatchCreateRows(c *gin.Context) {
	var rows []models.TableRow
	if err := c.ShouldBindJSON(&rows); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	if err := h.tableService.BatchCreateRows(c.Request.Context(), rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "批量创建行失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "批量创建行成功",
		"data":    rows,
	})
}

// BatchUpdateRows 批量更新行
func (h *TableHandler) BatchUpdateRows(c *gin.Context) {
	var rows []models.TableRow
	if err := c.ShouldBindJSON(&rows); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	if err := h.tableService.BatchUpdateRows(c.Request.Context(), rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "批量更新行失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "批量更新行成功",
		"data":    rows,
	})
}

// BatchDeleteRows 批量删除行
func (h *TableHandler) BatchDeleteRows(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
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
	
	if err := h.tableService.BatchDeleteRows(c.Request.Context(), req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "批量删除行失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "批量删除行成功",
		"data":    nil,
	})
}

