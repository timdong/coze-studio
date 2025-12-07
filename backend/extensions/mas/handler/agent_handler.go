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

// AgentHandler Agent 处理器
type AgentHandler struct {
	agentService *service.AgentService
}

// NewAgentHandler 创建 Agent 处理器
func NewAgentHandler(db *gorm.DB) *AgentHandler {
	return &AgentHandler{
		agentService: service.NewAgentService(db),
	}
}

// RegisterRoutes 注册路由
func (h *AgentHandler) RegisterRoutes(router *gin.RouterGroup) {
	agents := router.Group("/extensions/mas/agents")
	{
		agents.GET("", h.ListAgents)
		agents.POST("", h.CreateAgent)
		agents.GET("/:id", h.GetAgent)
		agents.PUT("/:id", h.UpdateAgent)
		agents.DELETE("/:id", h.DeleteAgent)
		agents.PUT("/:id/status", h.UpdateAgentStatus)
		agents.PUT("/:id/metrics", h.UpdateAgentMetrics)
	}
}

// ListAgents 获取 Agent 列表
func (h *AgentHandler) ListAgents(c *gin.Context) {
	workspaceID, _ := strconv.ParseUint(c.Query("workspace_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	// 构建过滤条件
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if agentType := c.Query("type"); agentType != "" {
		filters["type"] = agentType
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	
	agents, total, err := h.agentService.ListAgents(c.Request.Context(), uint(workspaceID), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "获取 Agent 列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取 Agent 列表成功",
		"data": gin.H{
			"agents": agents,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateAgent 创建 Agent
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	var agent models.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	if err := h.agentService.CreateAgent(c.Request.Context(), &agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "创建 Agent 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "创建 Agent 成功",
		"data":    agent,
	})
}

// GetAgent 获取 Agent 详情
func (h *AgentHandler) GetAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Agent ID",
			"data":    nil,
		})
		return
	}
	
	agent, err := h.agentService.GetAgent(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    40400,
			"success": false,
			"message": "Agent 不存在: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取 Agent 成功",
		"data":    agent,
	})
}

// UpdateAgent 更新 Agent
func (h *AgentHandler) UpdateAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Agent ID",
			"data":    nil,
		})
		return
	}
	
	var agent models.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	agent.ID = uint(id)
	if err := h.agentService.UpdateAgent(c.Request.Context(), &agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "更新 Agent 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "更新 Agent 成功",
		"data":    agent,
	})
}

// DeleteAgent 删除 Agent
func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Agent ID",
			"data":    nil,
		})
		return
	}
	
	if err := h.agentService.DeleteAgent(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "删除 Agent 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "删除 Agent 成功",
		"data":    nil,
	})
}

// UpdateAgentStatus 更新 Agent 状态
func (h *AgentHandler) UpdateAgentStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Agent ID",
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
	
	if err := h.agentService.UpdateAgentStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "更新 Agent 状态失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "更新 Agent 状态成功",
		"data":    nil,
	})
}

// UpdateAgentMetrics 更新 Agent 性能指标
func (h *AgentHandler) UpdateAgentMetrics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Agent ID",
			"data":    nil,
		})
		return
	}
	
	var metrics map[string]interface{}
	if err := c.ShouldBindJSON(&metrics); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	if err := h.agentService.UpdateAgentMetrics(c.Request.Context(), uint(id), metrics); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "更新 Agent 性能指标失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "更新 Agent 性能指标成功",
		"data":    nil,
	})
}

