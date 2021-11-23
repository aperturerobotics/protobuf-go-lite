package test_test

import (
	"testing"
	"time"

	. "github.com/TheThingsIndustries/protoc-gen-go-json/test/golang"
	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	testTime     = time.Date(2006, time.January, 2, 15, 4, 5, 123456789, time.FixedZone("07:00", 7*3600))
	testDuration = time.Hour + 2*time.Minute + 3*time.Second + 123456789
)

func TestTimeDuration(t *testing.T) {
	expectedTime := "2006-01-02T15:04:05.123456789+07:00"
	if actualTime := testTime.Format(time.RFC3339Nano); actualTime != expectedTime {
		t.Fatalf("expected timestamp %s, got %s", expectedTime, actualTime)
	}
	expectedDuration := "1h2m3.123456789s"
	if actualDuration := testDuration.String(); actualDuration != expectedDuration {
		t.Fatalf("expected timestamp %s, got %s", expectedDuration, actualDuration)
	}
}

func mustTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func mustDuration(t time.Duration) *durationpb.Duration {
	return durationpb.New(t)
}

func mustAny(pb proto.Message) *anypb.Any {
	any, err := anypb.New(pb)
	if err != nil {
		panic(err)
	}
	return any
}

var testMessagesWithWKTs = []struct {
	name         string
	msg          MessageWithWKTs
	expected     string
	expectedMask []string
}{
	{
		name:     "empty",
		msg:      MessageWithWKTs{},
		expected: `{}`,
	},
	{
		name: "zero",
		msg:  MessageWithWKTs{},
		expected: `{
			"double_value": null,
			"double_values": [],
			"float_value": null,
			"float_values": [],
			"int32_value": null,
			"int32_values": [],
			"int64_value": null,
			"int64_values": [],
			"uint32_value": null,
			"uint32_values": [],
			"uint64_value": null,
			"uint64_values": [],
			"bool_value": null,
			"bool_values": [],
			"string_value": null,
			"string_values": [],
			"bytes_value": null,
			"bytes_values": [],
			"empty_values": [],
			"timestamp_value": null,
			"timestamp_values": [],
			"duration_value": null,
			"duration_values": [],
			"field_mask_value": null,
			"field_mask_values": [],
			"value_values": [],
			"list_value_value": null,
			"list_value_values": [],
			"struct_value": null,
			"struct_values": [],
			"any_value": null,
			"any_values": []
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
			"bool_value",
			"bool_values",
			"string_value",
			"string_values",
			"bytes_value",
			"bytes_values",
			"empty_values",
			"timestamp_value",
			"timestamp_values",
			"duration_value",
			"duration_values",
			"field_mask_value",
			"field_mask_values",
			"value_values",
			"list_value_value",
			"list_value_values",
			"struct_value",
			"struct_values",
			"any_value",
			"any_values",
		},
	},
	{
		name: "full",
		msg: MessageWithWKTs{
			DoubleValue: &wrapperspb.DoubleValue{Value: 12.34},
			DoubleValues: []*wrapperspb.DoubleValue{
				{Value: 12.34},
				{Value: 56.78},
			},
			FloatValue: &wrapperspb.FloatValue{Value: 12.34},
			FloatValues: []*wrapperspb.FloatValue{
				{Value: 12.34},
				{Value: 56.78},
			},
			Int32Value: &wrapperspb.Int32Value{Value: -42},
			Int32Values: []*wrapperspb.Int32Value{
				{Value: 1},
				{Value: 2},
				{Value: -42},
			},
			Int64Value: &wrapperspb.Int64Value{Value: -42},
			Int64Values: []*wrapperspb.Int64Value{
				{Value: 1},
				{Value: 2},
				{Value: -42},
			},
			Uint32Value: &wrapperspb.UInt32Value{Value: 42},
			Uint32Values: []*wrapperspb.UInt32Value{
				{Value: 1},
				{Value: 2},
				{Value: 42},
			},
			Uint64Value: &wrapperspb.UInt64Value{Value: 42},
			Uint64Values: []*wrapperspb.UInt64Value{
				{Value: 1},
				{Value: 2},
				{Value: 42},
			},
			BoolValue: &wrapperspb.BoolValue{Value: true},
			BoolValues: []*wrapperspb.BoolValue{
				{Value: true},
				{Value: false},
			},
			StringValue: &wrapperspb.StringValue{Value: "foo"},
			StringValues: []*wrapperspb.StringValue{
				{Value: "foo"},
				{Value: "bar"},
			},
			BytesValue: &wrapperspb.BytesValue{Value: []byte("foo")},
			BytesValues: []*wrapperspb.BytesValue{
				{Value: []byte("foo")},
				{Value: []byte("bar")},
			},
			EmptyValue:     &emptypb.Empty{},
			EmptyValues:    []*emptypb.Empty{{}, {}},
			TimestampValue: mustTimestamp(testTime),
			TimestampValues: []*timestamppb.Timestamp{
				mustTimestamp(testTime),
				mustTimestamp(testTime.Truncate(10)),
				mustTimestamp(testTime.Truncate(100)),
				mustTimestamp(testTime.Truncate(1000)),
				mustTimestamp(testTime.Truncate(10000)),
				mustTimestamp(testTime.Truncate(100000)),
				mustTimestamp(testTime.Truncate(1000000)),
				mustTimestamp(testTime.Truncate(10000000)),
				mustTimestamp(testTime.Truncate(100000000)),
				mustTimestamp(testTime.Truncate(1000000000)),
			},
			DurationValue: mustDuration(testDuration),
			DurationValues: []*durationpb.Duration{
				mustDuration(testDuration),
				mustDuration(testDuration.Truncate(10)),
				mustDuration(testDuration.Truncate(100)),
				mustDuration(testDuration.Truncate(1000)),
				mustDuration(testDuration.Truncate(10000)),
				mustDuration(testDuration.Truncate(100000)),
				mustDuration(testDuration.Truncate(1000000)),
				mustDuration(testDuration.Truncate(10000000)),
				mustDuration(testDuration.Truncate(100000000)),
				mustDuration(testDuration.Truncate(1000000000)),
			},
			FieldMaskValue: &fieldmaskpb.FieldMask{Paths: []string{"foo.bar", "bar", "baz.qux"}},
			FieldMaskValues: []*fieldmaskpb.FieldMask{
				{Paths: []string{"foo.bar", "bar", "baz.qux"}},
			},
			ValueValue: &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "foo"}},
			ValueValues: []*structpb.Value{
				{Kind: &structpb.Value_NullValue{}},
				{Kind: &structpb.Value_NumberValue{NumberValue: 12.34}},
				{Kind: &structpb.Value_StringValue{StringValue: "foo"}},
				{Kind: &structpb.Value_BoolValue{BoolValue: true}},
			},
			ListValueValue: &structpb.ListValue{
				Values: []*structpb.Value{
					{Kind: &structpb.Value_NullValue{}},
					{Kind: &structpb.Value_NumberValue{NumberValue: 12.34}},
					{Kind: &structpb.Value_StringValue{StringValue: "foo"}},
					{Kind: &structpb.Value_BoolValue{BoolValue: true}},
				},
			},
			ListValueValues: []*structpb.ListValue{
				{
					Values: []*structpb.Value{
						{Kind: &structpb.Value_NullValue{}},
						{Kind: &structpb.Value_NumberValue{NumberValue: 12.34}},
						{Kind: &structpb.Value_StringValue{StringValue: "foo"}},
						{Kind: &structpb.Value_BoolValue{BoolValue: true}},
					},
				},
			},
			StructValue: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"null":   {Kind: &structpb.Value_NullValue{}},
					"number": {Kind: &structpb.Value_NumberValue{NumberValue: 12.34}},
					"string": {Kind: &structpb.Value_StringValue{StringValue: "foo"}},
					"bool":   {Kind: &structpb.Value_BoolValue{BoolValue: true}},
				},
			},
			StructValues: []*structpb.Struct{
				{Fields: map[string]*structpb.Value{"null": {Kind: &structpb.Value_NullValue{}}}},
				{Fields: map[string]*structpb.Value{"number": {Kind: &structpb.Value_NumberValue{NumberValue: 12.34}}}},
				{Fields: map[string]*structpb.Value{"string": {Kind: &structpb.Value_StringValue{StringValue: "foo"}}}},
				{Fields: map[string]*structpb.Value{"bool": {Kind: &structpb.Value_BoolValue{BoolValue: true}}}},
			},
			AnyValue: mustAny(&MessageWithMarshaler{Message: "hello"}),
			AnyValues: []*anypb.Any{
				mustAny(&MessageWithMarshaler{Message: "hello"}),
				mustAny(&MessageWithoutMarshaler{Message: "hello"}),
				mustAny(mustTimestamp(testTime)),
				mustAny(mustDuration(testDuration)),
				mustAny(&fieldmaskpb.FieldMask{Paths: []string{"foo.bar", "bar", "baz.qux"}}),
				mustAny(&structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "foo"}}),
			},
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
			"bool_value": true,
			"bool_values": [true, false],
			"string_value": "foo",
			"string_values": ["foo", "bar"],
			"bytes_value": "Zm9v",
			"bytes_values": ["Zm9v", "YmFy"],
			"empty_value": {},
			"empty_values": [{}, {}],
			"timestamp_value": "2006-01-02T08:04:05.123456789Z",
			"timestamp_values": [
				"2006-01-02T08:04:05.123456789Z",
				"2006-01-02T08:04:05.123456780Z",
				"2006-01-02T08:04:05.123456700Z",
				"2006-01-02T08:04:05.123456Z",
				"2006-01-02T08:04:05.123450Z",
				"2006-01-02T08:04:05.123400Z",
				"2006-01-02T08:04:05.123Z",
				"2006-01-02T08:04:05.120Z",
				"2006-01-02T08:04:05.100Z",
				"2006-01-02T08:04:05Z"
			],
			"duration_value": "3723.123456789s",
			"duration_values": [
				"3723.123456789s",
				"3723.123456780s",
				"3723.123456700s",
				"3723.123456s",
				"3723.123450s",
				"3723.123400s",
				"3723.123s",
				"3723.120s",
				"3723.100s",
				"3723s"
			],
			"field_mask_value": "foo.bar,bar,baz.qux",
			"field_mask_values": ["foo.bar,bar,baz.qux"],
			"value_value": "foo",
			"value_values": [null, 12.34, "foo", true],
			"list_value_value": [null, 12.34, "foo", true],
			"list_value_values": [[null, 12.34, "foo", true]],
			"struct_value": {
				"bool": true,
				"null": null,
				"number": 12.34,
				"string": "foo"
			},
			"struct_values": [
				{"null": null},
				{"number": 12.34},
				{"string": "foo"},
				{"bool": true}
			],
			"any_value": {
				"@type": "type.googleapis.com/thethings.json.test.MessageWithMarshaler",
				"message": "hello"
			},
			"any_values": [
				{
					"@type": "type.googleapis.com/thethings.json.test.MessageWithMarshaler",
					"message": "hello"
				},
				{
					"@type": "type.googleapis.com/thethings.json.test.MessageWithoutMarshaler",
					"message": "hello"
				},
				{
					"@type":"type.googleapis.com/google.protobuf.Timestamp",
					"value":"2006-01-02T08:04:05.123456789Z"
				},
				{
					"@type":"type.googleapis.com/google.protobuf.Duration",
					"value":"3723.123456789s"
				},
				{
					"@type":"type.googleapis.com/google.protobuf.FieldMask",
					"value": "foo.bar,bar,baz.qux"
				},
				{
					"@type":"type.googleapis.com/google.protobuf.Value",
					"value":"foo"
				}
			]
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
			"bool_value",
			"bool_values",
			"string_value",
			"string_values",
			"bytes_value",
			"bytes_values",
			"empty_value",
			"empty_values",
			"timestamp_value",
			"timestamp_values",
			"duration_value",
			"duration_values",
			"field_mask_value",
			"field_mask_values",
			"value_value",
			"value_values",
			"list_value_value",
			"list_value_values",
			"struct_value",
			"struct_values",
			"any_value",
			"any_values",
		},
	},
}

func TestMarshalMessageWithWKTs(t *testing.T) {
	for _, tt := range testMessagesWithWKTs {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithWKTs(t *testing.T) {
	for _, tt := range testMessagesWithWKTs {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}

var testMessagesWithOneofWKTs = []struct {
	name         string
	msg          MessageWithOneofWKTs
	expected     string
	expectedMask []string
}{
	{
		name:     "empty",
		msg:      MessageWithOneofWKTs{},
		expected: `{}`,
	},
	{
		name: "double_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_DoubleValue{},
		},
		expected:     `{"double_value": null}`,
		expectedMask: []string{"double_value"},
	},
	{
		name: "double_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_DoubleValue{DoubleValue: &wrapperspb.DoubleValue{Value: 12.34}},
		},
		expected:     `{"double_value": 12.34}`,
		expectedMask: []string{"double_value"},
	},
	{
		name: "float_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_FloatValue{},
		},
		expected:     `{"float_value": null}`,
		expectedMask: []string{"float_value"},
	},
	{
		name: "float_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_FloatValue{FloatValue: &wrapperspb.FloatValue{Value: 12.34}},
		},
		expected:     `{"float_value": 12.34}`,
		expectedMask: []string{"float_value"},
	},
	{
		name: "int32_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_Int32Value{},
		},
		expected:     `{"int32_value": null}`,
		expectedMask: []string{"int32_value"},
	},
	{
		name: "int32_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_Int32Value{Int32Value: &wrapperspb.Int32Value{Value: -42}},
		},
		expected:     `{"int32_value": -42}`,
		expectedMask: []string{"int32_value"},
	},
	{
		name: "int64_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_Int64Value{},
		},
		expected:     `{"int64_value": null}`,
		expectedMask: []string{"int64_value"},
	},
	{
		name: "int64_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_Int64Value{Int64Value: &wrapperspb.Int64Value{Value: -42}},
		},
		expected:     `{"int64_value": "-42"}`,
		expectedMask: []string{"int64_value"},
	},
	{
		name: "uint32_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_Uint32Value{},
		},
		expected:     `{"uint32_value": null}`,
		expectedMask: []string{"uint32_value"},
	},
	{
		name: "uint32_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_Uint32Value{Uint32Value: &wrapperspb.UInt32Value{Value: 42}},
		},
		expected:     `{"uint32_value": 42}`,
		expectedMask: []string{"uint32_value"},
	},
	{
		name: "uint64_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_Uint64Value{},
		},
		expected:     `{"uint64_value": null}`,
		expectedMask: []string{"uint64_value"},
	},
	{
		name: "uint64_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_Uint64Value{Uint64Value: &wrapperspb.UInt64Value{Value: 42}},
		},
		expected:     `{"uint64_value": "42"}`,
		expectedMask: []string{"uint64_value"},
	},
	{
		name: "bool_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_BoolValue{},
		},
		expected:     `{"bool_value": null}`,
		expectedMask: []string{"bool_value"},
	},
	{
		name: "bool_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_BoolValue{BoolValue: &wrapperspb.BoolValue{Value: true}},
		},
		expected:     `{"bool_value": true}`,
		expectedMask: []string{"bool_value"},
	},
	{
		name: "string_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_StringValue{},
		},
		expected:     `{"string_value": null}`,
		expectedMask: []string{"string_value"},
	},
	{
		name: "string_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_StringValue{StringValue: &wrapperspb.StringValue{Value: "foo"}},
		},
		expected:     `{"string_value": "foo"}`,
		expectedMask: []string{"string_value"},
	},
	{
		name: "bytes_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_BytesValue{},
		},
		expected:     `{"bytes_value": null}`,
		expectedMask: []string{"bytes_value"},
	},
	{
		name: "bytes_zero",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_BytesValue{BytesValue: &wrapperspb.BytesValue{Value: []byte{}}},
		},
		expected:     `{"bytes_value": ""}`,
		expectedMask: []string{"bytes_value"},
	},
	{
		name: "bytes_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_BytesValue{BytesValue: &wrapperspb.BytesValue{Value: []byte("foo")}},
		},
		expected:     `{"bytes_value": "Zm9v"}`,
		expectedMask: []string{"bytes_value"},
	},
	{
		name: "empty_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_EmptyValue{},
		},
		expected:     `{"empty_value": null}`,
		expectedMask: []string{"empty_value"},
	},
	{
		name: "empty_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_EmptyValue{EmptyValue: &emptypb.Empty{}},
		},
		expected:     `{"empty_value": {}}`,
		expectedMask: []string{"empty_value"},
	},
	{
		name: "timestamp_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_TimestampValue{},
		},
		expected:     `{"timestamp_value": null}`,
		expectedMask: []string{"timestamp_value"},
	},
	{
		name: "timestamp_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_TimestampValue{TimestampValue: mustTimestamp(testTime)},
		},
		expected:     `{"timestamp_value": "2006-01-02T08:04:05.123456789Z"}`,
		expectedMask: []string{"timestamp_value"},
	},
	{
		name: "duration_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_DurationValue{},
		},
		expected:     `{"duration_value": null}`,
		expectedMask: []string{"duration_value"},
	},
	{
		name: "duration_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_DurationValue{DurationValue: mustDuration(testDuration)},
		},
		expected:     `{"duration_value": "3723.123456789s"}`,
		expectedMask: []string{"duration_value"},
	},
	{
		name: "field_mask_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_FieldMaskValue{},
		},
		expected:     `{"field_mask_value": null}`,
		expectedMask: []string{"field_mask_value"},
	},
	{
		name: "field_mask_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_FieldMaskValue{FieldMaskValue: &fieldmaskpb.FieldMask{Paths: []string{"foo.bar", "bar", "baz.qux"}}},
		},
		expected:     `{"field_mask_value": "foo.bar,bar,baz.qux"}`,
		expectedMask: []string{"field_mask_value"},
	},
	{
		name: "value_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_ValueValue{ValueValue: &structpb.Value{Kind: &structpb.Value_NullValue{}}},
		},
		expected:     `{"value_value": null}`,
		expectedMask: []string{"value_value"},
	},
	{
		name: "value_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_ValueValue{ValueValue: &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "foo"}}},
		},
		expected:     `{"value_value": "foo"}`,
		expectedMask: []string{"value_value"},
	},
	{
		name: "list_value_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_ListValueValue{},
		},
		expected:     `{"list_value_value": null}`,
		expectedMask: []string{"list_value_value"},
	},
	{
		name: "list_value_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_ListValueValue{
				ListValueValue: &structpb.ListValue{
					Values: []*structpb.Value{
						{Kind: &structpb.Value_NullValue{}},
						{Kind: &structpb.Value_NumberValue{NumberValue: 12.34}},
						{Kind: &structpb.Value_StringValue{StringValue: "foo"}},
						{Kind: &structpb.Value_BoolValue{BoolValue: true}},
					},
				},
			},
		},
		expected:     `{"list_value_value": [null, 12.34, "foo", true]}`,
		expectedMask: []string{"list_value_value"},
	},
	{
		name: "struct_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_StructValue{},
		},
		expected:     `{"struct_value": null}`,
		expectedMask: []string{"struct_value"},
	},
	{
		name: "struct_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"null":   {Kind: &structpb.Value_NullValue{}},
						"number": {Kind: &structpb.Value_NumberValue{NumberValue: 12.34}},
						"string": {Kind: &structpb.Value_StringValue{StringValue: "foo"}},
						"bool":   {Kind: &structpb.Value_BoolValue{BoolValue: true}},
					},
				},
			},
		},
		expected:     `{"struct_value": {"bool": true, "null": null, "number": 12.34, "string": "foo"}}`,
		expectedMask: []string{"struct_value"},
	},
	{
		name: "any_null",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_AnyValue{},
		},
		expected:     `{"any_value": null}`,
		expectedMask: []string{"any_value"},
	},
	{
		name: "any_value",
		msg: MessageWithOneofWKTs{
			Value: &MessageWithOneofWKTs_AnyValue{
				AnyValue: mustAny(&MessageWithMarshaler{Message: "hello"}),
			},
		},
		expected: `{"any_value": {
			"@type": "type.googleapis.com/thethings.json.test.MessageWithMarshaler",
			"message": "hello"
		}}`,
		expectedMask: []string{"any_value"},
	},
}

func TestMarshalMessageWithOneofWKTs(t *testing.T) {
	for _, tt := range testMessagesWithOneofWKTs {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithOneofWKTs(t *testing.T) {
	for _, tt := range testMessagesWithOneofWKTs {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}

var testMessagesWithWKTMaps = []struct {
	name         string
	msg          MessageWithWKTMaps
	expected     string
	expectedMask []string
}{
	{
		name:     "empty",
		msg:      MessageWithWKTMaps{},
		expected: `{}`,
	},
	{
		name: "full",
		msg: MessageWithWKTMaps{
			StringDoubleMap:    map[string]*wrapperspb.DoubleValue{"value": {Value: 12.34}},
			StringFloatMap:     map[string]*wrapperspb.FloatValue{"value": {Value: 12.34}},
			StringInt32Map:     map[string]*wrapperspb.Int32Value{"value": {Value: -42}},
			StringInt64Map:     map[string]*wrapperspb.Int64Value{"value": {Value: -42}},
			StringUint32Map:    map[string]*wrapperspb.UInt32Value{"value": {Value: 42}},
			StringUint64Map:    map[string]*wrapperspb.UInt64Value{"value": {Value: 42}},
			StringBoolMap:      map[string]*wrapperspb.BoolValue{"yes": {Value: true}},
			StringStringMap:    map[string]*wrapperspb.StringValue{"value": {Value: "foo"}},
			StringBytesMap:     map[string]*wrapperspb.BytesValue{"value": {Value: []byte("foo")}},
			StringEmptyMap:     map[string]*emptypb.Empty{"value": {}},
			StringTimestampMap: map[string]*timestamppb.Timestamp{"value": mustTimestamp(testTime)},
			StringDurationMap:  map[string]*durationpb.Duration{"value": mustDuration(testDuration)},
			StringFieldMaskMap: map[string]*fieldmaskpb.FieldMask{"value": {Paths: []string{"foo.bar", "bar", "baz.qux"}}},
			StringValueMap:     map[string]*structpb.Value{"value": {Kind: &structpb.Value_StringValue{StringValue: "foo"}}},
			StringListValueMap: map[string]*structpb.ListValue{"value": {Values: []*structpb.Value{
				{Kind: &structpb.Value_NullValue{}},
				{Kind: &structpb.Value_NumberValue{NumberValue: 12.34}},
				{Kind: &structpb.Value_StringValue{StringValue: "foo"}},
				{Kind: &structpb.Value_BoolValue{BoolValue: true}},
			}}},
			StringStructMap: map[string]*structpb.Struct{
				"value": {Fields: map[string]*structpb.Value{
					"null":   {Kind: &structpb.Value_NullValue{}},
					"number": {Kind: &structpb.Value_NumberValue{NumberValue: 12.34}},
					"string": {Kind: &structpb.Value_StringValue{StringValue: "foo"}},
					"bool":   {Kind: &structpb.Value_BoolValue{BoolValue: true}},
				}},
			},
			StringAnyMap: map[string]*anypb.Any{
				"value": mustAny(&MessageWithMarshaler{Message: "hello"}),
			},
		},
		expected: `{
			"string_double_map": {"value": 12.34},
			"string_float_map": {"value": 12.34},
			"string_int32_map": {"value": -42},
			"string_int64_map": {"value": "-42"},
			"string_uint32_map": {"value": 42},
			"string_uint64_map": {"value": "42"},
			"string_bool_map": {"yes": true},
			"string_string_map": {"value": "foo"},
			"string_bytes_map": {"value": "Zm9v"},
			"string_empty_map": {"value": {}},
			"string_timestamp_map": {"value": "2006-01-02T08:04:05.123456789Z"},
			"string_duration_map": {"value": "3723.123456789s"},
			"string_field_mask_map": {"value": "foo.bar,bar,baz.qux"},
			"string_value_map": {"value": "foo"},
			"string_list_value_map": {"value": [null, 12.34, "foo", true]},
			"string_struct_map": {"value": {"bool": true, "null": null, "number": 12.34, "string": "foo"}},
			"string_any_map": {"value": {"@type": "type.googleapis.com/thethings.json.test.MessageWithMarshaler", "message": "hello"}}
		}`,
		expectedMask: []string{
			"string_double_map",
			"string_float_map",
			"string_int32_map",
			"string_int64_map",
			"string_uint32_map",
			"string_uint64_map",
			"string_bool_map",
			"string_string_map",
			"string_bytes_map",
			"string_empty_map",
			"string_timestamp_map",
			"string_duration_map",
			"string_field_mask_map",
			"string_value_map",
			"string_list_value_map",
			"string_struct_map",
			"string_any_map",
		},
	},
}

func TestMarshalMessageWithWKTMaps(t *testing.T) {
	for _, tt := range testMessagesWithWKTMaps {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithWKTMaps(t *testing.T) {
	for _, tt := range testMessagesWithWKTMaps {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}
