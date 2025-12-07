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
	"gorm.io/gorm"
	"github.com/coze-dev/coze-studio/backend/extensions/er-diagram/models"
)

// ERDiagramRepository ER 图数据访问层
type ERDiagramRepository struct {
	db *gorm.DB
}

// NewERDiagramRepository 创建 ER 图仓库
func NewERDiagramRepository(db *gorm.DB) *ERDiagramRepository {
	return &ERDiagramRepository{db: db}
}

// Create 创建 ER 图
func (r *ERDiagramRepository) Create(ctx context.Context, diagram *models.ERDiagram) error {
	return r.db.WithContext(ctx).Create(diagram).Error
}

// GetByID 根据 ID 获取 ER 图
func (r *ERDiagramRepository) GetByID(ctx context.Context, id uint) (*models.ERDiagram, error) {
	var diagram models.ERDiagram
	err := r.db.WithContext(ctx).
		Preload("Workspace").
		Preload("Creator").
		First(&diagram, id).Error
	if err != nil {
		return nil, err
	}
	return &diagram, nil
}

// List 获取 ER 图列表
func (r *ERDiagramRepository) List(ctx context.Context, workspaceID uint, page, pageSize int, filters map[string]interface{}) ([]models.ERDiagram, int64, error) {
	var diagrams []models.ERDiagram
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.ERDiagram{})
	
	// 过滤工作空间
	if workspaceID > 0 {
		query = query.Where("workspace_id = ?", workspaceID)
	}
	
	// 应用其他过滤条件
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if dbType, ok := filters["database_type"].(string); ok && dbType != "" {
		query = query.Where("database_type = ?", dbType)
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
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&diagrams).Error
	
	return diagrams, total, err
}

// Update 更新 ER 图
func (r *ERDiagramRepository) Update(ctx context.Context, diagram *models.ERDiagram) error {
	return r.db.WithContext(ctx).Save(diagram).Error
}

// Delete 删除 ER 图（软删除）
func (r *ERDiagramRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ERDiagram{}, id).Error
}

// GetByWorkspaceAndName 根据工作空间和名称获取 ER 图
func (r *ERDiagramRepository) GetByWorkspaceAndName(ctx context.Context, workspaceID uint, name string) (*models.ERDiagram, error) {
	var diagram models.ERDiagram
	err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND name = ?", workspaceID, name).
		First(&diagram).Error
	if err != nil {
		return nil, err
	}
	return &diagram, nil
}

