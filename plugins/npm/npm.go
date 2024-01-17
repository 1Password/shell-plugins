package npm

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func NPMCLI() schema.Executable {
	return schema.Executable{
		Name:    "NPM CLI",
		Runs:    []string{"npm"},
		DocsURL: sdk.URL("https://docs.npmjs.com/cli"),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
