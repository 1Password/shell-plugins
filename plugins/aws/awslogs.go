package aws

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func awslogsCLI() schema.Executable {
	return schema.Executable{
		Name:    "awslogs",
		Runs:    []string{"awslogs"},
		DocsURL: sdk.URL("https://github.com/jorgebastida/awslogs"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		// Known limitation when both credentials are configured: awslogs treats AWS_PROFILE as an explicitly
		// selected profile and then ignores credentials in the environment. With `--profile <access key profile>`,
		// the SSO Profile provisioner yields and keeps the profile as AWS_PROFILE, so awslogs looks for credentials
		// in that profile instead of using the Access Key that was provisioned. Setups with a single credential are
		// not affected.
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.AccessKey,
				Provisioner: CLIProvisioner{},
				Optional:    true,
			},
			{
				Name:        credname.SSOProfile,
				Provisioner: SSOCLIProvisioner{},
				Optional:    true,
			},
		},
	}
}
