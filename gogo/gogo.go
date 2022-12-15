// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package gogoplugin

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/TheThingsIndustries/protoc-gen-go-json/jsonplugin"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
)

// MarshalMessage marshals a message with the standard JSON marshaler.
func MarshalMessage(s *jsonplugin.MarshalState, v proto.Message) {
	rv := reflect.ValueOf(v)

	// If v is nil or typed nil, we write null.
	if v == nil || (rv.Kind() == reflect.Ptr && rv.IsNil()) {
		s.WriteNil()
		return
	}

	// Let gogo jsonpb marshal v.
	err := (&jsonpb.Marshaler{
		OrigName:    true,
		EnumsAsInts: s.Config().EnumsAsInts,
	}).Marshal(s, v)
	if err != nil {
		s.SetErrorf("failed to marshal %s to JSON: %w", proto.MessageName(v), err)
	}
}

// UnmarshalMessage unmarshals a message with the standard JSON unmarshaler.
func UnmarshalMessage(s *jsonplugin.UnmarshalState, v proto.Message) {
	// If we read null, don't do anything.
	if s.ReadNil() {
		return
	}

	// Read the raw object.
	data := s.ReadRawMessage()
	if s.Err() != nil {
		return
	}

	// Let gogo jsonpb unmarshal v.
	err := (&jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}).Unmarshal(bytes.NewBuffer(data), v)
	if err != nil {
		s.SetErrorf("failed to unmarshal %s from JSON: %w", proto.MessageName(v), err)
	}
}

// MarshalAny marshals a Any WKT.
func MarshalAny(s *jsonplugin.MarshalState, v *types.Any, legacyFieldmask bool) {
	if v == nil {
		s.WriteNil()
		return
	}

	// We first need to get the wrapped message out of the Any.
	// To do this, we instantiate an empty message of the type of the wrapped message.
	// Then we unmarshal the Any into that empty message.
	msg, err := types.EmptyAny(v)
	if err != nil {
		s.SetErrorf("unknown message type %q for Any: %w", v.GetTypeUrl(), err)
	}
	if err = types.UnmarshalAny(v, msg); err != nil {
		s.SetErrorf("failed to unmarshal wrapped message from Any: %w", err)
	}

	switch marshaler := msg.(type) {
	default:
		// If v doesn't implement jsonplugin.Marshaler, delegate to gogo jsonpb.
		MarshalMessage(s, v)
	case jsonplugin.Marshaler:
		// Instantiate a sub-marshaler with the same configuration and marshal the wrapped message to that.
		sub := s.Sub()
		marshaler.MarshalProtoJSON(sub)
		data, err := sub.Bytes()
		if err != nil {
			return
		}

		// We need to prepend the @type field to that object, so we read the { character.
		buf := bytes.NewBuffer(data)
		objectStart, err := buf.ReadByte()
		if err != nil {
			s.SetError(err)
			return
		}
		if objectStart != '{' {
			s.SetErrorf("marshaled Any is not an object")
			return
		}

		// We take a look at the next token, because if it's a ", we'll need a comma after we write the @type field.
		nextToken, err := buf.ReadByte()
		if err != nil {
			s.SetError(err)
			return
		}
		buf.UnreadByte()

		// Write the opening { and the type field to the main marshaler.
		s.WriteObjectStart()
		s.WriteObjectField("@type")
		s.WriteString(v.GetTypeUrl())

		// If the next token is a ", we have more fields, so we need to write a comma.
		// Otherwise, it's a } and we don't need a comma.
		if nextToken == '"' {
			s.WriteMore()
		}

		// Write the rest of the buffer (the sub-object without the { character).
		s.Write(buf.Bytes())
	case *types.Duration,
		*types.FieldMask,
		*types.Struct,
		*types.Value,
		*types.ListValue,
		*types.Timestamp:

		// Write the opening { and the type field to the main marshaler.
		s.WriteObjectStart()
		s.WriteObjectField("@type")
		s.WriteString(v.GetTypeUrl())

		// Write the comma, and the next field, which is always "value" for these types.
		s.WriteMore()
		s.WriteObjectField("value")

		// Write the value.
		switch msg := msg.(type) {
		case *types.Duration:
			MarshalDuration(s, msg)
		case *types.FieldMask:
			if legacyFieldmask {
				MarshalLegacyFieldMask(s, msg)
			} else {
				MarshalFieldMask(s, msg)
			}
		case *types.Struct:
			MarshalStruct(s, msg)
		case *types.Value:
			MarshalValue(s, msg)
		case *types.ListValue:
			MarshalListValue(s, msg)
		case *types.Timestamp:
			MarshalTimestamp(s, msg)
		}

		// Write the closing }.
		s.WriteObjectEnd()
	}
}

// UnmarshalAny unmarshals an Any WKT.
func UnmarshalAny(s *jsonplugin.UnmarshalState) *types.Any {
	if s.ReadNil() {
		return nil
	}

	// Read the raw object and create a sub-unmarshaler for it.
	data := s.ReadRawMessage()
	if s.Err() != nil {
		return nil
	}
	sub := s.Sub(data)

	// Read the first field in the object. This should be @type.
	if key := sub.ReadObjectField(); key != "@type" {
		s.SetErrorf("first field in Any is not @type, but %q", key)
		return nil
	}
	typeURL := sub.ReadString()
	if err := sub.Err(); err != nil {
		return nil
	}

	// Find the message type by the name that's in the type URL.
	slash := strings.LastIndex(typeURL, "/")
	if slash < 0 {
		s.SetErrorf("invalid type URL %q", typeURL)
		return nil
	}
	t := proto.MessageType(typeURL[slash+1:])
	if t == nil {
		s.SetErrorf("unknown message type %q", typeURL[slash+1:])
		return nil
	}

	// Allocate a new message of that type.
	msg := reflect.New(t.Elem()).Interface().(proto.Message)

	switch unmarshaler := msg.(type) {
	default:
		// Delegate unmarshaling to gogo jsonpb.
		var any types.Any
		if err := (&jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}).Unmarshal(bytes.NewBuffer(data), &any); err != nil {
			s.SetErrorf("failed to unmarshal Any from JSON: %w", err)
			return nil
		}
		return &any
	case jsonplugin.Unmarshaler:
		// Create another sub-unmarshaler for the raw data and unmarshal the message.
		sub = s.Sub(data)
		unmarshaler.UnmarshalProtoJSON(sub)
	case *types.Duration,
		*types.FieldMask,
		*types.Struct,
		*types.Value,
		*types.ListValue,
		*types.Timestamp:
		if field := sub.ReadObjectField(); field != "value" {
			s.SetErrorf("unexpected %q field in Any", field)
			return nil
		}
		switch msg.(type) {
		case *types.Duration:
			msg = UnmarshalDuration(sub)
		case *types.FieldMask:
			msg = UnmarshalFieldMask(sub)
		case *types.Struct:
			msg = UnmarshalStruct(sub)
		case *types.Value:
			msg = UnmarshalValue(sub)
		case *types.ListValue:
			msg = UnmarshalListValue(sub)
		case *types.Timestamp:
			msg = UnmarshalTimestamp(sub)
		}
	}

	if err := sub.Err(); err != nil {
		return nil
	}

	if field := sub.ReadObjectField(); field != "" {
		s.SetErrorf("unexpected %q field in Any", field)
		return nil
	}

	// Wrap the unmarshaled message in an Any and return that.
	v, err := types.MarshalAny(msg)
	if err != nil {
		s.SetError(err)
		return nil
	}
	return v
}

// MarshalDuration marshals a Duration WKT.
func MarshalDuration(s *jsonplugin.MarshalState, v *types.Duration) {
	if v == nil {
		s.WriteNil()
		return
	}
	d, err := types.DurationFromProto(v)
	if err != nil {
		s.SetErrorf("invalid duration: %w", err)
	}
	s.WriteDuration(d)
}

// UnmarshalDuration unmarshals a Duration WKT.
func UnmarshalDuration(s *jsonplugin.UnmarshalState) *types.Duration {
	if s.ReadNil() {
		return nil
	}
	d := s.ReadDuration()
	if s.Err() != nil {
		return nil
	}
	return types.DurationProto(*d)
}

// MarshalEmpty marshals an Empty WKT.
func MarshalEmpty(s *jsonplugin.MarshalState, _ *types.Empty) {
	s.WriteObjectStart()
	s.WriteObjectEnd()
}

// UnmarshalEmpty unmarshals a Empty WKT.
func UnmarshalEmpty(s *jsonplugin.UnmarshalState) *types.Empty {
	if s.ReadNil() {
		return nil
	}
	s.ReadObject(func(key string) {
		s.SetErrorf("unexpected key %q in Empty", key)
	})
	if s.Err() != nil {
		return nil
	}
	return &types.Empty{}
}

// MarshalFieldMask marshals a FieldMask WKT.
func MarshalFieldMask(s *jsonplugin.MarshalState, v *types.FieldMask) {
	if v == nil {
		s.WriteNil()
		return
	}
	s.WriteFieldMask(v)
}

func MarshalLegacyFieldMask(s *jsonplugin.MarshalState, v *types.FieldMask) {
	if v == nil {
		s.WriteNil()
		return
	}
	s.WriteLegacyFieldMask(v)
}

// UnmarshalFieldMask unmarshals a FieldMask WKT.
func UnmarshalFieldMask(s *jsonplugin.UnmarshalState) *types.FieldMask {
	if s.ReadNil() {
		return nil
	}
	m := s.ReadFieldMask()
	if s.Err() != nil {
		return nil
	}
	return &types.FieldMask{Paths: m.GetPaths()}
}

// MarshalStruct marshals a Struct WKT.
func MarshalStruct(s *jsonplugin.MarshalState, v *types.Struct) {
	if v == nil {
		s.WriteNil()
		return
	}
	MarshalMessage(s, v)
}

// UnmarshalStruct unmarshals a Struct WKT.
func UnmarshalStruct(s *jsonplugin.UnmarshalState) *types.Struct {
	if s.ReadNil() {
		return nil
	}
	var v types.Struct
	UnmarshalMessage(s, &v)
	if s.Err() != nil {
		return nil
	}
	return &v
}

// MarshalValue marshals a Value WKT.
func MarshalValue(s *jsonplugin.MarshalState, v *types.Value) {
	if v == nil {
		s.WriteNil()
		return
	}
	MarshalMessage(s, v)
}

// UnmarshalValue unmarshals a Value WKT.
func UnmarshalValue(s *jsonplugin.UnmarshalState) *types.Value {
	if s.ReadNil() {
		return &types.Value{Kind: &types.Value_NullValue{}}
	}
	var v types.Value
	UnmarshalMessage(s, &v)
	if s.Err() != nil {
		return nil
	}
	return &v
}

// MarshalListValue marshals a ListValue WKT.
func MarshalListValue(s *jsonplugin.MarshalState, v *types.ListValue) {
	if v == nil {
		s.WriteNil()
		return
	}
	MarshalMessage(s, v)
}

// UnmarshalListValue unmarshals a ListValue WKT.
func UnmarshalListValue(s *jsonplugin.UnmarshalState) *types.ListValue {
	if s.ReadNil() {
		return nil
	}
	var v types.ListValue
	UnmarshalMessage(s, &v)
	if s.Err() != nil {
		return nil
	}
	return &v
}

// MarshalTimestamp marshals a Timestamp WKT.
func MarshalTimestamp(s *jsonplugin.MarshalState, v *types.Timestamp) {
	if v == nil {
		s.WriteNil()
		return
	}
	t, err := types.TimestampFromProto(v)
	if err != nil {
		s.SetErrorf("invalid time: %w", err)
	}
	s.WriteTime(t)
}

// UnmarshalTimestamp unmarshals a Timestamp WKT.
func UnmarshalTimestamp(s *jsonplugin.UnmarshalState) *types.Timestamp {
	if s.ReadNil() {
		return nil
	}
	d := s.ReadTime()
	if s.Err() != nil {
		return nil
	}
	v, err := types.TimestampProto(*d)
	if err != nil {
		s.SetErrorf("invalid time: %w", err)
		return nil
	}
	return v
}
