// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The protoc-gen-go-lite binary is a protoc plugin that generates
// reflection-free Go code for proto2, proto3 and editions files. Its output
// matches protoc-gen-go's types and adds the methods of each generator
// feature it links.
//
// It takes these parameters:
//
//   - features: the generator features to emit, separated by '+', or "all"
//     (the default) for every registered feature: marshal, marshal_strict,
//     unmarshal, unmarshal_unsafe, size, clone, equal, json and text.
//   - codegen: "helper" (the default) makes generated methods call the
//     runtime's shared helpers; "unrolled" inlines that logic.
//   - allow-empty: write a file even when it has no content.
//   - buildTag: the go:build constraint to put on each generated file.
//   - registry: generate an init-time message registry with flattened custom
//     options.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aperturerobotics/protobuf-go-lite/compiler/protogen"
	"github.com/aperturerobotics/protobuf-go-lite/generator"
	"github.com/aperturerobotics/protobuf-go-lite/internal/version"

	_ "github.com/aperturerobotics/protobuf-go-lite/features/clone"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/equal"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/json"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/marshal"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/size"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/text"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/unmarshal"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Fprintf(os.Stdout, "%v %v\n", filepath.Base(os.Args[0]), version.String())
		os.Exit(0)
	}

	var cfg generator.Config
	var codegenMode string
	var features string
	var f flag.FlagSet

	f.BoolVar(&cfg.AllowEmpty, "allow-empty", false, "allow generation of empty files")
	f.StringVar(&codegenMode, "codegen", string(generator.CodegenModeHelper), "code generation mode: helper or unrolled")
	f.StringVar(&features, "features", "all", "list of features to generate (separated by '+')")
	f.StringVar(&cfg.BuildTag, "buildTag", "", "the go:build tag to set on generated files")
	f.BoolVar(&cfg.Registry, "registry", false, "generate init-time message registry with flattened custom options")

	protogen.Options{
		ParamFunc: f.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		if err := cfg.SetCodegenMode(codegenMode); err != nil {
			return err
		}
		featureNames := strings.Split(features, "+")
		gen, err := generator.NewGenerator(plugin, featureNames, &cfg)
		if err != nil {
			return err
		}
		gen.Generate()

		return nil
	})
}
