package sops

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SOPSCLI() schema.Executable {
	return schema.Executable{
		Name:    "SOPS CLI",
		Runs:    []string{"sops"},
		DocsURL: sdk.URL("https://github.com/getsops/sops"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.SecretKey,
			},
		},
	}
}
