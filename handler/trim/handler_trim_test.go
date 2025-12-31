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

package trim_test

import (
	"testing"

	"github.com/xgfone/go-structs"
	"github.com/xgfone/go-structs/handler/trim"
)

func TestTrimRunner(t *testing.T) {
	structs.Register("trim", trim.TrimRunner())

	type T1 struct {
		Name string `trim:""`
	}

	v1 := T1{Name: " xgfone "}
	if err := structs.Reflect(&v1); err != nil {
		t.Fatal(err)
	} else if v1.Name != "xgfone" {
		t.Errorf("expect %s, but got %s", "xgfone", v1.Name)
	}

	type T2 struct {
		Name string `trim:"#@"`
	}

	v2 := T2{Name: "@#xgfone@#"}
	if err := structs.Reflect(&v2); err != nil {
		t.Fatal(err)
	} else if v2.Name != "xgfone" {
		t.Errorf("expect %s, but got %s", "xgfone", v2.Name)
	}
}
