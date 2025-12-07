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

package exporters

import (
	"regexp"
	"strings"
)

// ParseDBML 解析 DBML 格式
func ParseDBML(dbml string) (*DrawDBFormat, error) {
	drawdb := NewDrawDB()
	
	// 移除注释
	dbml = removeDBMLComments(dbml)
	
	// 解析表定义
	tables := parseDBMLTables(dbml)
	drawdb.Tables = tables
	
	// 解析关系
	relationships := parseDBMLRelationships(dbml)
	drawdb.Relationships = relationships
	
	return drawdb, nil
}

// removeDBMLComments 移除 DBML 注释
func removeDBMLComments(dbml string) string {
	// 移除单行注释 //
	lines := strings.Split(dbml, "\n")
	var cleaned []string
	for _, line := range lines {
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = line[:idx]
		}
		cleaned = append(cleaned, strings.TrimSpace(line))
	}
	return strings.Join(cleaned, "\n")
}

// parseDBMLTables 解析 DBML 表定义
func parseDBMLTables(dbml string) []Table {
	var tables []Table
	
	// 匹配表定义: Table table_name { ... }
	tableRegex := regexp.MustCompile(`(?i)Table\s+(\w+)\s*\{([^}]+)\}`)
	matches := tableRegex.FindAllStringSubmatch(dbml, -1)
	
	for _, match := range matches {
		if len(match) >= 3 {
			tableName := match[1]
			tableBody := match[2]
			
			table := Table{
				Name:    tableName,
				Columns: parseDBMLColumns(tableBody),
			}
			
			// 解析表注释
			if commentMatch := regexp.MustCompile(`(?i)Note:\s*['"]([^'"]+)['"]`).FindStringSubmatch(tableBody); len(commentMatch) > 1 {
				table.Comment = commentMatch[1]
			}
			
			tables = append(tables, table)
		}
	}
	
	return tables
}

// parseDBMLColumns 解析 DBML 列定义
func parseDBMLColumns(tableBody string) []Column {
	var columns []Column
	
	// 匹配列定义: column_name type [settings]
	columnRegex := regexp.MustCompile(`(\w+)\s+(\w+(?:\([^)]+\))?)\s*([^\n]*)`)
	matches := columnRegex.FindAllStringSubmatch(tableBody, -1)
	
	for _, match := range matches {
		if len(match) >= 3 {
			col := Column{
				Name:     match[1],
				Type:     match[2],
				Nullable: true,
			}
			
			settings := match[3]
			
			// 检查主键
			if strings.Contains(settings, "pk") || strings.Contains(settings, "primary key") {
				// 主键会在表级别处理
			}
			
			// 检查 NOT NULL
			if strings.Contains(settings, "not null") {
				col.Nullable = false
			}
			
			// 检查默认值
			if defaultMatch := regexp.MustCompile(`default:\s*([^\s]+)`).FindStringSubmatch(settings); len(defaultMatch) > 1 {
				col.DefaultValue = strings.Trim(defaultMatch[1], `"'`)
			}
			
			// 检查注释
			if noteMatch := regexp.MustCompile(`note:\s*['"]([^'"]+)['"]`).FindStringSubmatch(settings); len(noteMatch) > 1 {
				col.Comment = noteMatch[1]
			}
			
			columns = append(columns, col)
		}
	}
	
	return columns
}

// parseDBMLRelationships 解析 DBML 关系定义
func parseDBMLRelationships(dbml string) []Relationship {
	var relationships []Relationship
	
	// 匹配关系定义: Ref: table1.column1 > table2.column2
	refRegex := regexp.MustCompile(`(?i)Ref:\s*(\w+)\.(\w+)\s*([><])\s*(\w+)\.(\w+)`)
	matches := refRegex.FindAllStringSubmatch(dbml, -1)
	
	for _, match := range matches {
		if len(match) >= 6 {
			relType := "one-to-many"
			if match[3] == ">" {
				relType = "one-to-many"
			} else if match[3] == "<" {
				relType = "many-to-one"
			}
			
			relationships = append(relationships, Relationship{
				FromTable:  match[1],
				FromColumn: match[2],
				ToTable:    match[4],
				ToColumn:   match[5],
				Type:       relType,
			})
		}
	}
	
	return relationships
}

