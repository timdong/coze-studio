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
	"context"
	"fmt"
	"reflect"
	"strings"
)

// Validator 验证器接口
type Validator interface {
	Validate(value interface{}) error
}

// FieldValidator 字段验证器
type FieldValidator struct {
	fieldName string
	rules     []Validator
}

// NewFieldValidator 创建字段验证器
func NewFieldValidator(fieldName string) *FieldValidator {
	return &FieldValidator{
		fieldName: fieldName,
		rules:     make([]Validator, 0),
	}
}

// AddRule 添加验证规则
func (fv *FieldValidator) AddRule(rule Validator) *FieldValidator {
	fv.rules = append(fv.rules, rule)
	return fv
}

// Validate 验证字段
func (fv *FieldValidator) Validate(value interface{}) error {
	for _, rule := range fv.rules {
		if err := rule.Validate(value); err != nil {
			return fmt.Errorf("%s: %w", fv.fieldName, err)
		}
	}
	return nil
}

// StructValidator 结构体验证器
type StructValidator struct {
	validators map[string]*FieldValidator
}

// NewStructValidator 创建结构体验证器
func NewStructValidator() *StructValidator {
	return &StructValidator{
		validators: make(map[string]*FieldValidator),
	}
}

// Field 添加字段验证器
func (sv *StructValidator) Field(fieldName string) *FieldValidator {
	if fv, exists := sv.validators[fieldName]; exists {
		return fv
	}
	fv := NewFieldValidator(fieldName)
	sv.validators[fieldName] = fv
	return fv
}

// Validate 验证结构体
func (sv *StructValidator) Validate(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("data must be a struct or pointer to struct")
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i).Interface()

		// 检查是否有该字段的验证器
		if validator, exists := sv.validators[field.Name]; exists {
			if err := validator.Validate(fieldValue); err != nil {
				return err
			}
		}
	}

	return nil
}

// ContextValidator 上下文验证器接口
type ContextValidator interface {
	Validate(ctx context.Context, value interface{}) error
}

// CompositeValidator 组合验证器
type CompositeValidator struct {
	validators []Validator
}

// NewCompositeValidator 创建组合验证器
func NewCompositeValidator() *CompositeValidator {
	return &CompositeValidator{
		validators: make([]Validator, 0),
	}
}

// Add 添加验证器
func (cv *CompositeValidator) Add(validator Validator) *CompositeValidator {
	cv.validators = append(cv.validators, validator)
	return cv
}

// Validate 执行所有验证
func (cv *CompositeValidator) Validate(value interface{}) error {
	var errors []error
	for _, validator := range cv.validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %v", errors)
	}
	return nil
}

// IsEmpty 检查值是否为空
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return strings.TrimSpace(v.String()) == ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

