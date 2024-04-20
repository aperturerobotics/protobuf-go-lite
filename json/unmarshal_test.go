// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"io"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testUnmarshal(t *testing.T, f func(s *UnmarshalState) interface{}, data string, expected interface{}) {
	t.Helper()
	if expectedValue := reflect.ValueOf(expected); expectedValue.Type().Kind() == reflect.Ptr {
		expected = expectedValue.Elem().Interface()
	}
	s := NewUnmarshalState([]byte(data), DefaultUnmarshalerConfig)
	actual := f(s)
	if actualValue := reflect.ValueOf(actual); actualValue.Type().Kind() == reflect.Ptr {
		actual = actualValue.Elem().Interface()
	}
	if err := s.Err(); err != nil && err != io.EOF {
		t.Error(err)
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("diff: %s", diff)
	}
}

func TestUnmarshaler(t *testing.T) {
	// float

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadFloat32()
	}, `-12.34`, float32(-12.34))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadFloat32()
	}, `"-12.34"`, float32(-12.34))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadFloat64()
	}, `-12.34`, -12.34)

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadFloat64()
	}, `"-12.34"`, -12.34)

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadFloat32Array()
	}, `[-12.34,56.78]`, []float32{-12.34, 56.78})

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadFloat64Array()
	}, `[-12.34,56.78]`, []float64{-12.34, 56.78})

	// int

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadInt32()
	}, `-12`, int32(-12))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadInt32()
	}, `"-12"`, int32(-12))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadInt64()
	}, `-12`, int64(-12))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadInt64()
	}, `"-12"`, int64(-12))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadInt32Array()
	}, `[-12,34]`, []int32{-12, 34})

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadInt64Array()
	}, `["-12","34"]`, []int64{-12, 34})

	// uint

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadUint32()
	}, `12`, uint32(12))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadUint32()
	}, `"12"`, uint32(12))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadUint64()
	}, `12`, uint64(12))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadUint64()
	}, `"12"`, uint64(12))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadUint32Array()
	}, `[12,34]`, []uint32{12, 34})

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadUint64Array()
	}, `["12","34"]`, []uint64{12, 34})

	// bool

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadBool()
	}, `true`, true)

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadBool()
	}, `false`, false)

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadBoolArray()
	}, `[true,false]`, []bool{true, false})

	// string

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadString()
	}, `"foo"`, "foo")

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadStringArray()
	}, `["foo","bar"]`, []string{"foo", "bar"})

	// bytes

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadBytes()
	}, `"Zm9vYg=="`, []byte("foob"))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadBytes()
	}, `"YXI="`, []byte("ar"))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadBytes()
	}, `"YXI"`, []byte("ar"))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadBytes()
	}, `"cXM/"`, []byte("qs?"))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadBytes()
	}, `"cXM_"`, []byte("qs?"))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadBytes()
	}, `"cHM+"`, []byte("ps>"))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadBytes()
	}, `"cHM-"`, []byte("ps>"))

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadBytesArray()
	}, `["Zm9vYg==","YXI=","cXM/","cHM+"]`, [][]byte{[]byte("foob"), []byte("ar"), []byte("qs?"), []byte("ps>")})

	// nil

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadNil()
	}, `null`, true)

	// time

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadTime()
	}, `"2006-01-02T08:04:05.123456789Z"`, testTime.UTC())
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadTime()
	}, `"2006-01-02T08:04:05.123456780Z"`, testTime.UTC().Truncate(10))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadTime()
	}, `"2006-01-02T08:04:05.123456700Z"`, testTime.UTC().Truncate(100))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadTime()
	}, `"2006-01-02T08:04:05.123456Z"`, testTime.UTC().Truncate(1000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadTime()
	}, `"2006-01-02T08:04:05.123450Z"`, testTime.UTC().Truncate(10000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadTime()
	}, `"2006-01-02T08:04:05.123400Z"`, testTime.UTC().Truncate(100000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadTime()
	}, `"2006-01-02T08:04:05.123Z"`, testTime.UTC().Truncate(1000000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadTime()
	}, `"2006-01-02T08:04:05.120Z"`, testTime.UTC().Truncate(10000000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadTime()
	}, `"2006-01-02T08:04:05.100Z"`, testTime.UTC().Truncate(100000000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadTime()
	}, `"2006-01-02T08:04:05Z"`, testTime.UTC().Truncate(1000000000))

	// duration

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadDuration()
	}, `"3723.123456789s"`, testDuration)
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadDuration()
	}, `"3723.123456780s"`, testDuration.Truncate(10))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadDuration()
	}, `"3723.123456700s"`, testDuration.Truncate(100))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadDuration()
	}, `"3723.123456s"`, testDuration.Truncate(1000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadDuration()
	}, `"3723.123450s"`, testDuration.Truncate(10000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadDuration()
	}, `"3723.123400s"`, testDuration.Truncate(100000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadDuration()
	}, `"3723.123s"`, testDuration.Truncate(1000000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadDuration()
	}, `"3723.120s"`, testDuration.Truncate(10000000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadDuration()
	}, `"3723.100s"`, testDuration.Truncate(100000000))
	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadDuration()
	}, `"3723s"`, testDuration.Truncate(1000000000))

	// field mask

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadFieldMask().GetPaths()
	}, `"foo,bar,baz.qux"`, []string{"foo", "bar", "baz.qux"})

	testUnmarshal(t, func(s *UnmarshalState) interface{} {
		return s.ReadFieldMask().GetPaths()
	}, `{"paths":["foo","bar","baz.qux"]}`, []string{"foo", "bar", "baz.qux"})
}
