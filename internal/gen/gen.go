// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package gen

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

// Version is the version of the generator.
var Version = "0.0.0-dev"

// Params are the parameters for the generator.
var Params struct {
	Lang string
}

const (
	bytesPackage   = protogen.GoImportPath("bytes")
	fmtPackage     = protogen.GoImportPath("fmt")
	strconvPackage = protogen.GoImportPath("strconv")

	gogoProtoTypesPackage = protogen.GoImportPath("github.com/gogo/protobuf/types")

	gogoProtoJSONPackage   = protogen.GoImportPath("github.com/gogo/protobuf/jsonpb")
	golangProtoJSONPackage = protogen.GoImportPath("google.golang.org/protobuf/encoding/protojson")

	jsonPluginPackage   = protogen.GoImportPath("github.com/TheThingsIndustries/protoc-gen-go-json/jsonplugin")
	gogoPluginPackage   = protogen.GoImportPath("github.com/TheThingsIndustries/protoc-gen-go-json/gogo")
	golangPluginPackage = protogen.GoImportPath("github.com/TheThingsIndustries/protoc-gen-go-json/golang")
)

type generator struct {
	gen  *protogen.Plugin
	file *protogen.File
	*protogen.GeneratedFile
}

// GenerateFile generates a file with JSON marshalers and unmarshalers.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	g := &generator{
		gen:  gen,
		file: file,
	}

	// If the file doesn't have marshalers or unmarshalers, we kan skip it.
	if !g.fileHasAnyMarshaler() {
		return nil
	}

	// Generate a new file that ends with `_json.pb.go`.
	filename := file.GeneratedFilenamePrefix + "_json.pb.go"
	g.GeneratedFile = gen.NewGeneratedFile(filename, file.GoImportPath)

	// The standard header for generated files.
	g.P("// Code generated by protoc-gen-go-json. DO NOT EDIT.")
	g.P("// versions:")
	g.P("// - protoc-gen-go-json v", Version)
	g.P("// - protoc             ", g.protocVersion())
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}

	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	g.generateFileContent()

	return g.GeneratedFile
}

func (g *generator) fileHasAnyMarshaler() bool {
	for _, enum := range g.file.Enums {
		if g.enumHasAnyMarshaler(enum) {
			return true
		}
	}
	for _, message := range g.file.Messages {
		if g.messageHasAnyMarshaler(message) {
			return true
		}
	}
	return false
}

func (g *generator) protocVersion() string {
	v := g.gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

func (g *generator) generateFileContent() {
	for _, enum := range g.file.Enums {
		g.genEnum(enum)
	}
	for _, message := range g.file.Messages {
		g.genMessage(message)
	}
}
