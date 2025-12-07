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
	"strings"
)

// ExportToSQL 将 DrawDB 格式导出为 SQL DDL
func ExportToSQL(drawdb *DrawDBFormat, dbType string) (string, error) {
	var sql strings.Builder
	
	for _, table := range drawdb.Tables {
		sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", escapeIdentifier(table.Name)))
		
		var columnDefs []string
		for _, col := range table.Columns {
			colDef := fmt.Sprintf("  %s %s", escapeIdentifier(col.Name), mapTypeToSQL(col.Type, dbType))
			
			if !col.Nullable {
				colDef += " NOT NULL"
			}
			
			if col.DefaultValue != "" {
				colDef += fmt.Sprintf(" DEFAULT %s", escapeValue(col.DefaultValue))
			}
			
			if col.Comment != "" {
				colDef += fmt.Sprintf(" COMMENT '%s'", escapeString(col.Comment))
			}
			
			columnDefs = append(columnDefs, colDef)
		}
		
		// 添加主键约束
		if len(table.PrimaryKeys) > 0 {
			pkCols := make([]string, len(table.PrimaryKeys))
			for i, pk := range table.PrimaryKeys {
				pkCols[i] = escapeIdentifier(pk)
			}
			columnDefs = append(columnDefs, fmt.Sprintf("  PRIMARY KEY (%s)", strings.Join(pkCols, ", ")))
		}
		
		sql.WriteString(strings.Join(columnDefs, ",\n"))
		sql.WriteString("\n);\n\n")
		
		// 添加表注释
		if table.Comment != "" {
			sql.WriteString(fmt.Sprintf("COMMENT ON TABLE %s IS '%s';\n\n", escapeIdentifier(table.Name), escapeString(table.Comment)))
		}
	}
	
	// 添加外键约束
	for _, rel := range drawdb.Relationships {
		if rel.FromTable != "" && rel.ToTable != "" {
			sql.WriteString(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s (%s);\n",
				escapeIdentifier(rel.FromTable),
				rel.FromTable,
				rel.ToTable,
				escapeIdentifier(rel.FromColumn),
				escapeIdentifier(rel.ToTable),
				escapeIdentifier(rel.ToColumn)))
		}
	}
	
	return sql.String(), nil
}

// mapTypeToSQL 映射类型到 SQL 类型
func mapTypeToSQL(typ, dbType string) string {
	typ = strings.ToUpper(typ)
	
	// PostgreSQL 类型映射
	if dbType == "postgresql" || dbType == "" {
		typeMap := map[string]string{
			"INTEGER": "INTEGER",
			"INT":     "INTEGER",
			"BIGINT":  "BIGINT",
			"SMALLINT": "SMALLINT",
			"VARCHAR": "VARCHAR(255)",
			"TEXT":    "TEXT",
			"BOOLEAN": "BOOLEAN",
			"DATE":    "DATE",
			"TIMESTAMP": "TIMESTAMP",
			"JSONB":   "JSONB",
			"UUID":   "UUID",
		}
		
		if mapped, ok := typeMap[typ]; ok {
			return mapped
		}
		return "TEXT"
	}
	
	// MySQL 类型映射
	if dbType == "mysql" {
		typeMap := map[string]string{
			"INTEGER": "INT",
			"INT":     "INT",
			"BIGINT":  "BIGINT",
			"SMALLINT": "SMALLINT",
			"VARCHAR": "VARCHAR(255)",
			"TEXT":    "TEXT",
			"BOOLEAN": "BOOLEAN",
			"DATE":    "DATE",
			"TIMESTAMP": "TIMESTAMP",
			"JSONB":   "JSON",
			"UUID":   "CHAR(36)",
		}
		
		if mapped, ok := typeMap[typ]; ok {
			return mapped
		}
		return "TEXT"
	}
	
	return "TEXT"
}

// escapeIdentifier 转义标识符
func escapeIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// escapeValue 转义值
func escapeValue(value string) string {
	// 如果是数字或布尔值，直接返回
	if value == "true" || value == "false" || isNumeric(value) {
		return value
	}
	return escapeString(value)
}

// escapeString 转义字符串
func escapeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// isNumeric 检查是否为数字
func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			if r != '.' && r != '-' {
				return false
			}
		}
	}
	return true
}

