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

package setter_test

import (
	"fmt"
	"strconv"

	"github.com/xgfone/go-structs"
	"github.com/xgfone/go-structs/handler/setter"
)

type _Int int

func (i *_Int) Set(v any) error {
	_v, err := strconv.ParseInt(v.(string), 10, 64)
	if err == nil {
		*i = _Int(_v)
	}
	return err
}

type _Str string

func (s *_Str) Set(v any) error {
	*s = _Str(v.(string))
	return nil
}

func (s *_Str) SetFormat(v any) error {
	return s.Set(v)
}

func ExampleSetterRunner() {
	// "set" is registered by default. Now, we register the customized
	// "setint" to pre-parse the tag value to int.
	structs.Register("setint", setter.SetterRunner(nil))

	var t struct {
		Str _Str `set:"abc"`
		Int _Int `setint:"123"`
	}

	if err := structs.Reflect(&t); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Int: %v\n", t.Int)
		fmt.Printf("Str: %v\n", t.Str)
	}

	// Output:
	// Int: 123
	// Str: abc
}

func ExampleSetFormatRunner() {
	var t struct {
		Str _Str `set:"abc"`
	}

	if err := structs.Reflect(&t); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Str: %v\n", t.Str)
	}

	// Output:
	// Str: abc
}
