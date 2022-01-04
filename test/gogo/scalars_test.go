package test_test

import (
	"testing"

	. "github.com/TheThingsIndustries/protoc-gen-go-json/test/gogo"
)

var testMessagesWithScalars = []struct {
	name         string
	msg          MessageWithScalars
	expected     string
	expectedMask []string
}{
	{
		name:     "empty",
		msg:      MessageWithScalars{},
		expected: `{}`,
	},
	{
		name: "zero",
		msg:  MessageWithScalars{},
		expected: `{
			"double_value": 0,
			"double_values": [],
			"float_value": 0,
			"float_values": [],
			"int32_value": 0,
			"int32_values": [],
			"int64_value": "0",
			"int64_values": [],
			"uint32_value": 0,
			"uint32_values": [],
			"uint64_value": "0",
			"uint64_values": [],
			"sint32_value": 0,
			"sint32_values": [],
			"sint64_value": "0",
			"sint64_values": [],
			"fixed32_value": 0,
			"fixed32_values": [],
			"fixed64_value": "0",
			"fixed64_values": [],
			"sfixed32_value": 0,
			"sfixed32_values": [],
			"sfixed64_value": "0",
			"sfixed64_values": [],
			"bool_value": false,
			"bool_values": [],
			"string_value": "",
			"string_values": [],
			"bytes_value": null,
			"bytes_values": [],
			"hex_bytes_value": null,
			"hex_bytes_values": []
		}`,
		expectedMask: []string{
			"double_value",
			"double_values",
			"float_value",
			"float_values",
			"int32_value",
			"int32_values",
			"int64_value",
			"int64_values",
			"uint32_value",
			"uint32_values",
			"uint64_value",
			"uint64_values",
			"sint32_value",
			"sint32_values",
			"sint64_value",
			"sint64_values",
			"fixed32_value",
			"fixed32_values",
			"fixed64_value",
			"fixed64_values",
			"sfixed32_value",
			"sfixed32_values",
			"sfixed64_value",
			"sfixed64_values",
			"bool_value",
			"bool_values",
			"string_value",
			"string_values",
			"bytes_value",
			"bytes_values",
			"hex_bytes_value",
			"hex_bytes_values",
		},
	},
	{
		name: "full",
		msg: MessageWithScalars{
			DoubleValue:    12.34,
			DoubleValues:   []float64{12.34, 56.78},
			FloatValue:     12.34,
			FloatValues:    []float32{12.34, 56.78},
			Int32Value:     -42,
			Int32Values:    []int32{1, 2, -42},
			Int64Value:     -42,
			Int64Values:    []int64{1, 2, -42},
			Uint32Value:    42,
			Uint32Values:   []uint32{1, 2, 42},
			Uint64Value:    42,
			Uint64Values:   []uint64{1, 2, 42},
			Sint32Value:    -42,
			Sint32Values:   []int32{1, 2, -42},
			Sint64Value:    -42,
			Sint64Values:   []int64{1, 2, -42},
			Fixed32Value:   42,
			Fixed32Values:  []uint32{1, 2, 42},
			Fixed64Value:   42,
			Fixed64Values:  []uint64{1, 2, 42},
			Sfixed32Value:  -42,
			Sfixed32Values: []int32{1, 2, -42},
			Sfixed64Value:  -42,
			Sfixed64Values: []int64{1, 2, -42},
			BoolValue:      true,
			BoolValues:     []bool{true, false},
			StringValue:    "foo",
			StringValues:   []string{"foo", "bar"},
			BytesValue:     []byte("foo"),
			BytesValues:    [][]byte{[]byte("foo"), []byte("bar")},
			HexBytesValue:  []byte("foo"),
			HexBytesValues: [][]byte{[]byte("foo"), []byte("bar")},
		},
		expected: `{
			"double_value": 12.34,
			"double_values": [12.34, 56.78],
			"float_value": 12.34,
			"float_values": [12.34, 56.78],
			"int32_value": -42,
			"int32_values": [1, 2, -42],
			"int64_value": "-42",
			"int64_values": ["1", "2", "-42"],
			"uint32_value": 42,
			"uint32_values": [1, 2, 42],
			"uint64_value": "42",
			"uint64_values": ["1", "2", "42"],
			"sint32_value": -42,
			"sint32_values": [1, 2, -42],
			"sint64_value": "-42",
			"sint64_values": ["1", "2", "-42"],
			"fixed32_value": 42,
			"fixed32_values": [1, 2, 42],
			"fixed64_value": "42",
			"fixed64_values": ["1", "2", "42"],
			"sfixed32_value": -42,
			"sfixed32_values": [1, 2, -42],
			"sfixed64_value": "-42",
			"sfixed64_values": ["1", "2", "-42"],
			"bool_value": true,
			"bool_values": [true, false],
			"string_value": "foo",
			"string_values": ["foo", "bar"],
			"bytes_value": "Zm9v",
			"bytes_values": ["Zm9v", "YmFy"],
			"hex_bytes_value": "666F6F",
			"hex_bytes_values": ["666F6F", "626172"]
		}`,
		expectedMask: []string{
			"double_value",
			"double_values",
			"float_value",
			"float_values",
			"int32_value",
			"int32_values",
			"int64_value",
			"int64_values",
			"uint32_value",
			"uint32_values",
			"uint64_value",
			"uint64_values",
			"sint32_value",
			"sint32_values",
			"sint64_value",
			"sint64_values",
			"fixed32_value",
			"fixed32_values",
			"fixed64_value",
			"fixed64_values",
			"sfixed32_value",
			"sfixed32_values",
			"sfixed64_value",
			"sfixed64_values",
			"bool_value",
			"bool_values",
			"string_value",
			"string_values",
			"bytes_value",
			"bytes_values",
			"hex_bytes_value",
			"hex_bytes_values",
		},
	},
}

func TestMarshalMessageWithScalars(t *testing.T) {
	for _, tt := range testMessagesWithScalars {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithScalars(t *testing.T) {
	for _, tt := range testMessagesWithScalars {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}

var testMessagesWithOneofScalars = []struct {
	name         string
	msg          MessageWithOneofScalars
	expected     string
	expectedMask []string
}{
	{
		name:     "empty",
		msg:      MessageWithOneofScalars{},
		expected: `{}`,
	},
	{
		name: "double_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_DoubleValue{},
		},
		expected:     `{"double_value": 0}`,
		expectedMask: []string{"double_value"},
	},
	{
		name: "double_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_DoubleValue{DoubleValue: 12.34},
		},
		expected:     `{"double_value": 12.34}`,
		expectedMask: []string{"double_value"},
	},
	{
		name: "float_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_FloatValue{},
		},
		expected:     `{"float_value": 0}`,
		expectedMask: []string{"float_value"},
	},
	{
		name: "float_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_FloatValue{FloatValue: 12.34},
		},
		expected:     `{"float_value": 12.34}`,
		expectedMask: []string{"float_value"},
	},
	{
		name: "int32_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Int32Value{},
		},
		expected:     `{"int32_value": 0}`,
		expectedMask: []string{"int32_value"},
	},
	{
		name: "int32_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Int32Value{Int32Value: -42},
		},
		expected:     `{"int32_value": -42}`,
		expectedMask: []string{"int32_value"},
	},
	{
		name: "int64_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Int64Value{},
		},
		expected:     `{"int64_value": "0"}`,
		expectedMask: []string{"int64_value"},
	},
	{
		name: "int64_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Int64Value{Int64Value: -42},
		},
		expected:     `{"int64_value": "-42"}`,
		expectedMask: []string{"int64_value"},
	},
	{
		name: "uint32_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Uint32Value{},
		},
		expected:     `{"uint32_value": 0}`,
		expectedMask: []string{"uint32_value"},
	},
	{
		name: "uint32_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Uint32Value{Uint32Value: 42},
		},
		expected:     `{"uint32_value": 42}`,
		expectedMask: []string{"uint32_value"},
	},
	{
		name: "uint64_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Uint64Value{},
		},
		expected:     `{"uint64_value": "0"}`,
		expectedMask: []string{"uint64_value"},
	},
	{
		name: "uint64_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Uint64Value{Uint64Value: 42},
		},
		expected:     `{"uint64_value": "42"}`,
		expectedMask: []string{"uint64_value"},
	},
	{
		name: "sint32_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Sint32Value{},
		},
		expected:     `{"sint32_value": 0}`,
		expectedMask: []string{"sint32_value"},
	},
	{
		name: "sint32_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Sint32Value{Sint32Value: -42},
		},
		expected:     `{"sint32_value": -42}`,
		expectedMask: []string{"sint32_value"},
	},
	{
		name: "sint64_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Sint64Value{},
		},
		expected:     `{"sint64_value": "0"}`,
		expectedMask: []string{"sint64_value"},
	},
	{
		name: "sint64_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Sint64Value{Sint64Value: -42},
		},
		expected:     `{"sint64_value": "-42"}`,
		expectedMask: []string{"sint64_value"},
	},
	{
		name: "fixed32_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Fixed32Value{},
		},
		expected:     `{"fixed32_value": 0}`,
		expectedMask: []string{"fixed32_value"},
	},
	{
		name: "fixed32_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Fixed32Value{Fixed32Value: 42},
		},
		expected:     `{"fixed32_value": 42}`,
		expectedMask: []string{"fixed32_value"},
	},
	{
		name: "fixed64_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Fixed64Value{},
		},
		expected:     `{"fixed64_value": "0"}`,
		expectedMask: []string{"fixed64_value"},
	},
	{
		name: "fixed64_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Fixed64Value{Fixed64Value: 42},
		},
		expected:     `{"fixed64_value": "42"}`,
		expectedMask: []string{"fixed64_value"},
	},

	{
		name: "sfixed32_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Sfixed32Value{},
		},
		expected:     `{"sfixed32_value": 0}`,
		expectedMask: []string{"sfixed32_value"},
	},
	{
		name: "sfixed32_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Sfixed32Value{Sfixed32Value: -42},
		},
		expected:     `{"sfixed32_value": -42}`,
		expectedMask: []string{"sfixed32_value"},
	},
	{
		name: "sfixed64_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Sfixed64Value{},
		},
		expected:     `{"sfixed64_value": "0"}`,
		expectedMask: []string{"sfixed64_value"},
	},
	{
		name: "sfixed64_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_Sfixed64Value{Sfixed64Value: -42},
		},
		expected:     `{"sfixed64_value": "-42"}`,
		expectedMask: []string{"sfixed64_value"},
	},

	{
		name: "bool_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_BoolValue{},
		},
		expected:     `{"bool_value": false}`,
		expectedMask: []string{"bool_value"},
	},
	{
		name: "bool_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_BoolValue{BoolValue: true},
		},
		expected:     `{"bool_value": true}`,
		expectedMask: []string{"bool_value"},
	},
	{
		name: "string_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_StringValue{},
		},
		expected:     `{"string_value": ""}`,
		expectedMask: []string{"string_value"},
	},
	{
		name: "string_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_StringValue{StringValue: "foo"},
		},
		expected:     `{"string_value": "foo"}`,
		expectedMask: []string{"string_value"},
	},
	{
		name: "bytes_null",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_BytesValue{},
		},
		expected:     `{"bytes_value": null}`,
		expectedMask: []string{"bytes_value"},
	},
	{
		name: "bytes_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_BytesValue{BytesValue: []byte{}},
		},
		expected:     `{"bytes_value": ""}`,
		expectedMask: []string{"bytes_value"},
	},
	{
		name: "bytes_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_BytesValue{BytesValue: []byte("foo")},
		},
		expected:     `{"bytes_value": "Zm9v"}`,
		expectedMask: []string{"bytes_value"},
	},
	{
		name: "hex_bytes_null",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_HexBytesValue{},
		},
		expected:     `{"hex_bytes_value": null}`,
		expectedMask: []string{"hex_bytes_value"},
	},
	{
		name: "hex_bytes_zero",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_HexBytesValue{HexBytesValue: []byte{}},
		},
		expected:     `{"hex_bytes_value": ""}`,
		expectedMask: []string{"hex_bytes_value"},
	},
	{
		name: "hex_bytes_value",
		msg: MessageWithOneofScalars{
			Value: &MessageWithOneofScalars_HexBytesValue{HexBytesValue: []byte("foo")},
		},
		expected:     `{"hex_bytes_value": "666F6F"}`,
		expectedMask: []string{"hex_bytes_value"},
	},
}

func TestMarshalMessageWithOneofScalars(t *testing.T) {
	for _, tt := range testMessagesWithOneofScalars {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithOneofScalars(t *testing.T) {
	for _, tt := range testMessagesWithOneofScalars {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}

var testMessagesWithScalarMaps = []struct {
	name         string
	msg          MessageWithScalarMaps
	expected     string
	expectedMask []string
}{
	{
		name:     "empty",
		msg:      MessageWithScalarMaps{},
		expected: `{}`,
	},
	{
		name: "full",
		msg: MessageWithScalarMaps{
			StringDoubleMap:   map[string]float64{"value": -42},
			StringFloatMap:    map[string]float32{"value": -42},
			StringInt32Map:    map[string]int32{"value": -42},
			Int32StringMap:    map[int32]string{-42: "answer"},
			StringInt64Map:    map[string]int64{"value": -42},
			Int64StringMap:    map[int64]string{-42: "answer"},
			StringUint32Map:   map[string]uint32{"value": 42},
			Uint32StringMap:   map[uint32]string{42: "answer"},
			StringUint64Map:   map[string]uint64{"value": 42},
			Uint64StringMap:   map[uint64]string{42: "answer"},
			StringSint32Map:   map[string]int32{"value": -42},
			Sint32StringMap:   map[int32]string{-42: "answer"},
			StringSint64Map:   map[string]int64{"value": -42},
			Sint64StringMap:   map[int64]string{-42: "answer"},
			StringFixed32Map:  map[string]uint32{"value": 42},
			Fixed32StringMap:  map[uint32]string{42: "answer"},
			StringFixed64Map:  map[string]uint64{"value": 42},
			Fixed64StringMap:  map[uint64]string{42: "answer"},
			StringSfixed32Map: map[string]int32{"value": -42},
			Sfixed32StringMap: map[int32]string{-42: "answer"},
			StringSfixed64Map: map[string]int64{"value": -42},
			Sfixed64StringMap: map[int64]string{-42: "answer"},
			StringBoolMap:     map[string]bool{"yes": true},
			BoolStringMap:     map[bool]string{true: "yes"},
			StringStringMap:   map[string]string{"value": "foo"},
			StringBytesMap:    map[string][]byte{"value": []byte("foo")},
			StringHexBytesMap: map[string][]byte{"value": []byte("foo")},
		},
		expected: `{
			"string_double_map": {"value": -42},
			"string_float_map": {"value": -42},
			"string_int32_map": {"value": -42},
			"int32_string_map": {"-42": "answer"},
			"string_int64_map": {"value": "-42"},
			"int64_string_map": {"-42": "answer"},
			"string_uint32_map": {"value": 42},
			"uint32_string_map": {"42": "answer"},
			"string_uint64_map": {"value": "42"},
			"uint64_string_map": {"42": "answer"},
			"string_sint32_map": {"value": -42},
			"sint32_string_map": {"-42": "answer"},
			"string_sint64_map": {"value": "-42"},
			"sint64_string_map": {"-42": "answer"},
			"string_fixed32_map": {"value": 42},
			"fixed32_string_map": {"42": "answer"},
			"string_fixed64_map": {"value": "42"},
			"fixed64_string_map": {"42": "answer"},
			"string_sfixed32_map": {"value": -42},
			"sfixed32_string_map": {"-42": "answer"},
			"string_sfixed64_map": {"value": "-42"},
			"sfixed64_string_map": {"-42": "answer"},
			"string_bool_map": {"yes": true},
			"bool_string_map": {"true": "yes"},
			"string_string_map": {"value": "foo"},
			"string_bytes_map": {"value": "Zm9v"},
			"string_hex_bytes_map": {"value": "666F6F"}
		}`,
		expectedMask: []string{
			"string_double_map",
			"string_float_map",
			"string_int32_map",
			"int32_string_map",
			"string_int64_map",
			"int64_string_map",
			"string_uint32_map",
			"uint32_string_map",
			"string_uint64_map",
			"uint64_string_map",
			"string_sint32_map",
			"sint32_string_map",
			"string_sint64_map",
			"sint64_string_map",
			"string_fixed32_map",
			"fixed32_string_map",
			"string_fixed64_map",
			"fixed64_string_map",
			"string_sfixed32_map",
			"sfixed32_string_map",
			"string_sfixed64_map",
			"sfixed64_string_map",
			"string_bool_map",
			"bool_string_map",
			"string_string_map",
			"string_bytes_map",
			"string_hex_bytes_map",
		},
	},
}

func TestMarshalMessageWithScalarMaps(t *testing.T) {
	for _, tt := range testMessagesWithScalarMaps {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithScalarMaps(t *testing.T) {
	for _, tt := range testMessagesWithScalarMaps {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}
