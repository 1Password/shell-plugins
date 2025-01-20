package cvelib

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func CVEServicesAPICLI() schema.Executable {
	return schema.Executable{
		Name:    "CVE Services API CLI",
		Runs:    []string{"cve"},
		DocsURL: sdk.URL("https://github.com/RedHatProductSecurity/cvelib"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
