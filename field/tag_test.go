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

package field

import (
	"fmt"
	"reflect"
)

func ExampleGetTag() {
	type T struct {
		Int   int
		Int8  int8  `json:""`
		Int16 int16 `json:"int16"`
		Int32 int32 `json:",arg"`
		Int64 int64 `json:"int64,arg"`
		Bool  bool  `json:"-"`
	}

	stype := reflect.TypeOf(T{})
	for i := 0; i < stype.NumField(); i++ {
		field := stype.Field(i)
		value, arg := GetTag(field, "json")
		fmt.Printf("fieldname=%s, tagvalue=%s, tagarg=%s\n", field.Name, value, arg)
	}

	// Output:
	// fieldname=Int, tagvalue=, tagarg=
	// fieldname=Int8, tagvalue=, tagarg=
	// fieldname=Int16, tagvalue=int16, tagarg=
	// fieldname=Int32, tagvalue=, tagarg=arg
	// fieldname=Int64, tagvalue=int64, tagarg=arg
	// fieldname=Bool, tagvalue=-, tagarg=
	//
}
