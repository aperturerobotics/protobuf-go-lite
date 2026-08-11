package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const registryAnnotationsProto = `syntax = "proto3";

package annotations;

import "google/protobuf/descriptor.proto";

option go_package = "registryfixture;registryfixture";

extend google.protobuf.MessageOptions {
  MsgOpts msg_opts = 50000;
}

enum Region {
  REGION_UNSPECIFIED = 0;
  REGION_US_EAST_2 = 1;
}

message Lifecycle {
  Region region = 1;
  repeated string tags = 2;
}

message MsgOpts {
  string tenant = 1;
  string table = 2;
  repeated string tenants = 3;
  Lifecycle lifecycle = 4;
  float weight = 5;
}
`

const registryEventsProto = `syntax = "proto3";

package demo;

import "annotations.proto";

option go_package = "registryfixture;registryfixture";

message ClickEvents {
  option (annotations.msg_opts) = {
    tenant: "BLINKIT"
    table: "click_events"
    tenants: "BLINKIT"
    tenants: "ZOMATO"
    lifecycle: {
      region: REGION_US_EAST_2
      tags: "a"
      tags: "b"
    }
    weight: 0.1
  };

  string id = 1;
}

message NestedHolder {
  message Inner {
    option (annotations.msg_opts).tenant = "ZOMATO";
    option (annotations.msg_opts).table = "inner";
    int32 value = 1;
  }
  Inner inner = 1;
}
`

func TestRegistryOptInGeneratesInitAndOptions(t *testing.T) {
	root := repoRoot(t)
	plugin := buildCurrentPlugin(t, root)
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "annotations.proto"), registryAnnotationsProto)
	writeFile(t, filepath.Join(dir, "events.proto"), registryEventsProto)
	outDir := filepath.Join(dir, "out")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(
		"protoc",
		"-I", dir,
		"-I", protobufSourceDir(t, root),
		"--plugin=protoc-gen-go-lite="+plugin,
		"--go-lite_out="+outDir,
		"--go-lite_opt=features=size+marshal+unmarshal,paths=source_relative,registry=true",
		"annotations.proto",
		"events.proto",
	)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generate registry fixture:\n%s", out)
	}

	eventsOut := string(readFile(t, filepath.Join(outDir, "events.pb.go")))
	assertContainsAll(t, eventsOut, "registry output", []string{
		`registry.Register(registry.Entry{`,
		`FullName: "demo.ClickEvents"`,
		`type.googleapis.com/demo.ClickEvents`,
		`annotations.msg_opts.tenant`,
		`"BLINKIT"`,
		`annotations.msg_opts.table`,
		`"click_events"`,
		`annotations.msg_opts.tenants`,
		`"ZOMATO"`,
		`annotations.msg_opts.lifecycle.region`,
		`"REGION_US_EAST_2"`,
		`annotations.msg_opts.lifecycle.tags`,
		`"a"`,
		`"b"`,
		`annotations.msg_opts.weight`,
		`FullName: "demo.NestedHolder"`,
		`FullName: "demo.NestedHolder.Inner"`,
		`return new(ClickEvents)`,
	})

	writeFile(t, filepath.Join(outDir, "go.mod"), "module registryfixture\n\ngo 1.23\n\nrequire github.com/aperturerobotics/protobuf-go-lite v0.0.0\n\nreplace github.com/aperturerobotics/protobuf-go-lite => "+root+"\n")
	writeFile(t, filepath.Join(outDir, "registry_runtime_test.go"), `package registryfixture

import (
	"testing"

	"github.com/aperturerobotics/protobuf-go-lite/registry"
)

func TestRuntimeRegistry(t *testing.T) {
	msg, ok := registry.NewByName("demo.ClickEvents")
	if !ok {
		t.Fatal("missing ClickEvents")
	}
	if _, ok := msg.(*ClickEvents); !ok {
		t.Fatalf("got %T", msg)
	}

	var weight string
	byKey := map[string]registry.Entry{}
	registry.Range(func(e registry.Entry) bool {
		tenant, tenantOK := e.Option("annotations.msg_opts.tenant")
		table, tableOK := e.Option("annotations.msg_opts.table")
		if tenantOK && tableOK {
			byKey[tenant+"."+table] = e
		}
		if e.FullName == "demo.ClickEvents" {
			weight, _ = e.Option("annotations.msg_opts.weight")
		}
		return true
	})
	if weight != "0.1" {
		t.Fatalf("weight=%q, want 0.1", weight)
	}
	if _, ok := byKey["BLINKIT.click_events"]; !ok {
		t.Fatalf("missing tenant.table key: %#v", byKey)
	}
	if _, ok := byKey["ZOMATO.inner"]; !ok {
		t.Fatalf("missing nested key: %#v", byKey)
	}
}
`)

	testCmd := exec.Command("go", "test", "-mod=mod", "./...")
	testCmd.Dir = outDir
	testOut, err := testCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated registry package should compile and pass:\n%s", testOut)
	}
}

func TestRegistryDisabledByDefault(t *testing.T) {
	root := repoRoot(t)
	plugin := buildCurrentPlugin(t, root)
	protoPath := writeTempProto(t, `syntax = "proto3";
package noreg;
option go_package = "noreg;noreg";
message Msg { int32 v = 1; }
`)
	outDir := t.TempDir()
	cmd := exec.Command(
		"protoc",
		"-I", filepath.Dir(protoPath),
		"--plugin=protoc-gen-go-lite="+plugin,
		"--go-lite_out="+outDir,
		"--go-lite_opt=features=size+marshal+unmarshal,paths=source_relative",
		protoPath,
	)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generate:\n%s", out)
	}
	content := string(readFile(t, filepath.Join(outDir, filepath.Base(strings.TrimSuffix(protoPath, ".proto")+".pb.go"))))
	if strings.Contains(content, "registry.Register") {
		t.Fatal("registry should be opt-in")
	}
}

func TestRegistryRequiresFeatures(t *testing.T) {
	root := repoRoot(t)
	plugin := buildCurrentPlugin(t, root)
	protoPath := writeTempProto(t, `syntax = "proto3";
package badreg;
option go_package = "badreg;badreg";
message Msg { int32 v = 1; }
`)
	outDir := t.TempDir()
	cmd := exec.Command(
		"protoc",
		"-I", filepath.Dir(protoPath),
		"--plugin=protoc-gen-go-lite="+plugin,
		"--go-lite_out="+outDir,
		"--go-lite_opt=features=size+marshal,paths=source_relative,registry=true",
		protoPath,
	)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected registry feature validation failure")
	}
	if !strings.Contains(string(out), "registry=true requires size, marshal, and unmarshal features") {
		t.Fatalf("unexpected error:\n%s", out)
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
