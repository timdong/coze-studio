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

// TaskHandler MAS Task 处理器
type TaskHandler struct {
	taskService *service.TaskService
}

// NewTaskHandler 创建 Task 处理器
func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{
		taskService: service.NewTaskService(db),
	}
}

// RegisterRoutes 注册路由
func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	tasks := router.Group("/extensions/mas/tasks")
	{
		tasks.GET("", h.ListTasks)
		tasks.POST("", h.CreateTask)
		tasks.GET("/:id", h.GetTask)
		tasks.PUT("/:id", h.UpdateTask)
		tasks.DELETE("/:id", h.DeleteTask)
		tasks.PUT("/:id/status", h.UpdateTaskStatus)
		tasks.PUT("/:id/start", h.StartTask)
		tasks.PUT("/:id/complete", h.CompleteTask)
		tasks.PUT("/:id/fail", h.FailTask)
		tasks.GET("/pending", h.GetPendingTasks)
	}
}

// ListTasks 获取 Task 列表
func (h *TaskHandler) ListTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	// 构建过滤条件
	filters := make(map[string]interface{})
	if agentID, _ := strconv.ParseUint(c.Query("agent_id"), 10, 32); agentID > 0 {
		filters["agent_id"] = uint(agentID)
	}
	if sessionID, _ := strconv.ParseUint(c.Query("session_id"), 10, 32); sessionID > 0 {
		filters["session_id"] = uint(sessionID)
	}
	if workspaceID, _ := strconv.ParseUint(c.Query("workspace_id"), 10, 32); workspaceID > 0 {
		filters["workspace_id"] = uint(workspaceID)
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if taskType := c.Query("type"); taskType != "" {
		filters["type"] = taskType
	}
	
	tasks, total, err := h.taskService.ListTasks(c.Request.Context(), filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "获取 Task 列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取 Task 列表成功",
		"data": gin.H{
			"tasks": tasks,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateTask 创建 Task
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.MASTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	if err := h.taskService.CreateTask(c.Request.Context(), &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "创建 Task 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "创建 Task 成功",
		"data":    task,
	})
}

// GetTask 获取 Task 详情
func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Task ID",
			"data":    nil,
		})
		return
	}
	
	task, err := h.taskService.GetTask(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    40400,
			"success": false,
			"message": "Task 不存在: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取 Task 成功",
		"data":    task,
	})
}

// UpdateTask 更新 Task
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Task ID",
			"data":    nil,
		})
		return
	}
	
	var task models.MASTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	task.ID = uint(id)
	if err := h.taskService.UpdateTask(c.Request.Context(), &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "更新 Task 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "更新 Task 成功",
		"data":    task,
	})
}

// DeleteTask 删除 Task
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Task ID",
			"data":    nil,
		})
		return
	}
	
	if err := h.taskService.DeleteTask(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "删除 Task 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "删除 Task 成功",
		"data":    nil,
	})
}

// UpdateTaskStatus 更新 Task 状态
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Task ID",
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
	
	if err := h.taskService.UpdateTaskStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "更新 Task 状态失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "更新 Task 状态成功",
		"data":    nil,
	})
}

// StartTask 开始执行 Task
func (h *TaskHandler) StartTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Task ID",
			"data":    nil,
		})
		return
	}
	
	if err := h.taskService.StartTask(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "启动 Task 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "启动 Task 成功",
		"data":    nil,
	})
}

// CompleteTask 完成 Task
func (h *TaskHandler) CompleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Task ID",
			"data":    nil,
		})
		return
	}
	
	var req struct {
		Output interface{} `json:"output"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// output 是可选的，忽略错误
	}
	
	if err := h.taskService.CompleteTask(c.Request.Context(), uint(id), req.Output); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "完成 Task 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "完成 Task 成功",
		"data":    nil,
	})
}

// FailTask 标记 Task 为失败
func (h *TaskHandler) FailTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"success": false,
			"message": "无效的 Task ID",
			"data":    nil,
		})
		return
	}
	
	var req struct {
		Error string `json:"error" binding:"required"`
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
	
	if err := h.taskService.FailTask(c.Request.Context(), uint(id), req.Error); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "标记 Task 失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "标记 Task 失败成功",
		"data":    nil,
	})
}

// GetPendingTasks 获取待处理的任务
func (h *TaskHandler) GetPendingTasks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	
	tasks, err := h.taskService.GetPendingTasks(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"success": false,
			"message": "获取待处理任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"success": true,
		"message": "获取待处理任务成功",
		"data": gin.H{
			"tasks": tasks,
		},
	})
}

