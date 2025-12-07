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

package event

import (
	"context"
	"errors"
	"testing"
)

func TestEventBus_Subscribe(t *testing.T) {
	bus := NewEventBus()
	handler := func(ctx context.Context, event *Event) error {
		return nil
	}

	bus.Subscribe(EventResourceCreated, handler)

	count := bus.GetHandlerCount(EventResourceCreated)
	if count != 1 {
		t.Errorf("expected 1 handler, got %d", count)
	}
}

func TestEventBus_Publish(t *testing.T) {
	bus := NewEventBus()
	called := false

	handler := func(ctx context.Context, event *Event) error {
		called = true
		return nil
	}

	bus.Subscribe(EventResourceCreated, handler)

	event := NewEvent(EventResourceCreated, "test-service", map[string]string{"id": "123"})
	err := bus.Publish(context.Background(), event)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !called {
		t.Error("expected handler to be called")
	}
}

func TestEventBus_Publish_NoSubscribers(t *testing.T) {
	bus := NewEventBus()

	event := NewEvent(EventResourceCreated, "test-service", map[string]string{"id": "123"})
	err := bus.Publish(context.Background(), event)

	if err != nil {
		t.Errorf("expected no error when no subscribers, got %v", err)
	}
}

func TestEventBus_Publish_HandlerError(t *testing.T) {
	bus := NewEventBus()
	testErr := errors.New("test error")

	handler := func(ctx context.Context, event *Event) error {
		return testErr
	}

	bus.Subscribe(EventResourceCreated, handler)

	event := NewEvent(EventResourceCreated, "test-service", map[string]string{"id": "123"})
	err := bus.Publish(context.Background(), event)

	if err == nil {
		t.Error("expected error from handler")
	}
}

func TestEventBus_Publish_NilEvent(t *testing.T) {
	bus := NewEventBus()

	err := bus.Publish(context.Background(), nil)
	if err == nil {
		t.Error("expected error for nil event")
	}
}

func TestEventBus_Unsubscribe(t *testing.T) {
	bus := NewEventBus()
	handler := func(ctx context.Context, event *Event) error {
		return nil
	}

	bus.Subscribe(EventResourceCreated, handler)
	bus.Unsubscribe(EventResourceCreated)

	count := bus.GetHandlerCount(EventResourceCreated)
	if count != 0 {
		t.Errorf("expected 0 handlers after unsubscribe, got %d", count)
	}
}

func TestEventBus_ListSubscribedEvents(t *testing.T) {
	bus := NewEventBus()
	handler := func(ctx context.Context, event *Event) error {
		return nil
	}

	bus.Subscribe(EventResourceCreated, handler)
	bus.Subscribe(EventResourceUpdated, handler)

	events := bus.ListSubscribedEvents()
	if len(events) != 2 {
		t.Errorf("expected 2 subscribed events, got %d", len(events))
	}
}

func TestEvent_WithUser(t *testing.T) {
	event := NewEvent(EventResourceCreated, "test-service", nil)
	event.WithUser(123, "testuser")

	if event.UserID != 123 {
		t.Errorf("expected user ID 123, got %d", event.UserID)
	}
	if event.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", event.Username)
	}
}

func TestEvent_WithWorkspace(t *testing.T) {
	event := NewEvent(EventResourceCreated, "test-service", nil)
	event.WithWorkspace(456)

	if event.WorkspaceID != 456 {
		t.Errorf("expected workspace ID 456, got %d", event.WorkspaceID)
	}
}

func TestEvent_WithMetadata(t *testing.T) {
	event := NewEvent(EventResourceCreated, "test-service", nil)
	event.WithMetadata("key1", "value1")
	event.WithMetadata("key2", "value2")

	if event.Metadata["key1"] != "value1" {
		t.Errorf("expected metadata key1='value1', got '%v'", event.Metadata["key1"])
	}
	if event.Metadata["key2"] != "value2" {
		t.Errorf("expected metadata key2='value2', got '%v'", event.Metadata["key2"])
	}
}

