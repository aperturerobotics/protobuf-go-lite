// Copyright © 2024 Aperture Robotics, LLC.
// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"fmt"
	"slices"

	"github.com/aperturerobotics/protobuf-go-lite/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (g *jsonGenerator) genMessage(message *protogen.Message) {
	// Generate marshalers and unmarshalers for all enums defined in the message.
	for _, enum := range message.Enums {
		g.genEnum(enum)
	}

	// Generate marshalers and unmarshalers for all sub-messages defined in the message.
	for _, message := range message.Messages {
		g.genMessage(message)
	}

	// skip early if the disable comment is present
	if hasDisableJsonComment(message.Comments.Leading) {
		return
	}

	// Check if the message has any optional fields and skip generation if so.
	anyOptional := slices.ContainsFunc(message.Fields, func(f *protogen.Field) bool {
		return f.Desc.HasOptionalKeyword()
	})

	if !anyOptional {
		g.genMessageMarshaler(message)
		g.genStdMessageMarshaler(message)

		g.genMessageUnmarshaler(message)
		g.genStdMessageUnmarshaler(message)
	} else {
		// We do not support marshaling this field, skip the entire message.
		g.P("// NOTE: protobuf-go-lite json only supports proto3 and not proto3opt (optional fields).")
		g.P()
	}
}

func fieldIsNullable(field *protogen.Field) bool {
	// In the supported subset of syntax (proto3 and not proto3opt) we only use pointers for messages.
	nullable := field.Desc.Kind() == protoreflect.MessageKind
	return nullable
}

func fieldGoName(field *protogen.Field) interface{} {
	var fieldGoName interface{} = field.GoName
	return fieldGoName
}

// libNameForField returns the name used in the protojson func that corresponds to the type of a given field.
func (g *jsonGenerator) libNameForField(field *protogen.Field) string {
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
