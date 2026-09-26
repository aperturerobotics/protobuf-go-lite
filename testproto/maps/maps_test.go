package testproto_maps

import (
	"strconv"
	"testing"

	"github.com/aperturerobotics/protobuf-go-lite/types/known/timestamppb"
)

// TestMapJSONKeyOrder verifies map fields encode in key order, so equal
// messages produce equal JSON regardless of Go map iteration order.
func TestMapJSONKeyOrder(t *testing.T) {
	msg := &MsgWithMaps{
		StringKeys: make(map[string]*timestamppb.Timestamp),
		IntKeys:    make(map[uint32]*timestamppb.Timestamp),
		BoolKeys:   map[bool]string{true: "yes", false: "no"},
	}
	for i := range 16 {
		msg.StringKeys["key-"+strconv.Itoa(i)] = &timestamppb.Timestamp{}
		msg.IntKeys[uint32(16-i)] = &timestamppb.Timestamp{}
	}

	want, err := msg.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	for range 8 {
		got, err := msg.CloneVT().MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Fatalf("JSON changed between marshals:\n%s\n%s", got, want)
		}
	}

	small := &MsgWithMaps{
		IntKeys:  map[uint32]*timestamppb.Timestamp{10: {}, 2: {}},
		BoolKeys: map[bool]string{true: "yes", false: "no"},
	}
	encoded, err := small.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	const expected = `{"intKeys":{"2":"1970-01-01T00:00:00Z","10":"1970-01-01T00:00:00Z"},"boolKeys":{"false":"no","true":"yes"}}`
	if string(encoded) != expected {
		t.Fatalf("JSON = %s, want %s", encoded, expected)
	}
}
