// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package gen

import (
	"fmt"

	"github.com/TheThingsIndustries/protoc-gen-go-json/annotations"
	"github.com/TheThingsIndustries/protoc-gen-go-json/internal/gogoproto"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) messageHasMarshaler(message *protogen.Message, visited ...*protogen.Message) bool {
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

	var generateMarshaler bool

	// If the file has the (thethings.json.file) option, and marshaler_all is set, we start with that.
	fileOpts := message.Desc.ParentFile().Options().(*descriptorpb.FileOptions)
	if proto.HasExtension(fileOpts, annotations.E_File) {
		if fileExt, ok := proto.GetExtension(fileOpts, annotations.E_File).(*annotations.FileOptions); ok {
			if fileExt.MarshalerAll != nil {
				generateMarshaler = *fileExt.MarshalerAll
			}
		}
	}

	for _, field := range message.Fields {
		// If the field has the (thethings.json.field) option, and marshaler_func is set, we need to generate a marshaler for the message.
		fieldOpts := field.Desc.Options().(*descriptorpb.FieldOptions)
		if proto.HasExtension(fieldOpts, annotations.E_Field) {
			if fieldExt, ok := proto.GetExtension(fieldOpts, annotations.E_Field).(*annotations.FieldOptions); ok {
				if fieldExt.MarshalerFunc != nil {
					generateMarshaler = true
				}
			}
		}

		// If the field is an enum, and the enum has a marshaler, we need to generate a marshaler for the message.
		if field.Enum != nil && g.enumHasMarshaler(field.Enum) {
			generateMarshaler = true
		}

		// If the field is a message, and that message has a marshaler, we need to generate a marshaler.
		if field.Message != nil && g.messageHasMarshaler(field.Message, append(visited, message)...) {
			generateMarshaler = true
		}
	}

	// If the message has the (thethings.json.message) option and is a wrapper, we need to generate a marshaler.
	// Finally, the marshaler field can still override to true or false if explicitly set.
	messageOpts := message.Desc.Options().(*descriptorpb.MessageOptions)
	if proto.HasExtension(messageOpts, annotations.E_Message) {
		if messageExt, ok := proto.GetExtension(messageOpts, annotations.E_Message).(*annotations.MessageOptions); ok {
			if messageExt.GetWrapper() {
				generateMarshaler = true
			}
			if messageExt.Marshaler != nil {
				generateMarshaler = *messageExt.Marshaler
			}
		}
	}

	return generateMarshaler
}

func (g *generator) genMessageMarshaler(message *protogen.Message) {
	g.P("// MarshalProtoJSON marshals the ", message.GoIdent, " message to JSON.")
	g.P("func (x *", message.GoIdent, ") MarshalProtoJSON(s *", jsonPluginPackage.Ident("MarshalState"), ") {")

	g.P("if x == nil {")
	g.P("s.WriteNil()")
	g.P("return")
	g.P("}")

	// If the message is a wrapper type, we operate directly on the first field (named Value) inside it.
	if messageIsWrapper(message) {
		field := message.Fields[0]
		switch field.Desc.Kind() {
		default:
			// Scalar types can be written by the library.
			g.P("s.Write", g.libNameForField(field), "(x.Value)")
		case protoreflect.EnumKind:
			if g.enumHasMarshaler(field.Enum) {
				// If the wrapped field is of type enum, and the enum has a marshaler, use that.
				g.P(`x.Value.MarshalProtoJSON(s)`)
			} else {
				// Otherwise we let the library write the enum.
				g.P("s.WriteEnum(int32(x.Value), ", field.Enum.GoIdent, "_name)")
			}
		}
		g.P("return")
		g.P("}") // end func (x *{message.GoIdent}) MarshalProtoJSON()
		g.P()
		return
	}

	g.P("s.WriteObjectStart()")

	// If the message doesn't have any fields, there's nothing to do.
	if len(message.Fields) == 0 {
		g.P("s.WriteObjectEnd()")
		g.P("}") // end func (x *{message.GoIdent}) MarshalProtoJSON()
		g.P()
		return
	}

	// wroteField keeps track of whether we wrote a field, so that we know when to add a comma before the next.
	g.P("var wroteField bool")

nextField:
	for _, field := range message.Fields {
		var (
			pluginPackage             = golangPluginPackage
			fieldGoName   interface{} = fieldGoName(field)
			nullable                  = fieldIsNullable(field)
			customtype                = fieldCustomType(field)
			marshalerFunc *protogen.GoIdent
		)
		fieldOpts := field.Desc.Options()
		if Params.Lang == "gogo" {
			pluginPackage = gogoPluginPackage
		} else {
			if proto.HasExtension(fieldOpts, annotations.E_Field) {
				marshalerFunc = parseGoIdent(proto.GetExtension(field.Desc.Options(), annotations.E_Field).(*annotations.FieldOptions).GetMarshalerFunc())
			}
		}

		if field.Desc.IsMap() {
			// If the field is a map, the field type is a MapEntry message.
			// In the MapEntry message, the first field is the key, and the second field is the value.
			key := field.Message.Fields[0]
			value := field.Message.Fields[1]

			// We emit the field if the map is not nil or if it's specified in the field mask.
			g.P("if x.", fieldGoName, ` != nil || s.HasField("`, field.Desc.Name(), `") {`)

			// Write a comma if this isn't the first field.
			g.P("s.WriteMoreIf(&wroteField)")

			// Write the field name and a colon.
			g.P(`s.WriteObjectField("`, field.Desc.Name(), `")`)

			// If the map value has a custom marshaler, call that and continue with the next field.
			if marshalerFunc != nil {
				g.P(*marshalerFunc, `(s.WithField("`, field.Desc.Name(), `"), x.`, fieldGoName, ")")
				g.P("}") // end if x.{fieldGoName} != nil {
				continue nextField
			}

			g.P("s.WriteObjectStart()")

			// wroteElement keeps track of whether we wrote an element of the map, so that we know when to add a comma before the next.
			g.P("var wroteElement bool")

			g.P("for k, v := range x.", fieldGoName, " {")

			// Write a comma if this isn't the first element of the map.
			g.P("s.WriteMoreIf(&wroteElement)")

			// Write the key and a a colon. Since they key can be of other types than string, we use the library to convert those.
			g.P("s.WriteObject", g.libNameForField(key), "Field(k)")

			switch value.Desc.Kind() {
			default:
				// Scalar types can be written by the library.
				g.P("s.Write", g.libNameForField(value), "(v)")
			case protoreflect.EnumKind:
				if g.enumHasMarshaler(value.Enum) {
					// If the map value is of type enum, and the enum has a marshaler, use that.
					g.P("v.MarshalProtoJSON(s)")
				} else {
					// Otherwise we write the enum with the standard settings.
					g.P("s.WriteEnum(int32(v), ", value.Enum.GoIdent, "_name)")
				}
			case protoreflect.MessageKind:
				switch {
				case g.messageHasMarshaler(value.Message):
					// If the map value is of type message, and the message has a marshaler, use that.
					g.P(`v.MarshalProtoJSON(s.WithField("`, field.Desc.Name(), `"))`)
				case messageIsWrapper(value.Message):
					// If the map value is a wrapper, write the wrapped value.
					g.writeWrapperValue(value.Message, "v")
				case messageIsWKT(value.Message):
					// If the map value is a WKT, write the WKT.
					g.writeWKTValue(field, value.Message, "v")
				default:
					// Otherwise delegate to the library.
					g.P("// NOTE: ", value.Message.GoIdent.GoName, " does not seem to implement MarshalProtoJSON.")
					g.P(pluginPackage.Ident("MarshalMessage"), "(s, v)")
				}
			}

			g.P("}") // end for k, v := range x.{fieldGoName} {
			g.P("s.WriteObjectEnd()")
			g.P("}") // end if x.{fieldGoName} != nil {

			continue nextField
		}

		if field.Desc.IsList() {
			// We emit the field if the list is not empty or if it's specified in the field mask.
			g.P("if len(x.", fieldGoName, `) > 0 || s.HasField("`, field.Desc.Name(), `") {`)

			// Write a comma if this isn't the first field.
			g.P("s.WriteMoreIf(&wroteField)")

			// Write the field name and a colon.
			g.P(`s.WriteObjectField("`, field.Desc.Name(), `")`)

			// If the field has a custom marshaler, call that and continue with the next field.
			if marshalerFunc != nil {
				g.P(*marshalerFunc, `(s.WithField("`, field.Desc.Name(), `"), x.`, fieldGoName, ")")
				g.P("}") // end if len(x.{fieldGoName}) > 0 {

				continue nextField
			}

			// If the field has a custom type, call MarshalProtoJSON for each element.
			if customtype != nil {
				g.P("s.WriteArrayStart()")

				// wroteElement keeps track of whether we wrote an element of the list, so that we know when to add a comma before the next.
				g.P("var wroteElement bool")

				g.P("for _, element := range x.", fieldGoName, " {")

				// Write a comma if this isn't the first element of the list.
				g.P("s.WriteMoreIf(&wroteElement)")

				g.P(`element.MarshalProtoJSON(s.WithField("`, field.Desc.Name(), `"))`)

				g.P("}") // end for _, element := range x.{fieldGoName} {
				g.P("s.WriteArrayEnd()")
				g.P("}") // end if len(x.{fieldGoName}) > 0 {

				continue nextField
			}

			switch field.Desc.Kind() {
			default:
				g.P("s.Write", g.libNameForField(field), "Array(x.", fieldGoName, ")")
			case protoreflect.EnumKind:
				g.P("s.WriteArrayStart()")

				// wroteElement keeps track of whether we wrote an element of the list, so that we know when to add a comma before the next.
				g.P("var wroteElement bool")

				g.P("for _, element := range x.", fieldGoName, " {")

				// Write a comma if this isn't the first element of the list.
				g.P("s.WriteMoreIf(&wroteElement)")

				if g.enumHasMarshaler(field.Enum) {
					// If the field is of type enum, and the enum has a marshaler, use that.
					g.P("element.MarshalProtoJSON(s)")
				} else {
					// Otherwise we write the enum with the standard settings.
					g.P("s.WriteEnum(int32(element), ", field.Enum.GoIdent, "_name)")
				}

				g.P("}") // end for _, element := range x.{fieldGoName} {
				g.P("s.WriteArrayEnd()")
			case protoreflect.MessageKind:
				g.P("s.WriteArrayStart()")

				// wroteElement keeps track of whether we wrote an element of the list, so that we know when to add a comma before the next.
				g.P("var wroteElement bool")

				g.P("for _, element := range x.", fieldGoName, " {")

				// Write a comma if this isn't the first element of the list.
				g.P("s.WriteMoreIf(&wroteElement)")

				switch {
				case g.messageHasMarshaler(field.Message):
					// If the list element is of type message, and the message has a marshaler, use that.
					g.P(`element.MarshalProtoJSON(s.WithField("`, field.Desc.Name(), `"))`)
				case messageIsWrapper(field.Message):
					// If the list element is a wrapper, write the wrapped value.
					g.writeWrapperValue(field.Message, "element")
				case messageIsWKT(field.Message):
					// If the list element is a WKT, write the WKT.
					g.writeWKTValue(field, field.Message, "element")
				default:
					// Otherwise delegate to the library.
					g.P("// NOTE: ", field.Message.GoIdent.GoName, " does not seem to implement MarshalProtoJSON.")
					g.P(pluginPackage.Ident("MarshalMessage"), "(s, ", ifThenElse(nullable, "", "&"), "element)")
				}

				g.P("}") // end for _, element := range x.{fieldGoName} {
				g.P("s.WriteArrayEnd()")
			}

			g.P("}") // end if len(x.{fieldGoName}) > 0 {

			continue nextField
		}

		// The identifier of the message is x, but in case of a oneof, we'll be operating on ov.
		messageOrOneofIdent := "x"

		// If this is the first field in a oneof, write the if statement that checks for nil
		// and start the switch statement for the oneof type.
		if field.Oneof != nil && field == field.Oneof.Fields[0] {
			// NOTE: we don't support field masks here (yet).
			g.P("if x.", field.Oneof.GoName, " != nil {")
			g.P("switch ov := x.", field.Oneof.GoName, ".(type) {")
		}

		if field.Oneof != nil {
			// If we're in a oneof, check if this is the field that's set in the oneof.
			g.P("case *", field.GoIdent.GoName, ":")
			messageOrOneofIdent = "ov"
		} else {
			// If we're not in a oneof, start "if not zero value".
			if nullable {
				// If this field is nullable, we emit it if it's not nil or if it's specified in the field mask.
				g.P("if ", messageOrOneofIdent, ".", fieldGoName, ` != nil || s.HasField("`, field.Desc.Name(), `") {`)
			} else {
				// If this field is not nullable, we emit it if it's not the zero value or if it's specified in the field mask.
				switch field.Desc.Kind() {
				case protoreflect.BoolKind:
					g.P("if ", messageOrOneofIdent, ".", fieldGoName, ` || s.HasField("`, field.Desc.Name(), `") {`)
				case protoreflect.EnumKind:
					g.P("if ", messageOrOneofIdent, ".", fieldGoName, ` != 0 || s.HasField("`, field.Desc.Name(), `") {`)
				case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
					protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
					protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
					protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
					protoreflect.FloatKind,
					protoreflect.DoubleKind:
					g.P("if ", messageOrOneofIdent, ".", fieldGoName, ` != 0 || s.HasField("`, field.Desc.Name(), `") {`)
				case protoreflect.StringKind:
					g.P("if ", messageOrOneofIdent, ".", fieldGoName, ` != "" || s.HasField("`, field.Desc.Name(), `") {`)
				case protoreflect.BytesKind:
					g.P("if len(", messageOrOneofIdent, ".", fieldGoName, `) > 0 || s.HasField("`, field.Desc.Name(), `") {`)
				case protoreflect.MessageKind:
					// For not-nullable messages we have a dummy check.
					g.P("if true { // (gogoproto.nullable) = false")
				}
			}
		}

		// Write a comma if this isn't the first field.
		g.P("s.WriteMoreIf(&wroteField)")

		// Write the field name and a colon.
		g.P(`s.WriteObjectField("`, field.Desc.Name(), `")`)

		if marshalerFunc != nil {
			// If the field has a custom marshaler, call that.
			g.P(*marshalerFunc, `(s.WithField("`, field.Desc.Name(), `"), x.`, fieldGoName, ")")
		} else if customtype != nil {
			// If the field has a custom type, call MarshalProtoJSON for it.
			g.P(messageOrOneofIdent, ".", fieldGoName, `.MarshalProtoJSON(s.WithField("`, field.Desc.Name(), `"))`)
		} else {
			switch field.Desc.Kind() {
			default:
				// Scalar types can be written by the library.
				g.P("s.Write", g.libNameForField(field), "(", messageOrOneofIdent, ".", fieldGoName, ")")
			case protoreflect.EnumKind:
				if g.enumHasMarshaler(field.Enum) {
					// If the field is of type enum, and the enum has a marshaler, use that.
					g.P(messageOrOneofIdent, ".", fieldGoName, ".MarshalProtoJSON(s)")
				} else {
					// Otherwise we write the enum with the standard settings.
					g.P("s.WriteEnum(int32(", messageOrOneofIdent, ".", fieldGoName, "), ", field.Enum.GoIdent, "_name)")
				}
			case protoreflect.MessageKind:
				switch {
				case g.messageHasMarshaler(field.Message):
					// If the field is of type message, and the message has a marshaler, use that.
					g.P(messageOrOneofIdent, ".", fieldGoName, `.MarshalProtoJSON(s.WithField("`, field.Desc.Name(), `"))`)
				case messageIsWrapper(field.Message):
					// If the field is a wrapper, write the wrapped value.
					g.writeWrapperValue(field.Message, fmt.Sprintf("%s.%s", messageOrOneofIdent, fieldGoName))
				case messageIsWKT(field.Message):
					// If the field is a WKT, write the WKT.
					g.writeWKTValue(field, field.Message, fmt.Sprintf("%s.%s", messageOrOneofIdent, fieldGoName))
				default:
					// Otherwise delegate to the library.
					g.P("// NOTE: ", field.Message.GoIdent.GoName, " does not seem to implement MarshalProtoJSON.")
					g.P(pluginPackage.Ident("MarshalMessage"), "(s, ", ifThenElse(nullable, "", "&"), messageOrOneofIdent, ".", fieldGoName, ")")
				}
			}
		}

		// If we're not in a oneof, end the "if not zero".
		if field.Oneof == nil {
			g.P("}") // end if x.{field.GoName} != zero value {
		}

		// If this is the last field in the oneof, close the switch and if statements.
		if field.Oneof != nil && field == field.Oneof.Fields[len(field.Oneof.Fields)-1] {
			g.P("}") // end switch v := x.{field.Oneof.GoName}.(type) {
			g.P("}") // end if x.{field.Oneof.GoName} != nil {
		}
	}

	g.P("s.WriteObjectEnd()")

	g.P("}") // end func (x *{message.GoIdent}) MarshalProtoJSON()
	g.P()
}

func (g *generator) writeWrapperValue(message *protogen.Message, ident string) {
	g.P("if ", ident, " == nil {")
	g.P("s.WriteNil()")
	g.P("} else {")
	// Wrapper values are scalar types, and scalar types can be written by the library.
	g.P("s.Write", g.libNameForField(message.Fields[0]), "(", ident, ".Value)")
	g.P("}") // end if {ident} == nil {
}

func (g *generator) writeWKTValue(field *protogen.Field, message *protogen.Message, ident string) {
	nullable := fieldIsNullable(field)
	if nullable {
		g.P("if ", ident, " == nil {")
		g.P("s.WriteNil()")
		g.P("} else {")
	}
	pluginPackage := golangPluginPackage
	if Params.Lang == "gogo" {
		pluginPackage = gogoPluginPackage
	}
	switch message.Desc.FullName() {
	case "google.protobuf.Any":
		g.P(pluginPackage.Ident("MarshalAny"), "(s, ", ident, ")")
	case "google.protobuf.Empty":
		g.P(pluginPackage.Ident("MarshalEmpty"), "(s, ", ident, ")")
	case "google.protobuf.FieldMask":
		g.P(pluginPackage.Ident("MarshalFieldMask"), "(s, ", ident, ")")
	case "google.protobuf.Struct":
		g.P(pluginPackage.Ident("MarshalStruct"), "(s, ", ident, ")")
	case "google.protobuf.Value":
		g.P(pluginPackage.Ident("MarshalValue"), "(s, ", ident, ")")
	case "google.protobuf.ListValue":
		g.P(pluginPackage.Ident("MarshalListValue"), "(s, ", ident, ")")
	case "google.protobuf.Timestamp":
		if Params.Lang == "gogo" && proto.HasExtension(field.Desc.Options(), gogoproto.E_Stdtime) && proto.GetExtension(field.Desc.Options(), gogoproto.E_Stdtime).(bool) {
			// If the file has the (gogoproto.stdtime) option, marshal the Go time directly.
			g.P("s.WriteTime(", ifThenElse(nullable, "*", ""), ident, ")")
		} else {
			// Otherwise delegate to the library.
			g.P(pluginPackage.Ident("MarshalTimestamp"), "(s, ", ident, ")")
		}
	case "google.protobuf.Duration":
		if Params.Lang == "gogo" && proto.HasExtension(field.Desc.Options(), gogoproto.E_Stdduration) && proto.GetExtension(field.Desc.Options(), gogoproto.E_Stdduration).(bool) {
			// If the file has the (gogoproto.stdduration) option, marshal the Go duration directly.
			g.P("s.WriteDuration(", ifThenElse(nullable, "*", ""), ident, ")")
		} else {
			// Otherwise delegate to the library.
			g.P(pluginPackage.Ident("MarshalDuration"), "(s, ", ident, ")")
		}
	default:
		g.gen.Error(fmt.Errorf("unsupported WKT %q", message.Desc.FullName()))
	}
	if nullable {
		g.P("}") // end if {ident} == nil {
	}
}
