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
			// Each entry below intentionally bypasses provisioning. Subcommands not listed here
			// (e.g. `aws configure import`, `aws configure export-credentials`) still flow through
			// the provisioner because they exchange or mutate credentials and benefit from the
			// 1Password-managed surface.
			//
			// `aws sso login`/`logout` write/erase the very token the SSO Profile provisioner reads
			// from disk; provisioning before they run creates a chicken-and-egg failure.
			needsauth.NotWhenContainsArgs("sso", "login"),
			needsauth.NotWhenContainsArgs("sso", "logout"),
			// `aws configure sso` and `aws configure sso-session` set up SSO configuration in
			// ~/.aws/config; they touch no AWS APIs and need no provisioned credentials.
			needsauth.NotWhenContainsArgs("configure", "sso"),
			needsauth.NotWhenContainsArgs("configure", "sso-session"),
			// `aws configure list` and `list-profiles` are read-only diagnostic queries against
			// the local config; `get` and `set` mutate ~/.aws/config or ~/.aws/credentials but
			// make no remote API call. Provisioning would be a no-op at best and a confusing
			// "missing credentials" error at worst when the user is just inspecting state.
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
