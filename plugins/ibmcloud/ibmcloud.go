package ibmcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func IBMCloudCLI() schema.Executable {
	return schema.Executable{
		Name:    "IBM Cloud CLI",
		Runs:    []string{"ibmcloud"},
		DocsURL: sdk.URL("https://cloud.ibm.com/docs/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("-u"),
			needsauth.NotWhenContainsArgs("-p"),
			needsauth.NotWhenContainsArgs("--apikey"),
			needsauth.NotWhenContainsArgs("--cr-token"),
			needsauth.NotWhenContainsArgs("--vpc-cri"),
			needsauth.NotWhenContainsArgs("--profile"),
			needsauth.NotWhenContainsArgs("--sso"),
			needsauth.NotWhenContainsArgs("-c"),
			needsauth.NotWhenContainsArgs("--no-account"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
