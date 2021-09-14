// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package gen

import (
	"github.com/TheThingsIndustries/protoc-gen-go-json/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) enumHasAnyMarshaler(enum *protogen.Enum) bool {
	return g.enumHasMarshaler(enum) || g.enumHasUnmarshaler(enum)
}

func (g *generator) genEnum(enum *protogen.Enum) {
	if g.enumHasMarshaler(enum) {
		g.genEnumMarshaler(enum)
	}

	if g.enumHasUnmarshaler(enum) {
		g.genEnumUnmarshaler(enum)
	}
}

func (*generator) enumHasCustomValues(enum *protogen.Enum) bool {
	for _, value := range enum.Values {
		// If the file has the (thethings.json.enum_value) option, and a value is set we have custom values.
		opts := value.Desc.Options().(*descriptorpb.EnumValueOptions)
		if ext, ok := proto.GetExtension(opts, annotations.E_EnumValue).(*annotations.EnumValueOptions); ok && ext.GetValue() != "" {
			return true
		}
	}
	return false
}

func (*generator) enumHasCustomAliases(enum *protogen.Enum) bool {
	for _, value := range enum.Values {
		// If the file has the (thethings.json.enum_value) option, and a value is set, or aliases are set, we have custom aliases.
		opts := value.Desc.Options().(*descriptorpb.EnumValueOptions)
		if ext, ok := proto.GetExtension(opts, annotations.E_EnumValue).(*annotations.EnumValueOptions); ok {
			if ext.GetValue() != "" {
				return true
			}
			if len(ext.GetAliases()) > 0 {
				return true
			}
		}
	}
	return false
}

func (*generator) enumValueAliases(value *protogen.EnumValue) []string {
	opts := value.Desc.Options().(*descriptorpb.EnumValueOptions)
	ext, ok := proto.GetExtension(opts, annotations.E_EnumValue).(*annotations.EnumValueOptions)
	if !ok {
		return nil
	}
	aliases := ext.GetAliases()
	// If the enum value has a custom value, add it to aliases if it's not already there.
	if value := ext.GetValue(); value != "" {
		var found bool
		for _, alias := range aliases {
			if alias == value {
				found = true
				break
			}
		}
		if !found {
			return append([]string{value}, aliases...)
		}
	}
	return aliases
}
