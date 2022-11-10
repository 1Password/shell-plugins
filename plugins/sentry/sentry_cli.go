package sentry

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SentryCLI() schema.Executable {
	return schema.Executable{
		Runs:      []string{"sentry-cli"},
		Name:      "Sentry CLI",
		DocsURL:   sdk.URL("https://docs.sentry.io/product/cli/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		UsesCredentials: []schema.CredentialUsage{
			{
				Name: credname.AuthToken,
			},
		},
	}
}
