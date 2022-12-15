// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package jsonplugin

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// Marshaler is the interface implemented by types that are supported by this plugin.
type Marshaler interface {
	MarshalProtoJSON(*MarshalState)
}

type marshalError struct {
	Err  error
	Path *path
}

func (e *marshalError) Error() string {
	if e.Path != nil {
		return fmt.Sprintf("marshal error at path %q: %v", e.Path, e.Err)
	}
	return fmt.Sprintf("marshal error: %v", e.Err)
}

// MarshalerConfig is the configuration for the Marshaler.
type MarshalerConfig struct {
	EnumsAsInts bool
}

// DefaultMarshalerConfig is the default configuration for the Marshaler.
var DefaultMarshalerConfig = MarshalerConfig{
	EnumsAsInts: true,
}

// Marshal marshals a message.
func (c MarshalerConfig) Marshal(m Marshaler) ([]byte, error) {
	s := NewMarshalState(c)
	m.MarshalProtoJSON(s)
	return s.Bytes()
}

// MarshalState is the internal state of the Marshaler.
type MarshalState struct {
	inner  *jsoniter.Stream
	config *MarshalerConfig

	err   *marshalError
	path  *path
	paths *pathSlice
}

// NewMarshalState creates a new MarshalState.
func NewMarshalState(config MarshalerConfig) *MarshalState {
	return &MarshalState{
		inner:  jsoniter.NewStream(jsoniterConfig, nil, 1024),
		config: &config,

		err:   &marshalError{},
		path:  nil,
		paths: &pathSlice{},
	}
}

// Config returns a copy of the marshaler configuration.
func (s *MarshalState) Config() MarshalerConfig {
	return *s.config
}

// Sub returns a sub-marshaler with a new buffer, but with the same configuration, error and path info.
func (s *MarshalState) Sub() *MarshalState {
	return &MarshalState{
		inner:  jsoniter.NewStream(jsoniterConfig, nil, 1024),
		config: s.config,

		err:   s.err,
		path:  s.path,
		paths: &pathSlice{},
	}
}

// Err returns an error from the marshaler, if any.
func (s *MarshalState) Err() error {
	if s.err.Err != nil {
		return s.err
	}
	if s.inner.Error != nil {
		return s.inner.Error
	}
	return nil
}

// SetError sets an error in the marshaler state.
// Subsequent operations become no-ops.
func (s *MarshalState) SetError(err error) {
	if s.Err() != nil {
		return
	}
	s.err.Err = err
	s.err.Path = s.path
}

// SetErrorf calls SetError with a formatted error.
func (s *MarshalState) SetErrorf(format string, a ...interface{}) {
	s.SetError(fmt.Errorf(format, a...))
}

// Bytes returns the buffer of the marshaler.
func (s *MarshalState) Bytes() ([]byte, error) {
	if err := s.Err(); err != nil {
		return nil, err
	}
	return s.inner.Buffer(), nil
}

// WithFieldMask returns a MarshalState for the given field mask.
func (s *MarshalState) WithFieldMask(paths ...string) *MarshalState {
	return &MarshalState{
		inner:  s.inner,
		config: s.config,

		err:   s.err,
		path:  s.path,
		paths: newPathSlice(paths...),
	}
}

// WithField returns a MarshalState for the given subfield.
func (s *MarshalState) WithField(field string) *MarshalState {
	return &MarshalState{
		inner:  s.inner,
		config: s.config,

		err:   s.err,
		path:  s.path.push(field),
		paths: s.paths,
	}
}

// HasField returns whether the field mask contains the given field.
func (s *MarshalState) HasField(field string) bool {
	return s.paths.contains(*s.path.push(field))
}

// Write writes raw data.
func (s *MarshalState) Write(v []byte) (nn int, err error) {
	if s.Err() != nil {
		return
	}
	return s.inner.Write(v)
}

// WriteFloat32 writes a float32 value.
func (s *MarshalState) WriteFloat32(v float32) {
	if s.Err() != nil {
		return
	}
	switch {
	case math.IsInf(float64(v), 1):
		s.inner.WriteString("Infinity")
		return
	case math.IsInf(float64(v), -1):
		s.inner.WriteString("-Infinity")
		return
	case math.IsNaN(float64(v)):
		s.inner.WriteString("NaN")
		return
	}
	s.inner.WriteFloat32(v)
}

// WriteFloat64 writes a float64 value.
func (s *MarshalState) WriteFloat64(v float64) {
	if s.Err() != nil {
		return
	}
	switch {
	case math.IsInf(v, 1):
		s.inner.WriteString("Infinity")
		return
	case math.IsInf(v, -1):
		s.inner.WriteString("-Infinity")
		return
	case math.IsNaN(v):
		s.inner.WriteString("NaN")
		return
	}
	s.inner.WriteFloat64(v)
}

// WriteFloat32Array writes an array of float32 values.
func (s *MarshalState) WriteFloat32Array(vs []float32) {
	if s.Err() != nil {
		return
	}
	s.WriteArrayStart()
	for i, v := range vs {
		if i > 0 {
			s.WriteMore()
		}
		s.WriteFloat32(v)
	}
	s.WriteArrayEnd()
}

// WriteFloat64Array writes an array of float64 values.
func (s *MarshalState) WriteFloat64Array(vs []float64) {
	if s.Err() != nil {
		return
	}
	s.WriteArrayStart()
	for i, v := range vs {
		if i > 0 {
			s.WriteMore()
		}
		s.WriteFloat64(v)
	}
	s.WriteArrayEnd()
}

// WriteInt32 writes an int32 value.
func (s *MarshalState) WriteInt32(v int32) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteInt32(v)
}

// WriteInt64 writes an int64 value as a string.
func (s *MarshalState) WriteInt64(v int64) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteString(strconv.FormatInt(v, 10))
}

// WriteInt32Array writes an array of int32 values.
func (s *MarshalState) WriteInt32Array(vs []int32) {
	if s.Err() != nil {
		return
	}
	s.WriteArrayStart()
	for i, v := range vs {
		if i > 0 {
			s.WriteMore()
		}
		s.WriteInt32(v)
	}
	s.WriteArrayEnd()
}

// WriteInt64Array writes an array of int64 values.
func (s *MarshalState) WriteInt64Array(vs []int64) {
	if s.Err() != nil {
		return
	}
	s.WriteArrayStart()
	for i, v := range vs {
		if i > 0 {
			s.WriteMore()
		}
		s.WriteInt64(v)
	}
	s.WriteArrayEnd()
}

// WriteUint32 writes a uint32 value.
func (s *MarshalState) WriteUint32(v uint32) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteUint32(v)
}

// WriteUint64 writes a uint64 value as a string.
func (s *MarshalState) WriteUint64(v uint64) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteString(strconv.FormatUint(v, 10))
}

// WriteUint32Array writes an array of uint32 values.
func (s *MarshalState) WriteUint32Array(vs []uint32) {
	if s.Err() != nil {
		return
	}
	s.WriteArrayStart()
	for i, v := range vs {
		if i > 0 {
			s.WriteMore()
		}
		s.WriteUint32(v)
	}
	s.WriteArrayEnd()
}

// WriteUint64Array writes an array of uint64 values.
func (s *MarshalState) WriteUint64Array(vs []uint64) {
	if s.Err() != nil {
		return
	}
	s.WriteArrayStart()
	for i, v := range vs {
		if i > 0 {
			s.WriteMore()
		}
		s.WriteUint64(v)
	}
	s.WriteArrayEnd()
}

// WriteBool writes a bool value.
func (s *MarshalState) WriteBool(v bool) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteBool(v)
}

// WriteBoolArray writes an array of bool values.
func (s *MarshalState) WriteBoolArray(vs []bool) {
	if s.Err() != nil {
		return
	}
	s.WriteArrayStart()
	for i, v := range vs {
		if i > 0 {
			s.WriteMore()
		}
		s.WriteBool(v)
	}
	s.WriteArrayEnd()
}

// WriteString writes a string.
func (s *MarshalState) WriteString(v string) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteString(v)
}

// WriteStringArray writes an array of string values.
func (s *MarshalState) WriteStringArray(vs []string) {
	if s.Err() != nil {
		return
	}
	s.WriteArrayStart()
	for i, v := range vs {
		if i > 0 {
			s.WriteMore()
		}
		s.WriteString(v)
	}
	s.WriteArrayEnd()
}

// WriteBytes writes a binary value.
func (s *MarshalState) WriteBytes(v []byte) {
	if s.Err() != nil {
		return
	}
	if v == nil {
		s.WriteNil()
		return
	}
	s.WriteString(base64.StdEncoding.EncodeToString(v))
}

// WriteBytesArray writes an array of binary values.
func (s *MarshalState) WriteBytesArray(vs [][]byte) {
	if s.Err() != nil {
		return
	}
	s.WriteArrayStart()
	for i, v := range vs {
		if i > 0 {
			s.WriteMore()
		}
		s.WriteBytes(v)
	}
	s.WriteArrayEnd()
}

// WriteNil writes a null.
func (s *MarshalState) WriteNil() {
	if s.Err() != nil {
		return
	}
	s.inner.WriteNil()
}

// WriteObjectStart writes the starting { of an object.
func (s *MarshalState) WriteObjectStart() {
	if s.Err() != nil {
		return
	}
	s.inner.WriteObjectStart()
}

// WriteObjectField writes a field name and colon.
func (s *MarshalState) WriteObjectField(field string) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteObjectField(field)
}

// WriteObjectBoolField writes a field name and colon.
func (s *MarshalState) WriteObjectBoolField(field bool) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteObjectField(strconv.FormatBool(field))
}

// WriteObjectInt32Field writes a field name and colon.
func (s *MarshalState) WriteObjectInt32Field(field int32) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteObjectField(strconv.FormatInt(int64(field), 10))
}

// WriteObjectUint32Field writes a field name and colon.
func (s *MarshalState) WriteObjectUint32Field(field uint32) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteObjectField(strconv.FormatUint(uint64(field), 10))
}

// WriteObjectInt64Field writes a field name and colon.
func (s *MarshalState) WriteObjectInt64Field(field int64) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteObjectField(strconv.FormatInt(field, 10))
}

// WriteObjectUint64Field writes a field name and colon.
func (s *MarshalState) WriteObjectUint64Field(field uint64) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteObjectField(strconv.FormatUint(field, 10))
}

// WriteObjectStringField writes a field name and colon.
func (s *MarshalState) WriteObjectStringField(field string) {
	if s.Err() != nil {
		return
	}
	s.inner.WriteObjectField(field)
}

// WriteObjectEnd writes the ending } of an object.
func (s *MarshalState) WriteObjectEnd() {
	if s.Err() != nil {
		return
	}
	s.inner.WriteObjectEnd()
}

// WriteArrayStart writes the starting [ of an array.
func (s *MarshalState) WriteArrayStart() {
	if s.Err() != nil {
		return
	}
	s.inner.WriteArrayStart()
}

// WriteArrayEnd writes the ending ] of an array.
func (s *MarshalState) WriteArrayEnd() {
	if s.Err() != nil {
		return
	}
	s.inner.WriteArrayEnd()
}

// WriteMore writes a comma.
func (s *MarshalState) WriteMore() {
	if s.Err() != nil {
		return
	}
	s.inner.WriteMore()
}

// WriteMoreIf writes a comma if b is false, and sets b to true.
func (s *MarshalState) WriteMoreIf(b *bool) {
	if s.Err() != nil {
		return
	}
	if *b {
		s.WriteMore()
	}
	*b = true
}

// WriteEnum writes an enum value.
// If config.EnumsAsInts is true or a string value is not found for the value, this writes a number.
func (s *MarshalState) WriteEnum(x int32, valueMaps ...map[int32]string) {
	if s.Err() != nil {
		return
	}
	if s.config.EnumsAsInts {
		s.WriteEnumNumber(x)
	} else {
		s.WriteEnumString(x, valueMaps...)
	}
}

// GetEnumString gets the string representation of the enum using the value maps.
// If none of the value maps contains a mapping for the enum value,
// it returns the numeric value as a string.
func GetEnumString(x int32, valueMaps ...map[int32]string) string {
	for _, valueMap := range valueMaps {
		if v, ok := valueMap[x]; ok {
			return v
		}
	}
	return strconv.FormatInt(int64(x), 10)
}

// WriteEnumString writes an enum value as a string.
func (s *MarshalState) WriteEnumString(x int32, valueMaps ...map[int32]string) {
	if s.Err() != nil {
		return
	}
	s.WriteString(GetEnumString(x, valueMaps...))
}

// WriteEnumNumber writes an enum value as a number.
func (s *MarshalState) WriteEnumNumber(x int32) {
	if s.Err() != nil {
		return
	}
	s.WriteInt32(x)
}

// WriteTime writes a time value.
func (s *MarshalState) WriteTime(x time.Time) {
	if s.Err() != nil {
		return
	}
	v := x.UTC().Format("2006-01-02T15:04:05.000000000")
	// According to the protobuf spec, nanoseconds in timestamps should be written as 3, 6 or 9 digits.
	v = strings.TrimSuffix(v, "000")
	v = strings.TrimSuffix(v, "000")
	v = strings.TrimSuffix(v, ".000")
	s.inner.WriteString(v + "Z")
}

// WriteDuration writes a duration value.
func (s *MarshalState) WriteDuration(x time.Duration) {
	if s.Err() != nil {
		return
	}
	v := fmt.Sprintf("%.09f", x.Seconds())
	// According to the protobuf spec, nanoseconds in durations should be written as 3, 6 or 9 digits.
	v = strings.TrimSuffix(v, "000")
	v = strings.TrimSuffix(v, "000")
	v = strings.TrimSuffix(v, ".000")
	s.inner.WriteString(v + "s")
}

// WriteFieldMask writes a field mask value.
func (s *MarshalState) WriteFieldMask(x FieldMask) {
	if s.Err() != nil {
		return
	}
	paths := x.GetPaths()
	s.inner.WriteString(strings.Join(paths, ","))
}

func (s *MarshalState) WriteLegacyFieldMask(x FieldMask) {
	if s.Err() != nil {
		return
	}
	paths := x.GetPaths()
	s.inner.WriteObjectStart()
	s.inner.WriteObjectField("paths")
	s.inner.WriteArrayStart()
	for i, path := range paths {
		if i != 0 {
			s.inner.WriteMore()
		}
		s.inner.WriteString(path)
	}
	s.inner.WriteArrayEnd()
	s.inner.WriteObjectEnd()
}
