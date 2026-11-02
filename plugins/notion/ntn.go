package notion

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func NotionCLI() schema.Executable {
	return schema.Executable{
		Name:    "Notion CLI",
		Runs:    []string{"ntn"},
		DocsURL: sdk.URL("https://developers.notion.com/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotForCommand("login"),
			needsauth.NotForCommand("logout"),
			needsauth.NotForCommand("update"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
