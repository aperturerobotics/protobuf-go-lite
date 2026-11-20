package generator

import "fmt"

// Config holds the plugin parameters that shape a generator run.
type Config struct {
	// CodegenMode chooses shared helpers or inlined code in generated methods.
	CodegenMode CodegenMode
	// AllowEmpty writes a file even when it has no content.
	AllowEmpty bool
	// BuildTag is the go:build constraint put on each generated file.
	BuildTag string
	// Registry generates init-time message constructor and option registration.
	Registry bool
}

// CodegenMode names how generated methods encode, decode, clone and compare
// fields.
type CodegenMode string

const (
	// CodegenModeHelper makes generated methods call the runtime's shared
	// helpers. It is the default and produces the smallest code.
	CodegenModeHelper CodegenMode = "helper"
	// CodegenModeUnrolled inlines the helpers' logic into each method.
	CodegenModeUnrolled CodegenMode = "unrolled"
)

// SetCodegenMode sets CodegenMode from a plugin parameter value. An empty
// value selects CodegenModeHelper; an unknown one is an error.
func (c *Config) SetCodegenMode(mode string) error {
	switch CodegenMode(mode) {
	case "", CodegenModeHelper:
		c.CodegenMode = CodegenModeHelper
		return nil
	case CodegenModeUnrolled:
		c.CodegenMode = CodegenModeUnrolled
		return nil
	default:
		return fmt.Errorf("unknown codegen mode: %q", mode)
	}
}

// HelperCodegen reports whether generated methods call the shared helpers.
// A nil Config uses the default helper mode.
func (c *Config) HelperCodegen() bool {
	return c == nil || c.CodegenMode != CodegenModeUnrolled
}
