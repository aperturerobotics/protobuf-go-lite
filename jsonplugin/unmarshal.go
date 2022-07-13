// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package jsonplugin

import (
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// Unmarshaler is the interface implemented by types that are supported by this plugin.
type Unmarshaler interface {
	UnmarshalProtoJSON(*UnmarshalState)
}

type unmarshalError struct {
	Err  error
	Path *path
}

func (e *unmarshalError) Error() string {
	if e.Path != nil {
		return fmt.Sprintf("unmarshal error at path %q: %v", e.Path, e.Err)
	}
	return fmt.Sprintf("unmarshal error: %v", e.Err)
}

// UnmarshalerConfig is the configuration for the Unmarshaler.
type UnmarshalerConfig struct{}

// DefaultUnmarshalerConfig is the default configuration for the Unmarshaler.
var DefaultUnmarshalerConfig = UnmarshalerConfig{}

// Unmarshal unmarshals a message.
func (c UnmarshalerConfig) Unmarshal(data []byte, m Unmarshaler) error {
	s := NewUnmarshalState(data, c)
	m.UnmarshalProtoJSON(s)
	return s.Err()
}

// UnmarshalState is the internal state of the Unmarshaler.
type UnmarshalState struct {
	inner  *jsoniter.Iterator
	config *UnmarshalerConfig

	err   *unmarshalError
	path  *path
	paths *pathSlice
}

// NewUnmarshalState creates a new UnmarshalState.
func NewUnmarshalState(data []byte, config UnmarshalerConfig) *UnmarshalState {
	return &UnmarshalState{
		inner:  jsoniter.ParseBytes(jsoniterConfig, data),
		config: &config,

		err:   &unmarshalError{},
		path:  nil,
		paths: &pathSlice{},
	}
}

// Config returns a copy of the unmarshaler configuration.
func (s *UnmarshalState) Config() UnmarshalerConfig {
	return *s.config
}

// Sub returns a subunmarshaler with a new buffer, but with the same configuration, error and path info.
func (s *UnmarshalState) Sub(data []byte) *UnmarshalState {
	return &UnmarshalState{
		inner:  jsoniter.ParseBytes(jsoniterConfig, data),
		config: s.config,

		err:   s.err,
		path:  s.path,
		paths: &pathSlice{},
	}
}

// Err returns an error from the marshaler, if any.
func (s *UnmarshalState) Err() error {
	if s.err.Err != nil {
		return s.err
	}
	if s.inner.Error != nil && s.inner.Error != io.EOF {
		return s.inner.Error
	}
	return nil
}

// SetError sets an error in the unmarshaler state.
// Subsequent operations become no-ops.
func (s *UnmarshalState) SetError(err error) {
	if s.Err() != nil {
		return
	}
	s.err.Err = err
	s.err.Path = s.path
}

// SetErrorf calls SetError with a formatted error.
func (s *UnmarshalState) SetErrorf(format string, a ...interface{}) {
	s.SetError(fmt.Errorf(format, a...))
}

// WithField returns a UnmarshalState for the given subfield.
func (s *UnmarshalState) WithField(field string, mask bool) *UnmarshalState {
	fm := s.paths
	if !mask {
		fm = &pathSlice{}
	}
	return &UnmarshalState{
		inner:  s.inner,
		config: s.config,

		err:   s.err,
		path:  s.path.push(field),
		paths: fm,
	}
}

// AddField registers a field in the field mask of the unmarshaler state.
func (s *UnmarshalState) AddField(field string) {
	s.paths.add(*s.path.push(field))
}

// FieldMask returns the field mask containing the unmarshaled fields.
func (s *UnmarshalState) FieldMask() FieldMask {
	return s.paths
}

// ReadFloat32 reads a float32 value. This also supports string encoding.
func (s *UnmarshalState) ReadFloat32() float32 {
	if s.Err() != nil {
		return 0
	}
	switch any := s.inner.ReadAny(); any.ValueType() {
	case jsoniter.NumberValue:
		return any.ToFloat32()
	case jsoniter.StringValue:
		f, err := strconv.ParseFloat(any.ToString(), 32)
		if err != nil {
			s.SetErrorf("invalid value for float32: %w", err)
			return 0
		}
		return float32(f)
	default:
		s.SetErrorf("invalid value type for float32: %s", valueTypeString(any.ValueType()))
		return 0
	}
}

// ReadWrappedFloat32 reads a wrapped float32 value. This also supports string encoding as well as {"value": ...}.
func (s *UnmarshalState) ReadWrappedFloat32() float32 {
	if s.Err() != nil {
		return 0
	}
	if s.inner.WhatIsNext() != jsoniter.ObjectValue {
		return s.ReadFloat32()
	}
	if key := s.ReadObjectField(); key != "value" {
		s.SetErrorf("first field in wrapped float32 is not value, but %q", key)
		return 0
	}
	v := s.ReadFloat32()
	if field := s.ReadObjectField(); field != "" {
		s.SetErrorf("unexpected %q field in wrapped float32", field)
		return 0
	}
	return v
}

// ReadFloat64 reads a float64 value. This also supports string encoding.
func (s *UnmarshalState) ReadFloat64() float64 {
	if s.Err() != nil {
		return 0
	}
	switch any := s.inner.ReadAny(); any.ValueType() {
	case jsoniter.NumberValue:
		return any.ToFloat64()
	case jsoniter.StringValue:
		f, err := strconv.ParseFloat(any.ToString(), 64)
		if err != nil {
			s.SetErrorf("invalid value for float64: %w", err)
			return 0
		}
		return float64(f)
	default:
		s.SetErrorf("invalid value type for float64: %s", valueTypeString(any.ValueType()))
		return 0
	}
}

// ReadWrappedFloat64 reads a wrapped float64 value. This also supports string encoding as well as {"value": ...}.
func (s *UnmarshalState) ReadWrappedFloat64() float64 {
	if s.Err() != nil {
		return 0
	}
	if s.inner.WhatIsNext() != jsoniter.ObjectValue {
		return s.ReadFloat64()
	}
	if key := s.ReadObjectField(); key != "value" {
		s.SetErrorf("first field in wrapped float64 is not value, but %q", key)
		return 0
	}
	v := s.ReadFloat64()
	if field := s.ReadObjectField(); field != "" {
		s.SetErrorf("unexpected %q field in wrapped float64", field)
		return 0
	}
	return v
}

// ReadFloat32Array reads an array of float32 values.
func (s *UnmarshalState) ReadFloat32Array() []float32 {
	var arr []float32
	s.ReadArray(func() {
		n := s.ReadFloat32()
		if s.Err() != nil {
			return
		}
		arr = append(arr, n)
	})
	if s.Err() != nil {
		return nil
	}
	return arr
}

// ReadFloat64Array reads an array of float64 values.
func (s *UnmarshalState) ReadFloat64Array() []float64 {
	var arr []float64
	s.ReadArray(func() {
		n := s.ReadFloat64()
		if s.Err() != nil {
			return
		}
		arr = append(arr, n)
	})
	if s.Err() != nil {
		return nil
	}
	return arr
}

// ReadInt32 reads a int32 value. This also supports string encoding.
func (s *UnmarshalState) ReadInt32() int32 {
	if s.Err() != nil {
		return 0
	}
	switch any := s.inner.ReadAny(); any.ValueType() {
	case jsoniter.NumberValue:
		return any.ToInt32()
	case jsoniter.StringValue:
		f, err := strconv.ParseInt(any.ToString(), 10, 32)
		if err != nil {
			s.SetErrorf("invalid value for int32: %w", err)
			return 0
		}
		return int32(f)
	default:
		s.SetErrorf("invalid value type for int32: %s", valueTypeString(any.ValueType()))
		return 0
	}
}

// ReadWrappedInt32 reads a wrapped int32 value. This also supports string encoding as well as {"value": ...}.
func (s *UnmarshalState) ReadWrappedInt32() int32 {
	if s.Err() != nil {
		return 0
	}
	if s.inner.WhatIsNext() != jsoniter.ObjectValue {
		return s.ReadInt32()
	}
	if key := s.ReadObjectField(); key != "value" {
		s.SetErrorf("first field in wrapped int32 is not value, but %q", key)
		return 0
	}
	v := s.ReadInt32()
	if field := s.ReadObjectField(); field != "" {
		s.SetErrorf("unexpected %q field in wrapped int32", field)
		return 0
	}
	return v
}

// ReadInt64 reads a int64 value. This also supports string encoding.
func (s *UnmarshalState) ReadInt64() int64 {
	if s.Err() != nil {
		return 0
	}
	switch any := s.inner.ReadAny(); any.ValueType() {
	case jsoniter.NumberValue:
		return any.ToInt64()
	case jsoniter.StringValue:
		f, err := strconv.ParseInt(any.ToString(), 10, 64)
		if err != nil {
			s.SetErrorf("invalid value for int64: %w", err)
			return 0
		}
		return f
	default:
		s.SetErrorf("invalid value type for int64: %s", valueTypeString(any.ValueType()))
		return 0
	}
}

// ReadWrappedInt64 reads a wrapped int64 value. This also supports string encoding as well as {"value": ...}.
func (s *UnmarshalState) ReadWrappedInt64() int64 {
	if s.Err() != nil {
		return 0
	}
	if s.inner.WhatIsNext() != jsoniter.ObjectValue {
		return s.ReadInt64()
	}
	if key := s.ReadObjectField(); key != "value" {
		s.SetErrorf("first field in wrapped int64 is not value, but %q", key)
		return 0
	}
	v := s.ReadInt64()
	if field := s.ReadObjectField(); field != "" {
		s.SetErrorf("unexpected %q field in wrapped int64", field)
		return 0
	}
	return v
}

// ReadInt32Array reads an array of int32 values.
func (s *UnmarshalState) ReadInt32Array() []int32 {
	var arr []int32
	s.ReadArray(func() {
		n := s.ReadInt32()
		if s.Err() != nil {
			return
		}
		arr = append(arr, n)
	})
	if s.Err() != nil {
		return nil
	}
	return arr
}

// ReadInt64Array reads an array of int64 values.
func (s *UnmarshalState) ReadInt64Array() []int64 {
	var arr []int64
	s.ReadArray(func() {
		n := s.ReadInt64()
		if s.Err() != nil {
			return
		}
		arr = append(arr, n)
	})
	if s.Err() != nil {
		return nil
	}
	return arr
}

// ReadUint32 reads a uint32 value. This also supports string encoding.
func (s *UnmarshalState) ReadUint32() uint32 {
	if s.Err() != nil {
		return 0
	}
	switch any := s.inner.ReadAny(); any.ValueType() {
	case jsoniter.NumberValue:
		return any.ToUint32()
	case jsoniter.StringValue:
		f, err := strconv.ParseUint(any.ToString(), 10, 32)
		if err != nil {
			s.SetErrorf("invalid value for uint32: %w", err)
			return 0
		}
		return uint32(f)
	default:
		s.SetErrorf("invalid value type for uint32: %s", valueTypeString(any.ValueType()))
		return 0
	}
}

// ReadWrappedUint32 reads a wrapped uint32 value. This also supports string encoding as well as {"value": ...}.
func (s *UnmarshalState) ReadWrappedUint32() uint32 {
	if s.Err() != nil {
		return 0
	}
	if s.inner.WhatIsNext() != jsoniter.ObjectValue {
		return s.ReadUint32()
	}
	if key := s.ReadObjectField(); key != "value" {
		s.SetErrorf("first field in wrapped uint32 is not value, but %q", key)
		return 0
	}
	v := s.ReadUint32()
	if field := s.ReadObjectField(); field != "" {
		s.SetErrorf("unexpected %q field in wrapped uint32", field)
		return 0
	}
	return v
}

// ReadUint64 reads a uint64 value. This also supports string encoding.
func (s *UnmarshalState) ReadUint64() uint64 {
	if s.Err() != nil {
		return 0
	}
	switch any := s.inner.ReadAny(); any.ValueType() {
	case jsoniter.NumberValue:
		return any.ToUint64()
	case jsoniter.StringValue:
		f, err := strconv.ParseUint(any.ToString(), 10, 64)
		if err != nil {
			s.SetErrorf("invalid value for uint64: %w", err)
			return 0
		}
		return f
	default:
		s.SetErrorf("invalid value type for uint64: %s", valueTypeString(any.ValueType()))
		return 0
	}
}

// ReadWrappedUint64 reads a wrapped uint64 value. This also supports string encoding as well as {"value": ...}.
func (s *UnmarshalState) ReadWrappedUint64() uint64 {
	if s.Err() != nil {
		return 0
	}
	if s.inner.WhatIsNext() != jsoniter.ObjectValue {
		return s.ReadUint64()
	}
	if key := s.ReadObjectField(); key != "value" {
		s.SetErrorf("first field in wrapped uint64 is not value, but %q", key)
		return 0
	}
	v := s.ReadUint64()
	if field := s.ReadObjectField(); field != "" {
		s.SetErrorf("unexpected %q field in wrapped uint64", field)
		return 0
	}
	return v
}

// ReadUint32Array reads an array of uint32 values.
func (s *UnmarshalState) ReadUint32Array() []uint32 {
	var arr []uint32
	s.ReadArray(func() {
		n := s.ReadUint32()
		if s.Err() != nil {
			return
		}
		arr = append(arr, n)
	})
	if s.Err() != nil {
		return nil
	}
	return arr
}

// ReadUint64Array reads an array of uint64 values.
func (s *UnmarshalState) ReadUint64Array() []uint64 {
	var arr []uint64
	s.ReadArray(func() {
		n := s.ReadUint64()
		if s.Err() != nil {
			return
		}
		arr = append(arr, n)
	})
	if s.Err() != nil {
		return nil
	}
	return arr
}

// ReadBool reads a bool value.
func (s *UnmarshalState) ReadBool() bool {
	if s.Err() != nil {
		return false
	}
	return s.inner.ReadBool()
}

// ReadWrappedBool reads a wrapped bool value. This also supports {"value": ...}.
func (s *UnmarshalState) ReadWrappedBool() bool {
	if s.Err() != nil {
		return false
	}
	if s.inner.WhatIsNext() != jsoniter.ObjectValue {
		return s.ReadBool()
	}
	if key := s.ReadObjectField(); key != "value" {
		s.SetErrorf("first field in wrapped bool is not value, but %q", key)
		return false
	}
	v := s.ReadBool()
	if field := s.ReadObjectField(); field != "" {
		s.SetErrorf("unexpected %q field in wrapped bool", field)
		return false
	}
	return v
}

// ReadBoolArray reads an array of bool values.
func (s *UnmarshalState) ReadBoolArray() []bool {
	var arr []bool
	s.ReadArray(func() {
		n := s.ReadBool()
		if s.Err() != nil {
			return
		}
		arr = append(arr, n)
	})
	if s.Err() != nil {
		return nil
	}
	return arr
}

// ReadString reads a string value.
func (s *UnmarshalState) ReadString() string {
	if s.Err() != nil {
		return ""
	}
	return s.inner.ReadString()
}

// ReadWrappedString reads a wrapped string value. This also supports {"value": ...}.
func (s *UnmarshalState) ReadWrappedString() string {
	if s.Err() != nil {
		return ""
	}
	if s.inner.WhatIsNext() != jsoniter.ObjectValue {
		return s.ReadString()
	}
	if key := s.ReadObjectField(); key != "value" {
		s.SetErrorf("first field in wrapped string is not value, but %q", key)
		return ""
	}
	v := s.ReadString()
	if field := s.ReadObjectField(); field != "" {
		s.SetErrorf("unexpected %q field in wrapped string", field)
		return ""
	}
	return v
}

// ReadStringArray reads an array of string values.
func (s *UnmarshalState) ReadStringArray() []string {
	var arr []string
	s.ReadArray(func() {
		n := s.ReadString()
		if s.Err() != nil {
			return
		}
		arr = append(arr, n)
	})
	if s.Err() != nil {
		return nil
	}
	return arr
}

var base64Replacer = strings.NewReplacer("_", "/", "-", "+")

// ReadBytes reads a string value.
func (s *UnmarshalState) ReadBytes() []byte {
	if s.Err() != nil {
		return nil
	}
	b64 := s.inner.ReadString()
	if s.Err() != nil {
		return nil
	}
	// According to the protobuf spec, we need to accept both padded and unpadded base64 strings.
	b64 = strings.TrimRight(b64, "=")
	// According to the protobuf spec, we need to accept both standard encoding and URL encoding.
	b64 = base64Replacer.Replace(b64)
	// What's left is raw standard encoding.
	v, err := base64.RawStdEncoding.DecodeString(b64)
	if err != nil {
		s.SetErrorf("invalid value: %w", err)
		return nil
	}
	return v
}

// ReadWrappedBytes reads a wrapped bytes value. This also supports {"value": ...}.
func (s *UnmarshalState) ReadWrappedBytes() []byte {
	if s.Err() != nil {
		return nil
	}
	if s.inner.WhatIsNext() != jsoniter.ObjectValue {
		return s.ReadBytes()
	}
	if key := s.ReadObjectField(); key != "value" {
		s.SetErrorf("first field in wrapped bytes is not value, but %q", key)
		return nil
	}
	v := s.ReadBytes()
	if field := s.ReadObjectField(); field != "" {
		s.SetErrorf("unexpected %q field in wrapped bytes", field)
		return nil
	}
	return v
}

// ReadBytesArray reads an array of []byte values.
func (s *UnmarshalState) ReadBytesArray() [][]byte {
	if s.Err() != nil {
		return nil
	}
	var arr [][]byte
	s.ReadArray(func() {
		n := s.ReadBytes()
		if s.Err() != nil {
			return
		}
		arr = append(arr, n)
	})
	if s.Err() != nil {
		return nil
	}
	return arr
}

// ReadNil reads a null, if there is one.
func (s *UnmarshalState) ReadNil() bool {
	return s.inner.ReadNil()
}

// ReadObjectField reads a single object field.
// An empty string indicates the end of the object.
func (s *UnmarshalState) ReadObjectField() string {
	if s.Err() != nil {
		return ""
	}
	return s.inner.ReadObject()
}

// ReadObject reads all object fields, and calls cb for each.
// cb must always read the value of the field.
func (s *UnmarshalState) ReadObject(cb func(key string)) {
	if s.Err() != nil {
		return
	}
	s.inner.ReadObjectCB(func(_ *jsoniter.Iterator, key string) bool {
		if s.Err() != nil {
			return false
		}
		cb(key)
		return true
	})
}

// ReadBoolMap reads an object where the keys are bool, and calls cb for each field.
// cb must always read the value of the field.
func (s *UnmarshalState) ReadBoolMap(cb func(key bool)) {
	s.ReadObject(func(keyStr string) {
		key, err := strconv.ParseBool(keyStr)
		if err != nil {
			s.SetErrorf("invalid map key %q for bool map", keyStr)
			return
		}
		cb(key)
	})
}

// ReadInt32Map reads an object where the keys are int32, and calls cb for each field.
// cb must always read the value of the field.
func (s *UnmarshalState) ReadInt32Map(cb func(key int32)) {
	s.ReadObject(func(keyStr string) {
		key, err := strconv.ParseInt(keyStr, 10, 32)
		if err != nil {
			s.SetErrorf("invalid map key %q for int32 map", keyStr)
			return
		}
		cb(int32(key))
	})
}

// ReadUint32Map reads an object where the keys are uint32, and calls cb for each field.
// cb must always read the value of the field.
func (s *UnmarshalState) ReadUint32Map(cb func(key uint32)) {
	s.ReadObject(func(keyStr string) {
		key, err := strconv.ParseUint(keyStr, 10, 32)
		if err != nil {
			s.SetErrorf("invalid map key %q for uint32 map", keyStr)
			return
		}
		cb(uint32(key))
	})
}

// ReadInt64Map reads an object where the keys are int64, and calls cb for each field.
// cb must always read the value of the field.
func (s *UnmarshalState) ReadInt64Map(cb func(key int64)) {
	s.ReadObject(func(keyStr string) {
		key, err := strconv.ParseInt(keyStr, 10, 64)
		if err != nil {
			s.SetErrorf("invalid map key %q for int64 map", keyStr)
			return
		}
		cb(key)
	})
}

// ReadUint64Map reads an object where the keys are uint64, and calls cb for each field.
// cb must always read the value of the field.
func (s *UnmarshalState) ReadUint64Map(cb func(key uint64)) {
	s.ReadObject(func(keyStr string) {
		key, err := strconv.ParseUint(keyStr, 10, 64)
		if err != nil {
			s.SetErrorf("invalid map key %q for uint64 map", keyStr)
			return
		}
		cb(key)
	})
}

// ReadStringMap reads an object where the keys are string, and calls cb for each field.
// cb must always read the value of the field.
func (s *UnmarshalState) ReadStringMap(cb func(key string)) {
	s.ReadObject(cb)
}

// ReadArray reads all array elements, and calls cb for each.
// cb must always read the value of the element.
func (s *UnmarshalState) ReadArray(cb func()) {
	if s.Err() != nil {
		return
	}
	s.inner.ReadArrayCB(func(_ *jsoniter.Iterator) bool {
		if s.Err() != nil {
			return false
		}
		cb()
		return true
	})
}

// ParseEnumString parses an enum from its string representation using the value maps.
// If none of the value maps contains a mapping for the string value,
// it attempts to parse the string as a numeric value.
func ParseEnumString(v string, valueMaps ...map[string]int32) (int32, error) {
	for _, valueMap := range valueMaps {
		if x, ok := valueMap[v]; ok {
			return x, nil
		}
	}
	x, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(x), nil
}

// ReadEnum reads an enum. It supports numeric values and string values.
func (s *UnmarshalState) ReadEnum(valueMaps ...map[string]int32) int32 {
	if s.Err() != nil {
		return 0
	}
	switch any := s.inner.ReadAny(); any.ValueType() {
	case jsoniter.NumberValue:
		return any.ToInt32()
	case jsoniter.StringValue:
		v := any.ToString()
		x, err := ParseEnumString(v, valueMaps...)
		if err != nil {
			s.SetErrorf("unknown value for enum: %q", v)
			return 0
		}
		return x
	default:
		s.SetErrorf("invalid value type for enum: %s", valueTypeString(any.ValueType()))
		return 0
	}
}

// ReadTime reads a time.
func (s *UnmarshalState) ReadTime() *time.Time {
	if s.Err() != nil {
		return nil
	}
	if s.ReadNil() {
		return nil
	}
	t, err := time.Parse("2006-01-02T15:04:05.999999999Z", s.inner.ReadString())
	if err != nil {
		s.SetErrorf("invalid time: %w", err)
		return nil
	}
	return &t
}

// ReadDuration reads a duration.
func (s *UnmarshalState) ReadDuration() *time.Duration {
	if s.Err() != nil {
		return nil
	}
	if s.ReadNil() {
		return nil
	}
	d, err := time.ParseDuration(s.inner.ReadString())
	if err != nil {
		s.SetErrorf("invalid duration: %w", err)
	}
	return &d
}

// ReadFieldMask reads a field mask value.
func (s *UnmarshalState) ReadFieldMask() FieldMask {
	if s.Err() != nil {
		return nil
	}
	next := s.inner.WhatIsNext()
	switch next {
	case jsoniter.StringValue:
		mask := newPathSlice(strings.Split(s.ReadString(), ",")...)
		if s.Err() != nil {
			return nil
		}
		return mask
	case jsoniter.ObjectValue:
		if field := s.ReadObjectField(); field != "paths" {
			s.SetErrorf("unexpected %q field in FieldMask object", field)
			return nil
		}
		mask := newPathSlice(s.ReadStringArray()...)
		if s.Err() != nil {
			return nil
		}
		if field := s.ReadObjectField(); field != "" {
			s.SetErrorf("unexpected %q field in FieldMask object", field)
			return nil
		}
		return mask
	}
	s.SetErrorf("invalid value type for field mask: %s", valueTypeString(next))
	return nil
}

// ReadRawMessage reads a raw JSON message.
func (s *UnmarshalState) ReadRawMessage() jsoniter.RawMessage {
	var msg jsoniter.RawMessage
	s.inner.ReadVal(&msg)
	if s.Err() != nil {
		return nil
	}
	return msg
}

// ReadAny reads any type and ignores it.
func (s *UnmarshalState) ReadAny() { s.inner.ReadAny() }
