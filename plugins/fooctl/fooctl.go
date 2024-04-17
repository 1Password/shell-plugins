package fooctl

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func FooctlCLI() schema.Executable {
	return schema.Executable{
		Name:    "Fooctl CLI",
		Runs:    []string{"fooctl"},
		DocsURL: sdk.URL("https://localhost/fooctl-cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			// needsauth.NotWhenContainsArgs("configure"),
			// needsauth.NotWhenContainsArgs("daemon"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Credentials,
			},
		},
	}
}
