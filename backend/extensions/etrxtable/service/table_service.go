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

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/extensions/etrxtable/models"
	"github.com/coze-dev/coze-studio/backend/extensions/etrxtable/repository"
)

// TableService 表格服务层
type TableService struct {
	tableRepo  *repository.TableRepository
	columnRepo *repository.TableColumnRepository
	rowRepo    *repository.TableRowRepository
	db         *gorm.DB
}

// NewTableService 创建表格服务
func NewTableService(db *gorm.DB) *TableService {
	return &TableService{
		tableRepo:  repository.NewTableRepository(db),
		columnRepo: repository.NewTableColumnRepository(db),
		rowRepo:    repository.NewTableRowRepository(db),
		db:         db,
	}
}

// CreateTable 创建表格
func (s *TableService) CreateTable(ctx context.Context, table *models.TableBody) error {
	// 在事务中创建表格
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建表格
		if err := tx.Create(table).Error; err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
		
		// 如果有列定义，创建列
		if len(table.Columns) > 0 {
			for i := range table.Columns {
				table.Columns[i].TableBodyID = table.ID
			}
			if err := tx.Create(&table.Columns).Error; err != nil {
				return fmt.Errorf("failed to create columns: %w", err)
			}
		}
		
		return nil
	})
}

// GetTable 获取表格
func (s *TableService) GetTable(ctx context.Context, id uint) (*models.TableBody, error) {
	return s.tableRepo.GetByID(ctx, id)
}

// ListTables 获取表格列表
func (s *TableService) ListTables(ctx context.Context, workspaceID uint, page, pageSize int) ([]models.TableBody, int64, error) {
	return s.tableRepo.List(ctx, workspaceID, page, pageSize)
}

// UpdateTable 更新表格
func (s *TableService) UpdateTable(ctx context.Context, table *models.TableBody) error {
	return s.tableRepo.Update(ctx, table)
}

// DeleteTable 删除表格
func (s *TableService) DeleteTable(ctx context.Context, id uint) error {
	// 在事务中删除表格及其关联数据
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除表格的列
		if err := tx.Where("table_body_id = ?", id).Delete(&models.TableColumn{}).Error; err != nil {
			return fmt.Errorf("failed to delete columns: %w", err)
		}
		
		// 删除表格的行
		if err := tx.Where("table_body_id = ?", id).Delete(&models.TableRow{}).Error; err != nil {
			return fmt.Errorf("failed to delete rows: %w", err)
		}
		
		// 删除表格
		if err := tx.Delete(&models.TableBody{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete table: %w", err)
		}
		
		return nil
	})
}

// CreateColumn 创建列
func (s *TableService) CreateColumn(ctx context.Context, column *models.TableColumn) error {
	return s.columnRepo.Create(ctx, column)
}

// GetColumns 获取表格的所有列
func (s *TableService) GetColumns(ctx context.Context, tableID uint) ([]models.TableColumn, error) {
	return s.columnRepo.GetByTableID(ctx, tableID)
}

// UpdateColumn 更新列
func (s *TableService) UpdateColumn(ctx context.Context, column *models.TableColumn) error {
	return s.columnRepo.Update(ctx, column)
}

// DeleteColumn 删除列
func (s *TableService) DeleteColumn(ctx context.Context, id uint) error {
	return s.columnRepo.Delete(ctx, id)
}

// CreateRow 创建行
func (s *TableService) CreateRow(ctx context.Context, row *models.TableRow) error {
	return s.rowRepo.Create(ctx, row)
}

// GetRows 获取表格的所有行
func (s *TableService) GetRows(ctx context.Context, tableID uint, page, pageSize int) ([]models.TableRow, int64, error) {
	return s.rowRepo.GetByTableID(ctx, tableID, page, pageSize)
}

// UpdateRow 更新行
func (s *TableService) UpdateRow(ctx context.Context, row *models.TableRow) error {
	return s.rowRepo.Update(ctx, row)
}

// DeleteRow 删除行
func (s *TableService) DeleteRow(ctx context.Context, id uint) error {
	return s.rowRepo.Delete(ctx, id)
}

// BatchCreateRows 批量创建行
func (s *TableService) BatchCreateRows(ctx context.Context, rows []models.TableRow) error {
	return s.rowRepo.BatchCreate(ctx, rows)
}

// BatchUpdateRows 批量更新行
func (s *TableService) BatchUpdateRows(ctx context.Context, rows []models.TableRow) error {
	return s.rowRepo.BatchUpdate(ctx, rows)
}

// BatchDeleteRows 批量删除行
func (s *TableService) BatchDeleteRows(ctx context.Context, ids []uint) error {
	return s.rowRepo.BatchDelete(ctx, ids)
}

