package test_test

import (
	"fmt"
	"testing"

	. "github.com/TheThingsIndustries/protoc-gen-go-json/test/gogo"
)

var testMessagesWithWithEnums = []struct {
	name         string
	msg          MessageWithEnums
	expected     string
	expectedMask []string
}{
	{
		name:     "empty",
		msg:      MessageWithEnums{},
		expected: `{}`,
	},
	{
		name: "zero",
		msg:  MessageWithEnums{},
		expected: `{
			"regular": 0,
			"regulars": [],
			"custom": "CUSTOM_UNKNOWN",
			"customs": [],
			"wrapped_custom": null,
			"wrapped_customs": []
		}`,
		expectedMask: []string{
			"regular",
			"regulars",
			"custom",
			"customs",
			"wrapped_custom",
			"wrapped_customs",
		},
	},
	{
		name: "full",
		msg: MessageWithEnums{
			Regular:  REGULAR_A,
			Regulars: []RegularEnum{REGULAR_A, REGULAR_B},
			Custom:   CustomEnum_CUSTOM_V1_0,
			Customs: []CustomEnum{
				CustomEnum_CUSTOM_V1_0,
				CustomEnum_CUSTOM_V1_0_1,
			},
			WrappedCustom: &CustomEnumValue{
				Value: CustomEnum_CUSTOM_V1_0,
			},
			WrappedCustoms: []*CustomEnumValue{
				{Value: CustomEnum_CUSTOM_V1_0},
				{Value: CustomEnum_CUSTOM_V1_0_1},
			},
		},
		expected: `{
			"regular": 1,
			"regulars": [1, 2],
			"custom": "1.0",
			"customs": ["1.0", "1.0.1"],
			"wrapped_custom": "1.0",
			"wrapped_customs": ["1.0", "1.0.1"]
		}`,
		expectedMask: []string{
			"regular",
			"regulars",
			"custom",
			"customs",
			"wrapped_custom",
			"wrapped_customs",
		},
	},
}

func TestMarshalMessageWithEnums(t *testing.T) {
	for _, tt := range testMessagesWithWithEnums {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithEnums(t *testing.T) {
	for _, tt := range testMessagesWithWithEnums {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}

var testMessagesWithWithOneofEnums = []struct {
	name         string
	msg          MessageWithOneofEnums
	expected     string
	expectedMask []string
}{
	{
		name:     "empty",
		msg:      MessageWithOneofEnums{},
		expected: `{}`,
	},
	{
		name: "regular_zero",
		msg: MessageWithOneofEnums{
			Value: &MessageWithOneofEnums_Regular{Regular: REGULAR_UNKNOWN},
		},
		expected:     `{"regular": 0}`,
		expectedMask: []string{"regular"},
	},
	{
		name: "regular",
		msg: MessageWithOneofEnums{
			Value: &MessageWithOneofEnums_Regular{Regular: REGULAR_A},
		},
		expected:     `{"regular": 1}`,
		expectedMask: []string{"regular"},
	},
	{
		name: "custom_zero",
		msg: MessageWithOneofEnums{
			Value: &MessageWithOneofEnums_Custom{Custom: CustomEnum_CUSTOM_UNKNOWN},
		},
		expected:     `{"custom": "CUSTOM_UNKNOWN"}`,
		expectedMask: []string{"custom"},
	},
	{
		name: "custom",
		msg: MessageWithOneofEnums{
			Value: &MessageWithOneofEnums_Custom{Custom: CustomEnum_CUSTOM_V1_0},
		},
		expected:     `{"custom": "1.0"}`,
		expectedMask: []string{"custom"},
	},
	{
		name: "wrapped_zero",
		msg: MessageWithOneofEnums{
			Value: &MessageWithOneofEnums_WrappedCustom{WrappedCustom: &CustomEnumValue{
				Value: CustomEnum_CUSTOM_UNKNOWN,
			}},
		},
		expected:     `{"wrapped_custom": "CUSTOM_UNKNOWN"}`,
		expectedMask: []string{"wrapped_custom"},
	},
	{
		name: "wrapped",
		msg: MessageWithOneofEnums{
			Value: &MessageWithOneofEnums_WrappedCustom{WrappedCustom: &CustomEnumValue{
				Value: CustomEnum_CUSTOM_V1_0,
			}},
		},
		expected:     `{"wrapped_custom": "1.0"}`,
		expectedMask: []string{"wrapped_custom"},
	},
}

func TestMarshalMessageWithOneofEnums(t *testing.T) {
	for _, tt := range testMessagesWithWithOneofEnums {
		t.Run(tt.name, func(t *testing.T) {
			expectMarshalEqual(t, &tt.msg, tt.expectedMask, []byte(tt.expected))
		})
	}
}

func TestUnmarshalMessageWithOneofEnums(t *testing.T) {
	for _, tt := range testMessagesWithWithOneofEnums {
		t.Run(tt.name, func(t *testing.T) {
			expectUnmarshalEqual(t, &tt.msg, []byte(tt.expected), tt.expectedMask)
		})
	}
}

func TestCustomEnum_TextMarshalers(t *testing.T) {
	for _, tt := range []struct {
		enum   CustomEnum
		values []string
	}{
		{CustomEnum_CUSTOM_UNKNOWN, []string{"CUSTOM_UNKNOWN", "UNKNOWN", "0"}},
		{CustomEnum_CUSTOM_V1_0, []string{"1.0", "1.0.0", "CUSTOM_V1_0", "V1_0", "1"}},
		{CustomEnum_CUSTOM_V1_0_1, []string{"1.0.1", "CUSTOM_V1_0_1", "V1_0_1", "2"}},
	} {
		t.Run(fmt.Sprintf("MarshalText_%s", tt.enum), func(t *testing.T) {
			txt, err := tt.enum.MarshalText()
			if err != nil {
				t.Fatal(err)
			}
			if string(txt) != tt.values[0] {
				t.Errorf("expected: %s, got: %s", tt.values[0], string(txt))
			}
		})
		t.Run(fmt.Sprintf("UnmarshalText_%s", tt.enum), func(t *testing.T) {
			for _, value := range tt.values {
				var enum CustomEnum
				err := enum.UnmarshalText([]byte(value))
				if err != nil {
					t.Fatal(err)
				}
				if enum != tt.enum {
					t.Errorf("expected %q to unmarshal as %s, got: %s", value, tt.enum, enum)
				}
			}
		})
	}
}
