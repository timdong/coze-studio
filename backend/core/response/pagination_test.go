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

package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coze-dev/coze-studio/backend/core/service"
	"etrxlite/tools/search"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccessPaginated(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		items := []string{"item1", "item2"}
		pagination := &service.PaginationResponse{
			Page:       1,
			PageSize:   20,
			Total:      100,
			TotalPages: 5,
		}
		SuccessPaginated(c, items, pagination, "success")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSuccessPaginatedFromSearch(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		result := &search.SearchResult{
			Page:       2,
			PageSize:   20,
			Total:      100,
			TotalPages: 5,
		}
		items := []string{"item1", "item2"}
		SuccessPaginatedFromSearch(c, result, items, "success")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSuccessPaginatedFromList(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		items := []string{"item1", "item2"}
		SuccessPaginatedFromList(c, items, 1, 20, 100, "success")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSuccessPaginatedWithOffset(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		items := []string{"item1", "item2"}
		SuccessPaginatedWithOffset(c, items, 20, 20, 100, "success")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

