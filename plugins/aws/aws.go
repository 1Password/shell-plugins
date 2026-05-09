package aws

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AWSCLI() schema.Executable {
	return schema.Executable{
		Name:    "AWS CLI",
		Runs:    []string{"aws"},
		DocsURL: sdk.URL("https://aws.amazon.com/cli/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			// `aws sso login` and `aws sso logout` write/erase the very token
			// the SSO Profile provisioner reads from disk; provisioning before
			// they run creates a chicken-and-egg failure.
			needsauth.NotWhenContainsArgs("sso", "login"),
			needsauth.NotWhenContainsArgs("sso", "logout"),
			// `aws configure sso` and `aws configure sso-session` set up the
			// SSO configuration in ~/.aws/config; no provisioned credentials needed.
			needsauth.NotWhenContainsArgs("configure", "sso"),
			needsauth.NotWhenContainsArgs("configure", "sso-session"),
			// Read-only / diagnostic configure subcommands don't make API calls.
			needsauth.NotWhenContainsArgs("configure", "list"),
			needsauth.NotWhenContainsArgs("configure", "list-profiles"),
			needsauth.NotWhenContainsArgs("configure", "get"),
			needsauth.NotWhenContainsArgs("configure", "set"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.AccessKey,
				Provisioner: CLIProvisioner{},
			},
			{
				Name:        credname.SSOProfile,
				Provisioner: SSOCLIProvisioner{},
			},
		},
	}
}
