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
	"github.com/coze-dev/coze-studio/backend/extensions/mas/models"
	"github.com/coze-dev/coze-studio/backend/extensions/mas/service"
)

// SessionHandler MAS Session 处理器
type SessionHandler struct {
	sessionService *service.SessionService
}

// NewSessionHandler 创建 Session 处理器
func NewSessionHandler(db *gorm.DB) *SessionHandler {
	return &SessionHandler{
		sessionService: service.NewSessionService(db),
	}
}

// RegisterRoutes 注册路由
func (h *SessionHandler) RegisterRoutes(router *gin.RouterGroup) {
	sessions := router.Group("/extensions/mas/sessions")
	{
		sessions.GET("", h.ListSessions)
		sessions.POST("", h.CreateSession)
		sessions.GET("/:id", h.GetSession)
		sessions.PUT("/:id", h.UpdateSession)
		sessions.DELETE("/:id", h.DeleteSession)
		sessions.PUT("/:id/status", h.UpdateSessionStatus)
		sessions.PUT("/:id/complete", h.CompleteSession)
	}
}

// ListSessions 获取 Session 列表
func (h *SessionHandler) ListSessions(c *gin.Context) {
	workspaceID, _ := strconv.ParseUint(c.Query("workspace_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	// 构建过滤条件
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if sessionType := c.Query("type"); sessionType != "" {
		filters["type"] = sessionType
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	
	sessions, total, err := h.sessionService.ListSessions(c.Request.Context(), uint(workspaceID), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "获取 Session 列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取 Session 列表成功",
		"data": gin.H{
			"sessions": sessions,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateSession 创建 Session
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var session models.MASSession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	if err := h.sessionService.CreateSession(c.Request.Context(), &session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "创建 Session 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "创建 Session 成功",
		"data":    session,
	})
}

// GetSession 获取 Session 详情
func (h *SessionHandler) GetSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Session ID",
			"data":    nil,
		})
		return
	}
	
	session, err := h.sessionService.GetSession(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    40400,
			"success": false,
			"message": "Session 不存在: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取 Session 成功",
		"data":    session,
	})
}

// UpdateSession 更新 Session
func (h *SessionHandler) UpdateSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Session ID",
			"data":    nil,
		})
		return
	}
	
	var session models.MASSession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	session.ID = uint(id)
	if err := h.sessionService.UpdateSession(c.Request.Context(), &session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "更新 Session 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "更新 Session 成功",
		"data":    session,
	})
}

// DeleteSession 删除 Session
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Session ID",
			"data":    nil,
		})
		return
	}
	
	if err := h.sessionService.DeleteSession(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "删除 Session 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "删除 Session 成功",
		"data":    nil,
	})
}

// UpdateSessionStatus 更新 Session 状态
func (h *SessionHandler) UpdateSessionStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Session ID",
			"data":    nil,
		})
		return
	}
	
	var req struct {
		Status string `json:"status" binding:"required"`
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
	
	if err := h.sessionService.UpdateSessionStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "更新 Session 状态失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "更新 Session 状态成功",
		"data":    nil,
	})
}

// CompleteSession 完成 Session
func (h *SessionHandler) CompleteSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Session ID",
			"data":    nil,
		})
		return
	}
	
	if err := h.sessionService.CompleteSession(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "完成 Session 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "完成 Session 成功",
		"data":    nil,
	})
}

