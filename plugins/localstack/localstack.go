package localstack

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func LocalStackCLI() schema.Executable {
	return schema.Executable{
		Name:    "LocalStack CLI",
		Runs:    []string{"localstack"},
		DocsURL: sdk.URL("https://docs.localstack.cloud/getting-started/installation/#localstack-cli"),
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
