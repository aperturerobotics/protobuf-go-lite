package generator

import "testing"

func TestValidateRegistryFeaturesAllRequiresRegisteredFeatures(t *testing.T) {
	features := defaultFeatures
	defaultFeatures = map[string]Feature{
		"marshal": nil,
		"size":    nil,
	}
	t.Cleanup(func() { defaultFeatures = features })

	if err := validateRegistryFeatures([]string{"all"}); err == nil {
		t.Fatal("registry accepted all without a registered unmarshal feature")
	}
}
