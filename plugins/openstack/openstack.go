package openstack

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OpenStackCLI() schema.Executable {
	return schema.Executable{
		Name:    "OpenStack CLI",
		Runs:    []string{"openstack"},
		DocsURL: sdk.URL("https://docs.openstack.org/python-openstackclient/latest/"),
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

func OSC() schema.Executable {
	return schema.Executable{
		Name:    "OSC",
		Runs:    []string{"osc"},
		DocsURL: sdk.URL("https://gtema.github.io/openstack/cli.html"),
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
