package localstack

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func LocalStackCLI() schema.Executable {
	return schema.Executable{
		Name:      "LocalStack CLI", // TODO: Check if this is correct
		Runs:      []string{"localstack"},
		DocsURL:   sdk.URL("https://localstack.com/docs/cli"), // TODO: Replace with actual URL
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
