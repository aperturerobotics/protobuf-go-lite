package jsononeof

import (
	"encoding/json"
	"testing"
)

// TestInterleavedOneofJSONRoundTrip verifies ordinary fields survive every choice.
func TestInterleavedOneofJSONRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		choice   isInterleaved_Choice
		expected string
	}{
		{name: "nil", expected: `{"before":"before","betweenMessage":{"value":"between"},"betweenScalar":3,"after":"after","afterMessage":{"value":"after-message"},"optionalZero":0}`},
		{name: "text", choice: &Interleaved_Text{Text: "selected"}, expected: `{"before":"before","text":"selected","betweenMessage":{"value":"between"},"betweenScalar":3,"after":"after","afterMessage":{"value":"after-message"},"optionalZero":0}`},
		{name: "message", choice: &Interleaved_ChildValue{ChildValue: &Child{Value: "selected"}}, expected: `{"before":"before","betweenMessage":{"value":"between"},"betweenScalar":3,"childValue":{"value":"selected"},"after":"after","afterMessage":{"value":"after-message"},"optionalZero":0}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zero := int32(0)
			want := &Interleaved{
				Before:         "before",
				Choice:         tt.choice,
				BetweenMessage: &Child{Value: "between"},
				BetweenScalar:  3,
				After:          "after",
				AfterMessage:   &Child{Value: "after-message"},
				OptionalZero:   &zero,
			}

			encoded, err := json.Marshal(want)
			if err != nil {
				t.Fatal(err)
			}

			if string(encoded) != tt.expected {
				t.Fatalf("JSON = %s, want %s", encoded, tt.expected)
			}

			got := &Interleaved{}
			if err := json.Unmarshal(encoded, got); err != nil {
				t.Fatal(err)
			}
			if !got.EqualVT(want) {
				t.Fatalf("round trip mismatch: got %#v, want %#v", got, want)
			}
		})
	}
}
