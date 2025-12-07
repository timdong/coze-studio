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

package validation

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"nil", nil, true},
		{"empty string", "", true},
		{"whitespace string", "   ", true},
		{"non-empty string", "test", false},
		{"empty slice", []string{}, true},
		{"non-empty slice", []string{"test"}, false},
		{"empty map", map[string]string{}, true},
		{"non-empty map", map[string]string{"key": "value"}, false},
		{"zero int", 0, false},
		{"non-zero int", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmpty(tt.value)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFieldValidator(t *testing.T) {
	fv := NewFieldValidator("test_field").
		AddRule(NewRequiredRule()).
		AddRule(NewMinLengthRule(3))

	// 测试空值
	err := fv.Validate("")
	if err == nil {
		t.Error("expected error for empty value")
	}

	// 测试太短的值
	err = fv.Validate("ab")
	if err == nil {
		t.Error("expected error for value too short")
	}

	// 测试有效值
	err = fv.Validate("abc")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStructValidator(t *testing.T) {
	type TestStruct struct {
		Name  string
		Email string
	}

	sv := NewStructValidator()
	sv.Field("Name").
		AddRule(NewRequiredRule()).
		AddRule(NewMinLengthRule(3))
	sv.Field("Email").
		AddRule(NewRequiredRule()).
		AddRule(NewEmailRule())

	// 测试有效结构体
	valid := TestStruct{
		Name:  "John",
		Email: "john@example.com",
	}
	if err := sv.Validate(valid); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// 测试无效结构体
	invalid := TestStruct{
		Name:  "Jo",
		Email: "invalid-email",
	}
	if err := sv.Validate(invalid); err == nil {
		t.Error("expected error for invalid struct")
	}
}

func TestCompositeValidator(t *testing.T) {
	cv := NewCompositeValidator().
		Add(NewRequiredRule()).
		Add(NewMinLengthRule(3))

	// 测试空值
	err := cv.Validate("")
	if err == nil {
		t.Error("expected error for empty value")
	}

	// 测试太短的值
	err = cv.Validate("ab")
	if err == nil {
		t.Error("expected error for value too short")
	}

	// 测试有效值
	err = cv.Validate("abc")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

