// Copyright 2017 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// package promql
package letheql

import (
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/letheql/parser"
)

func (String) Type() parser.ValueType { return parser.ValueTypeString }

// String represents a string value.
type String struct {
	T int64
	V string
}

func (s String) String() string {
	return s.V
}

// Result holds the resulting value of an execution or an error
// if any occurred.
type Result struct {
	Err      error
	Value    parser.Value
	Warnings model.Warnings
}

func (r *Result) String() string {
	if r.Err != nil {
		return r.Err.Error()
	}
	if r.Value == nil {
		return ""
	}
	return r.Value.String()
}
