package circleci

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func CircleCICLI() schema.Executable {
	return schema.Executable{
		Name:    "CircleCI CLI",
		Runs:    []string{"circleci"},
		DocsURL: sdk.URL("https://circleci.com/docs/local-cli/"),
		NeedsAuth: needsauth.For(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotForArgs("config"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAPIToken,
			},
		},
	}
}
