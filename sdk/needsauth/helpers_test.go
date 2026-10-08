package needsauth

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestNoArg(t *testing.T) {
	plugintest.TestNeedsAuth(t, NotWithoutArgs(), map[string]plugintest.NeedsAuthCase{
		"yes with args": {
			Args:              []string{"foo"},
			ExpectedNeedsAuth: true,
		},
		"not without command": {
			Args:              []string{},
			ExpectedNeedsAuth: false,
		},
	})
}

func TestHelp(t *testing.T) {
	plugintest.TestNeedsAuth(t, NotForHelp(), map[string]plugintest.NeedsAuthCase{
		"no for exact help flag": {
			Args:              []string{"--help"},
			ExpectedNeedsAuth: false,
		},
		"no for exact help short flag": {
			Args:              []string{"-h"},
			ExpectedNeedsAuth: false,
		},
		"no for help flag in command": {
			Args:              []string{"foo", "--help"},
			ExpectedNeedsAuth: false,
		},
		"no for help short flag in command": {
			Args:              []string{"foo", "-h"},
			ExpectedNeedsAuth: false,
		},
		"yes without help": {
			Args:              []string{"foo"},
			ExpectedNeedsAuth: true,
		},
	})
}

func TestVersion(t *testing.T) {
	plugintest.TestNeedsAuth(t, NotForVersion(), map[string]plugintest.NeedsAuthCase{
		"not for exact version flag": {
			Args:              []string{"--version"},
			ExpectedNeedsAuth: false,
		},
		"not for exact version short flag": {
			Args:              []string{"-v"},
			ExpectedNeedsAuth: false,
		},
		"yes when --version flag is part of a command": {
			Args:              []string{"deploy", "--version", "1.0.0"},
			ExpectedNeedsAuth: true,
		},
		"yes when -v is part of a command": {
			Args:              []string{"run", "-v", "$PWD:/app", "foo"},
			ExpectedNeedsAuth: true,
		},
	})
}

func TestContainsArgs(t *testing.T) {
	plugintest.TestNeedsAuth(t, NotWhenContainsArgs("--mode", "dry-run"), map[string]plugintest.NeedsAuthCase{
		"yes by default": {
			Args:              []string{"deploy"},
			ExpectedNeedsAuth: true,
		},
		"yes when only one of the args is present": {
			Args:              []string{"deploy", "--mode", "live"},
			ExpectedNeedsAuth: true,
		},
		"yes when both args are present, but not in sequence": {
			Args:              []string{"deploy", "--mode", "live", "--app-name", "dry-run"},
			ExpectedNeedsAuth: true,
		},
		"no when all args are present in sequence": {
			Args:              []string{"deploy", "--mode", "dry-run"},
			ExpectedNeedsAuth: false,
		},
	})
}

func TestForCommand(t *testing.T) {
	plugintest.TestNeedsAuth(t, ForCommand("config"), map[string]plugintest.NeedsAuthCase{
		"no by default": {
			Args:              []string{"deploy"},
			ExpectedNeedsAuth: false,
		},
		"no without args": {
			Args:              []string{},
			ExpectedNeedsAuth: false,
		},
		"yes for the bare command": {
			Args:              []string{"config"},
			ExpectedNeedsAuth: true,
		},
		"yes for a subcommand of the command": {
			Args:              []string{"config", "get", "token"},
			ExpectedNeedsAuth: true,
		},
		"yes for the command with flags": {
			Args:              []string{"config", "--global"},
			ExpectedNeedsAuth: true,
		},
		"no when the command appears after another command": {
			Args:              []string{"deploy", "config"},
			ExpectedNeedsAuth: false,
		},
		"no for a command that merely starts with the same characters": {
			Args:              []string{"configure"},
			ExpectedNeedsAuth: false,
		},
	})

	plugintest.TestNeedsAuth(t, ForCommand(), map[string]plugintest.NeedsAuthCase{
		"no without a command to match": {
			Args:              []string{"config"},
			ExpectedNeedsAuth: false,
		},
	})
}

func TestForNestedCommand(t *testing.T) {
	plugintest.TestNeedsAuth(t, ForCommand("shell-completions", "install"), map[string]plugintest.NeedsAuthCase{
		"yes for the exact nested command": {
			Args:              []string{"shell-completions", "install"},
			ExpectedNeedsAuth: true,
		},
		"yes for the nested command with args": {
			Args:              []string{"shell-completions", "install", "--shell", "bash"},
			ExpectedNeedsAuth: true,
		},
		"no for only the first part of the nested command": {
			Args:              []string{"shell-completions"},
			ExpectedNeedsAuth: false,
		},
		"no for the nested command's parts in the wrong order": {
			Args:              []string{"install", "shell-completions"},
			ExpectedNeedsAuth: false,
		},
	})
}

// NotForCommand negates ForCommand, so the prefix matching itself is covered above. What is
// worth pinning here is the inversion, and that NotForCommand() opts out of nothing, where
// NotForExactArgs() matches the empty arg list and is aliased as NotWithoutArgs().
func TestNotForCommand(t *testing.T) {
	plugintest.TestNeedsAuth(t, NotForCommand("config"), map[string]plugintest.NeedsAuthCase{
		"no for the command and anything nested under it": {
			Args:              []string{"config", "get", "token"},
			ExpectedNeedsAuth: false,
		},
		"yes when the command is not in command position": {
			Args:              []string{"deploy", "config"},
			ExpectedNeedsAuth: true,
		},
	})

	plugintest.TestNeedsAuth(t, NotForCommand(), map[string]plugintest.NeedsAuthCase{
		"yes without a command to opt out of": {
			Args:              []string{"config"},
			ExpectedNeedsAuth: true,
		},
	})
}

func TestComplexChain(t *testing.T) {
	// Example of a fictitious package manager that requires authentication for:
	// * The "publish" command, unless the "--dry-run" flag is present
	// * The "install" command
	// * The "auth" subcommands
	// * Not for "--version" and "--help", or when no args are specified
	// * Never when the "--local" flag is specified

	needsAuth := /*needsauth.*/ IfAll(
		/*needsauth.*/ NotWithoutArgs(),
		/*needsauth.*/ NotForHelpOrVersion(),
		/*needsauth.*/ NotWhenContainsArgs("--local"),
		/*needsauth.*/ IfAny(
			/*needsauth.*/ ForCommand("auth"),
			/*needsauth.*/ ForCommand("install"),
			/*needsauth.*/ IfAll(
				/*needsauth.*/ ForCommand("publish"),
				/*needsauth.*/ NotWhenContainsArgs("--dry-run"),
			),
		),
	)

	plugintest.TestNeedsAuth(t, needsAuth, map[string]plugintest.NeedsAuthCase{
		"no without args": {
			Args:              []string{},
			ExpectedNeedsAuth: false,
		},
		"no for help": {
			Args:              []string{"--help"},
			ExpectedNeedsAuth: false,
		},
		"yes for publish command": {
			Args:              []string{"publish", "my-app"},
			ExpectedNeedsAuth: true,
		},
		"no for publish help": {
			Args:              []string{"publish", "--help"},
			ExpectedNeedsAuth: false,
		},
		"no for publish dry run": {
			Args:              []string{"publish", "my-app", "--dry-run"},
			ExpectedNeedsAuth: false,
		},
		"no for publish local": {
			Args:              []string{"publish", "--local"},
			ExpectedNeedsAuth: false,
		},
		"yes for install command with args": {
			Args:              []string{"install"},
			ExpectedNeedsAuth: true,
		},
		"no for local install command": {
			Args:              []string{"install", "--local"},
			ExpectedNeedsAuth: false,
		},
		"yes for install command with a specific app version": {
			Args:              []string{"install", "--version", "1.0.0"},
			ExpectedNeedsAuth: true,
		},
		"no for version command": {
			Args:              []string{"--version"},
			ExpectedNeedsAuth: false,
		},
		"yes for auth command": {
			Args:              []string{"auth"},
			ExpectedNeedsAuth: true,
		},
		"yes for auth subcommand": {
			Args:              []string{"auth", "whoami"},
			ExpectedNeedsAuth: true,
		},
		"no for auth help": {
			Args:              []string{"auth", "--help"},
			ExpectedNeedsAuth: false,
		},
		"no for config command": {
			Args:              []string{"config"},
			ExpectedNeedsAuth: false,
		},
		"no for other install command": {
			Args:              []string{"shell-completions", "install", "--shell", "bash"},
			ExpectedNeedsAuth: false,
		},
	})
}
