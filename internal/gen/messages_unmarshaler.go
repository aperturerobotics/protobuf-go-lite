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

func (g *generator) messageHasUnmarshaler(message *protogen.Message, visited ...*protogen.Message) bool {
	// Since we're going to be looking at the fields of this message, it's possible that there will be cycles.
	// If that's the case, we'll return false here so that the caller can continue with the next field.
	for _, visited := range visited {
		if message == visited {
			return false
		}
	}

	// No code is generated for map entries, so we also don't need to generate unmarshalers.
	if message.Desc.IsMapEntry() {
		return false
	}

	var generateUnmarshaler bool

	// If the file has the (thethings.json.file) option, and unmarshaler_all is set, we start with that.
	fileOpts := message.Desc.ParentFile().Options().(*descriptorpb.FileOptions)
	if proto.HasExtension(fileOpts, annotations.E_File) {
		if fileExt, ok := proto.GetExtension(fileOpts, annotations.E_File).(*annotations.FileOptions); ok {
			if fileExt.UnmarshalerAll != nil {
				generateUnmarshaler = *fileExt.UnmarshalerAll
			}
		}
	}

	for _, field := range message.Fields {
		// If the field has the (thethings.json.field) option, and unmarshaler_func is set, we need to generate an unmarshaler for the message.
		fieldOpts := field.Desc.Options().(*descriptorpb.FieldOptions)
		if proto.HasExtension(fieldOpts, annotations.E_Field) {
			if fieldExt, ok := proto.GetExtension(fieldOpts, annotations.E_Field).(*annotations.FieldOptions); ok {
				if fieldExt.UnmarshalerFunc != nil {
					generateUnmarshaler = true
				}
			}
		}

		// If the field is an enum, and the enum has an unmarshaler, we need to generate an unmarshaler for the message.
		if field.Enum != nil && g.enumHasUnmarshaler(field.Enum) {
			generateUnmarshaler = true
		}

		// If the field is a message, and that message has an unmarshaler, we need to generate an unmarshaler.
		if field.Message != nil && g.messageHasUnmarshaler(field.Message, append(visited, message)...) {
			generateUnmarshaler = true
		}
	}

	// If the message has the (thethings.json.message) option and is a wrapper, we need to generate an unmarshaler.
	// Finally, the unmarshaler field can still override to true or false if explicitly set.
	messageOpts := message.Desc.Options().(*descriptorpb.MessageOptions)
	if proto.HasExtension(messageOpts, annotations.E_Message) {
		if messageExt, ok := proto.GetExtension(messageOpts, annotations.E_Message).(*annotations.MessageOptions); ok {
			if messageExt.GetWrapper() {
				generateUnmarshaler = true
			}
			if messageExt.Unmarshaler != nil {
				generateUnmarshaler = *messageExt.Unmarshaler
			}
		}
	}

	return generateUnmarshaler
}

func (g *generator) genMessageUnmarshaler(message *protogen.Message) {
	g.P("// UnmarshalProtoJSON unmarshals the ", message.GoIdent, " message from JSON.")
	g.P("func (x *", message.GoIdent, ") UnmarshalProtoJSON(s *", jsonPluginPackage.Ident("UnmarshalState"), ") {")

	// If we se a null, there's nothing to do.
	g.P("if s.ReadNil() {")
	g.P("return")
	g.P("}")

	// If the message doesn't have any fields, there's nothing to do.
	if len(message.Fields) == 0 {
		g.P("}") // end func (x *{message.GoIdent}) MarshalProtoJSON()
		g.P()
		return
	}

	// If the message is a wrapper type, we operate directly on the first field (named Value) inside it.
	if messageIsWrapper(message) {
		field := message.Fields[0]
		switch field.Desc.Kind() {
		default:
			// Scalar types can be read by the library.
			g.P("x.Value = s.Read", g.libNameForField(field), "()")
		case protoreflect.EnumKind:
			if g.enumHasUnmarshaler(field.Enum) {
				// If the wrapped field is of type enum, and the enum has an unmarshaler, use that.
				g.P(`x.Value.UnmarshalProtoJSON(s)`)
			} else {
				// Otherwise we let the library read the enum.
				g.P("x.Value = ", field.Enum.GoIdent, "(s.ReadEnum(", field.Enum.GoIdent, "_value))")
			}
		}
		g.P("return")
		g.P("}") // end func (x *{message.GoIdent}) MarshalProtoJSON()
		g.P()
		return
	}

	g.P("s.ReadObject(func(key string) {")
	g.P("switch key {")
	g.P("default:")
	g.P("s.ReadAny() // ignore unknown field")

nextField:
	for _, field := range message.Fields {
		var (
			pluginPackage               = golangPluginPackage
			fieldGoName     interface{} = fieldGoName(field)
			nullable                    = fieldIsNullable(field)
			customtype                  = fieldCustomType(field)
			unmarshalerFunc *protogen.GoIdent
		)
		fieldOpts := field.Desc.Options()
		if Params.Lang == "gogo" {
			pluginPackage = gogoPluginPackage
		}
		if proto.HasExtension(fieldOpts, annotations.E_Field) {
			if customtype == nil {
				unmarshalerFunc = parseGoIdent(proto.GetExtension(field.Desc.Options(), annotations.E_Field).(*annotations.FieldOptions).GetUnmarshalerFunc())
			}
		}

		// We need to match both the snake case field name and the camel case JSON name.
		// If those are the same, we only need to match one.
		if string(field.Desc.Name()) != field.Desc.JSONName() {
			g.P(`case "`, field.Desc.Name(), `", "`, field.Desc.JSONName(), `":`)
		} else {
			g.P(`case "`, field.Desc.Name(), `":`)
		}

		// For sub-messages, field mask handling will be handled by the unmarshaler of the sub-message.
		// For scalar types and fields that don't support field masks (lists, maps, fields without unmarshalers) we do field mask handling here.
		delegateMask := "true"
		if field.Message == nil || field.Desc.IsList() || field.Desc.IsMap() || !g.messageHasUnmarshaler(field.Message) || messageIsWKT(field.Message) || messageIsWrapper(field.Message) {
			delegateMask = "false"
			g.P(`s.AddField("`, field.Desc.Name(), `")`)
		}

		if field.Desc.IsMap() {
			// If the field has a custom unmarshaler, call that and continue with the next field.
			if unmarshalerFunc != nil {
				g.P("x.", fieldGoName, " = ", *unmarshalerFunc, `(s.WithField("`, field.Desc.Name(), `", `, delegateMask, `))`)
				continue nextField
			}

			// If the field is a map, the field type is a MapEntry message.
			// In the MapEntry message, the first field is the key, and the second field is the value.
			key := field.Message.Fields[0]
			value := field.Message.Fields[1]

			// Allocate an empty map[T(key)]T(value).
			g.P("x.", fieldGoName, " = make(map[", g.goTypeForField(key), "]", ifThenElse(fieldIsNullable(value), "*", ""), g.goTypeForField(value), ")")

			// Tell the library to read a map with keys of the given type, passing our handler func that will be called for each key.
			g.P("s.Read", g.libNameForField(key), "Map(func(key ", g.goTypeForField(key), ") {")

			switch value.Desc.Kind() {
			default:
				// Scalar types can be read by the library.
				g.P("x.", fieldGoName, "[key] = s.Read", g.libNameForField(value), "()")
			case protoreflect.EnumKind:
				if g.enumHasUnmarshaler(value.Enum) {
					// If the map value is of type enum, and the enum has an unmarshaler,
					// allocate a zero enum, call the unmarshaler, and set the map value for key to the enum.
					g.P("var v ", value.Enum.GoIdent)
					g.P(`v.UnmarshalProtoJSON(s)`)
					g.P("x.", fieldGoName, "[key] = v")
				} else {
					// Otherwise we let the library read the enum.
					g.P("x.", fieldGoName, "[key] = ", value.Enum.GoIdent, "(s.ReadEnum(", value.Enum.GoIdent, "_value))")
				}
			case protoreflect.MessageKind:
				switch {
				case g.messageHasUnmarshaler(value.Message):
					// If the map value is of type message, and the message has a marshaler,
					// allocate a zero message, call the unmarshaler and set the map value for the key to the message.
					g.P("var v ", value.Message.GoIdent)
					g.P(`v.UnmarshalProtoJSON(s)`)
					g.P("x.", fieldGoName, "[key] = &v")
				case messageIsWrapper(value.Message):
					// If the map value is a wrapper, and we read null, set the map value for the key to nil.
					// Otherwise read the wrapped value, and if successful, set the map value for the key to the wrapped value.
					g.P("if s.ReadNil() {")
					g.P("x.", fieldGoName, "[key] = nil")
					g.P("} else {")
					g.P("v := ", g.readWrapperValue(value.Message))
					g.P("if s.Err() != nil {")
					g.P("return")
					g.P("}")
					g.P("x.", fieldGoName, "[key] = &", value.Message.GoIdent, "{Value: v}")
					g.P("}") // end if s.ReadNil() {
				case messageIsWKT(value.Message):
					// If the map value is a WKT, read the WKT.
					g.P("v := ", g.readWKTValue(field, value.Message))
					g.P("if s.Err() != nil {")
					g.P("return")
					g.P("}")
					g.P("x.", fieldGoName, "[key] = ", ifThenElse(nullable, "", "*"), "v")
				default:
					// Otherwise, delegate to the library.
					g.P("// NOTE: ", value.Message.GoIdent.GoName, " does not seem to implement UnmarshalProtoJSON.")
					g.P("var v ", value.Message.GoIdent)
					g.P(pluginPackage.Ident("UnmarshalMessage"), "(s, &v)")
					g.P("x.", fieldGoName, "[key] = &v")
				}
			}

			g.P("})") // end s.Read{key}Map()
			continue nextField
		}

		if field.Desc.IsList() {
			// If the field has a custom unmarshaler, call that and continue with the next field.
			if unmarshalerFunc != nil {
				g.P("x.", fieldGoName, " = ", *unmarshalerFunc, `(s.WithField("`, field.Desc.Name(), `", `, delegateMask, `))`)
				continue nextField
			}
			if customtype != nil {
				// If the field has a custom type, for each element, we will
				// allocate a zero value, call the unmarshaler and append the value to the list.
				g.P("s.ReadArray(func() {")
				g.P("var v ", *customtype)
				g.P(`v.UnmarshalProtoJSON(s.WithField("`, field.Desc.Name(), `", `, delegateMask, `))`)
				g.P("x.", fieldGoName, " = append(x.", fieldGoName, ", v)")
				g.P("})")
				continue nextField
			}
			switch field.Desc.Kind() {
			default:
				// Lists of scalar types can be read by the library.
				g.P("x.", fieldGoName, " = s.Read", g.libNameForField(field), "Array()")
			case protoreflect.EnumKind:
				g.P("s.ReadArray(func() {")
				if g.enumHasUnmarshaler(field.Enum) {
					// If the list value is of type enum, and the enum has an unmarshaler,
					// allocate a zero enum, call the unmarshaler, and append the enum to the list.
					g.P("var v ", field.Enum.GoIdent)
					g.P(`v.UnmarshalProtoJSON(s)`)
					g.P("x.", fieldGoName, " = append(x.", fieldGoName, ", v)")
				} else {
					// Otherwise we let the library read the enum.
					g.P("x.", fieldGoName, " = append(x.", fieldGoName, ", ", field.Enum.GoIdent, "(s.ReadEnum(", field.Enum.GoIdent, "_value)))")
				}
				g.P("})") // end s.ReadArray()
			case protoreflect.MessageKind:
				g.P("s.ReadArray(func() {")
				switch {
				case g.messageHasUnmarshaler(field.Message):
					if nullable {
						// If we read nil, append nil and return so that we can continue with the next key.
						g.P("if s.ReadNil() {")
						g.P("x.", fieldGoName, " = append(x.", fieldGoName, ", nil)")
						g.P("return")
						g.P("}") // end if s.ReadNil() {
					}
					// Allocate a zero message, call the unmarshaler and append the message to the list.
					g.P("v := ", ifThenElse(nullable, "&", ""), field.Message.GoIdent, "{}")
					g.P(`v.UnmarshalProtoJSON(s.WithField("`, field.Desc.Name(), `", `, delegateMask, `))`)
					g.P("if s.Err() != nil {")
					g.P("return")
					g.P("}")
					g.P("x.", fieldGoName, " = append(x.", fieldGoName, ", v)")
				case messageIsWrapper(field.Message):
					if nullable {
						// If we read nil, append nil and return so that we can continue with the next key.
						g.P("if s.ReadNil() {")
						g.P("x.", fieldGoName, " = append(x.", fieldGoName, ", nil)")
						g.P("return")
						g.P("}") // end if s.ReadNil() {
					}
					// Read the wrapped value, and if successful, append the wrapped value to the list.
					g.P("v := ", g.readWrapperValue(field.Message))
					g.P("if s.Err() != nil {")
					g.P("return")
					g.P("}")
					g.P("x.", fieldGoName, " = append(x.", fieldGoName, ", &", field.Message.GoIdent, "{Value: v})")
				case messageIsWKT(field.Message):
					// If the list value is a WKT, read the WKT, and if successful, append it to the list.
					g.P("v := ", g.readWKTValue(field, field.Message))
					g.P("if s.Err() != nil {")
					g.P("return")
					g.P("}")
					g.P("x.", fieldGoName, " = append(x.", fieldGoName, ", ", ifThenElse(nullable, "", "*"), "v)")
				default:
					// Otherwise, delegate to the library.
					g.P("// NOTE: ", field.Message.GoIdent.GoName, " does not seem to implement UnmarshalProtoJSON.")
					g.P("var v ", field.Message.GoIdent)
					g.P(pluginPackage.Ident("UnmarshalMessage"), "(s, &v)")
					g.P("x.", fieldGoName, " = append(x.", fieldGoName, ", ", ifThenElse(nullable, "&", ""), "v)")
				}

				g.P("})") // end s.ReadArray()
			}

			continue nextField
		}

		// The identifier of the message is x, but in case of a oneof, we'll be operating on ov.
		messageOrOneofIdent := "x"

		// If this field is in a oneof, allocate a new oneof value wrapper.
		if field.Oneof != nil {
			g.P("ov := &", field.GoIdent.GoName, "{}")
			messageOrOneofIdent = "ov"
		}

		// If the field has a custom unmarshaler, call that
		if unmarshalerFunc != nil {
			g.P(messageOrOneofIdent, ".", fieldGoName, " = ", *unmarshalerFunc, `(s.WithField("`, field.Desc.Name(), `", `, delegateMask, `))`)
		} else if customtype != nil {
			if nullable {
				g.P("if !s.ReadNil() {")
				// Set the field to a newly allocated custom type.
				g.P(messageOrOneofIdent, ".", fieldGoName, " = &", *customtype, "{}")
				// Call UnmarshalProtoJSON on the field.
				g.P(messageOrOneofIdent, ".", fieldGoName, `.UnmarshalProtoJSON(s.WithField("`, field.Desc.Name(), `", `, delegateMask, `))`)
				g.P("}")
			} else {
				// Call UnmarshalProtoJSON on the field.
				g.P(messageOrOneofIdent, ".", fieldGoName, `.UnmarshalProtoJSON(s.WithField("`, field.Desc.Name(), `", `, delegateMask, `))`)
			}
		} else {
			switch field.Desc.Kind() {
			default:
				// Scalar types can be read by the library.
				g.P(messageOrOneofIdent, ".", fieldGoName, " = s.Read", g.libNameForField(field), "()")
			case protoreflect.EnumKind:
				if g.enumHasUnmarshaler(field.Enum) {
					// If the field is of type enum, and the enum has an unmarshaler, call the unmarshaler.
					g.P(messageOrOneofIdent, ".", fieldGoName, ".UnmarshalProtoJSON(s)")
				} else {
					// Otherwise we let the library read the enum.
					g.P(messageOrOneofIdent, ".", fieldGoName, " = ", field.Enum.GoIdent, "(s.ReadEnum(", field.Enum.GoIdent, "_value))")
				}
			case protoreflect.MessageKind:
				switch {
				case g.messageHasUnmarshaler(field.Message):
					g.P("if !s.ReadNil() {")
					if nullable {
						// Set the field (or enum wrapper) to a newly allocated custom type.
						g.P(messageOrOneofIdent, ".", fieldGoName, " = &", field.Message.GoIdent, "{}")
					}
					// Call UnmarshalProtoJSON on the field.
					g.P(messageOrOneofIdent, ".", fieldGoName, `.UnmarshalProtoJSON(s.WithField("`, field.Desc.Name(), `", `, delegateMask, `))`)
					g.P("}") // end if !s.ReadNil() {
				case messageIsWrapper(field.Message):
					g.P("if !s.ReadNil() {")
					// Read the wrapped value, and if successful, set the wrapped value in the field.
					g.P("v := ", g.readWrapperValue(field.Message))
					g.P("if s.Err() != nil {")
					g.P("return")
					g.P("}")
					g.P(messageOrOneofIdent, ".", fieldGoName, " = &", field.Message.GoIdent, "{Value: v}")
					g.P("}") // end if !s.ReadNil() {
				case messageIsWKT(field.Message):
					// Read the WKT, and if successful, set it in the field.
					g.P("v := ", g.readWKTValue(field, field.Message))
					g.P("if s.Err() != nil {")
					g.P("return")
					g.P("}")
					g.P(messageOrOneofIdent, ".", fieldGoName, " = ", ifThenElse(nullable, "", "*"), "v")
				default:
					// Otherwise, delegate to the library.
					g.P("// NOTE: ", field.Message.GoIdent.GoName, " does not seem to implement UnmarshalProtoJSON.")
					g.P("var v ", field.Message.GoIdent)
					g.P(pluginPackage.Ident("UnmarshalMessage"), "(s, &v)")
					g.P(messageOrOneofIdent, ".", fieldGoName, " = ", ifThenElse(nullable, "&", ""), "v")
				}
			}
		}

		if field.Oneof != nil {
			// Set the message field to the oneof wrapper.
			g.P("x.", field.Oneof.GoName, " = ov")
			continue nextField
		}
	}

	g.P("}")  // end switch key {
	g.P("})") // end s.ReadObject()
	g.P("}")  // end func (x *{message.GoIdent}) MarshalProtoJSON()
	g.P()
}

func (g *generator) readWrapperValue(message *protogen.Message) string {
	switch message.Desc.FullName() {
	case "google.protobuf.DoubleValue":
		return "s.ReadFloat64()"
	case "google.protobuf.FloatValue":
		return "s.ReadFloat32()"
	case "google.protobuf.Int64Value":
		return "s.ReadInt64()"
	case "google.protobuf.UInt64Value":
		return "s.ReadUint64()"
	case "google.protobuf.Int32Value":
		return "s.ReadInt32()"
	case "google.protobuf.UInt32Value":
		return "s.ReadUint32()"
	case "google.protobuf.BoolValue":
		return "s.ReadBool()"
	case "google.protobuf.StringValue":
		return "s.ReadString()"
	case "google.protobuf.BytesValue":
		return "s.ReadBytes()"
	default:
		g.gen.Error(fmt.Errorf("unsupported wrapper %q", message.Desc.FullName()))
		return ""
	}
}

func (g *generator) readWKTValue(field *protogen.Field, message *protogen.Message) string {
	pluginPackage := golangPluginPackage
	if Params.Lang == "gogo" {
		pluginPackage = gogoPluginPackage
	}
	switch message.Desc.FullName() {
	case "google.protobuf.Any":
		return g.QualifiedGoIdent(pluginPackage.Ident("UnmarshalAny")) + "(s)"
	case "google.protobuf.Empty":
		return g.QualifiedGoIdent(pluginPackage.Ident("UnmarshalEmpty")) + "(s)"
	case "google.protobuf.FieldMask":
		return g.QualifiedGoIdent(pluginPackage.Ident("UnmarshalFieldMask")) + "(s)"
	case "google.protobuf.Struct":
		return g.QualifiedGoIdent(pluginPackage.Ident("UnmarshalStruct")) + "(s)"
	case "google.protobuf.Value":
		return g.QualifiedGoIdent(pluginPackage.Ident("UnmarshalValue")) + "(s)"
	case "google.protobuf.ListValue":
		return g.QualifiedGoIdent(pluginPackage.Ident("UnmarshalListValue")) + "(s)"
	case "google.protobuf.Timestamp":
		if Params.Lang == "gogo" && proto.HasExtension(field.Desc.Options(), gogoproto.E_Stdtime) && proto.GetExtension(field.Desc.Options(), gogoproto.E_Stdtime).(bool) {
			return "s.ReadTime()"
		}
		return g.QualifiedGoIdent(pluginPackage.Ident("UnmarshalTimestamp")) + "(s)"
	case "google.protobuf.Duration":
		if Params.Lang == "gogo" && proto.HasExtension(field.Desc.Options(), gogoproto.E_Stdduration) && proto.GetExtension(field.Desc.Options(), gogoproto.E_Stdduration).(bool) {
			return "s.ReadDuration()"
		}
		return g.QualifiedGoIdent(pluginPackage.Ident("UnmarshalDuration")) + "(s)"
	default:
		g.gen.Error(fmt.Errorf("unsupported WKT %q", message.Desc.FullName()))
		return ""
	}
}
