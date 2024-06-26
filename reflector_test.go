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

package structs

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/xgfone/go-structs/handler"
)

func ExampleReflector() {
	parseInt := func(s string) (interface{}, error) { return strconv.ParseInt(s, 10, 64) }
	compareInt := func(isMin bool) handler.Runner {
		return func(_ interface{}, _, v reflect.Value, t reflect.StructField, a interface{}) error {
			value := v.Interface().(int64)
			if isMin {
				if min := a.(int64); value < min {
					return fmt.Errorf("%s: the value %d is less then %d", t.Name, value, min)
				}
			} else {
				if max := a.(int64); value > max {
					return fmt.Errorf("%s: the value %d is greater then %d", t.Name, value, max)
				}
			}
			return nil
		}
	}

	sf := NewReflector()
	sf.Register("min", handler.New(parseInt, compareInt(true)))
	sf.Register("max", handler.New(parseInt, compareInt(false)))
	sf.Register("default", handler.SimpleRunner(func(v reflect.Value, s interface{}) error {
		if !v.IsZero() {
			return nil
		}

		i, err := strconv.ParseInt(s.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid default value '%s': %s", s, err)
		}

		v.SetInt(i)
		return nil
	}))
	sf.Register("datamask", handler.SimpleRunner(func(v reflect.Value, s interface{}) error {
		switch s.(string) {
		case "username":
			name := v.Interface().(string)
			if r, _ := utf8.DecodeRuneInString(name); r != utf8.RuneError {
				v.SetString(string(r) + "**")
			} else {
				return fmt.Errorf("the name is not utf8")
			}

		case "password":
			v.SetString("******")

		default:
			return fmt.Errorf("unknown datamust type '%v'", s)
		}

		return nil
	}))

	/// Example 1: Check and validate the request arguments
	type Request struct {
		Page     int64 `default:"1" min:"1"`
		PageSize int64 `default:"10" min:"10" max:"100"`
	}
	request := Request{Page: 2}
	if err := sf.Reflect(&request); err != nil {
		fmt.Printf("reflect failed: %v\n", err)
	} else {
		fmt.Printf("Request.Page: %d\n", request.Page)
		fmt.Printf("Request.PageSize: %d\n", request.PageSize)
	}

	/// Example 2: Mask the response result data
	type Person struct {
		Username string `datamask:"username"`
		Password string `datamask:"password"`

		request Request
		Request Request `reflect:"-"` // Stop to reflect struct recursively.
	}
	type Response struct {
		Persons []Person
	}
	response := Response{Persons: []Person{
		{Username: "谢谢", Password: "123456789"},
	}}
	if err := sf.Reflect(&response); err != nil {
		fmt.Printf("reflect failed: %v\n", err)
	} else {
		fmt.Printf("Response.Username: %s\n", response.Persons[0].Username)
		fmt.Printf("Response.Password: %s\n", response.Persons[0].Password)
		fmt.Printf("Response.Request.Page: %d\n", response.Persons[0].Request.Page)
		fmt.Printf("Response.Request.PageSize: %d\n", response.Persons[0].Request.PageSize)
		fmt.Printf("Response.request.Page: %d\n", response.Persons[0].request.Page)
		fmt.Printf("Response.request.PageSize: %d\n", response.Persons[0].request.PageSize)
	}

	// Output:
	// Request.Page: 2
	// Request.PageSize: 10
	// Response.Username: 谢**
	// Response.Password: ******
	// Response.Request.Page: 0
	// Response.Request.PageSize: 0
	// Response.request.Page: 0
	// Response.request.PageSize: 0
}
