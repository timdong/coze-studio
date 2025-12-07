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

// EntGenerator Ent 代码生成器
type EntGenerator struct {
	typeMapper *mappers.TypeMapper
	packageName string
}

// NewEntGenerator 创建 Ent 生成器
func NewEntGenerator(packageName string) *EntGenerator {
	return &EntGenerator{
		typeMapper: mappers.NewTypeMapper("go", "ent"),
		packageName: packageName,
	}
}

// Generate 生成 Ent Schema 代码
func (g *EntGenerator) Generate(table *models.TableBody, columns []models.TableColumn) (string, error) {
	var code strings.Builder
	
	// 包声明
	if g.packageName == "" {
		g.packageName = "schema"
	}
	code.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	
	// 导入语句
	code.WriteString("import (\n")
	code.WriteString("\t\"entgo.io/ent\"\n")
	code.WriteString("\t\"entgo.io/ent/schema/field\"\n")
	
	// 检查是否需要 time 包
	needsTime := false
	for _, col := range columns {
		if col.Type == "date" || col.Type == "datetime" {
			needsTime = true
			break
		}
	}
	if needsTime {
		code.WriteString("\t\"time\"\n")
	}
	code.WriteString(")\n\n")
	
	// Schema 定义
	structName := toPascalCase(table.Name)
	code.WriteString(fmt.Sprintf("// %s represents the %s table\n", structName, table.Name))
	code.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	code.WriteString("\tent.Schema\n")
	code.WriteString("}\n\n")
	
	// Fields 方法
	code.WriteString(fmt.Sprintf("// Fields of the %s\n", structName))
	code.WriteString(fmt.Sprintf("func (%s) Fields() []ent.Field {\n", structName))
	code.WriteString("\treturn []ent.Field{\n")
	
	// ID 字段
	code.WriteString("\t\tfield.Uint(\"id\").\n")
	code.WriteString("\t\t\tUnique().\n")
	code.WriteString("\t\t\tImmutable(),\n")
	
	// 生成列字段
	for _, col := range columns {
		fieldName := col.Key
		mappedType := g.typeMapper.MapType(col.Type)
		
		code.WriteString(fmt.Sprintf("\t\tfield.%s(\"%s\")", g.mapToEntFieldType(mappedType), fieldName))
		
		if !col.Required {
			code.WriteString(".\n\t\t\tOptional()")
		}
		
		if col.Unique {
			code.WriteString(".\n\t\t\tUnique()")
		}
		
		if col.DefaultValue != "" {
			code.WriteString(fmt.Sprintf(".\n\t\t\tDefault(%s)", g.formatDefaultValue(col.DefaultValue, mappedType)))
		}
		
		if col.Description != "" {
			code.WriteString(fmt.Sprintf(".\n\t\t\tComment(\"%s\")", col.Description))
		}
		
		code.WriteString(",\n")
	}
	
	// 时间戳字段
	code.WriteString("\t\tfield.Time(\"created_at\").\n")
	code.WriteString("\t\t\tDefault(time.Now).\n")
	code.WriteString("\t\t\tImmutable(),\n")
	code.WriteString("\t\tfield.Time(\"updated_at\").\n")
	code.WriteString("\t\t\tDefault(time.Now).\n")
	code.WriteString("\t\t\tUpdateDefault(time.Now),\n")
	
	code.WriteString("\t}\n")
	code.WriteString("}\n")
	
	return code.String(), nil
}

// mapToEntFieldType 映射到 Ent 字段类型
func (g *EntGenerator) mapToEntFieldType(goType string) string {
	typeMap := map[string]string{
		"string": "String",
		"int":    "Int",
		"bool":   "Bool",
		"time.Time": "Time",
		"datatypes.JSON": "JSON",
	}
	
	if mapped, ok := typeMap[goType]; ok {
		return mapped
	}
	return "String"
}

// formatDefaultValue 格式化默认值
func (g *EntGenerator) formatDefaultValue(value, goType string) string {
	if goType == "string" {
		return fmt.Sprintf("\"%s\"", value)
	}
	if goType == "bool" {
		return value
	}
	if goType == "int" {
		return value
	}
	return fmt.Sprintf("\"%s\"", value)
}

