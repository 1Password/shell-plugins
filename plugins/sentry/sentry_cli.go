package sentry

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SentryCLI() schema.Executable {
	return schema.Executable{
		Name:      "Sentry CLI",
		Runs:      []string{"sentry-cli"},
		DocsURL:   sdk.URL("https://docs.sentry.io/product/cli/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AuthToken,
			},
		},
	}
}
