package specs

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTypeSupportJsonMarshalUnmarshal(t *testing.T) {
	b, err := json.Marshal(TypeSupportFull)
	if err != nil {
		t.Fatal("failed to marshal typeSupport:", err)
	}
	var typeSupport TypeSupport
	if err := json.Unmarshal(b, &typeSupport); err != nil {
		t.Fatal("failed to unmarshal typeSupport:", err)
	}
	if typeSupport != TypeSupportFull {
		t.Fatal("expected typeSupport to be full, but got:", typeSupport)
	}
}

func TestTypeSupportYamlMarshalUnmarsahl(t *testing.T) {
	b, err := yaml.Marshal(TypeSupportFull)
	if err != nil {
		t.Fatal("failed to marshal typeSupport:", err)
	}
	var typeSupport TypeSupport
	if err := yaml.Unmarshal(b, &typeSupport); err != nil {
		t.Fatal("failed to unmarshal typeSupport:", err)
	}
	if typeSupport != TypeSupportFull {
		t.Fatal("expected typeSupport to be full, but got:", typeSupport)
	}
}
