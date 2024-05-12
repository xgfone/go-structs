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

// Package setter provides a handler to set the struct field to a value.
package setter

import (
	"fmt"
	"reflect"

	"github.com/xgfone/go-structs/handler"
)

// SetFormatRunner is the same as SetterRunner, but asserts the struct field
// to one of the interfaces as follow:
//
//	interface { SetFormat(string) }
//	interface { SetFormat(string) error }
func SetFormatRunner() handler.Runner {
	return SetterRunner(handler.SimpleRunner(func(vf reflect.Value, arg any) (err error) {
		switch i := vf.Interface().(type) {
		case interface{ SetFormat(string) }:
			i.SetFormat(arg.(string))
		case interface{ SetFormat(string) error }:
			err = i.SetFormat(arg.(string))
		default:
			panic(fmt.Errorf("%T has not implemented the interface { SetFormat(string) } or { SetFormat(string) error }", i))
		}
		return
	}))
}

// SetterRunner returns a runner to set the struct field to something
// by the set function.
//
// If set is nil, use the default, which will assert the struct field
// to the interface { Set(any) error }.
func SetterRunner(setter handler.Runner) handler.Runner {
	if setter == nil {
		setter = setbyiface
	}

	return func(c any, r, vf reflect.Value, sf reflect.StructField, arg any) error {
		if !vf.CanSet() {
			return fmt.Errorf("the field '%s' cannnot be set", sf.Name)
		}

		ptr := vf
		if vf.Kind() != reflect.Pointer {
			ptr = vf.Addr()
		} else if vf.IsNil() {
			vf.Set(reflect.New(vf.Type().Elem()))
		}

		return setter(c, r, ptr, sf, arg)
	}
}

func setbyiface(_ any, _, fieldptr reflect.Value, sf reflect.StructField, arg any) error {
	if setter, ok := fieldptr.Interface().(interface{ Set(any) error }); ok {
		return setter.Set(arg)
	}
	panic(fmt.Errorf("%s(%T) has not implemented the interface setter.Setter", sf.Name, fieldptr.Interface()))
}
