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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestSuccess(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		Success(c, map[string]string{"key": "value"}, "test message")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreated(t *testing.T) {
	router := setupRouter()
	router.POST("/test", func(c *gin.Context) {
		Created(c, map[string]string{"id": "1"}, "created")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUpdated(t *testing.T) {
	router := setupRouter()
	router.PUT("/test", func(c *gin.Context) {
		Updated(c, map[string]string{"id": "1"}, "updated")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleted(t *testing.T) {
	router := setupRouter()
	router.DELETE("/test", func(c *gin.Context) {
		Deleted(c, "deleted")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBadRequest(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		BadRequest(c, "bad request", map[string]interface{}{"field": "error"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnauthorized(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		Unauthorized(c, "unauthorized")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestForbidden(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		Forbidden(c, "forbidden")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestNotFound(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		NotFound(c, "not found")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestValidationFailed(t *testing.T) {
	router := setupRouter()
	router.POST("/test", func(c *gin.Context) {
		ValidationFailed(c, "validation failed", map[string]interface{}{"field": "error"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestInternalError(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		InternalError(c, "internal error", assert.AnError)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestErrorResponse(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		err := NewError(ErrorCodeBadRequest, "test error", http.StatusBadRequest)
		ErrorResponse(c, err)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleError(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		err := NewError(ErrorCodeNotFound, "not found", http.StatusNotFound)
		HandleError(c, err)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSuccessWithPagination(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		items := []string{"item1", "item2"}
		pagination := &service.PaginationResponse{
			Page:       1,
			PageSize:   20,
			Total:      100,
			TotalPages: 5,
		}
		SuccessWithPagination(c, items, pagination, "success")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

