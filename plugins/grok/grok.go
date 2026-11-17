package grok

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func GrokCLI() schema.Executable {
	return schema.Executable{
		Name:    "Grok CLI",
		Runs:    []string{"grok"},
		DocsURL: sdk.URL("https://docs.x.ai/build/cli/reference"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWhenContainsArgs("completions"),
			needsauth.NotForExactArgs("login"),
			needsauth.NotForExactArgs("logout"),
			needsauth.NotForExactArgs("update"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
