package readme

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestRdmeNeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, ReadMeCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"yes for docs command": {
			Args:              []string{"docs"},
			ExpectedNeedsAuth: true,
		},
		"no for docs command help": {
			Args:              []string{"docs", "--help"},
			ExpectedNeedsAuth: false,
		},
		"no for open command": {
			Args:              []string{"open"},
			ExpectedNeedsAuth: false,
		},
		"no for version": {
			Args:              []string{"--version"},
			ExpectedNeedsAuth: false,
		},
		"yes for openapi with --version flag": {
			Args:              []string{"openapi", "--version", "2.0"},
			ExpectedNeedsAuth: true,
		},
	})
}
