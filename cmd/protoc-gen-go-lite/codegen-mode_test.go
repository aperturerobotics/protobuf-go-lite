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
}
`

func TestCodegenModeDefaultUsesHelperSize(t *testing.T) {
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
}

func TestCodegenModeUnrolledUsesPreviousSizeShape(t *testing.T) {
	out := generateCodegenModeFixture(t, "codegen=unrolled")
	if !strings.Contains(out, "protobuf_go_lite.SizeOfZigzag") {
		t.Fatalf("unrolled output missing SizeOfZigzag:\n%s", out)
	}
	if strings.Contains(out, "protobuf_go_lite.SizeZigzagNonZero") {
		t.Fatalf("unrolled output should not contain helper SizeZigzagNonZero:\n%s", out)
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

	opt := "features=size,paths=source_relative"
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
