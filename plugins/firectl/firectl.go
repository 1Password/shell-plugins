package firectl

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func FirectlCLI() schema.Executable {
	return schema.Executable{
		Name:    "Firectl CLI",
		Runs:    []string{"firectl"},
		DocsURL: sdk.URL("https://docs.fireworks.ai/tools-sdks/firectl/firectl"),
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
