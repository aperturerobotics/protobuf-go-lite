package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const jsonOptionalFieldMaskProto = `syntax = "proto3";

package jsonfixture;

option go_package = "jsonfixture;jsonfixture";

message Msg {
  optional int64 optional_value = 1;
  int64 value = 2;
}
`

const jsonOptionalFieldMaskTest = `package jsonfixture

import (
	"bytes"
	"testing"

	json "github.com/aperturerobotics/protobuf-go-lite/json"
)

func marshalWithMask(t *testing.T, msg *Msg, fields ...string) string {
	t.Helper()

	var buf bytes.Buffer
	state := json.NewMarshalState(json.DefaultMarshalerConfig, json.NewJsonStream(&buf)).WithFieldMask(fields...)
	msg.MarshalProtoJSON(state)
	if err := state.Err(); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestOptionalFieldMask(t *testing.T) {
	if got := marshalWithMask(t, &Msg{}, "optionalValue"); got != "{}" {
		t.Fatalf("absent optional = %s, want {}", got)
	}

	zero := int64(0)
	if got := marshalWithMask(t, &Msg{OptionalValue: &zero}); got != "{\"optionalValue\":\"0\"}" {
		t.Fatalf("present optional zero = %s", got)
	}

	nonzero := int64(7)
	if got := marshalWithMask(t, &Msg{OptionalValue: &nonzero}); got != "{\"optionalValue\":\"7\"}" {
		t.Fatalf("present optional nonzero = %s", got)
	}

	if got := marshalWithMask(t, &Msg{}, "value"); got != "{\"value\":\"0\"}" {
		t.Fatalf("selected non-optional zero = %s", got)
	}
}
`

func TestJSONOptionalFieldMask(t *testing.T) {
	root := repoRoot(t)
	plugin := buildCurrentPlugin(t, root)
	protoPath := writeTempProto(t, jsonOptionalFieldMaskProto)
	outDir := t.TempDir()

	cmd := exec.Command(
		"protoc",
		"-I", filepath.Dir(protoPath),
		"--plugin=protoc-gen-go-lite="+plugin,
		"--go-lite_out="+outDir,
		"--go-lite_opt=features=json,paths=source_relative",
		protoPath,
	)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generate JSON optional fixture:\n%s", out)
	}

	writeFile(t, filepath.Join(outDir, "go.mod"), "module jsonfixture\n\ngo 1.25\n\nrequire github.com/aperturerobotics/protobuf-go-lite v0.0.0\n\nreplace github.com/aperturerobotics/protobuf-go-lite => "+root+"\n")
	writeFile(t, filepath.Join(outDir, "fixture_test.go"), jsonOptionalFieldMaskTest)

	cmd = exec.Command("go", "test", "-mod=mod", "./...")
	cmd.Dir = outDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated JSON optional fixture:\n%s", out)
	}

	generatedPath := filepath.Join(outDir, filepath.Base(strings.TrimSuffix(protoPath, ".proto")+".pb.go"))
	generated, err := os.ReadFile(generatedPath)
	if err != nil {
		t.Fatal(err)
	}
	generatedText := string(generated)
	if strings.Contains(generatedText, `OptionalValue != nil || s.HasField("optionalValue")`) {
		t.Fatalf("optional scalar still forced by field mask:\n%s", generatedText)
	}
	if !strings.Contains(generatedText, `if x.OptionalValue != nil {`) {
		t.Fatalf("optional scalar presence guard missing:\n%s", generatedText)
	}
	if !strings.Contains(generatedText, `if x.Value != 0 || s.HasField("value") {`) {
		t.Fatalf("non-optional field-mask condition changed:\n%s", generatedText)
	}
}
