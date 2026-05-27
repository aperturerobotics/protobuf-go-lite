package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const edition2024Fixture = "testproto/editions2024/editions2024.proto"

func TestEdition2024DefaultFeaturesCompile(t *testing.T) {
	root := repoRoot(t)
	plugin := buildCurrentPlugin(t, root)
	outDir := t.TempDir()

	cmd := exec.Command(
		"protoc",
		"-I", ".",
		"--plugin=protoc-gen-go-lite="+plugin,
		"--go-lite_out="+outDir,
		"--go-lite_opt=features=all,paths=source_relative",
		edition2024Fixture,
	)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("current plugin should generate Edition 2024 output:\n%s", out)
	}
	assertGeneratedEditionFixture(t, root, outDir, "current plugin")
}

func TestEdition2024AdvertiseOnlyDefaultFeaturesCompile(t *testing.T) {
	root := repoRoot(t)
	plugin := buildAdvertiseOnlyPlugin(t, root)
	outDir := t.TempDir()

	cmd := exec.Command(
		"protoc",
		"-I", ".",
		"--plugin=protoc-gen-go-lite="+plugin,
		"--go-lite_out="+outDir,
		"--go-lite_opt=features=all,paths=source_relative",
		edition2024Fixture,
	)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("advertise-only shim should let protoc generate output:\n%s", out)
	}

	assertGeneratedEditionFixture(t, root, outDir, "advertise-only generated output")
}

func TestEdition2024RejectsClosedEnum(t *testing.T) {
	root := repoRoot(t)
	protoPath := writeTempProto(t, `edition = "2024";

package closedenum;

option go_package = "editionfixture/closedenum;closedenum";
option features.enum_type = CLOSED;

enum State {
  STATE_UNKNOWN = 0;
  STATE_READY = 1;
}

message Msg {
  State state = 1;
}
`)

	out, err := runAdvertiseOnlyProtoc(t, root, protoPath)
	if err == nil {
		t.Fatal("expected closed enum policy rejection")
	}
	if !strings.Contains(out, "closed Edition") {
		t.Fatalf("expected closed enum rejection, got:\n%s", out)
	}
}

func TestEdition2024RejectsLegacyBestEffortJSON(t *testing.T) {
	root := repoRoot(t)
	protoPath := writeTempProto(t, `edition = "2024";

package legacyjson;

option go_package = "editionfixture/legacyjson;legacyjson";
option features.json_format = LEGACY_BEST_EFFORT;

message Msg {
  string name = 1;
}
`)

	out, err := runAdvertiseOnlyProtoc(t, root, protoPath)
	if err == nil {
		t.Fatal("expected legacy JSON policy rejection")
	}
	if !strings.Contains(out, "LEGACY_BEST_EFFORT") {
		t.Fatalf("expected legacy JSON rejection, got:\n%s", out)
	}
}

func TestEdition2024RejectsOpaqueGoAPI(t *testing.T) {
	root := repoRoot(t)
	protoPath := writeTempProto(t, `edition = "2024";

package opaqueapi;

import "google/protobuf/go_features.proto";

option go_package = "editionfixture/opaqueapi;opaqueapi";
option features.(pb.go).api_level = API_OPAQUE;

message Msg {
  string name = 1;
}
`)

	out, err := runAdvertiseOnlyProtoc(t, root, protoPath, protobufSourceDir(t, root))
	if err == nil {
		t.Fatal("expected opaque API policy rejection")
	}
	if !strings.Contains(out, "API_OPAQUE") {
		t.Fatalf("expected opaque API rejection, got:\n%s", out)
	}
}

func assertGeneratedEditionFixture(t *testing.T, root, outDir, label string) {
	t.Helper()

	writeFile(t, filepath.Join(outDir, "go.mod"), "module editionfixture\n\ngo 1.25\n\nrequire github.com/aperturerobotics/protobuf-go-lite v0.0.0\n\nreplace github.com/aperturerobotics/protobuf-go-lite => "+root+"\n")
	writeFile(t, filepath.Join(outDir, "testproto/editions2024/editions2024_behavior_test.go"), edition2024BehaviorTest)

	cmd := exec.Command("go", "test", "-mod=mod", "./...")
	cmd.Dir = outDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s should compile and pass behavior tests:\n%s", label, out)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func buildCurrentPlugin(t *testing.T, root string) string {
	t.Helper()

	plugin := filepath.Join(t.TempDir(), "protoc-gen-go-lite")
	cmd := exec.Command("go", "build", "-o", plugin, "./cmd/protoc-gen-go-lite")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build current plugin:\n%s", out)
	}
	return plugin
}

func buildAdvertiseOnlyPlugin(t *testing.T, root string) string {
	t.Helper()

	dir := t.TempDir()
	mainPath := filepath.Join(dir, "main.go")
	writeFile(t, mainPath, advertiseOnlyPluginSource)

	plugin := filepath.Join(dir, "protoc-gen-go-lite")
	cmd := exec.Command("go", "build", "-o", plugin, mainPath)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build advertise-only plugin:\n%s", out)
	}
	return plugin
}

func writeTempProto(t *testing.T, dat string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "fixture.proto")
	writeFile(t, path, dat)
	return path
}

func runAdvertiseOnlyProtoc(t *testing.T, root, protoPath string, extraImports ...string) (string, error) {
	t.Helper()

	plugin := buildAdvertiseOnlyPlugin(t, root)
	args := []string{
		"-I", ".",
		"-I", filepath.Dir(protoPath),
		"--plugin=protoc-gen-go-lite=" + plugin,
		"--go-lite_out=" + t.TempDir(),
		"--go-lite_opt=features=all",
		protoPath,
	}
	for _, imp := range extraImports {
		args = append([]string{"-I", imp}, args...)
	}
	cmd := exec.Command("protoc", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func protobufSourceDir(t *testing.T, root string) string {
	t.Helper()

	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "google.golang.org/protobuf")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("find protobuf module dir:\n%s", out)
	}
	return filepath.Join(strings.TrimSpace(string(out)), "src")
}

func writeFile(t *testing.T, path, dat string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(dat), 0o600); err != nil {
		t.Fatal(err)
	}
}

const advertiseOnlyPluginSource = `package main

import (
	"flag"
	"strings"

	"github.com/aperturerobotics/protobuf-go-lite/compiler/protogen"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/clone"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/equal"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/json"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/marshal"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/size"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/text"
	_ "github.com/aperturerobotics/protobuf-go-lite/features/unmarshal"
	"github.com/aperturerobotics/protobuf-go-lite/generator"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var cfg generator.Config
	var features string
	var f flag.FlagSet

	f.BoolVar(&cfg.AllowEmpty, "allow-empty", false, "")
	f.StringVar(&features, "features", "all", "")
	f.StringVar(&cfg.BuildTag, "buildTag", "", "")

	protogen.Options{
		ParamFunc: f.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		gen, err := generator.NewGenerator(plugin, strings.Split(features, "+"), &cfg)
		if err != nil {
			return err
		}
		plugin.SupportedFeatures |= uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
		plugin.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO2
		plugin.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_2024
		gen.Generate()
		return nil
	})
}
`

const edition2024BehaviorTest = `package editions2024

import (
	stdjson "encoding/json"
	"strings"
	"testing"
)

func TestEdition2024GeneratedBehavior(t *testing.T) {
	required := int32(7)
	name := "nested"
	label := "group"
	explicit := "value"
	state := Edition2024Fixture_STATE_READY
	msg := &Edition2024Fixture{
		RequiredInt32:  &required,
		ExplicitString: &explicit,
		ExplicitState:  &state,
		NestedMessage:  &Edition2024Fixture_Nested{Name: &name, Value: 12},
		DelimitedGroup: &Edition2024Fixture_DelimitedGroup{Label: &label},
		NestedMap: map[string]*Edition2024Fixture_Nested{
			"key": &Edition2024Fixture_Nested{Name: &name, Value: 13},
		},
	}

	wire, err := msg.MarshalVT()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got Edition2024Fixture
	if err := got.UnmarshalVT(wire); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.DelimitedGroup == nil || got.DelimitedGroup.GetLabel() != "group" {
		t.Fatalf("delimited group round trip failed: %#v", got.DelimitedGroup)
	}
	var gotUnsafe Edition2024Fixture
	if err := gotUnsafe.UnmarshalVTUnsafe(wire); err != nil {
		t.Fatalf("unsafe unmarshal: %v", err)
	}
	if gotUnsafe.DelimitedGroup == nil || gotUnsafe.DelimitedGroup.GetLabel() != "group" {
		t.Fatalf("unsafe delimited group round trip failed: %#v", gotUnsafe.DelimitedGroup)
	}
	if !msg.EqualVT(msg.CloneVT()) {
		t.Fatal("clone/equal mismatch")
	}

	j, err := stdjson.Marshal(msg)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}
	var fromJSON Edition2024Fixture
	if err := stdjson.Unmarshal(j, &fromJSON); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if fromJSON.GetExplicitString() != "value" || fromJSON.GetExplicitState() != Edition2024Fixture_STATE_READY {
		t.Fatalf("json round trip mismatch: %#v", fromJSON)
	}

	invalidUTF8 := []byte{0x18, 0x01, 0x22, 0x01, 0xff}
	var invalid Edition2024Fixture
	if err := invalid.UnmarshalVT(invalidUTF8); err == nil || !strings.Contains(err.Error(), "invalid UTF-8") {
		t.Fatalf("safe unmarshal invalid UTF-8 error = %v", err)
	}
	var invalidUnsafe Edition2024Fixture
	if err := invalidUnsafe.UnmarshalVTUnsafe(invalidUTF8); err == nil || !strings.Contains(err.Error(), "invalid UTF-8") {
		t.Fatalf("unsafe unmarshal invalid UTF-8 error = %v", err)
	}
}
`
