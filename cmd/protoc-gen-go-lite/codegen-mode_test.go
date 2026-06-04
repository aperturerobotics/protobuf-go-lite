package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const codegenModeProto = `syntax = "proto3";

package codegenmode;

option go_package = "codegenfixture;codegenfixture";

message Msg {
  optional int32 value = 1;
  sint32 signed = 2;
  repeated int32 nums = 3;
  bytes payload = 4;
  repeated Child children = 5;
  map<string, Child> child_by_name = 6;
  Child child = 7;
  oneof choice {
    bytes choice_payload = 8;
    Child choice_child = 9;
  }
  bool enabled = 10;
  repeated fixed32 fixed_nums = 11;
  repeated bool flags = 12;
}

message Child {
  string name = 1;
}
`

func TestCodegenModeDefaultUsesHelperMethods(t *testing.T) {
	out := generateCodegenModeFixture(t)
	if !strings.Contains(out, "protobuf_go_lite.SizeVarintPtr") {
		t.Fatalf("default helper output missing SizeVarintPtr:\n%s", out)
	}
	if !strings.Contains(out, "protobuf_go_lite.SizeZigzagNonZero") {
		t.Fatalf("default helper output missing SizeZigzagNonZero:\n%s", out)
	}
	if strings.Contains(out, "protobuf_go_lite.SizeOfZigzag") {
		t.Fatalf("default helper output should not contain unrolled SizeOfZigzag:\n%s", out)
	}
	for _, expected := range []string{
		"protobuf_go_lite.ClonePtr",
		"protobuf_go_lite.CloneSlice",
		"protobuf_go_lite.CloneBytes",
		"protobuf_go_lite.CloneVTSlice",
		"protobuf_go_lite.CloneVTMap",
		"protobuf_go_lite.CloneVTValue",
		"protobuf_go_lite.EqualPtr",
		"protobuf_go_lite.EqualSlice",
		"protobuf_go_lite.EqualBytes",
		"protobuf_go_lite.IsEqualVT",
		"protobuf_go_lite.EqualVTSliceImplicit",
		"protobuf_go_lite.EqualVTMapImplicit",
		"protobuf_go_lite.EqualVTImplicit",
		"protobuf_go_lite.EncodeString",
		"protobuf_go_lite.EncodeBytes",
		"protobuf_go_lite.EncodeZigzag32",
		"protobuf_go_lite.EncodeVarintPacked",
		"protobuf_go_lite.DecodeVarintBool",
		"protobuf_go_lite.DecodeSint32",
		"protobuf_go_lite.DecodeString",
		"protobuf_go_lite.DecodeBytesAppend",
		"protobuf_go_lite.DecodeLengthDelimited",
		"protobuf_go_lite.PackedVarintElementCount",
		"protobuf_go_lite.PackedFixedElementCount",
		"protobuf_go_lite.SkipWithin",
		"protobuf_go_lite.TextStartMessage",
		"protobuf_go_lite.TextWriteFieldPrefix",
		"protobuf_go_lite.TextWriteString",
	} {
		if !strings.Contains(out, expected) {
			t.Fatalf("default helper output missing %s:\n%s", expected, out)
		}
	}
}

func TestCodegenModeUnrolledUsesPreviousMethodShape(t *testing.T) {
	out := generateCodegenModeFixture(t, "codegen=unrolled")
	if !strings.Contains(out, "protobuf_go_lite.SizeOfZigzag") {
		t.Fatalf("unrolled output missing SizeOfZigzag:\n%s", out)
	}
	if strings.Contains(out, "protobuf_go_lite.SizeZigzagNonZero") {
		t.Fatalf("unrolled output should not contain helper SizeZigzagNonZero:\n%s", out)
	}
	for _, helper := range []string{
		"protobuf_go_lite.ClonePtr",
		"protobuf_go_lite.CloneSlice",
		"protobuf_go_lite.CloneBytes",
		"protobuf_go_lite.CloneVTSlice",
		"protobuf_go_lite.CloneVTMap",
		"protobuf_go_lite.CloneVTValue",
		"protobuf_go_lite.EqualPtr",
		"protobuf_go_lite.EqualSlice",
		"protobuf_go_lite.EqualBytes",
		"protobuf_go_lite.IsEqualVT",
		"protobuf_go_lite.EqualVTSliceImplicit",
		"protobuf_go_lite.EqualVTMapImplicit",
		"protobuf_go_lite.EqualVTImplicit",
		"protobuf_go_lite.EncodeString",
		"protobuf_go_lite.EncodeBytes",
		"protobuf_go_lite.EncodeZigzag32",
		"protobuf_go_lite.EncodeVarintPacked",
		"protobuf_go_lite.DecodeVarintBool",
		"protobuf_go_lite.DecodeSint32",
		"protobuf_go_lite.DecodeString",
		"protobuf_go_lite.DecodeBytesAppend",
		"protobuf_go_lite.DecodeLengthDelimited",
		"protobuf_go_lite.PackedVarintElementCount",
		"protobuf_go_lite.PackedFixedElementCount",
		"protobuf_go_lite.SkipWithin",
		"protobuf_go_lite.TextStartMessage",
		"protobuf_go_lite.TextWriteFieldPrefix",
		"protobuf_go_lite.TextWriteString",
	} {
		if strings.Contains(out, helper) {
			t.Fatalf("unrolled output should not contain helper %s:\n%s", helper, out)
		}
	}
	for _, expected := range []string{
		"tmpVal := *rhs",
		"for i, vx := range this.Nums",
		"var sb strings.Builder",
	} {
		if !strings.Contains(out, expected) {
			t.Fatalf("unrolled output missing previous clone/equal shape %s:\n%s", expected, out)
		}
	}
}

func TestCodegenModeRejectsUnknown(t *testing.T) {
	root := repoRoot(t)
	plugin := buildCurrentPlugin(t, root)
	protoPath := writeTempProto(t, codegenModeProto)
	outDir := t.TempDir()

	cmd := exec.Command(
		"protoc",
		"-I", filepath.Dir(protoPath),
		"--plugin=protoc-gen-go-lite="+plugin,
		"--go-lite_out="+outDir,
		"--go-lite_opt=features=size,paths=source_relative,codegen=bogus",
		protoPath,
	)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected invalid codegen mode to fail")
	}
	if !strings.Contains(string(out), `unknown codegen mode: "bogus"`) {
		t.Fatalf("expected unknown codegen mode error, got:\n%s", out)
	}
}

func generateCodegenModeFixture(t *testing.T, opts ...string) string {
	t.Helper()

	root := repoRoot(t)
	plugin := buildCurrentPlugin(t, root)
	protoPath := writeTempProto(t, codegenModeProto)
	outDir := t.TempDir()

	opt := "features=size+equal+clone+marshal+unmarshal+text,paths=source_relative"
	if len(opts) != 0 {
		opt += "," + strings.Join(opts, ",")
	}

	cmd := exec.Command(
		"protoc",
		"-I", filepath.Dir(protoPath),
		"--plugin=protoc-gen-go-lite="+plugin,
		"--go-lite_out="+outDir,
		"--go-lite_opt="+opt,
		protoPath,
	)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generate codegen mode fixture:\n%s", out)
	}

	content, err := os.ReadFile(filepath.Join(outDir, filepath.Base(strings.TrimSuffix(protoPath, ".proto")+".pb.go")))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
