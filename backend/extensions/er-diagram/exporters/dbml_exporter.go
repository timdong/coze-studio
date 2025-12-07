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

// ExportToDBML 将 DrawDB 格式导出为 DBML
func ExportToDBML(drawdb *DrawDBFormat) (string, error) {
	var dbml strings.Builder
	
	for _, table := range drawdb.Tables {
		dbml.WriteString(fmt.Sprintf("Table %s {\n", table.Name))
		
		// 写入列定义
		for _, col := range table.Columns {
			colDef := fmt.Sprintf("  %s %s", col.Name, col.Type)
			
			// 添加主键标记
			for _, pk := range table.PrimaryKeys {
				if pk == col.Name {
					colDef += " [pk]"
					break
				}
			}
			
			if !col.Nullable {
				colDef += " [not null]"
			}
			
			if col.DefaultValue != "" {
				colDef += fmt.Sprintf(" [default: '%s']", col.DefaultValue)
			}
			
			if col.Comment != "" {
				colDef += fmt.Sprintf(" [note: '%s']", col.Comment)
			}
			
			dbml.WriteString(colDef + "\n")
		}
		
		// 写入表注释
		if table.Comment != "" {
			dbml.WriteString(fmt.Sprintf("  Note: '%s'\n", table.Comment))
		}
		
		dbml.WriteString("}\n\n")
	}
	
	// 写入关系定义
	for _, rel := range drawdb.Relationships {
		if rel.FromTable != "" && rel.ToTable != "" {
			arrow := ">"
			if rel.Type == "many-to-one" {
				arrow = "<"
			}
			dbml.WriteString(fmt.Sprintf("Ref: %s.%s %s %s.%s\n",
				rel.FromTable,
				rel.FromColumn,
				arrow,
				rel.ToTable,
				rel.ToColumn))
		}
	}
	
	return dbml.String(), nil
}

