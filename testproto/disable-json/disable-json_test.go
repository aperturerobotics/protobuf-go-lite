package disable_json

import (
	sjson "encoding/json"
	"testing"

	"github.com/aperturerobotics/protobuf-go-lite/json"
)

// TestDisableJson asserts that MessageDisableJson does not implement json.Unmarshal or json.Marshal.
func TestDisableJson(t *testing.T) {
	m := &MessageDisableJson{Body: &MessageDisableJson_Hello{Hello: true}}

	var i any = m
	switch i.(type) {
	case json.Marshaler:
		t.FailNow()
	case json.Unmarshaler:
		t.FailNow()
	case sjson.Marshaler:
		t.FailNow()
	case sjson.Unmarshaler:
		t.FailNow()
	}
}
