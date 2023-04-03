// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"

	plugin "github.com/TheThingsIndustries/protoc-gen-go-json/internal/gen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-json %v\n", plugin.Version)
		return
	}

	var flags flag.FlagSet
	flags.BoolVar(&plugin.Params.Std, "std", false, "generate standard library (un)marshalers")
	flags.BoolVar(&plugin.Params.LegacyFieldMaskMarshalling, "legacy_fieldmask_marshaling", false, "generate legacy fieldmask marshalers")
	flags.StringVar(&plugin.Params.Lang, "lang", "go", "language (go)")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			if f.Desc.Syntax() != protoreflect.Proto3 {
				return fmt.Errorf("the protoc-gen-go-json plugin only supports proto3 syntax")
			}
			plugin.GenerateFile(gen, f)
		}
		return nil
	})
}
