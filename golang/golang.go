// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package golangplugin

import (
	"bytes"
	"reflect"

	"github.com/TheThingsIndustries/protoc-gen-go-json/jsonplugin"
	protojson "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	anypb "google.golang.org/protobuf/types/known/anypb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// MarshalMessage marshals a message with the standard JSON marshaler.
func MarshalMessage(s *jsonplugin.MarshalState, v proto.Message) {
	rv := reflect.ValueOf(v)

	// If v is nil or typed nil, we write null.
	if v == nil || (rv.Kind() == reflect.Ptr && rv.IsNil()) {
		s.WriteNil()
		return
	}

	// Let protojson marshal v.
	b, err := protojson.MarshalOptions{
		UseProtoNames:  true,
		UseEnumNumbers: s.Config().EnumsAsInts,
	}.Marshal(v)
	if err != nil {
		s.SetErrorf("failed to marshal %s to JSON: %w", proto.MessageName(v), err)
	}
	s.Write(b)
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

	// Let protojson unmarshal v.
	err := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(data, v)
	if err != nil {
		s.SetErrorf("failed to unmarshal %s from JSON: %w", proto.MessageName(v), err)
	}
}

// MarshalAny marshals an Any WKT.
func MarshalAny(s *jsonplugin.MarshalState, v *anypb.Any, legacyFieldmask bool) {
	if v == nil {
		s.WriteNil()
		return
	}

	// We first need to get the wrapped message out of the Any.
	msg, err := v.UnmarshalNew()
	if err != nil {
		s.SetErrorf("failed to unmarshal wrapped message from Any: %w", err)
	}

	switch marshaler := msg.(type) {
	default:
		// If v doesn't implement jsonplugin.Marshaler, delegate to protojson.
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
	case *durationpb.Duration,
		*fieldmaskpb.FieldMask,
		*structpb.Struct,
		*structpb.Value,
		*structpb.ListValue,
		*timestamppb.Timestamp:

		// Write the opening { and the type field to the main marshaler.
		s.WriteObjectStart()
		s.WriteObjectField("@type")
		s.WriteString(v.GetTypeUrl())

		// Write the comma, and the next field, which is always "value" for these types.
		s.WriteMore()
		s.WriteObjectField("value")

		// Write the value.
		switch msg := msg.(type) {
		case *durationpb.Duration:
			MarshalDuration(s, msg)
		case *fieldmaskpb.FieldMask:
			if legacyFieldmask {
				MarshalLegacyFieldMask(s, msg)
			} else {
				MarshalFieldMask(s, msg)
			}
		case *structpb.Struct:
			MarshalStruct(s, msg)
		case *structpb.Value:
			MarshalValue(s, msg)
		case *structpb.ListValue:
			MarshalListValue(s, msg)
		case *timestamppb.Timestamp:
			MarshalTimestamp(s, msg)
		}

		// Write the closing }.
		s.WriteObjectEnd()
	}
}

// UnmarshalAny unmarshals an Any WKT.
func UnmarshalAny(s *jsonplugin.UnmarshalState) *anypb.Any {
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

	// Find the message type by the type URL.
	t, err := protoregistry.GlobalTypes.FindMessageByURL(typeURL)
	if err != nil {
		s.SetError(err)
		return nil
	}

	// Allocate a new message of that type.
	msg := t.New().Interface()

	switch unmarshaler := msg.(type) {
	default:
		// Delegate unmarshaling to protojson.
		var any anypb.Any
		if err := (&protojson.UnmarshalOptions{
			DiscardUnknown: true,
		}).Unmarshal(data, &any); err != nil {
			s.SetErrorf("failed to unmarshal Any from JSON: %w", err)
			return nil
		}
		return &any
	case jsonplugin.Unmarshaler:
		// Create another sub-unmarshaler for the raw data and unmarshal the message.
		sub = s.Sub(data)
		unmarshaler.UnmarshalProtoJSON(sub)
	case *durationpb.Duration,
		*fieldmaskpb.FieldMask,
		*structpb.Struct,
		*structpb.Value,
		*structpb.ListValue,
		*timestamppb.Timestamp:
		if field := sub.ReadObjectField(); field != "value" {
			s.SetErrorf("unexpected %q field in Any", field)
			return nil
		}
		switch msg.(type) {
		case *durationpb.Duration:
			msg = UnmarshalDuration(sub)
		case *fieldmaskpb.FieldMask:
			msg = UnmarshalFieldMask(sub)
		case *structpb.Struct:
			msg = UnmarshalStruct(sub)
		case *structpb.Value:
			msg = UnmarshalValue(sub)
		case *structpb.ListValue:
			msg = UnmarshalListValue(sub)
		case *timestamppb.Timestamp:
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
	v, err := anypb.New(msg)
	if err != nil {
		s.SetError(err)
		return nil
	}
	return v
}

// MarshalDuration marshals a Duration WKT.
func MarshalDuration(s *jsonplugin.MarshalState, v *durationpb.Duration) {
	if v == nil {
		s.WriteNil()
		return
	}
	s.WriteDuration(v.AsDuration())
}

// UnmarshalDuration unmarshals a Duration WKT.
func UnmarshalDuration(s *jsonplugin.UnmarshalState) *durationpb.Duration {
	if s.ReadNil() {
		return nil
	}
	d := s.ReadDuration()
	if s.Err() != nil {
		return nil
	}
	return durationpb.New(*d)
}

// MarshalEmpty marshals an Empty WKT.
func MarshalEmpty(s *jsonplugin.MarshalState, _ *emptypb.Empty) {
	s.WriteObjectStart()
	s.WriteObjectEnd()
}

// UnmarshalEmpty unmarshals a Empty WKT.
func UnmarshalEmpty(s *jsonplugin.UnmarshalState) *emptypb.Empty {
	if s.ReadNil() {
		return nil
	}
	s.ReadObject(func(key string) {
		s.SetErrorf("unexpected key %q in Empty", key)
	})
	if s.Err() != nil {
		return nil
	}
	return &emptypb.Empty{}
}

// MarshalFieldMask marshals a FieldMask WKT.
func MarshalFieldMask(s *jsonplugin.MarshalState, v *fieldmaskpb.FieldMask) {
	if v == nil {
		s.WriteNil()
		return
	}
	s.WriteFieldMask(v)
}

func MarshalLegacyFieldMask(s *jsonplugin.MarshalState, v *fieldmaskpb.FieldMask) {
	if v == nil {
		s.WriteNil()
		return
	}
	s.WriteLegacyFieldMask(v)
}

// UnmarshalFieldMask unmarshals a FieldMask WKT.
func UnmarshalFieldMask(s *jsonplugin.UnmarshalState) *fieldmaskpb.FieldMask {
	if s.ReadNil() {
		return nil
	}
	m := s.ReadFieldMask()
	if s.Err() != nil {
		return nil
	}
	return &fieldmaskpb.FieldMask{Paths: m.GetPaths()}
}

// MarshalStruct marshals a Struct WKT.
func MarshalStruct(s *jsonplugin.MarshalState, v *structpb.Struct) {
	if v == nil {
		s.WriteNil()
		return
	}
	MarshalMessage(s, v)
}

// UnmarshalStruct unmarshals a Struct WKT.
func UnmarshalStruct(s *jsonplugin.UnmarshalState) *structpb.Struct {
	if s.ReadNil() {
		return nil
	}
	var v structpb.Struct
	UnmarshalMessage(s, &v)
	if s.Err() != nil {
		return nil
	}
	return &v
}

// MarshalValue marshals a Value WKT.
func MarshalValue(s *jsonplugin.MarshalState, v *structpb.Value) {
	if v == nil {
		s.WriteNil()
		return
	}
	MarshalMessage(s, v)
}

// UnmarshalValue unmarshals a Value WKT.
func UnmarshalValue(s *jsonplugin.UnmarshalState) *structpb.Value {
	if s.ReadNil() {
		return &structpb.Value{Kind: &structpb.Value_NullValue{}}
	}
	var v structpb.Value
	UnmarshalMessage(s, &v)
	if s.Err() != nil {
		return nil
	}
	return &v
}

// MarshalListValue marshals a ListValue WKT.
func MarshalListValue(s *jsonplugin.MarshalState, v *structpb.ListValue) {
	if v == nil {
		s.WriteNil()
		return
	}
	MarshalMessage(s, v)
}

// UnmarshalListValue unmarshals a ListValue WKT.
func UnmarshalListValue(s *jsonplugin.UnmarshalState) *structpb.ListValue {
	if s.ReadNil() {
		return nil
	}
	var v structpb.ListValue
	UnmarshalMessage(s, &v)
	if s.Err() != nil {
		return nil
	}
	return &v
}

// MarshalTimestamp marshals a Timestamp WKT.
func MarshalTimestamp(s *jsonplugin.MarshalState, v *timestamppb.Timestamp) {
	if v == nil {
		s.WriteNil()
		return
	}
	s.WriteTime(v.AsTime())
}

// UnmarshalTimestamp unmarshals a Timestamp WKT.
func UnmarshalTimestamp(s *jsonplugin.UnmarshalState) *timestamppb.Timestamp {
	if s.ReadNil() {
		return nil
	}
	t := s.ReadTime()
	if s.Err() != nil {
		return nil
	}
	return timestamppb.New(*t)
}
