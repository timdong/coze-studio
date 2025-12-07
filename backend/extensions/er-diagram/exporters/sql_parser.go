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
	"fmt"
	"regexp"
	"strings"
)

// ParseSQL 解析 SQL DDL 语句
func ParseSQL(sql string) (*DrawDBFormat, error) {
	drawdb := NewDrawDB()
	
	// 移除注释和多余空白
	sql = removeSQLComments(sql)
	
	// 分割 CREATE TABLE 语句
	createTableRegex := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?["']?(\w+)["']?\s*\(`)
	matches := createTableRegex.FindAllStringSubmatch(sql, -1)
	
	for _, match := range matches {
		tableName := match[1]
		tableStart := strings.Index(sql, match[0])
		
		// 找到表定义的结束位置（匹配括号）
		tableDef := extractTableDefinition(sql[tableStart:])
		
		// 解析表定义
		table, err := parseTableDefinition(tableName, tableDef)
		if err != nil {
			return nil, fmt.Errorf("failed to parse table %s: %w", tableName, err)
		}
		
		drawdb.Tables = append(drawdb.Tables, *table)
	}
	
	// 解析外键关系
	drawdb.Relationships = parseForeignKeys(sql)
	
	return drawdb, nil
}

// removeSQLComments 移除 SQL 注释
func removeSQLComments(sql string) string {
	// 移除单行注释 --
	lines := strings.Split(sql, "\n")
	var cleaned []string
	for _, line := range lines {
		if idx := strings.Index(line, "--"); idx >= 0 {
			line = line[:idx]
		}
		cleaned = append(cleaned, strings.TrimSpace(line))
	}
	return strings.Join(cleaned, "\n")
}

// extractTableDefinition 提取表定义内容
func extractTableDefinition(sql string) string {
	// 找到第一个左括号
	start := strings.Index(sql, "(")
	if start == -1 {
		return ""
	}
	
	// 匹配括号对
	depth := 0
	for i := start; i < len(sql); i++ {
		if sql[i] == '(' {
			depth++
		} else if sql[i] == ')' {
			depth--
			if depth == 0 {
				return sql[start+1 : i]
			}
		}
	}
	return sql[start+1:]
}

// parseTableDefinition 解析表定义
func parseTableDefinition(tableName, tableDef string) (*Table, error) {
	table := &Table{
		Name:    tableName,
		Columns: []Column{},
	}
	
	// 分割列定义（按逗号，但要注意括号内的逗号）
	columnDefs := splitColumnDefinitions(tableDef)
	
	for _, colDef := range columnDefs {
		colDef = strings.TrimSpace(colDef)
		if colDef == "" {
			continue
		}
		
		// 跳过约束定义（PRIMARY KEY, FOREIGN KEY 等）
		if strings.HasPrefix(strings.ToUpper(colDef), "PRIMARY KEY") ||
			strings.HasPrefix(strings.ToUpper(colDef), "FOREIGN KEY") ||
			strings.HasPrefix(strings.ToUpper(colDef), "UNIQUE") ||
			strings.HasPrefix(strings.ToUpper(colDef), "INDEX") ||
			strings.HasPrefix(strings.ToUpper(colDef), "CONSTRAINT") {
			// 处理主键约束
			if strings.HasPrefix(strings.ToUpper(colDef), "PRIMARY KEY") {
				pkCols := extractPrimaryKeyColumns(colDef)
				table.PrimaryKeys = append(table.PrimaryKeys, pkCols...)
			}
			continue
		}
		
		// 解析列定义
		column, err := parseColumnDefinition(colDef)
		if err != nil {
			return nil, fmt.Errorf("failed to parse column: %w", err)
		}
		table.Columns = append(table.Columns, *column)
	}
	
	return table, nil
}

// splitColumnDefinitions 分割列定义
func splitColumnDefinitions(tableDef string) []string {
	var definitions []string
	var current strings.Builder
	depth := 0
	
	for _, char := range tableDef {
		if char == '(' {
			depth++
			current.WriteRune(char)
		} else if char == ')' {
			depth--
			current.WriteRune(char)
		} else if char == ',' && depth == 0 {
			def := strings.TrimSpace(current.String())
			if def != "" {
				definitions = append(definitions, def)
			}
			current.Reset()
		} else {
			current.WriteRune(char)
		}
	}
	
	// 添加最后一个定义
	def := strings.TrimSpace(current.String())
	if def != "" {
		definitions = append(definitions, def)
	}
	
	return definitions
}

// parseColumnDefinition 解析列定义
func parseColumnDefinition(colDef string) (*Column, error) {
	// 移除引号
	colDef = strings.ReplaceAll(colDef, `"`, "")
	colDef = strings.ReplaceAll(colDef, `'`, "")
	
	parts := strings.Fields(colDef)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid column definition: %s", colDef)
	}
	
	column := &Column{
		Name:     parts[0],
		Type:     parts[1],
		Nullable: true,
	}
	
	// 检查 NOT NULL
	if strings.Contains(strings.ToUpper(colDef), "NOT NULL") {
		column.Nullable = false
	}
	
	// 检查 DEFAULT
	if defaultMatch := regexp.MustCompile(`(?i)DEFAULT\s+([^\s,]+)`).FindStringSubmatch(colDef); len(defaultMatch) > 1 {
		column.DefaultValue = strings.Trim(defaultMatch[1], "'\"")
	}
	
	// 检查 COMMENT
	if commentMatch := regexp.MustCompile(`(?i)COMMENT\s+['"]([^'"]+)['"]`).FindStringSubmatch(colDef); len(commentMatch) > 1 {
		column.Comment = commentMatch[1]
	}
	
	return column, nil
}

// extractPrimaryKeyColumns 提取主键列
func extractPrimaryKeyColumns(pkDef string) []string {
	// 匹配 PRIMARY KEY (col1, col2, ...)
	re := regexp.MustCompile(`(?i)PRIMARY\s+KEY\s*\(([^)]+)\)`)
	matches := re.FindStringSubmatch(pkDef)
	if len(matches) > 1 {
		cols := strings.Split(matches[1], ",")
		var result []string
		for _, col := range cols {
			result = append(result, strings.TrimSpace(strings.Trim(col, `"'`)))
		}
		return result
	}
	return []string{}
}

// parseForeignKeys 解析外键关系
func parseForeignKeys(sql string) []Relationship {
	var relationships []Relationship
	
	// 匹配 FOREIGN KEY 约束
	fkRegex := regexp.MustCompile(`(?i)FOREIGN\s+KEY\s*\(["']?(\w+)["']?\)\s+REFERENCES\s+["']?(\w+)["']?\s*\(["']?(\w+)["']?\)`)
	matches := fkRegex.FindAllStringSubmatch(sql, -1)
	
	for _, match := range matches {
		if len(match) >= 4 {
			relationships = append(relationships, Relationship{
				FromTable:  "", // 需要从上下文获取
				FromColumn: match[1],
				ToTable:    match[2],
				ToColumn:   match[3],
				Type:       "one-to-many",
			})
		}
	}
	
	return relationships
}

