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

	"github.com/xgfone/go-structs/field"
	"github.com/xgfone/go-structs/handler"
	"github.com/xgfone/go-validation"
)

var DefaultRuleValidator RuleValidateFunc = validation.Validate

// RuleValidator is used to validate whether a value conforms with the rule.
type RuleValidateFunc func(value any, rule string) error

// ValidateStructFieldRunner returns a runner to validate
// whether a struct field value is valid, which is registered
// into DefaultReflector with the tag name "validate" by default.
//
// If ruleValidator is nil, use DefaultRuleValidator instead.
func ValidateStructFieldRunner(ruleValidator RuleValidateFunc) handler.Runner {
	return handler.FieldRunner(func(v reflect.Value, sf reflect.StructField, a any) (err error) {
		if ruleValidator == nil {
			err = DefaultRuleValidator(v.Interface(), a.(string))
		} else {
			err = ruleValidator(v.Interface(), a.(string))
		}

		if err != nil {
			err = fmt.Errorf("%s: %w", getStructFieldName(sf), err)
		}
		return
	})
}

func getStructFieldName(sf reflect.StructField) (name string) {
	name, _ = field.GetTag(sf, "json")
	if name == "" || name == "-" {
		name = sf.Name
	}
	return
}
