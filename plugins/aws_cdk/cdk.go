package aws_cdk

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AWSCDKCLI() schema.Executable {
	return schema.Executable{
		Name:      "AWS CDK CLI", // TODO: Check if this is correct
		Runs:      []string{"cdk"},
		DocsURL:   sdk.URL("https://aws_cdk.com/docs/cli"), // TODO: Replace with actual URL
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
