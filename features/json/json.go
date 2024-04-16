// Copyright (c) 2024 Aperture Robotics, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package json

import (
	"github.com/aperturerobotics/protobuf-go-lite/compiler/protogen"
	"github.com/aperturerobotics/protobuf-go-lite/generator"
	"github.com/aperturerobotics/protobuf-go-lite/internal/strs"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	fastjsonPackage = protogen.GoImportPath("github.com/valyala/fastjson")
	gabsPackage     = protogen.GoImportPath("github.com/Jeffail/gabs/v2")
)

func init() {
	generator.RegisterFeature("json", func(gen *generator.GeneratedFile) generator.FeatureGenerator {
		return &jsonMarshal{GeneratedFile: gen}
	})
}

type jsonMarshal struct {
	*generator.GeneratedFile
}

func (p *jsonMarshal) GenerateFile(file *protogen.File) bool {
	if file.Proto.GetSyntax() != "proto3" {
		p.P("// MarshalJSON generator only supports proto3 files.")
		p.P()
		p.P("// UnmarshalJSON generator only supports proto3 files.")
		p.P()
		p.P("// UnmarshalJSONValue generator only supports proto3 files.")
		p.P()
		return true
	}
	for _, message := range file.Messages {
		p.generateJSONMethods(message)
	}
	return true
}

func (p *jsonMarshal) generateJSONMethods(message *protogen.Message) {
	if message.Desc.IsMapEntry() {
		return
	}
	ccTypeName := message.GoIdent.GoName
	p.P(`func (m *`, ccTypeName, `) MarshalJSON() ([]byte, error) {`)
	p.P(`container := `, gabsPackage.Ident("New"), `()`)
	for _, field := range message.Fields {
		jsonName := strs.JSONCamelCase(string(field.Desc.Name()))
		fieldName := field.GoName
		nullable := field.Desc.HasPresence()
		repeated := field.Desc.Cardinality() == protoreflect.Repeated
		if nullable {
			p.P(`if m.`, fieldName, ` != nil {`)
			p.marshalField(field, `m.`+fieldName, jsonName)
			p.P(`}`)
		} else if repeated {
			p.P(`if len(m.`, fieldName, `) > 0 {`)
			p.P(`jsonFields := make([]interface{}, len(m.`, fieldName, `))`)
			p.P(`for i, val := range m.`, fieldName, ` {`)
			p.marshalRepeatedField(field, `val`, `jsonFields`, `i`)
			p.P(`}`)
			p.P(`container.Set(jsonFields, "`, jsonName, `")`)
			p.P(`}`)
		} else {
			switch field.Desc.Kind() {
			case protoreflect.StringKind:
				p.P(`if m.`, fieldName, ` != "" {`)
			case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.FloatKind, protoreflect.DoubleKind:
				p.P(`if m.`, fieldName, ` != 0 {`)
			case protoreflect.BoolKind:
				p.P(`if m.`, fieldName, ` {`)
			case protoreflect.BytesKind:
				p.P(`if len(m.`, fieldName, `) > 0 {`)
			case protoreflect.EnumKind:
				p.P(`if int(m.`, fieldName, `) != 0 {`)
			case protoreflect.MessageKind, protoreflect.GroupKind:
				// Always marshal non-nil message fields
				p.marshalField(field, `m.`+fieldName, jsonName)
				continue
			default:
				p.P(`if m.`, fieldName, ` != nil {`)
			}
			p.marshalField(field, `m.`+fieldName, jsonName)
			p.P(`}`)
		}
	}
	p.P(`return container.MarshalJSON()`)
	p.P(`}`)
	p.P()
	p.P(`func (m *`, ccTypeName, `) UnmarshalJSON(data []byte) error {`)
	p.P(`var p `, fastjsonPackage.Ident("Parser"))
	p.P(`v, err := p.ParseBytes(data)`)
	p.P(`if err != nil {`)
	p.P(`return err`)
	p.P(`}`)
	p.P(`return m.UnmarshalJSONValue(v)`)
	p.P(`}`)
	p.P()
	p.P(`func (m *`, ccTypeName, `) UnmarshalJSONValue(v *`, fastjsonPackage.Ident("Value"), `) error {`)
	p.P("if v == nil { return nil }")
	for _, field := range message.Fields {
		protoName := string(field.Desc.Name())
		jsonName := strs.JSONCamelCase(protoName)
		fieldName := field.GoName
		p.P(`if v.Exists("`, jsonName, `") {`)
		p.unmarshalField(field, `m.`+fieldName, jsonName)
		if protoName != jsonName {
			p.P(`} else if v.Exists("`, string(field.Desc.Name()), `") {`)
			p.unmarshalField(field, `m.`+fieldName, protoName)
		}
		p.P(`}`)
	}
	p.P(`return nil`)
	p.P(`}`)
	p.P()
}

func (p *jsonMarshal) marshalRepeatedField(field *protogen.Field, accessor, jsonFields, index string) {
	switch field.Desc.Kind() {
	case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.FloatKind, protoreflect.DoubleKind, protoreflect.BoolKind, protoreflect.StringKind, protoreflect.BytesKind, protoreflect.EnumKind:
		p.P(jsonFields, `[`, index, `] = `, accessor)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		p.P(`jsonData, err := `, accessor, `.MarshalJSON()`)
		p.P(`if err != nil {`)
		p.P(`return nil, err`)
		p.P(`}`)
		p.P(jsonFields, `[`, index, `] = jsonData`)
	default:
		p.P(`// Unsupported type `, field.Desc.Kind())
	}
}

func (p *jsonMarshal) marshalField(field *protogen.Field, accessor, jsonName string) {
	switch field.Desc.Kind() {
	case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.FloatKind, protoreflect.DoubleKind, protoreflect.BoolKind:
		p.P(`container.Set(`, accessor, `, "`, jsonName, `")`)
	case protoreflect.StringKind:
		p.P(`container.Set(`, accessor, `, "`, jsonName, `")`)
	case protoreflect.BytesKind:
		p.P(`container.Set(string(`, accessor, `), "`, jsonName, `")`)
	case protoreflect.EnumKind:
		p.P(`container.Set(`, accessor, `.String(), "`, jsonName, `")`)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		p.P(`if `, accessor, ` != nil {`)
		p.P(`jsonData, err := `, accessor, `.MarshalJSON()`)
		p.P(`if err != nil {`)
		p.P(`return nil, err`)
		p.P(`}`)
		p.P(`container.Set(jsonData, "`, jsonName, `")`)
		p.P(`}`)
	default:
		p.P(`// Unsupported type `, field.Desc.Kind())
	}
}

func (p *jsonMarshal) unmarshalField(field *protogen.Field, accessor, jsonName string) {
	repeated := field.Desc.Cardinality() == protoreflect.Repeated
	if repeated {
		p.P(`jsonArray := v.GetArray("`, jsonName, `")`)
		p.P(`if jsonArray != nil {`)
		p.P(accessor, ` = make([]*`, field.GoIdent.GoName, `, len(jsonArray))`)
		p.P(`for i, jsonValue := range jsonArray {`)
		p.unmarshalRepeatedField(field, accessor+`[i]`, `jsonValue`)
		p.P(`}`)
		p.P(`}`)
	} else {
		switch field.Desc.Kind() {
		case protoreflect.Int32Kind:
			p.P(accessor, ` = int32(v.GetInt("`, jsonName, `"))`)
		case protoreflect.Int64Kind:
			p.P(accessor, ` = v.GetInt64("`, jsonName, `")`)
		case protoreflect.Uint32Kind, protoreflect.Uint64Kind:
			p.P(accessor, ` = v.GetUint64("`, jsonName, `")`)
		case protoreflect.FloatKind:
			p.P(accessor, ` = float32(v.GetFloat64("`, jsonName, `"))`)
		case protoreflect.DoubleKind:
			p.P(accessor, ` = v.GetFloat64("`, jsonName, `")`)
		case protoreflect.BoolKind:
			p.P(accessor, ` = v.GetBool("`, jsonName, `")`)
		case protoreflect.StringKind:
			p.P(accessor, ` = string(v.GetStringBytes("`, jsonName, `"))`)
		case protoreflect.BytesKind:
			p.P(accessor, ` = v.GetStringBytes("`, jsonName, `")`)
		case protoreflect.EnumKind:
			enumName := field.Enum.GoIdent.GoName
			p.P(accessor, ` = `, enumName, `(v.GetInt("`, jsonName, `"))`)
		case protoreflect.MessageKind, protoreflect.GroupKind:
			p.P(`jsonValue := v.Get("`, jsonName, `")`)
			p.P(`if jsonValue == nil {`)
			p.P(accessor, " = nil")
			p.P(`} else {`)
			p.P(`err := `, accessor, `.UnmarshalJSONValue(jsonValue)`)
			p.P(`if err != nil {`)
			p.P(`return err`)
			p.P(`}`)
			p.P(`}`)
		default:
			p.P(`// Unsupported type `, field.Desc.Kind())
		}
	}
}

func (p *jsonMarshal) unmarshalRepeatedField(field *protogen.Field, accessor, jsonValue string) {
	switch field.Desc.Kind() {
	case protoreflect.Int32Kind:
		p.P(accessor, ` = int32(`, jsonValue, `.GetInt())`)
	case protoreflect.Int64Kind:
		p.P(accessor, ` = `, jsonValue, `.GetInt64()`)
	case protoreflect.Uint32Kind, protoreflect.Uint64Kind:
		p.P(accessor, ` = `, jsonValue, `.GetUint64()`)
	case protoreflect.FloatKind:
		p.P(accessor, ` = float32(`, jsonValue, `.GetFloat64())`)
	case protoreflect.DoubleKind:
		p.P(accessor, ` = `, jsonValue, `.GetFloat64()`)
	case protoreflect.BoolKind:
		p.P(accessor, ` = `, jsonValue, `.GetBool()`)
	case protoreflect.StringKind:
		p.P(accessor, ` = string(`, jsonValue, `.GetStringBytes())`)
	case protoreflect.BytesKind:
		p.P(accessor, ` = `, jsonValue, `.GetStringBytes()`)
	case protoreflect.EnumKind:
		enumName := field.Enum.GoIdent.GoName
		p.P(accessor, ` = `, enumName, `(`, jsonValue, `.GetInt())`)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		p.P(`if `, jsonValue, ` == nil {`)
		p.P(accessor, " = nil")
		p.P(`} else {`)
		p.P(`err := `, accessor, `.UnmarshalJSONValue(`, jsonValue, `)`)
		p.P(`if err != nil {`)
		p.P(`return err`)
		p.P(`}`)
		p.P(`}`)
	default:
		p.P(`// Unsupported type `, field.Desc.Kind())
	}
}
