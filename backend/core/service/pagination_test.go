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

package service

import (
	"testing"
)

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		pageSize       int
		expectedPage   int
		expectedSize   int
		expectedOffset int
	}{
		{
			name:           "normal case",
			page:           2,
			pageSize:       20,
			expectedPage:   2,
			expectedSize:   20,
			expectedOffset: 20,
		},
		{
			name:           "page less than 1",
			page:           0,
			pageSize:       20,
			expectedPage:   1,
			expectedSize:   20,
			expectedOffset: 0,
		},
		{
			name:           "pageSize less than 1",
			page:           1,
			pageSize:       0,
			expectedPage:   1,
			expectedSize:   20,
			expectedOffset: 0,
		},
		{
			name:           "pageSize greater than 100",
			page:           1,
			pageSize:       200,
			expectedPage:   1,
			expectedSize:   100,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, size, offset := NormalizePagination(tt.page, tt.pageSize)
			if page != tt.expectedPage {
				t.Errorf("expected page %d, got %d", tt.expectedPage, page)
			}
			if size != tt.expectedSize {
				t.Errorf("expected size %d, got %d", tt.expectedSize, size)
			}
			if offset != tt.expectedOffset {
				t.Errorf("expected offset %d, got %d", tt.expectedOffset, offset)
			}
		})
	}
}

func TestCalculateTotalPages(t *testing.T) {
	tests := []struct {
		name         string
		total        int
		pageSize     int
		expectedPages int
	}{
		{
			name:          "normal case",
			total:         100,
			pageSize:      20,
			expectedPages: 5,
		},
		{
			name:          "total is 0",
			total:         0,
			pageSize:      20,
			expectedPages: 0,
		},
		{
			name:          "remainder exists",
			total:         101,
			pageSize:      20,
			expectedPages: 6,
		},
		{
			name:          "total less than pageSize",
			total:         10,
			pageSize:      20,
			expectedPages: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages := CalculateTotalPages(tt.total, tt.pageSize)
			if pages != tt.expectedPages {
				t.Errorf("expected pages %d, got %d", tt.expectedPages, pages)
			}
		})
	}
}

func TestNewPaginationResponse(t *testing.T) {
	page := 2
	pageSize := 20
	total := 100

	response := NewPaginationResponse(page, pageSize, total)

	if response.Page != page {
		t.Errorf("expected page %d, got %d", page, response.Page)
	}
	if response.PageSize != pageSize {
		t.Errorf("expected pageSize %d, got %d", pageSize, response.PageSize)
	}
	if response.Total != total {
		t.Errorf("expected total %d, got %d", total, response.Total)
	}
	if response.TotalPages != 5 {
		t.Errorf("expected totalPages 5, got %d", response.TotalPages)
	}
}

