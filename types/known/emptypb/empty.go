package emptypb

import "github.com/aperturerobotics/protobuf-go-lite/json"

// MarshalJSON marshals the Empty to JSON.
func (x *Empty) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the Empty from JSON.
func (x *Empty) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals an Empty from JSON.
func (x *Empty) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	s.ReadObject(func(key string) {
		s.SetErrorf("unexpected key %q in Empty", key)
	})
	if s.Err() != nil {
		return
	}
	*x = Empty{}
}

// MarshalProtoJSON marshals an Empty to JSON.
func (x *Empty) MarshalProtoJSON(s *json.MarshalState) {
	s.WriteObjectStart()
	s.WriteObjectEnd()
}
