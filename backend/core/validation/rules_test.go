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
	"fmt"
	"testing"
)

func TestRequiredRule(t *testing.T) {
	rule := NewRequiredRule()

	if err := rule.Validate(nil); err == nil {
		t.Error("expected error for nil value")
	}

	if err := rule.Validate(""); err == nil {
		t.Error("expected error for empty string")
	}

	if err := rule.Validate("test"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMinLengthRule(t *testing.T) {
	rule := NewMinLengthRule(3)

	if err := rule.Validate("ab"); err == nil {
		t.Error("expected error for string too short")
	}

	if err := rule.Validate("abc"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := rule.Validate(123); err == nil {
		t.Error("expected error for non-string value")
	}
}

func TestMaxLengthRule(t *testing.T) {
	rule := NewMaxLengthRule(5)

	if err := rule.Validate("abcdef"); err == nil {
		t.Error("expected error for string too long")
	}

	if err := rule.Validate("abc"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLengthRule(t *testing.T) {
	rule := NewLengthRule(3, 5)

	if err := rule.Validate("ab"); err == nil {
		t.Error("expected error for string too short")
	}

	if err := rule.Validate("abcdef"); err == nil {
		t.Error("expected error for string too long")
	}

	if err := rule.Validate("abc"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestEmailRule(t *testing.T) {
	rule := NewEmailRule()

	validEmails := []string{
		"test@example.com",
		"user.name@example.co.uk",
		"user+tag@example.com",
	}

	invalidEmails := []string{
		"invalid",
		"@example.com",
		"user@",
		"user@example",
	}

	for _, email := range validEmails {
		if err := rule.Validate(email); err != nil {
			t.Errorf("unexpected error for valid email %s: %v", email, err)
		}
	}

	for _, email := range invalidEmails {
		if err := rule.Validate(email); err == nil {
			t.Errorf("expected error for invalid email %s", email)
		}
	}
}

func TestMinRule(t *testing.T) {
	rule := NewMinRule(10)

	if err := rule.Validate(5); err == nil {
		t.Error("expected error for value less than min")
	}

	if err := rule.Validate(15); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := rule.Validate(10); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMaxRule(t *testing.T) {
	rule := NewMaxRule(10)

	if err := rule.Validate(15); err == nil {
		t.Error("expected error for value greater than max")
	}

	if err := rule.Validate(5); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := rule.Validate(10); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRangeRule(t *testing.T) {
	rule := NewRangeRule(5, 10)

	if err := rule.Validate(3); err == nil {
		t.Error("expected error for value less than min")
	}

	if err := rule.Validate(15); err == nil {
		t.Error("expected error for value greater than max")
	}

	if err := rule.Validate(7); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInRule(t *testing.T) {
	rule := NewInRule("a", "b", "c")

	if err := rule.Validate("d"); err == nil {
		t.Error("expected error for value not in list")
	}

	if err := rule.Validate("a"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNotInRule(t *testing.T) {
	rule := NewNotInRule("a", "b", "c")

	if err := rule.Validate("a"); err == nil {
		t.Error("expected error for value in list")
	}

	if err := rule.Validate("d"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPatternRule(t *testing.T) {
	rule := NewPatternRule(`^[A-Z][a-z]+$`, "must start with uppercase letter")

	if err := rule.Validate("hello"); err == nil {
		t.Error("expected error for invalid pattern")
	}

	if err := rule.Validate("Hello"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCustomRule(t *testing.T) {
	rule := NewCustomRule(func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("must be a string")
		}
		if len(str) < 5 {
			return fmt.Errorf("must be at least 5 characters")
		}
		return nil
	})

	if err := rule.Validate("abc"); err == nil {
		t.Error("expected error for value too short")
	}

	if err := rule.Validate("abcdef"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

