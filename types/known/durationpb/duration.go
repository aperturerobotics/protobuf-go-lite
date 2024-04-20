package durationpb

import (
	"strconv"
	"strings"

	"github.com/aperturerobotics/protobuf-go-lite/json"
)

// MarshalJSON marshals the Duration to JSON.
func (x *Duration) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the Duration from JSON.
func (x *Duration) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// String formats the duration to a string.
func (d *Duration) String() string {
	var out strings.Builder
	secs, nanos := d.GetSeconds(), d.GetNanos()
	if secs != 0 {
		_, _ = out.WriteString("seconds:")
		_, _ = out.WriteString(strconv.FormatInt(secs, 10))
	}
	if nanos != 0 {
		if out.Len() != 0 {
			_, _ = out.WriteString(" ")
		}
		_, _ = out.WriteString("nanos:")
		_, _ = out.WriteString(strconv.Itoa(int(nanos)))
	}
	return out.String()
}

// UnmarshalProtoJSON unmarshals a Duration from JSON.
func (x *Duration) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	d := s.ReadDuration()
	if s.Err() != nil {
		return
	}
	*x = *New(*d)
}

// MarshalProtoJSON marshals a Duration to JSON.
func (x *Duration) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteDuration(x.AsDuration())
}
