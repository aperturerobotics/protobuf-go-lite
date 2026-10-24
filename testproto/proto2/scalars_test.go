package proto2

import (
	"errors"
	"io"
	"testing"
)

func TestPackedFixed32MalformedLength(t *testing.T) {
	wire := []byte{0x22, 0x03, 0x01, 0x02, 0x03}

	var safeOut Fixed32Message
	if err := safeOut.UnmarshalVT(wire); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("UnmarshalVT malformed packed fixed32 = %v, want ErrUnexpectedEOF", err)
	}

	var unsafeOut Fixed32Message
	if err := unsafeOut.UnmarshalVTUnsafe(wire); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("UnmarshalVTUnsafe malformed packed fixed32 = %v, want ErrUnexpectedEOF", err)
	}
}

// TestJSONPresence round-trips proto2 messages through JSON: a set optional
// field is written even when it holds the zero value, and an unset one is left
// out.
func TestJSONPresence(t *testing.T) {
	// A set zero-valued optional field and a repeated field are written.
	zero := int32(0)
	in := &Int32Message{OptionalField: &zero, RepeatedField: []int32{1, 2}}
	data, err := in.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"optionalField":0,"repeatedField":[1,2]}`; string(data) != want {
		t.Fatalf("MarshalJSON = %s, want %s", data, want)
	}

	// Reading it back restores the same presence.
	var out Int32Message
	if err := out.UnmarshalJSON(data); err != nil {
		t.Fatal(err)
	}
	if !out.EqualVT(in) {
		t.Fatalf("UnmarshalJSON = %v, want %v", &out, in)
	}

	// Closed proto2 enums are written as numbers by default.
	num := EnumMessage_NINE
	enumIn := &EnumMessage{RequiredField: &num, PackedField: []EnumMessage_Num{EnumMessage_TEN}}
	data, err = enumIn.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"requiredField":9,"packedField":[10]}`; string(data) != want {
		t.Fatalf("MarshalJSON = %s, want %s", data, want)
	}

	// Reading the enums back restores them.
	var enumOut EnumMessage
	if err := enumOut.UnmarshalJSON(data); err != nil {
		t.Fatal(err)
	}
	if !enumOut.EqualVT(enumIn) {
		t.Fatalf("UnmarshalJSON = %v, want %v", &enumOut, enumIn)
	}
}
