package motherduck

import (
	"os"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

// The plugin is only invoked if:
//  - environment variable motherduck_token is not set
//  - connection string contains 'md:' and does not contain 'motherduck_token='
func ForMotherDuckButTokenNotSet() sdk.NeedsAuthentication {
    return func(in sdk.NeedsAuthenticationInput) bool {
        // If environment variables are already set, we don't need to authenticate
        if envValue := os.Getenv("motherduck_token"); envValue != "" {
            return false
        }
        
        // Otherwise, check if the command uses MotherDuck
        if len(in.CommandArgs) == 0 {
            return false
        }

        for _, arg := range in.CommandArgs {
            if strings.Contains(arg, "md:") && !strings.Contains(arg, "motherduck_token=") {
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
