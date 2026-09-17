package exercism

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

// Regression test for https://github.com/1Password/shell-plugins/issues/619. The Exercism
// subcommands that only touch local setup must not require an Exercism API Key item, both in
// their bare form and with the arguments the CLI accepts for them.
func TestExercismCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, ExercismCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"no for completion": {
			Args:              []string{"completion"},
			ExpectedNeedsAuth: false,
		},
		"no for completion with shell arg": {
			Args:              []string{"completion", "bash"},
			ExpectedNeedsAuth: false,
		},
		"no for upgrade": {
			Args:              []string{"upgrade"},
			ExpectedNeedsAuth: false,
		},
		"no for workspace": {
			Args:              []string{"workspace"},
			ExpectedNeedsAuth: false,
		},
		"no for troubleshoot": {
			Args:              []string{"troubleshoot"},
			ExpectedNeedsAuth: false,
		},
		"no for troubleshoot with --full-api-key": {
			Args:              []string{"troubleshoot", "--full-api-key"},
			ExpectedNeedsAuth: false,
		},
		"no for workspace with a flag": {
			Args:              []string{"workspace", "--verbose"},
			ExpectedNeedsAuth: false,
		},
		"yes for submit": {
			Args:              []string{"submit", "bob.go"},
			ExpectedNeedsAuth: true,
		},
		"yes for download": {
			Args:              []string{"download", "--exercise=bob", "--track=go"},
			ExpectedNeedsAuth: true,
		},
		"yes for test": {
			Args:              []string{"test"},
			ExpectedNeedsAuth: true,
		},
		// The exempt names are ordinary words that can also appear as exercise names, track
		// names, or filenames, so they must only be exempt in command position.
		"yes for an exercise named after an exempt subcommand": {
			Args:              []string{"submit", "completion"},
			ExpectedNeedsAuth: true,
		},
		"yes for a subcommand that merely starts with an exempt name": {
			Args:              []string{"upgrades"},
			ExpectedNeedsAuth: true,
		},
		"no without args": {
			Args:              []string{},
			ExpectedNeedsAuth: false,
		},
		"no for --help": {
			Args:              []string{"--help"},
			ExpectedNeedsAuth: false,
		},
		"no for version": {
			Args:              []string{"version"},
			ExpectedNeedsAuth: false,
		},
	})
}
