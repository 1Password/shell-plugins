package oci

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func TestOCICLI(t *testing.T) {
	executable := OCICLI()

	if executable.Name != "oci" {
		t.Errorf("Name = %q, want %q", executable.Name, "oci")
	}
	if len(executable.Runs) != 1 || executable.Runs[0] != "oci" {
		t.Errorf("Runs = %v, want [oci]", executable.Runs)
	}
	if len(executable.Uses) != 1 || executable.Uses[0].Name != credname.APIKey {
		t.Errorf("Uses = %v, want credential %q", executable.Uses, credname.APIKey)
	}
}

func TestOCICLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, OCICLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"without arguments": {
			ExpectedNeedsAuth: false,
		},
		"help": {
			Args:              []string{"--help"},
			ExpectedNeedsAuth: false,
		},
		"version": {
			Args:              []string{"--version"},
			ExpectedNeedsAuth: false,
		},
		"command": {
			Args:              []string{"os", "ns", "get"},
			ExpectedNeedsAuth: true,
		},
	})
}
