// Copyright © 2024 Aperture Robotics, LLC.
// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"github.com/aperturerobotics/protobuf-go-lite/compiler/protogen"
	"github.com/aperturerobotics/protobuf-go-lite/generator"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	jsonPluginPackage = protogen.GoImportPath("github.com/aperturerobotics/protobuf-go-lite/json")
)

type jsonGenerator struct {
	gen  *protogen.Plugin
	file *protogen.File
	*generator.GeneratedFile
}

func init() {
	generator.RegisterFeature("json", func(gen *generator.GeneratedFile) generator.FeatureGenerator {
		return &jsonGenerator{GeneratedFile: gen}
	})
}

func (g *jsonGenerator) GenerateFile(file *protogen.File) bool {
	g.file = file

	// If the file doesn't have marshalers or unmarshalers, we can skip it.
	if !g.fileHasAnyMarshaler() {
		return false
	}

	// fields with pointers are not supported (proto2).
	if file.Desc.Syntax() != protoreflect.Proto3 {
		g.P("// NOTE: protobuf-go-lite json only supports proto3: ", file.Desc.Syntax().String(), " is not supported.")
		g.P()
		return true
	}

	g.generateFileContent()
	return true
}

func (g *jsonGenerator) fileHasAnyMarshaler() bool {
	return len(g.file.Enums) != 0 || len(g.file.Messages) != 0
}

func (g *jsonGenerator) generateFileContent() {
	for _, enum := range g.file.Enums {
		g.genEnum(enum)
	}
	for _, message := range g.file.Messages {
		g.genMessage(message)
	}
}
