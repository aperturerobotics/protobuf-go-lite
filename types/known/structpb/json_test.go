package structpb_test

import (
	"math"
	"strings"
	"testing"

	"github.com/aperturerobotics/protobuf-go-lite/types/known/structpb"
)

// TestValueJSONRoundTrip preserves the kind of every JSON value through reuse.
func TestValueJSONRoundTrip(t *testing.T) {
	inputs := []string{
		`{"":{"nested":[null,true,false,1.25,"1.25",{},[]]},"quoted":"a\"b"}`,
		`["NaN","Infinity","-Infinity",-3.5,0]`,
		`"42"`,
		`42`,
		`false`,
		`null`,
	}
	value := new(structpb.Value)
	for _, input := range inputs {
		if err := value.UnmarshalJSON([]byte(input)); err != nil {
			t.Fatalf("decode %s: %v", input, err)
		}

		encoded, err := value.MarshalJSON()
		if err != nil {
			t.Fatalf("encode %s: %v", input, err)
		}
		if string(encoded) != input {
			t.Errorf("round trip: got %s, want %s", encoded, input)
		}
	}
}

// TestStructuredJSONReuse replaces existing object fields and array elements.
func TestStructuredJSONReuse(t *testing.T) {
	object := new(structpb.Struct)
	for _, input := range []string{`{"old":1}`, `{"new":2}`, `{}`, `null`} {
		if err := object.UnmarshalJSON([]byte(input)); err != nil {
			t.Fatal(err)
		}

		encoded, err := object.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		want := input
		if input == "null" {
			want = "{}"
		}
		if string(encoded) != want {
			t.Errorf("object: got %s, want %s", encoded, want)
		}
	}

	list := new(structpb.ListValue)
	for _, input := range []string{`[1,2]`, `[3]`, `[]`, `null`} {
		if err := list.UnmarshalJSON([]byte(input)); err != nil {
			t.Fatal(err)
		}

		encoded, err := list.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		want := input
		if input == "null" {
			want = "[]"
		}
		if string(encoded) != want {
			t.Errorf("list: got %s, want %s", encoded, want)
		}
	}
}

// TestValueJSONErrors rejects absent kinds and numbers that cannot retain their kind.
func TestValueJSONErrors(t *testing.T) {
	for _, value := range []*structpb.Value{
		{},
		structpb.NewNumberValue(math.NaN()),
		structpb.NewNumberValue(math.Inf(1)),
		structpb.NewNumberValue(math.Inf(-1)),
	} {
		if _, err := value.MarshalJSON(); err == nil {
			t.Errorf("encoded invalid value: %#v", value)
		}
	}

	for _, input := range []string{``, `{"x":1,"x":2}`, `[1,]`, `{"x":}`, `1e1000`} {
		value := new(structpb.Value)
		if err := value.UnmarshalJSON([]byte(input)); err == nil {
			t.Errorf("decoded invalid JSON: %s", input)
		}
	}

	value := &structpb.Struct{Fields: map[string]*structpb.Value{
		"nested": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{{}}}),
	}}
	if _, err := value.MarshalJSON(); err == nil || !strings.Contains(err.Error(), "nested.0") {
		t.Errorf("missing nested error path: %v", err)
	}
}

// TestNullValueJSON uses null independently of enum marshaling settings.
func TestNullValueJSON(t *testing.T) {
	value := structpb.NullValue_NULL_VALUE
	if err := value.UnmarshalJSON([]byte("null")); err != nil {
		t.Fatal(err)
	}

	encoded, err := value.MarshalJSON()
	if err != nil || string(encoded) != "null" {
		t.Fatalf("got %s, %v", encoded, err)
	}
	if err := value.UnmarshalJSON([]byte("false")); err == nil {
		t.Fatal("accepted a boolean as null")
	}
}
