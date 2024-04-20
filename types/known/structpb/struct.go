package structpb

import (
	"errors"

	"github.com/aperturerobotics/protobuf-go-lite/json"
)

// ErrJSONNotSupported is returned when JSON marshaling or unmarshaling is unsupported
var ErrJSONNotSupported = errors.New("JSON marshal/unmarshal is not supported for Struct")

// MarshalJSON marshals the Struct to JSON.
func (x *Struct) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the Struct from JSON.
func (x *Struct) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a Struct from JSON.
func (x *Struct) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	// The Struct type is not yet supported.
	s.SetError(ErrJSONNotSupported)
}

// MarshalProtoJSON marshals a Struct to JSON.
func (x *Struct) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	// The Struct type is not yet supported.
	s.SetError(ErrJSONNotSupported)
}

// MarshalJSON marshals the Value to JSON.
func (x *Value) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the Value from JSON.
func (x *Value) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// MarshalProtoJSON marshals a Value to JSON.
func (x *Value) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	// The Value type is not yet supported.
	s.SetError(ErrJSONNotSupported)
}

// UnmarshalProtoJSON marshals a Struct to JSON.
func (x *Value) UnmarshalProtoJSON(s *json.UnmarshalState) {
	// The Struct type is not yet supported.
	s.SetError(ErrJSONNotSupported)
}

// MarshalJSON marshals the ListValue to JSON.
func (x *ListValue) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the ListValue from JSON.
func (x *ListValue) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// MarshalProtoJSON marshals a ListValue to JSON.
func (x *ListValue) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	// The ListValue type is not yet supported.
	s.SetError(ErrJSONNotSupported)
}

// UnmarshalProtoJSON marshals a ListValue to JSON.
func (x *ListValue) UnmarshalProtoJSON(s *json.UnmarshalState) {
	// The ListValue type is not supported.
	s.SetError(ErrJSONNotSupported)
}

// MarshalJSON marshals the NullValue to JSON.
func (x *NullValue) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the NullValue from JSON.
func (x *NullValue) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// MarshalProtoJSON marshals a NullValue to JSON.
func (x *NullValue) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	// The NullValue type is not yet supported.
	s.SetError(ErrJSONNotSupported)
}

// UnmarshalProtoJSON marshals a NullValue to JSON.
func (x *NullValue) UnmarshalProtoJSON(s *json.UnmarshalState) {
	// The NullValue type is not supported.
	s.SetError(ErrJSONNotSupported)
}
