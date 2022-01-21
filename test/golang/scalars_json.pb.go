// Code generated by protoc-gen-go-json. DO NOT EDIT.
// versions:
// - protoc-gen-go-json v0.0.0-dev
// - protoc             v3.17.3
// source: scalars.proto

package test

import (
	jsonplugin "github.com/TheThingsIndustries/protoc-gen-go-json/jsonplugin"
	types "github.com/TheThingsIndustries/protoc-gen-go-json/test/types"
)

// MarshalProtoJSON marshals the MessageWithScalars message to JSON.
func (x *MessageWithScalars) MarshalProtoJSON(s *jsonplugin.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteObjectStart()
	var wroteField bool
	if x.DoubleValue != 0 || s.HasField("double_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("double_value")
		s.WriteFloat64(x.DoubleValue)
	}
	if len(x.DoubleValues) > 0 || s.HasField("double_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("double_values")
		s.WriteFloat64Array(x.DoubleValues)
	}
	if x.FloatValue != 0 || s.HasField("float_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("float_value")
		s.WriteFloat32(x.FloatValue)
	}
	if len(x.FloatValues) > 0 || s.HasField("float_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("float_values")
		s.WriteFloat32Array(x.FloatValues)
	}
	if x.Int32Value != 0 || s.HasField("int32_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("int32_value")
		s.WriteInt32(x.Int32Value)
	}
	if len(x.Int32Values) > 0 || s.HasField("int32_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("int32_values")
		s.WriteInt32Array(x.Int32Values)
	}
	if x.Int64Value != 0 || s.HasField("int64_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("int64_value")
		s.WriteInt64(x.Int64Value)
	}
	if len(x.Int64Values) > 0 || s.HasField("int64_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("int64_values")
		s.WriteInt64Array(x.Int64Values)
	}
	if x.Uint32Value != 0 || s.HasField("uint32_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("uint32_value")
		s.WriteUint32(x.Uint32Value)
	}
	if len(x.Uint32Values) > 0 || s.HasField("uint32_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("uint32_values")
		s.WriteUint32Array(x.Uint32Values)
	}
	if x.Uint64Value != 0 || s.HasField("uint64_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("uint64_value")
		s.WriteUint64(x.Uint64Value)
	}
	if len(x.Uint64Values) > 0 || s.HasField("uint64_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("uint64_values")
		s.WriteUint64Array(x.Uint64Values)
	}
	if x.Sint32Value != 0 || s.HasField("sint32_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sint32_value")
		s.WriteInt32(x.Sint32Value)
	}
	if len(x.Sint32Values) > 0 || s.HasField("sint32_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sint32_values")
		s.WriteInt32Array(x.Sint32Values)
	}
	if x.Sint64Value != 0 || s.HasField("sint64_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sint64_value")
		s.WriteInt64(x.Sint64Value)
	}
	if len(x.Sint64Values) > 0 || s.HasField("sint64_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sint64_values")
		s.WriteInt64Array(x.Sint64Values)
	}
	if x.Fixed32Value != 0 || s.HasField("fixed32_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("fixed32_value")
		s.WriteUint32(x.Fixed32Value)
	}
	if len(x.Fixed32Values) > 0 || s.HasField("fixed32_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("fixed32_values")
		s.WriteUint32Array(x.Fixed32Values)
	}
	if x.Fixed64Value != 0 || s.HasField("fixed64_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("fixed64_value")
		s.WriteUint64(x.Fixed64Value)
	}
	if len(x.Fixed64Values) > 0 || s.HasField("fixed64_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("fixed64_values")
		s.WriteUint64Array(x.Fixed64Values)
	}
	if x.Sfixed32Value != 0 || s.HasField("sfixed32_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sfixed32_value")
		s.WriteInt32(x.Sfixed32Value)
	}
	if len(x.Sfixed32Values) > 0 || s.HasField("sfixed32_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sfixed32_values")
		s.WriteInt32Array(x.Sfixed32Values)
	}
	if x.Sfixed64Value != 0 || s.HasField("sfixed64_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sfixed64_value")
		s.WriteInt64(x.Sfixed64Value)
	}
	if len(x.Sfixed64Values) > 0 || s.HasField("sfixed64_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sfixed64_values")
		s.WriteInt64Array(x.Sfixed64Values)
	}
	if x.BoolValue || s.HasField("bool_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("bool_value")
		s.WriteBool(x.BoolValue)
	}
	if len(x.BoolValues) > 0 || s.HasField("bool_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("bool_values")
		s.WriteBoolArray(x.BoolValues)
	}
	if x.StringValue != "" || s.HasField("string_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_value")
		s.WriteString(x.StringValue)
	}
	if len(x.StringValues) > 0 || s.HasField("string_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_values")
		s.WriteStringArray(x.StringValues)
	}
	if len(x.BytesValue) > 0 || s.HasField("bytes_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("bytes_value")
		s.WriteBytes(x.BytesValue)
	}
	if len(x.BytesValues) > 0 || s.HasField("bytes_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("bytes_values")
		s.WriteBytesArray(x.BytesValues)
	}
	if len(x.HexBytesValue) > 0 || s.HasField("hex_bytes_value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("hex_bytes_value")
		types.MarshalHEX(s.WithField("hex_bytes_value"), x.HexBytesValue)
	}
	if len(x.HexBytesValues) > 0 || s.HasField("hex_bytes_values") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("hex_bytes_values")
		types.MarshalHEXArray(s.WithField("hex_bytes_values"), x.HexBytesValues)
	}
	s.WriteObjectEnd()
}

// MarshalJSON marshals the MessageWithScalars to JSON.
func (x MessageWithScalars) MarshalJSON() ([]byte, error) {
	return jsonplugin.DefaultMarshalerConfig.Marshal(&x)
}

// UnmarshalProtoJSON unmarshals the MessageWithScalars message from JSON.
func (x *MessageWithScalars) UnmarshalProtoJSON(s *jsonplugin.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	s.ReadObject(func(key string) {
		switch key {
		default:
			s.ReadAny() // ignore unknown field
		case "double_value", "doubleValue":
			s.AddField("double_value")
			x.DoubleValue = s.ReadFloat64()
		case "double_values", "doubleValues":
			s.AddField("double_values")
			if s.ReadNil() {
				x.DoubleValues = nil
				return
			}
			x.DoubleValues = s.ReadFloat64Array()
		case "float_value", "floatValue":
			s.AddField("float_value")
			x.FloatValue = s.ReadFloat32()
		case "float_values", "floatValues":
			s.AddField("float_values")
			if s.ReadNil() {
				x.FloatValues = nil
				return
			}
			x.FloatValues = s.ReadFloat32Array()
		case "int32_value", "int32Value":
			s.AddField("int32_value")
			x.Int32Value = s.ReadInt32()
		case "int32_values", "int32Values":
			s.AddField("int32_values")
			if s.ReadNil() {
				x.Int32Values = nil
				return
			}
			x.Int32Values = s.ReadInt32Array()
		case "int64_value", "int64Value":
			s.AddField("int64_value")
			x.Int64Value = s.ReadInt64()
		case "int64_values", "int64Values":
			s.AddField("int64_values")
			if s.ReadNil() {
				x.Int64Values = nil
				return
			}
			x.Int64Values = s.ReadInt64Array()
		case "uint32_value", "uint32Value":
			s.AddField("uint32_value")
			x.Uint32Value = s.ReadUint32()
		case "uint32_values", "uint32Values":
			s.AddField("uint32_values")
			if s.ReadNil() {
				x.Uint32Values = nil
				return
			}
			x.Uint32Values = s.ReadUint32Array()
		case "uint64_value", "uint64Value":
			s.AddField("uint64_value")
			x.Uint64Value = s.ReadUint64()
		case "uint64_values", "uint64Values":
			s.AddField("uint64_values")
			if s.ReadNil() {
				x.Uint64Values = nil
				return
			}
			x.Uint64Values = s.ReadUint64Array()
		case "sint32_value", "sint32Value":
			s.AddField("sint32_value")
			x.Sint32Value = s.ReadInt32()
		case "sint32_values", "sint32Values":
			s.AddField("sint32_values")
			if s.ReadNil() {
				x.Sint32Values = nil
				return
			}
			x.Sint32Values = s.ReadInt32Array()
		case "sint64_value", "sint64Value":
			s.AddField("sint64_value")
			x.Sint64Value = s.ReadInt64()
		case "sint64_values", "sint64Values":
			s.AddField("sint64_values")
			if s.ReadNil() {
				x.Sint64Values = nil
				return
			}
			x.Sint64Values = s.ReadInt64Array()
		case "fixed32_value", "fixed32Value":
			s.AddField("fixed32_value")
			x.Fixed32Value = s.ReadUint32()
		case "fixed32_values", "fixed32Values":
			s.AddField("fixed32_values")
			if s.ReadNil() {
				x.Fixed32Values = nil
				return
			}
			x.Fixed32Values = s.ReadUint32Array()
		case "fixed64_value", "fixed64Value":
			s.AddField("fixed64_value")
			x.Fixed64Value = s.ReadUint64()
		case "fixed64_values", "fixed64Values":
			s.AddField("fixed64_values")
			if s.ReadNil() {
				x.Fixed64Values = nil
				return
			}
			x.Fixed64Values = s.ReadUint64Array()
		case "sfixed32_value", "sfixed32Value":
			s.AddField("sfixed32_value")
			x.Sfixed32Value = s.ReadInt32()
		case "sfixed32_values", "sfixed32Values":
			s.AddField("sfixed32_values")
			if s.ReadNil() {
				x.Sfixed32Values = nil
				return
			}
			x.Sfixed32Values = s.ReadInt32Array()
		case "sfixed64_value", "sfixed64Value":
			s.AddField("sfixed64_value")
			x.Sfixed64Value = s.ReadInt64()
		case "sfixed64_values", "sfixed64Values":
			s.AddField("sfixed64_values")
			if s.ReadNil() {
				x.Sfixed64Values = nil
				return
			}
			x.Sfixed64Values = s.ReadInt64Array()
		case "bool_value", "boolValue":
			s.AddField("bool_value")
			x.BoolValue = s.ReadBool()
		case "bool_values", "boolValues":
			s.AddField("bool_values")
			if s.ReadNil() {
				x.BoolValues = nil
				return
			}
			x.BoolValues = s.ReadBoolArray()
		case "string_value", "stringValue":
			s.AddField("string_value")
			x.StringValue = s.ReadString()
		case "string_values", "stringValues":
			s.AddField("string_values")
			if s.ReadNil() {
				x.StringValues = nil
				return
			}
			x.StringValues = s.ReadStringArray()
		case "bytes_value", "bytesValue":
			s.AddField("bytes_value")
			x.BytesValue = s.ReadBytes()
		case "bytes_values", "bytesValues":
			s.AddField("bytes_values")
			if s.ReadNil() {
				x.BytesValues = nil
				return
			}
			x.BytesValues = s.ReadBytesArray()
		case "hex_bytes_value", "hexBytesValue":
			s.AddField("hex_bytes_value")
			x.HexBytesValue = types.UnmarshalHEX(s.WithField("hex_bytes_value", false))
		case "hex_bytes_values", "hexBytesValues":
			s.AddField("hex_bytes_values")
			if s.ReadNil() {
				x.HexBytesValues = nil
				return
			}
			x.HexBytesValues = types.UnmarshalHEXArray(s.WithField("hex_bytes_values", false))
		}
	})
}

// UnmarshalJSON unmarshals the MessageWithScalars from JSON.
func (x *MessageWithScalars) UnmarshalJSON(b []byte) error {
	return jsonplugin.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// MarshalProtoJSON marshals the MessageWithOneofScalars message to JSON.
func (x *MessageWithOneofScalars) MarshalProtoJSON(s *jsonplugin.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteObjectStart()
	var wroteField bool
	if x.Value != nil {
		switch ov := x.Value.(type) {
		case *MessageWithOneofScalars_DoubleValue:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("double_value")
			s.WriteFloat64(ov.DoubleValue)
		case *MessageWithOneofScalars_FloatValue:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("float_value")
			s.WriteFloat32(ov.FloatValue)
		case *MessageWithOneofScalars_Int32Value:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("int32_value")
			s.WriteInt32(ov.Int32Value)
		case *MessageWithOneofScalars_Int64Value:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("int64_value")
			s.WriteInt64(ov.Int64Value)
		case *MessageWithOneofScalars_Uint32Value:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("uint32_value")
			s.WriteUint32(ov.Uint32Value)
		case *MessageWithOneofScalars_Uint64Value:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("uint64_value")
			s.WriteUint64(ov.Uint64Value)
		case *MessageWithOneofScalars_Sint32Value:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("sint32_value")
			s.WriteInt32(ov.Sint32Value)
		case *MessageWithOneofScalars_Sint64Value:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("sint64_value")
			s.WriteInt64(ov.Sint64Value)
		case *MessageWithOneofScalars_Fixed32Value:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("fixed32_value")
			s.WriteUint32(ov.Fixed32Value)
		case *MessageWithOneofScalars_Fixed64Value:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("fixed64_value")
			s.WriteUint64(ov.Fixed64Value)
		case *MessageWithOneofScalars_Sfixed32Value:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("sfixed32_value")
			s.WriteInt32(ov.Sfixed32Value)
		case *MessageWithOneofScalars_Sfixed64Value:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("sfixed64_value")
			s.WriteInt64(ov.Sfixed64Value)
		case *MessageWithOneofScalars_BoolValue:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("bool_value")
			s.WriteBool(ov.BoolValue)
		case *MessageWithOneofScalars_StringValue:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("string_value")
			s.WriteString(ov.StringValue)
		case *MessageWithOneofScalars_BytesValue:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("bytes_value")
			s.WriteBytes(ov.BytesValue)
		case *MessageWithOneofScalars_HexBytesValue:
			s.WriteMoreIf(&wroteField)
			s.WriteObjectField("hex_bytes_value")
			types.MarshalHEX(s.WithField("hex_bytes_value"), ov.HexBytesValue)
		}
	}
	s.WriteObjectEnd()
}

// MarshalJSON marshals the MessageWithOneofScalars to JSON.
func (x MessageWithOneofScalars) MarshalJSON() ([]byte, error) {
	return jsonplugin.DefaultMarshalerConfig.Marshal(&x)
}

// UnmarshalProtoJSON unmarshals the MessageWithOneofScalars message from JSON.
func (x *MessageWithOneofScalars) UnmarshalProtoJSON(s *jsonplugin.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	s.ReadObject(func(key string) {
		switch key {
		default:
			s.ReadAny() // ignore unknown field
		case "double_value", "doubleValue":
			s.AddField("double_value")
			ov := &MessageWithOneofScalars_DoubleValue{}
			x.Value = ov
			ov.DoubleValue = s.ReadFloat64()
		case "float_value", "floatValue":
			s.AddField("float_value")
			ov := &MessageWithOneofScalars_FloatValue{}
			x.Value = ov
			ov.FloatValue = s.ReadFloat32()
		case "int32_value", "int32Value":
			s.AddField("int32_value")
			ov := &MessageWithOneofScalars_Int32Value{}
			x.Value = ov
			ov.Int32Value = s.ReadInt32()
		case "int64_value", "int64Value":
			s.AddField("int64_value")
			ov := &MessageWithOneofScalars_Int64Value{}
			x.Value = ov
			ov.Int64Value = s.ReadInt64()
		case "uint32_value", "uint32Value":
			s.AddField("uint32_value")
			ov := &MessageWithOneofScalars_Uint32Value{}
			x.Value = ov
			ov.Uint32Value = s.ReadUint32()
		case "uint64_value", "uint64Value":
			s.AddField("uint64_value")
			ov := &MessageWithOneofScalars_Uint64Value{}
			x.Value = ov
			ov.Uint64Value = s.ReadUint64()
		case "sint32_value", "sint32Value":
			s.AddField("sint32_value")
			ov := &MessageWithOneofScalars_Sint32Value{}
			x.Value = ov
			ov.Sint32Value = s.ReadInt32()
		case "sint64_value", "sint64Value":
			s.AddField("sint64_value")
			ov := &MessageWithOneofScalars_Sint64Value{}
			x.Value = ov
			ov.Sint64Value = s.ReadInt64()
		case "fixed32_value", "fixed32Value":
			s.AddField("fixed32_value")
			ov := &MessageWithOneofScalars_Fixed32Value{}
			x.Value = ov
			ov.Fixed32Value = s.ReadUint32()
		case "fixed64_value", "fixed64Value":
			s.AddField("fixed64_value")
			ov := &MessageWithOneofScalars_Fixed64Value{}
			x.Value = ov
			ov.Fixed64Value = s.ReadUint64()
		case "sfixed32_value", "sfixed32Value":
			s.AddField("sfixed32_value")
			ov := &MessageWithOneofScalars_Sfixed32Value{}
			x.Value = ov
			ov.Sfixed32Value = s.ReadInt32()
		case "sfixed64_value", "sfixed64Value":
			s.AddField("sfixed64_value")
			ov := &MessageWithOneofScalars_Sfixed64Value{}
			x.Value = ov
			ov.Sfixed64Value = s.ReadInt64()
		case "bool_value", "boolValue":
			s.AddField("bool_value")
			ov := &MessageWithOneofScalars_BoolValue{}
			x.Value = ov
			ov.BoolValue = s.ReadBool()
		case "string_value", "stringValue":
			s.AddField("string_value")
			ov := &MessageWithOneofScalars_StringValue{}
			x.Value = ov
			ov.StringValue = s.ReadString()
		case "bytes_value", "bytesValue":
			s.AddField("bytes_value")
			ov := &MessageWithOneofScalars_BytesValue{}
			x.Value = ov
			ov.BytesValue = s.ReadBytes()
		case "hex_bytes_value", "hexBytesValue":
			s.AddField("hex_bytes_value")
			ov := &MessageWithOneofScalars_HexBytesValue{}
			x.Value = ov
			ov.HexBytesValue = types.UnmarshalHEX(s.WithField("hex_bytes_value", false))
		}
	})
}

// UnmarshalJSON unmarshals the MessageWithOneofScalars from JSON.
func (x *MessageWithOneofScalars) UnmarshalJSON(b []byte) error {
	return jsonplugin.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// MarshalProtoJSON marshals the MessageWithScalarMaps message to JSON.
func (x *MessageWithScalarMaps) MarshalProtoJSON(s *jsonplugin.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteObjectStart()
	var wroteField bool
	if x.StringDoubleMap != nil || s.HasField("string_double_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_double_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringDoubleMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteFloat64(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringFloatMap != nil || s.HasField("string_float_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_float_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringFloatMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteFloat32(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringInt32Map != nil || s.HasField("string_int32_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_int32_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringInt32Map {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteInt32(v)
		}
		s.WriteObjectEnd()
	}
	if x.Int32StringMap != nil || s.HasField("int32_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("int32_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.Int32StringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectInt32Field(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringInt64Map != nil || s.HasField("string_int64_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_int64_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringInt64Map {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteInt64(v)
		}
		s.WriteObjectEnd()
	}
	if x.Int64StringMap != nil || s.HasField("int64_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("int64_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.Int64StringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectInt64Field(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringUint32Map != nil || s.HasField("string_uint32_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_uint32_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringUint32Map {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteUint32(v)
		}
		s.WriteObjectEnd()
	}
	if x.Uint32StringMap != nil || s.HasField("uint32_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("uint32_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.Uint32StringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectUint32Field(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringUint64Map != nil || s.HasField("string_uint64_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_uint64_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringUint64Map {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteUint64(v)
		}
		s.WriteObjectEnd()
	}
	if x.Uint64StringMap != nil || s.HasField("uint64_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("uint64_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.Uint64StringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectUint64Field(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringSint32Map != nil || s.HasField("string_sint32_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_sint32_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringSint32Map {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteInt32(v)
		}
		s.WriteObjectEnd()
	}
	if x.Sint32StringMap != nil || s.HasField("sint32_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sint32_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.Sint32StringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectInt32Field(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringSint64Map != nil || s.HasField("string_sint64_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_sint64_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringSint64Map {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteInt64(v)
		}
		s.WriteObjectEnd()
	}
	if x.Sint64StringMap != nil || s.HasField("sint64_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sint64_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.Sint64StringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectInt64Field(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringFixed32Map != nil || s.HasField("string_fixed32_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_fixed32_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringFixed32Map {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteUint32(v)
		}
		s.WriteObjectEnd()
	}
	if x.Fixed32StringMap != nil || s.HasField("fixed32_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("fixed32_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.Fixed32StringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectUint32Field(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringFixed64Map != nil || s.HasField("string_fixed64_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_fixed64_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringFixed64Map {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteUint64(v)
		}
		s.WriteObjectEnd()
	}
	if x.Fixed64StringMap != nil || s.HasField("fixed64_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("fixed64_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.Fixed64StringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectUint64Field(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringSfixed32Map != nil || s.HasField("string_sfixed32_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_sfixed32_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringSfixed32Map {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteInt32(v)
		}
		s.WriteObjectEnd()
	}
	if x.Sfixed32StringMap != nil || s.HasField("sfixed32_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sfixed32_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.Sfixed32StringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectInt32Field(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringSfixed64Map != nil || s.HasField("string_sfixed64_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_sfixed64_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringSfixed64Map {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteInt64(v)
		}
		s.WriteObjectEnd()
	}
	if x.Sfixed64StringMap != nil || s.HasField("sfixed64_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sfixed64_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.Sfixed64StringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectInt64Field(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringBoolMap != nil || s.HasField("string_bool_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_bool_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringBoolMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteBool(v)
		}
		s.WriteObjectEnd()
	}
	if x.BoolStringMap != nil || s.HasField("bool_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("bool_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.BoolStringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectBoolField(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringStringMap != nil || s.HasField("string_string_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_string_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringStringMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteString(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringBytesMap != nil || s.HasField("string_bytes_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_bytes_map")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringBytesMap {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			s.WriteBytes(v)
		}
		s.WriteObjectEnd()
	}
	if x.StringHexBytesMap != nil || s.HasField("string_hex_bytes_map") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("string_hex_bytes_map")
		types.MarshalStringHEXMap(s.WithField("string_hex_bytes_map"), x.StringHexBytesMap)
	}
	s.WriteObjectEnd()
}

// MarshalJSON marshals the MessageWithScalarMaps to JSON.
func (x MessageWithScalarMaps) MarshalJSON() ([]byte, error) {
	return jsonplugin.DefaultMarshalerConfig.Marshal(&x)
}

// UnmarshalProtoJSON unmarshals the MessageWithScalarMaps message from JSON.
func (x *MessageWithScalarMaps) UnmarshalProtoJSON(s *jsonplugin.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	s.ReadObject(func(key string) {
		switch key {
		default:
			s.ReadAny() // ignore unknown field
		case "string_double_map", "stringDoubleMap":
			s.AddField("string_double_map")
			if s.ReadNil() {
				x.StringDoubleMap = nil
				return
			}
			x.StringDoubleMap = make(map[string]float64)
			s.ReadStringMap(func(key string) {
				x.StringDoubleMap[key] = s.ReadFloat64()
			})
		case "string_float_map", "stringFloatMap":
			s.AddField("string_float_map")
			if s.ReadNil() {
				x.StringFloatMap = nil
				return
			}
			x.StringFloatMap = make(map[string]float32)
			s.ReadStringMap(func(key string) {
				x.StringFloatMap[key] = s.ReadFloat32()
			})
		case "string_int32_map", "stringInt32Map":
			s.AddField("string_int32_map")
			if s.ReadNil() {
				x.StringInt32Map = nil
				return
			}
			x.StringInt32Map = make(map[string]int32)
			s.ReadStringMap(func(key string) {
				x.StringInt32Map[key] = s.ReadInt32()
			})
		case "int32_string_map", "int32StringMap":
			s.AddField("int32_string_map")
			if s.ReadNil() {
				x.Int32StringMap = nil
				return
			}
			x.Int32StringMap = make(map[int32]string)
			s.ReadInt32Map(func(key int32) {
				x.Int32StringMap[key] = s.ReadString()
			})
		case "string_int64_map", "stringInt64Map":
			s.AddField("string_int64_map")
			if s.ReadNil() {
				x.StringInt64Map = nil
				return
			}
			x.StringInt64Map = make(map[string]int64)
			s.ReadStringMap(func(key string) {
				x.StringInt64Map[key] = s.ReadInt64()
			})
		case "int64_string_map", "int64StringMap":
			s.AddField("int64_string_map")
			if s.ReadNil() {
				x.Int64StringMap = nil
				return
			}
			x.Int64StringMap = make(map[int64]string)
			s.ReadInt64Map(func(key int64) {
				x.Int64StringMap[key] = s.ReadString()
			})
		case "string_uint32_map", "stringUint32Map":
			s.AddField("string_uint32_map")
			if s.ReadNil() {
				x.StringUint32Map = nil
				return
			}
			x.StringUint32Map = make(map[string]uint32)
			s.ReadStringMap(func(key string) {
				x.StringUint32Map[key] = s.ReadUint32()
			})
		case "uint32_string_map", "uint32StringMap":
			s.AddField("uint32_string_map")
			if s.ReadNil() {
				x.Uint32StringMap = nil
				return
			}
			x.Uint32StringMap = make(map[uint32]string)
			s.ReadUint32Map(func(key uint32) {
				x.Uint32StringMap[key] = s.ReadString()
			})
		case "string_uint64_map", "stringUint64Map":
			s.AddField("string_uint64_map")
			if s.ReadNil() {
				x.StringUint64Map = nil
				return
			}
			x.StringUint64Map = make(map[string]uint64)
			s.ReadStringMap(func(key string) {
				x.StringUint64Map[key] = s.ReadUint64()
			})
		case "uint64_string_map", "uint64StringMap":
			s.AddField("uint64_string_map")
			if s.ReadNil() {
				x.Uint64StringMap = nil
				return
			}
			x.Uint64StringMap = make(map[uint64]string)
			s.ReadUint64Map(func(key uint64) {
				x.Uint64StringMap[key] = s.ReadString()
			})
		case "string_sint32_map", "stringSint32Map":
			s.AddField("string_sint32_map")
			if s.ReadNil() {
				x.StringSint32Map = nil
				return
			}
			x.StringSint32Map = make(map[string]int32)
			s.ReadStringMap(func(key string) {
				x.StringSint32Map[key] = s.ReadInt32()
			})
		case "sint32_string_map", "sint32StringMap":
			s.AddField("sint32_string_map")
			if s.ReadNil() {
				x.Sint32StringMap = nil
				return
			}
			x.Sint32StringMap = make(map[int32]string)
			s.ReadInt32Map(func(key int32) {
				x.Sint32StringMap[key] = s.ReadString()
			})
		case "string_sint64_map", "stringSint64Map":
			s.AddField("string_sint64_map")
			if s.ReadNil() {
				x.StringSint64Map = nil
				return
			}
			x.StringSint64Map = make(map[string]int64)
			s.ReadStringMap(func(key string) {
				x.StringSint64Map[key] = s.ReadInt64()
			})
		case "sint64_string_map", "sint64StringMap":
			s.AddField("sint64_string_map")
			if s.ReadNil() {
				x.Sint64StringMap = nil
				return
			}
			x.Sint64StringMap = make(map[int64]string)
			s.ReadInt64Map(func(key int64) {
				x.Sint64StringMap[key] = s.ReadString()
			})
		case "string_fixed32_map", "stringFixed32Map":
			s.AddField("string_fixed32_map")
			if s.ReadNil() {
				x.StringFixed32Map = nil
				return
			}
			x.StringFixed32Map = make(map[string]uint32)
			s.ReadStringMap(func(key string) {
				x.StringFixed32Map[key] = s.ReadUint32()
			})
		case "fixed32_string_map", "fixed32StringMap":
			s.AddField("fixed32_string_map")
			if s.ReadNil() {
				x.Fixed32StringMap = nil
				return
			}
			x.Fixed32StringMap = make(map[uint32]string)
			s.ReadUint32Map(func(key uint32) {
				x.Fixed32StringMap[key] = s.ReadString()
			})
		case "string_fixed64_map", "stringFixed64Map":
			s.AddField("string_fixed64_map")
			if s.ReadNil() {
				x.StringFixed64Map = nil
				return
			}
			x.StringFixed64Map = make(map[string]uint64)
			s.ReadStringMap(func(key string) {
				x.StringFixed64Map[key] = s.ReadUint64()
			})
		case "fixed64_string_map", "fixed64StringMap":
			s.AddField("fixed64_string_map")
			if s.ReadNil() {
				x.Fixed64StringMap = nil
				return
			}
			x.Fixed64StringMap = make(map[uint64]string)
			s.ReadUint64Map(func(key uint64) {
				x.Fixed64StringMap[key] = s.ReadString()
			})
		case "string_sfixed32_map", "stringSfixed32Map":
			s.AddField("string_sfixed32_map")
			if s.ReadNil() {
				x.StringSfixed32Map = nil
				return
			}
			x.StringSfixed32Map = make(map[string]int32)
			s.ReadStringMap(func(key string) {
				x.StringSfixed32Map[key] = s.ReadInt32()
			})
		case "sfixed32_string_map", "sfixed32StringMap":
			s.AddField("sfixed32_string_map")
			if s.ReadNil() {
				x.Sfixed32StringMap = nil
				return
			}
			x.Sfixed32StringMap = make(map[int32]string)
			s.ReadInt32Map(func(key int32) {
				x.Sfixed32StringMap[key] = s.ReadString()
			})
		case "string_sfixed64_map", "stringSfixed64Map":
			s.AddField("string_sfixed64_map")
			if s.ReadNil() {
				x.StringSfixed64Map = nil
				return
			}
			x.StringSfixed64Map = make(map[string]int64)
			s.ReadStringMap(func(key string) {
				x.StringSfixed64Map[key] = s.ReadInt64()
			})
		case "sfixed64_string_map", "sfixed64StringMap":
			s.AddField("sfixed64_string_map")
			if s.ReadNil() {
				x.Sfixed64StringMap = nil
				return
			}
			x.Sfixed64StringMap = make(map[int64]string)
			s.ReadInt64Map(func(key int64) {
				x.Sfixed64StringMap[key] = s.ReadString()
			})
		case "string_bool_map", "stringBoolMap":
			s.AddField("string_bool_map")
			if s.ReadNil() {
				x.StringBoolMap = nil
				return
			}
			x.StringBoolMap = make(map[string]bool)
			s.ReadStringMap(func(key string) {
				x.StringBoolMap[key] = s.ReadBool()
			})
		case "bool_string_map", "boolStringMap":
			s.AddField("bool_string_map")
			if s.ReadNil() {
				x.BoolStringMap = nil
				return
			}
			x.BoolStringMap = make(map[bool]string)
			s.ReadBoolMap(func(key bool) {
				x.BoolStringMap[key] = s.ReadString()
			})
		case "string_string_map", "stringStringMap":
			s.AddField("string_string_map")
			if s.ReadNil() {
				x.StringStringMap = nil
				return
			}
			x.StringStringMap = make(map[string]string)
			s.ReadStringMap(func(key string) {
				x.StringStringMap[key] = s.ReadString()
			})
		case "string_bytes_map", "stringBytesMap":
			s.AddField("string_bytes_map")
			if s.ReadNil() {
				x.StringBytesMap = nil
				return
			}
			x.StringBytesMap = make(map[string][]byte)
			s.ReadStringMap(func(key string) {
				x.StringBytesMap[key] = s.ReadBytes()
			})
		case "string_hex_bytes_map", "stringHexBytesMap":
			s.AddField("string_hex_bytes_map")
			if s.ReadNil() {
				x.StringHexBytesMap = nil
				return
			}
			x.StringHexBytesMap = types.UnmarshalStringHEXMap(s.WithField("string_hex_bytes_map", false))
		}
	})
}

// UnmarshalJSON unmarshals the MessageWithScalarMaps from JSON.
func (x *MessageWithScalarMaps) UnmarshalJSON(b []byte) error {
	return jsonplugin.DefaultUnmarshalerConfig.Unmarshal(b, x)
}
