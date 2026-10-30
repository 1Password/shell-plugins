package bundler

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func BundleCLI() schema.Executable {
	return schema.Executable{
		Name:    "Bundler CLI",
		Runs:    []string{"bundle"},
		DocsURL: sdk.URL("https://bundler.io/man/bundle-install.1.html"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.ForCommand("install"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
