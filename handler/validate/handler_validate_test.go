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

package validate_test

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/xgfone/go-defaults/assists"
	"github.com/xgfone/go-structs"
	"github.com/xgfone/go-structs/handler/validate"
)

func ExampleNewValidatorHandler() {
	validator := func(value interface{}, rule string) error {
		value = reflect.Indirect(reflect.ValueOf(value)).Interface()
		switch {
		case strings.HasPrefix(rule, "min(") && strings.HasSuffix(rule, ")"):
			min, err := strconv.ParseInt(rule[len("min("):len(rule)-1], 10, 64)
			if err != nil {
				panic(err)
			} else if v := value.(int64); v < min {
				return fmt.Errorf("the integer %d is less than %d", v, min)
			}

		case strings.HasPrefix(rule, "max(") && strings.HasSuffix(rule, ")"):
			max, err := strconv.ParseInt(rule[len("max("):len(rule)-1], 10, 64)
			if err != nil {
				panic(err)
			} else if v := value.(int64); v > max {
				return fmt.Errorf("the integer %d is greater than %d", v, max)
			}

		default:
			panic(fmt.Errorf("unknown validate run '%s'", rule))
		}

		return nil
	}
	structs.Register("validate", validate.NewValidatorHandler(assists.RuleValidateFunc(validator)))

	type S struct {
		F1 int64  `validate:"min(100)"` // General Type
		F2 *int64 `validate:"max(200)"` // Pointer Type
		F3 []struct {
			F4 int64 `validate:"min(1)"`
			F5 int64
		}
	}

	var v S
	fmt.Println(structs.Reflect(v))  // Only returns the first error
	fmt.Println(structs.Reflect(&v)) // Only returns the first error

	v.F1 = 123
	v.F2 = new(int64)
	*v.F2 = 123
	v.F3 = []struct {
		F4 int64 `validate:"min(1)"`
		F5 int64
	}{
		{F4: 1},
		{F4: 2},
	}
	fmt.Println(structs.Reflect(v))
	fmt.Println(structs.Reflect(&v))

	// Output:
	// F1: the integer 0 is less than 100
	// F1: the integer 0 is less than 100
	// <nil>
	// <nil>
}
