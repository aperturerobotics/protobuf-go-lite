package main

import (
	"flag"
	"strings"

	_ "github.com/aperturerobotics/vtprotobuf-lite/features/clone"
	_ "github.com/aperturerobotics/vtprotobuf-lite/features/equal"
	_ "github.com/aperturerobotics/vtprotobuf-lite/features/marshal"
	_ "github.com/aperturerobotics/vtprotobuf-lite/features/size"
	_ "github.com/aperturerobotics/vtprotobuf-lite/features/unmarshal"
	"github.com/aperturerobotics/vtprotobuf-lite/generator"

	"github.com/aperturerobotics/protobuf-go-lite/compiler/protogen"
)

func main() {
	var cfg generator.Config
	var features string
	var f flag.FlagSet

	f.BoolVar(&cfg.AllowEmpty, "allow-empty", false, "allow generation of empty files")
	f.BoolVar(&cfg.Wrap, "wrap", false, "generate wrapper types")
	f.StringVar(&features, "features", "all", "list of features to generate (separated by '+')")
	f.StringVar(&cfg.BuildTag, "buildTag", "", "the go:build tag to set on generated files")

	protogen.Options{ParamFunc: f.Set}.Run(func(plugin *protogen.Plugin) error {
		gen, err := generator.NewGenerator(plugin, strings.Split(features, "+"), &cfg)
		if err != nil {
			return err
		}
		gen.Generate()
		return nil
	})
}
