package render

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func RenderCLI() schema.Executable {
	return schema.Executable{
		Name:    "Render CLI",
		Runs:    []string{"render"},
		DocsURL: sdk.URL("https://render.com/docs/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("config"),
			needsauth.NotWhenContainsArgs("jobs"),
			needsauth.NotWhenContainsArgs("deploys"),
			needsauth.NotWhenContainsArgs("repo"),
			needsauth.NotWhenContainsArgs("services"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
