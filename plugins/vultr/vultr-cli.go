package vultr

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func VultrCLI() schema.Executable {
	return schema.Executable{
		Name:      "Vultr CLI",
		Runs:      []string{"vultr-cli"},
		DocsURL:   sdk.URL("https://github.com/vultr/vultr-cli"),
		NeedsAuth: needsauth.NotWhenContainsArgs("--config"),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
