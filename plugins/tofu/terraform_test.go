package tofu

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

// The project credentials are optional, so when this rule declines the command
// still runs — just without the credentials the workspace needs.
func TestTofuCLIProjectCredentialsNeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, TofuCLI().Uses[0].NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"workspace select needs the project credentials": {
			Args:              []string{"workspace", "select", "staging"},
			ExpectedNeedsAuth: true,
		},
		"workspace list needs the project credentials": {
			Args:              []string{"workspace", "list"},
			ExpectedNeedsAuth: true,
		},
		"workspace new needs the project credentials": {
			Args:              []string{"workspace", "new", "staging"},
			ExpectedNeedsAuth: true,
		},
		"init needs the project credentials": {
			Args:              []string{"init"},
			ExpectedNeedsAuth: true,
		},
		"plan needs the project credentials": {
			Args:              []string{"plan"},
			ExpectedNeedsAuth: true,
		},
		"plan with flags needs the project credentials": {
			Args:              []string{"plan", "-out=tfplan"},
			ExpectedNeedsAuth: true,
		},
		"apply needs the project credentials": {
			Args:              []string{"apply"},
			ExpectedNeedsAuth: true,
		},
		"destroy needs the project credentials": {
			Args:              []string{"destroy"},
			ExpectedNeedsAuth: true,
		},
		"refresh needs the project credentials": {
			Args:              []string{"refresh"},
			ExpectedNeedsAuth: true,
		},
		"import needs the project credentials": {
			Args:              []string{"import"},
			ExpectedNeedsAuth: true,
		},
		"test needs the project credentials": {
			Args:              []string{"test"},
			ExpectedNeedsAuth: true,
		},
		"state subcommands need the project credentials": {
			Args:              []string{"state", "list"},
			ExpectedNeedsAuth: true,
		},
		"fmt is local and needs no credentials": {
			Args:              []string{"fmt"},
			ExpectedNeedsAuth: false,
		},
		"validate is local and needs no credentials": {
			Args:              []string{"validate"},
			ExpectedNeedsAuth: false,
		},
	})
}

// The executable-level rule gates whether the plugin engages at all. It admits
// every workspace invocation, so it is not what decides credential injection.
func TestTofuCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, TofuCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"engages for workspace select": {
			Args:              []string{"workspace", "select", "staging"},
			ExpectedNeedsAuth: true,
		},
		"skips auth for help": {
			Args:              []string{"--help"},
			ExpectedNeedsAuth: false,
		},
		"skips auth for version": {
			Args:              []string{"--version"},
			ExpectedNeedsAuth: false,
		},
		"skips auth without args": {
			Args:              []string{},
			ExpectedNeedsAuth: false,
		},
	})
}
