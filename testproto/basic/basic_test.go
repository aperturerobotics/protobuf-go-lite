package basic

import (
	"strings"
	"testing"
)

// NewMockBasicMsg buidls a new mock BasicMsg.
func NewMockBasicMsg() *BasicMsg {
	return &BasicMsg{
		Int32Field:         123,
		Int64Field:         45678901234,
		Uint32Field:        345,
		Uint64Field:        9876543210,
		Sint32Field:        -123,
		Sint64Field:        -4567890,
		Fixed32Field:       999,
		Fixed64Field:       88888888,
		Sfixed32Field:      -777,
		Sfixed64Field:      -6666666,
		FloatField:         1.23,
		DoubleField:        4.56,
		BoolField:          true,
		StringField:        "test string",
		BytesField:         []byte{0x01, 0x02, 0x03},
		RepeatedInt32Field: []int32{10, 20, 30},
		MapStringInt32Field: map[string]int32{
			"foo": 100,
			"bar": 200,
		},
		EnumField: BasicMsg_SECOND,
		NestedMessage: &BasicMsg_NestedMsg{
			NestedInt32:  321,
			NestedString: "nested test",
		},
		MyOneof: &BasicMsg_OneofInt32{
			OneofInt32: 777,
		},
	}
}

func TestBasicMarshalUnmarshal(t *testing.T) {
	basic := NewMockBasicMsg()

	data, err := basic.MarshalVT()
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(data) == 0 {
		t.Fatal("unexpected empty data")
	}

	out := &BasicMsg{}
	if err := out.UnmarshalVT(data); err != nil {
		t.Fatal(err.Error())
	}

	if !out.EqualVT(basic) {
		t.Fatal("message not equal after unmarshal")
	}
}

func TestBasicProtoTextMarshal(t *testing.T) {
	basic := NewMockBasicMsg()
	textData := basic.String()
	t.Logf("ProtoText format:\n%s", textData)
}

func TestBasicProtoTextFields(t *testing.T) {
	// Create a BasicMsg with specific fields
	basic := &BasicMsg{
		Int32Field:  123,
		StringField: "hello world",
		// BoolField:   false, // Omit this since false is default and won't appear
		EnumField:   BasicMsg_FIRST,
		FloatField:  3.14,
		DoubleField: 6.28,
		MyOneof: &BasicMsg_OneofString{
			OneofString: "oneof test",
		},
	}

	// Marshal to ProtoText format
	textData := basic.String()
	t.Logf("ProtoText format:\n%s", textData)

	// Verify that specific fields are present in the text output
	expectedFields := []string{
		"int32_field: 123",
		`string_field: "hello world"`,
		`enum_field: "FIRST"`,
		"float_field: 3.14",
		"double_field: 6.28",
		`oneof_string: "oneof test"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(textData, field) {
			t.Errorf("Expected %s in ProtoText output", field)
		}
	}

	// Verify that bool_field is not present in the output
	if strings.Contains(textData, "bool_field") {
		t.Errorf("Did not expect bool_field in ProtoText output")
	}
}

func TestBasicProtoTextOneof(t *testing.T) {
	// Test with oneof_string set
	basicString := &BasicMsg{
		MyOneof: &BasicMsg_OneofString{
			OneofString: "test oneof string",
		},
	}

	textDataString := basicString.String()
	t.Logf("ProtoText with oneof_string:\n%s", textDataString)

	// Test with oneof_int32 set
	basicInt32 := &BasicMsg{
		MyOneof: &BasicMsg_OneofInt32{
			OneofInt32: 42,
		},
	}

	textDataInt32 := basicInt32.String()
	t.Logf("ProtoText with oneof_int32:\n%s", textDataInt32)
}

func TestBasicProtoTextRepeatedAndMap(t *testing.T) {
	// Create a BasicMsg with repeated and map fields
	basic := &BasicMsg{
		RepeatedInt32Field: []int32{1, 2, 3},
		MapStringInt32Field: map[string]int32{
			"alpha": 100,
			"beta":  200,
		},
	}

	// Convert to string format
	textData := basic.String()
	t.Logf("String representation with repeated and map fields:\n%s", textData)

	// Expected string representation
	expected := `BasicMsg {repeated_int32_field: [1, 2, 3] map_string_int32_field: { "alpha": 100 "beta": 200 }}`

	// Trim spaces and newlines for comparison
	textData = strings.TrimSpace(textData)
	expected = strings.TrimSpace(expected)

	if textData != expected {
		t.Errorf("String representation mismatch.\nExpected: %s\nGot: %s", expected, textData)
	}
}

func TestBasicStringDefaultValues(t *testing.T) {
	// Create a BasicMsg with default values
	basic := &BasicMsg{}

	// Convert to string format
	textData := strings.TrimSpace(basic.String())
	t.Logf("String representation with default values:\n%s", textData)

	// Expected output
	expected := "BasicMsg {}"

	// Verify that the output matches the expected output
	if textData != expected {
		t.Errorf("Expected output for default values to be:\n%s\nGot:\n%s", expected, textData)
	}
}

func TestBasicStringNestedMessage(t *testing.T) {
	// Create a BasicMsg with a nested message
	basic := &BasicMsg{
		NestedMessage: &BasicMsg_NestedMsg{
			NestedInt32:  12345,
			NestedString: "nested message test",
		},
	}

	// Convert to string format
	textData := basic.String()
	t.Logf("String representation with nested message:\n%s", textData)

	// Verify that nested fields are present in the output
	if !strings.Contains(textData, "nested_int32: 12345") {
		t.Errorf("Expected nested_int32: 12345 in string output")
	}
	if !strings.Contains(textData, `nested_string: "nested message test"`) {
		t.Errorf(`Expected nested_string: "nested message test" in string output`)
	}
}

func TestBasicStringEnumField(t *testing.T) {
	// Create a BasicMsg with enum fields
	basic := &BasicMsg{
		EnumField: BasicMsg_SECOND,
	}

	// Convert to string format
	textData := basic.String()
	t.Logf("String representation with enum field:\n%s", textData)

	// Verify that enum field is present in the output
	if !strings.Contains(textData, `enum_field: "SECOND"`) {
		t.Errorf("Expected enum_field: \"SECOND\" in string output")
	}
}
