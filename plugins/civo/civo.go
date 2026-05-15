package civo

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func CivoCLI() schema.Executable {
	return schema.Executable{
		Name:    "Civo CLI",
		Runs:    []string{"civo"},
		DocsURL: sdk.URL("https://www.civo.com/docs/overview/civo-cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotForExactArgs("config"),
			needsauth.NotWhenContainsArgs("apikey", "save"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
