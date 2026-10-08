package motherduck

import (
	"os"
	"regexp"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

// motherDuckConnection matches MotherDuck connection strings like 'md:', 'md:my_db' or 'motherduck:my_db'.
// DuckDB treats the prefix case-insensitively. The prefix must not follow a word character,
// so text like 'cmd:' inside a query doesn't count.
var motherDuckConnection = regexp.MustCompile(`(?i)(^|[^[:alnum:]_])(md|motherduck):`)

// The plugin is only invoked if:
//   - neither the motherduck_token nor the MOTHERDUCK_TOKEN environment variable is set
//   - an argument contains an 'md:' or 'motherduck:' connection string that does not contain 'motherduck_token='
func ForMotherDuckButTokenNotSet() sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		// If a token is already set in the environment, we don't need to authenticate
		for _, envVar := range tokenEnvVars {
			if os.Getenv(envVar) != "" {
				return false
			}
		}

		// Otherwise, check if the command uses MotherDuck
		for _, arg := range in.CommandArgs {
			if motherDuckConnection.MatchString(arg) && !strings.Contains(arg, "motherduck_token=") {
				return true
			}
		}
		return false
	}
}

func DuckDBCLI() schema.Executable {
	return schema.Executable{
		Name:    "DuckDB CLI",
		Runs:    []string{"duckdb"},
		DocsURL: sdk.URL("https://duckdb.org/docs/api/cli/overview"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			ForMotherDuckButTokenNotSet(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
