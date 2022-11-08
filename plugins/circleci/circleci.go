package circleci

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func CircleCICLI() schema.Executable {
	return schema.Executable{
		Runs:    []string{"circleci"},
		Name:    "CircleCI CLI",
		DocsURL: sdk.URL("https://circleci.com/docs/local-cli/"),
		NeedsAuth: needsauth.For(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotForArgs("config"),
		),
		Credentials: []schema.CredentialType{
			PersonalAPIToken(),
		},
	}
}
