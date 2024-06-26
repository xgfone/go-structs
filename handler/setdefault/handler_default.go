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

// Package setdefault provides a handler to set the default of the struct field.
package setdefault

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/xgfone/go-defaults"
	"github.com/xgfone/go-structs/field"
	"github.com/xgfone/go-structs/handler"
	"github.com/xgfone/go-structs/handler/setter"
)

var (
	// ParseTime is used to parse a string to time.Time.
	ParseTime func(string) (time.Time, error) = parseTime

	// ParseDuration is used to parse a string to time.Duration.
	ParseDuration func(string) (time.Duration, error) = parseDuration
)

// SetDefaultRunnder returns a runner to set the default value
// of the struct field if it is ZERO, which is registered into DefaultReflector
// with the tag name "default" by default.
//
// For the type of the field, it only supports some base types as follow:
//
//	bool
//	string
//	float32
//	float64
//	int
//	int8
//	int16
//	int32
//	int64
//	uint
//	uint8
//	uint16
//	uint32
//	uint64
//	struct
//	struct slice
//	time.Time      // Format: A. Integer(UTC); B. String(RFC3339)
//	time.Duration  // Format: A. Integer(ms);  B. String(time.ParseDuration)
//
// And the pointer to the types above, and interface{ Set(interface{}) error }.
//
// If the field type is string or int64, and the tag value is like "now()"
// or "now(layout)", set the default value of the field to the current time
// by defaults.Now(). For example,
//
//	type T struct {
//	    StartTime string `default:"now()"`
//	    EndTime   int64  `default:"now()"`
//	}
//
// Notice: If the tag value starts with ".", it represents a field name and
// the default value of current field is set to the value of that field.
// But their types must be consistent, or panic.
func SetDefaultRunner() handler.Runner {
	return setter.SetterRunner(setdefault)
}

func setdefault(_ interface{}, root, fieldptr reflect.Value, sf reflect.StructField, arg interface{}) error {
	v := fieldptr.Elem()
	if !v.IsZero() {
		return nil
	}

	s := arg.(string)
	if len(s) > 0 && s[0] == '.' {
		if s = s[1:]; s == "" {
			return fmt.Errorf("%s: invalid default value", sf.Name)
		}

		fieldv, ok := field.GetValueByName(root, s)
		if !ok {
			panic(fmt.Errorf("not found the struct field '%s'", s))
		}

		if fieldv.Kind() == reflect.Pointer {
			fieldv = fieldv.Elem()
		}
		v.Set(fieldv)
		return nil
	}

	if i, ok := fieldptr.Interface().(interface{ Set(interface{}) error }); ok {
		return i.Set(s)
	}

	switch v.Kind() {
	case reflect.String:
		if strings.HasPrefix(s, "now(") && strings.HasSuffix(s, ")") {
			if layout := s[4 : len(s)-1]; layout == "" {
				s = defaults.Now().Format(time.RFC3339)
			} else {
				s = defaults.Now().Format(layout)
			}
		}
		v.SetString(s)

	case reflect.Bool:
		i, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(i)

	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(i)

	case reflect.Int64:
		if _, ok := v.Interface().(time.Duration); ok {
			i, err := ParseDuration(s)
			if err == nil {
				v.SetInt(int64(i))
			}
			return err
		}

		var e error
		var i int64
		if strings.HasPrefix(s, "now(") && strings.HasSuffix(s, ")") {
			i = defaults.Now().Unix()
		} else if i, e = strconv.ParseInt(s, 10, 64); e != nil {
			return e
		}
		v.SetInt(i)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)

	case reflect.Struct:
		if _, ok := v.Interface().(time.Time); !ok {
			return fmt.Errorf("%s: unsupported type %T", sf.Name, v.Interface())
		}

		i, err := ParseTime(s)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(i))

	default:
		return fmt.Errorf("%s: unsupported type %T", sf.Name, v.Interface())
	}

	return nil
}

func parseDuration(src string) (dst time.Duration, err error) {
	_len := len(src)
	if _len == 0 {
		return
	}

	switch src[_len-1] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var i int64
		i, err = strconv.ParseInt(src, 10, 64)
		dst = time.Duration(i) * time.Millisecond
	default:
		dst, err = time.ParseDuration(src)
	}

	return
}

func parseTime(value string) (time.Time, error) {
	loc := defaults.TimeLocation.Get()

	switch value {
	case "", "0000-00-00 00:00:00", "0000-00-00 00:00:00.000", "0000-00-00 00:00:00.000000":
		return time.Time{}.In(loc), nil
	}

	if isIntegerString(value) {
		i, err := strconv.ParseInt(value, 10, 64)
		return time.Unix(i, 0).In(loc), err
	}

	for _, layout := range defaults.TimeFormats.Get() {
		if t, err := time.ParseInLocation(layout, value, loc); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time '%s'", value)
}

func isIntegerString(s string) bool {
	_len := len(s)
	if _len == 0 {
		return false
	}

	switch s[0] {
	case '-', '+':
		s = s[1:]
		_len--
	}

	for i := 0; i < _len; i++ {
		switch s[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			return false
		}
	}
	return true
}
