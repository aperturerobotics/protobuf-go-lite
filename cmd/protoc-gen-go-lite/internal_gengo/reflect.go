// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal_gengo

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/aperturerobotics/protobuf-go-lite/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// stripSourceRetentionFieldsFromMessage walks the given message tree recursively
// and clears any fields with the field option: [retention = RETENTION_SOURCE]
func stripSourceRetentionFieldsFromMessage(m protoreflect.Message) {
	protorange.Range(m, func(ppv protopath.Values) error {
		m2, ok := ppv.Index(-1).Value.Interface().(protoreflect.Message)
		if !ok {
			return nil
		}
		m2.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
			fdo, ok := fd.Options().(*descriptorpb.FieldOptions)
			if ok && fdo.GetRetention() == descriptorpb.FieldOptions_RETENTION_SOURCE {
				m2.Clear(fd)
			}
			return true
		})
		return nil
	})
}

func genFileDescriptor(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo) {
	descProto := proto.Clone(f.Proto).(*descriptorpb.FileDescriptorProto)
	descProto.SourceCodeInfo = nil // drop source code information
	stripSourceRetentionFieldsFromMessage(descProto.ProtoReflect())
	b, err := proto.MarshalOptions{AllowPartial: true, Deterministic: true}.Marshal(descProto)
	if err != nil {
		gen.Error(err)
		return
	}

	g.P("var ", rawDescVarName(f), " = []byte{")
	for len(b) > 0 {
		n := 16
		if n > len(b) {
			n = len(b)
		}

		s := ""
		for _, c := range b[:n] {
			s += fmt.Sprintf("0x%02x,", c)
		}
		g.P(s)

		b = b[n:]
	}
	g.P("}")
	g.P()

	if f.needRawDesc {
		onceVar := rawDescVarName(f) + "Once"
		dataVar := rawDescVarName(f) + "Data"
		g.P("var (")
		g.P(onceVar, " ", syncPackage.Ident("Once"))
		g.P(dataVar, " = ", rawDescVarName(f))
		g.P(")")
		g.P()
	}
}

func fileVarName(f *protogen.File, suffix string) string {
	prefix := f.GoDescriptorIdent.GoName
	_, n := utf8.DecodeRuneInString(prefix)
	prefix = strings.ToLower(prefix[:n]) + prefix[n:]
	return prefix + "_" + suffix
}

func rawDescVarName(f *fileInfo) string {
	return fileVarName(f.File, "rawDesc")
}

func goTypesVarName(f *fileInfo) string {
	return fileVarName(f.File, "goTypes")
}

func depIdxsVarName(f *fileInfo) string {
	return fileVarName(f.File, "depIdxs")
}

func enumTypesVarName(f *fileInfo) string {
	return fileVarName(f.File, "enumTypes")
}

func messageTypesVarName(f *fileInfo) string {
	return fileVarName(f.File, "msgTypes")
}

func initFuncName(f *protogen.File) string {
	return fileVarName(f, "init")
}
