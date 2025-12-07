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
	"github.com/coze-dev/coze-studio/backend/extensions/code-generation/generators"
	"github.com/coze-dev/coze-studio/backend/extensions/etrxtable/models"
	"github.com/coze-dev/coze-studio/backend/extensions/etrxtable/repository"
)

// CodeGenService 代码生成服务层
type CodeGenService struct {
	tableRepo  *repository.TableRepository
	columnRepo *repository.TableColumnRepository
	db         *gorm.DB
}

// NewCodeGenService 创建代码生成服务
func NewCodeGenService(db *gorm.DB) *CodeGenService {
	return &CodeGenService{
		tableRepo:  repository.NewTableRepository(db),
		columnRepo: repository.NewTableColumnRepository(db),
		db:         db,
	}
}

// GenerateCodeRequest 代码生成请求
type GenerateCodeRequest struct {
	TableID      uint   `json:"table_id" binding:"required"`
	Language     string `json:"language" binding:"required"` // go, python, java
	Framework    string `json:"framework"`                   // ent, fastapi, spring
	PackageName  string `json:"package_name"`
	OutputFormat string `json:"output_format"` // zip, files
}

// GenerateCodeResponse 代码生成响应
type GenerateCodeResponse struct {
	Code     string            `json:"code"`
	Files    map[string]string `json:"files,omitempty"` // 文件名 -> 内容
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// GenerateCode 生成代码
func (s *CodeGenService) GenerateCode(ctx context.Context, req *GenerateCodeRequest) (*GenerateCodeResponse, error) {
	// 获取表格信息
	table, err := s.tableRepo.GetByID(ctx, req.TableID)
	if err != nil {
		return nil, fmt.Errorf("table not found: %w", err)
	}
	
	// 根据语言和框架生成代码
	switch req.Language {
	case "go":
		return s.generateGoCode(ctx, table, req)
	case "python":
		return s.generatePythonCode(ctx, table, req)
	case "java":
		return s.generateJavaCode(ctx, table, req)
	default:
		return nil, fmt.Errorf("unsupported language: %s", req.Language)
	}
}

// generateGoCode 生成 Go 代码
func (s *CodeGenService) generateGoCode(ctx context.Context, table *models.TableBody, req *GenerateCodeRequest) (*GenerateCodeResponse, error) {
	if req.Framework == "ent" {
		return s.generateGoEntCode(ctx, table, req)
	}
	
	// 默认生成 GORM 模型
	return s.generateGoGORMCode(ctx, table, req)
}

// generateGoEntCode 生成 Go Ent 代码
func (s *CodeGenService) generateGoEntCode(ctx context.Context, table *models.TableBody, req *GenerateCodeRequest) (*GenerateCodeResponse, error) {
	// 获取表格的列
	columns, err := s.columnRepo.GetByTableID(ctx, table.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns: %w", err)
	}
	
	// 使用 Ent 生成器
	packageName := req.PackageName
	if packageName == "" {
		packageName = "schema"
	}
	
	generator := generators.NewEntGenerator(packageName)
	code, err := generator.Generate(table, columns)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ent code: %w", err)
	}
	
	return &GenerateCodeResponse{
		Code: code,
		Metadata: map[string]interface{}{
			"table_name": table.Name,
			"language":   "go",
			"framework":  "ent",
			"columns":    len(columns),
		},
	}, nil
}

// generateGoGORMCode 生成 Go GORM 代码
func (s *CodeGenService) generateGoGORMCode(ctx context.Context, table *models.TableBody, req *GenerateCodeRequest) (*GenerateCodeResponse, error) {
	// 获取表格的列
	columns, err := s.columnRepo.GetByTableID(ctx, table.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns: %w", err)
	}
	
	// 使用 GORM 生成器
	packageName := req.PackageName
	if packageName == "" {
		packageName = "models"
	}
	
	generator := generators.NewGORMGenerator(packageName)
	code, err := generator.Generate(table, columns)
	if err != nil {
		return nil, fmt.Errorf("failed to generate GORM code: %w", err)
	}
	
	return &GenerateCodeResponse{
		Code: code,
		Metadata: map[string]interface{}{
			"table_name": table.Name,
			"language":   "go",
			"framework":  "gorm",
			"columns":    len(columns),
		},
	}, nil
}

// generatePythonCode 生成 Python 代码
func (s *CodeGenService) generatePythonCode(ctx context.Context, table *models.TableBody, req *GenerateCodeRequest) (*GenerateCodeResponse, error) {
	if req.Framework == "fastapi" {
		return s.generateFastAPICode(ctx, table, req)
	}
	
	return nil, fmt.Errorf("unsupported Python framework: %s", req.Framework)
}

// generateFastAPICode 生成 FastAPI 代码
func (s *CodeGenService) generateFastAPICode(ctx context.Context, table *models.TableBody, req *GenerateCodeRequest) (*GenerateCodeResponse, error) {
	// 获取表格的列
	columns, err := s.columnRepo.GetByTableID(ctx, table.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns: %w", err)
	}
	
	// 使用 FastAPI 生成器
	generator := generators.NewFastAPIGenerator()
	code, err := generator.Generate(table, columns)
	if err != nil {
		return nil, fmt.Errorf("failed to generate FastAPI code: %w", err)
	}
	
	return &GenerateCodeResponse{
		Code: code,
		Metadata: map[string]interface{}{
			"table_name": table.Name,
			"language":   "python",
			"framework":  "fastapi",
			"columns":    len(columns),
		},
	}, nil
}

// generateJavaCode 生成 Java 代码
func (s *CodeGenService) generateJavaCode(ctx context.Context, table *models.TableBody, req *GenerateCodeRequest) (*GenerateCodeResponse, error) {
	if req.Framework == "spring" {
		return s.generateSpringCode(ctx, table, req)
	}
	
	return nil, fmt.Errorf("unsupported Java framework: %s", req.Framework)
}

// generateSpringCode 生成 Spring Boot 代码
func (s *CodeGenService) generateSpringCode(ctx context.Context, table *models.TableBody, req *GenerateCodeRequest) (*GenerateCodeResponse, error) {
	// 获取表格的列
	columns, err := s.columnRepo.GetByTableID(ctx, table.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns: %w", err)
	}
	
	// 使用 Spring 生成器
	packageName := req.PackageName
	if packageName == "" {
		packageName = "com.example.model"
	}
	
	generator := generators.NewSpringGenerator(packageName)
	code, err := generator.Generate(table, columns)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Spring code: %w", err)
	}
	
	return &GenerateCodeResponse{
		Code: code,
		Metadata: map[string]interface{}{
			"table_name":   table.Name,
			"language":     "java",
			"framework":    "spring",
			"package_name": packageName,
			"columns":      len(columns),
		},
	}, nil
}

