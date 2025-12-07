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

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/api/model/data/database/table"
	"github.com/coze-dev/coze-studio/backend/application/memory"
	"github.com/coze-dev/coze-studio/backend/application/singleagent"
)

// ListDatabase 列出数据库
// @router /api/memory/database/list [POST]
func ListDatabase(c *gin.Context) {
	var req table.ListDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.ListDatabase(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDatabaseByID 根据 ID 获取数据库
// @router /api/memory/database/get_by_id [POST]
func GetDatabaseByID(c *gin.Context) {
	var req table.SingleDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.GetDatabaseByID(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AddDatabase 添加数据库
// @router /api/memory/database/add [POST]
func AddDatabase(c *gin.Context) {
	var req table.AddDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.AddDatabase(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateDatabase 更新数据库
// @router /api/memory/database/update [POST]
func UpdateDatabase(c *gin.Context) {
	var req table.UpdateDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.UpdateDatabase(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteDatabase 删除数据库
// @router /api/memory/database/delete [POST]
func DeleteDatabase(c *gin.Context) {
	var req table.DeleteDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.DeleteDatabase(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// BindDatabase 绑定数据库到 Bot
// @router /api/memory/database/bind_to_bot [POST]
func BindDatabase(c *gin.Context) {
	var req table.BindDatabaseToBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.BindDatabase(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UnBindDatabase 解绑数据库与 Bot
// @router /api/memory/database/unbind_to_bot [POST]
func UnBindDatabase(c *gin.Context) {
	var req table.BindDatabaseToBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.UnBindDatabase(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListDatabaseRecords 列出数据库记录
// @router /api/memory/database/list_records [POST]
func ListDatabaseRecords(c *gin.Context) {
	var req table.ListDatabaseRecordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.ListDatabaseRecords(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateDatabaseRecords 更新数据库记录
// @router /api/memory/database/update_records [POST]
func UpdateDatabaseRecords(c *gin.Context) {
	var req table.UpdateDatabaseRecordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.UpdateDatabaseRecords(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetOnlineDatabaseId 获取在线数据库 ID
// @router /api/memory/database/get_online_database_id [POST]
func GetOnlineDatabaseId(c *gin.Context) {
	var req table.GetOnlineDatabaseIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.GetOnlineDatabaseId(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ResetBotTable 重置 Bot 表
// @router /api/memory/database/table/reset [POST]
func ResetBotTable(c *gin.Context) {
	var req table.ResetBotTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.ResetBotTable(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDatabaseTemplate 获取数据库模板
// @router /api/memory/database/get_template [POST]
func GetDatabaseTemplate(c *gin.Context) {
	var req table.GetDatabaseTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.GetDatabaseTemplate(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetConnectorName 获取连接器名称
// @router /api/memory/database/get_connector_name [POST]
func GetConnectorName(c *gin.Context) {
	var req table.GetSpaceConnectorListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.GetConnectorName(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetBotDatabase 获取 Bot 数据库
// @router /api/memory/database/table/list_new [POST]
func GetBotDatabase(c *gin.Context) {
	var req table.GetBotTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.GetBotDatabase(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateDatabaseBotSwitch 更新数据库 Bot 开关
// @router /api/memory/database/update_bot_switch [POST]
func UpdateDatabaseBotSwitch(c *gin.Context) {
	var req table.UpdateDatabaseBotSwitchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := singleagent.SingleAgentSVC.UpdatePromptDisable(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDatabaseTableSchema 获取数据库表结构
// @router /api/memory/table_schema/get [POST]
func GetDatabaseTableSchema(c *gin.Context) {
	var req table.GetTableSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.GetDatabaseTableSchema(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SubmitDatabaseInsertTask 提交数据库插入任务
// @router /api/memory/table_file/submit [POST]
func SubmitDatabaseInsertTask(c *gin.Context) {
	var req table.SubmitDatabaseInsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.SubmitDatabaseInsertTask(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DatabaseFileProgressData 获取数据库文件进度数据
// @router /api/memory/table_file/get_progress [POST]
func DatabaseFileProgressData(c *gin.Context) {
	var req table.GetDatabaseFileProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.DatabaseFileProgressData(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ValidateDatabaseTableSchema 验证数据库表结构
// @router /api/memory/table_schema/validate [POST]
func ValidateDatabaseTableSchema(c *gin.Context) {
	var req table.ValidateTableSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40000,
			"msg":  err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := memory.DatabaseApplicationSVC.ValidateDatabaseTableSchema(ctx, &req)
	if err != nil {
		httputil.InternalErrorGin(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

