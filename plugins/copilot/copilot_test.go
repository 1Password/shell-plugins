package copilot

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestCopilotCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, CopilotCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"interactive session": {
			Args:              []string{},
			ExpectedNeedsAuth: true,
		},
		"agent invocation": {
			Args:              []string{"--agent", "blub"},
			ExpectedNeedsAuth: true,
		},
		"login": {
			Args:              []string{"login"},
			ExpectedNeedsAuth: false,
		},
		"version": {
			Args:              []string{"version"},
			ExpectedNeedsAuth: false,
		},
		"completion": {
			Args:              []string{"completion", "bash"},
			ExpectedNeedsAuth: false,
		},
	})
}
