package oci

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OCICLI() schema.Executable {
	return schema.Executable{
		Name:    "oci",
		Runs:    []string{"oci"},
		DocsURL: sdk.URL("https://docs.oracle.com/en-us/iaas/tools/oci-cli/latest/oci_cli_docs/"),
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
