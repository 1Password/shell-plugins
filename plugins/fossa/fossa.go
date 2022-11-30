package fossa

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func FOSSACLI() schema.Executable {
	return schema.Executable{
		Runs:      []string{"fossa"},
		Name:      "FOSSA CLI",
		DocsURL:   sdk.URL("https://github.com/fossas/fossa-cli"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{{
			Name: credname.APIKey,
		}},
	}
}
