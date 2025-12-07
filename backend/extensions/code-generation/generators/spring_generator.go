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

// SpringGenerator Spring Boot 代码生成器
type SpringGenerator struct {
	typeMapper *mappers.TypeMapper
	packageName string
}

// NewSpringGenerator 创建 Spring 生成器
func NewSpringGenerator(packageName string) *SpringGenerator {
	return &SpringGenerator{
		typeMapper: mappers.NewTypeMapper("java", "spring"),
		packageName: packageName,
	}
}

// Generate 生成 Spring Boot Entity 代码
func (g *SpringGenerator) Generate(table *models.TableBody, columns []models.TableColumn) (string, error) {
	var code strings.Builder
	
	// 包声明
	if g.packageName == "" {
		g.packageName = "com.example.model"
	}
	code.WriteString(fmt.Sprintf("package %s;\n\n", g.packageName))
	
	// 导入语句
	imports := g.typeMapper.GetImportStatements([]string{})
	code.WriteString(strings.Join(imports, "\n"))
	code.WriteString("\n\n")
	
	structName := toPascalCase(table.Name)
	
	// 类定义
	code.WriteString(fmt.Sprintf("@Entity\n"))
	code.WriteString(fmt.Sprintf("@Table(name = \"%s\")\n", table.Name))
	code.WriteString(fmt.Sprintf("public class %s {\n\n", structName))
	
	// ID 字段
	code.WriteString("\t@Id\n")
	code.WriteString("\t@GeneratedValue(strategy = GenerationType.IDENTITY)\n")
	code.WriteString("\tprivate Long id;\n\n")
	
	// 生成列字段
	for _, col := range columns {
		if col.Key == "id" {
			continue
		}
		
		fieldName := toCamelCase(col.Name)
		mappedType := g.typeMapper.MapType(col.Type)
		
		code.WriteString(fmt.Sprintf("\t@Column(name = \"%s\"", col.Key))
		if !col.Required {
			code.WriteString(", nullable = true")
		}
		if col.Unique {
			code.WriteString(", unique = true")
		}
		code.WriteString(")\n")
		
		if col.Description != "" {
			code.WriteString(fmt.Sprintf("\t// %s\n", col.Description))
		}
		
		code.WriteString(fmt.Sprintf("\tprivate %s %s;\n\n", mappedType, fieldName))
	}
	
	// 时间戳字段
	code.WriteString("\t@Column(name = \"created_at\")\n")
	code.WriteString("\tprivate LocalDateTime createdAt;\n\n")
	code.WriteString("\t@Column(name = \"updated_at\")\n")
	code.WriteString("\tprivate LocalDateTime updatedAt;\n\n")
	
	// Getters and Setters
	code.WriteString("\t// Getters and Setters\n\n")
	
	// ID getter/setter
	code.WriteString("\tpublic Long getId() {\n")
	code.WriteString("\t\treturn id;\n")
	code.WriteString("\t}\n\n")
	code.WriteString("\tpublic void setId(Long id) {\n")
	code.WriteString("\t\tthis.id = id;\n")
	code.WriteString("\t}\n\n")
	
	// 列字段的 getter/setter
	for _, col := range columns {
		if col.Key == "id" {
			continue
		}
		
		fieldName := toCamelCase(col.Name)
		mappedType := g.typeMapper.MapType(col.Type)
		
		// Getter
		code.WriteString(fmt.Sprintf("\tpublic %s get%s() {\n", mappedType, toPascalCase(col.Name)))
		code.WriteString(fmt.Sprintf("\t\treturn %s;\n", fieldName))
		code.WriteString("\t}\n\n")
		
		// Setter
		code.WriteString(fmt.Sprintf("\tpublic void set%s(%s %s) {\n", toPascalCase(col.Name), mappedType, fieldName))
		code.WriteString(fmt.Sprintf("\t\tthis.%s = %s;\n", fieldName, fieldName))
		code.WriteString("\t}\n\n")
	}
	
	// 时间戳 getter/setter
	code.WriteString("\tpublic LocalDateTime getCreatedAt() {\n")
	code.WriteString("\t\treturn createdAt;\n")
	code.WriteString("\t}\n\n")
	code.WriteString("\tpublic void setCreatedAt(LocalDateTime createdAt) {\n")
	code.WriteString("\t\tthis.createdAt = createdAt;\n")
	code.WriteString("\t}\n\n")
	code.WriteString("\tpublic LocalDateTime getUpdatedAt() {\n")
	code.WriteString("\t\treturn updatedAt;\n")
	code.WriteString("\t}\n\n")
	code.WriteString("\tpublic void setUpdatedAt(LocalDateTime updatedAt) {\n")
	code.WriteString("\t\tthis.updatedAt = updatedAt;\n")
	code.WriteString("\t}\n")
	
	code.WriteString("}\n")
	
	return code.String(), nil
}

// toCamelCase 转换为 camelCase
func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	
	var result strings.Builder
	for i, part := range parts {
		if len(part) > 0 {
			if i == 0 {
				result.WriteString(strings.ToLower(part))
			} else {
				result.WriteString(strings.ToUpper(part[:1]))
				if len(part) > 1 {
					result.WriteString(part[1:])
				}
			}
		}
	}
	
	return result.String()
}

