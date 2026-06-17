package kiro

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func KiroCLI() schema.Executable {
	return schema.Executable{
		Name:      "Kiro CLI",
		Runs:      []string{"kiro-cli"},
		DocsURL:   sdk.URL("https://kiro.dev/docs/cli/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("login"),
			needsauth.NotWhenContainsArgs("logout"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
