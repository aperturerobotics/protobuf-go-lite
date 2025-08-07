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
	slicesPackage  = protogen.GoImportPath("slices")
	mapsPackage    = protogen.GoImportPath("maps")
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
	g.P("sb.WriteString(\"", message.Desc.Name(), " {\")")

	// initialSbLen is the sb contents length before anything is written
	initialSbLen := len(message.Desc.Name()) + 2 // +2 for " {"

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
				g.genField(initialSbLen, oneofField, "body."+oneofField.GoName)
			}
			g.P("}")
		} else {
			accessor := "x." + field.GoName
			g.genField(initialSbLen, field, accessor)
		}
	}
	g.P("sb.WriteString(\"}\")")
	g.P("return sb.String()")
	g.P("}")

	g.P()

	g.P("func (x *", message.GoIdent, ") String() string {")
	g.P("return x.MarshalProtoText()")
	g.P("}")
}

func (g *textGenerator) genField(sbInitialLen int, field *protogen.Field, accessor string) {
	maybeAddSpace := func() {
		g.P("if sb.Len() > ", sbInitialLen, " {")
		g.P("sb.WriteString(\" \")")
		g.P("}")
	}

	if field.Desc.IsList() {
		g.P("if len(", accessor, ") > 0 {")
		maybeAddSpace()
		g.P("sb.WriteString(\"", field.Desc.Name(), ": [\")")
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
			maybeAddSpace()
			g.P("sb.WriteString(\"", field.Desc.Name(), ": {\")")
			g.P("for _, k := range ", g.QualifiedGoIdent(slicesPackage.Ident("Sorted")), "(", g.QualifiedGoIdent(mapsPackage.Ident("Keys")), "(", accessor, ")) {")
			g.P("v := ", accessor, "[k]")
			g.P("sb.WriteString(\" \")")
			g.genFieldValue(field.Message.Fields[0], "k", false)
			g.P("sb.WriteString(\": \")")
			g.genFieldValue(field.Message.Fields[1], "v", false)
			g.P("}")
			g.P("sb.WriteString(\" }\")")
			g.P("}")
		} else {
			g.P("if ", accessor, " != nil {")
			maybeAddSpace()
			g.P("sb.WriteString(\"", field.Desc.Name(), ": \")")
			g.P("sb.WriteString(", accessor, ".MarshalProtoText())")
			g.P("}")
		}
	case protoreflect.StringKind, protoreflect.BytesKind:
		if field.Desc.HasOptionalKeyword() || g.file.Desc.Syntax() == protoreflect.Proto2 {
			g.P("if ", accessor, " != nil {")
		} else {
			emptyCheck := "\"\""
			if field.Desc.Kind() == protoreflect.BytesKind {
				emptyCheck = "nil"
			}
			g.P("if ", accessor, " != ", emptyCheck, " {")
		}
		maybeAddSpace()
		g.P("sb.WriteString(\"", field.Desc.Name(), ": \")")
		g.genFieldValue(field, accessor, false)
		g.P("}")
	case protoreflect.EnumKind, protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind, protoreflect.BoolKind:
		zeroValue := "0"
		if field.Desc.Kind() == protoreflect.BoolKind {
			zeroValue = "false"
		}
		if field.Desc.HasOptionalKeyword() || g.file.Desc.Syntax() == protoreflect.Proto2 {
			g.P("if ", accessor, " != nil {")
		} else {
			g.P("if ", accessor, " != ", zeroValue, " {")
		}
		maybeAddSpace()
		g.P("sb.WriteString(\"", field.Desc.Name(), ": \")")
		g.genFieldValue(field, accessor, false)
		g.P("}")
	default:
		// For unsupported field types, do nothing
		return
	}
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
			g.P("sb.WriteString(\"\\\"\")")
			g.P("sb.WriteString(", accessor, ".String())")
			g.P("sb.WriteString(\"\\\"\")")
		} else {
			g.P("sb.WriteString(\"\\\"\")")
			g.P("sb.WriteString(", field.Enum.GoIdent, "(", accessor, ").String())")
			g.P("sb.WriteString(\"\\\"\")")
		}
	case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
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
