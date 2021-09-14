// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package gen

import (
	"fmt"
	"strconv"

	"github.com/TheThingsIndustries/protoc-gen-go-json/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) enumHasMarshaler(enum *protogen.Enum) bool {
	var generateMarshaler bool

	// If the file has the (thethings.json.file) option, and marshaler_all is set, we start with that.
	fileOpts := enum.Desc.ParentFile().Options().(*descriptorpb.FileOptions)
	if proto.HasExtension(fileOpts, annotations.E_File) {
		if fileExt, ok := proto.GetExtension(fileOpts, annotations.E_File).(*annotations.FileOptions); ok {
			if fileExt.MarshalerAll != nil {
				generateMarshaler = *fileExt.MarshalerAll
			}
		}
	}

	if g.enumHasCustomValues(enum) {
		generateMarshaler = true
	}

	// If the enum has the (thethings.json.enum) option and wants to always marshal as a number/string, we need to generate a marshaler.
	// Finally, the marshaler field can still override to true or false if explicitly set.
	enumOpts := enum.Desc.Options().(*descriptorpb.EnumOptions)
	if proto.HasExtension(enumOpts, annotations.E_Enum) {
		if enumExt, ok := proto.GetExtension(enumOpts, annotations.E_Enum).(*annotations.EnumOptions); ok {
			if enumExt.GetMarshalAsNumber() || enumExt.GetMarshalAsString() {
				generateMarshaler = true
			}
			if enumExt.Marshaler != nil {
				generateMarshaler = *enumExt.Marshaler
			}
		}
	}

	return generateMarshaler
}

func (g *generator) genEnumMarshaler(enum *protogen.Enum) {
	ext, _ := proto.GetExtension(enum.Desc.Options().(*descriptorpb.EnumOptions), annotations.E_Enum).(*annotations.EnumOptions)

	// If the enum has custom values, we create a map[int32]string that maps the number to those values.
	hasCustomValues := g.enumHasCustomValues(enum)
	if hasCustomValues {
		g.P("// ", enum.GoIdent, "_customname contains custom string values that override ", enum.GoIdent, "_name.")
		g.P("var ", enum.GoIdent, "_customname = map[int32]string{")
		for _, value := range enum.Values {
			opts := value.Desc.Options().(*descriptorpb.EnumValueOptions)
			if ext, ok := proto.GetExtension(opts, annotations.E_EnumValue).(*annotations.EnumValueOptions); ok {
				if customValue := ext.GetValue(); customValue != "" {
					g.P(value.Desc.Number(), " : ", strconv.Quote(customValue), ",")
				}
			}
		}
		g.P("}")
	}

	g.P("// MarshalProtoJSON marshals the ", enum.GoIdent, " to JSON.")
	g.P("func (x ", enum.GoIdent, ") MarshalProtoJSON(s *", jsonPluginPackage.Ident("MarshalState"), ") {")
	if ext.GetMarshalAsNumber() {
		if ext.GetMarshalAsString() {
			g.gen.Error(fmt.Errorf("%s has both marshal_as_numer and marshal_as_string", enum.Desc.Name()))
		}
		g.P("s.WriteEnumNumber(int32(x))")
	} else {
		// WriteEnum writes the enum according to the EnumsAsInts config.
		// If we really want strings, we use WriteEnumString.
		fun := "WriteEnum"
		if ext.GetMarshalAsString() {
			fun = "WriteEnumString"
		}
		if hasCustomValues {
			// We write the enum, passing both the original mapping, and our custom mapping to the marshaler.
			g.P("s.", fun, "(int32(x), ", enum.GoIdent, "_customname, ", enum.GoIdent, "_name)")
		} else {
			// We write the enum, passing only the original mapping to the marshaler.
			g.P("s.", fun, "(int32(x), ", enum.GoIdent, "_name)")
		}
	}
	g.P("}")
	g.P()
}
