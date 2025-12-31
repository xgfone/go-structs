// Copyright 2025 xgfone
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

// Package trim provides a handler to trim the leading and trailing
// whitespace or the specified characters of the struct field.
package trim

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/xgfone/go-structs/handler"
	"github.com/xgfone/go-structs/handler/setter"
)

// TrimRunner returns a runner to trim the leading and trailing
// whitespace or the specified characters of the struct field.
func TrimRunner() handler.Runner {
	return setter.SetterRunner(trim)
}

func trim(_ any, _, fieldptr reflect.Value, sf reflect.StructField, a any) (err error) {
	value := fieldptr.Elem()
	if value.Kind() != reflect.String {
		panic(fmt.Errorf("%s: expect string, but got %T", sf.Name, value.Interface()))
	}

	var news string
	olds := value.String()
	if arg := a.(string); arg == "" {
		news = strings.TrimSpace(olds)
	} else {
		news = strings.Trim(olds, arg)
	}

	if olds != news {
		value.SetString(news)
	}
	return
}
