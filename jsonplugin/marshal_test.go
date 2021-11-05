// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package jsonplugin_test

import (
	"math"
	"testing"
	"time"

	. "github.com/TheThingsIndustries/protoc-gen-go-json/jsonplugin"
	"github.com/google/go-cmp/cmp"
)

var (
	testTime     = time.Date(2006, time.January, 2, 15, 4, 5, 123456789, time.FixedZone("07:00", 7*3600))
	testDuration = time.Hour + 2*time.Minute + 3*time.Second + 123456789
)

func testMarshal(t *testing.T, f func(s *MarshalState), expected string) {
	t.Helper()

	s := NewMarshalState(MarshalerConfig{})
	f(s)
	data, err := s.Bytes()
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(expected, string(data)); diff != "" {
		t.Errorf("diff: %s", diff)
	}
}

func TestMarshaler(t *testing.T) {
	// float

	testMarshal(t, func(s *MarshalState) {
		s.WriteFloat32(-12.34)
	}, `-12.34`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteFloat64(-12.34)
	}, `-12.34`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteFloat32Array([]float32{-12.34, 56.78})
	}, `[-12.34,56.78]`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteFloat64Array([]float64{-12.34, 56.78})
	}, `[-12.34,56.78]`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteFloat32(float32(math.NaN()))
	}, `"NaN"`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteFloat64(math.NaN())
	}, `"NaN"`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteFloat32(float32(math.Inf(-1)))
	}, `"-Infinity"`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteFloat64(math.Inf(-1))
	}, `"-Infinity"`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteFloat32(float32(math.Inf(1)))
	}, `"Infinity"`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteFloat64(math.Inf(1))
	}, `"Infinity"`)

	// int

	testMarshal(t, func(s *MarshalState) {
		s.WriteInt32(-12)
	}, `-12`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteInt64(-12)
	}, `"-12"`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteInt32Array([]int32{-12, 34})
	}, `[-12,34]`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteInt64Array([]int64{-12, 34})
	}, `["-12","34"]`)

	// uint

	testMarshal(t, func(s *MarshalState) {
		s.WriteUint32(12)
	}, `12`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteUint64(12)
	}, `"12"`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteUint32Array([]uint32{12, 34})
	}, `[12,34]`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteUint64Array([]uint64{12, 34})
	}, `["12","34"]`)

	// bool

	testMarshal(t, func(s *MarshalState) {
		s.WriteBool(true)
	}, `true`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteBool(false)
	}, `false`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteBoolArray([]bool{true, false})
	}, `[true,false]`)

	// string

	testMarshal(t, func(s *MarshalState) {
		s.WriteString("foo")
	}, `"foo"`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteStringArray([]string{"foo", "bar"})
	}, `["foo","bar"]`)

	// bytes

	testMarshal(t, func(s *MarshalState) {
		s.WriteBytes([]byte("foob"))
	}, `"Zm9vYg=="`)

	testMarshal(t, func(s *MarshalState) {
		s.WriteBytesArray([][]byte{[]byte("foob"), []byte("ar"), []byte("qs?"), []byte("ps>")})
	}, `["Zm9vYg==","YXI=","cXM/","cHM+"]`)

	// nil

	testMarshal(t, func(s *MarshalState) {
		s.WriteNil()
	}, `null`)

	// time

	testMarshal(t, func(s *MarshalState) { s.WriteTime(testTime) }, `"2006-01-02T08:04:05.123456789Z"`)
	testMarshal(t, func(s *MarshalState) { s.WriteTime(testTime.Truncate(10)) }, `"2006-01-02T08:04:05.123456780Z"`)
	testMarshal(t, func(s *MarshalState) { s.WriteTime(testTime.Truncate(100)) }, `"2006-01-02T08:04:05.123456700Z"`)
	testMarshal(t, func(s *MarshalState) { s.WriteTime(testTime.Truncate(1000)) }, `"2006-01-02T08:04:05.123456Z"`)
	testMarshal(t, func(s *MarshalState) { s.WriteTime(testTime.Truncate(10000)) }, `"2006-01-02T08:04:05.123450Z"`)
	testMarshal(t, func(s *MarshalState) { s.WriteTime(testTime.Truncate(100000)) }, `"2006-01-02T08:04:05.123400Z"`)
	testMarshal(t, func(s *MarshalState) { s.WriteTime(testTime.Truncate(1000000)) }, `"2006-01-02T08:04:05.123Z"`)
	testMarshal(t, func(s *MarshalState) { s.WriteTime(testTime.Truncate(10000000)) }, `"2006-01-02T08:04:05.120Z"`)
	testMarshal(t, func(s *MarshalState) { s.WriteTime(testTime.Truncate(100000000)) }, `"2006-01-02T08:04:05.100Z"`)
	testMarshal(t, func(s *MarshalState) { s.WriteTime(testTime.Truncate(1000000000)) }, `"2006-01-02T08:04:05Z"`)

	// duration

	testMarshal(t, func(s *MarshalState) { s.WriteDuration(testDuration) }, `"3723.123456789s"`)
	testMarshal(t, func(s *MarshalState) { s.WriteDuration(testDuration.Truncate(10)) }, `"3723.123456780s"`)
	testMarshal(t, func(s *MarshalState) { s.WriteDuration(testDuration.Truncate(100)) }, `"3723.123456700s"`)
	testMarshal(t, func(s *MarshalState) { s.WriteDuration(testDuration.Truncate(1000)) }, `"3723.123456s"`)
	testMarshal(t, func(s *MarshalState) { s.WriteDuration(testDuration.Truncate(10000)) }, `"3723.123450s"`)
	testMarshal(t, func(s *MarshalState) { s.WriteDuration(testDuration.Truncate(100000)) }, `"3723.123400s"`)
	testMarshal(t, func(s *MarshalState) { s.WriteDuration(testDuration.Truncate(1000000)) }, `"3723.123s"`)
	testMarshal(t, func(s *MarshalState) { s.WriteDuration(testDuration.Truncate(10000000)) }, `"3723.120s"`)
	testMarshal(t, func(s *MarshalState) { s.WriteDuration(testDuration.Truncate(100000000)) }, `"3723.100s"`)
	testMarshal(t, func(s *MarshalState) { s.WriteDuration(testDuration.Truncate(1000000000)) }, `"3723s"`)
}
