// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package gen

import (
	"strconv"
	"strings"

	"github.com/TheThingsIndustries/protoc-gen-go-json/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) enumHasUnmarshaler(enum *protogen.Enum) bool {
	var generateUnmarshaler bool

	// If the file has the (thethings.json.file) option, and unmarshaler_all is set, we start with that.
	fileOpts := enum.Desc.ParentFile().Options().(*descriptorpb.FileOptions)
	if proto.HasExtension(fileOpts, annotations.E_File) {
		if fileExt, ok := proto.GetExtension(fileOpts, annotations.E_File).(*annotations.FileOptions); ok {
			if fileExt.UnmarshalerAll != nil {
				generateUnmarshaler = *fileExt.UnmarshalerAll
			}
		}
	}

	// NOTE: enumHasCustomAliases already checks if there are custom values, but just in case, we also check it.
	if g.enumHasCustomValues(enum) || g.enumHasCustomAliases(enum) {
		generateUnmarshaler = true
	}

	// If the enum has the (thethings.json.enum) option and wants to always marshal as a number/string or has a prefix, we need to generate an unmarshaler.
	// Finally, the unmarshaler field can still override to true or false if explicitly set.
	enumOpts := enum.Desc.Options().(*descriptorpb.EnumOptions)
	if proto.HasExtension(enumOpts, annotations.E_Enum) {
		if enumExt, ok := proto.GetExtension(enumOpts, annotations.E_Enum).(*annotations.EnumOptions); ok {
			if enumExt.GetMarshalAsNumber() || enumExt.GetMarshalAsString() || enumExt.GetPrefix() != "" {
				generateUnmarshaler = true
			}
			if enumExt.Unmarshaler != nil {
				generateUnmarshaler = *enumExt.Unmarshaler
			}
		}
	}
	return generateUnmarshaler
}

func (g *generator) genEnumUnmarshaler(enum *protogen.Enum) {
	ext, _ := proto.GetExtension(enum.Desc.Options().(*descriptorpb.EnumOptions), annotations.E_Enum).(*annotations.EnumOptions)

	// If the enum has a prefix, or any aliases, we create a map[string]int32 that maps the non-prefixed values and aliases to the number.
	prefix := strings.TrimSuffix(ext.GetPrefix(), "_") + "_"
	hasCustomAliases := g.enumHasCustomAliases(enum) || prefix != "_"
	if hasCustomAliases {
		g.P("// ", enum.GoIdent, "_customvalue contains custom string values that extend ", enum.GoIdent, "_value.")
		g.P("var ", enum.GoIdent, "_customvalue = map[string]int32{")
		for _, value := range enum.Values {
			if prefix != "_" {
				if strings.HasPrefix(string(value.Desc.Name()), prefix) {
					g.P(strconv.Quote(strings.TrimPrefix(string(value.Desc.Name()), prefix)), ": ", value.Desc.Number(), ",")
				}
			}
			for _, customValue := range g.enumValueAliases(value) {
				g.P(strconv.Quote(customValue), ": ", value.Desc.Number(), ",")
			}
		}
		g.P("}")
	}

	g.P("// UnmarshalProtoJSON unmarshals the ", enum.GoIdent, " from JSON.")
	g.P("func (x *", enum.GoIdent, ") UnmarshalProtoJSON(s *", jsonPluginPackage.Ident("UnmarshalState"), ") {")
	if hasCustomAliases {
		// We read the enum, passing both the original mapping, and our custom mapping to the unmarshaler.
		g.P("v := s.ReadEnum(", enum.GoIdent, "_value, ", enum.GoIdent, "_customvalue)")
	} else {
		// We read the enum, passing only the original mapping to the unmarshaler.
		g.P("v := s.ReadEnum(", enum.GoIdent, "_value)")
	}
	g.P("if err := s.Err(); err != nil {")
	g.P(`s.SetErrorf("could not read `, enum.Desc.Name(), ` enum: %v", err)`)
	g.P("return")
	g.P("}")
	g.P("*x = ", enum.GoIdent, "(v)")
	g.P("}")
	g.P()
}
