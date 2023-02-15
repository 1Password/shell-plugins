package confluent

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ConfluentCLI() schema.Executable {
	return schema.Executable{
		Name:    "Confluent CLI",
		Runs:    []string{"confluent"},
		DocsURL: sdk.URL("https://docs.confluent.io/confluent-cli/current/overview.html"),
		NeedsAuth: needsauth.IfAll(
			needsauth.ForCommand("login"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Credentials,
			},
		},
	}
}
