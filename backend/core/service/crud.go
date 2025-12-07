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

	"etrxlite/tools/crud"
)

// CRUDService CRUD操作服务基类
type CRUDService struct {
	*BaseServiceImpl
	crud      *crud.EntCRUD
	tableName string
}

// NewCRUDService 创建CRUD服务
func NewCRUDService(base *BaseServiceImpl, crudTool *crud.EntCRUD, tableName string) *CRUDService {
	return &CRUDService{
		BaseServiceImpl: base,
		crud:            crudTool,
		tableName:      tableName,
	}
}

// Create 创建资源
func (s *CRUDService) Create(ctx context.Context, data map[string]interface{}) (uint, error) {
	id, err := s.crud.Create(ctx, s.tableName, data)
	if err != nil {
		return 0, fmt.Errorf("failed to create %s: %w", s.tableName, err)
	}
	return uint(id), nil
}

// GetByID 根据ID获取资源
func (s *CRUDService) GetByID(ctx context.Context, id uint, result interface{}) error {
	if err := s.crud.GetByID(ctx, s.tableName, id, result); err != nil {
		return fmt.Errorf("failed to get %s: %w", s.tableName, err)
	}
	return nil
}

// Update 更新资源
func (s *CRUDService) Update(ctx context.Context, id uint, data map[string]interface{}) error {
	if err := s.crud.Update(ctx, s.tableName, id, data); err != nil {
		return fmt.Errorf("failed to update %s: %w", s.tableName, err)
	}
	return nil
}

// Delete 删除资源
func (s *CRUDService) Delete(ctx context.Context, id uint) error {
	if err := s.crud.Delete(ctx, s.tableName, id); err != nil {
		return fmt.Errorf("failed to delete %s: %w", s.tableName, err)
	}
	return nil
}

// CreateAndGet 创建并获取资源（常用模式）
func (s *CRUDService) CreateAndGet(ctx context.Context, data map[string]interface{}, result interface{}) error {
	id, err := s.Create(ctx, data)
	if err != nil {
		return err
	}
	return s.GetByID(ctx, id, result)
}

// GetTableName 获取表名
func (s *CRUDService) GetTableName() string {
	return s.tableName
}

// GetCRUD 获取CRUD工具
func (s *CRUDService) GetCRUD() *crud.EntCRUD {
	return s.crud
}

