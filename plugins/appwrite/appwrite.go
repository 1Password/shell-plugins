package appwrite

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AppwriteCLI() schema.Executable {
	return schema.Executable{
		Name:    "Appwrite CLI",
		Runs:    []string{"appwrite"},
		DocsURL: sdk.URL("https://appwrite.io/docs/command-line-ci"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("client"),
			needsauth.NotWhenContainsArgs("login"),
			needsauth.NotWhenContainsArgs("logout"),
			needsauth.NotForExactArgs("deploy"),
			needsauth.NotForExactArgs("projects"),
			needsauth.NotForExactArgs("storage"),
			needsauth.NotForExactArgs("teams"),
			needsauth.NotForExactArgs("users"),
			needsauth.NotForExactArgs("account"),
			needsauth.NotForExactArgs("avatars"),
			needsauth.NotForExactArgs("functions"),
			needsauth.NotForExactArgs("databases"),
			needsauth.NotForExactArgs("health"),
			needsauth.NotForExactArgs("locale"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
