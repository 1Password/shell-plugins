package buildkite

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func BuildkiteCLI() schema.Executable {
	return schema.Executable{
		Name:      "Buildkite CLI",
		Runs:      []string{"bk"},
		DocsURL:   sdk.URL("https://buildkite.com/docs/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotForExactArgs("configure"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
