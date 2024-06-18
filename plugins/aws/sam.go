package aws

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AWSSAMCLI() schema.Executable {
	return schema.Executable{
		Name:    "AWS SAM CLI",
		Runs:    []string{"sam"},
		DocsURL: sdk.URL("https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/reference-sam-cli.html"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.AccessKey,
				Provisioner: CLIProvisioner{},
			},
		},
	}
}
