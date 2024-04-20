package timestamppb

import (
	"strconv"
	"strings"

	"github.com/aperturerobotics/protobuf-go-lite/json"
)

// MarshalJSON marshals the Timestamp to JSON.
func (x *Timestamp) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the Timestamp from JSON.
func (x *Timestamp) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a Timestamp from JSON.
func (x *Timestamp) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	t := s.ReadTime()
	if s.Err() != nil {
		return
	}
	*x = *New(*t)
}

// MarshalProtoJSON marshals a Timestamp to JSON.
func (x *Timestamp) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteTime(x.AsTime())
}

// String formats the timestamp to a string.
func (t *Timestamp) String() string {
	var out strings.Builder
	secs, nanos := t.GetSeconds(), t.GetNanos()
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
