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
	"encoding/json"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/extensions/etrxtable/models"
)

// TableRepositoryCache 表格仓库缓存层
type TableRepositoryCache struct {
	repo  *TableRepository
	cache CacheInterface
}

// CacheInterface 缓存接口
type CacheInterface interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
}

// NewTableRepositoryCache 创建带缓存的表格仓库
func NewTableRepositoryCache(repo *TableRepository, cache CacheInterface) *TableRepositoryCache {
	return &TableRepositoryCache{
		repo:  repo,
		cache: cache,
	}
}

// cacheKey 生成缓存键
func (r *TableRepositoryCache) cacheKey(id uint) string {
	return fmt.Sprintf("table:body:%d", id)
}

func (r *TableRepositoryCache) cacheKeyList(workspaceID uint, page, pageSize int) string {
	return fmt.Sprintf("table:list:workspace:%d:page:%d:size:%d", workspaceID, page, pageSize)
}

// GetByID 根据 ID 获取表格（带缓存）
func (r *TableRepositoryCache) GetByID(ctx context.Context, id uint) (*models.TableBody, error) {
	// 尝试从缓存获取
	cacheKey := r.cacheKey(id)
	if r.cache != nil {
		cached, err := r.cache.Get(ctx, cacheKey)
		if err == nil && len(cached) > 0 {
			var table models.TableBody
			if err := json.Unmarshal(cached, &table); err == nil {
				return &table, nil
			}
		}
	}

	// 从数据库获取
	table, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if r.cache != nil {
		data, err := json.Marshal(table)
		if err == nil {
			_ = r.cache.Set(ctx, cacheKey, data, 5*time.Minute)
		}
	}

	return table, nil
}

// List 获取表格列表（带缓存）
func (r *TableRepositoryCache) List(ctx context.Context, workspaceID uint, page, pageSize int) ([]models.TableBody, int64, error) {
	// 尝试从缓存获取
	cacheKey := r.cacheKeyList(workspaceID, page, pageSize)
	if r.cache != nil {
		cached, err := r.cache.Get(ctx, cacheKey)
		if err == nil && len(cached) > 0 {
			var result struct {
				Tables []models.TableBody `json:"tables"`
				Total  int64               `json:"total"`
			}
			if err := json.Unmarshal(cached, &result); err == nil {
				return result.Tables, result.Total, nil
			}
		}
	}

	// 从数据库获取
	tables, total, err := r.repo.List(ctx, workspaceID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 写入缓存
	if r.cache != nil {
		result := struct {
			Tables []models.TableBody `json:"tables"`
			Total  int64               `json:"total"`
		}{
			Tables: tables,
			Total:  total,
		}
		data, err := json.Marshal(result)
		if err == nil {
			_ = r.cache.Set(ctx, cacheKey, data, 2*time.Minute)
		}
	}

	return tables, total, nil
}

// Create 创建表格（清除相关缓存）
func (r *TableRepositoryCache) Create(ctx context.Context, table *models.TableBody) error {
	err := r.repo.Create(ctx, table)
	if err == nil && r.cache != nil {
		// 清除列表缓存
		_ = r.cache.DeletePattern(ctx, "table:list:workspace:*")
	}
	return err
}

// Update 更新表格（清除相关缓存）
func (r *TableRepositoryCache) Update(ctx context.Context, table *models.TableBody) error {
	err := r.repo.Update(ctx, table)
	if err == nil && r.cache != nil {
		// 清除单个表格缓存
		_ = r.cache.Delete(ctx, r.cacheKey(table.ID))
		// 清除列表缓存
		_ = r.cache.DeletePattern(ctx, "table:list:workspace:*")
	}
	return err
}

// Delete 删除表格（清除相关缓存）
func (r *TableRepositoryCache) Delete(ctx context.Context, id uint) error {
	err := r.repo.Delete(ctx, id)
	if err == nil && r.cache != nil {
		// 清除单个表格缓存
		_ = r.cache.Delete(ctx, r.cacheKey(id))
		// 清除列表缓存
		_ = r.cache.DeletePattern(ctx, "table:list:workspace:*")
	}
	return err
}

