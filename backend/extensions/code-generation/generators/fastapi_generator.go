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

package generators

import (
	"fmt"
	"strings"
	
	"github.com/coze-dev/coze-studio/backend/extensions/code-generation/mappers"
	"github.com/coze-dev/coze-studio/backend/extensions/etrxtable/models"
)

// FastAPIGenerator FastAPI 代码生成器
type FastAPIGenerator struct {
	typeMapper *mappers.TypeMapper
}

// NewFastAPIGenerator 创建 FastAPI 生成器
func NewFastAPIGenerator() *FastAPIGenerator {
	return &FastAPIGenerator{
		typeMapper: mappers.NewTypeMapper("python", "fastapi"),
	}
}

// Generate 生成 FastAPI 模型代码
func (g *FastAPIGenerator) Generate(table *models.TableBody, columns []models.TableColumn) (string, error) {
	var code strings.Builder
	
	// 导入语句
	imports := g.typeMapper.GetImportStatements([]string{})
	code.WriteString(strings.Join(imports, "\n"))
	code.WriteString("\n\n")
	
	structName := toPascalCase(table.Name)
	
	// Base 模型
	code.WriteString(fmt.Sprintf("class %sBase(BaseModel):\n", structName))
	for _, col := range columns {
		if col.Key == "id" {
			continue // ID 在完整模型中
		}
		
		fieldName := toSnakeCase(col.Name)
		mappedType := g.typeMapper.MapType(col.Type)
		
		fieldDef := fmt.Sprintf("\t%s: %s", fieldName, mappedType)
		if !col.Required {
			fieldDef = fmt.Sprintf("\t%s: Optional[%s] = None", fieldName, mappedType)
		}
		
		if col.Description != "" {
			fieldDef += fmt.Sprintf("  # %s", col.Description)
		}
		
		code.WriteString(fieldDef + "\n")
	}
	code.WriteString("\n")
	
	// Create 模型
	code.WriteString(fmt.Sprintf("class %sCreate(%sBase):\n", structName, structName))
	code.WriteString("\tpass\n\n")
	
	// Update 模型
	code.WriteString(fmt.Sprintf("class %sUpdate(%sBase):\n", structName, structName))
	code.WriteString("\tpass\n\n")
	
	// 完整模型
	code.WriteString(fmt.Sprintf("class %s(%sBase):\n", structName, structName))
	code.WriteString("\tid: int\n")
	code.WriteString("\tcreated_at: datetime\n")
	code.WriteString("\tupdated_at: datetime\n")
	code.WriteString("\n")
	code.WriteString("\tclass Config:\n")
	code.WriteString("\t\torm_mode = True\n")
	
	return code.String(), nil
}

// toSnakeCase 转换为 snake_case
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}
	
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	
	return strings.ToLower(result.String())
}

