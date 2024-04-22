// Copyright © 2024 Aperture Robotics, LLC.
// SPDX-License-Identifier: MIT

package text

import (
	"strings"

	"github.com/aperturerobotics/protobuf-go-lite/compiler/protogen"
	"github.com/aperturerobotics/protobuf-go-lite/generator"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	stringsPackage = protogen.GoImportPath("strings")
	strconvPackage = protogen.GoImportPath("strconv")
	base64Package  = protogen.GoImportPath("encoding/base64")
)

var disableTextComment = "protobuf-go-lite:disable-text"

// hasDisableTextComment checks if a comments section has the disable text comment.
func hasDisableTextComment(comments protogen.Comments) bool {
	for _, line := range strings.Split(strings.TrimSuffix(string(comments), "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == disableTextComment {
			return true
		}
	}
	return false
}

type textGenerator struct {
	file *protogen.File
	*generator.GeneratedFile
}

func init() {
	generator.RegisterFeature("text", func(gen *generator.GeneratedFile) generator.FeatureGenerator {
		return &textGenerator{GeneratedFile: gen}
	})
}

func (g *textGenerator) GenerateFile(file *protogen.File) bool {
	g.file = file
	// If the file doesn't have marshalers or unmarshalers, we can skip it.
	if !g.fileHasAnyMarshaler() {
		return false
	}
	g.generateFileContent()
	return true
}

func (g *textGenerator) fileHasAnyMarshaler() bool {
	return len(g.file.Enums) != 0 || len(g.file.Messages) != 0
}

func (g *textGenerator) generateFileContent() {
	for _, enum := range g.file.Enums {
		g.genEnum(enum)
	}
	for _, message := range g.file.Messages {
		g.genMessage(message)
	}
}

func (g *textGenerator) genEnum(enum *protogen.Enum) {
	if !hasDisableTextComment(enum.Comments.Leading) {
		// Generate enum text marshaling code
		g.P("func (x ", enum.GoIdent, ") MarshalProtoText() string {")
		g.P("return x.String()")
		g.P("}")
	}
}

func (g *textGenerator) genMessage(message *protogen.Message) {
	// Generate marshalers and unmarshalers for all enums defined in the message.
	for _, enum := range message.Enums {
		g.genEnum(enum)
	}
	// Generate marshalers and unmarshalers for all sub-messages defined in the message.
	for _, message := range message.Messages {
		g.genMessage(message)
	}
	// skip early if the disable comment is present
	if hasDisableTextComment(message.Comments.Leading) {
		return
	}
	// Generate message text marshaling code
	g.P("func (x *", message.GoIdent, ") MarshalProtoText() string {")
	g.P("var sb ", g.QualifiedGoIdent(stringsPackage.Ident("Builder")))
	g.P("sb.WriteString(\"", message.Desc.Name(), " { \")")
	handledOneOfs := make(map[string]struct{})
	for _, field := range message.Fields {
		if oneof := field.Oneof; oneof != nil && !field.Desc.HasOptionalKeyword() {
			if _, ok := handledOneOfs[oneof.GoName]; ok {
				continue
			}
			handledOneOfs[oneof.GoName] = struct{}{}
			g.P("switch body := x.", oneof.GoName, ".(type) {")
			for _, oneofField := range oneof.Fields {
				g.P("case *", oneofField.GoIdent, ":")
				g.genField(oneofField, "body."+oneofField.GoName)
			}
			g.P("}")
		} else {
			accessor := "x." + field.GoName
			g.genField(field, accessor)
		}
	}
	g.P("sb.WriteString(\"}\")")
	g.P("return sb.String()")
	g.P("}")

	g.P("func (x *", message.GoIdent, ") String() string {")
	g.P("return x.MarshalProtoText()")
	g.P("}")
}

func (g *textGenerator) genField(field *protogen.Field, accessor string) {
	if field.Desc.IsList() {
		g.P("if len(", accessor, ") > 0 {")
		g.P("sb.WriteString(\" ", field.Desc.Name(), ": [\")")
		g.P("for i, v := range ", accessor, " {")
		g.P("if i > 0 {")
		g.P("sb.WriteString(\", \")")
		g.P("}")
		g.genFieldValue(field, "v", true)
		g.P("}")
		g.P("sb.WriteString(\"]\")")
		g.P("}")
		return
	}

	switch field.Desc.Kind() {
	case protoreflect.MessageKind:
		if field.Desc.IsMap() {
			g.P("if len(", accessor, ") > 0 {")
			g.P("sb.WriteString(\" ", field.Desc.Name(), ": {\")")
			g.P("for k, v := range ", accessor, " {")
			g.P("sb.WriteString(\" \")")
			g.genFieldValue(field.Message.Fields[0], "k", false)
			g.P("sb.WriteString(\": \")")
			g.genFieldValue(field.Message.Fields[1], "v", false)
			g.P("}")
			g.P("sb.WriteString(\" }\")")
		} else {
			g.P("if ", accessor, " != nil {")
			g.P("sb.WriteString(\" ", field.Desc.Name(), ": \")")
			g.P("sb.WriteString(", accessor, ".MarshalProtoText())")
		}
	case protoreflect.StringKind:
		if field.Desc.HasOptionalKeyword() || g.file.Desc.Syntax() == protoreflect.Proto2 {
			g.P("if ", accessor, " != nil {")
		} else {
			g.P("if ", accessor, " != \"\" {")
		}
		g.P("sb.WriteString(\" ", field.Desc.Name(), ": \")")
		g.genFieldValue(field, accessor, false)
	case protoreflect.EnumKind, protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind, protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind, protoreflect.FloatKind, protoreflect.DoubleKind:
		if field.Desc.HasOptionalKeyword() || g.file.Desc.Syntax() == protoreflect.Proto2 {
			g.P("if ", accessor, " != nil {")
		} else {
			g.P("if ", accessor, " != 0 {")
		}
		g.P("sb.WriteString(\" ", field.Desc.Name(), ": \")")
		g.genFieldValue(field, accessor, false)
	case protoreflect.BytesKind:
		g.P("if len(", accessor, ") > 0 {")
		g.P("sb.WriteString(\" ", field.Desc.Name(), ": \")")
		g.genFieldValue(field, accessor, false)
	case protoreflect.BoolKind:
		if field.Desc.HasOptionalKeyword() || g.file.Desc.Syntax() == protoreflect.Proto2 {
			g.P("if ", accessor, " != nil {")
		} else {
			g.P("if ", accessor, " {")
		}
		g.P("sb.WriteString(\" ", field.Desc.Name(), ": \")")
		g.genFieldValue(field, accessor, false)
	default:
		return
	}
	g.P("}")
}

func (g *textGenerator) genFieldValue(field *protogen.Field, accessor string, isList bool) {
	isPointer := field.Desc.HasOptionalKeyword() || (g.file.Desc.Syntax() == protoreflect.Proto2 && !isList)
	switch field.Desc.Kind() {
	case protoreflect.MessageKind:
		g.P("sb.WriteString(", accessor, ".MarshalProtoText())")
	case protoreflect.StringKind:
		if isPointer {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("Quote")), "(*", accessor, "))")
		} else {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("Quote")), "(", accessor, "))")
		}
	case protoreflect.BytesKind:
		g.P("sb.WriteString(\"\\\"\")")
		g.P("sb.WriteString(", g.QualifiedGoIdent(base64Package.Ident("StdEncoding.EncodeToString")), "(", accessor, "))")
		g.P("sb.WriteString(\"\\\"\")")
	case protoreflect.EnumKind:
		if isPointer {
			g.P("sb.WriteString(", accessor, ".String())")
		} else {
			g.P("sb.WriteString(", field.Enum.GoIdent, "(", accessor, ").String())")
		}
	case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Sint32Kind, protoreflect.Sint64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		if isPointer {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("FormatInt")), "(int64(*", accessor, "), 10))")
		} else {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("FormatInt")), "(int64(", accessor, "), 10))")
		}
	case protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		if isPointer {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("FormatUint")), "(uint64(*", accessor, "), 10))")
		} else {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("FormatUint")), "(uint64(", accessor, "), 10))")
		}
	case protoreflect.FloatKind:
		if isPointer {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("FormatFloat")), "(float64(*", accessor, "), 'g', -1, 32))")
		} else {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("FormatFloat")), "(float64(", accessor, "), 'g', -1, 32))")
		}
	case protoreflect.DoubleKind:
		if isPointer {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("FormatFloat")), "(*", accessor, ", 'g', -1, 64))")
		} else {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("FormatFloat")), "(", accessor, ", 'g', -1, 64))")
		}
	case protoreflect.BoolKind:
		if isPointer {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("FormatBool")), "(*", accessor, "))")
		} else {
			g.P("sb.WriteString(", g.QualifiedGoIdent(strconvPackage.Ident("FormatBool")), "(", accessor, "))")
		}
	}
}
