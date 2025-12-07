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

package mappers

import (
	"strings"
)

// TypeMapper 类型映射器
type TypeMapper struct {
	language string
	framework string
}

// NewTypeMapper 创建类型映射器
func NewTypeMapper(language, framework string) *TypeMapper {
	return &TypeMapper{
		language: language,
		framework: framework,
	}
}

// MapType 映射类型
func (m *TypeMapper) MapType(tableType string) string {
	tableType = strings.ToLower(tableType)
	
	switch m.language {
	case "go":
		return m.mapGoType(tableType, m.framework)
	case "python":
		return m.mapPythonType(tableType)
	case "java":
		return m.mapJavaType(tableType)
	default:
		return "string"
	}
}

// mapGoType 映射 Go 类型
func (m *TypeMapper) mapGoType(tableType, framework string) string {
	typeMap := map[string]string{
		"text":       "string",
		"number":     "int",
		"date":       "time.Time",
		"datetime":   "time.Time",
		"boolean":    "bool",
		"select":     "string",
		"multiselect": "string", // 或 []string
		"email":      "string",
		"url":        "string",
		"phone":      "string",
		"file":       "string",
		"image":      "string",
		"json":       "datatypes.JSON",
		"formula":    "string",
	}
	
	if mapped, ok := typeMap[tableType]; ok {
		return mapped
	}
	return "string"
}

// mapPythonType 映射 Python 类型
func (m *TypeMapper) mapPythonType(tableType string) string {
	typeMap := map[string]string{
		"text":       "str",
		"number":     "int",
		"date":       "date",
		"datetime":   "datetime",
		"boolean":    "bool",
		"select":     "str",
		"multiselect": "List[str]",
		"email":      "EmailStr",
		"url":        "HttpUrl",
		"phone":      "str",
		"file":       "str",
		"image":      "str",
		"json":       "dict",
		"formula":    "str",
	}
	
	if mapped, ok := typeMap[tableType]; ok {
		return mapped
	}
	return "str"
}

// mapJavaType 映射 Java 类型
func (m *TypeMapper) mapJavaType(tableType string) string {
	typeMap := map[string]string{
		"text":       "String",
		"number":     "Integer",
		"date":       "LocalDate",
		"datetime":   "LocalDateTime",
		"boolean":    "Boolean",
		"select":     "String",
		"multiselect": "List<String>",
		"email":      "String",
		"url":        "String",
		"phone":      "String",
		"file":       "String",
		"image":      "String",
		"json":       "String", // 或使用 JSONObject
		"formula":    "String",
	}
	
	if mapped, ok := typeMap[tableType]; ok {
		return mapped
	}
	return "String"
}

// GetImportStatements 获取导入语句
func (m *TypeMapper) GetImportStatements(types []string) []string {
	var imports []string
	
	switch m.language {
	case "go":
		imports = m.getGoImports(types, m.framework)
	case "python":
		imports = m.getPythonImports(types)
	case "java":
		imports = m.getJavaImports(types)
	}
	
	return imports
}

// getGoImports 获取 Go 导入语句
func (m *TypeMapper) getGoImports(types []string, framework string) []string {
	imports := []string{}
	
	hasTime := false
	hasJSON := false
	
	for _, t := range types {
		if strings.Contains(t, "time.Time") {
			hasTime = true
		}
		if strings.Contains(t, "datatypes.JSON") {
			hasJSON = true
		}
	}
	
	if hasTime {
		imports = append(imports, `"time"`)
	}
	
	if hasJSON {
		imports = append(imports, `"gorm.io/datatypes"`)
	}
	
	if framework == "gorm" {
		imports = append(imports, `"gorm.io/gorm"`)
	} else if framework == "ent" {
		imports = append(imports, `"entgo.io/ent"`, `"entgo.io/ent/schema/field"`)
	}
	
	return imports
}

// getPythonImports 获取 Python 导入语句
func (m *TypeMapper) getPythonImports(types []string) []string {
	imports := []string{`from pydantic import BaseModel`}
	
	hasDate := false
	hasDateTime := false
	hasList := false
	
	for _, t := range types {
		if strings.Contains(t, "date") {
			hasDate = true
		}
		if strings.Contains(t, "datetime") {
			hasDateTime = true
		}
		if strings.Contains(t, "List[") {
			hasList = true
		}
	}
	
	if hasDate {
		imports = append(imports, `from datetime import date`)
	}
	if hasDateTime {
		imports = append(imports, `from datetime import datetime`)
	}
	if hasList {
		imports = append(imports, `from typing import List`)
	}
	
	return imports
}

// getJavaImports 获取 Java 导入语句
func (m *TypeMapper) getJavaImports(types []string) []string {
	imports := []string{`javax.persistence.*`}
	
	hasDate := false
	hasDateTime := false
	hasList := false
	
	for _, t := range types {
		if strings.Contains(t, "LocalDate") {
			hasDate = true
		}
		if strings.Contains(t, "LocalDateTime") {
			hasDateTime = true
		}
		if strings.Contains(t, "List<") {
			hasList = true
		}
	}
	
	if hasDate || hasDateTime {
		imports = append(imports, `java.time.LocalDate`, `java.time.LocalDateTime`)
	}
	if hasList {
		imports = append(imports, `java.util.List`)
	}
	
	return imports
}

