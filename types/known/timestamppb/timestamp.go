package timestamppb

import (
	"errors"
	"strconv"
	"strings"
	time "time"

	"github.com/aperturerobotics/protobuf-go-lite/json"
)

// ErrEmptyTimestamp is returned from Validate if the timestamp was empty.
var ErrEmptyTimestamp = errors.New("empty timestamp")

// ToTimestamp constructs a new Timestamp from the provided time.Time.
func ToTimestamp(t time.Time) *Timestamp {
	return New(t)
}

// FromUnixMilli constructs a new Timestamp from the provided unix milliseconds value.
//
// Example: time.Now().UnixMilli()
func FromUnixMilli(timestampUnixMilli int64) *Timestamp {
	return New(time.UnixMilli(timestampUnixMilli))
}

// GetEmpty checks if the timestamp is empty.
func (x *Timestamp) GetEmpty() bool {
	return x.SizeVT() == 0
}

// Validate is an alias to CheckValid.
func (x *Timestamp) Validate(allowEmpty bool) error {
	isEmpty := x.GetEmpty()
	if isEmpty {
		if allowEmpty {
			return nil
		}
		return ErrEmptyTimestamp
	}

	return x.CheckValid()
}

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

// MarshalProtoText formats the Timestamp to a proto text string.
func (t *Timestamp) MarshalProtoText() string {
	return t.String()
}
