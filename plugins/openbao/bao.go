package openbao

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OpenBaoCLI() schema.Executable {
	return schema.Executable{
		Name:      "OpenBao CLI", // TODO: Check if this is correct
		Runs:      []string{"bao"},
		DocsURL:   sdk.URL("https://openbao.com/docs/cli"), // TODO: Replace with actual URL
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AuthToken,
			},
		},
	}
}
