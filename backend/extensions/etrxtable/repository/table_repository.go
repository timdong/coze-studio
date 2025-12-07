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

package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/extensions/etrxtable/models"
)

// TableRepository 表格数据访问层
type TableRepository struct {
	db *gorm.DB
}

// NewTableRepository 创建表格仓库
func NewTableRepository(db *gorm.DB) *TableRepository {
	return &TableRepository{db: db}
}

// Create 创建表格
func (r *TableRepository) Create(ctx context.Context, table *models.TableBody) error {
	return r.db.WithContext(ctx).Create(table).Error
}

// GetByID 根据 ID 获取表格
func (r *TableRepository) GetByID(ctx context.Context, id uint) (*models.TableBody, error) {
	var table models.TableBody
	err := r.db.WithContext(ctx).
		Preload("Workspace").
		Preload("Creator").
		Preload("Columns", func(db *gorm.DB) *gorm.DB {
			// 优化：只加载必要的列，按 sort_order 排序
			return db.Order("sort_order ASC")
		}).
		First(&table, id).Error
	if err != nil {
		return nil, err
	}
	return &table, nil
}

// List 获取表格列表
func (r *TableRepository) List(ctx context.Context, workspaceID uint, page, pageSize int) ([]models.TableBody, int64, error) {
	var tables []models.TableBody
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.TableBody{})
	
	// 过滤工作空间
	if workspaceID > 0 {
		query = query.Where("workspace_id = ?", workspaceID)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.
		Preload("Workspace").
		Preload("Creator").
		// 优化：不预加载 Columns 和 Rows，减少查询开销
		// 如果前端需要，可以通过单独的 API 获取
		Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&tables).Error
	
	return tables, total, err
}

// Update 更新表格
func (r *TableRepository) Update(ctx context.Context, table *models.TableBody) error {
	return r.db.WithContext(ctx).Save(table).Error
}

// Delete 删除表格（软删除）
func (r *TableRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TableBody{}, id).Error
}

// GetByWorkspaceAndName 根据工作空间和名称获取表格
func (r *TableRepository) GetByWorkspaceAndName(ctx context.Context, workspaceID uint, name string) (*models.TableBody, error) {
	var table models.TableBody
	err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND name = ?", workspaceID, name).
		First(&table).Error
	if err != nil {
		return nil, err
	}
	return &table, nil
}

// CountByWorkspace 统计工作空间的表格数量
func (r *TableRepository) CountByWorkspace(ctx context.Context, workspaceID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.TableBody{}).
		Where("workspace_id = ?", workspaceID).
		Count(&count).Error
	return count, err
}

// TableColumnRepository 表格列数据访问层
type TableColumnRepository struct {
	db *gorm.DB
}

// NewTableColumnRepository 创建表格列仓库
func NewTableColumnRepository(db *gorm.DB) *TableColumnRepository {
	return &TableColumnRepository{db: db}
}

// Create 创建列
func (r *TableColumnRepository) Create(ctx context.Context, column *models.TableColumn) error {
	return r.db.WithContext(ctx).Create(column).Error
}

// GetByTableID 获取表格的所有列
func (r *TableColumnRepository) GetByTableID(ctx context.Context, tableID uint) ([]models.TableColumn, error) {
	var columns []models.TableColumn
	err := r.db.WithContext(ctx).
		Where("table_body_id = ?", tableID).
		Order("sort_order ASC").
		Find(&columns).Error
	return columns, err
}

// Update 更新列
func (r *TableColumnRepository) Update(ctx context.Context, column *models.TableColumn) error {
	return r.db.WithContext(ctx).Save(column).Error
}

// Delete 删除列
func (r *TableColumnRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TableColumn{}, id).Error
}

// BatchCreate 批量创建列
func (r *TableColumnRepository) BatchCreate(ctx context.Context, columns []models.TableColumn) error {
	return r.db.WithContext(ctx).Create(&columns).Error
}

// TableRowRepository 表格行数据访问层
type TableRowRepository struct {
	db *gorm.DB
}

// NewTableRowRepository 创建表格行仓库
func NewTableRowRepository(db *gorm.DB) *TableRowRepository {
	return &TableRowRepository{db: db}
}

// Create 创建行
func (r *TableRowRepository) Create(ctx context.Context, row *models.TableRow) error {
	return r.db.WithContext(ctx).Create(row).Error
}

// GetByTableID 获取表格的所有行
func (r *TableRowRepository) GetByTableID(ctx context.Context, tableID uint, page, pageSize int) ([]models.TableRow, int64, error) {
	var rows []models.TableRow
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.TableRow{}).
		Where("table_body_id = ?", tableID)
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.
		Preload("Creator").
		Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&rows).Error
	
	return rows, total, err
}

// Update 更新行
func (r *TableRowRepository) Update(ctx context.Context, row *models.TableRow) error {
	return r.db.WithContext(ctx).Save(row).Error
}

// Delete 删除行
func (r *TableRowRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TableRow{}, id).Error
}

// BatchCreate 批量创建行
func (r *TableRowRepository) BatchCreate(ctx context.Context, rows []models.TableRow) error {
	if len(rows) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&rows).Error
}

// BatchUpdate 批量更新行
func (r *TableRowRepository) BatchUpdate(ctx context.Context, rows []models.TableRow) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, row := range rows {
			if err := tx.Save(&row).Error; err != nil {
				return fmt.Errorf("failed to update row %d: %w", row.ID, err)
			}
		}
		return nil
	})
}

// BatchDelete 批量删除行
func (r *TableRowRepository) BatchDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Delete(&models.TableRow{}, ids).Error
}

