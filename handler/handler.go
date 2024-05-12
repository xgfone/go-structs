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

// Package handler provides a handler interface.
package handler

import "reflect"

// Handler is an interface to handle the struct field.
type Handler interface {
	Parse(string) (any, error) // Used to optimize: pre-parse and cache the tag value.
	Run(ctx any, rootStructValue, fieldValue reflect.Value, fieldType reflect.StructField, arg any) error
}

// Parser is the function to pre-parse the field tag value.
type Parser func(string) (any, error)

// Runner is the function to handle the struct field.
type Runner func(ctx any, root reflect.Value, vf reflect.Value, sf reflect.StructField, arg any) error

// SimpleRunner converts a simple function to Runner.
func SimpleRunner(f func(field reflect.Value, arg any) error) Runner {
	return func(_ any, _, vf reflect.Value, _ reflect.StructField, arg any) error {
		return f(vf, arg)
	}
}

// FieldRunner converts a field value function to Runner.
func FieldRunner(f func(reflect.Value, reflect.StructField, any) error) Runner {
	return func(_ any, _, vf reflect.Value, sf reflect.StructField, arg any) error {
		return f(vf, sf, arg)
	}
}

// Parse implements the interface Handler, which does nothing
// and returns the original string input as the parsed result.
func (f Runner) Parse(s string) (any, error) { return s, nil }

// Run implements the interface Handler.
func (f Runner) Run(ctx any, r, v reflect.Value, t reflect.StructField, arg any) error {
	return f(ctx, r, v, t, arg)
}

// New returns a new Handler from the parse and run functions.
//
// parse may be nil, but run must not be nil.
func New(parse Parser, run Runner) Handler {
	if parse == nil {
		parse = noopparse
	}
	return handler{parse: parse, run: run}
}

type handler struct {
	parse Parser
	run   Runner
}

func noopparse(s string) (any, error)         { return s, nil }
func (h handler) Parse(s string) (any, error) { return h.parse(s) }
func (h handler) Run(c any, r, v reflect.Value, t reflect.StructField, a any) error {
	return h.run(c, r, v, t, a)
}
