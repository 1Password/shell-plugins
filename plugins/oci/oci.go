package oci

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OracleCloudCLI() schema.Executable {
	return schema.Executable{
		Name:      "Oracle Cloud CLI",
		Runs:      []string{"oci"},
		DocsURL:   sdk.URL("https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessKey,
			},
		},
	}
}
