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

// ValidateStructFieldRunner returns a runner to validate
// whether a struct field value is valid, which is registered
// into DefaultReflector with the tag name "validate" by default.
//
// If ruleValidator is nil, use defaults.RuleValidator instead.
func ValidateStructFieldRunner(ruleValidator assists.RuleValidator) handler.Runner {
	return handler.FieldRunner(func(v reflect.Value, sf reflect.StructField, a any) (err error) {
		if ruleValidator == nil {
			err = defaults.ValidateWithRule(v.Interface(), a.(string))
		} else {
			err = ruleValidator.Validate(v.Interface(), a.(string))
		}

		if err != nil {
			err = fmt.Errorf("%s: %w", getStructFieldName(sf), err)
		}
		return
	})
}

func getStructFieldName(sf reflect.StructField) (name string) {
	name, _ = defaults.GetStructFieldName(sf)
	if name == "" {
		name = sf.Name
	}
	return
}
