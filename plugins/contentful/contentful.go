package contentful

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ContentfulCLI() schema.Executable {
	return schema.Executable{
		Name:    "Contentful CLI",
		Runs:    []string{"contentful"},
		DocsURL: sdk.URL("https://www.contentful.com/developers/docs/tutorials/cli/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
