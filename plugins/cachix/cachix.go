package cachix

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func CachixCLI() schema.Executable {
	return schema.Executable{
		Name:      "Cachix CLI",
		Runs:      []string{"cachix"},
		DocsURL:   sdk.URL("https://docs.cachix.org"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AuthToken,
			},
		},
	}
}
