package generator_test

import (
	"testing"

	"github.com/aperturerobotics/protobuf-go-lite/compiler/protogen"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/marshal"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/size"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/unmarshal"
	"github.com/aperturerobotics/protobuf-go-lite/generator"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestNewGeneratorValidatesRegistryFeatures(t *testing.T) {
	plugin, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{})
	if err != nil {
		t.Fatal(err)
	}

	for _, features := range [][]string{
		{"size", "marshal"},
		{"size", "unmarshal"},
		{"marshal", "unmarshal"},
	} {
		_, err := generator.NewGenerator(plugin, features, &generator.Config{Registry: true})
		if err == nil {
			t.Fatalf("NewGenerator(%q) accepted registry without all message features", features)
		}
	}

	if _, err := generator.NewGenerator(plugin, []string{"size", "marshal", "unmarshal"}, &generator.Config{Registry: true}); err != nil {
		t.Fatalf("NewGenerator rejected complete registry features: %v", err)
	}
}
