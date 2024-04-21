package wrapperspb

import (
	"math"
	"strconv"

	"github.com/aperturerobotics/protobuf-go-lite/json"
)

// MarshalJSON marshals the DoubleValue to JSON.
func (x *DoubleValue) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the DoubleValue from JSON.
func (x *DoubleValue) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a DoubleValue from JSON.
func (x *DoubleValue) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	*x = DoubleValue{Value: s.ReadFloat64()}
}

// MarshalProtoJSON marshals a DoubleValue to JSON.
func (x *DoubleValue) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteFloat64(x.Value)
}

// String formats the DoubleValue to a string.
func (x *DoubleValue) String() string {
	return strconv.FormatFloat(x.Value, 'g', -1, 64)
}

// MarshalJSON marshals the FloatValue to JSON.
func (x *FloatValue) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the FloatValue from JSON.
func (x *FloatValue) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a FloatValue from JSON.
func (x *FloatValue) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	*x = FloatValue{Value: float32(s.ReadFloat64())}
}

// MarshalProtoJSON marshals a FloatValue to JSON.
func (x *FloatValue) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteFloat64(float64(x.Value))
}

// String formats the FloatValue to a string.
func (x *FloatValue) String() string {
	return strconv.FormatFloat(float64(x.Value), 'g', -1, 32)
}

// MarshalJSON marshals the Int64Value to JSON.
func (x *Int64Value) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the Int64Value from JSON.
func (x *Int64Value) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a Int64Value from JSON.
func (x *Int64Value) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	*x = Int64Value{Value: s.ReadInt64()}
}

// MarshalProtoJSON marshals a Int64Value to JSON.
func (x *Int64Value) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteInt64(x.Value)
}

// String formats the Int64Value to a string.
func (x *Int64Value) String() string {
	return strconv.FormatInt(x.Value, 10)
}

// MarshalJSON marshals the UInt64Value to JSON.
func (x *UInt64Value) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the UInt64Value from JSON.
func (x *UInt64Value) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a UInt64Value from JSON.
func (x *UInt64Value) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	*x = UInt64Value{Value: s.ReadUint64()}
}

// MarshalProtoJSON marshals a UInt64Value to JSON.
func (x *UInt64Value) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteUint64(x.Value)
}

// String formats the UInt64Value to a string.
func (x *UInt64Value) String() string {
	return strconv.FormatUint(x.Value, 10)
}

// MarshalJSON marshals the Int32Value to JSON.
func (x *Int32Value) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the Int32Value from JSON.
func (x *Int32Value) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a Int32Value from JSON.
func (x *Int32Value) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	v := s.ReadInt64()
	if v < math.MinInt32 || v > math.MaxInt32 {
		s.SetErrorf("value out of range for int32: %v", v)
		return
	}
	*x = Int32Value{Value: int32(v)}
}

// MarshalProtoJSON marshals a Int32Value to JSON.
func (x *Int32Value) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteInt64(int64(x.Value))
}

// String formats the Int32Value to a string.
func (x *Int32Value) String() string {
	return strconv.Itoa(int(x.Value))
}

// MarshalJSON marshals the UInt32Value to JSON.
func (x *UInt32Value) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the UInt32Value from JSON.
func (x *UInt32Value) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a UInt32Value from JSON.
func (x *UInt32Value) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	v := s.ReadUint64()
	if v > math.MaxUint32 {
		s.SetErrorf("value out of range for uint32: %v", v)
		return
	}
	*x = UInt32Value{Value: uint32(v)}
}

// MarshalProtoJSON marshals a UInt32Value to JSON.
func (x *UInt32Value) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteUint64(uint64(x.Value))
}

// String formats the UInt32Value to a string.
func (x *UInt32Value) String() string {
	return strconv.FormatUint(uint64(x.Value), 10)
}

// MarshalJSON marshals the BoolValue to JSON.
func (x *BoolValue) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the BoolValue from JSON.
func (x *BoolValue) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a BoolValue from JSON.
func (x *BoolValue) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	*x = BoolValue{Value: s.ReadBool()}
}

// MarshalProtoJSON marshals a BoolValue to JSON.
func (x *BoolValue) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteBool(x.Value)
}

// String formats the BoolValue to a string.
func (x *BoolValue) String() string {
	return strconv.FormatBool(x.Value)
}

// MarshalJSON marshals the StringValue to JSON.
func (x *StringValue) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the StringValue from JSON.
func (x *StringValue) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a StringValue from JSON.
func (x *StringValue) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	*x = StringValue{Value: s.ReadString()}
}

// MarshalProtoJSON marshals a StringValue to JSON.
func (x *StringValue) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteString(x.Value)
}

// String formats the StringValue to a string.
func (x *StringValue) String() string {
	return x.Value
}

// MarshalJSON marshals the BytesValue to JSON.
func (x *BytesValue) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalJSON unmarshals the BytesValue from JSON.
func (x *BytesValue) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// UnmarshalProtoJSON unmarshals a BytesValue from JSON.
func (x *BytesValue) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	*x = BytesValue{Value: s.ReadBytes()}
}

// MarshalProtoJSON marshals a BytesValue to JSON.
func (x *BytesValue) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteBytes(x.Value)
}

// String formats the BytesValue to a string.
func (x *BytesValue) String() string {
	return string(x.Value)
}
