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
	"reflect"
	"strings"
)

// GetTag returns the value and arg about the tag from the struct field.
func GetTag(sf reflect.StructField, tag string) (value, arg string) {
	value, arg, _ = LookupTag(sf, tag)
	return
}

// LookupTag returns the value and arg about the tag from the struct field.
func LookupTag(sf reflect.StructField, tag string) (value, arg string, ok bool) {
	if tag == "" {
		panic("LookupTag: tag must not be empty")
	}

	if value, ok = sf.Tag.Lookup(tag); !ok {
		return
	}

	if index := strings.IndexByte(value, ','); index > -1 {
		arg = strings.TrimSpace(value[index+1:])
		value = strings.TrimSpace(value[:index])
	}

	return
}
