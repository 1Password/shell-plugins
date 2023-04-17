package zendesk

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ZendeskCLI() schema.Executable {
	return schema.Executable{
		Name:    "Zendesk CLI",
		Runs:    []string{"zcli"},
		DocsURL: sdk.URL("https://developer.zendesk.com/documentation/apps/getting-started/using-zcli/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("profiles"),
			needsauth.NotWhenContainsArgs("login"),
			needsauth.NotWhenContainsArgs("logout"),
			needsauth.NotWhenContainsArgs("autocomplete"),
			needsauth.NotWhenContainsArgs("apps:bump"),
			needsauth.NotWhenContainsArgs("apps:clean"),
			needsauth.NotWhenContainsArgs("apps:new"),
			needsauth.NotWhenContainsArgs("apps:server"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
