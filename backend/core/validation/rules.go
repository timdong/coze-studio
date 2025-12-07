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
	"reflect"
	"regexp"
	"strings"
)

// RequiredRule 必填规则
type RequiredRule struct{}

// NewRequiredRule 创建必填规则
func NewRequiredRule() *RequiredRule {
	return &RequiredRule{}
}

func (r *RequiredRule) Validate(value interface{}) error {
	if IsEmpty(value) {
		return fmt.Errorf("field is required")
	}
	return nil
}

// MinLengthRule 最小长度规则
type MinLengthRule struct {
	Min int
}

// NewMinLengthRule 创建最小长度规则
func NewMinLengthRule(min int) *MinLengthRule {
	return &MinLengthRule{Min: min}
}

func (r *MinLengthRule) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}
	if len(str) < r.Min {
		return fmt.Errorf("field must be at least %d characters", r.Min)
	}
	return nil
}

// MaxLengthRule 最大长度规则
type MaxLengthRule struct {
	Max int
}

// NewMaxLengthRule 创建最大长度规则
func NewMaxLengthRule(max int) *MaxLengthRule {
	return &MaxLengthRule{Max: max}
}

func (r *MaxLengthRule) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}
	if len(str) > r.Max {
		return fmt.Errorf("field must be at most %d characters", r.Max)
	}
	return nil
}

// LengthRule 长度范围规则
type LengthRule struct {
	Min int
	Max int
}

// NewLengthRule 创建长度范围规则
func NewLengthRule(min, max int) *LengthRule {
	return &LengthRule{Min: min, Max: max}
}

func (r *LengthRule) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}
	length := len(str)
	if length < r.Min || length > r.Max {
		return fmt.Errorf("field must be between %d and %d characters", r.Min, r.Max)
	}
	return nil
}

// EmailRule 邮箱规则
type EmailRule struct {
	pattern *regexp.Regexp
}

// NewEmailRule 创建邮箱规则
func NewEmailRule() *EmailRule {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return &EmailRule{pattern: pattern}
}

func (r *EmailRule) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}
	if !r.pattern.MatchString(str) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// URLRule URL规则
type URLRule struct {
	pattern *regexp.Regexp
}

// NewURLRule 创建URL规则
func NewURLRule() *URLRule {
	pattern := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return &URLRule{pattern: pattern}
}

func (r *URLRule) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}
	if !r.pattern.MatchString(str) {
		return fmt.Errorf("invalid URL format")
	}
	return nil
}

// MinRule 最小值规则
type MinRule struct {
	Min float64
}

// NewMinRule 创建最小值规则
func NewMinRule(min float64) *MinRule {
	return &MinRule{Min: min}
}

func (r *MinRule) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	var num float64

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num = float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		num = v.Float()
	default:
		return fmt.Errorf("value must be a number")
	}

	if num < r.Min {
		return fmt.Errorf("value must be at least %v", r.Min)
	}
	return nil
}

// MaxRule 最大值规则
type MaxRule struct {
	Max float64
}

// NewMaxRule 创建最大值规则
func NewMaxRule(max float64) *MaxRule {
	return &MaxRule{Max: max}
}

func (r *MaxRule) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	var num float64

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num = float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		num = v.Float()
	default:
		return fmt.Errorf("value must be a number")
	}

	if num > r.Max {
		return fmt.Errorf("value must be at most %v", r.Max)
	}
	return nil
}

// RangeRule 范围规则
type RangeRule struct {
	Min float64
	Max float64
}

// NewRangeRule 创建范围规则
func NewRangeRule(min, max float64) *RangeRule {
	return &RangeRule{Min: min, Max: max}
}

func (r *RangeRule) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	var num float64

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num = float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		num = v.Float()
	default:
		return fmt.Errorf("value must be a number")
	}

	if num < r.Min || num > r.Max {
		return fmt.Errorf("value must be between %v and %v", r.Min, r.Max)
	}
	return nil
}

// InRule 在列表中规则
type InRule struct {
	Values []interface{}
}

// NewInRule 创建在列表中规则
func NewInRule(values ...interface{}) *InRule {
	return &InRule{Values: values}
}

func (r *InRule) Validate(value interface{}) error {
	for _, v := range r.Values {
		if reflect.DeepEqual(value, v) {
			return nil
		}
	}
	return fmt.Errorf("value must be one of: %v", r.Values)
}

// NotInRule 不在列表中规则
type NotInRule struct {
	Values []interface{}
}

// NewNotInRule 创建不在列表中规则
func NewNotInRule(values ...interface{}) *NotInRule {
	return &NotInRule{Values: values}
}

func (r *NotInRule) Validate(value interface{}) error {
	for _, v := range r.Values {
		if reflect.DeepEqual(value, v) {
			return fmt.Errorf("value must not be one of: %v", r.Values)
		}
	}
	return nil
}

// PatternRule 正则表达式规则
type PatternRule struct {
	Pattern *regexp.Regexp
	Message string
}

// NewPatternRule 创建正则表达式规则
func NewPatternRule(pattern string, message string) *PatternRule {
	regex := regexp.MustCompile(pattern)
	return &PatternRule{
		Pattern: regex,
		Message: message,
	}
}

func (r *PatternRule) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}
	if !r.Pattern.MatchString(str) {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("value does not match required pattern")
	}
	return nil
}

// CustomRule 自定义规则
type CustomRule struct {
	ValidateFunc func(interface{}) error
}

// NewCustomRule 创建自定义规则
func NewCustomRule(validateFunc func(interface{}) error) *CustomRule {
	return &CustomRule{ValidateFunc: validateFunc}
}

func (r *CustomRule) Validate(value interface{}) error {
	return r.ValidateFunc(value)
}

// TrimmedStringRule 去除空白字符后验证字符串规则
type TrimmedStringRule struct {
	Rule Validator
}

// NewTrimmedStringRule 创建去除空白字符后验证字符串规则
func NewTrimmedStringRule(rule Validator) *TrimmedStringRule {
	return &TrimmedStringRule{Rule: rule}
}

func (r *TrimmedStringRule) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}
	trimmed := strings.TrimSpace(str)
	return r.Rule.Validate(trimmed)
}

