package generator

import "testing"

func TestConfigSetCodegenMode(t *testing.T) {
	var cfg Config
	if err := cfg.SetCodegenMode(""); err != nil {
		t.Fatalf("default mode: %v", err)
	}
	if cfg.CodegenMode != CodegenModeHelper {
		t.Fatalf("default CodegenMode = %q, want %q", cfg.CodegenMode, CodegenModeHelper)
	}
	if !cfg.HelperCodegen() {
		t.Fatal("default mode should use helper codegen")
	}

	if err := cfg.SetCodegenMode(string(CodegenModeUnrolled)); err != nil {
		t.Fatalf("unrolled mode: %v", err)
	}
	if cfg.CodegenMode != CodegenModeUnrolled {
		t.Fatalf("CodegenMode = %q, want %q", cfg.CodegenMode, CodegenModeUnrolled)
	}
	if cfg.HelperCodegen() {
		t.Fatal("unrolled mode should not use helper codegen")
	}

	if err := cfg.SetCodegenMode(string(CodegenModeHelper)); err != nil {
		t.Fatalf("helper mode: %v", err)
	}
	if cfg.CodegenMode != CodegenModeHelper {
		t.Fatalf("CodegenMode = %q, want %q", cfg.CodegenMode, CodegenModeHelper)
	}
	if !cfg.HelperCodegen() {
		t.Fatal("helper mode should use helper codegen")
	}
}

func TestConfigSetCodegenModeRejectsUnknown(t *testing.T) {
	var cfg Config
	if err := cfg.SetCodegenMode("bogus"); err == nil {
		t.Fatal("expected unknown codegen mode to fail")
	}
}
