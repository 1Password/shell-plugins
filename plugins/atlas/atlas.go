package atlas

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func MongoDBAtlasCLI() schema.Executable {
	return schema.Executable{
		Name:    "MongoDB Atlas CLI",
		Runs:    []string{"atlas"},
		DocsURL: sdk.URL("https://www.mongodb.com/docs/atlas/cli/stable/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
