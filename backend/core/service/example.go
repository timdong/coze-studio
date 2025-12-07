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

package service

// 这个文件展示了如何使用基础服务框架重构现有服务
// 这是一个示例，不会在实际代码中使用

/*
示例：重构后的EtrxtableService（完整版）

package services

import (
	"context"
	"fmt"
	
	"github.com/coze-dev/coze-studio/backend/core/errors"
	"github.com/coze-dev/coze-studio/backend/core/event"
	"github.com/coze-dev/coze-studio/backend/core/service"
	"etrxlite/platform/ent"
	"etrxlite/tools/query"
)

// EtrxtableService 重构后的etrxtable表格服务
type EtrxtableService struct {
	*service.CRUDService
}

// NewEtrxtableService 创建重构后的etrxtable服务
func NewEtrxtableService(manager *Manager) *EtrxtableService {
	base := service.NewBaseService(manager)
	crudService := service.NewCRUDService(base, manager.GetCRUD(), "table_bodies")
	return &EtrxtableService{
		CRUDService: crudService,
	}
}

// HealthCheck 健康检查（可以覆盖默认实现）
func (s *EtrxtableService) HealthCheck(ctx context.Context) error {
	// 先调用基类的健康检查
	if err := s.BaseServiceImpl.HealthCheck(ctx); err != nil {
		return err
	}
	// 可以添加额外的健康检查逻辑
	return nil
}

// CreateTable 创建表格（使用基类的CreateAndGet方法，并发布事件和记录审计日志）
func (s *EtrxtableService) CreateTable(ctx context.Context, data map[string]interface{}) (*ent.TableBody, error) {
	var table ent.TableBody
	if err := s.CreateAndGet(ctx, data, &table); err != nil {
		// 使用错误包装
		return nil, errors.WrapError("table-service", "CreateTable", "创建表格失败", err)
	}
	
	// 发布资源创建事件
	_ = s.PublishResourceCreated(ctx, "table", table.ID, table)
	
	// 记录审计日志
	_ = s.LogCreate(ctx, "table", fmt.Sprintf("%d", table.ID), data)
	
	return &table, nil
}

// GetTable 获取单个表格（使用基类的GetByID方法）
func (s *EtrxtableService) GetTable(ctx context.Context, id uint) (*ent.TableBody, error) {
	var table ent.TableBody
	if err := s.GetByID(ctx, id, &table); err != nil {
		return nil, errors.WrapError("table-service", "GetTable", "获取表格失败", err)
	}
	return &table, nil
}

// GetTables 获取表格列表（使用查询构建器和分页工具）
func (s *EtrxtableService) GetTables(ctx context.Context, workspaceID uint, page, pageSize int) ([]*ent.TableBody, *service.PaginationResponse, error) {
	normalizedPage, normalizedPageSize, _ := service.NormalizePagination(page, pageSize)
	
	// 使用查询构建器
	builder := query.NewBuilder("table_bodies").
		Filter("workspace_id", workspaceID).
		OrderBy("created_at", "DESC").
		Paginate(normalizedPage, normalizedPageSize)
	
	result, err := builder.Execute(ctx, s.GetSearch())
	if err != nil {
		return nil, nil, errors.WrapError("table-service", "GetTables", "获取表格列表失败", err)
	}
	
	// 转换结果...
	tables := make([]*ent.TableBody, 0, len(result.Items))
	for _, item := range result.Items {
		if table, ok := item.(*ent.TableBody); ok {
			tables = append(tables, table)
		}
	}
	
	// 创建分页响应
	pagination := service.NewPaginationResponse(normalizedPage, normalizedPageSize, int(result.Total))
	
	return tables, pagination, nil
}

// UpdateTable 更新表格（使用事务、事件和审计日志）
func (s *EtrxtableService) UpdateTable(ctx context.Context, id uint, data map[string]interface{}) (*ent.TableBody, error) {
	// 使用事务管理器
	txManager := service.GetTransactionManager(s.GetManager())
	
	var table ent.TableBody
	err := txManager.WithTransaction(ctx, func(tx *ent.Tx) error {
		// 更新操作
		if err := s.Update(ctx, id, data); err != nil {
			return err
		}
		
		// 获取更新后的数据
		return s.GetByID(ctx, id, &table)
	})
	
	if err != nil {
		return nil, errors.WrapError("table-service", "UpdateTable", "更新表格失败", err)
	}
	
	// 发布资源更新事件
	_ = s.PublishResourceUpdated(ctx, "table", id, table)
	
	// 记录审计日志
	_ = s.LogUpdate(ctx, "table", fmt.Sprintf("%d", id), data)
	
	return &table, nil
}

// DeleteTable 删除表格（使用事务、事件和审计日志）
func (s *EtrxtableService) DeleteTable(ctx context.Context, id uint) error {
	// 先获取表格信息（用于事件和审计日志）
	var table ent.TableBody
	if err := s.GetByID(ctx, id, &table); err != nil {
		return errors.WrapError("table-service", "DeleteTable", "获取表格失败", err)
	}
	
	// 使用事务管理器
	txManager := service.GetTransactionManager(s.GetManager())
	err := txManager.WithTransaction(ctx, func(tx *ent.Tx) error {
		return s.Delete(ctx, id)
	})
	
	if err != nil {
		return errors.WrapError("table-service", "DeleteTable", "删除表格失败", err)
	}
	
	// 发布资源删除事件
	_ = s.PublishResourceDeleted(ctx, "table", id, table)
	
	// 记录审计日志
	_ = s.LogDelete(ctx, "table", fmt.Sprintf("%d", id), nil)
	
	return nil
}
*/

