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

package ctxutil

import (
	"context"
	"testing"
)

func TestWithUserID(t *testing.T) {
	ctx := context.Background()
	userID := uint(123)

	ctx = WithUserID(ctx, userID)

	retrievedID, ok := GetUserID(ctx)
	if !ok {
		t.Error("expected user ID to be set")
	}
	if retrievedID != userID {
		t.Errorf("expected user ID %d, got %d", userID, retrievedID)
	}
}

func TestGetUserID_NotSet(t *testing.T) {
	ctx := context.Background()

	_, ok := GetUserID(ctx)
	if ok {
		t.Error("expected user ID not to be set")
	}
}

func TestMustGetUserID(t *testing.T) {
	ctx := context.Background()
	userID := uint(456)

	ctx = WithUserID(ctx, userID)

	retrievedID := MustGetUserID(ctx)
	if retrievedID != userID {
		t.Errorf("expected user ID %d, got %d", userID, retrievedID)
	}
}

func TestMustGetUserID_NotSet(t *testing.T) {
	ctx := context.Background()

	retrievedID := MustGetUserID(ctx)
	if retrievedID != 0 {
		t.Errorf("expected user ID 0, got %d", retrievedID)
	}
}

func TestWithUsername(t *testing.T) {
	ctx := context.Background()
	username := "testuser"

	ctx = WithUsername(ctx, username)

	retrievedUsername, ok := GetUsername(ctx)
	if !ok {
		t.Error("expected username to be set")
	}
	if retrievedUsername != username {
		t.Errorf("expected username %s, got %s", username, retrievedUsername)
	}
}

func TestGetUsername_NotSet(t *testing.T) {
	ctx := context.Background()

	_, ok := GetUsername(ctx)
	if ok {
		t.Error("expected username not to be set")
	}
}

func TestMustGetUsername(t *testing.T) {
	ctx := context.Background()
	username := "testuser2"

	ctx = WithUsername(ctx, username)

	retrievedUsername := MustGetUsername(ctx)
	if retrievedUsername != username {
		t.Errorf("expected username %s, got %s", username, retrievedUsername)
	}
}

func TestWithWorkspaceID(t *testing.T) {
	ctx := context.Background()
	workspaceID := uint(789)

	ctx = WithWorkspaceID(ctx, workspaceID)

	retrievedID, ok := GetWorkspaceID(ctx)
	if !ok {
		t.Error("expected workspace ID to be set")
	}
	if retrievedID != workspaceID {
		t.Errorf("expected workspace ID %d, got %d", workspaceID, retrievedID)
	}
}

func TestGetWorkspaceID_NotSet(t *testing.T) {
	ctx := context.Background()

	_, ok := GetWorkspaceID(ctx)
	if ok {
		t.Error("expected workspace ID not to be set")
	}
}

func TestMustGetWorkspaceID(t *testing.T) {
	ctx := context.Background()
	workspaceID := uint(101112)

	ctx = WithWorkspaceID(ctx, workspaceID)

	retrievedID := MustGetWorkspaceID(ctx)
	if retrievedID != workspaceID {
		t.Errorf("expected workspace ID %d, got %d", workspaceID, retrievedID)
	}
}

func TestWithRequestID(t *testing.T) {
	ctx := context.Background()
	requestID := "req-123"

	ctx = WithRequestID(ctx, requestID)

	retrievedID, ok := GetRequestID(ctx)
	if !ok {
		t.Error("expected request ID to be set")
	}
	if retrievedID != requestID {
		t.Errorf("expected request ID %s, got %s", requestID, retrievedID)
	}
}

func TestGetRequestID_NotSet(t *testing.T) {
	ctx := context.Background()

	_, ok := GetRequestID(ctx)
	if ok {
		t.Error("expected request ID not to be set")
	}
}

func TestMustGetRequestID(t *testing.T) {
	ctx := context.Background()
	requestID := "req-456"

	ctx = WithRequestID(ctx, requestID)

	retrievedID := MustGetRequestID(ctx)
	if retrievedID != requestID {
		t.Errorf("expected request ID %s, got %s", requestID, retrievedID)
	}
}

func TestFromGinContext(t *testing.T) {
	ctx := context.Background()
	userID := 123
	username := "testuser"
	workspaceID := 456

	ctx = FromGinContext(ctx, userID, username, workspaceID)

	retrievedUserID, ok := GetUserID(ctx)
	if !ok {
		t.Error("expected user ID to be set")
	}
	if retrievedUserID != uint(userID) {
		t.Errorf("expected user ID %d, got %d", userID, retrievedUserID)
	}

	retrievedUsername, ok := GetUsername(ctx)
	if !ok {
		t.Error("expected username to be set")
	}
	if retrievedUsername != username {
		t.Errorf("expected username %s, got %s", username, retrievedUsername)
	}

	retrievedWorkspaceID, ok := GetWorkspaceID(ctx)
	if !ok {
		t.Error("expected workspace ID to be set")
	}
	if retrievedWorkspaceID != uint(workspaceID) {
		t.Errorf("expected workspace ID %d, got %d", workspaceID, retrievedWorkspaceID)
	}
}

func TestFromGinContext_WithUint(t *testing.T) {
	ctx := context.Background()
	userID := uint(789)
	workspaceID := uint(101112)

	ctx = FromGinContext(ctx, userID, nil, workspaceID)

	retrievedUserID, ok := GetUserID(ctx)
	if !ok {
		t.Error("expected user ID to be set")
	}
	if retrievedUserID != userID {
		t.Errorf("expected user ID %d, got %d", userID, retrievedUserID)
	}

	retrievedWorkspaceID, ok := GetWorkspaceID(ctx)
	if !ok {
		t.Error("expected workspace ID to be set")
	}
	if retrievedWorkspaceID != workspaceID {
		t.Errorf("expected workspace ID %d, got %d", workspaceID, retrievedWorkspaceID)
	}
}

