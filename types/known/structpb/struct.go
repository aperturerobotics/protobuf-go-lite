package structpb

import (
	"errors"
	"maps"
	"math"
	"slices"
	"strconv"

	jsoniter "github.com/aperturerobotics/json-iterator-lite"
	"github.com/aperturerobotics/protobuf-go-lite/json"
)

// ErrJSONNotSupported was returned before structured JSON values were supported.
// Deprecated: structured JSON codecs no longer return this error.
var ErrJSONNotSupported = errors.New("JSON marshal/unmarshal is not supported for Struct")

// MarshalJSON marshals the Struct to JSON.
func (x *Struct) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the Struct from JSON.
func (x *Struct) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a Struct from JSON.
func (x *Struct) UnmarshalProtoJSON(s *json.UnmarshalState) {
	x.Fields = nil
	if s.ReadNil() {
		return
	}

	x.Fields = make(map[string]*Value)
	s.ReadObject(func(key string) {
		field := s.WithField(key, false)
		if _, exists := x.Fields[key]; exists {
			field.SetError(errors.New("duplicate Struct field"))
			return
		}

		value := new(Value)
		value.UnmarshalProtoJSON(field)
		x.Fields[key] = value
	})
}

// MarshalProtoJSON marshals a Struct to JSON.
func (x *Struct) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}

	s.WriteObjectStart()
	for i, key := range slices.Sorted(maps.Keys(x.Fields)) {
		if i > 0 {
			s.WriteMore()
		}

		s.WriteObjectField(key)
		x.Fields[key].MarshalProtoJSON(s.WithField(key))
	}
	s.WriteObjectEnd()
}

// MarshalJSON marshals the Value to JSON.
func (x *Value) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the Value from JSON.
func (x *Value) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// MarshalProtoJSON marshals a Value to JSON.
func (x *Value) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}

	switch kind := x.Kind.(type) {
	case *Value_NullValue:
		s.WriteNil()
	case *Value_NumberValue:
		// Encoding nonfinite numbers as strings would change the Value's kind.
		if math.IsNaN(kind.NumberValue) || math.IsInf(kind.NumberValue, 0) {
			s.SetError(errors.New("nonfinite Value number"))
			return
		}

		s.WriteFloat64(kind.NumberValue)
	case *Value_StringValue:
		s.WriteString(kind.StringValue)
	case *Value_BoolValue:
		s.WriteBool(kind.BoolValue)
	case *Value_StructValue:
		kind.StructValue.MarshalProtoJSON(s)
	case *Value_ListValue:
		kind.ListValue.MarshalProtoJSON(s)
	default:
		s.SetError(errors.New("Value has no kind"))
	}
}

// UnmarshalProtoJSON replaces the Value with the next JSON value.
func (x *Value) UnmarshalProtoJSON(s *json.UnmarshalState) {
	x.Kind = nil
	switch s.WhatIsNext() {
	case jsoniter.NilValue:
		s.ReadNil()
		x.Kind = &Value_NullValue{}
	case jsoniter.NumberValue:
		x.Kind = &Value_NumberValue{NumberValue: s.ReadFloat64()}
	case jsoniter.StringValue:
		x.Kind = &Value_StringValue{StringValue: s.ReadString()}
	case jsoniter.BoolValue:
		x.Kind = &Value_BoolValue{BoolValue: s.ReadBool()}
	case jsoniter.ObjectValue:
		value := new(Struct)
		value.UnmarshalProtoJSON(s)
		x.Kind = &Value_StructValue{StructValue: value}
	case jsoniter.ArrayValue:
		value := new(ListValue)
		value.UnmarshalProtoJSON(s)
		x.Kind = &Value_ListValue{ListValue: value}
	default:
		s.SetError(errors.New("expected JSON value"))
	}
}

// MarshalJSON marshals the ListValue to JSON.
func (x *ListValue) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the ListValue from JSON.
func (x *ListValue) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// MarshalProtoJSON marshals a ListValue to JSON.
func (x *ListValue) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}

	s.WriteArrayStart()
	for i, value := range x.Values {
		if i > 0 {
			s.WriteMore()
		}

		value.MarshalProtoJSON(s.WithField(strconv.Itoa(i)))
	}
	s.WriteArrayEnd()
}

// UnmarshalProtoJSON replaces the ListValue with the next JSON array.
func (x *ListValue) UnmarshalProtoJSON(s *json.UnmarshalState) {
	x.Values = nil
	if s.ReadNil() {
		return
	}

	s.ReadArray(func() {
		value := new(Value)
		value.UnmarshalProtoJSON(s.WithField(strconv.Itoa(len(x.Values)), false))
		x.Values = append(x.Values, value)
	})
}

// MarshalJSON marshals the NullValue to JSON.
func (x *NullValue) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the NullValue from JSON.
func (x *NullValue) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// MarshalProtoJSON marshals a NullValue to JSON.
func (x *NullValue) MarshalProtoJSON(s *json.MarshalState) {
	s.WriteNil()
}

// UnmarshalProtoJSON reads JSON null into the singleton enum value.
func (x *NullValue) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if !s.ReadNil() {
		s.SetError(errors.New("expected JSON null"))
		return
	}

	*x = NullValue_NULL_VALUE
}
