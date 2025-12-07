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
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"github.com/coze-dev/coze-studio/backend/extensions/er-diagram/exporters"
	"github.com/coze-dev/coze-studio/backend/extensions/er-diagram/models"
	"github.com/coze-dev/coze-studio/backend/extensions/er-diagram/repository"
)

// ERDiagramService ER 图服务层
type ERDiagramService struct {
	diagramRepo *repository.ERDiagramRepository
	db          *gorm.DB
}

// NewERDiagramService 创建 ER 图服务
func NewERDiagramService(db *gorm.DB) *ERDiagramService {
	return &ERDiagramService{
		diagramRepo: repository.NewERDiagramRepository(db),
		db:          db,
	}
}

// CreateDiagram 创建 ER 图
func (s *ERDiagramService) CreateDiagram(ctx context.Context, diagram *models.ERDiagram) error {
	// 验证必填字段
	if diagram.Name == "" {
		return fmt.Errorf("diagram name is required")
	}
	if len(diagram.DiagramData) == 0 {
		return fmt.Errorf("diagram_data is required")
	}
	if diagram.WorkspaceID == 0 {
		return fmt.Errorf("workspace_id is required")
	}
	if diagram.CreatedBy == 0 {
		return fmt.Errorf("created_by is required")
	}
	
	// 验证 DiagramData 是否为有效的 JSON（使用 datatypes.JSON，GORM 会自动处理）
	var data interface{}
	if err := json.Unmarshal(diagram.DiagramData, &data); err != nil {
		return fmt.Errorf("invalid diagram_data JSON: %w", err)
	}
	
	// 设置默认值
	if diagram.DatabaseType == "" {
		diagram.DatabaseType = "postgresql"
	}
	if len(diagram.LinkedTableIDs) == 0 {
		diagram.LinkedTableIDs = []byte("[]")
	}
	if len(diagram.SyncConfig) == 0 {
		diagram.SyncConfig = []byte("{}")
	}
	
	return s.diagramRepo.Create(ctx, diagram)
}

// GetDiagram 获取 ER 图
func (s *ERDiagramService) GetDiagram(ctx context.Context, id uint) (*models.ERDiagram, error) {
	return s.diagramRepo.GetByID(ctx, id)
}

// ListDiagrams 获取 ER 图列表
func (s *ERDiagramService) ListDiagrams(ctx context.Context, workspaceID uint, page, pageSize int, filters map[string]interface{}) ([]models.ERDiagram, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	return s.diagramRepo.List(ctx, workspaceID, page, pageSize, filters)
}

// UpdateDiagram 更新 ER 图
func (s *ERDiagramService) UpdateDiagram(ctx context.Context, diagram *models.ERDiagram) error {
	// 验证 ER 图是否存在
	existing, err := s.diagramRepo.GetByID(ctx, diagram.ID)
	if err != nil {
		return fmt.Errorf("diagram not found: %w", err)
	}
	
	// DiagramData 使用 datatypes.JSON，GORM 会自动验证 JSON 格式
	if len(diagram.DiagramData) == 0 {
		diagram.DiagramData = existing.DiagramData
	}
	
	// 保留一些不可修改的字段
	diagram.CreatedAt = existing.CreatedAt
	diagram.WorkspaceID = existing.WorkspaceID
	diagram.CreatedBy = existing.CreatedBy
	
	return s.diagramRepo.Update(ctx, diagram)
}

// DeleteDiagram 删除 ER 图
func (s *ERDiagramService) DeleteDiagram(ctx context.Context, id uint) error {
	// 验证 ER 图是否存在
	_, err := s.diagramRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("diagram not found: %w", err)
	}
	
	return s.diagramRepo.Delete(ctx, id)
}

// ExportToSQL 导出为 SQL
func (s *ERDiagramService) ExportToSQL(ctx context.Context, id uint) (string, error) {
	diagram, err := s.diagramRepo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("diagram not found: %w", err)
	}
	
	// 解析 DrawDB 格式
	drawdb, err := exporters.ParseDrawDB(string(diagram.DiagramData))
	if err != nil {
		return "", fmt.Errorf("failed to parse diagram data: %w", err)
	}
	
	// 导出为 SQL
	dbType := diagram.DatabaseType
	if dbType == "" {
		dbType = "postgresql"
	}
	
	sql, err := exporters.ExportToSQL(drawdb, dbType)
	if err != nil {
		return "", fmt.Errorf("failed to export to SQL: %w", err)
	}
	
	return sql, nil
}

// ExportToDBML 导出为 DBML
func (s *ERDiagramService) ExportToDBML(ctx context.Context, id uint) (string, error) {
	diagram, err := s.diagramRepo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("diagram not found: %w", err)
	}
	
	// 解析 DrawDB 格式
	drawdb, err := exporters.ParseDrawDB(string(diagram.DiagramData))
	if err != nil {
		return "", fmt.Errorf("failed to parse diagram data: %w", err)
	}
	
	// 导出为 DBML
	dbml, err := exporters.ExportToDBML(drawdb)
	if err != nil {
		return "", fmt.Errorf("failed to export to DBML: %w", err)
	}
	
	return dbml, nil
}

// ImportFromSQL 从 SQL 导入
func (s *ERDiagramService) ImportFromSQL(ctx context.Context, sql string) (*models.ERDiagram, error) {
	// 解析 SQL
	drawdb, err := exporters.ParseSQL(sql)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SQL: %w", err)
	}
	
	// 转换为 JSON
	diagramData, err := drawdb.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to convert to JSON: %w", err)
	}
	
	// 创建 ER 图对象
	diagram := &models.ERDiagram{
		DiagramData:    []byte(diagramData),
		DatabaseType:   "postgresql", // 默认 PostgreSQL
		LinkedTableIDs: []byte("[]"),
		SyncConfig:     []byte("{}"),
	}
	
	return diagram, nil
}

// ImportFromDBML 从 DBML 导入
func (s *ERDiagramService) ImportFromDBML(ctx context.Context, dbml string) (*models.ERDiagram, error) {
	// 解析 DBML
	drawdb, err := exporters.ParseDBML(dbml)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DBML: %w", err)
	}
	
	// 转换为 JSON
	diagramData, err := drawdb.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to convert to JSON: %w", err)
	}
	
	// 创建 ER 图对象
	diagram := &models.ERDiagram{
		DiagramData:    []byte(diagramData),
		DatabaseType:   "postgresql", // 默认 PostgreSQL
		LinkedTableIDs: []byte("[]"),
		SyncConfig:     []byte("{}"),
	}
	
	return diagram, nil
}

