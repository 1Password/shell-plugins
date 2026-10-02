package argocd

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestArgocdCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, ArgocdCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"yes for app list": {
			Args:              []string{"app", "list"},
			ExpectedNeedsAuth: true,
		},
		"no for login": {
			Args:              []string{"login", "localhost:8080"},
			ExpectedNeedsAuth: false,
		},
		"no for relogin": {
			Args:              []string{"relogin"},
			ExpectedNeedsAuth: false,
		},
	})
}
