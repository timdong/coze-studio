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

// GORMGenerator GORM 代码生成器
type GORMGenerator struct {
	typeMapper *mappers.TypeMapper
	packageName string
}

// NewGORMGenerator 创建 GORM 生成器
func NewGORMGenerator(packageName string) *GORMGenerator {
	return &GORMGenerator{
		typeMapper: mappers.NewTypeMapper("go", "gorm"),
		packageName: packageName,
	}
}

// Generate 生成 GORM 模型代码
func (g *GORMGenerator) Generate(table *models.TableBody, columns []models.TableColumn) (string, error) {
	var code strings.Builder
	
	// 包声明
	if g.packageName == "" {
		g.packageName = "models"
	}
	code.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	
	// 收集所有使用的类型
	var usedTypes []string
	for _, col := range columns {
		mappedType := g.typeMapper.MapType(col.Type)
		usedTypes = append(usedTypes, mappedType)
	}
	usedTypes = append(usedTypes, "time.Time", "gorm.io/gorm")
	
	// 导入语句
	imports := g.typeMapper.GetImportStatements(usedTypes)
	if len(imports) > 0 {
		code.WriteString("import (\n")
		for _, imp := range imports {
			code.WriteString(fmt.Sprintf("\t%s\n", imp))
		}
		code.WriteString(")\n\n")
	}
	
	// 结构体定义
	structName := toPascalCase(table.Name)
	code.WriteString(fmt.Sprintf("// %s represents the %s table\n", structName, table.Name))
	code.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	
	// 基础字段
	code.WriteString("\tID        uint           `gorm:\"primaryKey\" json:\"id\"`\n")
	
	// 生成列字段
	for _, col := range columns {
		fieldName := toPascalCase(col.Name)
		mappedType := g.typeMapper.MapType(col.Type)
		
		// GORM 标签
		var gormTags []string
		gormTags = append(gormTags, fmt.Sprintf("column:%s", col.Key))
		
		if col.Required {
			gormTags = append(gormTags, "not null")
		}
		
		if col.Unique {
			gormTags = append(gormTags, "unique")
		}
		
		if col.DefaultValue != "" {
			gormTags = append(gormTags, fmt.Sprintf("default:%s", col.DefaultValue))
		}
		
		// JSON 标签
		jsonTag := fmt.Sprintf("json:\"%s", col.Key)
		if !col.Required {
			jsonTag += ",omitempty"
		}
		jsonTag += "\""
		
		// 注释
		var comment string
		if col.Description != "" {
			comment = fmt.Sprintf(" // %s", col.Description)
		}
		
		code.WriteString(fmt.Sprintf("\t%s %s `gorm:\"%s\" %s`%s\n",
			fieldName,
			mappedType,
			strings.Join(gormTags, ";"),
			jsonTag,
			comment))
	}
	
	// 时间戳字段
	code.WriteString("\tCreatedAt time.Time      `gorm:\"autoCreateTime\" json:\"created_at\"`\n")
	code.WriteString("\tUpdatedAt time.Time      `gorm:\"autoUpdateTime\" json:\"updated_at\"`\n")
	code.WriteString("\tDeletedAt gorm.DeletedAt `gorm:\"index\" json:\"deleted_at,omitempty\"`\n")
	
	code.WriteString("}\n\n")
	
	// TableName 方法
	code.WriteString(fmt.Sprintf("func (%s) TableName() string {\n", structName))
	code.WriteString(fmt.Sprintf("\treturn \"%s\"\n", table.Name))
	code.WriteString("}\n")
	
	return code.String(), nil
}

// toPascalCase 转换为 PascalCase
func toPascalCase(s string) string {
	if s == "" {
		return s
	}
	
	// 处理下划线和连字符
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]))
			if len(part) > 1 {
				result.WriteString(part[1:])
			}
		}
	}
	
	return result.String()
}

