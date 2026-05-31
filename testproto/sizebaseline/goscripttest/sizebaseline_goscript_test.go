package goscripttest

import (
	"testing"

	sizebaseline "github.com/aperturerobotics/protobuf-go-lite/testproto/sizebaseline"
)

func TestGoScriptSizeBaseline(t *testing.T) {
	required := int32(7)
	explicit := int32(11)
	signed := int32(-65)
	msg := &sizebaseline.SizeBaseline{
		ExplicitInt32:  &explicit,
		ExplicitSint32: &signed,
		RequiredInt32:  &required,
		PackedInt32:    []int32{1, 2, 300},
		ExpandedInt32:  []int32{4, 5},
		BytesValue:     []byte("bytes"),
	}

	wire, err := msg.MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := msg.SizeVT(), len(wire); got != want {
		t.Fatalf("SizeVT() = %d, want %d", got, want)
	}

	var out sizebaseline.SizeBaseline
	if err := out.UnmarshalVT(wire); err != nil {
		t.Fatal(err)
	}
	if !out.EqualVT(msg) {
		t.Fatal("GoScript size baseline round trip mismatch")
	}
}
