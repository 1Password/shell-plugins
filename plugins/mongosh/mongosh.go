package mongosh

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func MongoDBShellCLI() schema.Executable {
	return schema.Executable{
		Name:      "MongoDB Shell CLI",
		Runs:      []string{"mongosh"},
		DocsURL:   sdk.URL("https://mongosh.com/docs/cli"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DatabaseCredentials,
			},
		},
	}
}
