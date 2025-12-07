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
	"encoding/json"
	"fmt"
)

// DrawDBFormat DrawDB 格式结构
type DrawDBFormat struct {
	Tables       []Table       `json:"tables"`
	Relationships []Relationship `json:"relationships"`
}

// Table 表结构
type Table struct {
	Name        string   `json:"name"`
	Comment     string   `json:"comment,omitempty"`
	Columns     []Column `json:"columns"`
	PrimaryKeys []string `json:"primary_keys,omitempty"`
	Indexes     []Index  `json:"indexes,omitempty"`
}

// Column 列结构
type Column struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Nullable     bool   `json:"nullable"`
	DefaultValue string `json:"default_value,omitempty"`
	Comment      string `json:"comment,omitempty"`
}

// Index 索引结构
type Index struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
}

// Relationship 关系结构
type Relationship struct {
	FromTable  string `json:"from_table"`
	FromColumn string `json:"from_column"`
	ToTable    string `json:"to_table"`
	ToColumn   string `json:"to_column"`
	Type       string `json:"type"` // one-to-one, one-to-many, many-to-many
}

// ParseDrawDB 解析 DrawDB JSON 字符串
func ParseDrawDB(data string) (*DrawDBFormat, error) {
	var drawdb DrawDBFormat
	if err := json.Unmarshal([]byte(data), &drawdb); err != nil {
		return nil, fmt.Errorf("failed to parse DrawDB format: %w", err)
	}
	return &drawdb, nil
}

// ToJSON 转换为 JSON 字符串
func (d *DrawDBFormat) ToJSON() (string, error) {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal DrawDB format: %w", err)
	}
	return string(data), nil
}

// NewDrawDB 创建新的 DrawDB 格式
func NewDrawDB() *DrawDBFormat {
	return &DrawDBFormat{
		Tables:       []Table{},
		Relationships: []Relationship{},
	}
}

