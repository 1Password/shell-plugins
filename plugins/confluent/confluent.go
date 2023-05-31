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
			needsauth.NotForHelpOrVersion(),
			needsauth.NotForExactArgs("local"),
			needsauth.NotForExactArgs("update"),
			needsauth.NotForExactArgs("prompt"),
			needsauth.NotForExactArgs("plugin"),
			needsauth.NotForExactArgs("logout"),
			needsauth.NotForExactArgs("context"),
			needsauth.NotForExactArgs("completion"),
			needsauth.NotForExactArgs("cloud-signup"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.UserLogin,
			},
		},
	}
}
