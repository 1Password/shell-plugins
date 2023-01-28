package npm

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func NPMCLI() schema.Executable {
	return schema.Executable{
		Name:      "NPM install/publish",
		Runs:      []string{"npm"},
		DocsURL:   sdk.URL("https://docs.npmjs.com/cli"),
		NeedsAuth: needsauth.IfAny(
			needsauth.WhenContainsArgs("i"),
			needsauth.WhenContainsArgs("install"),
			needsauth.WhenContainsArgs("ci"),
			needsauth.WhenContainsArgs("clean-install"),
			needsauth.WhenContainsArgs("publish"),
			needsauth.WhenContainsArgs("unpublish"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
