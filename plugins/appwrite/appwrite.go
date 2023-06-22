package appwrite

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AppwriteCLI() schema.Executable {
	return schema.Executable{
		Name:    "Appwrite CLI", // TODO: Check if this is correct
		Runs:    []string{"appwrite"},
		DocsURL: sdk.URL("https://appwrite.io/docs/"), // TODO: Replace with actual URL
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
