package test_test

import (
	"testing"

	. "github.com/TheThingsIndustries/protoc-gen-go-json/test/gogo"
	types "github.com/TheThingsIndustries/protoc-gen-go-json/test/types"
)

var testMessagesWithGoGoOptions = []struct {
	name         string
	msg          MessageWithGoGoOptions
	expected     string
	expectedMask []string
}{
	{
		name: "empty",
		msg: MessageWithGoGoOptions{
			NonNullableEUIWithCustomNameAndType: types.EUI64{1, 2, 3, 4, 5, 6, 7, 8},
			NonNullableDuration:                 testDuration,
			NonNullableTimestamp:                testTime,
		},
		expected: `{
			"non_nullable_eui_with_custom_name_and_type": "0102030405060708",
			"non_nullable_duration": "3723.123456789s",
			"non_nullable_timestamp": "2006-01-02T08:04:05.123456789Z"
		}`,
		expectedMask: []string{
			"non_nullable_eui_with_custom_name_and_type",
			"non_nullable_duration",
			"non_nullable_timestamp",
		},
	},
	{
		name: "zero",
		msg: MessageWithGoGoOptions{
			NonNullableEUIWithCustomNameAndType: types.EUI64{1, 2, 3, 4, 5, 6, 7, 8},
			NonNullableDuration:                 testDuration,
			NonNullableTimestamp:                testTime,
		},
		expected: `{
			"eui_with_custom_name": null,
			"eui_with_custom_name_and_type": null,
			"non_nullable_eui_with_custom_name_and_type": "0102030405060708",
			"euis_with_custom_name_and_type": [],
			"duration": null,
			"non_nullable_duration": "3723.123456789s",
			"timestamp": null,
			"non_nullable_timestamp": "2006-01-02T08:04:05.123456789Z"
		}`,
		expectedMask: []string{
			"eui_with_custom_name",
			"eui_with_custom_name_and_type",
			"non_nullable_eui_with_custom_name_and_type",
			"euis_with_custom_name_and_type",
			"duration",
			"non_nullable_duration",
			"timestamp",
			"non_nullable_timestamp",
		},
	},
	{
		name: "full",
		msg: MessageWithGoGoOptions{
			EUIWithCustomName:                   []byte{1, 2, 3, 4, 5, 6, 7, 8},
			EUIWithCustomNameAndType:            &types.EUI64{1, 2, 3, 4, 5, 6, 7, 8},
			NonNullableEUIWithCustomNameAndType: types.EUI64{1, 2, 3, 4, 5, 6, 7, 8},
			EUIsWithCustomNameAndType: []types.EUI64{
				{1, 2, 3, 4, 5, 6, 7, 8},
			},
			Duration:             &testDuration,
			NonNullableDuration:  testDuration,
			Timestamp:            &testTime,
			NonNullableTimestamp: testTime,
		},
		expected: `{
			"eui_with_custom_name": "AQIDBAUGBwg=",
			"eui_with_custom_name_and_type": "0102030405060708",
			"non_nullable_eui_with_custom_name_and_type": "0102030405060708",
			"euis_with_custom_name_and_type": ["0102030405060708"],
			"duration": "3723.123456789s",
			"non_nullable_duration": "3723.123456789s",
			"timestamp": "2006-01-02T08:04:05.123456789Z",
			"non_nullable_timestamp": "2006-01-02T08:04:05.123456789Z"
		}`,
		expectedMask: []string{
			"eui_with_custom_name",
			"eui_with_custom_name_and_type",
			"non_nullable_eui_with_custom_name_and_type",
			"euis_with_custom_name_and_type",
			"duration",
			"non_nullable_duration",
			"timestamp",
			"non_nullable_timestamp",
		},
	},
}

func TestMarshalMessageWithGoGoOptions(t *testing.T) {
	for _, tt := range testMessagesWithGoGoOptions {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithGoGoOptions(t *testing.T) {
	for _, tt := range testMessagesWithGoGoOptions {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}

var testMessagesWithNullable = []struct {
	name         string
	msg          MessageWithNullable
	expected     string
	expectedMask []string
}{
	{
		name: "empty",
		msg: MessageWithNullable{
			Sub:      SubMessage{Field: "foo"},
			OtherSub: SubMessageWithoutMarshalers{OtherField: "foo"},
		},
		expected: `{
			"sub": {"field": "foo"},
			"other_sub": {"other_field": "foo"}
		}`,
		expectedMask: []string{
			"sub.field",
			"other_sub", // NOTE: no marshaler.
		},
	},
	{
		name: "full",
		msg: MessageWithNullable{
			Sub: SubMessage{Field: "foo"},
			Subs: []SubMessage{
				{Field: "foo"},
				{Field: "bar"},
			},
			OtherSub: SubMessageWithoutMarshalers{OtherField: "foo"},
			OtherSubs: []SubMessageWithoutMarshalers{
				{OtherField: "foo"},
				{OtherField: "bar"},
			},
		},
		expected: `{
			"sub": {"field": "foo"},
			"subs": [
				{"field": "foo"},
				{"field": "bar"}
			],
			"other_sub": {"other_field": "foo"},
			"other_subs": [
				{"other_field": "foo"},
				{"other_field": "bar"}
			]
		}`,
		expectedMask: []string{
			"sub.field",
			"subs",
			"other_sub", // NOTE: no marshaler.
			"other_subs",
		},
	},
}

func TestMarshalMessageWithNullable(t *testing.T) {
	for _, tt := range testMessagesWithNullable {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithNullable(t *testing.T) {
	for _, tt := range testMessagesWithNullable {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}

var testMessagesWithEmbedded = []struct {
	name         string
	msg          MessageWithEmbedded
	expected     string
	expectedMask []string
}{
	{
		name:     "empty",
		msg:      MessageWithEmbedded{},
		expected: `{}`,
	},
	{
		name: "zero",
		msg: MessageWithEmbedded{
			SubMessage: &SubMessage{},
		},
		expected: `{
			"sub": {"field": ""}
		}`,
		expectedMask: []string{
			"sub.field",
		},
	},
	{
		name: "full",
		msg: MessageWithEmbedded{
			SubMessage:                  &SubMessage{Field: "foo"},
			SubMessageWithoutMarshalers: &SubMessageWithoutMarshalers{OtherField: "foo"},
		},
		expected: `{
			"sub": {"field": "foo"},
			"other_sub": {"other_field": "foo"}
		}`,
		expectedMask: []string{
			"sub.field",
			"other_sub", // NOTE: no marshaler.
		},
	},
}

func TestMarshalMessageWithEmbedded(t *testing.T) {
	for _, tt := range testMessagesWithEmbedded {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithEmbedded(t *testing.T) {
	for _, tt := range testMessagesWithEmbedded {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}

var testMessagesWithNullableEmbedded = []struct {
	name         string
	msg          MessageWithNullableEmbedded
	expected     string
	expectedMask []string
}{
	{
		name: "full",
		msg: MessageWithNullableEmbedded{
			SubMessage:                  SubMessage{Field: "foo"},
			SubMessageWithoutMarshalers: SubMessageWithoutMarshalers{OtherField: "foo"},
		},
		expected: `{
			"sub": {"field": "foo"},
			"other_sub": {"other_field": "foo"}
		}`,
		expectedMask: []string{
			"sub.field",
			"other_sub", // NOTE: no marshaler.
		},
	},
}

func TestMarshalMessageWithNullableEmbedded(t *testing.T) {
	for _, tt := range testMessagesWithNullableEmbedded {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithNullableEmbedded(t *testing.T) {
	for _, tt := range testMessagesWithNullableEmbedded {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}
