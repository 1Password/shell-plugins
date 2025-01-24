package motherduck

import (
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func NotWhenAnyArgsContain(argsSequence ...string) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		if len(argsSequence) == 0 {
			return true
		}

		if len(argsSequence) > len(in.CommandArgs) {
			return true
		}

		for i := range in.CommandArgs {
			if i+len(argsSequence) > len(in.CommandArgs) {
				return true
			}

			matches := true
			for i, argsToCompare := range in.CommandArgs[i : i+len(argsSequence)] {
				if !strings.Contains(argsToCompare, argsSequence[i]) {
					matches = false
				}
			}

			// If the argsToSkip are found in the command-line args, return that the command
			// does not not require authentication
			if matches {
				return false
			}
		}
		return true
	}
}

func DuckDBCLI() schema.Executable {
	return schema.Executable{
		Name:    "DuckDB CLI",
		Runs:    []string{"duckdb"},
		DocsURL: sdk.URL("https://duckdb.org/docs/api/cli/overview"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			NotWhenAnyArgsContain("motherduck_token="),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
