package upstash

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func UpstashCLI() schema.Executable {
	return schema.Executable{
		Name:    "Upstash CLI",
		Runs:    []string{"upstash"},
		DocsURL: sdk.URL("https://github.com/upstash/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("--config"),
			needsauth.NotWhenContainsArgs("-c"),
			needsauth.NotWhenContainsArgs("auth"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
