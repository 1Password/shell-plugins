package motherduck

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

// unsetTokenEnvVars keeps tokens from the developer's own shell from leaking into the tests.
func unsetTokenEnvVars(t *testing.T) {
	t.Setenv("motherduck_token", "")
	t.Setenv("MOTHERDUCK_TOKEN", "")
}

func TestDuckDBCLINeedsAuth(t *testing.T) {
	unsetTokenEnvVars(t)

	plugintest.TestNeedsAuth(t, DuckDBCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"no args opens an in-memory database": {
			Args:              []string{},
			ExpectedNeedsAuth: false,
		},
		"local database file": {
			Args:              []string{"local.ddb"},
			ExpectedNeedsAuth: false,
		},
		"local database file with a command": {
			Args:              []string{"local.ddb", "-c", "select 1"},
			ExpectedNeedsAuth: false,
		},
		"default MotherDuck database": {
			Args:              []string{"md:"},
			ExpectedNeedsAuth: true,
		},
		"named MotherDuck database": {
			Args:              []string{"md:my_db"},
			ExpectedNeedsAuth: true,
		},
		"MotherDuck database with a command": {
			Args:              []string{"md:my_db", "-c", "select 1"},
			ExpectedNeedsAuth: true,
		},
		"MotherDuck attached from a command": {
			Args:              []string{"-c", "ATTACH 'md:'"},
			ExpectedNeedsAuth: true,
		},
		"token passed in the connection string": {
			Args:              []string{"md:my_db?motherduck_token=abc"},
			ExpectedNeedsAuth: false,
		},
	})
}

func TestDuckDBCLINeedsAuthWithTokenEnvVar(t *testing.T) {
	unsetTokenEnvVars(t)
	t.Setenv("motherduck_token", "abc")

	plugintest.TestNeedsAuth(t, DuckDBCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"MotherDuck database": {
			Args:              []string{"md:my_db"},
			ExpectedNeedsAuth: false,
		},
	})
}
