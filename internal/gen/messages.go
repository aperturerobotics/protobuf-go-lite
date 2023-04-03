// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package gen

import (
	"fmt"
	"strings"

	"github.com/TheThingsIndustries/protoc-gen-go-json/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) messageHasFieldMask(message *protogen.Message, visited ...*protogen.Message) bool {
	// Since we're going to be looking at the fields of this message, it's possible that there will be cycles.
	// If that's the case, we'll return false here so that the caller can continue with the next field.
	for _, visited := range visited {
		if message == visited {
			return false
		}
	}

	// No code is generated for map entries, so we also don't need to generate marshalers.
	if message.Desc.IsMapEntry() {
		return false
	}

	var generate bool

	for _, field := range message.Fields {
		if strings.HasPrefix(string(field.Desc.FullName()), "google.protobuf.FieldMask") {
			generate = true
		}

		// If the field is a message, and that message has a field mask, we need to generate a marshaler.
		if field.Message != nil && g.messageHasFieldMask(field.Message, append(visited, message)...) {
			generate = true
		}
	}
	return generate
}

func (g *generator) messageHasAnyMarshaler(message *protogen.Message) bool {
	// We have a marshaler if the message itself has a marshaler or unmarshaler.
	if g.messageHasMarshaler(message) || g.messageHasUnmarshaler(message) || g.messageHasFieldMask(message) {
		return true
	}

	// We have a marshaler if any of the enums defined in the message has a marshaler or unmarshaler.
	for _, enum := range message.Enums {
		if g.enumHasAnyMarshaler(enum) {
			return true
		}
	}

	// We have a marshaler if any of the sub-messages defined in the message has any marshaler.
	for _, message := range message.Messages {
		if g.messageHasAnyMarshaler(message) || g.messageHasFieldMask(message) {
			return true
		}
	}

	return false
}

func (g *generator) genMessage(message *protogen.Message) {
	// Generate marshalers and unmarshalers for all enums defined in the message.
	for _, enum := range message.Enums {
		g.genEnum(enum)
	}

	// Generate marshalers and unmarshalers for all sub-messages defined in the message.
	for _, message := range message.Messages {
		g.genMessage(message)
	}

	// Generate marshaler for the message itself, if it has one.
	if g.messageHasMarshaler(message) || g.messageHasFieldMask(message) {
		g.genMessageMarshaler(message)
		if Params.Std {
			g.genStdMessageMarshaler(message)
		}
	}

	// Generate unmarshaler for the message itself, if it has one.
	if g.messageHasUnmarshaler(message) || g.messageHasFieldMask(message) {
		g.genMessageUnmarshaler(message)
		if Params.Std {
			g.genStdMessageUnmarshaler(message)
		}
	}
}

func fieldIsNullable(field *protogen.Field) bool {
	// Typically, only message fields are nullable (use pointers).
	nullable := field.Desc.Kind() == protoreflect.MessageKind
	return nullable
}

func fieldGoName(field *protogen.Field) interface{} {
	var fieldGoName interface{} = field.GoName
	return fieldGoName
}

func fieldCustomType(field *protogen.Field) *protogen.GoIdent {
	return nil
}

func ifThenElse(condition bool, ifTrue, ifFalse string) string {
	if condition {
		return ifTrue
	}
	return ifFalse
}

// goTypeForField returns the name of the Go type that corresponds to the type of a given field.
func (g *generator) goTypeForField(field *protogen.Field) interface{} {
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.EnumKind:
		return field.Enum.GoIdent
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64"
	case protoreflect.FloatKind:
		return "float32"
	case protoreflect.DoubleKind:
		return "float64"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "[]byte"
	case protoreflect.MessageKind:
		return field.Message.GoIdent
	default:
		g.gen.Error(fmt.Errorf("unsupported field kind %q", field.Desc.Kind()))
		return ""
	}
}

// libNameForField returns the name used in the protojson func that corresponds to the type of a given field.
func (g *generator) libNameForField(field *protogen.Field) string {
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return "Bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "Int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "Uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "Int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "Uint64"
	case protoreflect.FloatKind:
		return "Float32"
	case protoreflect.DoubleKind:
		return "Float64"
	case protoreflect.StringKind:
		return "String"
	case protoreflect.BytesKind:
		return "Bytes"
	default:
		g.gen.Error(fmt.Errorf("unsupported field kind %q", field.Desc.Kind()))
		return ""
	}
}

// parseGoIdent parses a custom type and returns a GoIdent for it.
// If it's unable to parse the custom type, it returns nil.
func parseGoIdent(customtype string) *protogen.GoIdent {
	if customtype == "" {
		return nil
	}
	i := strings.LastIndex(customtype, ".")
	ident := protogen.GoImportPath(customtype[:i]).Ident(customtype[i+1:])
	return &ident
}

// messageIsWrapper returns true if the given message is a wrapper type.
// This is the case for well known wrapper types (google.protobuf.XXXValue)
// and for messages that have the (thethings.json.message) option with wrapper = true.
func messageIsWrapper(message *protogen.Message) bool {
	switch message.Desc.FullName() {
	case "google.protobuf.DoubleValue",
		"google.protobuf.FloatValue",
		"google.protobuf.Int64Value",
		"google.protobuf.UInt64Value",
		"google.protobuf.Int32Value",
		"google.protobuf.UInt32Value",
		"google.protobuf.BoolValue",
		"google.protobuf.StringValue",
		"google.protobuf.BytesValue":
		return true
	}
	opts := message.Desc.Options().(*descriptorpb.MessageOptions)
	if ext, hasExt := proto.GetExtension(opts, annotations.E_Message).(*annotations.MessageOptions); hasExt {
		return ext.GetWrapper() && len(message.Fields) == 1 && message.Fields[0].GoName == "Value"
	}
	return false
}

// messageIsWKT returns true if the given message is a well-known type.
func messageIsWKT(message *protogen.Message) bool {
	return strings.HasPrefix(string(message.Desc.FullName()), "google.protobuf.")
}
