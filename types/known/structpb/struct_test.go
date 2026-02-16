// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structpb_test

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/reflect/protoreflect"

	spb "github.com/aperturerobotics/protobuf-go-lite/types/known/structpb"
)

func TestToStruct(t *testing.T) {
	tests := []struct {
		in      map[string]any
		wantPB  *spb.Struct
		wantErr error
	}{{
		in:     nil,
		wantPB: new(spb.Struct),
	}, {
		in:     make(map[string]any),
		wantPB: new(spb.Struct),
	}, {
		in: map[string]any{
			"nil":     nil,
			"bool":    bool(false),
			"int":     int(-123),
			"int32":   int32(math.MinInt32),
			"int64":   int64(math.MinInt64),
			"uint":    uint(123),
			"uint32":  uint32(math.MaxInt32),
			"uint64":  uint64(math.MaxInt64),
			"float32": float32(123.456),
			"float64": float64(123.456),
			"string":  string("hello, world!"),
			"bytes":   []byte("\xde\xad\xbe\xef"),
			"map":     map[string]any{"k1": "v1", "k2": "v2"},
			"slice":   []any{"one", "two", "three"},
		},
		wantPB: &spb.Struct{Fields: map[string]*spb.Value{
			"nil":     spb.NewNullValue(),
			"bool":    spb.NewBoolValue(false),
			"int":     spb.NewNumberValue(float64(-123)),
			"int32":   spb.NewNumberValue(float64(math.MinInt32)),
			"int64":   spb.NewNumberValue(float64(math.MinInt64)),
			"uint":    spb.NewNumberValue(float64(123)),
			"uint32":  spb.NewNumberValue(float64(math.MaxInt32)),
			"uint64":  spb.NewNumberValue(float64(math.MaxInt64)),
			"float32": spb.NewNumberValue(float64(float32(123.456))),
			"float64": spb.NewNumberValue(float64(float64(123.456))),
			"string":  spb.NewStringValue("hello, world!"),
			"bytes":   spb.NewStringValue("3q2+7w=="),
			"map": spb.NewStructValue(&spb.Struct{Fields: map[string]*spb.Value{
				"k1": spb.NewStringValue("v1"),
				"k2": spb.NewStringValue("v2"),
			}}),
			"slice": spb.NewListValue(&spb.ListValue{Values: []*spb.Value{
				spb.NewStringValue("one"),
				spb.NewStringValue("two"),
				spb.NewStringValue("three"),
			}}),
		}},
	}, {
		in:      map[string]any{"\xde\xad\xbe\xef": "<invalid UTF-8>"},
		wantErr: cmpopts.AnyError,
	}, {
		in:      map[string]any{"<invalid UTF-8>": "\xde\xad\xbe\xef"},
		wantErr: cmpopts.AnyError,
	}, {
		in:      map[string]any{"key": protoreflect.Name("named string")},
		wantErr: cmpopts.AnyError,
	}}

	for _, tt := range tests {
		_, gotErr := spb.NewStruct(tt.in)
		if diff := cmp.Diff(tt.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
			t.Errorf("NewStruct(%v) error mismatch (-want +got):\n%s", tt.in, diff)
		}
	}
}

func TestFromStruct(t *testing.T) {
	tests := []struct {
		in   *spb.Struct
		want map[string]any
	}{{
		in:   nil,
		want: make(map[string]any),
	}, {
		in:   new(spb.Struct),
		want: make(map[string]any),
	}, {
		in:   &spb.Struct{Fields: make(map[string]*spb.Value)},
		want: make(map[string]any),
	}, {
		in: &spb.Struct{Fields: map[string]*spb.Value{
			"nil":     spb.NewNullValue(),
			"bool":    spb.NewBoolValue(false),
			"int":     spb.NewNumberValue(float64(-123)),
			"int32":   spb.NewNumberValue(float64(math.MinInt32)),
			"int64":   spb.NewNumberValue(float64(math.MinInt64)),
			"uint":    spb.NewNumberValue(float64(123)),
			"uint32":  spb.NewNumberValue(float64(math.MaxInt32)),
			"uint64":  spb.NewNumberValue(float64(math.MaxInt64)),
			"float32": spb.NewNumberValue(float64(float32(123.456))),
			"float64": spb.NewNumberValue(float64(float64(123.456))),
			"string":  spb.NewStringValue("hello, world!"),
			"bytes":   spb.NewStringValue("3q2+7w=="),
			"map": spb.NewStructValue(&spb.Struct{Fields: map[string]*spb.Value{
				"k1": spb.NewStringValue("v1"),
				"k2": spb.NewStringValue("v2"),
			}}),
			"slice": spb.NewListValue(&spb.ListValue{Values: []*spb.Value{
				spb.NewStringValue("one"),
				spb.NewStringValue("two"),
				spb.NewStringValue("three"),
			}}),
		}},
		want: map[string]any{
			"nil":     nil,
			"bool":    bool(false),
			"int":     float64(-123),
			"int32":   float64(math.MinInt32),
			"int64":   float64(math.MinInt64),
			"uint":    float64(123),
			"uint32":  float64(math.MaxInt32),
			"uint64":  float64(math.MaxInt64),
			"float32": float64(float32(123.456)),
			"float64": float64(float64(123.456)),
			"string":  string("hello, world!"),
			"bytes":   string("3q2+7w=="),
			"map":     map[string]any{"k1": "v1", "k2": "v2"},
			"slice":   []any{"one", "two", "three"},
		},
	}}

	for _, tt := range tests {
		got := tt.in.AsMap()
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("AsMap(%v) mismatch (-want +got):\n%s", tt.in, diff)
		}
		/*
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("Marshal error: %v", err)
			}
			wantJSON, err := tt.in.MarshalJSON()
			if err != nil {
				t.Errorf("Marshal error: %v", err)
			}
			if diff := cmp.Diff(wantJSON, gotJSON, equateJSON); diff != "" {
				t.Errorf("MarshalJSON(%v) mismatch (-want +got):\n%s", tt.in, diff)
			}
		*/
	}
}

func TestFromListValue(t *testing.T) {
	tests := []struct {
		in   *spb.ListValue
		want []any
	}{{
		in:   nil,
		want: make([]any, 0),
	}, {
		in:   new(spb.ListValue),
		want: make([]any, 0),
	}, {
		in:   &spb.ListValue{Values: make([]*spb.Value, 0)},
		want: make([]any, 0),
	}, {
		in: &spb.ListValue{Values: []*spb.Value{
			spb.NewNullValue(),
			spb.NewBoolValue(false),
			spb.NewNumberValue(float64(-123)),
			spb.NewNumberValue(float64(math.MinInt32)),
			spb.NewNumberValue(float64(math.MinInt64)),
			spb.NewNumberValue(float64(123)),
			spb.NewNumberValue(float64(math.MaxInt32)),
			spb.NewNumberValue(float64(math.MaxInt64)),
			spb.NewNumberValue(float64(float32(123.456))),
			spb.NewNumberValue(float64(float64(123.456))),
			spb.NewStringValue("hello, world!"),
			spb.NewStringValue("3q2+7w=="),
			spb.NewStructValue(&spb.Struct{Fields: map[string]*spb.Value{
				"k1": spb.NewStringValue("v1"),
				"k2": spb.NewStringValue("v2"),
			}}),
			spb.NewListValue(&spb.ListValue{Values: []*spb.Value{
				spb.NewStringValue("one"),
				spb.NewStringValue("two"),
				spb.NewStringValue("three"),
			}}),
		}},
		want: []any{
			nil,
			bool(false),
			float64(-123),
			float64(math.MinInt32),
			float64(math.MinInt64),
			float64(123),
			float64(math.MaxInt32),
			float64(math.MaxInt64),
			float64(float32(123.456)),
			float64(float64(123.456)),
			string("hello, world!"),
			string("3q2+7w=="),
			map[string]any{"k1": "v1", "k2": "v2"},
			[]any{"one", "two", "three"},
		},
	}}

	for _, tt := range tests {
		got := tt.in.AsSlice()
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("AsSlice(%v) mismatch (-want +got):\n%s", tt.in, diff)
		}
		/*
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("Marshal error: %v", err)
			}
			wantJSON, err := tt.in.MarshalJSON()
			if err != nil {
				t.Errorf("Marshal error: %v", err)
			}
			if diff := cmp.Diff(wantJSON, gotJSON, equateJSON); diff != "" {
				t.Errorf("MarshalJSON(%v) mismatch (-want +got):\n%s", tt.in, diff)
			}
		*/
	}
}

func TestFromValue(t *testing.T) {
	tests := []struct {
		in   *spb.Value
		want any
	}{{
		in:   nil,
		want: nil,
	}, {
		in:   new(spb.Value),
		want: nil,
	}, {
		in:   &spb.Value{Kind: (*spb.Value_NullValue)(nil)},
		want: nil,
	}, {
		in:   spb.NewNullValue(),
		want: nil,
	}, {
		in:   &spb.Value{Kind: &spb.Value_NullValue{NullValue: math.MinInt32}},
		want: nil,
	}, {
		in:   &spb.Value{Kind: (*spb.Value_BoolValue)(nil)},
		want: nil,
	}, {
		in:   spb.NewBoolValue(false),
		want: bool(false),
	}, {
		in:   &spb.Value{Kind: (*spb.Value_NumberValue)(nil)},
		want: nil,
	}, {
		in:   spb.NewNumberValue(float64(math.MinInt32)),
		want: float64(math.MinInt32),
	}, {
		in:   spb.NewNumberValue(float64(math.MinInt64)),
		want: float64(math.MinInt64),
	}, {
		in:   spb.NewNumberValue(float64(123)),
		want: float64(123),
	}, {
		in:   spb.NewNumberValue(float64(math.MaxInt32)),
		want: float64(math.MaxInt32),
	}, {
		in:   spb.NewNumberValue(float64(math.MaxInt64)),
		want: float64(math.MaxInt64),
	}, {
		in:   spb.NewNumberValue(float64(float32(123.456))),
		want: float64(float32(123.456)),
	}, {
		in:   spb.NewNumberValue(float64(float64(123.456))),
		want: float64(float64(123.456)),
	}, {
		in:   spb.NewNumberValue(math.NaN()),
		want: string("NaN"),
	}, {
		in:   spb.NewNumberValue(math.Inf(-1)),
		want: string("-Infinity"),
	}, {
		in:   spb.NewNumberValue(math.Inf(+1)),
		want: string("Infinity"),
	}, {
		in:   &spb.Value{Kind: (*spb.Value_StringValue)(nil)},
		want: nil,
	}, {
		in:   spb.NewStringValue("hello, world!"),
		want: string("hello, world!"),
	}, {
		in:   spb.NewStringValue("3q2+7w=="),
		want: string("3q2+7w=="),
	}, {
		in:   &spb.Value{Kind: (*spb.Value_StructValue)(nil)},
		want: nil,
	}, {
		in:   &spb.Value{Kind: &spb.Value_StructValue{}},
		want: make(map[string]any),
	}, {
		in: spb.NewListValue(&spb.ListValue{Values: []*spb.Value{
			spb.NewStringValue("one"),
			spb.NewStringValue("two"),
			spb.NewStringValue("three"),
		}}),
		want: []any{"one", "two", "three"},
	}, {
		in:   &spb.Value{Kind: (*spb.Value_ListValue)(nil)},
		want: nil,
	}, {
		in:   &spb.Value{Kind: &spb.Value_ListValue{}},
		want: make([]any, 0),
	}, {
		in: spb.NewListValue(&spb.ListValue{Values: []*spb.Value{
			spb.NewStringValue("one"),
			spb.NewStringValue("two"),
			spb.NewStringValue("three"),
		}}),
		want: []any{"one", "two", "three"},
	}}

	for _, tt := range tests {
		got := tt.in.AsInterface()
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("AsInterface(%v) mismatch (-want +got):\n%s", tt.in, diff)
		}
		/*
			gotJSON, gotErr := json.Marshal(got)
			if gotErr != nil {
				t.Errorf("Marshal error: %v", gotErr)
			}
			wantJSON, wantErr := tt.in.MarshalJSON()
			if diff := cmp.Diff(wantJSON, gotJSON, equateJSON); diff != "" && wantErr == nil {
				t.Errorf("MarshalJSON(%v) mismatch (-want +got):\n%s", tt.in, diff)
			}
		*/
	}
}
