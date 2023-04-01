// Copyright 2023 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package validate provides a handler to validate the struct field.
package validate

import (
	"fmt"
	"reflect"

	"github.com/xgfone/go-defaults"
	"github.com/xgfone/go-defaults/assists"
	"github.com/xgfone/go-structs/handler"
)

// NewValidatorHandler returns a handler to validate whether the field value
// is valid, which is registered into DefaultReflector
// with the tag name "validate" by default.
//
// Because the reflector walks the sub-fields of the struct slice field
// recursively, so the validation rule "array(structure)" should not be used.
//
// If ruleValidator is nil, use defaults.RuleValidator instead.
func NewValidatorHandler(ruleValidator assists.RuleValidator) handler.Handler {
	return validator{ruleValidator}
}

type validator struct {
	assists.RuleValidator
}

func (h validator) Parse(s string) (interface{}, error) { return s, nil }
func (h validator) Run(c interface{}, r, v reflect.Value, t reflect.StructField, a interface{}) (err error) {
	if h.RuleValidator == nil {
		err = defaults.ValidateWithRule(v.Interface(), a.(string))
	} else {
		err = h.RuleValidator.Validate(v.Interface(), a.(string))
	}

	if err != nil {
		err = fmt.Errorf("%s: %w", getStructFieldName(t), err)
	}
	return err
}

func getStructFieldName(sf reflect.StructField) (name string) {
	name, _ = defaults.GetStructFieldName(sf)
	if name == "" {
		name = sf.Name
	}
	return
}
